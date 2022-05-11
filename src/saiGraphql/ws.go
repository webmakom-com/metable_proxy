package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var broadcast = make(chan string)
var upgrader = websocket.Upgrader{}

func wss(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
	}

	defer func(ws *websocket.Conn) {
		err := ws.Close()

		if err != nil {
			fmt.Println(err)
		}
	}(ws)

	for {
		msg := <-broadcast
		err := ws.WriteJSON(msg)

		if err != nil {
			return
		}

		fmt.Printf("%s sent: %s\n", ws.RemoteAddr(), string(msg))

	}
}
