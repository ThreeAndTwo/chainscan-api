package chainscan_api

import (
	"github.com/ThreeAndTwo/chainscan-api/types"
	"os"
	"strings"
	"testing"
)

func TestNewDataSource(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		alias    types.PlatformForDataSource
		url      string
		apiKey   string
		contract string
		tps      int
	}{
		{
			name:     "test eth",
			source:   "etherscan",
			alias:    "",
			url:      "https://api.etherscan.io/api",
			apiKey:   os.Getenv("apiKey"),
			contract: "0xAf5191B0De278C7286d6C7CC6ab6BB8A73bA2Cd6",
			tps:      1,
		},
		{
			name:     "test bsc",
			source:   "bsc",
			alias:    types.EtherScan,
			url:      "https://api.bscscan.com/api",
			apiKey:   os.Getenv("apiKey"),
			contract: "0xB0D502E938ed5f4df2E681fE6E419ff29631d62b",
			tps:      1,
		},
		{
			name:     "test cmc",
			source:   "",
			alias:    types.CoinMarketCap,
			apiKey:   os.Getenv("apiKey"),
			contract: "0xB0D502E938ed5f4df2E681fE6E419ff29631d62b",
			tps:      1,
		},
		{
			name:     "test coinGecko",
			source:   "bsc",
			alias:    types.CoinGecko,
			apiKey:   os.Getenv("apiKey"),
			contract: "0xB0D502E938ed5f4df2E681fE6E419ff29631d62b",
			tps:      1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, err := NewDataSource(tt.source, tt.alias, tt.url, tt.apiKey, tt.tps)
			if err != nil {
				t.Errorf("err: %s", err)
			}

			coin, err := source.GetMarketInfoForCoin()
			if err != nil && !strings.Contains(err.Error(), "unSupport") {
				t.Fatalf("unknown error: %s, for %s platform", err, tt.name)
			}

			t.Logf("coin: %v", coin)

			sourceCode, err := source.GetSourceCode(tt.contract)
			if err != nil && !strings.Contains(err.Error(), "unSupport") {
				t.Fatalf("get source code error: %s, for %s platform", err, tt.name)
			}

			t.Logf("sourceCode: %v", sourceCode)

			abiData, err := source.GetABIData(tt.contract)
			if err != nil && !strings.Contains(err.Error(), "unSupport") {
				t.Fatalf("get abi data error: %s, for %s platform", err, tt.name)
			}
			t.Logf("AbiData: %s", abiData)

			IsVerified, err := source.IsVerifyCode(tt.contract)
			if err != nil && !strings.Contains(err.Error(), "unSupport") {
				t.Fatalf("get IsVerifyCode error: %s, for %s platform", err, tt.name)
			}

			t.Logf("status: %v", IsVerified)

			info, err := source.GetTokenInfo(tt.contract)
			if err != nil {
				t.Fatalf("get tokenInfo error: %s, for %s platform", err, tt.name)
			}
			t.Logf("tokenInfo: %v", info)
		})
	}
}
