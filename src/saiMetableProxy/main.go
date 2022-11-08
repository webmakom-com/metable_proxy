package main

import (
	"github.com/saiset-co/saiMetableProxy/internal"
	"github.com/saiset-co/saiService"
)

func main() {
	internal.Service.GlobalService = saiService.NewService("saiMetableProxy")

	internal.Service.GlobalService.RegisterConfig("config.yml")

	internal.Service.GlobalService.RegisterHandlers(internal.Service.Handler)

	internal.Init()

	internal.Service.GlobalService.RegisterInitTask(internal.Service.Init)

	internal.Service.GlobalService.RegisterTasks([]func(){
		internal.Service.Process,
	})

	internal.Service.GlobalService.Start()
}
