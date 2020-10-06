package client

import "github.com/figment-networks/avalanche-indexer/util"

// PlatformClient talks to the P chain
type PlatformClient struct {
	rpc
}

func (c PlatformClient) GetCurrentValidators() (resp *ValidatorsResponse, err error) {
	err = c.call("platform.getCurrentValidators", nil, &resp)
	return
}

func (c PlatformClient) GetPendingValidators() (resp *ValidatorsResponse, err error) {
	err = c.call("platform.getPendingValidators", nil, &resp)
	return
}

func (c PlatformClient) GetBalance(address string) (resp *Balance, err error) {
	err = c.call("platform.getBalance", map[string]string{"address": address}, &resp)
	return
}

func (c PlatformClient) GetMinStake() (resp *MinStakeResponse, err error) {
	err = c.call("platform.getMinStake", nil, &resp)
	return
}

func (c PlatformClient) GetBlockchains() (resp *BlockchainsResponse, err error) {
	err = c.call("platform.getBlockchains", nil, &resp)
	return
}

func (c PlatformClient) GetCurrentHeight() (int64, error) {
	result := map[string]string{}
	err := c.call("platform.getHeight", nil, &result)
	if err != nil {
		return 0, err
	}
	return util.ParseInt64(result["height"])
}
