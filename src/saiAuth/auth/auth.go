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
)

type AuthManager struct {
	Config   config.Configuration
	Database saiStorageUtil.Database
}

type Token struct {
	Name        string
	Permissions any
	Expiration  int64
}

func NewAuthManager(c config.Configuration) AuthManager {
	return AuthManager{
		Config:   c,
		Database: saiStorageUtil.Storage(c.Storage.Url, c.Storage.Auth.Email, c.Storage.Auth.Password),
	}
}

func (am AuthManager) Register(r map[string]interface{}) interface{} {
	if am.isAuthRequestWrong(r) {
		fmt.Println("Wrong auth request")
		return false
	}

	if am.isUserExists(r) {
		fmt.Println("User exists")
		return false
	}

	if r["role"] == nil || r["role"] == "Admin" {
		r["role"] = "User"
	}

	r["password"] = am.createPass(r["password"].(string))

	err, result := am.Database.Put("users", r, am.Config.Token)

	if err != nil {
		fmt.Println(err)
		return false
	}

	return string(result)
}

func (am AuthManager) Login(r map[string]interface{}) interface{} {
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

	if users[0]["role"] == nil {
		fmt.Println("Missing role")
		return false
	}

	if am.Config.Roles[users[0]["role"].(string)].Id > 0 {
		t := am.createToken(am.Config.Roles[users[0]["role"].(string)].Permissions)

		if t == nil {
			return false
		}

		return t.Name
	}

	return false
}

func (am AuthManager) Access(r any) interface{} {
	return false
}

func (am AuthManager) createPass(pass string) string {
	hasher := sha256.New()
	hasher.Write([]byte(pass))
	hasher.Write([]byte(am.Config.Salt))

	return hex.EncodeToString(hasher.Sum(nil))
}

func (am AuthManager) createToken(permissions interface{}) *Token {
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

func (am AuthManager) isAuthRequestWrong(r map[string]interface{}) bool {
	return r["name"] == nil || r["password"] == nil
}

func (am AuthManager) isUserExists(r map[string]interface{}) bool {
	err, result := am.Database.Get("auth", bson.M{"name": r["name"]}, bson.M{}, am.Config.Token)

	if err != nil {
		fmt.Println(err)
		return true
	}

	return string(result) != "null"
}
