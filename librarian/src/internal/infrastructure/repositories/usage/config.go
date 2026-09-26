package usage

import (
	"encoding/json"
	"os"
	"path/filepath"

	domainusage "github.com/sentzunhat/hawp/librarian/src/internal/domain/usage"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
)

func LoadConfig(path string) domainusage.Config {
	data, err := os.ReadFile(path)
	if err != nil {
		return domainusage.Config{}
	}
	var cfg domainusage.Config
	if json.Unmarshal(data, &cfg) != nil {
		return domainusage.Config{}
	}
	return cfg
}

func SaveConfig(path string, cfg domainusage.Config) error {
	if err := filesystem.RejectSymlinksInPath(path); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	root := filepath.Dir(path)
	if err := filesystem.RejectSymlinkAncestors(root, path); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return filesystem.AtomicWriteFile(root, path, data, 0o644)
}
