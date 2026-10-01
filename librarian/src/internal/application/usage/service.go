// Package usage provides application-layer use cases for usage-log queries.
package usage

import (
	"fmt"

	domainusage "github.com/sentzunhat/hawp/librarian/src/internal/domain/usage"
)

// StoreOpener is supplied by the platform composition root. Application code
// depends on the domain store port and does not select a persistence adapter.
type StoreOpener func(path string) (domainusage.Store, error)

// Service contains the usage-log use cases backed by an injected store opener.
type Service struct {
	openStore StoreOpener
}

// NewService constructs usage use cases with a platform-provided store opener.
func NewService(openStore StoreOpener) Service {
	return Service{openStore: openStore}
}

func (s Service) open(path string) (domainusage.Store, error) {
	if s.openStore == nil {
		return nil, fmt.Errorf("usage store opener is not configured")
	}
	return s.openStore(path)
}

// RecentLog returns the n most recent usage log entries from the store at dbPath.
func (s Service) RecentLog(dbPath string, n int) ([]domainusage.Entry, error) {
	store, err := s.open(dbPath)
	if err != nil {
		return nil, err
	}
	defer store.Close()
	return store.Recent(n)
}

// GetReport returns a full usage report from the store at dbPath.
func (s Service) GetReport(dbPath string) (domainusage.Report, error) {
	store, err := s.open(dbPath)
	if err != nil {
		return domainusage.Report{}, err
	}
	defer store.Close()
	return store.GetReport()
}

// ClearLog deletes all entries from the store at dbPath.
func (s Service) ClearLog(dbPath string) error {
	store, err := s.open(dbPath)
	if err != nil {
		return err
	}
	defer store.Close()
	return store.Clear()
}

// GetTotals returns aggregate totals from the store at dbPath.
func (s Service) GetTotals(dbPath string) (domainusage.Totals, error) {
	store, err := s.open(dbPath)
	if err != nil {
		return domainusage.Totals{}, err
	}
	defer store.Close()
	return store.GetTotals()
}
