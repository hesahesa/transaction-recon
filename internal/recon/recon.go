package recon

import (
	"time"
	"transaction-recon/internal/model"
)

type ReconKey struct {
	Date   time.Time
	Amount float64
}

func Reconcile(
	systemTrxs []model.SystemTransaction,
	bankStatements []model.BankStatementRecord,
	startDate, endDate time.Time,
) *model.ReconciliationSummary {

	summary := &model.ReconciliationSummary{
		MissingInSystem: make(map[string][]model.BankStatementRecord),
	}

	systemMap := make(map[ReconKey][]model.SystemTransaction)

	for _, trx := range systemTrxs {
		if isWithinRange(trx.TransactionTime, startDate, endDate) {
			key := ReconKey{
				Date:   truncateToDate(trx.TransactionTime),
				Amount: trx.Amount,
			}
			systemMap[key] = append(systemMap[key], trx)
			summary.TotalProcessed++
		}
	}

	for _, bankTrx := range bankStatements {
		if isWithinRange(bankTrx.Date, startDate, endDate) {
			summary.TotalProcessed++
			key := ReconKey{
				Date:   truncateToDate(bankTrx.Date),
				Amount: bankTrx.Amount,
			}

			trxs, ok := systemMap[key]
			if ok && len(trxs) > 0 {
				summary.TotalMatched++
				systemMap[key] = trxs[1:]
			} else {
				summary.TotalUnmatched++
				summary.MissingInSystem[bankTrx.BankName] = append(summary.MissingInSystem[bankTrx.BankName], bankTrx)
			}
		}
	}

	for _, trxs := range systemMap {
		for _, trx := range trxs {
			summary.TotalUnmatched++
			summary.MissingInBank = append(summary.MissingInBank, trx)
		}
	}

	return summary
}

// start and end are inclusive
func isWithinRange(t, start, end time.Time) bool {
	target := truncateToDate(t)
	return !target.Before(start) && !target.After(end)
}

func truncateToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
