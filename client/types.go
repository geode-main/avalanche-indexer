package client

import "time"

type Peer struct {
	ID           string `json:"nodeID"`
	IP           string `json:"ip"`
	PublicIP     string `json:"publicIP"`
	Version      string `json:"version"`
	LastSent     string `json:"lastSent"`
	LastReceived string `json:"lastReceived"`
}

func (p Peer) Metadata() map[string]interface{} {
	return map[string]interface{}{
		"ip":            p.IP,
		"public_ip":     p.PublicIP,
		"version":       p.Version,
		"last_sent":     p.LastSent,
		"last_received": p.LastReceived,
	}
}

type BlockchainsResponse struct {
	Blockchains []Blockchain `json:"blockchains"`
}

type Blockchain struct {
	ID       string `json:"id"`
	SubnetID string `json:"subnetID"`
	Name     string `json:"name"`
	VMID     string `json:"vmId"`
}

type ValidatorsResponse struct {
	Validators []Validator `json:"validators"`
	Delegators []Delegator `json:"delegators"`
}

type MinStakeResponse struct {
	MinValidatorStake string `json:"minValidatorStake"`
	MinDelegatorStake string `json:"minDelegatorStake"`
}

type Validator struct {
	TxID            string      `json:"txID"`
	StartTime       string      `json:"startTime"`
	EndTime         string      `json:"endTime"`
	StakeAmount     string      `json:"stakeAmount"`
	NodeID          string      `json:"nodeID"`
	RewardOwner     RewardOwner `json:"rewardOwner"`
	PotentialReward string      `json:"potentialReward"`
	DelegationFee   string      `json:"delegationFee"`
	Uptime          string      `json:"uptime"`
	Connected       bool        `json:"connected"`
	Delegators      []Delegator `json:"delegators"`
}

type Delegator struct {
	TxID            string      `json:"txID"`
	StartTime       string      `json:"startTime"`
	EndTime         string      `json:"endTime"`
	StakeAmount     string      `json:"stakeAmount"`
	NodeID          string      `json:"nodeID"`
	RewardOwner     RewardOwner `json:"rewardOwner"`
	PotentialReward string      `json:"potentialReward"`
	DelegationFee   string      `json:"delegationFee"`
	Uptime          string      `json:"uptime"`
	Connected       bool        `json:"connected"`
}

type RewardOwner struct {
	Locktime  string   `json:"locktime"`
	Threshold string   `json:"threshold"`
	Addresses []string `json:"addresses"`
}

type TxFeeResponse struct {
	CreationTxFee string `json:"creationTxFee"`
	TxFee         string `json:"txFee"`
}

type Balance struct {
	Balance            string `json:"balance"`
	Unlocked           string `json:"unlocked"`
	LockedStakeable    string `json:"lockedStakeable"`
	LockedNotStakeable string `json:"lockedNotStakeable"`
	Staked             string `json:"staked"`
}

type RewardUTXOsResponse struct {
	NumFetched string   `json:"numFetched"`
	Encoding   string   `json:"encoding"`
	UTXOs      []string `json:"utxos"`
}

type Container struct {
	ID        string    `json:"id"`
	Index     string    `json:"index"`
	Bytes     string    `json:"bytes"`
	Encoding  string    `json:"encoding"`
	Timestamp time.Time `json:"timestamp"`
}

type ContainersResponse struct {
	Containers []Container `json:"containers"`
}

type AvmBalance struct {
	Asset   string `json:"asset"`
	Balance string `json:"balance"`
}

type GetAllBalancesResponse struct {
	Balances []AvmBalance `json:"balances"`
}

type GetStakeResponse struct {
	Staked string `json:"staked"`
}
