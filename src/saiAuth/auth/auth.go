package auth

type Role string

type User struct {
	Login string `json:"login"`
	Pass  string `json:"pass"`
	Roles []Role
}

type Permission struct {
	Collection     string `json:"collection"`
	Method         string `json:"method"`
	RequiredParams string `json:"required_params"` // check request contains this or something multiple or wildcard
}
