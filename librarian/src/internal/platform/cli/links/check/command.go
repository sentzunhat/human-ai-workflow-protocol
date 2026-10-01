package check

import (
	"os"

	applinks "github.com/sentzunhat/hawp/librarian/src/internal/application/links"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/exitcode"
)

func Run(cwd string) error {
	root, err := repo.FindBacklogRepoRoot(cwd)
	if err != nil {
		return err
	}
	result := applinks.Check(root)
	if code := applinks.Render(os.Stdout, os.Stderr, result); code != 0 {
		return exitcode.Error{Code: code}
	}
	return nil
}
