package reader

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"time"
	"transaction-recon/internal/model"
)

func ReadSystemTransactions(filePath string) ([]model.SystemTransaction, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)

	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	headerMap := make(map[string]int)
	for i, h := range header {
		headerMap[h] = i
	}

	var transactions []model.SystemTransaction
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		amount, _ := strconv.ParseFloat(record[headerMap["amount"]], 64)
		trxType := model.TransactionType(record[headerMap["type"]])

		normalizedAmount := amount
		if trxType == model.TypeDebit {
			normalizedAmount = -amount
		}

		// Assume the format is 2023-01-01T09:00:00Z
		t, err := time.Parse(time.RFC3339, record[headerMap["transactionTime"]])
		if err != nil {
			continue
		}

		transactions = append(transactions, model.SystemTransaction{
			TrxID:           record[headerMap["trxID"]],
			Amount:          normalizedAmount,
			Type:            trxType,
			TransactionTime: t,
		})
	}

	return transactions, nil
}

func ReadBankStatement(filePath string, bankName string) ([]model.BankStatementRecord, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	headerMap := make(map[string]int)
	for i, h := range header {
		headerMap[h] = i
	}

	var records []model.BankStatementRecord
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		amount, _ := strconv.ParseFloat(record[headerMap["amount"]], 64)
		date, _ := time.Parse("2006-01-02", record[headerMap["date"]])

		records = append(records, model.BankStatementRecord{
			UniqueID: record[headerMap["unique_identifier"]],
			Amount:   amount,
			Date:     date,
			BankName: bankName,
		})
	}

	return records, nil
}
