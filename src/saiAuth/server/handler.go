package server

import (
	"encoding/json"
	"fmt"
)

type HandlerRequest struct {
	Token  string
	Method string
	Body   []byte
}

func (s Server) Register(h HandlerRequest) interface{} {
	return s.AuthManager.Register(h.getInterface(), h.Token)
}

func (s Server) Login(h HandlerRequest) interface{} {
	return s.AuthManager.Login(h.getInterface())
}

func (s Server) Access(h HandlerRequest) interface{} {
	return s.AuthManager.Access(h.getInterface(), h.Token)
}

func (s Server) Auth(h HandlerRequest) interface{} {
	return s.AuthManager.Auth(h.getInterface(), h.Token)
}

func (h HandlerRequest) getInterface() map[string]interface{} {
	var r = new(map[string]interface{})
	err := json.Unmarshal(h.Body, r)

	if err != nil {
		fmt.Println(err)
	}

	return *r
}
