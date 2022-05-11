package websocket

import (
	"bytes"
	"fmt"
	"github.com/webmakom-com/hv/src/saiContractExplorer/config"
	"net/http"
)

type Manager struct {
	Config config.Configuration
}

func NewWebSocketManager(c config.Configuration) Manager {
	return Manager{
		Config: c,
	}
}

func (w Manager) SendMessage(message string, token string) {
	url := w.Config.WebSocket.Url + "?method=broadcast&message=" + token + "|" + message
	req, err := http.NewRequest("GET", url, new(bytes.Buffer))

	if err != nil {
		fmt.Println("Websocket error:", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err = client.Do(req)

	if err != nil {
		fmt.Println("Websocket error:", err)
	}

	client.CloseIdleConnections()
}