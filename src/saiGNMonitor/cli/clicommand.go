package cli

import (
	"github.com/thatisuday/commando"
	"github.com/webmakom-com/saiGNMonitor/config"
	"github.com/webmakom-com/saiGNMonitor/server"
)

func InitCli() {
	commando.
		SetExecutableName("sai-gn-monitor").
		SetVersion("1.0.0")

	commando.
		Register("start").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			config.Load()
			server.NewServer().Start()
		})

	commando.Parse(nil)
}
