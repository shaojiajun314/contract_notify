package config

import (
	"math/big"
	collectorConfig "contract_notify/collector/conf"
)

var Defaults = &Config{
	Collector:    DefaultCollectorConfig(),
}

type Config struct {
	Collector    *collectorConfig.Config `toml:"collector"`
}

func DefaultCollectorConfig() *collectorConfig.Config {
	return &collectorConfig.Config{
		CollectIterBlockStep: new(big.Int).SetUint64(1000000000),
		Chains: []collectorConfig.ChainConfig{
			collectorConfig.ChainConfig{
				Network: "5",
				RPCUrl: "https://eth-goerli.g.alchemy.com/v2/5c8YbVzERLkTRMTYSotAYp07LR2cGbOk",
				CachedBlockLength: 10,
				StartBlocknum: 7853766,
			},
			// collectorConfig.ChainConfig{
			// 	Network: "4",
			// 	RPCUrl: "https://rpc.ankr.com/eth_rinkeby",
			// 	CachedBlockLength: 10,
			// 	StartBlocknum: 11769822,
			// },
			// collectorConfig.ChainConfig{
			// 	Network: "97",
			// 	RPCUrl: "https://bsctestapi.terminet.io/rpc",
			// 	CachedBlockLength: 10,
			// 	StartBlocknum: 25159159,
			// },
		},
  	}
}
