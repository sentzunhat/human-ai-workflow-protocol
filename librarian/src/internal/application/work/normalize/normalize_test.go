package work

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestNormalizeRejectsSymlinkedWorkTreeBeforeReading(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires elevated privileges on some Windows environments")
	}
	root := t.TempDir()
	outside := t.TempDir()
	for rel, content := range map[string]string{
		"BACKLOG.md":       duplicateBacklogHeader + duplicateBacklogFooter,
		"active/secret.md": "# OUTSIDE_NORMALIZE_SENTINEL\n",
	} {
		path := filepath.Join(outside, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, ".hawp"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, ".hawp", "work")); err != nil {
		t.Fatal(err)
	}

	for _, opts := range []NormalizeOptions{
		{RepoRoot: root},
		{RepoRoot: root, MigrateFolders: true},
	} {
		var stdout, stderr bytes.Buffer
		if code := Normalize(&stdout, &stderr, opts); code != 1 {
			t.Fatalf("Normalize(%+v) = %d, want 1", opts, code)
		}
		if !strings.Contains(stderr.String(), "symlink") {
			t.Fatalf("Normalize(%+v) error = %q, want symlink rejection", opts, stderr.String())
		}
		if strings.Contains(stdout.String(), "OUTSIDE_NORMALIZE_SENTINEL") || strings.Contains(stderr.String(), "OUTSIDE_NORMALIZE_SENTINEL") {
			t.Fatalf("Normalize(%+v) read content outside the repository", opts)
		}
	}
}

func TestNormalizeApplyMigrateFoldersUsesProductionFileCollector(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		".hawp/work/BACKLOG.md": `# Backlog

## Active Work

| ID | Type | Title | Status | Plan File | Updated |
| --- | --- | --- | --- | --- | --- |
| ` + "`legacy-item`" + ` | task | legacy | inbox | [plan](active/legacy-item.md) | 2026-09-18 |

## Blocked / Parked

| ID | Type | Title | Reason | Detail | Updated |
| --- | --- | --- | --- | --- | --- |

## Recently Closed

| ID | Type | Title | Closed | Detail |
| --- | --- | --- | --- | --- |
`,
		".hawp/work/active/legacy-item.md": "# Legacy item\n\n**Plan file:** work/active/legacy-item.md\n",
	}
	for rel, content := range files {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	code := Normalize(io.Discard, io.Discard, NormalizeOptions{
		RepoRoot:       root,
		Apply:          true,
		MigrateFolders: true,
		ForceDirty:     true,
	})
	if code != 0 {
		t.Fatalf("Normalize returned exit code %d", code)
	}
	if _, err := os.Stat(filepath.Join(root, ".hawp/work/active/legacy-item/plan.md")); err != nil {
		t.Fatalf("migrated plan missing: %v", err)
	}
}
