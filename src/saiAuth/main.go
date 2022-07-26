package main

import (
	"github.com/webmakom-com/saiAuth/config"
	"github.com/webmakom-com/saiAuth/server"
)

func main() {
	cfg := config.Load()
	srv := server.NewServer(cfg, false)

	if cfg.SocketServer.Host != "" {
		go srv.SocketStart()
	}

	srv.StartHttps()
	//srv.Start()
}
