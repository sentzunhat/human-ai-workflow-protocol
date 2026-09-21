package providers

import (
	"fmt"
	"os"

	appprovidersync "github.com/sentzunhat/hawp/librarian/src/internal/application/providersync"
	domainprovidersync "github.com/sentzunhat/hawp/librarian/src/internal/domain/providersync"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/exitcode"
)

func RunMaterialize(cwd string) error {
	root, err := repo.FindBacklogRepoRoot(cwd)
	if err != nil {
		return err
	}
	result, err := appprovidersync.Materialize(root)
	if err != nil {
		return err
	}
	for _, file := range result.Updated {
		fmt.Printf("updated %s\n", file)
	}
	for _, file := range result.Unchanged {
		fmt.Printf("unchanged %s\n", file)
	}
	fmt.Printf("\nprovider materialize complete: %d/%d file(s) updated\n", len(result.Updated), len(result.Updated)+len(result.Unchanged))
	return nil
}

func RunValidate(cwd string) error {
	root, err := repo.FindBacklogRepoRoot(cwd)
	if err != nil {
		return err
	}
	result, err := appprovidersync.Validate(root)
	if err != nil {
		return err
	}
	if !result.OK() {
		fmt.Fprintln(os.Stderr, "provider materialization validation failed")
		if len(result.Missing) > 0 {
			fmt.Fprintln(os.Stderr, "\nmissing materialized outputs:")
			for _, file := range result.Missing {
				fmt.Fprintf(os.Stderr, "- %s\n", file)
			}
		}
		if len(result.Stale) > 0 {
			fmt.Fprintln(os.Stderr, "\nstale materialized outputs:")
			for _, file := range result.Stale {
				fmt.Fprintf(os.Stderr, "- %s\n", file)
			}
		}
		fmt.Fprintln(os.Stderr, "\nrun `hawp providers materialize`")
		return exitcode.Error{Code: 1}
	}
	fmt.Printf("provider validation passed: %d materialized file(s) are current\n", len(domainprovidersync.MaterializationTargets))
	return nil
}

func RunSync(cwd string) error {
	if err := RunMaterialize(cwd); err != nil {
		return err
	}
	return RunValidate(cwd)
}
