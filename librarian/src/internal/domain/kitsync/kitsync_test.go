package kitsync

import (
	"os"
	"path/filepath"
	"testing"
)

// sampleManifest mirrors the real core/providers/manifest.yaml shape,
// including the always-refresh github entries with no install/update
// fields at all.
const sampleManifest = `
providers:
  claude:
    source: providers/.claude
    installs_to:
      - dest: .claude/rules/
        from: rules/
        pattern: "hawp-*.md"
        install: refresh
        update: refresh
      - dest: CLAUDE.md
        from: CLAUDE.md.seed
        install: seed-if-missing
        update: skip
  github:
    source: providers/.github
    installs_to:
      - dest: .github/instructions/
        from: instructions/
      - dest: .github/copilot-instructions.md
        from: copilot-instructions.md
        install: seed_if_missing
        update: refresh
`

func writeTree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func parseSample(t *testing.T) *Manifest {
	t.Helper()
	dir := writeTree(t, map[string]string{"manifest.yaml": sampleManifest})
	manifest, err := ParseManifest(filepath.Join(dir, "manifest.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return manifest
}

func TestParseManifestMatchesRealShape(t *testing.T) {
	manifest := parseSample(t)
	claude, ok := manifest.Providers["claude"]
	if !ok || len(claude.InstallsTo) != 2 {
		t.Fatalf("claude provider = %+v", claude)
	}
	if claude.InstallsTo[0].Pattern != "hawp-*.md" || claude.InstallsTo[0].Update != "refresh" {
		t.Errorf("claude rules rule = %+v", claude.InstallsTo[0])
	}
	if claude.InstallsTo[1].Update != "skip" {
		t.Errorf("CLAUDE.md rule should be update:skip, got %+v", claude.InstallsTo[1])
	}
}

func TestShouldRefreshOnUpdate(t *testing.T) {
	cases := []struct {
		update string
		want   bool
	}{
		{"refresh", true},
		{"skip", false},
		{"", true}, // missing field (github instructions/) defaults to refresh
	}
	for _, c := range cases {
		rule := InstallRule{Update: c.update}
		if got := rule.ShouldRefreshOnUpdate(); got != c.want {
			t.Errorf("ShouldRefreshOnUpdate(update=%q) = %v, want %v", c.update, got, c.want)
		}
	}
}

func TestDetectProvidersFindsPatternMarkedFiles(t *testing.T) {
	manifest := parseSample(t)

	repoRoot := writeTree(t, map[string]string{
		".claude/rules/hawp-core.md": "# core\n",
		".claude/rules/other.md":     "# non-hawp file, should not itself trigger detection\n",
	})
	detected := DetectProviders(repoRoot, manifest)
	if len(detected) != 1 || detected[0] != "claude" {
		t.Fatalf("detected = %v, want [claude]", detected)
	}
}

func TestDetectProvidersIgnoresUnrelatedGithubFolder(t *testing.T) {
	manifest := parseSample(t)

	// A generic .github/ (e.g. real CI workflows) must NOT be mistaken
	// for the HAWP github provider — github's manifest entries have no
	// `pattern`, so DetectProviders should not fire on them at all
	// (avoiding false positives from ordinary repo structure).
	repoRoot := writeTree(t, map[string]string{
		".github/workflows/ci.yml": "name: CI\n",
	})
	detected := DetectProviders(repoRoot, manifest)
	if len(detected) != 0 {
		t.Fatalf("detected = %v, want none (github has no pattern-marked rule)", detected)
	}
}

func TestSyncKitCopiesWholeTree(t *testing.T) {
	bundleKit := writeTree(t, map[string]string{
		"start-here.md": "# start\n",
		"usage/init.md": "# init\n",
	})
	repoRoot := t.TempDir()

	written, err := SyncKit(bundleKit, repoRoot)
	if err != nil {
		t.Fatal(err)
	}
	if written != 2 {
		t.Fatalf("written = %d, want 2", written)
	}
	content, err := os.ReadFile(filepath.Join(repoRoot, ".hawp", "kit", "usage", "init.md"))
	if err != nil || string(content) != "# init\n" {
		t.Fatalf("kit file not synced correctly: %v %q", err, content)
	}
}

func TestApplyProviderUpdateRefreshesAndSkips(t *testing.T) {
	manifest := parseSample(t)

	bundleRoot := writeTree(t, map[string]string{
		"providers/.claude/rules/hawp-core.md":    "# updated core content\n",
		"providers/.claude/rules/not-matched.txt": "should not be copied\n",
		"providers/.claude/CLAUDE.md.seed":        "# seed, should never be applied on update\n",
	})
	repoRoot := writeTree(t, map[string]string{
		".claude/rules/hawp-core.md": "# stale old content\n",
		"CLAUDE.md":                  "# user's customized CLAUDE.md, must survive\n",
	})

	written, skipped, err := ApplyProviderUpdate(bundleRoot, repoRoot, manifest, "claude")
	if err != nil {
		t.Fatal(err)
	}
	if written != 1 {
		t.Fatalf("written = %d, want 1 (only the pattern-matched file)", written)
	}
	if len(skipped) != 1 || skipped[0] != "claude:CLAUDE.md" {
		t.Fatalf("skipped = %v, want [claude:CLAUDE.md]", skipped)
	}

	refreshed, err := os.ReadFile(filepath.Join(repoRoot, ".claude", "rules", "hawp-core.md"))
	if err != nil || string(refreshed) != "# updated core content\n" {
		t.Fatalf("rule file not refreshed: %v %q", err, refreshed)
	}
	if _, err := os.Stat(filepath.Join(repoRoot, ".claude", "rules", "not-matched.txt")); !os.IsNotExist(err) {
		t.Error("non-matching pattern file should not have been copied")
	}
	untouched, err := os.ReadFile(filepath.Join(repoRoot, "CLAUDE.md"))
	if err != nil || string(untouched) != "# user's customized CLAUDE.md, must survive\n" {
		t.Fatalf("CLAUDE.md (update:skip) was modified: %v %q", err, untouched)
	}
}

func TestApplyProviderUpdateUnknownProvider(t *testing.T) {
	manifest := parseSample(t)
	if _, _, err := ApplyProviderUpdate(t.TempDir(), t.TempDir(), manifest, "nonexistent"); err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

// --- ApplyProviderInstall ---

func TestApplyProviderInstallCreatesAllFiles(t *testing.T) {
	manifest := parseSample(t)

	// Fresh repo: no .claude/ or CLAUDE.md yet.
	bundleRoot := writeTree(t, map[string]string{
		"providers/.claude/rules/hawp-core.md": "# core\n",
		"providers/.claude/CLAUDE.md.seed":     "# starter CLAUDE.md\n",
	})
	repoRoot := t.TempDir()

	written, seeded, err := ApplyProviderInstall(bundleRoot, repoRoot, manifest, "claude")
	if err != nil {
		t.Fatal(err)
	}
	// hawp-core.md (refresh) + CLAUDE.md.seed (seed-if-missing, dest missing) = 2
	if written != 2 {
		t.Fatalf("written = %d, want 2", written)
	}
	// CLAUDE.md is seed-if-missing — should appear in seeded list
	if len(seeded) != 1 || seeded[0] != "CLAUDE.md" {
		t.Fatalf("seeded = %v, want [CLAUDE.md]", seeded)
	}

	if _, err := os.ReadFile(filepath.Join(repoRoot, ".claude", "rules", "hawp-core.md")); err != nil {
		t.Fatalf("rule file not installed: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(repoRoot, "CLAUDE.md"))
	if err != nil || string(content) != "# starter CLAUDE.md\n" {
		t.Fatalf("CLAUDE.md not seeded: %v %q", err, content)
	}
}

func TestApplyProviderInstallSeedIfMissingSkipsExisting(t *testing.T) {
	manifest := parseSample(t)

	bundleRoot := writeTree(t, map[string]string{
		"providers/.claude/rules/hawp-core.md": "# updated core\n",
		"providers/.claude/CLAUDE.md.seed":     "# seed\n",
	})
	// Repo already has CLAUDE.md with user content.
	repoRoot := writeTree(t, map[string]string{
		"CLAUDE.md": "# my custom content\n",
	})

	written, _, err := ApplyProviderInstall(bundleRoot, repoRoot, manifest, "claude")
	if err != nil {
		t.Fatal(err)
	}
	// only hawp-core.md written; CLAUDE.md already exists → skipped
	if written != 1 {
		t.Fatalf("written = %d, want 1 (CLAUDE.md already exists)", written)
	}
	content, err := os.ReadFile(filepath.Join(repoRoot, "CLAUDE.md"))
	if err != nil || string(content) != "# my custom content\n" {
		t.Fatalf("CLAUDE.md (existing) was overwritten: %v %q", err, content)
	}
}

func TestApplyProviderInstallUnknownProvider(t *testing.T) {
	manifest := parseSample(t)
	if _, _, err := ApplyProviderInstall(t.TempDir(), t.TempDir(), manifest, "nonexistent"); err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

func TestIsSeedIfMissing(t *testing.T) {
	cases := []struct {
		install string
		want    bool
	}{
		{"seed-if-missing", true},
		{"seed_if_missing", true},
		{"refresh", false},
		{"", false},
	}
	for _, c := range cases {
		rule := InstallRule{Install: c.install}
		if got := rule.IsSeedIfMissing(); got != c.want {
			t.Errorf("IsSeedIfMissing(install=%q) = %v, want %v", c.install, got, c.want)
		}
	}
}

func TestAllProviderNames(t *testing.T) {
	manifest := parseSample(t)
	names := manifest.AllProviderNames()
	if len(names) != 2 { // claude + github from sampleManifest
		t.Fatalf("AllProviderNames() = %v, want 2 entries", names)
	}
}
