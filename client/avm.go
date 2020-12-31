package client

// AvmClient talks to the X chain
type AvmClient struct {
	rpc
}

// GetAllBalances returns all assets balances for a given address
func (c AvmClient) GetAllBalances(address string) (result GetAllBalancesResponse, err error) {
	err = c.rpc.call(
		"avm.getAllBalances",
		map[string]string{"address": address},
		&result,
	)
	return
}
