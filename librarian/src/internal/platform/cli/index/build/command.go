package build

import (
	"fmt"
	appindex "github.com/sentzunhat/hawp/librarian/src/internal/application/index"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
)

func Run(args []string, cwd string) error {
	opts, err := parseIndexBuildArgs(args)
	if err != nil {
		return err
	}

	root, err := repo.FindBacklogRepoRoot(cwd)
	if err != nil {
		return err
	}
	service := appindex.NewBuildService(root)
	result, err := service.Execute(opts.scope)
	if err != nil {
		return err
	}
	fmt.Print(result.String())

	if opts.exportPath != "" {
		if err := result.Export(opts.exportPath); err != nil {
			return err
		}
		fmt.Printf("\nExported %d document(s) to %s\n", len(result.Documents), opts.exportPath)
	}
	return nil
}
