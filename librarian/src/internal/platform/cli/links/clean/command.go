package clean

import (
	"os"

	applinks "github.com/sentzunhat/hawp/librarian/src/internal/application/links"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/exitcode"
)

func Run(args []string, cwd string) error {
	opts, err := parseArgs(args)
	if err != nil {
		return err
	}

	root, err := repo.FindBacklogRepoRoot(cwd)
	if err != nil {
		return err
	}

	result, err := applinks.Clean(root, opts.apply)
	if err != nil {
		return err
	}
	if code := applinks.RenderClean(os.Stdout, result); code != 0 {
		return exitcode.Error{Code: code}
	}
	return nil
}
