package client

import (
	"context"

	"github.com/ava-labs/coreth/ethclient"
	"github.com/ethereum/go-ethereum/eth/tracers"
)

var (
	tracerTimeout = "180s"
)

type EvmClient struct {
	*ethclient.Client
	rpc         rpc
	traceConfig *tracers.TraceConfig
}

func (c *EvmClient) TraceTransaction(ctx context.Context, hash string) (*Call, error) {
	result := &Call{}
	args := []interface{}{hash, c.traceConfig}

	err := c.rpc.call("debug_traceTransaction", args, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
