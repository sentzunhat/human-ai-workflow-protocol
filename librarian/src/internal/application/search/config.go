package search

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// SearchConfig holds user-configurable search settings.
// Loaded from ~/.hawp/config/search.json (home) then .hawp/config/search.json
// (project); project values override home values.
type SearchConfig struct {
	Index IndexConfig `json:"index"`
}

// IndexConfig controls which paths are walked when running `hawp search index`.
type IndexConfig struct {
	// Paths lists directories or files (relative to repo root) to include in
	// the search index. Defaults to [".hawp/kit", ".hawp/work"] when empty.
	Paths []string `json:"paths"`
}

// DefaultSearchConfig returns the built-in defaults.
func DefaultSearchConfig() SearchConfig {
	return SearchConfig{
		Index: IndexConfig{
			Paths: []string{".hawp/kit", ".hawp/work"},
		},
	}
}

// LoadSearchConfig loads config with home → project priority. Missing files
// are silently skipped; malformed files return an error.
func LoadSearchConfig(hawpHome, projectRoot string) (SearchConfig, error) {
	cfg := DefaultSearchConfig()

	if hawpHome != "" {
		if err := mergeSearchConfigFile(filepath.Join(hawpHome, "config", "search.json"), &cfg); err != nil {
			return cfg, err
		}
	}

	if projectRoot != "" {
		if err := mergeSearchConfigFile(filepath.Join(projectRoot, ".hawp", "config", "search.json"), &cfg); err != nil {
			return cfg, err
		}
	}

	return cfg, nil
}

func mergeSearchConfigFile(path string, cfg *SearchConfig) error {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	var fileCfg SearchConfig
	if err := json.Unmarshal(data, &fileCfg); err != nil {
		return err
	}
	if len(fileCfg.Index.Paths) > 0 {
		cfg.Index.Paths = fileCfg.Index.Paths
	}
	return nil
}
