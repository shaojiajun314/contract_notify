package eventdata

import (
	"fmt"
	// "encoding/hex"

	"contract_notify/db"
	"contract_notify/types"
	"contract_notify/common"
)

func WriteHeadEventDataChainHeight(db db.KeyValueWriter, height uint64) error {
	if err := db.Put(HeadEventDataChainHeightKey, EncodeUint64(height)); err != nil {
		return fmt.Errorf("failed to store HeadEventDataChainHeight. err: %v", err)
	}
	return nil
}

func ReadHeadEventDataChainHeight(db db.Reader) (uint64, error) {
	if data, err := db.Get(HeadEventDataChainHeightKey); err != nil {
		return 0, fmt.Errorf("failed to get ReadHeadEventDataChainHeight. err: %v", err)
	} else {
		return DecodeUint64(data), nil
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////
//             block height -> data list
//////////////////////////////////////////////////////////////////////////////////////////////////////////////

func saveEventTaskData(db db.Database, height uint64, eventTaskDatas types.EventTaskDatas) error {
	var e error
	var dataBytes []byte
	if dataBytes, e = eventTaskDatas.Marshal(); e != nil {
		return fmt.Errorf("failed to Marshal eventTaskDatas, err: %v", e)
	}
	fmt.Println("this 1")
	if e := db.Put(blockEventTaskDataKey(height), dataBytes); e != nil {
		return fmt.Errorf("failed to store eventTaskDatas, err: %v", e)
	}
	return nil
}

func getEventTaskData(db db.Database, height uint64) (*types.EventTaskDatas, error) {
	var e error
	var dataBytes []byte

	fmt.Println("this 1")
	if dataBytes, e = db.Get(blockEventTaskDataKey(height)); e != nil { // has
		return nil, fmt.Errorf("failed to getEventTaskData.Get, err: %v", e)
	}
	eventTaskDatas := &types.EventTaskDatas{}
	if e = eventTaskDatas.Unmarshal(dataBytes); e != nil {
		return nil, fmt.Errorf("failed to getEventTaskData.Unmarshal, err: %v", e)
	}
	return eventTaskDatas, nil
}


//////////////////////////////////////////////////////////////////////////////////////////////////////////////
//             address -> data list
//////////////////////////////////////////////////////////////////////////////////////////////////////////////
func getEventTaskDataLastIndexByAddress(db db.Database, taskAddress common.Address) (uint64, error) {
	key := addressEventTaskDataLastIndexKey(taskAddress)
	if has, e := db.Has(key); e != nil {
		return 0, fmt.Errorf("failed to getEventTaskDataLastIndexByAddress.hasAddress, err: %v", e)
	} else if !has {
		return 0, nil
	}
	if bs, e := db.Get(key); e != nil {
		return 0, fmt.Errorf("failed to getEventTaskDataLastIndexByAddress, err: %v", e)
	} else {
		return DecodeUint64(bs), nil
	}
}

func saveEventTaskDataLastIndexByAddress(db db.Database, taskAddress common.Address, index uint64) (error) {
	key := addressEventTaskDataLastIndexKey(taskAddress)
	if e := db.Put(key, EncodeUint64(index)); e != nil {
		return fmt.Errorf("failed to saveEventTaskDataLastIndexByAddress, err: %v", e)
	}
	return nil
}

func saveEventTaskDataByAddress(db db.Database, taskAddress common.Address, index uint64, height uint64) error {
	if e := db.Put(addressEventTaskDataWithIndexKey(taskAddress, index), EncodeUint64(height)); e != nil {
		return fmt.Errorf("failed to saveEventTaskDataByAddress, err: %v", e)
	}
	return nil
}


//////////////////////////////////////////////////////////////////////////////////////////////////////////////
//             address-k-v -> data list
//////////////////////////////////////////////////////////////////////////////////////////////////////////////

func getFilterkEventTaskDataByArgLastIndex(db db.Database, taskAddress common.Address, key string, value interface{}) (uint64, error) {
	dbkey := filterkEventTaskDataByArgLastIndexKey(taskAddress, key, value)
	if has, e := db.Has(dbkey); e != nil {
		return 0, fmt.Errorf("failed to getFilterkEventTaskDataByArgLastIndex.hasAddress, err: %v", e)
	} else if !has {
		return 0, nil
	}
	if bs, e := db.Get(dbkey); e != nil {
		return 0, fmt.Errorf("failed to getFilterkEventTaskDataByArgLastIndex, err: %v", e)
	} else {
		return DecodeUint64(bs), nil
	}
}

func saveFilterkEventTaskDataByArgLastIndex(db db.Database, taskAddress common.Address, key string, value interface{}, index uint64) (error) {
	if e := db.Put(filterkEventTaskDataByArgLastIndexKey(taskAddress, key, value), EncodeUint64(index)); e != nil {
		return fmt.Errorf("failed to saveFilterkEventTaskDataByArgLastIndex, err: %v", e)
	}
	return nil
}

func savefilterkEventTaskDataByArg(db db.Database, taskAddress common.Address, key string, value interface{}, index uint64, height uint64) error {
	if e := db.Put(filterkEventTaskDataByArgWithIndexKey(taskAddress, key, value, index), EncodeUint64(height)); e != nil {
		return fmt.Errorf("failed to savefilterkEventTaskDataByArg, err: %v", e)
	}
	return nil
}