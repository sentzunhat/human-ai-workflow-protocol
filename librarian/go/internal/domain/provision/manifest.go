package provision

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Manifest records what has been provisioned into ~/.hawp/ so `init` can
// verify and repair without blindly re-downloading.
type Manifest struct {
	SchemaVersion  int                    `json:"schemaVersion"`
	RuntimeVersion string                 `json:"runtimeVersion,omitempty"`
	ModelName      string                 `json:"modelName,omitempty"`
	ModelVersion   string                 `json:"modelVersion,omitempty"`
	Assets         map[string]AssetRecord `json:"assets"`
}

// AssetRecord is what the manifest remembers about one installed asset.
type AssetRecord struct {
	SHA256      string `json:"sha256"`
	Size        int64  `json:"size"`
	InstalledAt string `json:"installedAt"`
}

const manifestSchemaVersion = 1

// LoadManifest reads manifest.json from root, returning an empty manifest
// (not an error) when the file does not exist yet.
func LoadManifest(root string) (*Manifest, error) {
	path := filepath.Join(root, "manifest.json")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Manifest{SchemaVersion: manifestSchemaVersion, Assets: map[string]AssetRecord{}}, nil
	}
	if err != nil {
		return nil, err
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	if manifest.Assets == nil {
		manifest.Assets = map[string]AssetRecord{}
	}
	return &manifest, nil
}

// Save writes the manifest to root/manifest.json.
func (m *Manifest) Save(root string) error {
	m.SchemaVersion = manifestSchemaVersion
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(root, "manifest.json"), append(data, '\n'), 0o644)
}

// Satisfies reports whether the manifest already records this exact asset
// (by name and hash) as installed.
func (m *Manifest) Satisfies(asset Asset) bool {
	record, ok := m.Assets[asset.Name]
	return ok && record.SHA256 == asset.SHA256
}
