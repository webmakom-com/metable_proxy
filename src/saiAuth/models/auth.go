package models

import (
	"github.com/webmakom-com/saiAuth/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	AccessTokenType  = "access_token"
	RefreshTokenType = "refresh_token"
)

// Access token representation for unmarshal
type AccessToken struct {
	ID          string                         `json:"_id,omitempty"`
	Type        string                         `json:"type,omitempty"`
	Name        string                         `json:"name"`
	Expiration  int64                          `json:"expiration"`
	InternalID  string                         `json:"internal_id,omitempty"`
	User        map[string]interface{}         `json:"user,omitempty"`
	Permissions []map[string]config.Permission `json:"permissions,omitempty"`
}

// User representation inside access token
type User struct {
	ID         string `json:"_id,omitempty"`
	InternalID string `json:"internal_id,omitempty"`
}

// Refresh token representation
type RefreshToken struct {
	ID          string                         `json:"_id,omitempty"`
	Type        string                         `json:"type,omitempty"`
	Name        string                         `json:"name"`
	Expiration  int64                          `json:"expiration"`
	InternalID  string                         `json:"internal_id,omitempty"`
	AccessToken *AccessToken                   `json:"access_token,omitempty"`
	Permissions []map[string]config.Permission `json:"permissions,omitempty"`
}

// Response after login method
type LoginResponse struct {
	*AccessToken  `json:"at"`
	*RefreshToken `json:"rt"`
	User          map[string]interface{} `json:"user,omitempty"`
}

type AccessTokenWithObjectID struct {
	ID          primitive.ObjectID             `json:"_id,omitempty"`
	Type        string                         `json:"type"`
	Name        string                         `json:"name"`
	Expiration  int64                          `json:"expiration"`
	InternalID  string                         `json:"internal_id,omitempty"`
	User        *User                          `json:"user,omitempty"`
	Permissions []map[string]config.Permission `json:"permissions,omitempty"`
}
