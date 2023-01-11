package eventdata

import (
	"bytes"
	"encoding/json"
	"contract_notify/types"
	"contract_notify/common"

	"github.com/ethereum/go-ethereum/rlp"
)


func (mngr *EventTaskDataManager)EventTaskDataByHeight(height uint64) (*types.EventTaskDatas, error) {
	return getEventTaskData(mngr.db, height)
}

func (mngr *EventTaskDataManager)HeightSet(prefix []byte, offset uint64, limit int) (*map[uint64]bool, error) {
	iter := mngr.db.NewIterator(prefix, EncodeUint64(offset))
	defer iter.Release()
	heightSet := make(map[uint64]bool)
	for i:=(0); i < limit && iter.Next(); i++ {
		v := iter.Value()
		heightSet[DecodeUint64(v)] = true
	}
	return &heightSet, nil
}

func (mngr *EventTaskDataManager)FilterEventTaskDataByTaskAddress(address common.Address, offset uint64, limit int) ([]*types.EventTaskData, error) {
	prefix := addressEventTaskDataKey(address)
	heightSet, e := mngr.HeightSet(prefix, offset, limit)
	if e != nil {
		return nil, e
	}
	targetAddressBytes := address.Bytes()
	ret := []*types.EventTaskData{}
	for height, _ := range *heightSet {
		eventTaskDatas, e := mngr.EventTaskDataByHeight(height)
		if e != nil {
			return nil, e
		}
		for _, eventTaskData := range []*types.EventTaskData(*eventTaskDatas) {
			if bytes.Equal(eventTaskData.TaskHash.Bytes(), targetAddressBytes) {
				ret = append(ret, eventTaskData)
			}
		}
	}
	return ret, nil
}


func (mngr *EventTaskDataManager)FilterEventTaskDataByTaskAddressKV(address common.Address, k string, v interface{}, offset uint64, limit int) ([]*types.EventTaskData, error) {
	prefix := filterkEventTaskDataByArgKey(address, k, v) 
	heightSet, e := mngr.HeightSet(prefix, offset, limit)
	if e != nil {
		return nil, e
	}
	targetAddressBytes := address.Bytes()
	bf := new(bytes.Buffer)
	rlp.Encode(bf, v)
	targetValue := bf.Bytes()
	ret := []*types.EventTaskData{}
	for height, _ := range *heightSet {
		eventTaskDatas, e := mngr.EventTaskDataByHeight(height)
		if e != nil {
			return nil, e
		}
		for _, eventTaskData := range []*types.EventTaskData(*eventTaskDatas) {
			if bytes.Equal(eventTaskData.TaskHash.Bytes(), targetAddressBytes) {
				var kv map[string]interface{}
				if e = json.Unmarshal(eventTaskData.Data.ParsedEvent, &kv); e != nil {
					return nil, e
				}
				value, has := kv[k]
				if !has {
					continue
				}
				bf := new(bytes.Buffer)
				rlp.Encode(bf, value)
				if !bytes.Equal(bf.Bytes(), targetValue) {
					continue
				}
				ret = append(ret, eventTaskData)
			}
		}
		
	}
	return ret, nil
}