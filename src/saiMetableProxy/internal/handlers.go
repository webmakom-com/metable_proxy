// Package internal implementing handlers
package internal

import (
	"encoding/json"
	"fmt"
	"github.com/saiset-co/saiService"
	"go.mongodb.org/mongo-driver/bson"
)

// GetNFTList Get nft value by id
var GetNFTList = saiService.HandlerElement{
	Name:        "GetNFTList",
	Description: "Get nft list",
	Function: func(data interface{}, token string) (interface{}, error) {
		err, response := Service.Storage.Get("transactions", bson.M{"$or": bson.A{
			bson.M{"From": "0x4c504a8fba715b05512eff6ac25934dfdc34373c"},
			bson.M{"To": "0x4c504a8fba715b05512eff6ac25934dfdc34373c"},
		}, "Operation": "Mint"}, bson.M{})
		if err != nil {
			return nil, fmt.Errorf("can not get transactions from the storage : %w", err)
		}

		result := new(interface{})
		json.Unmarshal(response, &result)

		return result, nil
	},
}

// GetCoursesNFTList Get courses nft value by id
var GetCoursesNFTList = saiService.HandlerElement{
	Name:        "GetCoursesNFTList",
	Description: "Get courses nft list",
	Function: func(data interface{}, token string) (interface{}, error) {
		err, response := Service.Storage.Get("transactions", bson.M{"$or": bson.A{
			bson.M{"From": "0x5fb589b1b4e7129a3747c2786a2ad668cc0e7eb8"},
			bson.M{"To": "0x5fb589b1b4e7129a3747c2786a2ad668cc0e7eb8"},
		}, "Operation": "Mint"}, bson.M{})
		if err != nil {
			return nil, fmt.Errorf("can not get transactions from the storage : %w", err)
		}

		result := new(interface{})
		json.Unmarshal(response, &result)

		return result, nil
	},
}

// GetNFTValue Get nft value by id
var GetNFTValue = saiService.HandlerElement{
	Name:        "GetNFTValue",
	Description: "Get nft value by id",
	Function: func(data interface{}, token string) (interface{}, error) {
		idp, ok := data.(map[string]interface{})["id"].(float64)

		if !ok {
			return nil, fmt.Errorf("wrong data parameter")
		}

		id := int(idp)

		err, response1 := Service.Storage.Get("transactions", bson.M{"To": "0x4c504a8fba715b05512eff6ac25934dfdc34373c", "Operation": "Mint"}, bson.M{"skip": id - 1, "limit": 1})
		if err != nil {
			return nil, fmt.Errorf("can not get transactions from the storage : %w", err)
		}

		err, response2 := Service.Storage.Get("transactions", bson.M{"To": "0x4c504a8fba715b05512eff6ac25934dfdc34373c", "Input.tokenId": id}, bson.M{})
		if err != nil {
			return nil, fmt.Errorf("can not get transactions from the storage : %w", err)
		}

		var wrappedResult1 map[string][]interface{}
		var wrappedResult2 map[string][]interface{}

		jsonErr1 := json.Unmarshal(response1, &wrappedResult1)

		if jsonErr1 != nil {
			fmt.Println(string(response1))
			fmt.Println(jsonErr1)
			return nil, jsonErr1
		}

		jsonErr2 := json.Unmarshal(response2, &wrappedResult2)

		if jsonErr2 != nil {
			fmt.Println(string(response2))
			fmt.Println(jsonErr2)
			return nil, jsonErr2
		}

		response := append(wrappedResult1["result"], wrappedResult2["result"]...)

		return response, nil
	},
}

// getNFTByWalletAddress Get nft value by wallet address
var getNFTByWalletAddress = saiService.HandlerElement{
	Name:        "getNFTByWalletAddress",
	Description: "Get nft value by wallet address",
	Function: func(data interface{}, token string) (interface{}, error) {
		wallet, ok := data.(map[string]interface{})["wallet"].(string)

		if !ok {
			return nil, fmt.Errorf("wrong data parameter")
		}

		err, response := Service.Storage.Get("transactions", bson.M{"$or": bson.A{
			bson.M{"From": wallet},
			bson.M{"To": wallet},
		}}, bson.M{})
		if err != nil {
			return nil, fmt.Errorf("can not get transactions from the storage : %w", err)
		}

		result := new(interface{})
		json.Unmarshal(response, &result)

		return result, nil
	},
}

// getUtilityTokenBalanceByWallet Get utility token balance by wallet address
var getUtilityTokenBalanceByWallet = saiService.HandlerElement{
	Name:        "getUtilityTokenBalanceByWallet",
	Description: "Get utility token balance by wallet address",
	Function: func(data interface{}, token string) (interface{}, error) {
		var balance float64
		wallet, ok := data.(map[string]interface{})["wallet"].(string)

		if !ok {
			return nil, fmt.Errorf("wrong data parameter")
		}

		err, response := Service.Storage.Get("transactions", bson.M{"$or": bson.A{
			bson.M{"From": wallet},
			bson.M{"To": wallet},
		}}, bson.M{})

		if err != nil {
			return nil, fmt.Errorf("can not get transactions from the storage : %w", err)
		}

		var result map[string][]interface{}
		json.Unmarshal(response, &result)

		for _, v := range result["result"] {
			transaction, ok := v.(map[string]interface{})

			if !ok {
				return nil, fmt.Errorf("can not get transactions from the response : %v", v)
			}
			input, ok := transaction["Input"].(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("can not get input from the transaction type: %v", transaction["Input"])
			}

			amountFloat, ok := input["amount"].(float64)
			if !ok {
				continue
			}

			if transaction["From"] == wallet {
				balance = balance - amountFloat
			}

		}

		return balance, nil
	},
}

// getGovernanceTokenBalanceByWallet Get governance token balance by wallet address
var getGovernanceTokenBalanceByWallet = saiService.HandlerElement{
	Name:        "getGovernanceTokenBalanceByWallet",
	Description: "Get governance token balance by wallet address",
	Function: func(data interface{}, token string) (interface{}, error) {
		return "get governance token balance by wallet", nil
	},
}

// getGovernanceTokenStakingAmountByWallet Get nft value by wallet address
var getGovernanceTokenStakingAmountByWallet = saiService.HandlerElement{
	Name:        "getGovernanceTokenStakingAmountByWallet",
	Description: "Get governance token staking amount by wallet address",
	Function: func(data interface{}, token string) (interface{}, error) {
		return "get governance token staking amount by wallet", nil
	},
}
