package main

import (
	"log"

	"github.com/webmakom-com/saiAuth/config"
	"github.com/webmakom-com/saiAuth/server"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	logger, err := zap.NewDevelopment(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		log.Fatal("error creating logger : ", err.Error())
	}
	logger.Debug("Logger started", zap.String("mode", "debug"))
	srv := server.NewServer(cfg, false, logger)

	if cfg.SocketServer.Host != "" {
		go srv.SocketStart()
	}

	//srv.StartHttps()
	srv.Start()
}
