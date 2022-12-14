package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/webmakom-com/saiContractExplorer/config"
)

type Server struct {
	Host      string
	Port      string
	Websocket bool
}

func NewServer(c config.Configuration, w bool) Server {
	return Server{
		Host:      c.HttpServer.Host,
		Port:      c.HttpServer.Port,
		Websocket: w,
	}
}

func (s Server) Start() {
	http.HandleFunc("/", s.handleConnections)

	if s.Websocket {
		http.HandleFunc("/ws", s.handleWSConnections)
	}

	err := http.ListenAndServe(s.Host+":"+s.Port, nil)

	if err != nil {
		log.Println("Server error: ", err)
	}
}

func (s Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		return
	}

	method := strings.Join(r.Form["method"], "")
	switch method {
	default:
		{
			_, _ = w.Write([]byte("I'm alive"))
		}
	}
}
