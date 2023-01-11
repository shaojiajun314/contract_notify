package eventdata

import (
	"fmt"
	"time"
	"encoding/json"

	"contract_notify/db"
	"contract_notify/types"
	"contract_notify/blockchain"

	"contract_notify/db/memorydb"
	"contract_notify/blockchain/chain_db"
)


type EventTaskDataManager struct {
	chain  *blockchain.BlockChain
	db     db.Database
}

func NewEventTaskDataManager(chain *blockchain.BlockChain, db db.Database) *EventTaskDataManager {
	
	// todo for test
	db = chain_db.NewDatabase(memorydb.New())

	if has, e := db.Has(HeadEventDataChainHeightKey); e != nil {
		panic(e)
	} else if !has {
		db.Put(HeadEventDataChainHeightKey, EncodeUint64(0))
	}

	return &EventTaskDataManager{
		db: db,
		chain: chain,
	}
}

func (mngr *EventTaskDataManager)Run() error {
	var e error
	var head uint64
	var currenHeight uint64
	var block *types.Block
	if head, e = ReadHeadEventDataChainHeight(mngr.db); e != nil {
		return e
	}
	for {
		if currenHeight <= head {
			currenHeight = mngr.chain.CurrentHeader().Height
			if currenHeight <= head {
				time.Sleep(time.Second * 10)
				continue
			}
		}
		head ++ // 0 is genesis block
		if block, e = mngr.chain.GetBlockByHeight(head); e != nil {
			return e
		}
		eventTaskDatas := block.Transactions.FilterEventTaskData()

		// warnning for test
		eventTaskDatas = test().FilterEventTaskData()
		
		if e = mngr.save(head, eventTaskDatas); e != nil {
			return e
		}
		if e = WriteHeadEventDataChainHeight(mngr.db, head); e != nil {
			return e
		}
	}
}

func (mngr *EventTaskDataManager)save(height uint64, eventTaskDatas types.EventTaskDatas) error {
	if e := mngr.saveEventTaskData(height, eventTaskDatas); e != nil {
		return fmt.Errorf("failed to saveEventTaskData. err: %v", e)
	}
	for _, item := range eventTaskDatas {
		if e := mngr.saveEventTaskDataLookupEntry(height, item); e != nil {
			return fmt.Errorf("failed to saveEventTaskDataLookupEntry. err: %v", e)
		}
	}
	return nil
}


//////////////////////////////////////////////////////////////////////////////////////////////////////////////
//             block height -> data list
//////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (mngr *EventTaskDataManager)saveEventTaskData(height uint64, eventTaskDatas types.EventTaskDatas) error {
	return saveEventTaskData(mngr.db, height, eventTaskDatas)
}

func (mngr *EventTaskDataManager)saveEventTaskDataLookupEntry(height uint64, eventTaskData *types.EventTaskData) error {
	// dataBytes, e := eventTaskData.Marshal()
	// if e != nil {
	// 	return fmt.Errorf("failed to savefilterkEventTaskDataByArg.Marshal, err: %v", e)
	// }

	lastIndex, e := getEventTaskDataLastIndexByAddress(mngr.db, eventTaskData.TaskHash)
	if e != nil {
		return e
	}
	lastIndex ++
	if e := saveEventTaskDataByAddress(mngr.db, eventTaskData.TaskHash, lastIndex, height);e != nil {
		return e
	}
	if e := saveEventTaskDataLastIndexByAddress(mngr.db, eventTaskData.TaskHash, lastIndex); e != nil {
		return e
	}

	var parsedEvent map[string]interface{}
	if len(eventTaskData.Data.ParsedEvent) == 0 {
		return nil
	}
	if e := json.Unmarshal(eventTaskData.Data.ParsedEvent, &parsedEvent); e!=nil { // todo
		return fmt.Errorf("ParsedEvent Unmarshal error: %v", e)
	}
	for k, v := range parsedEvent {

		lastIndex, e := getFilterkEventTaskDataByArgLastIndex(mngr.db, eventTaskData.TaskHash, k, v)
		if e != nil {
			return e
		}
		lastIndex ++ 
		if e := savefilterkEventTaskDataByArg(mngr.db, eventTaskData.TaskHash, k, v, lastIndex, height); e != nil {
			return e
		}
		if e := saveFilterkEventTaskDataByArgLastIndex(mngr.db, eventTaskData.TaskHash, k, v, lastIndex); e != nil {
			return e
		}
	}
	return nil
}
