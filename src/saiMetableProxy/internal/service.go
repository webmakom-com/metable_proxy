package internal

import (
	"github.com/saiset-co/saiMetableProxy/utils"
	"github.com/saiset-co/saiService"
)

// Init here we add all implemented handlers, create name of service and register config
func Init() {
	Service.Storage = NewDB()

	Service.Handler[GetNFTList.Name] = GetNFTList
	Service.Handler[GetCoursesNFTList.Name] = GetCoursesNFTList
	Service.Handler[GetNFTValue.Name] = GetNFTValue
	Service.Handler[getNFTByWalletAddress.Name] = getNFTByWalletAddress
	Service.Handler[getUtilityTokenBalanceByWallet.Name] = getUtilityTokenBalanceByWallet
	Service.Handler[getGovernanceTokenBalanceByWallet.Name] = getGovernanceTokenBalanceByWallet
	Service.Handler[getGovernanceTokenStakingAmountByWallet.Name] = getGovernanceTokenStakingAmountByWallet
}

type InternalService struct {
	Storage       utils.Database
	Handler       saiService.Handler  // handlers to define in this specified microservice
	GlobalService *saiService.Service // saiService reference
}

// Service global handler for registering handlers
var Service = &InternalService{
	Handler: saiService.Handler{},
}

func (s *InternalService) Init() {

}

func (s *InternalService) Process() {
	select {}
}
