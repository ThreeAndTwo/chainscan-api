package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThreeAndTwo/chainscan-api/datasource"
	"github.com/ThreeAndTwo/chainscan-api/types"
	"github.com/imroc/req"
	"golang.org/x/time/rate"
	"strings"
	"time"
)

// api document: https://www.coingecko.com/en/api/documentation

type coingecko struct {
	source      string
	url         string
	apiKey      string
	rateLimiter *rate.Limiter
	market      *types.MarketMap
}

func NewCoinGecko(source, url, apiKey string, rate *rate.Limiter, market *types.MarketMap) *coingecko {
	if url == "" {
		url = "https://api.coingecko.com/api/v3/"
	}
	return &coingecko{source: source, url: url, apiKey: apiKey, rateLimiter: rate, market: market}
}

func (c *coingecko) check() bool {
	return c.url != ""
}

// GetMarketInfoForCoin /asset_platforms
func (c *coingecko) GetMarketInfoForCoin() ([]*types.MarketInfo, error) {
	_ = c.rateLimiter.Wait(context.Background())
	if !c.check() {
		return nil, fmt.Errorf("config mismatched for %s", c.source)
	}

	url := c.url + "asset_platforms"
	net := datasource.NewNet(url, req.Header{}, req.Param{}, datasource.GET)

	resp, err := net.Request()
	if err != nil {
		return nil, err
	}

	var cgMarket []*types.CoinGeckoMarket
	err = json.Unmarshal(resp, &cgMarket)
	if err != nil {
		return nil, err
	}

	var marketInfo []*types.MarketInfo
	for _, market := range cgMarket {
		_mi := &types.MarketInfo{
			ID: market.Id,
		}

		if market.Shortname != "" {
			_mi.Name = market.Shortname
		} else {
			_mi.Name = market.Name
		}
		marketInfo = append(marketInfo, _mi)
	}
	return marketInfo, nil
}

func (c *coingecko) GetSourceCode(contract string) ([]*types.EtherSourceCode, error) {
	return nil, fmt.Errorf("unSupport on CoinGecko")
}

func (c *coingecko) getMarketId() error {
	_ = c.rateLimiter.Wait(context.Background())
	c.market.Lock.Lock()
	defer c.market.Lock.Unlock()

	if time.Now().Unix()-c.market.LastUpdatedAt.Unix() > 86400 || len(c.market.Market[string(types.CoinGecko)]) == 0 {
		// force update
		markInfo, err := c.GetMarketInfoForCoin()
		if err != nil {
			return err
		}

		if c.market.Market[string(types.CoinGecko)] == nil {
			c.market.Market[string(types.CoinGecko)] = make(map[string]*types.MarketInfo)
		}

		for _, coin := range markInfo {
			c.market.Market[string(types.CoinGecko)][strings.ToLower(coin.Name)] = coin
		}
	}

	return nil
}

// GetTokenInfo /coins/binance-smart-chain/contract/0xb0d502e938ed5f4df2e681fe6e419ff29631d62b
func (c *coingecko) GetTokenInfo(contract string) (*types.TokenInfo, error) {
	if !c.check() {
		return nil, fmt.Errorf("config mismatched for %s", c.source)
	}

	_ = c.rateLimiter.Wait(context.Background())
	if err := c.getMarketId(); err != nil {
		return nil, err
	}

	if _, ok := c.market.Market[string(types.CoinGecko)][c.source]; !ok {
		return nil, fmt.Errorf("martket ID not exist")
	}

	url := c.url + "coins/" + c.market.Market[string(types.CoinGecko)][c.source].ID + "/contract/" + strings.ToLower(contract)
	net := datasource.NewNet(url, req.Header{}, req.Param{}, datasource.GET)
	resp, err := net.Request()
	if err != nil {
		return nil, err
	}

	var cgti = &types.CoinGeckoTokenInfo{}
	err = json.Unmarshal(resp, cgti)
	if err != nil {
		return nil, err
	}

	tokenInfo := &types.TokenInfo{
		Name:        cgti.Name,
		Symbol:      cgti.Symbol,
		Decimals:    "unknown",
		Type:        "ERC20",
		Twitter:     "https://twitter.com/" + cgti.Links.TwitterScreenName,
		Reddit:      "",
		Telegram:    cgti.Links.TelegramChannelIdentifier,
		Description: cgti.Description.En,
	}

	if len(cgti.Links.Homepage) != 0 {
		tokenInfo.Website = cgti.Links.Homepage[0]
	}

	if len(cgti.Links.ReposUrl.Github) != 0 {
		tokenInfo.Github = cgti.Links.ReposUrl.Github[0]
	}
	return tokenInfo, err
}

func (c *coingecko) GetABIData(contact string) (string, error) {
	return "", fmt.Errorf("unSupport on CoinGecko")
}

func (c *coingecko) IsVerifyCode(contact string) (bool, error) {
	return false, fmt.Errorf("unSupport on CoinGecko")
}
