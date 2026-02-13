package reader

import (
	"os"
	"testing"
)

func TestReadSystemTransactions(t *testing.T) {
	content := `trxID,amount,type,transactionTime
S1,100.50,CREDIT,2023-01-01T10:00:00Z
S2,50.00,DEBIT,2023-01-01T11:00:00Z
`
	tmpfile, _ := os.CreateTemp("", "system*.csv")
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte(content))
	tmpfile.Close()

	trxs, err := ReadSystemTransactions(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to read system transactions: %v", err)
	}

	if len(trxs) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(trxs))
	}

	if trxs[0].Amount != 100.50 {
		t.Errorf("Expected CREDIT 100.50 to be 100.50, got %f", trxs[0].Amount)
	}

	if trxs[1].Amount != -50.00 {
		t.Errorf("Expected DEBIT 50.00 to be -50.00, got %f", trxs[1].Amount)
	}
}

func TestReadBankStatement(t *testing.T) {
	content := `unique_identifier,amount,date
B1,100.50,2023-01-01
B2,-50.00,2023-01-01
`
	tmpfile, _ := os.CreateTemp("", "bank*.csv")
	defer os.Remove(tmpfile.Name())
	tmpfile.Write([]byte(content))
	tmpfile.Close()

	records, err := ReadBankStatement(tmpfile.Name(), "MyBank")
	if err != nil {
		t.Fatalf("Failed to read bank statement: %v", err)
	}

	if len(records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(records))
	}

	if records[0].Amount != 100.50 {
		t.Errorf("Expected 100.50, got %f", records[0].Amount)
	}

	if records[1].Amount != -50.00 {
		t.Errorf("Expected -50.00, got %f", records[1].Amount)
	}
}
