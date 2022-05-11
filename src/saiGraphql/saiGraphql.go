package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tkanos/gonfig"
	"log"
	"net/http"
	"time"
)

type configuration struct {
	Query      string
	Graph      string
	Log        string
	Sleep      int64
	HttpServer struct {
		Host string
		Port string
	}
	Storage struct {
		Url  string
		Auth struct {
			Email    string
			Password string
		}
	}
}

type GraphAnswerType struct {
	Data struct {
		Pairs []PairType `json:"pairs"`
	}
}

type PairType struct {
	Id         string `json:"id"`
	Reserve0   string `json:"reserve0"`
	Reserve1   string `json:"reserve1"`
	ReserveUSD string `json:"reserveUSD"`
	Token0     struct {
		Decimals string `json:"decimals"`
		Id       string `json:"id"`
		Name     string `json:"name"`
		Symbol   string `json:"symbol"`
	}
	Token0Price string `json:"token0Price"`
	Token1      struct {
		Decimals string `json:"decimals"`
		Id       string `json:"id"`
		Name     string `json:"name"`
		Symbol   string `json:"symbol"`
	}
	Token1Price       string `json:"token1Price"`
	TrackedReserveETH string `json:"trackedReserveETH"`
	TxCount           string `json:"txCount"`
	VolumeUSD         string `json:"volumeUSD"`
}

type MessageDataType struct {
	Pair    string `json:"pair"`
	Market  string `json:"market"`
	Amounts map[string]string `json:"amounts"`
}

type MessageType struct {
	Type string `json:"type"`
	Data []MessageDataType `json:"data"`
}

var config configuration
var answer GraphAnswerType

func main() {
	configErr := gonfig.GetConf("config.json", &config)

	if configErr != nil {
		fmt.Println("Config missed!! ")
		panic(configErr)
	}

	go process()

	fmt.Println("Server start: http://" + config.HttpServer.Host + ":" + config.HttpServer.Port)
	http.HandleFunc("/", api)
	http.HandleFunc("/ws", wss)

	serverErr := http.ListenAndServe(config.HttpServer.Host+":"+config.HttpServer.Port, nil)

	if serverErr != nil {
		fmt.Println("Server error: ", serverErr)
	}
}

func formatMessage(answer GraphAnswerType) MessageType {
	var messageData []MessageDataType

	for _, pair := range answer.Data.Pairs {
		var amounts = map[string]string{
			pair.Token0.Symbol: pair.Reserve0,
			pair.Token1.Symbol: pair.Reserve1,
		}

		data := MessageDataType {
			pair.Token0.Symbol + "/" + pair.Token1.Symbol,
			"uniswap",
			amounts,
		}

		messageData = append(messageData, data)
	}

	return MessageType {
		"LIQUIDITY_POOL_AMOUNTS",
		messageData,
	}
}

func process() {
	for {
		response, err := http.Post(config.Graph, "application/json", bytes.NewBuffer([]byte(config.Query)))

		if err != nil {
			log.Fatal(err)
		}

		json.NewDecoder(response.Body).Decode(&answer)
		message := formatMessage(answer)
		jsonString, _ := json.Marshal(message)

		fmt.Println("Message: ", string(jsonString))

		if config.Log != "" {
			//WebSocket send
			//err, _ = saiStorageUtil.Storage(config.Storage.Url, config.Storage.Auth.Email, config.Storage.Auth.Password).Put("log", message)
			//
			//if err != nil {
			//	log.Println(err)
			//}
		}

		time.Sleep(time.Second * time.Duration(config.Sleep))
	}
}
