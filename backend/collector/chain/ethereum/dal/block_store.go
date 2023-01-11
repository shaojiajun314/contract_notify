package ethdal

import (
	"errors"
	"math/big"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"contract_notify/collector/types"
)


type BlockStoreDal struct {
	Network string
	parser  types.Parser
} // todo for store


func NewBlockStoreDal(
	network string,
	parser types.Parser,
) types.BlockStoreDal {
	return BlockStoreDal{
		Network: network,
		parser: parser,
	}
}


func (bsd BlockStoreDal) Insert(
	blocknumber *big.Int,
	logs interface{},
) error {
	// todo
	TestDataStore[blocknumber] = logs.(*[]ethtypes.Log)
	return nil
}


func (bsd BlockStoreDal) Remove(
	blocknumber *big.Int,
) error {
	delete(TestDataStore, blocknumber)
	return nil
}


func (bsd BlockStoreDal) Fetch(
	fromNumber *big.Int,
	toBlocknum *big.Int,
) (map[*big.Int]*[]types.Log, error) {
	ret := make(map[*big.Int]*[]types.Log)
	ok := false
	for blocknumber, logs := range TestDataStore {
		if blocknumber.Cmp(fromNumber) == 0 {
			ok = true
		}
		if blocknumber.Cmp(fromNumber) < 0 {
			continue
		}
		if blocknumber.Cmp(toBlocknum) > 0 {
			continue
		}
		logsRet := []types.Log{}
		for _, l := range *logs {
			if parsedLog, hit, e := bsd.parser.ParseLog(l); hit {
				logsRet = append(logsRet, *parsedLog)
			} else if (e!=nil) {
				return ret, e
			}
		}
		ret[blocknumber] = &logsRet
	}
	if !ok {
		return ret, errors.New("from number block is deleted")
	}
	return ret, nil
}

func (bsd BlockStoreDal) MinBlockNum() (*big.Int, error) {
	// todo
	// return new(big.Int).SetUint64(0), nil
	return nil, nil
}
	
func (bsd BlockStoreDal) MaxBlockNum() (*big.Int, error) {
	// todo
	// return new(big.Int).SetUint64(0), nil
	return nil, nil
}


// Warnning remove； for test

var TestDataStore map[*big.Int]*[]ethtypes.Log

func init(){
	TestDataStore = make(map[*big.Int]*[]ethtypes.Log)
}