package block

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/onrik/ethrpc"
	"github.com/webmakom-com/saiContractExplorer/config"
	"github.com/webmakom-com/saiContractExplorer/utils"
	"github.com/webmakom-com/saiContractExplorer/utils/saiStorageUtil"
	"github.com/webmakom-com/saiContractExplorer/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

type Manager struct {
	abis      map[string]*abi.ABI
	config    config.Configuration
	storage   saiStorageUtil.Database
	websocket websocket.Manager
}

type BlockResult struct {
	Result string
}

type BlockData struct {
	Transactions []ethrpc.Transaction
}

type BlockInfo struct {
	Result BlockData
}

type Block struct {
	Id int `json:"id"`
}

var startBlock int
var manager Manager

func NewBlockManager(c config.Configuration) Manager {
	manager = Manager{
		config:    c,
		storage:   saiStorageUtil.Storage(c.Storage.Url, c.Storage.Auth.Email, c.Storage.Auth.Password),
		websocket: websocket.NewWebSocketManager(c),
		abis:      map[string]*abi.ABI{},
	}

	for _, contract := range c.Contracts {
		_abi, err := abi.JSON(strings.NewReader(contract.Data.ABI))

		if err != nil {
			log.Fatal(err)
		}

		manager.abis[contract.Data.Address] = &_abi
	}

	return manager
}

func (m Manager) GetLastBlock(id int) (Block, error) {
	block := Block{Id: id}
	pwd, err := os.Getwd()

	if err != nil {
		log.Println("Can't read current directory:", err)
		return block, nil
	}

	data, err := ioutil.ReadFile(pwd + "/block.data")

	if err != nil {
		log.Println("Can't open file:", err)
		return block, nil
	}

	lastBlock, strErr := strconv.Atoi(string(data))

	if strErr != nil {
		log.Println("Data from file can't be converted to int:", err)
		return block, nil
	}

	if lastBlock > 0 {
		startBlock = lastBlock + 1
	} else if m.config.StartBlock > 0 {
		startBlock = m.config.StartBlock
	} else {
		startBlock = id
	}

	return Block{Id: startBlock}, nil
}

func (m Manager) SetLastBlock(blk Block) {
	pwd, dirErr := os.Getwd()

	if dirErr != nil {
		log.Println("Can't read current directory:", dirErr)
		return
	}

	lastBlock := strconv.Itoa(blk.Id)
	writeErr := ioutil.WriteFile(pwd+"/block.data", []byte(lastBlock), 0777)

	if writeErr != nil {
		log.Println("Can't write file:", writeErr)
	}
}

func (m Manager) HandleTransactions(trs []ethrpc.Transaction) {
	for _, contract := range m.config.Contracts {
		for j := 0; j < len(trs); j++ {
			if strings.ToLower(trs[j].From) != strings.ToLower(contract.Data.Address) && strings.ToLower(trs[j].To) != strings.ToLower(contract.Data.Address) {
				continue
			}

			raw, _ := json.Marshal(trs[j])

			data := bson.M{
				"Hash":   trs[j].Hash,
				"From":   trs[j].From,
				"To":     trs[j].To,
				"Amount": trs[j].Value,
			}

			decodedSig, decodeSigErr := hex.DecodeString(trs[j].Input[2:10])

			if decodeSigErr != nil {
				log.Println("Decode sig error:", decodeSigErr)
				continue
			}

			method, methodErr := m.abis[contract.Data.Address].MethodById(decodedSig)

			if methodErr != nil {
				log.Println("Get method error:", methodErr)
				log.Println("ABI:", m.abis[contract.Data.Address])
				continue
			}

			decodedData, decodeDataErr := hex.DecodeString(trs[j].Input[2:])

			if decodeDataErr != nil {
				log.Println("Decode sig error:", decodeDataErr)
				continue
			}

			decodedInput := map[string]interface{}{}
			decodeInputErr := method.Inputs.UnpackIntoMap(decodedInput, decodedData[4:])

			if decodeInputErr != nil {
				log.Println("Decode input error:", decodeInputErr)
				continue
			}

			data["Operation"] = method.Name
			data["Input"] = decodedInput

			if utils.InArray(method, contract.Operations) != -1 {
				m.websocket.SendMessage(string(raw), m.config.WebSocket.Token)
			}

			storageErr, _ := m.storage.Put("transactions", data, m.config.Storage.Token)

			if storageErr != nil {
				log.Println("Storage error:", storageErr)
				continue
			}

			log.Printf("%d transaction from %s to %s has been updated.\n", trs[j].TransactionIndex, trs[j].From, trs[j].To)
		}
	}
}

func (m Manager) EthBlockNumber() (int, error) {
	data := bson.M{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"params":  bson.A{},
		"id":      1,
	}

	jsonStr, mErr := json.Marshal(data)

	if mErr != nil {
		log.Println("Marshaling error: ", mErr)
	}

	resp, err := http.Post(m.config.Geth.Web.Addresses[0], "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Println("Geth http: ", err)
	}

	var answer BlockResult
	jsonErr := json.NewDecoder(resp.Body).Decode(&answer)
	defer resp.Body.Close()

	if jsonErr != nil {
		log.Println("Wrong answer format from the geth: ", jsonErr)
	}

	numberStr := strings.Replace(answer.Result, "0x", "", -1)
	n, intError := strconv.ParseInt(numberStr, 16, 64)

	if intError != nil {
		log.Println("Wrong block number answer: ", intError)
	}

	return int(n), nil
}

func (m Manager) EthGetBlockByNumber(bid int, full bool) (BlockData, error) {
	data := bson.M{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  bson.A{bid, full},
		"id":      1,
	}

	jsonStr, mErr := json.Marshal(data)

	if mErr != nil {
		log.Println("Marshaling error: ", mErr)
	}

	resp, err := http.Post(m.config.Geth.Web.Addresses[0], "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Println("Geth http: ", err)
	}

	var answer BlockInfo
	jsonErr := json.NewDecoder(resp.Body).Decode(&answer)
	defer resp.Body.Close()

	if jsonErr != nil {
		log.Println("Wrong answer format from the geth: ", jsonErr)
	}

	return answer.Result, nil
}
