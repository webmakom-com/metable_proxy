package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/webmakom-com/saiStorage/mongo"
	"github.com/webmakom-com/saiStorage/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (s Server) handleWebSocketRequest(msg []byte) {

}

type jsonRequestType struct {
	Collection string        `json:"collection"`
	Select     bson.M        `json:"select"`
	Options    mongo.Options `json:"options"`
	Data       bson.M        `json:"data"`
}

func (s Server) handleServerRequest(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/get":
		{
			s.get(w, r, "get")
		}
	case "/save":
		{
			s.save(w, r, "save")
		}
	case "/update":
		{
			s.update(w, r, "update")
		}
	case "/upsert":
		{
			s.upsert(w, r, "upsert")
		}
	case "/remove":
		{
			s.remove(w, r, "remove")
		}
	}
}

func (s Server) get(w http.ResponseWriter, r *http.Request, method string) {
	var request jsonRequestType
	decoder := json.NewDecoder(r.Body)
	decoderErr := decoder.Decode(&request)

	if decoderErr != nil {
		fmt.Printf("Wrong JSON: %v", decoderErr)
		return
	}

	if s.Config.UsePermissionAuth {
		err := s.checkPermissionRequest(r, request.Collection, method)
		if err != nil {
			fmt.Println(err)
			w.Write([]byte(err.Error()))
			return
		}
	}

	result, mongoErr := s.Client.Find(request.Collection, request.Select, request.Options)

	if mongoErr != nil {
		fmt.Println("Mongo error:", mongoErr)
		return
	}

	_, writeErr := w.Write(utils.ConvertInterfaceToJson(result))

	if writeErr != nil {
		fmt.Println("Write error:", writeErr)
		return
	}
}

func (s Server) save(w http.ResponseWriter, r *http.Request, method string) {
	var request jsonRequestType
	decoder := json.NewDecoder(r.Body)
	decoderErr := decoder.Decode(&request)

	if decoderErr != nil {
		fmt.Printf("Wrong JSON: %v", decoderErr)
		return
	}

	if s.Config.UsePermissionAuth {
		err := s.checkPermissionRequest(r, request.Collection, method)
		if err != nil {
			fmt.Println(err)
			w.Write([]byte(err.Error()))
			return
		}
	}

	uuid := uuid.New()
	request.Data["internal_id"] = uuid.String()

	mongoErr := s.Client.Insert(request.Collection, request.Data)

	if mongoErr != nil {
		fmt.Println("Mongo error:", mongoErr)
		return
	}

	_, writeErr := w.Write(utils.ConvertInterfaceToJson(bson.M{"Status": "Ok"}))

	if writeErr != nil {
		fmt.Println("Write error:", writeErr)
		return
	}
}

func (s Server) update(w http.ResponseWriter, r *http.Request, method string) {
	var request jsonRequestType
	decoder := json.NewDecoder(r.Body)
	decoderErr := decoder.Decode(&request)

	if decoderErr != nil {
		fmt.Printf("Wrong JSON: %v", decoderErr)
		return
	}

	if s.Config.UsePermissionAuth {
		err := s.checkPermissionRequest(r, request.Collection, method)
		if err != nil {
			fmt.Println(err)
			w.Write([]byte(err.Error()))
			return
		}
	}

	mongoErr := s.Client.Update(request.Collection, request.Select, bson.M{"$set": request.Data})

	if mongoErr != nil {
		fmt.Println("Mongo error:", mongoErr)
		return
	}

	_, writeErr := w.Write(utils.ConvertInterfaceToJson(bson.M{"Status": "Ok"}))

	if writeErr != nil {
		fmt.Println("Write error:", writeErr)
		return
	}
}

func (s Server) upsert(w http.ResponseWriter, r *http.Request, method string) {
	var request jsonRequestType
	decoder := json.NewDecoder(r.Body)
	decoderErr := decoder.Decode(&request)

	if decoderErr != nil {
		fmt.Printf("Wrong JSON: %v", decoderErr)
		return
	}

	if s.Config.UsePermissionAuth {
		err := s.checkPermissionRequest(r, request.Collection, method)
		if err != nil {
			fmt.Println(err)
			w.Write([]byte(err.Error()))
			return
		}
	}

	mongoErr := s.Client.Upsert(request.Collection, request.Select, request.Data)

	if mongoErr != nil {
		fmt.Println("Mongo error:", mongoErr)
		return
	}

	_, writeErr := w.Write(utils.ConvertInterfaceToJson(bson.M{"Status": "Ok"}))

	if writeErr != nil {
		fmt.Println("Write error:", writeErr)
		return
	}
}

func (s Server) remove(w http.ResponseWriter, r *http.Request, method string) {
	var request jsonRequestType
	decoder := json.NewDecoder(r.Body)
	decoderErr := decoder.Decode(&request)

	if decoderErr != nil {
		fmt.Printf("Wrong JSON: %v", decoderErr)
		return
	}

	if s.Config.UsePermissionAuth {
		err := s.checkPermissionRequest(r, request.Collection, method)
		if err != nil {
			fmt.Println(err)
			w.Write([]byte(err.Error()))
			return
		}
	}

	mongoErr := s.Client.Remove(request.Collection, request.Select)

	if mongoErr != nil {
		fmt.Println("Mongo error:", mongoErr)
		return
	}

	_, writeErr := w.Write(utils.ConvertInterfaceToJson(bson.M{"Status": "Ok"}))

	if writeErr != nil {
		fmt.Println("Write error:", writeErr)
		return
	}
}
