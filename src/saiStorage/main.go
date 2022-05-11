package main

import (
	"github.com/webmakom-com/hv/src/saiStorage/config"
	"github.com/webmakom-com/hv/src/saiStorage/mongo"
	"github.com/webmakom-com/hv/src/saiStorage/server"
)

func main()  {
	cfg := config.Load()
	srv := server.NewServer(cfg, false)
	mSrv := mongo.NewMongoServer(cfg)

	go mSrv.Start()

	srv.Start()
}
