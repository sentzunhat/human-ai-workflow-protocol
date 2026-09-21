package normalizecmd

import (
	appwork "github.com/sentzunhat/hawp/librarian/src/internal/application/work/normalize"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/exitcode"
	"os"
)

// Run normalizes a work tree resolved from cwd when not explicit.
func Run(args []string, cwd string) error {
	opts, explicitRepoRoot, err := parseNormalizeArgs(args)
	if err != nil {
		return err
	}
	if explicitRepoRoot != "" {
		opts.RepoRoot = explicitRepoRoot
	} else {
		root, err := repo.FindBacklogRepoRoot(cwd)
		if err != nil {
			return err
		}
		opts.RepoRoot = root
	}
	if code := appwork.Normalize(os.Stdout, os.Stderr, opts); code != 0 {
		return exitcode.Error{Code: code}
	}
	return nil
}
