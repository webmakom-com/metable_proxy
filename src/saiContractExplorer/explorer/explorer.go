package explorer

import (
	"fmt"
	"github.com/webmakom-com/hv/src/saiContractExplorer/block"
	"github.com/webmakom-com/hv/src/saiContractExplorer/config"
	"github.com/webmakom-com/hv/src/saiContractExplorer/utils"
	"time"
)

type Explorer struct {
	Config config.Configuration
}

func NewExplorer(c config.Configuration) Explorer {
	return Explorer{
		Config: c,
	}
}

func (e Explorer) Process()  {
	client, err := utils.NewGethClient(e.Config)
	blockManager := block.NewBlockManager(e.Config)

	if err != nil {
		panic(err)
	}

	for {
		bid, err := client.EthBlockNumber()

		if err != nil {
			fmt.Println("Can't get last block:", err)
			continue
		}

		blk, err := blockManager.GetLastBlock(bid)

		if err != nil {
			fmt.Println("Can't get last block:", err)
			continue
		}

		for i := blk.Id; i <= bid; i++ {
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

			blockManager.HandleTransactions(blkInfo.Transactions)
		}

		blk.Id = bid
		blockManager.SetLastBlock(blk)
		time.Sleep(time.Duration(e.Config.Sleep) * time.Second)
	}
}