// Package doccmd implements hawp work status|evidence|decision|note.
package doccmd

import (
	"fmt"
	"path/filepath"

	appdoc "github.com/sentzunhat/hawp/librarian/src/internal/application/work/doc"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
)

// Run creates a secondary work document of the given docType.
func Run(docType string, args []string, cwd string) error {
	opts, err := parseDocArgs(docType, args)
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

	result, err := appdoc.CreateWorkDoc(docType, opts.title, workDir, opts.workItemID)
	if err != nil {
		return err
	}

	fmt.Printf("Created %s document %s\n", result.DocType, result.UUID)
	fmt.Printf("  File: %s\n", result.FilePath)
	fmt.Println()
	fmt.Println("Open the file and fill in the content.")
	return nil
}
