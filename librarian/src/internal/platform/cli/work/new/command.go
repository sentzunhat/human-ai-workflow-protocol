package newcmd

import (
	"fmt"
	"path/filepath"

	appuuid "github.com/sentzunhat/hawp/librarian/src/internal/application/uuidgen"
	appwork "github.com/sentzunhat/hawp/librarian/src/internal/application/work/intake"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
)

// Run scaffolds a work item. Investigation and planning remain human/agent
// work; this command only records the mechanical intake.
func Run(args []string, cwd string) error {
	opts, err := parseNewArgs(args)
	if err != nil {
		return err
	}

	workDir := filepath.Join(opts.hawpRoot, "work")
	if opts.hawpRoot == "" {
		root, err := repo.FindBacklogRepoRoot(cwd)
		if err != nil {
			return err
		}
		workDir = filepath.Join(root, ".hawp", "work")
	}
	if !repo.Exists(workDir) {
		return fmt.Errorf("could not resolve .hawp/work directory: %s", workDir)
	}

	result, err := appwork.NewItem(workDir, opts.itemType, opts.title, opts.input)
	if err != nil {
		return err
	}

	fmt.Printf("Created work item %s (%s)\n", appuuid.Short(result.UUID), result.Type)
	fmt.Printf("  Plan file: %s\n", result.PlanFilePath)
	fmt.Printf("  Backlog row added: %s\n", result.BacklogPath)
	fmt.Println()
	fmt.Println("Next: investigate and fill in the plan file (see .hawp/kit/usage/intake-workflow.md Step 2),")
	fmt.Println("then move the backlog status to analyzing/plan-ready as you go.")
	return nil
}
