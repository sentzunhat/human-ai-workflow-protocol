package kitsync

import (
	"fmt"
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
func DetectProviders(fc FileCopier, repoRoot string, manifest *Manifest) ([]string, error) {
	var detected []string
	for name, provider := range manifest.Providers {
		installed, err := providerInstalled(fc, repoRoot, provider)
		if err != nil {
			return nil, fmt.Errorf("detect provider %s: %w", name, err)
		}
		if installed {
			detected = append(detected, name)
		}
	}
	sort.Strings(detected)
	return detected, nil
}

func providerInstalled(fc FileCopier, repoRoot string, provider Provider) (bool, error) {
	for _, rule := range provider.InstallsTo {
		if rule.Pattern == "" {
			continue
		}
		dir, err := resolveWithinRoot(repoRoot, rule.Dest, "provider detection destination")
		if err != nil {
			return false, err
		}
		if err := fc.RejectSymlinkAncestors(repoRoot, dir); err != nil {
			return false, err
		}
		entries, err := fc.ReadDir(dir)
		if err != nil {
			if fc.IsNotExist(err) {
				continue
			}
			return false, fmt.Errorf("read provider destination %s: %w", dir, err)
		}
		for _, entry := range entries {
			if matched, _ := filepath.Match(rule.Pattern, entry.Name()); matched {
				return true, nil
			}
		}
	}
	return false, nil
}
