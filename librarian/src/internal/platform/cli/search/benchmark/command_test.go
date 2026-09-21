package benchmark

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestReshapeReportTotalsOnlySuccessfulRows(t *testing.T) {
	report := formatReshapeReport([]reshapeBenchResult{
		{Request: "successful", RawTokens: 100, ShapedTokens: 40},
		{Request: "failed", RawTokens: 100, ErrNote: "provider unavailable"},
	}, "ollama", "mistral")
	if !strings.Contains(report, "TOTAL (1/2 succeeded)** | **100** | **40** | **+60** | **+60%") {
		t.Fatalf("total should include successful rows only:\n%s", report)
	}
	if !strings.Contains(report, "Date: "+time.Now().UTC().Format("2006-01-02")) {
		t.Fatalf("report date should be generated at runtime:\n%s", report)
	}
}

func TestDownstreamReportTotalsOnlySuccessfulRows(t *testing.T) {
	report := formatDownstreamReport([]downstreamBenchResult{
		{Intent: "successful", RawTokens: 100, ShapedTokens: 40},
		{Intent: "failed", RawTokens: 100, ErrNote: "provider unavailable"},
	}, "ollama", "mistral")
	if !strings.Contains(report, "TOTAL (1/2 succeeded)** | **100** | **40** | **+60** | **+60%") {
		t.Fatalf("total should include successful rows only:\n%s", report)
	}
	if !strings.Contains(report, "Date: "+time.Now().UTC().Format("2006-01-02")) {
		t.Fatalf("report date should be generated at runtime:\n%s", report)
	}
}

func TestRunRejectsSymlinkedSearchDatabaseAncestry(t *testing.T) {

	root := t.TempDir()
	workDir := filepath.Join(root, ".hawp", "work")
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workDir, "BACKLOG.md"), []byte("# Backlog\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	external := t.TempDir()
	dbParent := filepath.Join(root, ".hawp", "db")
	if err := os.Symlink(external, dbParent); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	err := Run(nil, root)
	if err == nil {
		t.Fatal("accepted symlinked search database ancestry")
	}
	if !strings.Contains(err.Error(), "search index path") || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("unexpected error: %v", err)
	}
}
