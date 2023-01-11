package conf

import (
  "math/big"
)



type ChainConfig struct {
    Network             string `toml:"network"`
    RPCUrl              string `toml:"rpcurl"`
    CachedBlockLength   uint64 `toml:"cached_block_length"`
    StartBlocknum       uint64 `toml:"start_blocknumber"`
}


type Config struct {
    Chains                  []ChainConfig  `toml:"chains"`
    CollectIterBlockStep    *big.Int       `toml:"collect_step"`
}


var CollectIterBlockStep *big.Int


func RegisteCollectIterBlockStep(n *big.Int) {
  if CollectIterBlockStep != nil {
    panic("duplicate registe")
  }
  CollectIterBlockStep = n
}