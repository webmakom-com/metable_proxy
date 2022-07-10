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
	Address struct {
		Url string
	}
	Storage struct {
		Token string
		Url   string
		Auth  struct {
			Email    string
			Password string
		}
	}
	StartBlock int
	WebSocket  struct {
		Token string
		Url   string
	}
	Contracts []struct {
		Data struct {
			Address string
			ABI     string
		}
		Operations []string
	}
	Geth  []string
	Sleep int
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
