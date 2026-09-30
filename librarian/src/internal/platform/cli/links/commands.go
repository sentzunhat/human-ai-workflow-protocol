package links

import (
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/links/check"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/links/clean"
)

func RunCheck(cwd string) error                { return check.Run(cwd) }
func RunClean(args []string, cwd string) error { return clean.Run(args, cwd) }
