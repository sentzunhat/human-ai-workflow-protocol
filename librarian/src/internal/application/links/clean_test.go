package links

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCleanRelinksWhenUniqueMatchFound(t *testing.T) {
	// setup-guide.md moved from docs/ to docs/guides/ — the link should be
	// repaired to point at the new location, not neutralized.
	root := buildRepo(t, map[string]string{
		"README.md":                  "See the [setup guide](docs/setup-guide.md) for details.\n",
		"docs/guides/setup-guide.md": "# Setup\n",
	})

	result, err := Clean(root, true)
	if err != nil {
		t.Fatalf("Clean failed: %v", err)
	}
	if len(result.Changes) != 1 {
		t.Fatalf("Changes = %+v, want exactly one", result.Changes)
	}
	if result.Changes[0].Action != "relinked" {
		t.Errorf("Action = %q, want relinked", result.Changes[0].Action)
	}
	if result.Changes[0].New != "[setup guide](docs/guides/setup-guide.md)" {
		t.Errorf("New = %q, want the repaired link pointing at the new location", result.Changes[0].New)
	}

	raw, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(raw)
	if !strings.Contains(got, "[setup guide](docs/guides/setup-guide.md)") {
		t.Errorf("file should contain the repaired link, got: %q", got)
	}
	if strings.Contains(got, "(docs/setup-guide.md)") {
		t.Error("file should no longer contain the broken link")
	}
}

func TestCleanNeutralizesWhenAmbiguousMatches(t *testing.T) {
	// Two files share the base name "setup-guide.md" — Clean must not guess
	// which one was intended, and should fall back to neutralizing instead.
	root := buildRepo(t, map[string]string{
		"README.md":             "See the [setup guide](missing-dir/setup-guide.md) for details.\n",
		"docs/a/setup-guide.md": "# A\n",
		"docs/b/setup-guide.md": "# B\n",
	})

	result, err := Clean(root, true)
	if err != nil {
		t.Fatalf("Clean failed: %v", err)
	}
	if len(result.Changes) != 1 {
		t.Fatalf("Changes = %+v, want exactly one", result.Changes)
	}
	if result.Changes[0].Action != "neutralized" {
		t.Errorf("Action = %q, want neutralized (ambiguous match must not be guessed)", result.Changes[0].Action)
	}
	if result.Changes[0].New != "setup guide" {
		t.Errorf("New = %q, want plain visible text", result.Changes[0].New)
	}
}

func TestCleanPreservesAnchorOnRelink(t *testing.T) {
	root := buildRepo(t, map[string]string{
		"README.md":                  "See [setup](docs/setup-guide.md#install) for install steps.\n",
		"docs/guides/setup-guide.md": "# Setup\n",
	})

	result, err := Clean(root, true)
	if err != nil {
		t.Fatalf("Clean failed: %v", err)
	}
	if len(result.Changes) != 1 || result.Changes[0].Action != "relinked" {
		t.Fatalf("Changes = %+v, want exactly one relink", result.Changes)
	}
	if result.Changes[0].New != "[setup](docs/guides/setup-guide.md#install)" {
		t.Errorf("New = %q, want the anchor preserved on the repaired link", result.Changes[0].New)
	}
}

func TestCleanDryRunReportsWithoutWriting(t *testing.T) {
	root := buildRepo(t, map[string]string{
		"README.md": "See the [setup guide](missing.md) for details.\n",
	})

	result, err := Clean(root, false)
	if err != nil {
		t.Fatalf("Clean failed: %v", err)
	}

	if result.Applied {
		t.Error("dry-run Clean should report Applied=false")
	}
	if len(result.Changes) != 1 {
		t.Fatalf("Changes = %+v, want exactly one", result.Changes)
	}
	if result.Changes[0].Raw != "[setup guide](missing.md)" {
		t.Errorf("Changes[0].Raw = %q, want the exact broken link substring", result.Changes[0].Raw)
	}
	if result.Changes[0].New != "setup guide" {
		t.Errorf("Changes[0].New = %q, want the visible text only (no match found anywhere)", result.Changes[0].New)
	}

	// File on disk must be untouched in dry-run mode.
	raw, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "[setup guide](missing.md)") {
		t.Error("dry-run Clean must not modify the file on disk")
	}
}

func TestCleanApplyWritesNeutralizedLink(t *testing.T) {
	root := buildRepo(t, map[string]string{
		"README.md": "See the [setup guide](missing.md) for details.\n",
	})

	result, err := Clean(root, true)
	if err != nil {
		t.Fatalf("Clean failed: %v", err)
	}
	if !result.Applied {
		t.Error("apply Clean should report Applied=true")
	}

	raw, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(raw)
	if strings.Contains(got, "[setup guide](missing.md)") {
		t.Error("apply Clean should have removed the broken link syntax")
	}
	if !strings.Contains(got, "See the setup guide for details.") {
		t.Errorf("apply Clean should keep the visible text in place, got: %q", got)
	}
}

func TestCleanLeavesValidLinksAlone(t *testing.T) {
	root := buildRepo(t, map[string]string{
		"README.md":               "[ok](.hawp/kit/start-here.md)\n[broken](missing.md)\n",
		".hawp/kit/start-here.md": "# start\n",
	})

	result, err := Clean(root, true)
	if err != nil {
		t.Fatalf("Clean failed: %v", err)
	}
	if len(result.Changes) != 1 {
		t.Fatalf("Changes = %+v, want exactly one (the broken link only)", result.Changes)
	}

	raw, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(raw)
	if !strings.Contains(got, "[ok](.hawp/kit/start-here.md)") {
		t.Error("Clean should not touch valid links")
	}
	if strings.Contains(got, "[broken](missing.md)") {
		t.Error("Clean should have fixed the broken link")
	}
}

func TestCleanNoOpWhenNoFailures(t *testing.T) {
	root := buildRepo(t, map[string]string{
		"README.md":               "[ok](.hawp/kit/start-here.md)\n",
		".hawp/kit/start-here.md": "# start\n",
	})

	result, err := Clean(root, true)
	if err != nil {
		t.Fatalf("Clean failed: %v", err)
	}
	if len(result.Changes) != 0 {
		t.Errorf("Changes = %+v, want none", result.Changes)
	}
}

func TestCleanSkipsArchivalDirectories(t *testing.T) {
	// Clean must respect the same archival skip-list as Check — frozen
	// history in .hawp/work/closed is allowed to reference removed paths.
	root := buildRepo(t, map[string]string{
		".hawp/work/closed/2026/01/01/old.md": "[stale](../does-not-exist.md)\n",
		".hawp/work/BACKLOG.md":               "# backlog\n",
	})

	result, err := Clean(root, true)
	if err != nil {
		t.Fatalf("Clean failed: %v", err)
	}
	if len(result.Changes) != 0 {
		t.Errorf("Changes = %+v, want none (archival directory should be skipped)", result.Changes)
	}

	raw, err := os.ReadFile(filepath.Join(root, ".hawp/work/closed/2026/01/01/old.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "[stale](../does-not-exist.md)") {
		t.Error("archival file should be left completely untouched")
	}
}

func TestRenderCleanExitCodes(t *testing.T) {
	var buf strings.Builder
	if code := RenderClean(&buf, CleanResult{FilesChecked: 3, Applied: false}); code != 0 {
		t.Errorf("no changes should exit 0, got %d", code)
	}

	buf.Reset()
	dryRun := CleanResult{FilesChecked: 1, Applied: false, Changes: []CleanChange{{RelFile: "a.md", Raw: "[x](y.md)", New: "x", Action: "neutralized"}}}
	if code := RenderClean(&buf, dryRun); code != 1 {
		t.Errorf("dry-run with pending changes should exit 1, got %d", code)
	}
	if !strings.Contains(buf.String(), "--apply") {
		t.Error("dry-run output should mention --apply")
	}

	buf.Reset()
	applied := CleanResult{FilesChecked: 1, Applied: true, Changes: []CleanChange{{RelFile: "a.md", Raw: "[x](y.md)", New: "x", Action: "neutralized"}}}
	if code := RenderClean(&buf, applied); code != 0 {
		t.Errorf("applied changes should exit 0, got %d", code)
	}
}
