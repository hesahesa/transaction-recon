package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"transaction-recon/internal/model"
	"transaction-recon/internal/reader"
	"transaction-recon/internal/recon"
)

func main() {
	systemFile := flag.String("system", "", "Path to system transactions CSV")
	bankFiles := flag.String("bank", "", "Comma-separated paths to bank statement CSVs")
	startDateStr := flag.String("start", "", "Start date (YYYY-MM-DD)")
	endDateStr := flag.String("end", "", "End date (YYYY-MM-DD)")

	flag.Parse()

	if *systemFile == "" || *bankFiles == "" || *startDateStr == "" || *endDateStr == "" {
		fmt.Println("Usage: recon --system <path> --bank <path1,path2> --start <YYYY-MM-DD> --end <YYYY-MM-DD>")
		os.Exit(1)
	}

	startDate, err := time.Parse("2006-01-02", *startDateStr)
	if err != nil {
		log.Fatalf("Invalid start date: %v", err)
	}
	endDate, err := time.Parse("2006-01-02", *endDateStr)
	if err != nil {
		log.Fatalf("Invalid end date: %v", err)
	}

	// Read System Transactions
	systemTrxs, err := reader.ReadSystemTransactions(*systemFile)
	if err != nil {
		log.Fatalf("Error reading system file: %v", err)
	}

	// Read Bank Statements
	var allBankRecords []model.BankStatementRecord
	for _, path := range strings.Split(*bankFiles, ",") {
		path = strings.TrimSpace(path)
		bankName := filepath.Base(path)
		records, err := reader.ReadBankStatement(path, bankName)
		if err != nil {
			log.Fatalf("Error reading bank file %s: %v", bankName, err)
		}
		allBankRecords = append(allBankRecords, records...)
	}

	// Reconcile
	summary := recon.Reconcile(systemTrxs, allBankRecords, startDate, endDate)

	// Output Summary
	printSummary(summary)
}

func printSummary(s *model.ReconciliationSummary) {
	fmt.Println("========================================")
	fmt.Println(" RECONCILIATION SUMMARY")
	fmt.Println("========================================")
	fmt.Printf("Total Transactions Processed: %d\n", s.TotalProcessed)
	fmt.Printf("Total Matched Transactions:   %d\n", s.TotalMatched)
	fmt.Printf("Total Unmatched Transactions: %d\n", s.TotalUnmatched)
	fmt.Println("----------------------------------------")

	fmt.Println("Details of Unmatched Transactions:")

	if len(s.MissingInBank) > 0 {
		fmt.Printf("\n[System Transactions missing in Bank Statements (%d)]\n", len(s.MissingInBank))
		for _, trx := range s.MissingInBank {
			fmt.Printf("- ID: %s | Date: %s | Amount: %.2f | Type: %s\n",
				trx.TrxID, trx.TransactionTime.Format("2006-01-02"), trx.Amount, trx.Type)
		}
	}

	if len(s.MissingInSystem) > 0 {
		for bank, records := range s.MissingInSystem {
			fmt.Printf("\n[Bank Statement records missing in System - %s (%d)]\n", bank, len(records))
			for _, rec := range records {
				fmt.Printf("- ID: %s | Date: %s | Amount: %.2f\n",
					rec.UniqueID, rec.Date.Format("2006-01-02"), rec.Amount)
			}
		}
	}
}
