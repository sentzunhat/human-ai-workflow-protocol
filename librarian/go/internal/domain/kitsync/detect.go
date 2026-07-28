package kitsync

import (
	"os"
	"path/filepath"
	"sort"
)

// DetectProviders returns the names of providers that appear to be
// installed in repoRoot, based on the manifest's own destination
// markers — no separate stored marker file is needed. A provider is
// considered detected when at least one of its pattern-filtered
// destination directories contains a matching file (a HAWP-specific
// signal: e.g. ".claude/rules/hawp-*.md" existing, not just ".claude/"
// existing, which avoids false positives from unrelated tooling that
// happens to use the same directory name).
func DetectProviders(repoRoot string, manifest *Manifest) []string {
	var detected []string
	for name, provider := range manifest.Providers {
		if providerInstalled(repoRoot, provider) {
			detected = append(detected, name)
		}
	}
	sort.Strings(detected)
	return detected
}

func providerInstalled(repoRoot string, provider Provider) bool {
	for _, rule := range provider.InstallsTo {
		if rule.Pattern == "" {
			continue
		}
		dir := filepath.Join(repoRoot, rule.Dest)
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if matched, _ := filepath.Match(rule.Pattern, entry.Name()); matched {
				return true
			}
		}
	}
	return false
}
