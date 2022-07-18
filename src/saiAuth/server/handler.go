package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/webmakom-com/saiAuth/utils"
)

type HandlerMessage struct {
	Path string
	Body []byte
}

func (s Server) handleWebSocketRequest(msg []byte) {
	handlerMessage := new(HandlerMessage)
	err := json.Unmarshal(msg, handlerMessage)

	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	s.handleServerRequest(*handlerMessage)
}

func (s Server) handleSocketServerRequest(msg SocketMessage) {
	handlerMessage := HandlerMessage{
		Path: msg.Path,
		Body: msg.Body,
	}

	s.handleServerRequest(handlerMessage)
}

func (s Server) handleHttpServerRequest(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	handlerMessage := HandlerMessage{
		Path: r.URL.Path,
		Body: bytes,
	}

	result := s.handleServerRequest(handlerMessage)
	_, writeErr := w.Write(utils.ConvertInterfaceToJson(result))

	if writeErr != nil {
		fmt.Println("Write error:", writeErr)
		return
	}
}

func (s Server) handleServerRequest(handlerMessage HandlerMessage) interface{} {
	return nil
}
