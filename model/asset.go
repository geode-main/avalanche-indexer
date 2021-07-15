package model

type Asset struct {
	ID                int    `json:"-"`
	AssetID           string `json:"id"`
	Type              string `json:"type"`
	Name              string `json:"name"`
	Symbol            string `json:"symbol"`
	Denomination      int    `json:"denomination"`
	TransactionsCount *int   `json:"transactions_count,omitempty" gorm:"-"`
}

func (Asset) TableName() string {
	return "assets"
}
