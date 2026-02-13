package model

import "time"

type TransactionType string

const (
	TypeDebit  TransactionType = "DEBIT"
	TypeCredit TransactionType = "CREDIT"
)

type SystemTransaction struct {
	TrxID           string          `csv:"trxID"`
	Amount          float64         `csv:"amount"`
	Type            TransactionType `csv:"type"`
	TransactionTime time.Time       `csv:"transactionTime"`
}

type BankStatementRecord struct {
	UniqueID string    `csv:"unique_identifier"`
	Amount   float64   `csv:"amount"`
	Date     time.Time `csv:"date"`
	BankName string    `csv:"-"`
}

type ReconciliationSummary struct {
	TotalProcessed  int
	TotalMatched    int
	TotalUnmatched  int
	MissingInBank   []SystemTransaction
	MissingInSystem map[string][]BankStatementRecord // Grouped by bank
}
