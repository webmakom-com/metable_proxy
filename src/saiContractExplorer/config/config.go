package config

import (
	"log"

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
	Geth struct {
		Web struct {
			Enabled   bool
			Addresses []string
		}
		Socket struct {
			Enabled   bool
			Addresses []string
		}
	}
	Sleep int
}

func Load() Configuration {
	var config Configuration
	err := gonfig.GetConf("config.json", &config)

	if err != nil {
		log.Println("Configuration problem:", err)
		panic(err)
	}

	return config
}
