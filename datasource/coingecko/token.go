package coingecko

import (
	"github.com/ThreeAndTwo/chainscan-api/types"
	"golang.org/x/time/rate"
)

type coingecko struct {
	url         string
	apiKey      string
	rateLimiter *rate.Limiter
}

// https://pro-api.coinmarketcap.com/v2/cryptocurrency/info

func NewCoinGecko(url string, apiKey string, rate *rate.Limiter) *coingecko {
	return &coingecko{url: url, apiKey: apiKey, rateLimiter: rate}
}

func (e *coingecko) GetTokenInfo() (*types.TokenInfo, error) {

}

func (e *coingecko) GetABIData() (string, error) {

}

func (e *coingecko) IsVerifyCode() (bool, error) {

}
