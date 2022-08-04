package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var clients = make(map[string]*websocket.Conn)
var broadcast = make(chan []byte)
var upgrader = websocket.Upgrader{}

func (s Server) WSProcess() {
	for {
		msg := <-broadcast

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

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
	}

	defer func(ws *websocket.Conn) {
		cErr := ws.Close()

		if cErr != nil {
			log.Println(cErr)
		}
	}(ws)

	clients[ws.RemoteAddr().String()] = ws

	for {
		_, msg, rErr := ws.ReadMessage()

		if rErr != nil {
			return
		}

		broadcast <- msg
	}
}
