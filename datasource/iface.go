package datasource

import "github.com/ThreeAndTwo/chainscan-api/types"

type IDataSource interface {
	GetMarketInfoForCoin() ([]*types.MarketInfo, error)
	GetTokenInfo(string) (*types.TokenInfo, error)
	GetSourceCode(string) ([]*types.EtherSourceCode, error)
	GetABIData(string) (string, error)
	IsVerifyCode(string) (bool, error)
}
