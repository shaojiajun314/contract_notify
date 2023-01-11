package server

import (
	"math/big"
	"contract_notify/config"
	"contract_notify/server/api"
	"contract_notify/collector"
	collectorConf "contract_notify/collector/conf"
	"contract_notify/rpc"
)

type Node struct {
	httpServer  *rpc.HttpServer
	collector   *mepcollector.Manager
}

func New(conf *config.Config) (*Node, error) {
	httpServer, err := rpc.NewHttpServer("0.0.0.0:8000", false)
	if err != nil {
		return nil, err
	}

	// todo
	collector := mepcollector.NewManager()
	collector.RegisteIngestor(
		"5",
		"https://eth-goerli.g.alchemy.com/v2/5c8YbVzERLkTRMTYSotAYp07LR2cGbOk",
		10,
		7853766,
	)
	collectorConf.RegisteCollectIterBlockStep(new(big.Int).SetUint64(1000000))

	node := &Node{
		httpServer: httpServer,
		collector: collector,
	}

	node.registerAPI()

	return node, nil
}

func (n *Node) Start() {
	if err := n.httpServer.Start(); err != nil {
		panic(err)
	}
}


func (n *Node) registerAPI() {
	n.httpServer.RegisterApis([]rpc.API{
		{
			Namespace: "collector",
			Service:   api.NewCollectorAPI(n.collector),
		},
	})
}
