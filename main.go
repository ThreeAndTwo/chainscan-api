package chainscan_api

import (
	"fmt"
	"github.com/ThreeAndTwo/chainscan-api/datasource"
	"github.com/ThreeAndTwo/chainscan-api/datasource/coingecko"
	"github.com/ThreeAndTwo/chainscan-api/datasource/coinmarketcap"
	"github.com/ThreeAndTwo/chainscan-api/datasource/etherscan"
	"github.com/ThreeAndTwo/chainscan-api/types"
	"golang.org/x/time/rate"
	"time"
)

func NewDataSource(source string, alias types.PlatformForDataSource, url, apiKey string, tps int) (datasource.IDataSource, error) {
	platform := types.PlatformForDataSource("")
	if alias == "" {
		platform = types.PlatformForDataSource(source)
	} else {
		platform = alias
	}

	if tps <= 0 {
		tps = 1
	}

	marketMap := &types.MarketMap{
		Market:        make(map[string]map[string]*types.MarketInfo),
		LastUpdatedAt: time.Now(),
	}
	rateLimiter := rate.NewLimiter(rate.Every(time.Second*1), tps)

	switch platform {
	case types.EtherScan:
		return etherscan.NewEther(source, url, apiKey, rateLimiter), nil
	case types.CoinMarketCap:
		return coinmarketcap.NewCmc(source, url, apiKey, rateLimiter, marketMap), nil
	case types.CoinGecko:
		return coingecko.NewCoinGecko(source, url, apiKey, rateLimiter, marketMap), nil
	default:
		return nil, fmt.Errorf("unknown datasource for %s source. plz check it", source)
	}
}
