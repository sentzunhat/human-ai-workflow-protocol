package usage_test

import (
	"path/filepath"
	"testing"

	domainusage "github.com/sentzunhat/hawp/librarian/src/internal/domain/usage"
	usageinfra "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repositories/usage"
	appusage "github.com/sentzunhat/hawp/librarian/src/internal/application/usage"
)

// seed writes one entry into a temp store so queries return non-empty results.
func seed(t *testing.T, dbPath string) {
	t.Helper()
	s, err := usageinfra.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if err := s.Write("hawp_search", []byte(`{"query":"test"}`), []byte(`{"content":"result"}`), false); err != nil {
		t.Fatal(err)
	}
}

func TestRecentLog(t *testing.T) {
	cases := []struct {
		name    string
		dbPath  func(t *testing.T) string
		n       int
		wantLen int
		wantErr bool
	}{
		{
			name:    "success returns entries",
			dbPath:  func(t *testing.T) string { p := filepath.Join(t.TempDir(), "usage.db"); seed(t, p); return p },
			n:       10,
			wantLen: 1,
		},
		{
			name:    "empty store returns nil slice",
			dbPath:  func(t *testing.T) string { return filepath.Join(t.TempDir(), "usage.db") },
			n:       10,
			wantLen: 0,
		},
		{
			name:    "bad path returns error",
			dbPath:  func(t *testing.T) string { return "/dev/null/nonexistent/usage.db" },
			n:       10,
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			entries, err := appusage.RecentLog(tc.dbPath(t), tc.n)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(entries) != tc.wantLen {
				t.Fatalf("len(entries) = %d, want %d", len(entries), tc.wantLen)
			}
		})
	}
}

func TestGetReport(t *testing.T) {
	cases := []struct {
		name      string
		dbPath    func(t *testing.T) string
		wantCalls int
		wantErr   bool
	}{
		{
			name:      "success returns report",
			dbPath:    func(t *testing.T) string { p := filepath.Join(t.TempDir(), "usage.db"); seed(t, p); return p },
			wantCalls: 1,
		},
		{
			name:      "empty store returns zero-call report",
			dbPath:    func(t *testing.T) string { return filepath.Join(t.TempDir(), "usage.db") },
			wantCalls: 0,
		},
		{
			name:    "bad path returns error",
			dbPath:  func(t *testing.T) string { return "/dev/null/nonexistent/usage.db" },
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rep, err := appusage.GetReport(tc.dbPath(t))
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if rep.Calls != tc.wantCalls {
				t.Fatalf("rep.Calls = %d, want %d", rep.Calls, tc.wantCalls)
			}
		})
	}
}

func TestClearLog(t *testing.T) {
	cases := []struct {
		name    string
		dbPath  func(t *testing.T) string
		wantErr bool
	}{
		{
			name:   "success clears entries",
			dbPath: func(t *testing.T) string { p := filepath.Join(t.TempDir(), "usage.db"); seed(t, p); return p },
		},
		{
			name:    "bad path returns error",
			dbPath:  func(t *testing.T) string { return "/dev/null/nonexistent/usage.db" },
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dbPath := tc.dbPath(t)
			err := appusage.ClearLog(dbPath)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			// verify store is empty after clear
			entries, err := appusage.RecentLog(dbPath, 10)
			if err != nil {
				t.Fatalf("RecentLog after clear: %v", err)
			}
			if len(entries) != 0 {
				t.Fatalf("expected 0 entries after clear, got %d", len(entries))
			}
		})
	}
}

func TestGetTotals(t *testing.T) {
	cases := []struct {
		name       string
		dbPath     func(t *testing.T) string
		wantTotals domainusage.Totals
		wantErr    bool
	}{
		{
			name:       "success returns totals",
			dbPath:     func(t *testing.T) string { p := filepath.Join(t.TempDir(), "usage.db"); seed(t, p); return p },
			wantTotals: domainusage.Totals{Calls: 1},
		},
		{
			name:       "empty store returns zero totals",
			dbPath:     func(t *testing.T) string { return filepath.Join(t.TempDir(), "usage.db") },
			wantTotals: domainusage.Totals{},
		},
		{
			name:    "bad path returns error",
			dbPath:  func(t *testing.T) string { return "/dev/null/nonexistent/usage.db" },
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			totals, err := appusage.GetTotals(tc.dbPath(t))
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if totals.Calls != tc.wantTotals.Calls {
				t.Fatalf("totals.Calls = %d, want %d", totals.Calls, tc.wantTotals.Calls)
			}
		})
	}
}
