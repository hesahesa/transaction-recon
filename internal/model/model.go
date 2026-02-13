package model

import "time"

type TransactionType string

const (
	TypeDebit  TransactionType = "DEBIT"
	TypeCredit TransactionType = "CREDIT"
)

type SystemTransaction struct {
	TrxID           string
	Amount          float64
	Type            TransactionType
	TransactionTime time.Time
}

type BankStatementRecord struct {
	UniqueID string
	Amount   float64
	Date     time.Time
	BankName string
}

type ReconciliationSummary struct {
	TotalProcessed  int
	TotalMatched    int
	TotalUnmatched  int
	MissingInBank   []SystemTransaction
	MissingInSystem map[string][]BankStatementRecord // Grouped by bank
}
