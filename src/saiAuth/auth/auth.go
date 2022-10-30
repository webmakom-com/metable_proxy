package auth

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/webmakom-com/saiAuth/config"
	"github.com/webmakom-com/saiAuth/models"
	"github.com/webmakom-com/saiAuth/utils/saiStorageUtil"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

const (
	placeholder = "$"
)

type Manager struct {
	Config   config.Configuration
	Database saiStorageUtil.Database
	Logger   *zap.Logger
}

type FindResult struct {
	Count int64                    `json:"count,omitempty"`
	Users []map[string]interface{} `json:"result,omitempty"`
}

type LoginResult struct {
	Token string `json:"token"`
	User  map[string]interface{}
}

type Selection struct {
	Field string
	Value interface{}
}

func NewAuthManager(c config.Configuration, logger *zap.Logger) Manager {
	return Manager{
		Config:   c,
		Database: saiStorageUtil.Storage(c.Storage.Url, c.Storage.Auth.Email, c.Storage.Auth.Password),
		Logger:   logger,
	}
}

func (am Manager) Register(r map[string]interface{}, t string) interface{} {
	if t != am.Config.Token {
		am.Logger.Error("Wrong auth request : wrong config token")
		return false
	}

	if am.isAuthRequestWrong(r) {
		am.Logger.Error("Wrong auth request : password field was not found")
		return false
	}

	r["password"] = am.createPass(r["password"].(string))

	if am.isUserExists(r) {
		am.Logger.Error("User already exists")
		return false
	}

	if r["roles"] == nil || slices.Contains(r["roles"].([]string), "Admin") {
		r["roles"] = [1]string{"User"}
	}

	err, result := am.Database.Put("users", r, am.Config.Token)

	if err != nil {
		am.Logger.Error("REGISTER - DB.PUT", zap.Error(err))
		return false
	}

	return string(result)
}

func (am Manager) Login(r map[string]interface{}, token string) interface{} {
	// handle login method with empty body
	if r == nil {
		am.Logger.Debug("GOT EMPTY BODY")
		return am.HandleRefreshToken(token)
	}

	if am.isAuthRequestWrong(r) {
		am.Logger.Error("Wrong auth request")
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
	fmt.Println(string(result))
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

	var perms []map[string]config.Permission

	for _, role := range roles {
		roleName := role.(string)

		if am.Config.Roles[roleName].Exists {
			rolePerm, mapErr := Map(am.Config.Roles[roleName].Permissions)

			if mapErr != nil {
				fmt.Println(mapErr)
				return false
			}

			perms = append(perms, rolePerm)
		}
	}

	t := am.createToken(perms, users[0])

	if t == nil {
		return false
	}

	delete(users[0], "password")

	return t
}

func (am Manager) Access(r map[string]interface{}, t string) interface{} {
	if t == am.Config.Token {
		return true
	}

	if am.isAccessRequestWrong(r) {
		fmt.Println("Wrong access request")
		return false
	}

	err, result := am.Database.Get("tokens", bson.M{"name": t}, bson.M{}, am.Config.Token)

	if err != nil {
		fmt.Println(err)
		return false
	}
	var (
		wrappedResult  map[string]interface{}
		tokens         []models.AccessToken
		emptySelection bool
	)

	selection := handleSelect(r)
	if selection == nil {
		emptySelection = true
	}

	fmt.Printf("got selection : %+v\n", selection) // DEBUG

	jsonErr := json.Unmarshal(result, &wrappedResult)

	if jsonErr != nil {
		fmt.Println(string(result))
		fmt.Println(jsonErr)
		return false
	}

	fmt.Println(string(result))

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

	token := tokens[0]

	fmt.Printf("got token : %+v\n", token) //DEBUG

	// for _, perm := range token.Permissions {
	// 	if perm[r["collection"].(string)].Exists &&
	// 		perm[r["collection"].(string)].Methods[r["method"].(string)] {
	// 		if emptySelection {
	// 			if perm[r["collection"].(string)].Required == nil {
	// 				return true
	// 			}
	// 		} else {
	// 			if perm[r["collection"].(string)].Required[selection.Field] == nil {
	// 				continue
	// 			} else {
	// 				if perm[r["collection"].(string)].Required[selection.Field] == selection.Value {
	// 					return true
	// 				}
	// 			}

	// 		}
	// 	}
	// }

	if emptySelection {
		for _, perms := range tokens[0].Permissions {
			if perms[r["collection"].(string)].Exists &&
				perms[r["collection"].(string)].Methods[r["method"].(string)] &&
				perms[r["collection"].(string)].Required == nil {
				return true
			}
		}
	} else {
		for _, perms := range tokens[0].Permissions {

			if perms[r["collection"].(string)].Required[selection.Field] == nil {
				continue
			}
			fmt.Printf("REQUIRED FIELD : %s\n", perms[r["collection"].(string)].Required[selection.Field].(string)) //DEBUG
			fmt.Printf("SELECTION FIELD : %s\n", selection.Value.(string))
			fmt.Println(perms[r["collection"].(string)].Required[selection.Field].(string) == selection.Value.(string)) //DEBUG

			if perms[r["collection"].(string)].Exists &&
				perms[r["collection"].(string)].Methods[r["method"].(string)] &&
				perms[r["collection"].(string)].Required[selection.Field].(string) == selection.Value.(string) {
				return true
			}
		}
	}

	return true
}

func (am Manager) createPass(pass string) string {
	hasher := sha256.New()
	hasher.Write([]byte(pass))
	hasher.Write([]byte(am.Config.Salt))

	return hex.EncodeToString(hasher.Sum(nil))
}

func (am Manager) replacePlaceholders(permissions []map[string]config.Permission, object map[string]interface{}) []map[string]config.Permission {
	for _, permMap := range permissions {
		for _, permission := range permMap {
			for reqKey, reqValue := range permission.Required {
				for k, v := range object {
					if reqKey == k && reqValue == placeholder {
						permission.Required[reqKey] = v
					}
				}
			}
		}
	}

	return permissions
}

func (am Manager) createToken(permissions []map[string]config.Permission, object map[string]interface{}) *models.LoginResponse {
	// access token creating
	at := &models.AccessToken{
		User: &models.User{},
	}
	at.Type = models.AccessTokenType

	hasher := sha256.New()
	hasher.Write(uuid.New().NodeID())
	hasher.Write([]byte(time.Now().String()))
	at.Name = hex.EncodeToString(hasher.Sum(nil))
	at.Permissions = am.replacePlaceholders(permissions, object)
	at.Expiration = time.Now().Unix() + am.Config.AccessTokenExp
	am.Logger.Sugar().Debugf("LOGIN - CREATE TOKEN - USER OBJECT INSIDE CREATE TOKEN FUNC : %+v\n", object)
	at.User.ID = object["_id"].(string)
	at.User.InternalID = object["internal_id"].(string)

	tokenErr, _ := am.Database.Put("tokens", at, am.Config.Token)

	if tokenErr != nil {
		fmt.Println(tokenErr)
		return nil
	}

	accessToken, err := am.getAccessToken(at.Name)
	if err != nil {
		return nil
	}
	at.ID = accessToken.ID

	am.Logger.Sugar().Debugf("AUTH - CREATE TOKEN - GET ACCESS TOKEN : [%+v\n]", accessToken)

	rt := &models.RefreshToken{
		AccessToken: at,
	}
	rt.Type = models.RefreshTokenType

	hasher = sha256.New()
	hasher.Write(uuid.New().NodeID())
	hasher.Write([]byte(time.Now().String()))
	rt.Name = hex.EncodeToString(hasher.Sum(nil))

	rtUpdatePerm := map[string]config.Permission{}
	rtUpdatePerm["tokens"] = config.Permission{
		Exists: true,
		Methods: map[string]bool{
			"update": true,
		},
		Required: map[string]any{
			"name": at.Name,
		},
	}

	rtSavePerm := map[string]config.Permission{}
	rtSavePerm["tokens"] = config.Permission{
		Exists: true,
		Methods: map[string]bool{
			"save":   true,
			"update": true,
		},
		Required: map[string]any{
			"type": models.RefreshTokenType,
		},
	}
	rt.Permissions = append(rt.Permissions, rtUpdatePerm, rtSavePerm)

	rt.Expiration = time.Now().Unix() + am.Config.RefreshTokenExp

	tokenErr, _ = am.Database.Put("tokens", rt, am.Config.Token)

	if tokenErr != nil {
		fmt.Println(tokenErr)
		return nil
	}

	return &models.LoginResponse{
		AccessToken: &models.AccessToken{
			Name:       at.Name,
			Expiration: at.Expiration,
		},
		RefreshToken: &models.RefreshToken{
			Name:       rt.Name,
			Expiration: rt.Expiration,
		},
	}
}

func (am Manager) isAuthRequestWrong(r map[string]interface{}) bool {
	return r["password"] == nil
}

func (am Manager) isAccessRequestWrong(r map[string]interface{}) bool {
	return r["collection"] == nil || r["method"] == nil
}

func (am Manager) isUserExists(r map[string]interface{}) bool {
	err, result := am.Database.Get("users", r, bson.M{}, am.Config.Token)

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

func Map(m map[string]config.Permission) (map[string]config.Permission, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err := enc.Encode(m)
	if err != nil {
		return nil, err
	}
	var copy map[string]config.Permission
	err = dec.Decode(&copy)
	if err != nil {
		return nil, err
	}
	return copy, nil
}

func handleSelect(r map[string]interface{}) *Selection {
	selection := &Selection{}

	fmt.Printf("got selection from handle select : %+v\n", r) //DEBUG
	s, ok := r["select"]
	if !ok {
		return nil
	} else {
		for k, v := range s.(map[string]interface{}) {
			selection.Field = k
			switch v.(type) {
			case string:
				selection.Value = v.(string)
			case []string:
				selection.Value = v.([]string)
			}
		}
	}
	return selection
}

func (am Manager) HandleRefreshToken(refreshToken string) bool {
	err, result := am.Database.Get("tokens", bson.M{"name": refreshToken}, bson.M{}, am.Config.Token)
	if err != nil {
		am.Logger.Error("LOGIN - HANDLE REFRESH TOKEN - DB.GET", zap.Error(err))
	}

	var (
		wrappedResult map[string]interface{}
		tokens        []models.RefreshToken
	)

	jsonErr := json.Unmarshal(result, &wrappedResult)

	if jsonErr != nil {
		am.Logger.Error("LOGIN - HANDLE REFRESH TOKEN - UNMARHSAL UNWRAPPED RESULT", zap.Error(err))
		return false
	}

	tokensMarshalled, err := json.Marshal(wrappedResult["result"])

	if err != nil {
		am.Logger.Error("LOGIN - HANDLE REFRESH TOKEN - MARHSAL TOKENS", zap.Error(err))
		return false
	}

	jsonErr = json.Unmarshal(tokensMarshalled, &tokens)

	if jsonErr != nil {
		am.Logger.Error("LOGIN - HANDLE REFRESH TOKEN - UNMARHSAL REFRESH TOKENS", zap.Error(err))
		return false
	}

	if len(tokens) == 0 {
		am.Logger.Error("LOGIN - HANDLE REFRESH TOKEN - RESULT LENGTH == 0")
		return false
	}

	token := tokens[0]

	if token.Type != models.RefreshTokenType {
		am.Logger.Error("LOGIN - HANDLE REFRESH TOKEN - INCORRECT TYPE OF TOKEN")
		return false
	}

	am.Logger.Sugar().Debugf("GOT REFRESH TOKEN : [%+v]", token)

	hasher := sha256.New()
	hasher.Write(uuid.New().NodeID())
	hasher.Write([]byte(time.Now().String()))
	accessTokenName := hex.EncodeToString(hasher.Sum(nil))
	accessTokenExpiration := time.Now().Unix() + am.Config.AccessTokenExp

	//update access token
	filter := bson.M{"name": token.AccessToken.Name}
	update := bson.M{"name": accessTokenName, "expiration": accessTokenExpiration}
	err, _ = am.Database.Update("tokens", filter, update, am.Config.Token)
	if err != nil {
		am.Logger.Error("LOGIN - HANDLE REFRESH TOKEN - UPDATE ACCESS TOKEN", zap.Error(err))
		return false
	}

	// !!!!!!!!!!!!!! create new refresh token

	token.AccessToken.Expiration = accessTokenExpiration
	token.AccessToken.Name = accessTokenName
	rt := &models.RefreshToken{
		AccessToken: token.AccessToken,
	}
	rt.Type = models.RefreshTokenType

	hasher = sha256.New()
	hasher.Write(uuid.New().NodeID())
	hasher.Write([]byte(time.Now().String()))
	rt.Name = hex.EncodeToString(hasher.Sum(nil))

	rtUpdatePerm := map[string]config.Permission{}
	rtUpdatePerm["tokens"] = config.Permission{
		Exists: true,
		Methods: map[string]bool{
			"update": true,
		},
		Required: map[string]any{
			"name": rt.AccessToken.Name,
		},
	}

	rtSavePerm := map[string]config.Permission{}
	rtSavePerm["tokens"] = config.Permission{
		Exists: true,
		Methods: map[string]bool{
			"save":   true,
			"update": true,
		},
		Required: map[string]any{
			"type": models.RefreshTokenType,
		},
	}
	rt.Permissions = append(rt.Permissions, rtUpdatePerm, rtSavePerm)

	rt.Expiration = time.Now().Unix() + am.Config.RefreshTokenExp

	tokenErr, _ := am.Database.Put("tokens", rt, am.Config.Token)

	if tokenErr != nil {
		am.Logger.Error("LOGIN - HANDLE REFRESH TOKEN - PUT NEW REFRESH TOKEN", zap.Error(err))
		return false
	}
	return true

}

func (am Manager) getAccessToken(t string) (*models.AccessToken, error) {
	err, result := am.Database.Get("tokens", bson.M{"name": t}, bson.M{}, am.Config.Token)

	if err != nil {
		return nil, err
	}
	var (
		wrappedResult map[string]interface{}
		tokens        []models.AccessToken
	)

	am.Logger.Sugar().Debugf("result from db : %s", string(result))

	jsonErr := json.Unmarshal(result, &wrappedResult)

	if jsonErr != nil {
		am.Logger.Error("AUTH - GET ACCESS TOKEN - UNMARSHAL WRAPPED RESULT", zap.Error(err))
		return nil, err
	}

	tokensMarshalled, err := json.Marshal(wrappedResult["result"])

	if err != nil {
		am.Logger.Error("AUTH - GET ACCESS TOKEN - MARSHAL RESULT", zap.Error(err))
		return nil, err
	}

	jsonErr = json.Unmarshal(tokensMarshalled, &tokens)

	if jsonErr != nil {
		am.Logger.Error("AUTH - GET ACCESS TOKEN - UNMARSHAL  RESULT", zap.Error(err))
		return nil, err
	}

	if len(tokens) == 0 {
		am.Logger.Error("AUTH - GET ACCESS TOKEN - UNMARSHAL WRAPPED RESULT - EMPTY RESULT")
		return nil, errors.New("empty result")
	}

	token := tokens[0]

	return &token, nil
}
