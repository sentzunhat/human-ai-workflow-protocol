package distribution

import (
	"fmt"
	"os"

	appdistribution "github.com/sentzunhat/hawp/librarian/src/internal/application/distribution"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/providers"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/exitcode"
)

func RunBuild(cwd string) error {
	root, err := repo.FindBacklogRepoRoot(cwd)
	if err != nil {
		return err
	}
	result, err := appdistribution.Build(root)
	if err != nil {
		return err
	}
	for _, file := range result.Updated {
		fmt.Printf("updated %s\n", file)
	}
	for _, file := range result.Unchanged {
		fmt.Printf("unchanged %s\n", file)
	}
	fmt.Printf("\ndistribution build complete: %d/%d file(s) updated\n", len(result.Updated), len(result.Updated)+len(result.Unchanged))
	return nil
}

func RunValidate(cwd string) error {
	root, err := repo.FindBacklogRepoRoot(cwd)
	if err != nil {
		return err
	}
	result, err := appdistribution.Validate(root)
	if err != nil {
		return err
	}
	if !result.OK() {
		fmt.Fprintln(os.Stderr, "distribution validation failed")
		if len(result.LegacyPresent) > 0 {
			fmt.Fprintln(os.Stderr, "\nlegacy root-level guides must be removed:")
			for _, file := range result.LegacyPresent {
				fmt.Fprintf(os.Stderr, "- %s\n", file)
			}
		}
		if len(result.Missing) > 0 {
			fmt.Fprintln(os.Stderr, "\nmissing generated outputs:")
			for _, file := range result.Missing {
				fmt.Fprintf(os.Stderr, "- %s\n", file)
			}
		}
		if len(result.Stale) > 0 {
			fmt.Fprintln(os.Stderr, "\nstale generated outputs:")
			for _, file := range result.Stale {
				fmt.Fprintf(os.Stderr, "- %s\n", file)
			}
		}
		if len(result.PathLeaks) > 0 {
			fmt.Fprintln(os.Stderr, "\ninvalid downstream path leaks (core/.hawp/) in source kit files:")
			for _, leak := range result.PathLeaks {
				fmt.Fprintf(os.Stderr, "- %s:%d %s\n", leak.File, leak.Line, leak.Text)
			}
		}
		fmt.Fprintln(os.Stderr, "\nrun `hawp distribution build`")
		return exitcode.Error{Code: 1}
	}
	fmt.Println("distribution validation passed: generated outputs are current")
	return nil
}

func RunSync(cwd string) error {
	if err := providers.RunSync(cwd); err != nil {
		return err
	}
	if err := RunBuild(cwd); err != nil {
		return err
	}
	return RunValidate(cwd)
}
