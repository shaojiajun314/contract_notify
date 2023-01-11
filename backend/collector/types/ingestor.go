package types
import (
	"math/big"

)


type Ingestor interface {
	Run() error
	// GetStore() BlockStoreDal
}


type IngestorState struct {
	HeadNumber 		*big.Int
	TailNumber 		*big.Int
	// RwMutex 		*sync.RWMutex
}


type NewIngestor func(
	network string,
	rpcUrl string,
	cachedBlockLength uint64,
	startBlocknum uint64,
	IngestorState *IngestorState,
) (*Ingestor, error)

