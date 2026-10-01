// Package provision provides infrastructure-layer adapters for domain provision operations.
package provision

import (
	"io/fs"
	"os"
	"path/filepath"

	domainprovision "github.com/sentzunhat/hawp/librarian/src/internal/domain/provision"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
)

// LoadManifest wraps the domain function with os.ReadFile as the concrete reader.
func LoadManifest(root string) (*domainprovision.Manifest, error) {
	return domainprovision.LoadManifest(root, os.ReadFile)
}

// Save writes manifest data through the guarded atomic writer. Replacing the
// directory entry avoids following a substituted final symlink or truncating
// an unrelated hard-linked inode.
func Save(m *domainprovision.Manifest, root string) error {
	manifestPath := filepath.Join(root, "manifest.json")
	if err := filesystem.RejectSymlinksInPath(manifestPath); err != nil {
		return err
	}
	return m.Save(root, func(path string, data []byte, perm fs.FileMode) error {
		return filesystem.AtomicWriteFile(root, path, data, perm)
	})
}
