package client

type AvmBalance struct {
	Asset   string `json:"asset"`
	Balance string `json:"balance"`
}

type GetAllBalancesResponse struct {
	Balances []AvmBalance `json:"balances"`
}
