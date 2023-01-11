package types

import (
	"math/big"

)


type FilterQuery struct {
	Address string		`json:"address"`
	Events 	[]string	`json:"events"`
	ABI 	string		`json:"abi"`
}

type CollectArg struct {
	ChainID  		uint64 			`json:"chain_id"`// network
	RPCURL   		string			`json:"rpc"`
	FilterQuery 	[]FilterQuery	`json:"filters"`
	StartBlocknum 	uint64			`json:"start_number"`
}

type Adapter interface {
	LastBlocknum() (*big.Int, error)
	Logs(fromHeight *big.Int, toHeight *big.Int) (*[]Log, error)
}


type AdapterConstructor func(rpcUrl string, parserInstance Parser) (*Adapter, error)