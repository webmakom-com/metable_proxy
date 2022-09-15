package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/webmakom-com/saiAuth/config"
	"github.com/webmakom-com/saiAuth/utils/saiStorageUtil"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slices"
)

const (
	placeholder = "$"
)

// struct for replacePlaceholders func
type outerObject struct {
	Key   string
	Value string
}

type Manager struct {
	Config   config.Configuration
	Database saiStorageUtil.Database
}

type Token struct {
	Name        string
	Permissions []map[string]config.Permission
	Expiration  int64
}

type FindResult struct {
	Count int64                    `json:"count,omitempty"`
	Users []map[string]interface{} `json:"result,omitempty"`
}

type LoginResult struct {
	Token string `json:"token"`
	User  map[string]interface{}
}

func NewAuthManager(c config.Configuration) Manager {
	return Manager{
		Config:   c,
		Database: saiStorageUtil.Storage(c.Storage.Url, c.Storage.Auth.Email, c.Storage.Auth.Password),
	}
}

func (am Manager) Register(r map[string]interface{}, t string) interface{} {
	if !am.Access(r, t).(bool) {
		fmt.Println("Unauthorized request")
		return false
	}

	if am.isAuthRequestWrong(r) {
		fmt.Println("Wrong auth request")
		return false
	}

	if am.isUserExists(r) {
		fmt.Println("User exists")
		return false
	}

	if r["roles"] == nil || slices.Contains(r["roles"].([]string), "Admin") {
		r["roles"] = [1]string{"User"}
	}

	r["password"] = am.createPass(r["password"].(string))

	err, result := am.Database.Put("users", r, am.Config.Token)

	if err != nil {
		fmt.Println(err)
		return false
	}

	return string(result)
}

func (am Manager) Login(r map[string]interface{}) interface{} {
	if am.isAuthRequestWrong(r) {
		fmt.Println("Wrong auth request")
		return false
	}

	var (
		wrappedResult map[string]interface{}
		users         []map[string]interface{}
	)

	r["password"] = am.createPass(r["password"].(string))
	err, result := am.Database.Get("users", r, bson.M{}, am.Config.Token)

	if err != nil {
		fmt.Println(err)
		return false
	}

	jsonErr := json.Unmarshal(result, &wrappedResult)

	if jsonErr != nil {
		fmt.Println(string(result))
		fmt.Println(jsonErr)
		return false
	}

	usersMarshalled, err := json.Marshal(wrappedResult["result"])

	if err != nil {
		fmt.Println(string(usersMarshalled))
		fmt.Println(err)
		return false
	}

	jsonErr = json.Unmarshal(usersMarshalled, &users)

	if jsonErr != nil {
		fmt.Println(string(result))
		fmt.Println(jsonErr)
		return false
	}

	if len(users) == 0 {
		fmt.Println("Missing user")
		return false
	}

	if users[0]["roles"] == nil {
		fmt.Println("Missing roles")
		return false
	}

	roles := users[0]["roles"].([]interface{})

	for _, role := range roles {
		var perms []map[string]config.Permission
		roleName := role.(string)

		if am.Config.Roles[roleName].Exists {
			perms = append(perms, am.Config.Roles[roleName].Permissions)
		}

		t := am.createToken(perms, users[0])

		if t == nil {
			return false
		}

		return &LoginResult{
			Token: t.Name,
			User:  users[0],
		}
	}

	return false
}

func (am Manager) Access(r map[string]interface{}, t string) interface{} {
	if t == am.Config.Token {
		return true
	}

	if am.isAccessRequestWrong(r) {
		fmt.Println("Wrong access request")
		return false
	}

	err, result := am.Database.Get("tokens", bson.M{"Name": t}, bson.M{}, am.Config.Token)

	if err != nil {
		fmt.Println(err)
		return false
	}
	var (
		wrappedResult map[string]interface{}
		tokens        []Token
	)

	jsonErr := json.Unmarshal(result, &wrappedResult)

	if jsonErr != nil {
		fmt.Println(string(result))
		fmt.Println(jsonErr)
		return false
	}

	tokensMarshalled, err := json.Marshal(wrappedResult["result"])

	if err != nil {
		fmt.Println(string(tokensMarshalled))
		fmt.Println(err)
		return false
	}

	jsonErr = json.Unmarshal(tokensMarshalled, &tokens)

	if jsonErr != nil {
		fmt.Println(jsonErr)
		return false
	}

	if len(tokens) == 0 {
		fmt.Println("Unauthorized request")
		return false
	}

	for _, perms := range tokens[0].Permissions {
		if perms[r["collection"].(string)].Exists &&
			perms[r["collection"].(string)].Methods[r["method"].(string)] {
			return true
		}
	}

	return false
}

func (am Manager) createPass(pass string) string {
	hasher := sha256.New()
	hasher.Write([]byte(pass))
	hasher.Write([]byte(am.Config.Salt))

	return hex.EncodeToString(hasher.Sum(nil))
}

func (am Manager) replacePlaceholders(permissions []map[string]config.Permission, object map[string]interface{}) []map[string]config.Permission {
	var outerObjects []*outerObject
	permBytes, err := json.Marshal(permissions)

	if err != nil {
		return permissions
	}

	keys := findPlaceholders(permBytes)

	for k, v := range object {
		for _, key := range keys {
			if k == key {
				obj := &outerObject{
					Key:   key,
					Value: v.(string), // suppose that values to replace placeholders of string type
				}
				outerObjects = append(outerObjects, obj)
			}
		}
	}

	for _, permMap := range permissions {
		for _, permission := range permMap {
			for reqKey, reqValue := range permission.Required {
				for _, object := range outerObjects {
					if reqKey == object.Key {
						if reqValue == placeholder {
							permission.Required[reqKey] = object.Value
						}
					}
				}
			}
		}

	}
	return permissions
}

func (am Manager) createToken(permissions []map[string]config.Permission, object map[string]interface{}) *Token {
	var t = new(Token)

	hasher := sha256.New()
	hasher.Write(uuid.New().NodeID())
	hasher.Write([]byte(time.Now().String()))
	t.Name = hex.EncodeToString(hasher.Sum(nil))
	t.Permissions = am.replacePlaceholders(permissions, object)
	t.Expiration = time.Now().Unix() + 3600

	tokenErr, _ := am.Database.Put("tokens", t, am.Config.Token)

	if tokenErr != nil {
		fmt.Println(tokenErr)
		return nil
	}

	return t
}

func (am Manager) isAuthRequestWrong(r map[string]interface{}) bool {
	return r["password"] == nil
}

func (am Manager) isAccessRequestWrong(r map[string]interface{}) bool {
	return r["collection"] == nil || r["method"] == nil
}

func (am Manager) isUserExists(r map[string]interface{}) bool {
	err, result := am.Database.Get("auth", bson.M{"name": r["name"]}, bson.M{}, am.Config.Token)
	if err != nil {
		fmt.Println(err)
		return true
	}

	var wrappedResult map[string]interface{}

	jsonErr := json.Unmarshal(result, &wrappedResult)
	if jsonErr != nil {
		fmt.Println(string(result))
		fmt.Println(jsonErr)
		return false
	}

	resultMarshalled, err := json.Marshal(wrappedResult["result"])

	if err != nil {
		fmt.Println(string(resultMarshalled))
		fmt.Println(err)
		return false
	}

	return string(resultMarshalled) != "null"
}

func findPlaceholders(permBytes []byte) (keys []string) {
	slice := strings.Split(string(permBytes), ",")
	for _, s1 := range slice {
		if strings.Contains(s1, placeholder) {
			s2 := strings.Split(s1, ":")
			if len(s2) == 2 {
				uncroppedKey := s2[0]
				key := strings.Trim(uncroppedKey, "\"")
				keys = append(keys, key)

			} else {
				uncroppedKey := s2[len(s2)-2]
				key := strings.Trim(uncroppedKey, "{")
				key = strings.Trim(key, "\"")
				keys = append(keys, key)
			}
		}

	}
	return keys
}
