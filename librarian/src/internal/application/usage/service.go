// Package usage provides application-layer use cases for usage-log queries.
// Each function opens the store, performs one operation, and closes the store,
// returning only domain types to the caller.
package usage

import (
	domainusage "github.com/sentzunhat/hawp/librarian/src/internal/domain/usage"
	usageinfra "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repositories/usage"
)

// RecentLog returns the n most recent usage log entries from the store at dbPath.
func RecentLog(dbPath string, n int) ([]domainusage.Entry, error) {
	store, err := usageinfra.Open(dbPath)
	if err != nil {
		return nil, err
	}
	defer store.Close()
	return store.Recent(n)
}

// GetReport returns a full usage report from the store at dbPath.
func GetReport(dbPath string) (domainusage.Report, error) {
	store, err := usageinfra.Open(dbPath)
	if err != nil {
		return domainusage.Report{}, err
	}
	defer store.Close()
	return store.GetReport()
}

// ClearLog deletes all entries from the store at dbPath.
func ClearLog(dbPath string) error {
	store, err := usageinfra.Open(dbPath)
	if err != nil {
		return err
	}
	defer store.Close()
	return store.Clear()
}

// GetTotals returns aggregate totals from the store at dbPath.
func GetTotals(dbPath string) (domainusage.Totals, error) {
	store, err := usageinfra.Open(dbPath)
	if err != nil {
		return domainusage.Totals{}, err
	}
	defer store.Close()
	return store.GetTotals()
}
