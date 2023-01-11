package api

import (
	"context"
	// "bytes"
	"io/ioutil"
	"net/http"


	// "contract_notify/types"
	// "contract_notify/common"
	// "contract_notify/blockchain"
	// hexutil "contract_notify/common/hex"
	"encoding/json"
	"fmt"
	"contract_notify/collector"
	"contract_notify/collector/types"
)

const PushDeerUrl string = "http://43.206.196.16:8800/message/push?pushkey=%v&text=%v"

type CollectorAPI struct {
	manager *mepcollector.Manager
}

func NewCollectorAPI(manager *mepcollector.Manager) *CollectorAPI {
	return &CollectorAPI{
		manager: manager,
	}
}

type CollectorForm struct {
	Network     	string  		`json:"network" mapstructure:"network"`
	ContractAddress string 			`json:"contract_address" mapstructure:"contract_address"`
	Abi          	string 			`json:"abi" mapstructure:"abi"`
	Events 			[]string 		`json:"events" mapstructure:"events"`
	RPC 			string			`json:"rpc" mapstructure:"rpc"`
	StartBlocknum 	uint64			`json:"start_blocknumber" mapstructure:"start_blocknumber"`
	PushKey 		string 			`json:"push_key" mapstructure:"start_blocknumber"`
}

func (c *CollectorAPI) Register(ctx context.Context, form *CollectorForm) error {
	collectorInstance, e := c.manager.NewCollector(
		form.Network,
		form.RPC,
  		[]types.FilterQuery {
  			types.FilterQuery {
  				Address: form.ContractAddress,
				Events: form.Events,
				ABI: form.Abi,
  			},
  		},
        form.StartBlocknum,
	)
	if e != nil {
		return e
	}
	go func() {
		for {
			fmt.Println(1, collectorInstance)
			log, err := collectorInstance.Next()
			fmt.Println(2)
			if err != nil {
				panic(err)
			}
			jsonData, _ := json.Marshal(&log)

			// jsonData1, _ := json.Marshal(&map[string]string{
			// 	// "pushkey": "PDU1TX2ISARYGIG7howSBIKYPVJpzSksbtRwR",
			// 	"pushkey": form.PushKey,
			// 	"text": string(jsonData),
			// })
			var request *http.Request
			var resp *http.Response

			if request, e = http.NewRequest(
				"GET",
				fmt.Sprintf(PushDeerUrl, form.PushKey, string(jsonData)),
				// bytes.NewReader(jsonData1),
				nil,
			); e != nil {
				panic(e)
			}
			if resp, e = http.DefaultClient.Do(request); e != nil {
				panic(e)
			}
			defer resp.Body.Close()
			var content []byte
			if content, e = ioutil.ReadAll(resp.Body); e!=nil{
				panic(e)
			}
			fmt.Println(string(content), "aaaaa")

			fmt.Println("event 1: ", string(jsonData))
		}
	}()
	return nil
}