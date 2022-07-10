package main

import (
	"github.com/webmakom-com/saiContractExplorer/config"
	"github.com/webmakom-com/saiContractExplorer/explorer"
	"github.com/webmakom-com/saiContractExplorer/server"
)

func main() {
	cfg := config.Load()
	srv := server.NewServer(cfg, true)
	exp := explorer.NewExplorer(cfg)

	go srv.WSProcess()
	go exp.Process()

	srv.Start()
}
