package indexcmd

import (
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/index/build"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/index/ingest"
)

func RunBuild(args []string, cwd string) error { return build.Run(args, cwd) }
func RunSearchIndex(cwd string) error          { return ingest.Run(cwd) }
