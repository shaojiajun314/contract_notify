package eventdata

import (
	"fmt"
	"encoding/json"
	"encoding/binary"

	"contract_notify/common"
)

var (
	HeadEventDataChainHeightKey 				= []byte("HeadEventDataChainHeight")
	blockEventTaskDataPrefix   	 				= []byte("blockEventTaskData")

	addressEventTaskDataLastIndexPrefix  		= []byte("addressEventTaskDataLastIndex")
	addressEventTaskDataIndexPrefix   	 		= []byte("addressEventTaskData")

	filterkEventTaskDataByArgPrefix  			= []byte("filterkEventTaskDataByArg")
	filterkEventTaskDataByArgLastIndexPrefix 	= []byte("filterkEventTaskDataByArgLastIndex")
)


func EncodeUint64(data uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, data)
	return enc
}

func DecodeUint64(data []byte) uint64 {
	return binary.BigEndian.Uint64(data)
}


func blockEventTaskDataKey(height uint64) []byte {
	return append(blockEventTaskDataPrefix, EncodeUint64(height)...)
}


//////////////////////////////////////////////////////////////////////////////////////////////////////////////
//             address -> data list
//////////////////////////////////////////////////////////////////////////////////////////////////////////////

func addressEventTaskDataLastIndexKey(taskAddress common.Address) []byte {
	return append(addressEventTaskDataLastIndexPrefix, taskAddress.Bytes()...)
}

func addressEventTaskDataKey(taskAddress common.Address) []byte {
	return append(addressEventTaskDataIndexPrefix, taskAddress.Bytes()...)
}


func addressEventTaskDataWithIndexKey(taskAddress common.Address, index uint64) []byte {
	return append(addressEventTaskDataKey(taskAddress), EncodeUint64(index)...)
}


//////////////////////////////////////////////////////////////////////////////////////////////////////////////
//             address-k-v -> data list
//////////////////////////////////////////////////////////////////////////////////////////////////////////////

func filterkEventTaskDataByArgLastIndexKey(taskAddress common.Address, k string, v interface{}) []byte {
	bs, e := json.Marshal(v) // todo
	if e != nil {
		panic(fmt.Sprintf("error filterkEventTaskDataByArgLastIndexKey Marshal err: %v", e))
	}
	return append(filterkEventTaskDataByArgLastIndexPrefix, append(append(taskAddress.Bytes(), []byte(k)...), bs...)...)
}

func filterkEventTaskDataByArgKey(taskAddress common.Address, k string, v interface{}) []byte {
	bs, e := json.Marshal(v) // todo
	if e != nil {
		panic(fmt.Sprintf("error filterkEventTaskDataByArgWithIndexKey Marshal err: %v", e))
	}
	return append(filterkEventTaskDataByArgPrefix, append(append(taskAddress.Bytes(), []byte(k)...), bs...)...)
}

func filterkEventTaskDataByArgWithIndexKey(taskAddress common.Address, k string, v interface{}, index uint64) []byte {
	return append(filterkEventTaskDataByArgKey(taskAddress, k, v), EncodeUint64(index)...)
}
