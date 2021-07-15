package model

type Chain struct {
	ID      int    `json:"-"`
	ChainID string `json:"id"`
	Name    string `json:"name"`
	VM      string `json:"vm,omitempty"`
	Subnet  string `json:"subnet,omitempty"`
	Network uint32 `json:"network,omitempty"`
}

func (Chain) TableName() string {
	return "chains"
}
