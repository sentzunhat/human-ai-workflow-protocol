package normalize

import (
	"os"
	"path/filepath"

	appkit "github.com/sentzunhat/hawp/librarian/src/internal/application/kit"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/exitcode"
)

func Run(args []string, cwd string) error {
	o, err := parseKitNormalizeArgs(args)
	if err != nil {
		return err
	}
	root, err := repo.FindBacklogRepoRoot(cwd)
	if err != nil {
		return err
	}
	p := resolveKitPath(root, cwd, o.kitPath)
	if code := appkit.Normalize(os.Stdout, os.Stderr, appkit.NormalizeOptions{KitPath: p, RepoRoot: root, Apply: o.apply}); code != 0 {
		return exitcode.Error{Code: code}
	}
	return nil
}

func resolveKitPath(root, cwd, configured string) string {
	if configured == "" {
		return filepath.Join(root, ".hawp", "kit")
	}
	if filepath.IsAbs(configured) {
		return filepath.Clean(configured)
	}
	return filepath.Clean(filepath.Join(cwd, configured))
}
