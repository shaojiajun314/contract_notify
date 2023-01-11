package ethereum

import (
	"fmt"
	"time"
	"errors"
	"context"
	"math/big"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"

	"contract_notify/collector/types"
	"contract_notify/collector/chain/ethereum/log_parser"
)


type Adapter struct {
	client 				*ethclient.Client
    parser				parser.Parser// AddressTopicsMap
    // logger 				log.Logger
}


func NewAdapter(
	rpcUrl string,
	parserInstance types.Parser,
	// logger log.Logger,
) (*types.Adapter, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("new eth adapter error, %v", err))
	}
	var adapterRet types.Adapter
	adapter := Adapter {
		client: client,
		parser: parserInstance.(parser.Parser),
		// logger: logger,
	}
	adapterRet = adapter
	return &adapterRet, nil
}



func (adapter Adapter) LastBlocknum() (*big.Int, error) {
	curH, err := adapter.client.BlockNumber(context.Background())
	for i:=0; err != nil && i < 5; i++ {
		time.Sleep(time.Duration(2)*time.Second)
        curH, err = adapter.client.BlockNumber(context.Background())
	}

	if err != nil {
		return nil, errors.New(fmt.Sprintf("eth BlockNumber error: %v", err))
	}
	return new(big.Int).SetUint64(curH), nil
}


func (adapter Adapter) Logs(fromHeight *big.Int, toHeight *big.Int) (*[]types.Log, error) {
	query := ethereum.FilterQuery{
		FromBlock: fromHeight,
		ToBlock:   toHeight,
		Addresses: adapter.parser.FilterContractAddresses,
		Topics: adapter.parser.FilterContractTopics,
	}

	logs, err := adapter.client.FilterLogs(context.Background(), query)
	for i:=0; err != nil && i < 5; i++ {
		time.Sleep(time.Duration(2)*time.Second)
        logs, err = adapter.client.FilterLogs(context.Background(), query)
	}
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"error FilterLogs, %v", err,
		))
	}

	parsedEvents := []types.Log{}
	for _, log := range logs {
		if parsedEvent, hit, e := adapter.parser.ParseLog(log); hit {
			parsedEvents = append(parsedEvents, *parsedEvent)
		} else if (e != nil) {
			return nil, e
		}
	}

	return &parsedEvents, nil
}