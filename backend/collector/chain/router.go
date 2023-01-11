package chain

import (
	"contract_notify/collector/types"
	"contract_notify/collector/chain/ethereum"
	"contract_notify/collector/chain/ethereum/ingestor"
	"contract_notify/collector/chain/ethereum/dal"
	"contract_notify/collector/chain/ethereum/log_parser"
)


type ChainComponent struct {
	AdapterConstructor 		types.AdapterConstructor
	IngestorConstructor  	types.NewIngestor
	StoreDalConstructor  	types.NewBlockStoreDal
	ParserConstructor 		types.NewParser
}


var DynChainMap map[string]ChainComponent


func init() {
	DynChainMap = map[string]ChainComponent {
		"1": ChainComponent{
			AdapterConstructor: ethereum.NewAdapter,
			IngestorConstructor: ethingestor.NewIngestor,
			StoreDalConstructor: ethdal.NewBlockStoreDal,
			ParserConstructor: parser.NewParser,
		},
		"3": ChainComponent{
			AdapterConstructor: ethereum.NewAdapter,
			IngestorConstructor: ethingestor.NewIngestor,
			StoreDalConstructor: ethdal.NewBlockStoreDal,
			ParserConstructor: parser.NewParser,
		},
		"4": ChainComponent{
			AdapterConstructor: ethereum.NewAdapter,
			IngestorConstructor: ethingestor.NewIngestor,
			StoreDalConstructor: ethdal.NewBlockStoreDal,
			ParserConstructor: parser.NewParser,
		},
		"5": ChainComponent{
			AdapterConstructor: ethereum.NewAdapter,
			IngestorConstructor: ethingestor.NewIngestor,
			StoreDalConstructor: ethdal.NewBlockStoreDal,
			ParserConstructor: parser.NewParser,
		},
		"42": ChainComponent{
			AdapterConstructor: ethereum.NewAdapter,
			IngestorConstructor: ethingestor.NewIngestor,
			StoreDalConstructor: ethdal.NewBlockStoreDal,
			ParserConstructor: parser.NewParser,
		},
		"97": ChainComponent{
			AdapterConstructor: ethereum.NewAdapter,
			IngestorConstructor: ethingestor.NewIngestor,
			StoreDalConstructor: ethdal.NewBlockStoreDal,
			ParserConstructor: parser.NewParser,
		},
	}
}


