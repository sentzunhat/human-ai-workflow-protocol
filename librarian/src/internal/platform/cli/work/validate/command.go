package validatecmd

import (
	"fmt"
	"os"
	"path/filepath"

	appwork "github.com/sentzunhat/hawp/librarian/src/internal/application/work/validation"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/exitcode"
)

// Run validates a work directory resolved from cwd when not explicit.
func Run(args []string, cwd string) error {
	workDir, err := parseValidateArgs(args)
	if err != nil {
		return err
	}

	explicit := workDir != ""
	root := ""
	if !explicit {
		root, err = repo.FindBacklogRepoRoot(cwd)
		if err != nil {
			return err
		}
		workDir = filepath.Join(root, ".hawp", "work")
	} else {
		if !filepath.IsAbs(workDir) {
			workDir = filepath.Join(cwd, workDir)
		}
		workDir = filepath.Clean(workDir)
	}
	if !repo.Exists(workDir) {
		return fmt.Errorf("could not resolve .hawp/work directory: %s", workDir)
	}

	if explicit {
		if err := filesystem.RejectSymlinksInPath(workDir); err != nil {
			return fmt.Errorf("unsafe work validation path %s: %w", workDir, err)
		}
	} else if err := filesystem.RejectSymlinksInTree(root, workDir); err != nil {
		return fmt.Errorf("unsafe work validation path %s: %w", workDir, err)
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
