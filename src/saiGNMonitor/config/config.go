package config

import (
	"fmt"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	HttpServer struct {
		Host string
		Port int
	}
	GlassNode struct {
		Url      string
		Key      string
		Period   int
		Interval int
		Retries  int
	}
}

var config Configuration

func Load() {
	configErr := gonfig.GetConf("config.json", &config)

	if configErr != nil {
		fmt.Println("Config load error: ", configErr)
		panic(configErr)
	}
}

func Get() Configuration {
	return config
}
