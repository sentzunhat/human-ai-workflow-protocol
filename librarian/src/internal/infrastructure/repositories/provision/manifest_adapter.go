// Package provision provides infrastructure-layer adapters for domain provision operations.
package provision

import (
	"os"

	domainprovision "github.com/sentzunhat/hawp/librarian/src/internal/domain/provision"
)

// LoadManifest wraps the domain function with os.ReadFile as the concrete reader.
func LoadManifest(root string) (*domainprovision.Manifest, error) {
	return domainprovision.LoadManifest(root, os.ReadFile)
}

// Save writes manifest data using os.WriteFile as the concrete writer.
func Save(m *domainprovision.Manifest, root string) error {
	return m.Save(root, os.WriteFile)
}
