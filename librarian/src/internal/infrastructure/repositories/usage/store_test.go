package usage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"

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

func TestExtractQueryTextTruncatesAtValidUTF8Boundary(t *testing.T) {
	query := strings.Repeat("a", 255) + "é" + "suffix"
	entry := extractQueryText([]byte(`{"query":"` + query + `"}`))
	if entry == nil {
		t.Fatal("expected query text")
	}
	if !utf8.ValidString(*entry) {
		t.Fatalf("query text is invalid UTF-8: %q", *entry)
	}
	if len(*entry) != 255 {
		t.Fatalf("truncated query has %d bytes, want 255", len(*entry))
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

func TestOpenRejectsSymlinkedDatabasePaths(t *testing.T) {
	for _, tc := range []struct {
		name  string
		setup func(t *testing.T, root string) string
	}{
		{
			name: "parent directory",
			setup: func(t *testing.T, root string) string {
				t.Helper()
				outside := t.TempDir()
				link := filepath.Join(root, "db")
				if err := os.Symlink(outside, link); err != nil {
					t.Skipf("symlink capability unavailable: %v", err)
				}
				return filepath.Join(link, "usage.db")
			},
		},
		{
			name: "database file",
			setup: func(t *testing.T, root string) string {
				t.Helper()
				outside := filepath.Join(t.TempDir(), "outside.db")
				if err := os.WriteFile(outside, []byte("original"), 0o600); err != nil {
					t.Fatal(err)
				}
				path := filepath.Join(root, "usage.db")
				if err := os.Symlink(outside, path); err != nil {
					t.Skipf("symlink capability unavailable: %v", err)
				}
				return path
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			path := tc.setup(t, root)
			if store, err := Open(path); err == nil {
				store.Close()
				t.Fatal("Open accepted a symlinked database path")
			}
		})
	}
}

func TestOpenRejectsHardLinkedDatabaseAndSidecars(t *testing.T) {
	for _, suffix := range []string{"", "-wal", "-shm", "-journal"} {
		t.Run(suffix, func(t *testing.T) {
			root := t.TempDir()
			outside := filepath.Join(t.TempDir(), "outside.db")
			if err := os.WriteFile(outside, []byte("original"), 0o600); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(root, "usage.db") + suffix
			if err := os.Link(outside, path); err != nil {
				t.Skipf("hard links unavailable: %v", err)
			}
			if store, err := Open(filepath.Join(root, "usage.db")); err == nil {
				store.Close()
				t.Fatal("Open accepted a hard-linked database file")
			}
			if got, err := os.ReadFile(outside); err != nil || string(got) != "original" {
				t.Fatalf("external hard-link target changed: data=%q err=%v", got, err)
			}
		})
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

func TestSaveConfigRejectsSymlinkedDirectory(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	link := filepath.Join(root, "config")
	if err := os.Symlink(outside, link); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(link, "usage.json")
	if err := SaveConfig(path, domainusage.Config{Enabled: true}); err == nil {
		t.Fatal("SaveConfig followed a symlinked config directory")
	}
	if _, err := os.Stat(filepath.Join(outside, "usage.json")); !os.IsNotExist(err) {
		t.Fatalf("outside config was created: %v", err)
	}
}

func TestSaveConfigRejectsSymlinkedFile(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.json")
	if err := os.WriteFile(outside, []byte("original"), 0o644); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "usage.json")
	if err := os.Symlink(outside, path); err != nil {
		t.Fatal(err)
	}

	if err := SaveConfig(path, domainusage.Config{Enabled: true}); err == nil {
		t.Fatal("SaveConfig followed a symlinked config file")
	}
	got, err := os.ReadFile(outside)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "original" {
		t.Fatalf("outside config changed to %q", got)
	}
}
