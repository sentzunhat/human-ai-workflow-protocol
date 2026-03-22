// Package links checks local markdown links across the repository's
// documentation roots (.hawp, docs, README.md). Ported from
// librarian/scripts/check-markdown-links.mjs — with the link regex fixed:
// the TS version's lookbehind matched nothing, so it validated zero links.
package links

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/markdown"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
)

var roots = []string{".hawp", "docs", "README.md"}

// Archival directories under .hawp/work are excluded, matching the
// dead-links policy in work validate: frozen history may reference paths
// that no longer exist.
var skipDirs = map[string]struct{}{
	".hawp/work/closed":   {},
	".hawp/work/evidence": {},
	".hawp/work/notes":    {},
	".hawp/work/status":   {},
}

// Result summarizes a link-check run.
type Result struct {
	FilesChecked int
	Failures     []string        // "relative/file.md -> target"
	Details      []FailureDetail // same failures, structured — used by Clean
}

// FailureDetail is one broken local link with enough detail to fix it: which
// file it's in and the exact markdown substring to replace. Raw is always
// exactly "[Text](Target)" — the link regex only matches that literal shape,
// so no whitespace/offset ambiguity is possible.
type FailureDetail struct {
	File    string // absolute path to the file containing the broken link
	RelFile string // repo-relative path
	Text    string
	Target  string
	Raw     string // "[Text](Target)" — the exact substring to replace
}

// Check scans all markdown under the documentation roots and verifies local
// link targets exist. Image links are skipped (matching the TS intent).
func Check(repoRoot string) Result {
	var files []string
	for _, root := range roots {
		full := filepath.Join(repoRoot, root)
		info, err := os.Stat(full)
		if err != nil {
			continue
		}
		if info.IsDir() {
			files = append(files, collectAll(repoRoot, full)...)
		} else if strings.HasSuffix(full, ".md") {
			files = append(files, full)
		}
	}

	result := Result{FilesChecked: len(files)}
	for _, file := range files {
		raw, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		content := markdown.BlankFences(string(raw))
		for _, link := range markdown.ExtractLinks(content) {
			if link.Image {
				continue
			}
			target := strings.Trim(link.Href, "<>")
			if target == "" || strings.HasPrefix(target, "#") || isExternal(target) {
				continue
			}
			pathPart := markdown.PathPart(target)
			var resolved string
			if strings.HasPrefix(pathPart, "/") {
				resolved = filepath.Join(repoRoot, "."+pathPart)
			} else {
				resolved = filepath.Join(filepath.Dir(file), pathPart)
			}
			if !repo.Exists(resolved) {
				relFile := repo.ToRepoRelative(repoRoot, file)
				result.Failures = append(result.Failures, relFile+" -> "+target)
				result.Details = append(result.Details, FailureDetail{
					File:    file,
					RelFile: relFile,
					Text:    link.Text,
					Target:  link.Href,
					Raw:     "[" + link.Text + "](" + link.Href + ")",
				})
			}
		}
	}
	return result
}

// isExternal reports scheme-prefixed (mailto:, https:, …) or
// protocol-relative (//host) targets.
func isExternal(target string) bool {
	if strings.HasPrefix(target, "//") {
		return true
	}
	for i, r := range target {
		switch {
		case r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z':
			continue
		case i > 0 && (r >= '0' && r <= '9' || r == '+' || r == '.' || r == '-'):
			continue
		case i > 0 && r == ':':
			return true
		default:
			return false
		}
	}
	return false
}

// collectAll gathers every .md file (including README.md) under dir,
// skipping archival directories.
func collectAll(repoRoot, dir string) []string {
	if _, skip := skipDirs[repo.ToRepoRelative(repoRoot, dir)]; skip {
		return nil
	}
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	for _, entry := range entries {
		full := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			files = append(files, collectAll(repoRoot, full)...)
		} else if strings.EqualFold(filepath.Ext(entry.Name()), ".md") {
			files = append(files, full)
		}
	}
	return files
}

// Render writes the human-readable result and returns the exit code.
func Render(out, errOut io.Writer, result Result) int {
	if len(result.Failures) > 0 {
		fmt.Fprintln(errOut, "Broken local Markdown links:")
		for _, failure := range result.Failures {
			fmt.Fprintln(errOut, "- "+failure)
		}
		return 1
	}
	fmt.Fprintf(out, "Checked %d Markdown file(s): local links are valid.\n", result.FilesChecked)
	return 0
}

// CleanChange is one broken link fixed (or that would be, in dry-run) —
// either relinked to a file found elsewhere in the repo, or neutralized
// (link syntax dropped, visible text kept) when no safe repair target exists.
type CleanChange struct {
	RelFile string
	Raw     string // the exact "[Text](Target)" that was/would be replaced
	New     string // the full replacement: a new "[Text](NewTarget)" link, or just Text if neutralized
	Action  string // "relinked" or "neutralized"
}

// CleanResult summarizes a Clean run.
type CleanResult struct {
	FilesChecked int
	Changes      []CleanChange
	Applied      bool // false = dry-run (nothing written), true = files were written
}

// buildBasenameIndex walks repoRoot once and indexes every file by its base
// name, so Clean can look up "did this file move somewhere else?" without a
// fresh filesystem walk per broken link. Skips version-control and
// dependency directories — large, irrelevant to a moved-doc search, and
// would slow the walk for no benefit.
func buildBasenameIndex(repoRoot string) map[string][]string {
	skipNames := map[string]bool{
		".git": true, "node_modules": true, "vendor": true,
		"dist": true, "build": true, ".next": true, "target": true,
	}
	index := make(map[string][]string)
	_ = filepath.WalkDir(repoRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // unreadable entry — skip it, don't abort the whole walk
		}
		if d.IsDir() {
			if skipNames[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		index[d.Name()] = append(index[d.Name()], path)
		return nil
	})
	return index
}

// findRelinkTarget looks for exactly one file in the repo whose base name
// matches the broken link's target (an anchor suffix, if present, is
// preserved on the rewritten link). Returns ok=false — meaning "don't
// guess" — when zero or more than one candidate exists; an ambiguous match
// is exactly as unsafe to auto-apply as no match at all.
func findRelinkTarget(index map[string][]string, sourceFile, target string) (newHref string, ok bool) {
	pathPart, anchor, _ := strings.Cut(target, "#")
	base := filepath.Base(pathPart)
	candidates := index[base]
	if len(candidates) != 1 {
		return "", false
	}

	rel, err := filepath.Rel(filepath.Dir(sourceFile), candidates[0])
	if err != nil {
		return "", false
	}
	rel = filepath.ToSlash(rel)
	if anchor != "" {
		rel += "#" + anchor
	}
	return rel, true
}

// Clean finds the same broken local links Check does and repairs each one:
//
//  1. Relink: search the repo for exactly one file with the target's base
//     name (e.g. a doc that moved from docs/a.md to docs/sub/a.md) and
//     rewrite the link to point at it. This is the preferred fix — it keeps
//     the reference working rather than just removing evidence it existed.
//  2. Neutralize: when no unique match is found (moved-and-renamed, deleted,
//     or ambiguous — multiple files share that base name, and guessing which
//     one is wrong is worse than not guessing), drop the link syntax and
//     keep the visible text — "[setup guide](gone.md)" becomes plain
//     "setup guide". A dangling link is worse than a plain-text mention: a
//     stale pointer looks like it should work and doesn't, indefinitely.
//
// Never touches the archival directories Check already skips
// (.hawp/work/closed, evidence, notes, status) — frozen history is allowed
// to reference removed paths by design.
//
// apply=false (the default from the CLI) only computes what would change;
// apply=true writes the fix to disk. Mirrors work normalize's
// dry-run-by-default convention.
func Clean(repoRoot string, apply bool) (CleanResult, error) {
	checked := Check(repoRoot)
	index := buildBasenameIndex(repoRoot)

	byFile := make(map[string][]FailureDetail)
	for _, d := range checked.Details {
		byFile[d.File] = append(byFile[d.File], d)
	}

	result := CleanResult{FilesChecked: checked.FilesChecked, Applied: apply}
	for file, details := range byFile {
		raw, err := os.ReadFile(file)
		if err != nil {
			return result, fmt.Errorf("read %s: %w", file, err)
		}
		content := string(raw)
		for _, d := range details {
			if !strings.Contains(content, d.Raw) {
				// Already fixed by an earlier pass, or the raw text isn't
				// unique enough to safely replace blindly — skip rather
				// than risk corrupting unrelated content.
				continue
			}

			var newRaw, action string
			if newHref, ok := findRelinkTarget(index, file, d.Target); ok {
				newRaw = "[" + d.Text + "](" + newHref + ")"
				action = "relinked"
			} else {
				newRaw = d.Text
				action = "neutralized"
			}

			content = strings.Replace(content, d.Raw, newRaw, 1)
			result.Changes = append(result.Changes, CleanChange{
				RelFile: d.RelFile,
				Raw:     d.Raw,
				New:     newRaw,
				Action:  action,
			})
		}
		if apply && content != string(raw) {
			if err := os.WriteFile(file, []byte(content), 0644); err != nil {
				return result, fmt.Errorf("write %s: %w", file, err)
			}
		}
	}

	return result, nil
}

// RenderClean writes the human-readable Clean result and returns the exit
// code (0 = nothing to do or successfully applied; 1 = dry-run found
// changes pending, matching the dry-run-by-default convention where "there
// is unapplied work" is a non-zero, non-error signal to scripts).
func RenderClean(out io.Writer, result CleanResult) int {
	if len(result.Changes) == 0 {
		fmt.Fprintf(out, "Checked %d Markdown file(s): no broken local links to clean.\n", result.FilesChecked)
		return 0
	}

	verbed := "Would fix"
	if result.Applied {
		verbed = "Fixed"
	}
	fmt.Fprintf(out, "%s %d broken local link(s):\n", verbed, len(result.Changes))
	for _, c := range result.Changes {
		fmt.Fprintf(out, "- [%s] %s: %s -> %s\n", c.Action, c.RelFile, c.Raw, c.New)
	}
	if !result.Applied {
		fmt.Fprintln(out, "\nRun with --apply to write these changes.")
		return 1
	}
	return 0
}
