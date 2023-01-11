package collector

import (
    "errors"
    "contract_notify/collector/types"
)


type CachedLogs struct {
    logs    []types.Log
    length  int
    cursor  int
}

func NewCachedLogs(logs *[]types.Log) *CachedLogs {
    return &CachedLogs {
        logs: *logs,
        length: len(*logs),
        cursor: 0,
    }
}

func (cl *CachedLogs) Next() (*types.Log, error) {
    if cl == nil {
        return nil, errors.New("nil pointer")
    }
    if cl.cursor >= cl.length {
        return nil, errors.New("this cache has no more log that can be consumes")
    }
    element := cl.logs[cl.cursor]
    cl.cursor = cl.cursor + 1
    return &element, nil
}

