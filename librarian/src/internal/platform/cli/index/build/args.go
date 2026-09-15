package build

import (
	"flag"
	"fmt"
	"io"

	domainindex "github.com/sentzunhat/hawp/librarian/src/internal/domain/index"
)

type indexBuildOptions struct {
	scope      domainindex.DocumentScope
	exportPath string
}

func parseIndexBuildArgs(args []string) (indexBuildOptions, error) {
	opts := indexBuildOptions{scope: domainindex.ScopeAll}

	flags := flag.NewFlagSet("index build", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var scopeStr string
	flags.StringVar(&scopeStr, "scope", string(domainindex.ScopeAll), "all|work|kit")
	flags.StringVar(&opts.exportPath, "export", "", "export path for results")
	flags.Bool("no-update-check", false, "suppress update notice")

	if err := flags.Parse(args); err != nil {
		return indexBuildOptions{}, fmt.Errorf("index build arguments: %w", err)
	}
	if flags.NArg() != 0 {
		return indexBuildOptions{}, fmt.Errorf("unexpected index build argument %q", flags.Arg(0))
	}

	opts.scope = domainindex.DocumentScope(scopeStr)
	switch opts.scope {
	case domainindex.ScopeAll, domainindex.ScopeWork, domainindex.ScopeKit:
	default:
		return indexBuildOptions{}, fmt.Errorf("unknown --scope %q (want all|work|kit)", scopeStr)
	}

	return opts, nil
}
