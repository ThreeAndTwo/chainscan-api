package coinmarketcap

import (
	"encoding/json"
	"fmt"
	"github.com/ThreeAndTwo/chainscan-api/datasource"
	"github.com/ThreeAndTwo/chainscan-api/types"
	"golang.org/x/time/rate"
)

type cmc struct {
	url         string
	apiKey      string
	rateLimiter *rate.Limiter
}

func NewCmc(url string, apiKey string, rate *rate.Limiter) *cmc {
	return &cmc{url: url, apiKey: apiKey, rateLimiter: rate}
}

func (e *cmc) GetTokenInfo() (*types.TokenInfo, error) {
	param := make(map[string]string)
	param["module"] = "token"
	param["action"] = "tokeninfo"
	//param["address"] = e.contract
	reqParam := datasource.InitParam(param)

	header := make(map[string]string)
	header["X-CMC_PRO_API_KEY"] = e.apiKey
	header["Accept"] = "application/json"
	reqHeader, useJson := datasource.InitHeader(header)

	net := datasource.NewNet(e.url, e.apiKey, reqHeader, reqParam, datasource.POST)

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

func (e *cmc) GetABIData() (string, error) {

}

func (e *cmc) IsVerifyCode() (bool, error) {

}
