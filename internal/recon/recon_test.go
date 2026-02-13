package recon

import (
	"testing"
	"time"
	"transaction-recon/internal/model"
)

func TestReconcile(t *testing.T) {
	startDate, _ := time.Parse("2006-01-02", "2023-01-01")
	endDate, _ := time.Parse("2006-01-02", "2023-01-02")

	systemTrxs := []model.SystemTransaction{
		{TrxID: "S1", Amount: 100.0, TransactionTime: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)},
		{TrxID: "S2", Amount: -50.0, TransactionTime: time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC)},
		{TrxID: "S3", Amount: 200.0, TransactionTime: time.Date(2023, 1, 3, 11, 0, 0, 0, time.UTC)}, // Out of range
	}

	bankRecords := []model.BankStatementRecord{
		{UniqueID: "B1", Amount: 100.0, Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), BankName: "BankA"},
		{UniqueID: "B2", Amount: -60.0, Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), BankName: "BankA"}, // Mismatch amount
		{UniqueID: "B3", Amount: 300.0, Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), BankName: "BankB"}, // Missing in system
	}

	summary := Reconcile(systemTrxs, bankRecords, startDate, endDate)

	if summary.TotalProcessed != 5 { // 2 sys + 3 bank
		t.Errorf("Expected 5 processed, got %d", summary.TotalProcessed)
	}

	if summary.TotalMatched != 1 { // S1 matches B1
		t.Errorf("Expected 1 matched, got %d", summary.TotalMatched)
	}

	if len(summary.MissingInBank) != 1 { // S2's equivalent is missing
		t.Errorf("Expected 1 missing in bank, got %d", len(summary.MissingInBank))
	}

	if len(summary.MissingInSystem["BankA"]) != 1 { // B2
		t.Errorf("Expected 1 missing in system (BankA), got %d", len(summary.MissingInSystem["BankA"]))
	}

	if len(summary.MissingInSystem["BankB"]) != 1 { // B3
		t.Errorf("Expected 1 missing in system (BankB), got %d", len(summary.MissingInSystem["BankB"]))
	}
}

func TestReconcile_Duplicates(t *testing.T) {
	startDate, _ := time.Parse("2006-01-02", "2023-01-01")
	endDate, _ := time.Parse("2006-01-02", "2023-01-01")

	// Two identical system trxs
	systemTrxs := []model.SystemTransaction{
		{TrxID: "S1", Amount: 100.0, TransactionTime: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)},
		{TrxID: "S2", Amount: 100.0, TransactionTime: time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC)},
	}

	// Only one bank record
	bankRecords := []model.BankStatementRecord{
		{UniqueID: "B1", Amount: 100.0, Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), BankName: "BankA"},
	}

	summary := Reconcile(systemTrxs, bankRecords, startDate, endDate)

	if summary.TotalMatched != 1 {
		t.Errorf("Expected 1 balanced match, got %d", summary.TotalMatched)
	}

	if len(summary.MissingInBank) != 1 {
		t.Errorf("Expected 1 remaining in system, got %d", len(summary.MissingInBank))
	}
}
