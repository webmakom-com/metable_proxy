package config

import (
	"sync"

	valid "github.com/asaskevich/govalidator"
	configinternal "github.com/saiset-co/saiEthIndexer/internal/config-internal"
)

type Configuration struct {
	Common   configinternal.Common `json:"common"` // built-in framework config
	Specific `json:"specific"`
	EthContracts
}

// Specific - specific for current microservice settings
type Specific struct {
	GethServer string `json:"geth_server"`
	Storage    `json:"storage"`
	StartBlock int      `json:"start_block"`
	Operations []string `json:"operations"`
	Sleep      int      `json:"sleep"`
	WebSocket  `json:"websocket"`
}

// settings for saiStorage
type Storage struct {
	Token    string `json:"token"`
	URL      string `json:"url"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Contract struct {
	Address    string `json:"address" valid:",required"`
	ABI        string `json:"abi" valid:",required"`
	StartBlock int    `json:"start_block" valid:",required"`
}

func (r *Contract) Validate() error {
	_, err := valid.ValidateStruct(r)
	return err
}

type EthContracts struct {
	Mutex     *sync.RWMutex `json:"-"`
	Contracts []Contract    `json:"contracts"`
}

type WebSocket struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}
