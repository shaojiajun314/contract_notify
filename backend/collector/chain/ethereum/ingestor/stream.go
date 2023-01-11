package ethingestor

import (
	"fmt"
	"time"
	"sync"
	"errors"
	"context"
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)


type LogStream struct {
	client 				*ethclient.Client
    currentBlocknum 	*big.Int
    remoteBlocknum 		*big.Int
}


func NeBlockLogsStream(
	network string,
	rpcUrl string,
	startBlocknum *big.Int,
) (*LogStream, error) {
	client, e := ethclient.Dial(rpcUrl)
	if e != nil {
		return nil, e
	}
	return &LogStream{
		client: client,
		currentBlocknum: startBlocknum, 
		remoteBlocknum: new(big.Int).SetUint64(0),
	}, nil
}


func (stream *LogStream) Next() (*big.Int, *[]ethtypes.Log, error) {
	for {
        if stream.remoteBlocknum.Cmp(stream.currentBlocknum) <= 0 {
        	blocknumber, e := stream.LastBlocknum()
        	if e != nil {
        		return nil, nil, e
        	}
        	stream.remoteBlocknum = blocknumber
            if stream.remoteBlocknum.Cmp(stream.currentBlocknum) <= 0 {
                time.Sleep(time.Duration(10)*time.Second)
                continue
            }
        }
        scanBlocknum := new(big.Int).Add(stream.currentBlocknum, new(big.Int).SetUint64(1))
        logs, e := stream.LogsFromBlock(scanBlocknum)
        if e != nil {
        	return nil, nil, e
        }
        stream.currentBlocknum = scanBlocknum;
        return scanBlocknum, logs, nil
        
    }
}


func (stream *LogStream) LastBlocknum()(*big.Int, error) {
	curH, err := stream.client.BlockNumber(context.Background())
	for i:=0; err != nil && i < 5; i++ {
		time.Sleep(time.Duration(2)*time.Second)
        curH, err = stream.client.BlockNumber(context.Background())
	}

	if err != nil {
		return nil, errors.New(fmt.Sprintf("eth BlockNumber error: %v", err))
	}
	return new(big.Int).SetUint64(curH), nil
}


func (stream *LogStream) LogsFromBlock(blocknum *big.Int) (*[]ethtypes.Log, error) {
	
	block, e := stream.client.BlockByNumber(context.Background(), blocknum)
	for i:=0; e != nil && i < 5; i++ {
		time.Sleep(time.Duration(2)*time.Second)
        block, e = stream.client.BlockByNumber(context.Background(), blocknum)
	}

	if e != nil {
		return nil, errors.New(fmt.Sprintf(
			"error FilterLogs, %v", e,
		))
	}
	logs := []ethtypes.Log{}
	mutex := new(sync.Mutex)
	txes := block.Transactions()
    txLength := len(txes)
 	ch := make(chan error)
	for _, tx := range txes {
		go func(txArg common.Hash) {
			rp, err := stream.client.TransactionReceipt(context.Background(), txArg)
			for i:=0; err != nil && i < 5; i++ {
				time.Sleep(time.Duration(2)*time.Second)
			    rp, err = stream.client.TransactionReceipt(context.Background(), txArg)
			}

			if err != nil {
				ch <- err
				return
			}
			mutex.Lock()
			defer mutex.Unlock()
			for _, l := range rp.Logs {
				logs = append(logs, *l)
			}
			ch <- nil
		}(tx.Hash())
	}
	errs := ""
	for index := 0; index < txLength; index ++ {
		eChRet := <- ch
		if e != nil {
			errs = errs + fmt.Sprintf(
				"error get receipt, %v; ", eChRet,
			)
		}
	}
	if errs != "" {
		return nil, errors.New(errs)
	}
	return &logs, nil

}