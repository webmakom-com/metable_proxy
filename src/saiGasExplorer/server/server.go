package server

import (
	"fmt"
	"github.com/webmakom-com/hv/src/saiGasExplorer/config"
	"github.com/webmakom-com/hv/src/saiGasExplorer/explorer"
	"github.com/webmakom-com/hv/src/saiGasExplorer/utils/saiStorageUtil"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strings"
)

type Server struct {
	Host string
	Port string
	Websocket bool
	Token string
	Storage saiStorageUtil.Database
	Explorer explorer.Explorer
}

func NewServer(c config.Configuration, w bool) Server {
	return Server{
		Token: c.Storage.Token,
		Host: c.HttpServer.Host,
		Port: c.HttpServer.Port,
		Websocket: w,
		Storage: saiStorageUtil.Storage(c.Storage.Url, c.Storage.Auth.Email, c.Storage.Auth.Password),
		Explorer: explorer.NewExplorer(c),
	}
}

func (s Server) Start() {
	http.HandleFunc("/", s.handleConnections)

	if s.Websocket {
		http.HandleFunc("/ws", s.handleWSConnections)
	}

	err := http.ListenAndServe(s.Host + ":" + s.Port, nil)

	if err != nil {
		fmt.Println("Server error: ", err)
	}
}

func (s Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		return
	}

	block := strings.Join(r.Form["block"], "")
	count := strings.Join(r.Form["count"], "")

	switch r.URL.Path {
	case "/get":
		{
			res, getErr := s.get(block, count)

			if getErr != nil {
				b, _ := bson.Marshal(bson.M{"Status": "Fail"})
				_, _ = w.Write(b)
				return
			}

			b, _ := bson.Marshal(bson.M{"Status": "Ok", "Result": res})
			_, _ = w.Write(b)
		}
	case "/log":
		{
			s.log(block, count)
			b, _ := bson.Marshal(bson.M{"Status": "Ok"})
			_, _ = w.Write(b)
		}
	}
}

func (s Server) get(block string, count string) (interface{}, error) {
	res, explorerErr := s.Explorer.Process(block, count)

	if explorerErr != nil {
		fmt.Println("Explorer error:", explorerErr)
		return res, explorerErr
	}

	return res, nil
}

func (s Server) log(block string, count string) string {
	res, explorerErr := s.Explorer.Process(block, count)

	if explorerErr != nil {
		fmt.Println("Explorer error:", explorerErr)
		return "Fail"
	}

	storageErr, _ := s.Storage.Put("gas", res, s.Token)

	if storageErr != nil {
		fmt.Println("Database error:", storageErr)
		return "Fail"
	}

	return "Success"
}
