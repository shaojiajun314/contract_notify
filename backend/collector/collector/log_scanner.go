package collector

import (
	"fmt"
	"errors"
	"math/big"

	"contract_notify/collector/types"
	"contract_notify/collector/chain"
)


type LogScanner struct {
	adapter 			types.Adapter
    remoteBlocknum 		*big.Int
}


func NeLogScanner(
	network string,
	rpcUrl string,
	parser types.Parser,
	startBlocknum *big.Int,
) (*LogScanner, error) {
	component, has := chain.DynChainMap[network]
	if !has {
		return nil, errors.New(fmt.Sprintf("error chain network(%v)", network))
	}
	adapter, e := component.AdapterConstructor(
		rpcUrl,
		parser,
	)
	if e != nil {
		return nil, errors.New(fmt.Sprintf("error new adapter (%v)", e))
	}
  	ls := LogScanner {
		adapter: *adapter,
		remoteBlocknum: new(big.Int).SetUint64(0),
	}
	return &ls, nil
}



func (ls *LogScanner) NextLogs(fromBlocknum *big.Int, toBlocknum *big.Int) (*[]types.Log, error) {
	if fromBlocknum.Cmp(toBlocknum) > 0 {
        panic(fmt.Sprintf(
        	"error fromBlocknum should be less than toBlocknum, %v, %v",
        	fromBlocknum,
        	toBlocknum,
        ))
    }
    if toBlocknum.Cmp(ls.remoteBlocknum) > 0 {
    	rpcBlocknum, e := ls.adapter.LastBlocknum()
    	if e != nil {
    		return nil, e
    	}
    	ls.remoteBlocknum = rpcBlocknum
    	if rpcBlocknum.Cmp(toBlocknum) < 0 {
	        panic("error to_blocknum， too large")
	    }
    }
    logs, e := ls.adapter.Logs(fromBlocknum, toBlocknum)
    return logs, e
}