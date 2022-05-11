package explorer

import (
	"fmt"
	"github.com/webmakom-com/hv/src/saiGasExplorer/config"
	"github.com/webmakom-com/hv/src/saiGasExplorer/utils"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
)

type Explorer struct {
	Config config.Configuration
}

func NewExplorer(c config.Configuration) Explorer {
	return Explorer{
		Config: c,
	}
}

func (e Explorer) Process(ids string, counts string) ([]interface{}, error) {
	var data []interface{}
	client, err := utils.NewGethClient(e.Config)
	id, _ := strconv.Atoi(ids)
	count, _ := strconv.Atoi(counts)
	bid := id + count

	if err != nil {
		fmt.Println("Geth error:", err)
		return data, err
	}

	for i := id; i <= bid; i++ {
		blkInfo, blockInfoErr := client.EthGetBlockByNumber(i, true)

		if blockInfoErr != nil {
			fmt.Println("Can't get block data!! ")
			i--
			continue
		}

		if len(blkInfo.Transactions) == 0 {
			fmt.Printf("Block %d from %d: No transactions found.\n", i, bid)
			continue
		}

		fmt.Printf("Block %d from %d: %d transactions found.\n", i, bid, len(blkInfo.Transactions))

		for j := 0; j < len(blkInfo.Transactions); j++ {
			item := bson.M{
				"Block":     i,
				"Hash":      blkInfo.Transactions[j].Hash,
				"From":      blkInfo.Transactions[j].From,
				"To":        blkInfo.Transactions[j].To,
				"Gas":    	 blkInfo.Transactions[j].Gas,
				"GasPrice":  blkInfo.Transactions[j].GasPrice,
			}

			data = append(data, item)
		}
	}

	return data, nil
}
