package explorer

import (
	"log"
	"time"

	"github.com/webmakom-com/saiContractExplorer/block"
	"github.com/webmakom-com/saiContractExplorer/config"
	"github.com/webmakom-com/saiContractExplorer/utils"
)

type Explorer struct {
	Config config.Configuration
}

func NewExplorer(c config.Configuration) Explorer {
	return Explorer{
		Config: c,
	}
}

func (e Explorer) WProcess() {
	blockManager := block.NewBlockManager(e.Config)

	for {
		bid, err := blockManager.EthBlockNumber()

		if err != nil {
			log.Println("Can't get last block:", err)
			continue
		}

		blk, err := blockManager.GetLastBlock(bid)

		if err != nil {
			log.Println("Can't get last block:", err)
			continue
		}

		for i := blk.Id; i <= bid; i++ {
			blkInfo, blockInfoErr := blockManager.EthGetBlockByNumber(i, true)

			if blockInfoErr != nil {
				log.Println("Can't get block data!! ")
				i--
				continue
			}

			if len(blkInfo.Transactions) == 0 {
				log.Printf("Block %d from %d: No transactions found.\n", i, bid)
				continue
			}

			log.Printf("Block %d from %d: %d transactions found.\n", i, bid, len(blkInfo.Transactions))

			blockManager.HandleTransactions(blkInfo.Transactions)
		}

		blk.Id = bid
		blockManager.SetLastBlock(blk)
		time.Sleep(time.Duration(e.Config.Sleep) * time.Second)
	}
}

func (e Explorer) Process() {
	client, err := utils.NewGethClient(e.Config)
	blockManager := block.NewBlockManager(e.Config)

	if err != nil {
		panic(err)
	}

	for {
		bid, err := client.EthBlockNumber()

		if err != nil {
			log.Println("Can't get last block:", err)
			continue
		}

		blk, err := blockManager.GetLastBlock(bid)

		if err != nil {
			log.Println("Can't get last block:", err)
			continue
		}

		for i := blk.Id; i <= bid; i++ {
			blkInfo, blockInfoErr := client.EthGetBlockByNumber(i, true)

			if blockInfoErr != nil {
				log.Println("Can't get block data!! ")
				i--
				continue
			}

			if len(blkInfo.Transactions) == 0 {
				log.Printf("Block %d from %d: No transactions found.\n", i, bid)
				continue
			}

			log.Printf("Block %d from %d: %d transactions found.\n", i, bid, len(blkInfo.Transactions))

			blockManager.HandleTransactions(blkInfo.Transactions)
		}

		blk.Id = bid
		blockManager.SetLastBlock(blk)
		time.Sleep(time.Duration(e.Config.Sleep) * time.Second)
	}
}
