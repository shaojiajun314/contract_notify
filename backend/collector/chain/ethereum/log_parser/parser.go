package parser

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"contract_notify/collector/types"
)


type Topic struct {
	hash 		common.Hash
	eventName	string
}

type TopicsConf struct {
	BoundContract 	*bind.BoundContract
	Topics 			[]Topic
}

type AddressTopicsMap map[common.Address]TopicsConf


func (tc TopicsConf) parseLog(log ethtypes.Log) (map[string]interface{}, bool, error) {
	bc := tc.BoundContract
	topics := tc.Topics
	received := make(map[string]interface{})
	if len(log.Topics) == 0 {
		return received, false, nil
	}
	topicHash := log.Topics[0]
	for _, topic := range topics {
		if topic.hash != topicHash {
			continue
		}
		if err := bc.UnpackLogIntoMap(received, topic.eventName, log); err != nil {
			// panic("todo")
			return received, false, errors.New("error abi file")
		}
		return received, true, nil
	}
	return received, false, nil
}


type Parser struct {
	FilterContractAddresses []common.Address
	FilterContractTopics 	[][]common.Hash
	FilerAddressesTopicMap 	AddressTopicsMap
	EventParamsMap  		map[string]map[string][]string
}


func (p Parser) ParseLog(logInterface interface{}) (*types.Log, bool, error) {
	log := logInterface.(ethtypes.Log)
	topicsConf := p.FilerAddressesTopicMap[log.Address]
	if parsedEvent, hit, e := topicsConf.parseLog(log); hit {
		return &types.Log {
			Address: log.Address.Hex(),
			Topics: func(ts []common.Hash) []string {
				ret := []string{}
				for _, t := range log.Topics {
					ret = append(
						ret,
						t.Hex(),
					)
				}
				return ret
			}(log.Topics),
			Data: hexutil.Encode(log.Data),
			BlockHash: log.BlockHash.Hex(),
			BlockNumber: log.BlockNumber,
			TransactionHash: log.TxHash.Hex(),
			TxIndex: log.TxIndex,
			LogIndex: log.Index,
			ParsedEvent: parsedEvent,
		}, true, nil
	} else if (e != nil){
		return nil, false, e
	} else {
		return nil, false, nil
	}
}


func (p Parser) GetEventParamsMap() map[string]map[string][]string {
	return p.EventParamsMap
}




func NewParser(addressEventsMap types.AddressEventsMap) (types.Parser, error) {
	var p types.Parser
	addresses, topics, addressesTopicMap, eventParamsMap, e := getFilterAddressAndTopics(addressEventsMap)
	if e != nil {
		return p, e
	}
	p = Parser {
		FilterContractAddresses: addresses,
		FilterContractTopics: topics,
		FilerAddressesTopicMap: addressesTopicMap,
		EventParamsMap: eventParamsMap,
	}
	return p, nil
}



