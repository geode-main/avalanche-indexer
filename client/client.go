package client

import (
	"fmt"

	"github.com/ava-labs/coreth/ethclient"
	"github.com/ethereum/go-ethereum/eth/tracers"
)

type Client struct {
	Avm      AvmClient
	Evm      EvmClient
	Platform PlatformClient
	Info     InfoClient
	Ipc      IpcClient
	Index    IndexClient
}

func New(endpoint string) *Client {
	return &Client{
		Avm:      NewAvmClient(endpoint),
		Evm:      NewEVmClient(endpoint),
		Platform: NewPlatformClient(endpoint),
		Info:     NewInfoClient(endpoint),
		Ipc:      NewIpcClient(endpoint),
		Index:    NewIndexClient(endpoint),
	}
}

func NewAvmClient(endpoint string) AvmClient {
	return AvmClient{initRPC(endpoint, prefixAvm)}
}

func NewPlatformClient(endpoint string) PlatformClient {
	return PlatformClient{initRPC(endpoint, prefixPlatform)}
}

func NewIpcClient(endpoint string) IpcClient {
	return IpcClient{rpc: initRPC(endpoint, prefixIpc)}
}

func NewInfoClient(endpoint string) InfoClient {
	return InfoClient{initRPC(endpoint, prefixInfo)}
}

func NewIndexClient(endpoint string) IndexClient {
	return IndexClient{initRPC(endpoint, prefixIndex)}
}

func NewEVmClient(endpoint string) EvmClient {
	c, _ := ethclient.Dial(fmt.Sprintf("%s%s", endpoint, prefixEvm))

	return EvmClient{
		Client: c,
		rpc:    initRPC(endpoint, prefixEvm),
		traceConfig: &tracers.TraceConfig{
			Timeout: &tracerTimeout,
			Tracer:  &jsTracer,
		},
	}
}
