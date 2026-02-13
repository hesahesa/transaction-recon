# Transaction Reconciliation Service

A Go-based command-line tool designed to identify unmatched transactions between internal system records and external bank statements.

## Features
- matches transactions based on **Date** and **Amount**.
- Handles multiple bank statement files simultaneously.
- Identifies system transactions missing from bank records.
- Identifies bank records missing from the system, grouped by source bank file.
- Filters reconciliation by a specified date range.

## Assumptions
- The date in the bank statement csv and in the CLI parameters are assumed in UTC time zone for simplicity.
- The currency in the csvs are the same
- Because this tool do the matching by Date and Amount, the total discrepancies (sum of absolute differences in amount between matched transactions) will be 0

## Tech Stack
- **Language**: Go (Golang)
- **Data Format**: CSV
- **Dependencies**: Standard library only (no external database or services required).

## Installation
Ensure you have Go installed on your system.

```bash
git clone <repository-url>
cd transaction-recon
go mod download
```

## Usage
Run the service using the following command:

```bash
go run main.go \
  --system <path_to_system_csv> \
  --bank <comma_separated_bank_csv_paths> \
  --start <YYYY-MM-DD> \
  --end <YYYY-MM-DD>
```

### Example
```bash
go run main.go \
  --system samples/system.csv \
  --bank samples/bank_a.csv \
  --start 2023-01-01 \
  --end 2023-01-05
```

## Data Format Requirements

### System Transactions (CSV)
| Column | Description |
| --- | --- |
| `trxID` | Unique identifier (string) |
| `amount` | Transaction amount (decimal) |
| `type` | `DEBIT` or `CREDIT` |
| `transactionTime` | Date and time (e.g., `2023-01-01T09:00:00Z`) |

### Bank Statements (CSV)
| Column | Description |
| --- | --- |
| `unique_identifier` | Unique identifier (string) |
| `amount` | Transaction amount (decimal) |
| `date` | Transaction date (`YYYY-MM-DD`) |

## Running Tests
To run unit tests for reconciliation logic and CSV parsing:

```bash
go test ./...
```
