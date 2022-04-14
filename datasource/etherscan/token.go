package etherscan

import (
	"encoding/json"
	"fmt"
	"github.com/ThreeAndTwo/chainscan-api/datasource"
	"github.com/ThreeAndTwo/chainscan-api/types"
	"github.com/imroc/req"
	"golang.org/x/time/rate"
)

type ether struct {
	source      string
	url         string
	apiKey      string
	contract    string
	rateLimiter *rate.Limiter
}

func NewEther(source, url, apiKey, contract string, rate *rate.Limiter) *ether {
	return &ether{source: source, url: url, apiKey: apiKey, contract: contract, rateLimiter: rate}
}

func (e *ether) GetTokenInfo() (*types.TokenInfo, error) {
	param := make(map[string]string)
	param["module"] = "token"
	param["action"] = "tokeninfo"
	param["address"] = e.contract
	param["apikey"] = e.apiKey
	reqParam := datasource.InitParam(param)

	net := datasource.NewNet(e.url, e.apiKey, req.Header{}, reqParam, datasource.GET)
	resp, err := net.Request()
	if err != nil {
		return nil, err
	}

	res := &types.EtherResult{}
	err = json.Unmarshal([]byte(resp), res)
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
		Github:      ethInfo.Github,
		Description: ethInfo.Description,
	}
	return tokenInfo, err
}

// TODO: get sourcecode

func (e *ether) GetABIData() (string, error) {
	abi, err := e.getAbiData()
	if err != nil {
		return "", err
	}

	if abi.Status != "1" {
		return "", fmt.Errorf("request service error, %s", abi)
	}

	return abi.Result.(string), nil
}

func (e *ether) getAbiData() (*types.EtherResult, error) {
	param := make(map[string]string)
	param["module"] = "contract"
	param["action"] = "getabi"
	param["address"] = e.contract
	param["apikey"] = e.apiKey
	reqParam := datasource.InitParam(param)

	net := datasource.NewNet(e.url, e.apiKey, req.Header{}, reqParam, datasource.GET)
	resp, err := net.Request()
	if err != nil {
		return nil, err
	}

	res := &types.EtherResult{}
	err = json.Unmarshal([]byte(resp), res)
	if err != nil {
		return nil, err
	}

	return res, err

}

func (e *ether) IsVerifyCode() (bool, error) {
	abi, err := e.getAbiData()
	if err != nil {
		return false, err
	}
	return abi.Status == "1", nil
}
