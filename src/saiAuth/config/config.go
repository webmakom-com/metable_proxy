package config

import (
	"fmt"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	HttpServer struct {
		Host string
		Port string
	}
	HttpsServer struct {
		Host string
		Port string
	}
	Address struct {
		Url string
	}
	SocketServer struct {
		Host string
		Port string
	}
	Salt    string
	Token   string
	Storage struct {
		Token string
		Url   string
		Auth  struct {
			Email    string
			Password string
		}
	}
	Operations []string
	StartBlock int
	WebSocket  struct {
		Token string
		Url   string
	}
	Contract struct {
		Address string
		ABI     string
	}
	Geth  []string
	Sleep int
	Roles map[string]struct {
		Id          int
		Permissions map[string]struct {
			Methods  map[string]bool
			Required map[string]any
		}
	}
}

func Load() Configuration {
	var config Configuration
	err := gonfig.GetConf("config.json", &config)

	if err != nil {
		fmt.Println("Configuration problem:", err)
		panic(err)
	}

	return config
}
