package mepcollector


import (
	"sync"
	"errors"
	"math/big"

	"contract_notify/collector/types"
	"contract_notify/collector/chain"
	"contract_notify/collector/collector"
)


type Manager struct {
	IngestorMap 				map[string]*types.Ingestor // network -> Ingestor instance
	SateMap	 					map[string]*types.IngestorState // network -> state
	IngestorMapRwMutex 			*sync.RWMutex
}

func NewManager() *Manager {
	return &Manager {
		IngestorMapRwMutex: new(sync.RWMutex),
		IngestorMap: make(map[string]*types.Ingestor),
		SateMap: make(map[string]*types.IngestorState),
	}
}


func (manager *Manager) newIngestor (
	network string,
	rpcUrl string,
	cachedBlockLength uint64,
	startBlocknum uint64,
) (*types.Ingestor, error) {
	component, has := chain.DynChainMap[network]
	if !has {
		panic("erroor network ingestor")
	}
	manager.IngestorMapRwMutex.Lock()
  	defer manager.IngestorMapRwMutex.Unlock()
  	if _, has :=  manager.IngestorMap[network]; has{
  		panic("duplicate newIngestor // todo")
  	}

	ingestorState := &types.IngestorState {
		HeadNumber: new(big.Int).SetUint64(0),
		TailNumber: new(big.Int).SetUint64(0),
	}

	instance, e := component.IngestorConstructor(
		network,
		rpcUrl,
		cachedBlockLength,
		startBlocknum,
		ingestorState,
	)
	if e != nil {
   		return nil, e
   	}
   manager.SateMap[network] = ingestorState
   manager.IngestorMap[network] = instance
   return instance, nil
}


func (manager *Manager) RegisteIngestor (
	network string,
	rpcUrl string,
	cachedBlockLength uint64,
	startBlocknum uint64,
) {
	ingestorInstance, e := manager.newIngestor(
		network,
		rpcUrl,
		cachedBlockLength,
		startBlocknum,
	)
	if e != nil {
		// manager.logger.Error(
		// 	"RegisteIngestor error",
		// 	e,
		// )
		panic(e)
	}
	go func(){
		e := (*ingestorInstance).Run()
		panic(e)
	}()
}


func (manager *Manager) NewCollector (
	network string,
	rpcurl string,
	filters []types.FilterQuery,
	startBlocknum uint64,
	// notify string,
	// isRPCNotify bool,
	// publickey []byte,
) (*collector.CollectorInstance, error) {
	manager.IngestorMapRwMutex.RLock()
  	defer manager.IngestorMapRwMutex.RUnlock()

	_, has :=  manager.IngestorMap[network]
	if !has {
		// manager.logger.Error(
		// 	"Ingestor is unregisted",
		// 	network,
		// )
		return nil, errors.New("Ingestor is unregisted")
	}
	state, _ := manager.SateMap[network]
	collector, e := collector.NewCollectorInstance(
		network,
		rpcurl,
		filters,
		startBlocknum,
		state,
		// notify,
		// isRPCNotify,
		// publickey,
	)
	if e != nil {
		// manager.logger.Error(
		// 	"error new collector",
		// 	network,
		// 	rpcurl,
		// 	startBlocknum,
		// 	notify,
		// 	isRPCNotify,
		// 	e,
		// )
		return nil, e
	}
	return collector, nil
}