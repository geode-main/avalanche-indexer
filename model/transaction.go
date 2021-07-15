package model

import (
	"encoding/base64"
	"time"

	"github.com/figment-networks/avalanche-indexer/model/types"
	"github.com/figment-networks/avalanche-indexer/util"
)

type (
	Transaction struct {
		ID            string  `json:"id"`
		ReferenceTxID *string `json:"reference_tx_id,omitempty"`

		Chain       string    `json:"chain"`
		Type        string    `json:"type"`
		Block       *string   `json:"block,omitempty"`
		BlockHeight *uint64   `json:"block_height,omitempty"`
		Timestamp   time.Time `json:"timestamp"`
		Status      string    `json:"status,omitempty"`

		Memo             *string   `json:"memo,omitempty"`
		MemoText         *string   `json:"memo_text,omitempty"`
		Nonce            *uint64   `json:"nonce,omitempty"`
		Fee              uint64    `json:"fee"`
		SourceChain      *string   `json:"source_chain,omitempty"`
		DestinationChain *string   `json:"destination_chain,omitempty"`
		Metadata         types.Map `json:"metadata,omitempty" gorm:"type:text"`

		Inputs        []Output          `json:"inputs,omitempty" sql:"-" gorm:"-"`
		InputAmounts  map[string]uint64 `json:"input_amounts,omitempty" sql:"-" gorm:"-"`
		Outputs       []Output          `json:"outputs,omitempty" sql:"-" gorm:"-"`
		OutputAmounts map[string]uint64 `json:"output_amounts,omitempty" sql:"-" gorm:"-"`
	}

	TransactionTypeCount struct {
		Type       string `json:"type"`
		TotalCount int    `json:"total_count"`
	}

	Credential struct {
		Address   string `json:"address"`
		PublicKey string `json:"public_key"`
		Signature string `json:"signature"`
	}
)

func (Transaction) TableName() string {
	return "transactions"
}

func (tx *Transaction) SetReferenceID(value string) {
	tx.ReferenceTxID = &value
}

func (tx *Transaction) SetRawMemo(data []byte) {
	if data == nil {
		return
	}

	plain := util.TxMemo(data)
	encoded := base64.StdEncoding.EncodeToString(data)

	tx.Memo = &encoded
	tx.MemoText = &plain
}

func (tx *Transaction) UpdateAmounts() {
	if tx.InputAmounts == nil {
		tx.InputAmounts = map[string]uint64{}
	}
	if tx.OutputAmounts == nil {
		tx.OutputAmounts = map[string]uint64{}
	}

	for _, in := range tx.Inputs {
		tx.InputAmounts[in.Asset] += in.Amount
	}

	for _, out := range tx.Outputs {
		tx.OutputAmounts[out.Asset] += out.Amount
	}
}

func (tx *Transaction) UsesUTXOs() bool {
	return tx.Type != TxTypeEvm
}
