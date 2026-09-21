package usage

import (
	"encoding/json"
	"os"
	"path/filepath"

	domainusage "github.com/sentzunhat/hawp/librarian/src/internal/domain/usage"
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
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
