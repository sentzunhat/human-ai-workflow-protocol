package cli_test

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	platformcli "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli"
)

func TestRunHelp(t *testing.T) {
	cases := [][]string{nil, {"--help"}, {"-h"}}
	for _, args := range cases {
		if err := platformcli.Run(args); err != nil {
			t.Errorf("Run(%v) returned error: %v", args, err)
		}
	}
}

func TestRunUnknownCommand(t *testing.T) {
	err := platformcli.Run([]string{"does-not-exist"})
	if err == nil {
		t.Fatal("Run with unknown command returned nil, want error")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Errorf("error %q does not mention unknown command", err.Error())
	}
	if !strings.Contains(err.Error(), "USAGE") {
		t.Errorf("error %q does not include help text", err.Error())
	}
}

func TestRunDBInit(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	if err := platformcli.Run([]string{"db", "init"}); err != nil {
		t.Fatalf("Run(db init) returned error: %v", err)
	}
}

func TestRunUUIDRoutes(t *testing.T) {
	for _, args := range [][]string{{"uuid"}, {"uuid", "--short"}} {
		if err := platformcli.Run(args); err != nil {
			t.Errorf("Run(%v) returned error: %v", args, err)
		}
	}
}

func TestRunValidateRoutesWithExplicitPaths(t *testing.T) {
	err := platformcli.Run([]string{"kit", "validate", "--kit-path", t.TempDir()})
	var exitErr platformcli.ExitError
	if !errors.As(err, &exitErr) || exitErr.Code != 1 {
		t.Errorf("kit validate on empty kit: err = %v, want ExitError{1}", err)
	}

	if err := platformcli.Run([]string{"work", "validate", "--work-root", "/nonexistent-hawp-work"}); err == nil {
		t.Error("work validate with bad root should error")
	}
}

func TestRunIndexBuild(t *testing.T) {
	cases := [][]string{
		{"index", "build"},
		{"index", "build", "--scope", "work"},
		{"index", "build", "--scope", "kit"},
	}
	for _, args := range cases {
		if err := platformcli.Run(args); err != nil {
			t.Errorf("Run(%v) returned error: %v", args, err)
		}
	}
}

func TestRunUsage(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	if err := platformcli.Run([]string{"usage", "enable"}); err != nil {
		t.Fatalf("usage enable: %v", err)
	}
	if err := platformcli.Run([]string{"usage", "enable", "--log-bodies"}); err != nil {
		t.Fatalf("usage enable --log-bodies: %v", err)
	}
	if err := platformcli.Run([]string{"usage"}); err != nil {
		t.Fatalf("usage totals: %v", err)
	}
	if err := platformcli.Run([]string{"usage", "log"}); err != nil {
		t.Fatalf("usage log: %v", err)
	}
	if err := platformcli.Run([]string{"usage", "disable"}); err != nil {
		t.Fatalf("usage disable: %v", err)
	}
	if err := platformcli.Run([]string{"usage"}); err != nil {
		t.Fatalf("usage totals (disabled): %v", err)
	}
}

func TestRunUsageReport(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	if err := platformcli.Run([]string{"usage", "enable"}); err != nil {
		t.Fatalf("usage enable: %v", err)
	}
	if err := platformcli.Run([]string{"usage", "report"}); err != nil {
		t.Fatalf("usage report (empty): %v", err)
	}
	dir := t.TempDir()
	exportPath := dir + "/report.md"
	if err := platformcli.Run([]string{"usage", "report", "--export", exportPath}); err != nil {
		t.Fatalf("usage report --export: %v", err)
	}
	data, err := os.ReadFile(exportPath)
	if err != nil {
		t.Fatalf("read exported report: %v", err)
	}
	if string(data) == "" {
		t.Error("exported report should be non-empty")
	}
}

func TestRunUsageClearCancelled(t *testing.T) {
	old := os.Stdin
	t.Cleanup(func() { os.Stdin = old })
	f, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	os.Stdin = f

	t.Setenv("HOME", t.TempDir())
	if err := platformcli.Run([]string{"usage", "enable"}); err != nil {
		t.Fatalf("usage enable: %v", err)
	}
	if err := platformcli.Run([]string{"usage", "clear"}); err != nil {
		t.Fatalf("usage clear (cancelled): %v", err)
	}
}

func TestRunSearchWithContext(t *testing.T) {
	cases := [][]string{
		{"search", "test"},
		{"search", "test", "--limit", "5"},
		{"search", "test", "--context"},
		{"search", "test", "--context", "--format", "json"},
		{"search", "test", "--context", "--format", "markdown"},
		{"search", "test", "--context", "--max-tokens", "1000"},
		{"search", "test", "--context", "--format", "json", "--max-tokens", "500"},
	}
	for _, args := range cases {
		_ = platformcli.Run(args)
	}
}

func TestRunWorkNormalizeSupportsHawpRootOverride(t *testing.T) {
	currentDir := t.TempDir()
	targetRepo := t.TempDir()

	targetFiles := map[string]string{
		".hawp/work/BACKLOG.md": `# Backlog

## Active Work

| # | Status | Title | Plan File | Next action |
| --- | --- | --- | --- | --- |
| 049 | in-progress | legacy | active/049.md | next |

## Blocked / Parked

| # | Status | Title | Plan File | Next action |
| --- | --- | --- | --- | --- |

## Recently Closed

| # | Title | Closed | Plan File |
| --- | --- | --- | --- |
`,
		".hawp/work/active/049.md": "# legacy plan\n",
	}
	for rel, content := range targetFiles {
		full := filepath.Join(targetRepo, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(currentDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(oldWd)
	})

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	rOut, wOut, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	rErr, wErr, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = wOut
	os.Stderr = wErr
	t.Cleanup(func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	})

	runErr := platformcli.Run([]string{
		"work", "normalize", "--dry-run", "--migrate-folders",
		"--hawp-root", filepath.Join(targetRepo, ".hawp"),
	})

	_ = wOut.Close()
	_ = wErr.Close()
	stdoutBytes, _ := io.ReadAll(rOut)
	_, _ = io.ReadAll(rErr)

	if runErr != nil {
		t.Fatalf("Run returned error: %v", runErr)
	}
	if !strings.Contains(string(stdoutBytes), ".hawp/work/active/049/plan.md") {
		t.Fatalf("expected target repo migration preview, got %q", string(stdoutBytes))
	}
}
