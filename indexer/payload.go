package indexer

import (
	"time"

	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/util"
	"github.com/figment-networks/indexing-engine/pipeline"
)

type PayloadFactory struct {
}

type Payload struct {
	ID        string
	SyncTime  time.Time
	Processed bool

	// Fetched properties
	NodeVersion       string
	NetworkName       string
	Height            int64
	Peers             []client.Peer
	Blockchains       []client.Blockchain
	CurrentValidators []client.Validator
	CurrentDelegators []client.Delegator
	PendingValidators []client.Validator
	PendingDelegators []client.Delegator
	MinStake          *client.MinStakeResponse
	RawTxFee          *client.TxFeeResponse
	Balances          map[string]*client.Balance

	// Calculated properties
	ActiveStakeAmount     int64
	ActiveValidatorShare  map[string]float64
	PendingStakeAmount    int64
	PendingValidatorShare map[string]float64
	MinValidatorStake     int64
	MinDelegatorStake     int64
	TxFee                 int64
	CreationTxFee         int64

	// Mapped properties
	NetworkMetric *model.NetworkMetric
	Validators    []model.Validator
	ValidatorSeq  []model.ValidatorSeq
	Delegations   []model.Delegation
}

func NewPayload() *Payload {
	return &Payload{
		ID:       util.UUID(),
		SyncTime: time.Now(),
	}
}

func (p *Payload) MarkAsProcessed() {
	p.Processed = true
}

func NewPayloadFactory() PayloadFactory {
	return PayloadFactory{}
}

func (p PayloadFactory) GetPayload(int64) pipeline.Payload {
	return NewPayload()
}
