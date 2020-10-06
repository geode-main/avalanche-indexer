package client

type Client struct {
	Avm      AvmClient
	Platform PlatformClient
	Info     InfoClient
	Ipc      IpcClient
}

func New(endpoint string) *Client {
	return &Client{
		Avm:      NewAvmClient(endpoint),
		Platform: NewPlatformClient(endpoint),
		Info:     NewInfoClient(endpoint),
		Ipc:      NewIpcClient(endpoint),
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
