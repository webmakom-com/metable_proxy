package mongo

import (
	"fmt"
	"github.com/slayer/autorestart"
	"github.com/webmakom-com/hv/src/saiStorage/config"
	"os/exec"
)

type Server struct {
	Config config.Configuration
}

func NewMongoServer(c config.Configuration) Server {
	return Server{
		Config: c,
	}
}

func (m Server) Start()  {
	autorestart.WatchFilename = "/usr/bin/mongod"
	autorestart.StartWatcher()

	startMongoCmd := exec.Command("/usr/bin/mongod")
	err := startMongoCmd.Start()

	if err != nil {
		fmt.Println("Mongo has been failed to start:", err)
		return
	}

	fmt.Println("Mongo has been started. PID:", startMongoCmd.Process.Pid)
}
