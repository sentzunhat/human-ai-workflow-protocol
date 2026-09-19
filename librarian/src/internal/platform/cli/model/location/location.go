package location

import (
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
	"os"
)

// Root resolves ~/.hawp/models, matching internal/domain/provision's
// layout so pulled models live alongside the init-provisioned ones.
func Root() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filesystem.ResolveHawpHome(home).Models, nil
}
