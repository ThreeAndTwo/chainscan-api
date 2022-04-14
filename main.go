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

func NewDataSource(source, alias, url, apiKey, contract string, tps int) (datasource.IDataSource, error) {
	platform := types.PlatformForDataSource("")
	if alias == "" {
		platform = types.PlatformForDataSource(source)
	} else {
		platform = types.PlatformForDataSource(alias)
	}

	rateLimiter := rate.NewLimiter(rate.Every(time.Second*1), tps)

	switch platform {
	case types.EtherScan:
		return etherscan.NewEther(source, url, apiKey, contract, rateLimiter), nil
	case types.CoinMarketCap:
		return coinmarketcap.NewCmc(url, apiKey, rateLimiter), nil
	case types.CoinGecko:
		return coingecko.NewCoinGecko(url, apiKey, rateLimiter), nil
	default:
		return nil, fmt.Errorf("unknown datasource for %s source. plz check it", source)
	}
}
