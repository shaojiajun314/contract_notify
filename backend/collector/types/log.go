package types

import (
    "contract_notify/common"
)


type Log struct {
    Address         string                      `json:"address"`
    Topics          []string                    `json:"topics"`
    Data            string                      `json:"data"`
    BlockHash       string                      `json:"blockHash"`
    BlockNumber     uint64                      `json:"blockNumber"`
    TransactionHash string                      `json:"transactionHash"`
    ParsedEvent     map[string]interface{}      `json:"parsedEvent"`
    TxIndex         uint                        `json:"tx_index"`
    LogIndex        uint                        `json:"log_index"`
    ChainID         string                      `json:"chain_id"`
}

type ThirdResponse struct {
    Sign        string          `json:"sign"`
    Passphrase  string          `json:"passphrase"`
    Address     common.Address  `json:"address"`
}

type SignedLog struct {
    Log         Log             `json:"log"`
    ThirdResp   *ThirdResponse  `json:"third_response"`
}


type EventsConf struct {
    Events  []string
    ABI     string
}


type AddressEventsMap map[string]*EventsConf


type Parser interface {
     ParseLog(logInterface interface{}) (*Log, bool, error)
     GetEventParamsMap() map[string]map[string][]string // address -> eventName -> []string params
}

type NewParser func(addressEventsMap AddressEventsMap) (Parser, error)