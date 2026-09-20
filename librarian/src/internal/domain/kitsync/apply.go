package kitsync

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

type resolvedRule struct {
	rule   InstallRule
	source string
	dest   string
}

// resolveWithinRoot accepts only relative manifest paths whose cleaned target
// remains below root. Manifest content arrives from a downloaded release
// bundle, so it must not be allowed to select arbitrary filesystem paths.
func resolveWithinRoot(root, value, field string) (string, error) {
	if value == "" {
		return "", fmt.Errorf("%s must not be empty", field)
	}
	pathValue := filepath.FromSlash(value)
	if filepath.IsAbs(pathValue) {
		return "", fmt.Errorf("%s %q must be relative", field, value)
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve %s root: %w", field, err)
	}
	target := filepath.Join(root, pathValue)
	relative, err := filepath.Rel(root, target)
	if err != nil {
		return "", fmt.Errorf("check %s %q: %w", field, value, err)
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("%s %q escapes root %s", field, value, root)
	}
	return target, nil
}

func resolveRules(bundleRoot, repoRoot string, provider Provider) ([]resolvedRule, error) {
	sourceBase, err := resolveWithinRoot(bundleRoot, provider.Source, "provider source")
	if err != nil {
		return nil, err
	}
	rules := make([]resolvedRule, 0, len(provider.InstallsTo))
	for _, rule := range provider.InstallsTo {
		source, err := resolveWithinRoot(sourceBase, rule.From, "rule source")
		if err != nil {
			return nil, err
		}
		dest, err := resolveWithinRoot(repoRoot, rule.Dest, "rule destination")
		if err != nil {
			return nil, err
		}
		rules = append(rules, resolvedRule{rule: rule, source: source, dest: dest})
	}
	return rules, nil
}

// SyncKit refreshes repoRoot/.hawp/kit/ from bundleKitDir, copying every
// file wholesale (kit content is canonical/generated, never hand-edited
// downstream — same assumption the existing distribution update guides
// already make).
func SyncKit(fc FileCopier, bundleKitDir, repoRoot string) (int, error) {
	destRoot := filepath.Join(repoRoot, ".hawp", "kit")
	return copyTree(fc, bundleKitDir, destRoot, "")
}

// ApplyProviderUpdate applies providerName's update rules from bundleRoot
// (the extracted release bundle's top-level directory, containing both
// "kit/" and "providers/" — provider.Source, e.g. "providers/.claude",
// is already relative to this root) into repoRoot, per the manifest.
// update:refresh overwrites from the provider pack; update:seed-if-missing
// writes only when the destination is absent; update:skip leaves the path
// untouched.
func ApplyProviderUpdate(fc FileCopier, bundleRoot, repoRoot string, manifest *Manifest, providerName string) (int, []string, error) {
	provider, ok := manifest.Providers[providerName]
	if !ok {
		return 0, nil, fmt.Errorf("unknown provider %q", providerName)
	}

	rules, err := resolveRules(bundleRoot, repoRoot, provider)
	if err != nil {
		return 0, nil, fmt.Errorf("provider %s manifest paths: %w", providerName, err)
	}

	written := 0
	var skipped []string

	for _, resolved := range rules {
		rule := resolved.rule
		switch rule.UpdateMode() {
		case "skip":
			skipped = append(skipped, providerName+":"+rule.Dest)
			continue
		case "seed-if-missing":
			srcPath := resolved.source
			destPath := resolved.dest

			info, err := fc.Stat(srcPath)
			if err != nil {
				return written, skipped, fmt.Errorf("provider %s rule %s: source %s: %w", providerName, rule.Dest, srcPath, err)
			}

			if info.IsDir() {
				count, err := seedTree(fc, srcPath, destPath, rule.Pattern)
				if err != nil {
					return written, skipped, err
				}
				written += count
			} else {
				if _, statErr := fc.Stat(destPath); fc.IsNotExist(statErr) {
					if err := copyFile(fc, srcPath, destPath); err != nil {
						return written, skipped, err
					}
					written++
				}
			}
			continue
		}

		srcPath := resolved.source
		destPath := resolved.dest

		info, err := fc.Stat(srcPath)
		if err != nil {
			return written, skipped, fmt.Errorf("provider %s rule %s: source %s: %w", providerName, rule.Dest, srcPath, err)
		}

		if !info.IsDir() {
			if err := copyFile(fc, srcPath, destPath); err != nil {
				return written, skipped, err
			}
			written++
			continue
		}

		count, err := copyTree(fc, srcPath, destPath, rule.Pattern)
		if err != nil {
			return written, skipped, err
		}
		written += count
	}

	return written, skipped, nil
}

// copyTree copies every file from srcDir into destDir, optionally
// filtered by an fnmatch-style pattern on the base filename (recurses
// into subdirectories; the pattern only applies to file names, not
// directory names, so nested files still match e.g. "hawp-*.md").
func copyTree(fc FileCopier, srcDir, destDir, pattern string) (int, error) {
	entries, err := fc.ReadDir(srcDir)
	if err != nil {
		return 0, err
	}
	if err := fc.MkdirAll(destDir); err != nil {
		return 0, err
	}

	written := 0
	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			count, err := copyTree(fc, srcPath, destPath, pattern)
			if err != nil {
				return written, err
			}
			written += count
			continue
		}

		if pattern != "" {
			if matched, _ := filepath.Match(pattern, entry.Name()); !matched {
				continue
			}
		}
		if err := copyFile(fc, srcPath, destPath); err != nil {
			return written, err
		}
		written++
	}
	return written, nil
}

// ApplyProviderInstall places providerName's files into repoRoot for a
// first-time installation. Rules with install:seed-if-missing only write
// the file when the destination does not already exist (e.g. CLAUDE.md,
// AGENTS.md — the user will customise them). Rules without that flag are
// always written (refresh behaviour).
func ApplyProviderInstall(fc FileCopier, bundleRoot, repoRoot string, manifest *Manifest, providerName string) (int, []string, error) {
	provider, ok := manifest.Providers[providerName]
	if !ok {
		return 0, nil, fmt.Errorf("unknown provider %q", providerName)
	}

	rules, err := resolveRules(bundleRoot, repoRoot, provider)
	if err != nil {
		return 0, nil, fmt.Errorf("provider %s manifest paths: %w", providerName, err)
	}

	written := 0
	var seeded []string

	for _, resolved := range rules {
		rule := resolved.rule
		srcPath := resolved.source
		destPath := resolved.dest

		info, err := fc.Stat(srcPath)
		if err != nil {
			return written, seeded, fmt.Errorf("provider %s rule %s: source %s: %w", providerName, rule.Dest, srcPath, err)
		}

		if rule.IsSeedIfMissing() {
			if info.IsDir() {
				count, err := seedTree(fc, srcPath, destPath, rule.Pattern)
				if err != nil {
					return written, seeded, err
				}
				written += count
			} else {
				if _, statErr := fc.Stat(destPath); fc.IsNotExist(statErr) {
					if err := copyFile(fc, srcPath, destPath); err != nil {
						return written, seeded, err
					}
					written++
					seeded = append(seeded, rule.Dest)
				}
			}
			continue
		}

		if !info.IsDir() {
			if err := copyFile(fc, srcPath, destPath); err != nil {
				return written, seeded, err
			}
			written++
			continue
		}
		count, err := copyTree(fc, srcPath, destPath, rule.Pattern)
		if err != nil {
			return written, seeded, err
		}
		written += count
	}

	return written, seeded, nil
}

// seedTree copies only files that do not already exist at their destination.
func seedTree(fc FileCopier, srcDir, destDir, pattern string) (int, error) {
	entries, err := fc.ReadDir(srcDir)
	if err != nil {
		return 0, err
	}
	if err := fc.MkdirAll(destDir); err != nil {
		return 0, err
	}
	written := 0
	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())
		if entry.IsDir() {
			count, err := seedTree(fc, srcPath, destPath, pattern)
			if err != nil {
				return written, err
			}
			written += count
			continue
		}
		if pattern != "" {
			if matched, _ := filepath.Match(pattern, entry.Name()); !matched {
				continue
			}
		}
		if _, statErr := fc.Stat(destPath); !fc.IsNotExist(statErr) {
			continue // already exists
		}
		if err := copyFile(fc, srcPath, destPath); err != nil {
			return written, err
		}
		written++
	}
	return written, nil
}

func copyFile(fc FileCopier, srcPath, destPath string) error {
	src, err := fc.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	if err := fc.MkdirAll(filepath.Dir(destPath)); err != nil {
		return err
	}
	temp, tempName, err := fc.CreateTemp(filepath.Dir(destPath), ".kitsync-*")
	if err != nil {
		return err
	}
	defer fc.Remove(tempName)

	if _, err := io.Copy(temp, src); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	return fc.Rename(tempName, destPath)
}
