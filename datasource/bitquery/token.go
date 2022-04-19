package bitquery

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThreeAndTwo/chainscan-api/datasource"
	"github.com/ThreeAndTwo/chainscan-api/types"
	"golang.org/x/time/rate"
)

type bitQuery struct {
	source      string
	url         string
	apiKey      string
	rateLimiter *rate.Limiter
}

func (b *bitQuery) GetMarketInfoForCoin() ([]*types.MarketInfo, error) {
	return nil, fmt.Errorf("unSupport on bitQuery")
}

func (b *bitQuery) GetTokenInfo(contract string) (*types.TokenInfo, error) {
	_ = b.rateLimiter.Wait(context.Background())
	params := &types.BitQueryParams{
		Query: queryTokenInfo,
		Variables: fmt.Sprintf(`{
    "limit": 10,
    "offset": 0,
    "network": "%s",
    "address": "%s",
    "from": "2008-01-01",
    "till": "2022-04-13T23:59:59"
  }`, b.source, contract),
	}

	header := make(map[string]string)
	header["X-API-KEY"] = b.apiKey
	header["Content-Type"] = "application/json"
	reqHeader, _ := datasource.InitHeader(header)

	bParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	mapParam := make(map[string]interface{})
	err = json.Unmarshal(bParams, &mapParam)
	if err != nil {
		return nil, err
	}

	reqParams := datasource.InitParam(mapParam)
	net := datasource.NewNet(b.url, reqHeader, reqParams, datasource.POST)
	net.SetJson(true)

	resp, err := net.Request()
	if err != nil {
		return nil, err
	}

	res := &types.BitQueryTokenResult{}
	err = json.Unmarshal(resp, res)
	if err != nil {
		return nil, err
	}

	if len(res.Data.Ethereum.SmartContractCalls) == 0 {
		return nil, fmt.Errorf("not data")
	}

	return &types.TokenInfo{
		Creator:       res.Data.Ethereum.SmartContractCalls[0].CreatedBy,
		DeployedBlock: res.Data.Ethereum.SmartContractCalls[0].CratedBlock,
		DeployedTx:    res.Data.Ethereum.SmartContractCalls[0].CreatedTx,
		DeployedAt:    res.Data.Ethereum.SmartContractCalls[0].Created,
	}, nil
}

func (b *bitQuery) GetSourceCode(contract string) ([]*types.EtherSourceCode, error) {
	return nil, fmt.Errorf("unSupport on bitQuery")
}

func (b *bitQuery) GetABIData(contract string) (string, error) {
	return "", fmt.Errorf("unSupport on bitQuery")
}

func (b *bitQuery) IsVerifyCode(contract string) (bool, error) {
	return false, fmt.Errorf("unSupport on bitQuery")
}

func NewBitQuery(source string, url string, apiKey string, rateLimiter *rate.Limiter) *bitQuery {
	if url == "" {
		url = "https://graphql.bitquery.io/"
	}
	return &bitQuery{source: source, url: url, apiKey: apiKey, rateLimiter: rateLimiter}
}
