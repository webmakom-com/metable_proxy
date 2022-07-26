package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/webmakom-com/saiAuth/config"
	"github.com/webmakom-com/saiAuth/utils/saiStorageUtil"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slices"
)

type Manager struct {
	Config   config.Configuration
	Database saiStorageUtil.Database
}

type Token struct {
	Name        string
	Permissions []map[string]config.Permission
	Expiration  int64
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

	var users []map[string]interface{}

	r["password"] = am.createPass(r["password"].(string))
	err, result := am.Database.Get("users", r, bson.M{}, am.Config.Token)

	if err != nil {
		fmt.Println(err)
		return false
	}

	jsonErr := json.Unmarshal(result, &users)

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

	for _, role := range users[0]["roles"].([]string) {
		var perms []map[string]config.Permission

		if am.Config.Roles[role].Exists {
			perms = append(perms, am.Config.Roles[role].Permissions)
		}

		t := am.createToken(perms)

		if t == nil {
			return false
		}

		return t.Name
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

	var tokens []Token

	jsonErr := json.Unmarshal(result, &tokens)

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

func (am Manager) createToken(permissions []map[string]config.Permission) *Token {
	var t = new(Token)

	hasher := sha256.New()
	hasher.Write(uuid.New().NodeID())
	hasher.Write([]byte(time.Now().String()))
	t.Name = hex.EncodeToString(hasher.Sum(nil))
	t.Permissions = permissions
	t.Expiration = time.Now().Unix() + 3600

	tokenErr, _ := am.Database.Put("tokens", t, am.Config.Token)

	if tokenErr != nil {
		fmt.Println(tokenErr)
		return nil
	}

	return t
}

func (am Manager) isAuthRequestWrong(r map[string]interface{}) bool {
	return r["name"] == nil || r["password"] == nil
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

	return string(result) != "null"
}
