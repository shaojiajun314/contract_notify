package ethingestor


import (
	"fmt"
	"errors"
	"math/big"

	"contract_notify/collector/types"
	"contract_notify/collector/chain/ethereum/dal"
	"contract_notify/collector/chain/ethereum/log_parser"
)


type Ingestor struct {
	Network 			string
	RPCUrl				string
	CachedBlockLength 	*big.Int
	state 				*types.IngestorState
	Store  				types.BlockStoreDal
}


func NewIngestor(
	network string,
	rpcUrl string,
	cachedBlockLength uint64,
	startBlocknum uint64,
	IngestorState *types.IngestorState,
) (*types.Ingestor, error) {
	store := ethdal.NewBlockStoreDal(network, parser.Parser{}) // todo
	var th *big.Int
	hh, e := getHeadBlocknum(store, network)
	if e != nil {
		return nil, e
	}
	if hh == nil {
		var s uint64
		if startBlocknum > 0 {
			s = startBlocknum - 1
		} else {
			s = 0
		}
		hh = new(big.Int).SetUint64(s)
		th = new(big.Int).SetUint64(s)
	} else {
		th, er := getTailBlocknum(store, network)
		if er != nil {
			return nil, er
		}
		if th == nil {
			return nil, errors.New(fmt.Sprintf("some error hit the chain(%v) store", network))
		}
	}
	IngestorState.HeadNumber = hh
	IngestorState.TailNumber = th
	ingestorInstance := types.Ingestor(
		Ingestor {
			Network: network,
			RPCUrl: rpcUrl,
			CachedBlockLength: new(big.Int).SetUint64(cachedBlockLength),
			Store: store,
			state: IngestorState,
		},
	)
	return &ingestorInstance, nil
}


func (ingestor Ingestor) Run() error {
	stream, e:= NeBlockLogsStream(
		ingestor.Network,
		ingestor.RPCUrl,
		ingestor.state.HeadNumber,
	)
	if e != nil {
		panic(e)
	}
	one := new(big.Int).SetUint64(1)
	for {
		number, logs, e := stream.Next()
		if e != nil {
			return e
		}
		// fmt.Println(logs, " block_logs");
		ingestor.Store.Insert(number, logs)
		fmt.Println("insert logs blocknum: ", number, ingestor.CachedBlockLength)
		if new(big.Int).Sub(number, ingestor.state.TailNumber).Cmp(ingestor.CachedBlockLength) > 0 {
			removingNum := new(big.Int).Sub(number, ingestor.CachedBlockLength)
			ingestor.Store.Remove(removingNum)
			fmt.Println("remove logs blocknum: ", removingNum)
			// ingestor.state.RwMutex.Lock()
			ingestor.state.TailNumber = new(big.Int).Add(removingNum, one)
			ingestor.state.HeadNumber = number
			// ingestor.state.RwMutex.Unlock()
		} else {
			// ingestor.state.RwMutex.Lock()
			// ingestor.state.RwMutex.Unlock()
		}
		ingestor.state.HeadNumber = number
	}
}