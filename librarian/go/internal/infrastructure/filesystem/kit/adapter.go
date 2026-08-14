// Package kit provides the filesystem adapter for the kit capability.
package kit

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/kit/source"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/markdown"
)

// Adapter acquires and mutates a kit workspace on disk.
type Adapter struct{}

func NewAdapter() Adapter { return Adapter{} }

func (Adapter) Read(kitPath string) (source.Snapshot, error) {
	entries, err := readEntries(kitPath, kitPath)
	if err != nil {
		if os.IsNotExist(err) {
			return source.Snapshot{}, nil
		}
		return source.Snapshot{}, err
	}
	snapshot := source.Snapshot{Entries: entries}
	for _, path := range markdown.CollectFiles(kitPath, false) {
		raw, readErr := os.ReadFile(path)
		if readErr != nil {
			continue
		}
		masked := markdown.BlankFences(string(raw))
		file := source.File{Path: path, RelPath: toSlashRelative(kitPath, path), Content: string(raw)}
		for _, link := range markdown.ExtractLinks(masked) {
			file.Links = append(file.Links, source.Link{Href: link.Href, Offset: link.Offset})
		}
		snapshot.Files = append(snapshot.Files, file)
	}
	return snapshot, nil
}

func readEntries(root, dir string) ([]source.Entry, error) {
	directory, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	entries := make([]source.Entry, 0, len(directory))
	for _, entry := range directory {
		path := filepath.Join(dir, entry.Name())
		entries = append(entries, source.Entry{Path: path, RelPath: toSlashRelative(root, path), Name: entry.Name(), IsDir: entry.IsDir()})
		if entry.IsDir() {
			children, childErr := readEntries(root, path)
			if childErr != nil {
				return nil, childErr
			}
			entries = append(entries, children...)
		}
	}
	return entries, nil
}

func (Adapter) ApplyRenames(renames []source.Rename) (conflictFrom, conflictTo string, err error) {
	sorted := append([]source.Rename(nil), renames...)
	sort.Slice(sorted, func(i, j int) bool { return len(sorted[i].From) > len(sorted[j].From) })
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

func (Adapter) ApplyLinkUpdates(updates []source.LinkUpdate) (int, error) {
	perFile := map[string][]source.LinkUpdate{}
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
		sort.Slice(fileUpdates, func(i, j int) bool { return fileUpdates[i].Start > fileUpdates[j].Start })
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

func toSlashRelative(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return path
	}
	return strings.ReplaceAll(rel, "\\", "/")
}
