package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/webmakom-com/saiAuth/auth"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/webmakom-com/saiAuth/config"
	"github.com/webmakom-com/saiAuth/utils"
	saiWebSocket "github.com/webmakom-com/saiAuth/websocket"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Server struct {
	Config      config.Configuration
	Websocket   bool
	AuthManager auth.Manager
}

type SocketMessage struct {
	Path  string `json:"path"`
	Body  []byte `json:"body"`
	Token string `json:"token"`
}

var ws saiWebSocket.Manager
var socketMessage SocketMessage
var clients = make(map[string]*websocket.Conn)
var broadcast = make(chan []byte)
var upgrader = websocket.Upgrader{}

func NewServer(c config.Configuration, w bool) Server {
	return Server{
		Config:      c,
		Websocket:   w,
		AuthManager: auth.NewAuthManager(c),
	}
}

func (s Server) SocketStart() {
	ln, nErr := net.Listen("tcp", s.Config.SocketServer.Host+":"+s.Config.SocketServer.Port)

	if nErr != nil {
		fmt.Println(nErr)
	}

	conn, cErr := ln.Accept()

	if nErr != nil {
		fmt.Println(cErr)
	}

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
		r.HandleFunc("/ws/{any}", s.handleWSConnections)
		ws = saiWebSocket.NewWebSocketManager(s.Config)
	}

	r.HandleFunc("/{any}", s.handleHttpConnections)

	fmt.Println("Server has been started!")
	httpErr := http.ListenAndServe(s.Config.HttpServer.Host+":"+s.Config.HttpServer.Port, nil)

	if httpErr != nil {
		fmt.Println("Htpp server error: ", httpErr)
	}
}

func (s Server) StartHttps() {
	r := mux.NewRouter()
	http.Handle("/", r)

	if s.Websocket {
		r.HandleFunc("/ws/{any}", s.handleWSConnections)
		ws = saiWebSocket.NewWebSocketManager(s.Config)
	}

	r.HandleFunc("/{any}", s.handleHttpConnections)

	fmt.Println("Htpps server has been started!")

	httpsErr := http.ListenAndServeTLS(s.Config.HttpsServer.Host+":"+s.Config.HttpsServer.Port, "server.crt", "server.key", nil)

	if httpsErr != nil {
		fmt.Println("Server error: ", httpsErr)
	}
}

func (s Server) handleHttpConnections(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := r.ParseForm()

	if err != nil {
		fmt.Println(err)
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

func (s Server) handleWebSocketRequest(msg []byte) {
	handlerMessage := new(HandlerRequest)
	err := json.Unmarshal(msg, handlerMessage)

	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	s.handleServerRequest(*handlerMessage)
}

func (s Server) handleSocketServerRequest(msg SocketMessage) {
	handlerMessage := HandlerRequest{
		Method: msg.Path,
		Body:   msg.Body,
		Token:  msg.Token,
	}

	s.handleServerRequest(handlerMessage)
}

func (s Server) handleHttpServerRequest(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	headers := r.Header
	token, ok := headers["Token"]

	if !ok {
		return
	}

	handlerMessage := HandlerRequest{
		Method: strings.Trim(r.URL.Path, "/"),
		Body:   bytes,
		Token:  token[0],
	}

	result := s.handleServerRequest(handlerMessage)
	_, writeErr := w.Write(utils.ConvertInterfaceToJson(result))

	if writeErr != nil {
		fmt.Println("Write error:", writeErr)
		return
	}
}

func (s Server) handleServerRequest(h HandlerRequest) interface{} {
	val := reflect.ValueOf(s)
	in := []reflect.Value{reflect.ValueOf(h)}

	if !val.IsValid() {
		return "Reflect failed"
	}

	stringUtil := cases.Title(language.Und)
	method := val.MethodByName(stringUtil.String(h.Method))

	if !method.IsValid() {
		return "Method not found"
	}

	r := method.Call(in)

	return r[0].Interface()
}
