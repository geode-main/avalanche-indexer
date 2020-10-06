package client

type InfoClient struct {
	rpc
}

type infoPeersResponse struct {
	Peers []Peer `json:"peers"`
}

func (c InfoClient) BlockchainID(alias string) (string, error) {
	data, err := c.callRaw("info.getBlockchainID", map[string]string{"alias": alias})
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c InfoClient) NetworkID() (string, error) {
	data, err := c.callRaw("info.getNetworkID", nil)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c InfoClient) NetworkName() (string, error) {
	resp := map[string]string{}
	err := c.call("info.getNetworkName", nil, &resp)
	return resp["networkName"], err
}

func (c InfoClient) NodeID() (string, error) {
	data, err := c.callRaw("info.getNodeID", nil)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c InfoClient) NodeVersion() (string, error) {
	resp := map[string]string{}
	if err := c.call("info.getNodeVersion", nil, &resp); err != nil {
		return "", err
	}
	return resp["version"], nil
}

func (c InfoClient) Peers() ([]Peer, error) {
	resp := infoPeersResponse{}
	if err := c.call("info.peers", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Peers, nil
}

func (c InfoClient) TxFee() (*TxFeeResponse, error) {
	resp := &TxFeeResponse{}
	return resp, c.call("info.getTxFee", nil, resp)
}
