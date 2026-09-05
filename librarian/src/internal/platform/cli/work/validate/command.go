package validatecmd

import (
	"fmt"
	appwork "github.com/sentzunhat/hawp/librarian/src/internal/application/work/validation"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/exitcode"
	"os"
	"path/filepath"
)

// Run validates a work directory resolved from cwd when not explicit.
func Run(args []string, cwd string) error {
	workDir, err := parseValidateArgs(args)
	if err != nil {
		return err
	}
	if workDir == "" {
		root, err := repo.FindBacklogRepoRoot(cwd)
		if err != nil {
			return err
		}
		workDir = filepath.Join(root, ".hawp", "work")
	}
	if !repo.Exists(workDir) {
		return fmt.Errorf("could not resolve .hawp/work directory: %s", workDir)
	}
	report, err := appwork.Validate(workDir)
	if err != nil {
		return err
	}
	if code := appwork.Render(os.Stdout, workDir, report); code != 0 {
		return exitcode.Error{Code: code}
	}
	return nil
}
