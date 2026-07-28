package kit

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/markdown"
)

// FileRename is one planned rename (absolute paths).
type FileRename struct {
	From string
	To   string
}

// LinkUpdate is one planned href rewrite inside a markdown file.
type LinkUpdate struct {
	File  string
	From  string
	To    string
	Start int
	End   int
}

var nonAlnumRe = regexp.MustCompile(`[^a-z0-9]+`)

// NormalizeFileName returns the lowercase-hyphen form of fileName, or ""
// when the name is already normalized (or is an allowed exact name).
func NormalizeFileName(fileName string) string {
	if fileName == "README.md" {
		return ""
	}
	lastDot := strings.LastIndex(fileName, ".")
	hasExt := lastDot > 0 && lastDot < len(fileName)-1
	stem, ext := fileName, ""
	if hasExt {
		stem, ext = fileName[:lastDot], strings.ToLower(fileName[lastDot:])
	}

	normalizedStem := strings.Trim(nonAlnumRe.ReplaceAllString(strings.ToLower(stem), "-"), "-")
	if normalizedStem == "" {
		return ""
	}
	normalized := normalizedStem + ext
	if normalized == fileName {
		return ""
	}
	return normalized
}

// PlanFileRenames walks the kit and plans lowercase-hyphen renames.
func PlanFileRenames(kitPath string) []FileRename {
	var renames []FileRename
	var walk func(dir string)
	walk = func(dir string) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, entry := range entries {
			full := filepath.Join(dir, entry.Name())
			if entry.IsDir() {
				walk(full)
				continue
			}
			if normalized := NormalizeFileName(entry.Name()); normalized != "" {
				renames = append(renames, FileRename{From: full, To: filepath.Join(dir, normalized)})
			}
		}
	}
	walk(kitPath)
	return renames
}

// PlanLinkUpdates finds relative links whose targets are being renamed and
// plans the href rewrites. Fenced code blocks are ignored.
func PlanLinkUpdates(kitPath string, renameMap map[string]string) []LinkUpdate {
	var updates []LinkUpdate
	for _, file := range markdown.CollectFiles(kitPath, false) {
		raw, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		masked := markdown.BlankFences(string(raw))
		fileDir := filepath.Dir(file)

		for _, link := range markdown.ExtractLinks(masked) {
			href := link.Href
			if href == "" || strings.HasPrefix(href, "http") ||
				strings.HasPrefix(href, "/") || strings.HasPrefix(href, "#") {
				continue
			}
			pathPart, anchor := href, ""
			if i := strings.Index(href, "#"); i >= 0 {
				pathPart, anchor = href[:i], href[i:]
			}
			if pathPart == "" {
				continue
			}
			target := filepath.Join(fileDir, pathPart)
			renamedTarget, ok := renameMap[target]
			if !ok {
				continue
			}
			rel, err := filepath.Rel(fileDir, renamedTarget)
			if err != nil {
				continue
			}
			nextHref := strings.ReplaceAll(rel, "\\", "/") + anchor
			if nextHref == href {
				continue
			}
			hrefStart := link.Offset + strings.Index(masked[link.Offset:], "("+href) + 1
			updates = append(updates, LinkUpdate{
				File: file, From: href, To: nextHref,
				Start: hrefStart, End: hrefStart + len(href),
			})
		}
	}
	return updates
}

// ApplyRenames performs the renames longest-path-first, refusing to
// overwrite existing targets. Returns the first conflict path pair, if any.
func ApplyRenames(renames []FileRename) (conflictFrom, conflictTo string, err error) {
	sorted := append([]FileRename(nil), renames...)
	sort.Slice(sorted, func(i, j int) bool {
		return len(sorted[i].From) > len(sorted[j].From)
	})
	for _, rename := range sorted {
		if _, statErr := os.Stat(rename.To); statErr == nil {
			return rename.From, rename.To, nil
		}
		if renameErr := os.Rename(rename.From, rename.To); renameErr != nil {
			return "", "", renameErr
		}
	}
	return "", "", nil
}

// ApplyLinkUpdates rewrites hrefs bottom-up per file and returns how many
// files changed.
func ApplyLinkUpdates(updates []LinkUpdate) (int, error) {
	perFile := map[string][]LinkUpdate{}
	for _, update := range updates {
		perFile[update.File] = append(perFile[update.File], update)
	}
	changed := 0
	for file, fileUpdates := range perFile {
		raw, err := os.ReadFile(file)
		if err != nil {
			return changed, err
		}
		next := string(raw)
		sort.Slice(fileUpdates, func(i, j int) bool {
			return fileUpdates[i].Start > fileUpdates[j].Start
		})
		for _, update := range fileUpdates {
			next = next[:update.Start] + update.To + next[update.End:]
		}
		if next != string(raw) {
			if err := os.WriteFile(file, []byte(next), 0o644); err != nil {
				return changed, err
			}
			changed++
		}
	}
	return changed, nil
}
