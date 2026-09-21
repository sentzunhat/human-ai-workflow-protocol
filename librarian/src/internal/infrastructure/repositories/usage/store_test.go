package usage

import (
	"path/filepath"
	"testing"

	domainusage "github.com/sentzunhat/hawp/librarian/src/internal/domain/usage"
)

func openTemp(t *testing.T) domainusage.Store {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "usage.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(s.Close)
	return s
}

func TestOpenWriteAndReport(t *testing.T) {
	s := openTemp(t)
	if err := s.Write("hawp_search", []byte(`{"query":"kubernetes deployments"}`), []byte(`{"content":"result"}`), false); err != nil {
		t.Fatal(err)
	}
	entries, err := s.Recent(10)
	if err != nil || len(entries) != 1 {
		t.Fatalf("Recent: %v, entries=%d", err, len(entries))
	}
	if entries[0].QueryText == nil || *entries[0].QueryText != "kubernetes deployments" {
		t.Fatalf("query text = %v", entries[0].QueryText)
	}
	if entries[0].InputBody != nil {
		t.Fatal("input body should be omitted")
	}
	report, err := s.GetReport()
	if err != nil {
		t.Fatal(err)
	}
	if report.Calls != 1 || len(report.ByTool) != 1 || len(report.TopEntries) != 1 {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestSchemaIsIdempotentAndClear(t *testing.T) {
	path := filepath.Join(t.TempDir(), "usage.db")
	first, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := first.Write("tool", []byte(`{"title":"item"}`), []byte(`{}`), true); err != nil {
		t.Fatal(err)
	}
	first.Close()
	second, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	if err := second.Clear(); err != nil {
		t.Fatal(err)
	}
	totals, err := second.GetTotals()
	if err != nil || totals.Calls != 0 {
		t.Fatalf("totals after clear: %+v, %v", totals, err)
	}
}

func TestConfigRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config", "usage.json")
	want := domainusage.Config{Enabled: true, LogBodies: true}
	if err := SaveConfig(path, want); err != nil {
		t.Fatal(err)
	}
	if got := LoadConfig(path); got != want {
		t.Fatalf("config = %+v, want %+v", got, want)
	}
}
