package etherscan

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThreeAndTwo/chainscan-api/datasource"
	"github.com/ThreeAndTwo/chainscan-api/types"
	"github.com/imroc/req"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/time/rate"
)

type ether struct {
	source      string
	url         string
	apiKey      string
	rateLimiter *rate.Limiter
}

func NewEther(source, url, apiKey string, rate *rate.Limiter) *ether {
	if url[len(url)-1:] != "?" {
		url += "?"
	}
	return &ether{source: source, url: url, apiKey: apiKey, rateLimiter: rate}
}

func (e *ether) GetMarketInfoForCoin() ([]*types.MarketInfo, error) {
	return nil, fmt.Errorf("unSupport for %s source", e.source)
}

func (e *ether) check() bool {
	return e.url != "" && e.apiKey != ""
}

func (e *ether) GetTokenInfo(contract string) (*types.TokenInfo, error) {
	if !e.check() {
		return nil, fmt.Errorf("config mismatched for %s", e.source)
	}

	_ = e.rateLimiter.Wait(context.Background())
	url := e.url + "module=token&action=tokeninfo&address=" + contract + "&apiKey=" + e.apiKey
	net := datasource.NewNet(url, req.Header{}, req.Param{}, datasource.GET)
	resp, err := net.Request()
	if err != nil {
		return nil, err
	}

	res := &types.EtherResult{}
	err = json.Unmarshal(resp, res)
	if err != nil {
		return nil, err
	}

	if res.Status != "1" {
		return nil, fmt.Errorf("request service error, %s", resp)
	}

	ethInfo := res.Result.(types.EtherTokenInfo)
	tokenInfo := &types.TokenInfo{
		Name:        ethInfo.TokenName,
		Symbol:      ethInfo.Symbol,
		Decimals:    ethInfo.Divisor,
		Type:        ethInfo.TokenType,
		Website:     ethInfo.Website,
		Twitter:     ethInfo.Twitter,
		Reddit:      ethInfo.Reddit,
		Telegram:    ethInfo.Telegram,
		Discord:     ethInfo.Discord,
		Github:      ethInfo.Github,
		Description: ethInfo.Description,
	}
	return tokenInfo, err
}

func (e *ether) GetSourceCode(contract string) ([]*types.EtherSourceCode, error) {
	if !e.check() {
		return nil, fmt.Errorf("config mismatched for %s", e.source)
	}
	_ = e.rateLimiter.Wait(context.Background())

	url := e.url + "module=contract&action=getsourcecode&address=" + contract + "&apiKey=" + e.apiKey
	net := datasource.NewNet(url, req.Header{}, req.Param{}, datasource.GET)
	resp, err := net.Request()
	if err != nil {
		return nil, err
	}

	res := &types.EtherResult{}
	err = json.Unmarshal(resp, res)
	if err != nil {
		return nil, err
	}

	if res.Status != "1" {
		return nil, fmt.Errorf("request service error for %s scan", e.source)
	}

	var sourceCode []*types.EtherSourceCode
	for _, _codeRes := range res.Result.([]interface{}) {
		code := &types.EtherSourceCode{}
		if err = mapstructure.Decode(_codeRes, code); err != nil {
			return nil, err
		}

		sourceCode = append(sourceCode, code)
	}
	return sourceCode, err
}

func (e *ether) GetABIData(contract string) (string, error) {
	if !e.check() {
		return "", fmt.Errorf("config mismatched for %s", e.source)
	}
	_ = e.rateLimiter.Wait(context.Background())

	abi, err := e.getAbiData(contract)
	if err != nil {
		return "", err
	}

	if abi.Status != "1" {
		return "", fmt.Errorf("request service error, %s", abi)
	}

	return abi.Result.(string), nil
}

func (e *ether) getAbiData(contract string) (*types.EtherResult, error) {
	url := e.url + "module=contract&action=getabi&address=" + contract + "&apiKey=" + e.apiKey
	net := datasource.NewNet(url, req.Header{}, req.Param{}, datasource.GET)
	resp, err := net.Request()
	if err != nil {
		return nil, err
	}

	res := &types.EtherResult{}
	err = json.Unmarshal(resp, res)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (e *ether) IsVerifyCode(contract string) (bool, error) {
	if !e.check() {
		return false, fmt.Errorf("config mismatched for %s", e.source)
	}
	_ = e.rateLimiter.Wait(context.Background())

	abi, err := e.getAbiData(contract)
	if err != nil {
		return false, err
	}
	return abi.Status == "1", nil
}
