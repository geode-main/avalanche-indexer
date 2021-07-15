package client

import "strconv"

type InfoClient struct {
	rpc
}

type infoNodeVersionResponse struct {
	Version         string `json:"version"`
	DatabaseVersion string `json:"databaseVersion"`
	GitCommit       string `json:"gitCommit"`
	VmVersions      struct {
		AVM      string `json:"avm"`
		EVM      string `json:"evm"`
		Platform string `json:"platform"`
	} `json:"vmVersions"`
}

type infoPeersResponse struct {
	Peers []Peer `json:"peers"`
}

func (c InfoClient) BlockchainID(alias string) (string, error) {
	resp := map[string]string{}

	err := c.call("info.getBlockchainID", map[string]string{"alias": alias}, &resp)
	if err != nil {
		return "", err
	}

	return resp["blockchainID"], nil
}

func (c InfoClient) NetworkID() (int, error) {
	resp := map[string]string{}

	err := c.call("info.getNetworkID", nil, &resp)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(resp["networkID"])
}

func (c InfoClient) NetworkName() (string, error) {
	resp := map[string]string{}
	err := c.call("info.getNetworkName", nil, &resp)

	return resp["networkName"], err
}

func (c InfoClient) NodeID() (string, error) {
	data, err := c.callRaw(c.endpoint, "info.getNodeID", nil)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c InfoClient) NodeVersion() (string, error) {
	resp := infoNodeVersionResponse{}

	if err := c.call("info.getNodeVersion", nil, &resp); err != nil {
		return "", err
	}
	return resp.Version, nil
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
