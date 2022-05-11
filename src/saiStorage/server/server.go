package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/webmakom-com/hv/src/saiStorage/config"
	"github.com/webmakom-com/hv/src/saiStorage/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

type Server struct {
	Config    config.Configuration
	Host      string
	Port      string
	Websocket bool
}

var ws websocket.Manager

func NewServer(c config.Configuration, w bool) Server {
	return Server{
		Config:    c,
		Host:      c.HttpServer.Host,
		Port:      c.HttpServer.Port,
		Websocket: w,
	}
}

func (s Server) Start() {
	r := mux.NewRouter()
	http.Handle("/", r)

	if s.Websocket {
		r.HandleFunc("/ws", s.handleWSConnections)
		ws = websocket.NewWebSocketManager(s.Config)
	}

	r.HandleFunc("/{any}", s.handleConnections)

	fmt.Println("Server has been started!")
	err := http.ListenAndServe(s.Host + ":" + s.Port, nil)

	if err != nil {
		fmt.Println("Server error: ", err)
	}
}

func (s Server) hasAccess(r *http.Request) bool {
	headers := r.Header
	token, ok := headers["Token"]

	if !ok {
		return false
	}

	if len(token) > 0 && token[0] == s.Config.Token {
		return true
	}

	return false
}

func (s Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := r.ParseForm()

	if err != nil {
		return
	}

	if !s.hasAccess(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		fmt.Println("Unauthorized access")
		message, _ := json.Marshal(bson.M{"message": "Unauthorized access"})
		_, _ = w.Write(message)
		return
	}

	s.handleServerRequest(w, r)
}
