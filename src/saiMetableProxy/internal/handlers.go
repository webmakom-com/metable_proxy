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
		return "get nft info", nil
	},
}

// getNFTByWalletAddress Get nft value by wallet address
var getNFTByWalletAddress = saiService.HandlerElement{
	Name:        "getNFTByWalletAddress",
	Description: "Get nft value by wallet address",
	Function: func(data interface{}, token string) (interface{}, error) {
		return "get nft by wallet address", nil
	},
}

// getUtilityTokenBalanceByWallet Get utility token balance by wallet address
var getUtilityTokenBalanceByWallet = saiService.HandlerElement{
	Name:        "getUtilityTokenBalanceByWallet",
	Description: "Get utility token balance by wallet address",
	Function: func(data interface{}, token string) (interface{}, error) {
		return "get utility token balance by wallet", nil
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
