package validate

import (
	"fmt"
	"os"
	"path/filepath"

	appkit "github.com/sentzunhat/hawp/librarian/src/internal/application/kit"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/exitcode"
)

func Run(args []string, cwd string) error {
	o, err := parseKitValidateArgs(args)
	if err != nil {
		return err
	}

	p := o.kitPath
	explicit := p != ""
	if !explicit {
		root, err := repo.FindBacklogRepoRoot(cwd)
		if err != nil {
			return err
		}
		p = filepath.Join(root, ".hawp", "kit")
		if err := filesystem.RejectSymlinksInTree(root, p); err != nil {
			return fmt.Errorf("unsafe kit validation path %s: %w", p, err)
		}
	} else {
		if !filepath.IsAbs(p) {
			p = filepath.Join(cwd, p)
		}
		p = filepath.Clean(p)
		// An explicit kit root may intentionally live outside the repository,
		// but its full ancestor chain and descendants must still be symlink-free.
		if _, err := os.Lstat(p); err == nil {
			if err := filesystem.RejectSymlinksInPath(p); err != nil {
				return fmt.Errorf("unsafe explicit kit validation path %s: %w", p, err)
			}
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("inspect explicit kit validation path %s: %w", p, err)
		}
	}

	result := appkit.Validate(p)
	if code := appkit.Render(os.Stdout, os.Stderr, result); code != 0 {
		return exitcode.Error{Code: code}
	}
	return nil
}
