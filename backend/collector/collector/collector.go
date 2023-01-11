package collector

import (
	"time"
	"fmt"
	"math/big"

	"contract_notify/collector/conf"
	"contract_notify/collector/types"
	"contract_notify/collector/chain"
)


type CollectorInstance struct {
	scanner 			*LogScanner
	cachedLogs 			*CachedLogs
	ingestorState 		*types.IngestorState
	store  				types.BlockStoreDal
	scanningBlocknum 	*big.Int
	parser 				types.Parser
	Network 			string
}


func transferFilersToArgs(
	filters []types.FilterQuery,
) (map[string][]string, types.AddressEventsMap) {
	addressEventsMap := make(map[string][]string)
	scannerAddressEventsMap := make(types.AddressEventsMap)
	for _, fq := range filters {
		if _, has := addressEventsMap[fq.Address]; !has {
			addressEventsMap[fq.Address] = fq.Events
			scannerAddressEventsMap[fq.Address] = &types.EventsConf {
				Events: fq.Events,
				ABI: fq.ABI,
			}
		} else {
			for _, e := range fq.Events {
				if !arrayContains(e, addressEventsMap[fq.Address]) {
					addressEventsMap[fq.Address] = append(
						addressEventsMap[fq.Address],
						e,
					)
					scannerAddressEventsMap[fq.Address].Events = append(
						scannerAddressEventsMap[fq.Address].Events,
						e,
					)
				}
			}
		}
		
	}
	return addressEventsMap, scannerAddressEventsMap
}


func NewCollectorInstance(
	network string,
	rpcurl string,
	filters []types.FilterQuery,
	startBlocknum uint64,
	ingestorState *types.IngestorState,
) (*CollectorInstance, error) {
	startBlocknumBignum := new(big.Int).SetUint64(startBlocknum)
	_, scannerAddressEventsMap := transferFilersToArgs(filters)
	component, _ := chain.DynChainMap[network]
	parser, e := component.ParserConstructor(scannerAddressEventsMap)
	if e != nil {
		return nil, e
	}
	scanner, e:= NeLogScanner(
		network,
		rpcurl,
		parser,
		startBlocknumBignum,
	)
	if e != nil {
		return nil, e
	}
	return &CollectorInstance {
		scanner: scanner,
		ingestorState: ingestorState,
		store: component.StoreDalConstructor(network, parser),
		scanningBlocknum: startBlocknumBignum,
		parser: parser,
		Network: network,
	}, nil
}

func (cI *CollectorInstance)GetEventParamsMap() map[string](map[string][]string) {
	return cI.parser.GetEventParamsMap()
}

func (cI *CollectorInstance) nextStep() (*types.Log, error) {
	for {
		if element, e := cI.cachedLogs.Next(); e == nil {
			element.ChainID = cI.Network
			return element, nil
		}
		headNumber := cI.ingestorState.HeadNumber
		tailNumber := cI.ingestorState.TailNumber
		if headNumber.Cmp(cI.scanningBlocknum) <= 0 || headNumber.Cmp(new(big.Int).SetUint64(0)) == 0 {
			time.Sleep(time.Duration(10)*time.Second)
			continue
		}
		fmt.Println(cI.scanningBlocknum, conf.CollectIterBlockStep, "cI.scanningBlocknum, conf.CollectIterBlockStep")
		loadToBlocknumber := new(big.Int).Add(cI.scanningBlocknum, conf.CollectIterBlockStep)
		if cI.scanningBlocknum.Cmp(tailNumber) >= 0 {

			if loadToBlocknumber.Cmp(headNumber) > 0 {
				loadToBlocknumber = headNumber
			}
			lagestNum, logs, er := cI.loadDataFromDal(
				loadToBlocknumber,
			)
			if er == nil {
				cI.scanningBlocknum = new(big.Int).Add(lagestNum, new(big.Int).SetUint64(1))
				cI.cachedLogs = NewCachedLogs(logs)
				continue
			} else {
			}
		}
		if loadToBlocknumber.Cmp(tailNumber) > 0 {
			loadToBlocknumber = tailNumber
		}
		logs, err := cI.scanner.NextLogs(
			cI.scanningBlocknum,
			loadToBlocknumber,
		)
		if err != nil {
			return nil, err
		}
		cI.scanningBlocknum = new(big.Int).Add(loadToBlocknumber, new(big.Int).SetUint64(1))
		cI.cachedLogs = NewCachedLogs(logs)
	}
}


func (cI *CollectorInstance) Next() (*types.Log, error) {
	var log *types.Log
	var e error
	if log, e = cI.nextStep(); e != nil {
		return nil, e
	}
	return log, nil
}


func (cI *CollectorInstance) loadDataFromDal(toNumber *big.Int) (*big.Int, *[]types.Log, error) {
	blocknumLogsMap, e := cI.store.Fetch(
		cI.scanningBlocknum,
		toNumber,
	)
	if e != nil {
		return nil, nil, e
	}
	logsRet := []types.Log{}
	lagestNum := new(big.Int).SetUint64(0)
	for num, logs := range blocknumLogsMap {
		if num.Cmp(lagestNum) > 0 {
			lagestNum = num;
		}
		logsRet = append(logsRet, *logs...)
	}
	return lagestNum, &logsRet, nil
}



func arrayContains(e string, array []string) bool {
	ret := false
	for _,  i := range array {
		if e == i {
			ret = true
			break
		}
	}
	return ret
}