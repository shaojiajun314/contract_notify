package parser

import (
	"fmt"
	"errors"
	"strings"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"contract_notify/collector/types"
)



func checkTopicIsExsited(t common.Hash, ts []common.Hash) bool {
	ret := false
	for _,  i := range ts {
		if t == i {
			ret = true
			break
		}
	}
	return ret
}


func getEventNameByABI(parsedAbi abi.ABI, name string) (string, *abi.Event, bool) {
	for n, e := range parsedAbi.Events {
		if e.Sig == name {
			return n, &e, true
		}
	}
	return "", nil, false
}


func getFilterAddressAndTopics(
	addressEventsMap types.AddressEventsMap,
) ([]common.Address, [][]common.Hash, AddressTopicsMap, map[string]map[string][]string, error) {
	filterContractAddresses := []common.Address{}
	filterContractTopics := []common.Hash{}
	filerAddressesTopicMap := make(AddressTopicsMap)
	eventParamsMap := make(map[string](map[string][]string))
	for a, eventsWithAbi := range addressEventsMap {
		eventParamsMap[a] = make(map[string]([]string))
		address := common.HexToAddress(a)
		filterContractAddresses = append(
			filterContractAddresses,
			address,
		)
		topics := []Topic{}
		parsedAbi, e := abi.JSON(strings.NewReader(eventsWithAbi.ABI))
		if e != nil {
			return filterContractAddresses, [][]common.Hash{}, filerAddressesTopicMap, eventParamsMap, e
		}
		for _, e := range eventsWithAbi.Events {
			topic := common.HexToHash(crypto.Keccak256Hash([]byte(e)).Hex())
			if !checkTopicIsExsited(topic, filterContractTopics) {
					filterContractTopics = append(
					filterContractTopics,
					topic,
				)
			}
			name, event, got := getEventNameByABI(parsedAbi, e)
			if !got {
				return filterContractAddresses, [][]common.Hash{}, filerAddressesTopicMap, eventParamsMap, errors.New(fmt.Sprintf(
						"event(%v) does not find in abi", e,
					),
				)
			}
			eventParamsMap[a][e] = []string{}
			for _, arg := range event.Inputs {
				eventParamsMap[a][topic.Hex()] = append(eventParamsMap[a][topic.Hex()], arg.Name)
			}
			topics = append(topics, Topic{
				hash: topic,
				eventName: name,
			})
			
		}
		boundConract := bind.NewBoundContract(
			address,
			parsedAbi,
			nil, nil, nil,
		)
		filerAddressesTopicMap[address] = TopicsConf{
			Topics: topics,
			BoundContract: boundConract,
		}
	}
	filterContractTopicsRet := [][]common.Hash{}
	for _, t := range filterContractTopics{
		filterContractTopicsRet = append(filterContractTopicsRet, []common.Hash{t})
	}
	return filterContractAddresses, filterContractTopicsRet, filerAddressesTopicMap, eventParamsMap, nil
}
