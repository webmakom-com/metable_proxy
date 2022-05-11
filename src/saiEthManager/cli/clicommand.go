package cli

import (
	"github.com/thatisuday/commando"
	"github.com/webmakom-com/saiEthManager/config"
	"github.com/webmakom-com/saiEthManager/server"
)

func InitCli() {
	commando.
		SetExecutableName("sai-eth-manager").
		SetVersion("1.0.0")

	commando.
		Register("start").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			config.Load()
			server.NewServer().Start()
		})

	commando.Parse(nil)
}
