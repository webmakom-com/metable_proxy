package utils

import (
	"github.com/onrik/ethrpc"
	"github.com/webmakom-com/saiContractExplorer/config"
	"log"
)

func NewGethClient(c config.Configuration) (*ethrpc.EthRPC, error) {
	if len(c.Geth.Socket.Addresses) < 0 {
		panic("Geth configuration missed!")
	}

	client := ethrpc.New(c.Geth.Socket.Addresses[0])

	_, clientErr := client.Web3ClientVersion()

	if clientErr != nil {
		var iclientErr error

		for i := 0; i < len(c.Geth.Socket.Addresses); i++ {
			client = ethrpc.New(c.Geth.Socket.Addresses[i])
			_, iclientErr = client.Web3ClientVersion()

			if iclientErr == nil {
				return client, nil
			}
		}

		if iclientErr != nil {
			log.Println("Geth client problem:", iclientErr)
			return client, clientErr
		}
	}

	return client, nil
}
