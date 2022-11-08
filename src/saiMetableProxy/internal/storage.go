package internal

import (
	"log"
	"strconv"

	"github.com/saiset-co/saiMetableProxy/utils"
)

func NewDB() utils.Database {
	host, ok := Service.GlobalService.GetConfig("specific.storage.host", "").(string)
	if !ok {
		log.Fatalf("configuration : invalid storage host provided, url : %s", Service.GlobalService.GetConfig("specific.storage.host", ""))
	}

	port, ok := Service.GlobalService.GetConfig("specific.storage.port", "").(int)
	if !ok {
		log.Fatalf("configuration : invalid storage port provided, url : %s", Service.GlobalService.GetConfig("specific.storage.port", ""))
	}

	token, ok := Service.GlobalService.GetConfig("specific.storage.token", "").(string)
	if !ok {
		log.Fatalf("configuration : invalid storage email provided, email : %s", Service.GlobalService.GetConfig("specific.storage.token", ""))
	}

	return utils.Storage(host+":"+strconv.Itoa(port), token)
}
