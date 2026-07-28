package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

// HawpProject is the resolved .hawp/ layout for a project.
// Contains project-specific data, work tracking, and patterns.
type HawpProject struct {
	Root   string // .hawp/ directory in the project root
	Data   string // .hawp/.data — auto-created runtime data (NOT in git)
	Config string // .hawp/.data/config — project-specific config
	DB     string // .hawp/.data/db — project-specific database

	Work string // .hawp/work — work tracking (in git)
	Kit  string // .hawp/kit — patterns and standards (in git)
}

// ResolveHawpProject builds the .hawp/ layout for a project.
func ResolveHawpProject(projectRoot string) HawpProject {
	hawpRoot := filepath.Join(projectRoot, ".hawp")
	dataRoot := filepath.Join(hawpRoot, ".data")
	return HawpProject{
		Root:   hawpRoot,
		Data:   dataRoot,
		Config: filepath.Join(dataRoot, "config"),
		DB:     filepath.Join(dataRoot, "db"),
		Work:   filepath.Join(hawpRoot, "work"),
		Kit:    filepath.Join(hawpRoot, "kit"),
	}
}

// EnsureDataFolders creates .hawp/.data/ structure on first run.
// Does NOT require git tracking — created automatically.
// Returns true if folders were created, false if they already existed.
func (p *HawpProject) EnsureDataFolders() (bool, error) {
	// Check if .data folder already exists
	if _, err := os.Stat(p.Data); err == nil {
		return false, nil // Already exists
	}

	// Create directory structure
	dirs := []string{p.Data, p.Config, p.DB}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return false, fmt.Errorf("failed to create %s: %w", dir, err)
		}
	}

	return true, nil
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
