package location

import (
	"os"

	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
)

// Root resolves ~/.hawp/models, matching internal/domain/provision's
// layout so pulled models live alongside the init-provisioned ones.
func Root() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return rootForHome(home)
}

func rootForHome(home string) (string, error) {
	models := filesystem.ResolveHawpHome(home).Models
	if err := filesystem.RejectSymlinksInPath(models); err != nil {
		return "", err
	}
	return models, nil
}
