package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

// HawpProject is the resolved .hawp/ layout for a project.
// Contains project-specific runtime data, work tracking, and patterns.
type HawpProject struct {
	Root   string // .hawp/ directory in the project root
	DB     string // .hawp/db — project-specific database
	Config string // .hawp/config — project-specific config

	Work string // .hawp/work — work tracking (in git)
	Kit  string // .hawp/kit — patterns and standards (in git)
}

// ResolveHawpProject builds the .hawp/ layout for a project.
func ResolveHawpProject(projectRoot string) HawpProject {
	hawpRoot := filepath.Join(projectRoot, ".hawp")
	return HawpProject{
		Root:   hawpRoot,
		DB:     filepath.Join(hawpRoot, "db"),
		Config: filepath.Join(hawpRoot, "config"),
		Work:   filepath.Join(hawpRoot, "work"),
		Kit:    filepath.Join(hawpRoot, "kit"),
	}
}

// EnsureRuntimeFolders creates the runtime directories used by search/config on
// first use. Returns true if folders were created, false if they already existed.
func (p *HawpProject) EnsureRuntimeFolders() (bool, error) {
	dirs := []string{p.DB, p.Config}
	created := false
	for _, dir := range dirs {
		if _, err := os.Stat(dir); err == nil {
			continue
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return false, fmt.Errorf("failed to create %s: %w", dir, err)
		}
		created = true
	}
	return created, nil
}

// EnsureDataFolders is kept as a compatibility wrapper for older callers that
// still refer to the previous helper name.
func (p *HawpProject) EnsureDataFolders() (bool, error) {
	return p.EnsureRuntimeFolders()
}

// CheckSearchIndexExists checks if the search index has been built.
// Returns false if not found, true if index.sqlite exists.
func (p *HawpProject) CheckSearchIndexExists() bool {
	indexPath := filepath.Join(p.DB, "index.sqlite")
	_, err := os.Stat(indexPath)
	return err == nil
}

// GetSearchIndexPath returns the path to the SQLite search index.
func (p *HawpProject) GetSearchIndexPath() string {
	return filepath.Join(p.DB, "index.sqlite")
}

// GetEmbeddingsCachePath returns the path to the project embeddings cache.
func (p *HawpProject) GetEmbeddingsCachePath() string {
	return filepath.Join(p.DB, "embeddings")
}

// GetProjectConfigPath returns the path to the project-specific config file.
func (p *HawpProject) GetProjectConfigPath() string {
	return filepath.Join(p.Config, "context.json")
}

// SuggestSearchIndexing returns a helpful message if the index is missing.
func (p *HawpProject) SuggestSearchIndexing() string {
	if !p.CheckSearchIndexExists() {
		return `
Search index not found. Build it with:
  hawp search index

This is a one-time setup for this project.
Rebuilds on demand if the codebase changes.
`
	}
	return ""
}
