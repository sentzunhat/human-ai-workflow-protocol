package kitsync

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// SyncKit refreshes repoRoot/.hawp/kit/ from bundleKitDir, copying every
// file wholesale (kit content is canonical/generated, never hand-edited
// downstream — same assumption the existing distribution update guides
// already make).
func SyncKit(bundleKitDir, repoRoot string) (int, error) {
	destRoot := filepath.Join(repoRoot, ".hawp", "kit")
	return copyTree(bundleKitDir, destRoot, "")
}

// ApplyProviderUpdate refreshes providerName's update:refresh rules from
// bundleRoot (the extracted release bundle's top-level directory,
// containing both "kit/" and "providers/" — provider.Source, e.g.
// "providers/.claude", is already relative to this root) into repoRoot,
// per the manifest. update:skip rules (e.g. a customized CLAUDE.md) are
// left untouched.
func ApplyProviderUpdate(bundleRoot, repoRoot string, manifest *Manifest, providerName string) (int, []string, error) {
	provider, ok := manifest.Providers[providerName]
	if !ok {
		return 0, nil, fmt.Errorf("unknown provider %q", providerName)
	}

	written := 0
	var skipped []string
	sourceBase := filepath.Join(bundleRoot, provider.Source)

	for _, rule := range provider.InstallsTo {
		if !rule.ShouldRefreshOnUpdate() {
			skipped = append(skipped, providerName+":"+rule.Dest)
			continue
		}

		srcPath := filepath.Join(sourceBase, rule.From)
		destPath := filepath.Join(repoRoot, rule.Dest)

		info, err := os.Stat(srcPath)
		if err != nil {
			return written, skipped, fmt.Errorf("provider %s rule %s: source %s: %w", providerName, rule.Dest, srcPath, err)
		}

		if !info.IsDir() {
			if err := copyFile(srcPath, destPath); err != nil {
				return written, skipped, err
			}
			written++
			continue
		}

		count, err := copyTree(srcPath, destPath, rule.Pattern)
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
func copyTree(srcDir, destDir, pattern string) (int, error) {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return 0, err
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return 0, err
	}

	written := 0
	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			count, err := copyTree(srcPath, destPath, pattern)
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
		if err := copyFile(srcPath, destPath); err != nil {
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
func ApplyProviderInstall(bundleRoot, repoRoot string, manifest *Manifest, providerName string) (int, []string, error) {
	provider, ok := manifest.Providers[providerName]
	if !ok {
		return 0, nil, fmt.Errorf("unknown provider %q", providerName)
	}

	written := 0
	var seeded []string
	sourceBase := filepath.Join(bundleRoot, provider.Source)

	for _, rule := range provider.InstallsTo {
		srcPath := filepath.Join(sourceBase, rule.From)
		destPath := filepath.Join(repoRoot, rule.Dest)

		info, err := os.Stat(srcPath)
		if err != nil {
			return written, seeded, fmt.Errorf("provider %s rule %s: source %s: %w", providerName, rule.Dest, srcPath, err)
		}

		if rule.IsSeedIfMissing() {
			if info.IsDir() {
				count, err := seedTree(srcPath, destPath, rule.Pattern)
				if err != nil {
					return written, seeded, err
				}
				written += count
			} else {
				if _, statErr := os.Stat(destPath); os.IsNotExist(statErr) {
					if err := copyFile(srcPath, destPath); err != nil {
						return written, seeded, err
					}
					written++
					seeded = append(seeded, rule.Dest)
				}
			}
			continue
		}

		if !info.IsDir() {
			if err := copyFile(srcPath, destPath); err != nil {
				return written, seeded, err
			}
			written++
			continue
		}
		count, err := copyTree(srcPath, destPath, rule.Pattern)
		if err != nil {
			return written, seeded, err
		}
		written += count
	}

	return written, seeded, nil
}

// seedTree copies only files that do not already exist at their destination.
func seedTree(srcDir, destDir, pattern string) (int, error) {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return 0, err
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return 0, err
	}
	written := 0
	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())
		if entry.IsDir() {
			count, err := seedTree(srcPath, destPath, pattern)
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
		if _, statErr := os.Stat(destPath); !os.IsNotExist(statErr) {
			continue // already exists
		}
		if err := copyFile(srcPath, destPath); err != nil {
			return written, err
		}
		written++
	}
	return written, nil
}

func copyFile(srcPath, destPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return err
	}
	temp, err := os.CreateTemp(filepath.Dir(destPath), ".kitsync-*")
	if err != nil {
		return err
	}
	tempPath := temp.Name()
	defer os.Remove(tempPath)

	if _, err := io.Copy(temp, src); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	return os.Rename(tempPath, destPath)
}
