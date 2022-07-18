package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/webmakom-com/saiAuth/config"
	saiWebSocket "github.com/webmakom-com/saiAuth/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

type Server struct {
	Config    config.Configuration
	Host      string
	Port      string
	Websocket bool
}

type SocketMessage struct {
	Path string `json:"path"`
	Body []byte `json:"body"`
}

var ws saiWebSocket.Manager
var socketMessage SocketMessage
var clients = make(map[string]*websocket.Conn)
var broadcast = make(chan []byte)
var upgrader = websocket.Upgrader{}

func NewServer(c config.Configuration, w bool) Server {
	return Server{
		Config:    c,
		Host:      c.HttpServer.Host,
		Port:      c.HttpServer.Port,
		Websocket: w,
	}
}

func (s Server) SocketStart() {
	ln, _ := net.Listen("tcp", s.Config.SocketServer.Port)
	conn, _ := ln.Accept()

	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		json.Unmarshal([]byte(message), &socketMessage)
		s.handleSocketServerRequest(socketMessage)
	}
}

func (s Server) Start() {
	r := mux.NewRouter()
	http.Handle("/", r)

	if s.Websocket {
		r.HandleFunc("/ws", s.handleWSConnections)
		ws = saiWebSocket.NewWebSocketManager(s.Config)
	}

	r.HandleFunc("/{any}", s.handleHttpConnections)

	fmt.Println("Server has been started!")
	err := http.ListenAndServe(s.Host+":"+s.Port, nil)

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

func (s Server) handleHttpConnections(w http.ResponseWriter, r *http.Request) {
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

	s.handleHttpServerRequest(w, r)
}

func (s Server) WSProcess() {
	for {
		msg := <-broadcast
		s.handleWebSocketRequest(msg)

		for k, client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			time.Sleep(3 * time.Millisecond)

			if err != nil {
				_ = client.Close()
				delete(clients, k)
			}
		}
	}
}

func (s Server) handleWSConnections(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	wss, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
	}

	defer func(wss *websocket.Conn) {
		cErr := wss.Close()

		if cErr != nil {
			fmt.Println(cErr)
		}
	}(wss)

	clients[wss.RemoteAddr().String()] = wss

	for {
		_, msg, rErr := wss.ReadMessage()

		if rErr != nil {
			return
		}

		broadcast <- msg
	}
}
