package utils

import (
	"fmt"

	"github.com/onrik/ethrpc"
	"github.com/webmakom-com/saiContractExplorer/config"
)

func NewGethClient(c config.Configuration) (*ethrpc.EthRPC, error) {
	if len(c.Geth) < 0 {
		panic("Geth configuration missed!")
	}

	client := ethrpc.New(c.Geth[0])

	_, clientErr := client.Web3ClientVersion()

	if clientErr != nil {
		var iclientErr error

		for i := 0; i < len(c.Geth); i++ {
			client = ethrpc.New(c.Geth[i])
			_, iclientErr = client.Web3ClientVersion()

			if iclientErr == nil {
				return client, nil
			}
		}

		if iclientErr != nil {
			fmt.Println("Geth client problem!! ")
			return client, clientErr
		}
	}

	return client, nil
}
