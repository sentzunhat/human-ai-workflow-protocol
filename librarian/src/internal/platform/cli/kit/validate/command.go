package validate

import (
	appkit "github.com/sentzunhat/hawp/librarian/src/internal/application/kit"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/exitcode"
	"os"
	"path/filepath"
)

func Run(args []string, cwd string) error {
	o, err := parseKitValidateArgs(args)
	if err != nil {
		return err
	}
	p := o.kitPath
	if p == "" {
		root, err := repo.FindBacklogRepoRoot(cwd)
		if err != nil {
			return err
		}
		p = filepath.Join(root, ".hawp", "kit")
	}
	result := appkit.Validate(p)
	if code := appkit.Render(os.Stdout, os.Stderr, result); code != 0 {
		return exitcode.Error{Code: code}
	}
	return nil
}
