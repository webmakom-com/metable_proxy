package server

import (
	"fmt"
	"net/http"

	"github.com/webmakom-com/saiGNMonitor/config"
	"github.com/webmakom-com/saiGNMonitor/service"
)

type Server struct {
	Url string
}

var urlName = "%s:%d"

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func NewServer() *Server {
	return &Server{
		fmt.Sprintf(urlName, config.Get().HttpServer.Host, config.Get().HttpServer.Port),
	}
}

func (s *Server) Start() {
	sce := service.NewService()
	go sce.Start()

	http.HandleFunc("/api", s.api)
	fmt.Println("Server has been started: http://" + s.Url)
	serverErr := http.ListenAndServe(s.Url, nil)
	if serverErr != nil {
		panic(serverErr)
	}
}

func (s *Server) api(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
}
