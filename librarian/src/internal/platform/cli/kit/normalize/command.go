package normalize

import (
	appkit "github.com/sentzunhat/hawp/librarian/src/internal/application/kit"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/exitcode"
	"os"
	"path/filepath"
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
	p := o.kitPath
	if p == "" {
		p = filepath.Join(root, ".hawp", "kit")
	}
	if code := appkit.Normalize(os.Stdout, os.Stderr, appkit.NormalizeOptions{KitPath: p, RepoRoot: root, Apply: o.apply}); code != 0 {
		return exitcode.Error{Code: code}
	}
	return nil
}
