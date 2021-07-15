package model

type TransactionInput struct {
	ID   string
	TxID string
}

func (TransactionInput) TableName() string {
	return "transaction_inputs"
}
