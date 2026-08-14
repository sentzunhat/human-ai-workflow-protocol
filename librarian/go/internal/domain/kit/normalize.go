package kit

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/kit/source"
)

var nonAlnumRe = regexp.MustCompile(`[^a-z0-9]+`)

// FileRename and LinkUpdate keep the internal plan shapes available from the
// kit capability while their storage-facing definitions live at the boundary.
type FileRename = source.Rename
type LinkUpdate = source.LinkUpdate

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

// PlanFileRenames returns lowercase-hyphen rename operations from an
// acquired workspace snapshot.
func PlanFileRenames(snapshot source.Snapshot) []source.Rename {
	var renames []source.Rename
	for _, entry := range snapshot.Entries {
		if normalized := NormalizeFileName(entry.Name); normalized != "" {
			renames = append(renames, source.Rename{From: entry.Path, To: filepath.Join(filepath.Dir(entry.Path), normalized)})
		}
	}
	return renames
}

// PlanLinkUpdates finds links whose direct targets are being renamed. Fenced
// code blocks were excluded by the adapter before this pure rule runs.
func PlanLinkUpdates(snapshot source.Snapshot, renameMap map[string]string) []source.LinkUpdate {
	var updates []source.LinkUpdate
	for _, file := range snapshot.Files {
		fileDir := filepath.Dir(file.Path)
		for _, link := range file.Links {
			href := link.Href
			if !isLocalHref(href) {
				continue
			}
			pathPart, anchor := TrimAnchor(href), ""
			if len(pathPart) < len(href) {
				anchor = href[len(pathPart):]
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
			hrefStart := link.Offset + strings.Index(file.Content[link.Offset:], "("+href) + 1
			updates = append(updates, source.LinkUpdate{File: file.Path, From: href, To: nextHref, Start: hrefStart, End: hrefStart + len(href)})
		}
	}
	return updates
}
