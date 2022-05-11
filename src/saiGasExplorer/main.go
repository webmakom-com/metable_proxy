package main

import (
	"github.com/webmakom-com/hv/src/saiGasExplorer/config"
	"github.com/webmakom-com/hv/src/saiGasExplorer/server"
)

func main()  {
	cfg := config.Load()
	srv := server.NewServer(cfg, true)

	srv.Start()
}
