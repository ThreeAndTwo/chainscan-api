package datasource

import "github.com/ThreeAndTwo/chainscan-api/types"

type IDataSource interface {
	GetTokenInfo() (*types.TokenInfo, error)
	GetABIData() (string, error)
	IsVerifyCode() (bool, error)
}
