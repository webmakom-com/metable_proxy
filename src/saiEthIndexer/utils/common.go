package utils

import (
	"encoding/json"
	"reflect"

	"github.com/saiset-co/saiEthIndexer/config"
)

func InArray(val interface{}, array interface{}) (index int) {
	values := reflect.ValueOf(array)

	if reflect.TypeOf(array).Kind() == reflect.Slice || values.Len() > 0 {
		for i := 0; i < values.Len(); i++ {
			if reflect.DeepEqual(val, values.Index(i).Interface()) {
				return i
			}
		}
	}

	return -1
}

func ConvertInterfaceToJson(obj interface{}) []byte {
	jsonResult, _ := json.Marshal(obj)
	return jsonResult
}

func RemoveContract(slice []config.Contract, s int) []config.Contract {
	return append(slice[:s], slice[s+1:]...)
}

func RemoveAddress(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
