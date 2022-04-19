package coinmarketcap

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThreeAndTwo/chainscan-api/datasource"
	"github.com/ThreeAndTwo/chainscan-api/types"
	"github.com/imroc/req"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/time/rate"
	"strings"
)

// api document: https://pro-api.coinmarketcap.com/v2/cryptocurrency/info

type cmc struct {
	source      string
	url         string
	apiKey      string
	rateLimiter *rate.Limiter
	market      *types.MarketMap
}

func NewCmc(source, url, apiKey string, rate *rate.Limiter, market *types.MarketMap) *cmc {
	if url == "" {
		url = "https://pro-api.coinmarketcap.com"
	}
	return &cmc{source: source, url: url, apiKey: apiKey, rateLimiter: rate, market: market}
}

func (c *cmc) check() bool {
	return c.url != "" && c.apiKey != ""
}

// GetMarketInfoForCoin /v1/cryptocurrency/map
func (c *cmc) GetMarketInfoForCoin() ([]*types.MarketInfo, error) {
	if !c.check() {
		return nil, fmt.Errorf("config mismatched for %s", c.source)
	}

	_ = c.rateLimiter.Wait(context.Background())

	header := make(map[string]string)
	header["X-CMC_PRO_API_KEY"] = c.apiKey
	header["Accept"] = "application/json"
	reqHeader, _ := datasource.InitHeader(header)
	url := c.url + "/v1/cryptocurrency/map"

	net := datasource.NewNet(url, reqHeader, req.Param{}, datasource.GET)
	resp, err := net.Request()
	if err != nil {
		return nil, err
	}

	res := &types.CmcResult{}
	err = json.Unmarshal(resp, res)
	if err != nil {
		return nil, err
	}

	var marketInfo []*types.MarketInfo
	for _, _market := range res.Data.([]interface{}) {
		data := &types.MarketInfo{}
		if err = mapstructure.Decode(_market, data); err != nil {
			return nil, err
		}

		marketInfo = append(marketInfo, data)
	}
	return marketInfo, err
}

func (c *cmc) GetSourceCode(contract string) ([]*types.EtherSourceCode, error) {
	return nil, fmt.Errorf("unSupport for CoinMarketCap")
}

func (c *cmc) GetTokenInfo(contract string) (*types.TokenInfo, error) {
	if !c.check() {
		return nil, fmt.Errorf("config mismatched for %s", c.source)
	}

	_ = c.rateLimiter.Wait(context.Background())

	header := make(map[string]string)
	header["X-CMC_PRO_API_KEY"] = c.apiKey
	header["Accept"] = "application/json"
	reqHeader, _ := datasource.InitHeader(header)
	url := c.url + "/v2/cryptocurrency/info?address=" + strings.ToLower(contract)
	net := datasource.NewNet(url, reqHeader, req.Param{}, datasource.GET)

	resp, err := net.Request()
	if err != nil {
		return nil, err
	}

	res := &types.CmcResult{}
	err = json.Unmarshal(resp, res)
	if err != nil {
		return nil, err
	}

	if res.Status.ErrorCode != 0 {
		return nil, fmt.Errorf("request service error, %s", resp)
	}

	_tokenInfo := make(map[string]types.CmcTokenInfo)
	err = mapstructure.Decode(res.Data, &_tokenInfo)
	if err != nil {
		return nil, err
	}

	key := ""
	for k := range _tokenInfo {
		key = k
	}

	tokenInfo := &types.TokenInfo{
		Name:        _tokenInfo[key].Name,
		Symbol:      _tokenInfo[key].Symbol,
		Decimals:    "",
		Type:        "ERC20",
		Description: _tokenInfo[key].Description,
	}

	if len(_tokenInfo[key].Urls.Website) != 0 {
		tokenInfo.Website = _tokenInfo[key].Urls.Website[0]
	}

	if len(_tokenInfo[key].Urls.Twitter) != 0 {
		tokenInfo.Twitter = _tokenInfo[key].Urls.Twitter[0]
	}

	if len(_tokenInfo[key].Urls.Reddit) != 0 {
		tokenInfo.Reddit = _tokenInfo[key].Urls.Reddit[0]
	}

	if len(_tokenInfo[key].Urls.Chat) != 0 {
		for _, chatType := range _tokenInfo[key].Urls.Chat {
			if strings.Contains(chatType, "discord.com") {
				tokenInfo.Discord = chatType
			}
			if strings.Contains(chatType, "https://t.me") {
				tokenInfo.Telegram = chatType
			}
		}
	}

	if len(_tokenInfo[key].Urls.SourceCode) != 0 {
		tokenInfo.Github = _tokenInfo[key].Urls.SourceCode[0]
	}

	return tokenInfo, err
}

func (c *cmc) GetABIData(contact string) (string, error) {
	return "", fmt.Errorf("unSupport for CoinMarketCap")
}

func (c *cmc) IsVerifyCode(contact string) (bool, error) {
	return false, fmt.Errorf("unSupport for CoinMarketCap")
}
