package ethingestor

import (
	"math/big"
	"contract_notify/collector/types"
)


func getHeadBlocknum(store types.BlockStoreDal, network string) (*big.Int, error) {
	return store.MaxBlockNum()
}


func getTailBlocknum(store types.BlockStoreDal, network string) (*big.Int, error) {
	return store.MinBlockNum()
}