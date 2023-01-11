package types

import (
	"math/big"
)


type BlockStoreDal interface {
	Insert(blocknumber *big.Int, logs interface{}) error

	Remove(blocknumber *big.Int) error

	Fetch(fromNumber *big.Int, toBlocknum *big.Int) (map[*big.Int]*[]Log, error)

	MinBlockNum() (*big.Int, error)

	MaxBlockNum() (*big.Int, error)
}


type NewBlockStoreDal func(network string, parser Parser) BlockStoreDal