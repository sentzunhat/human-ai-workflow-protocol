package ingest

import (
	"fmt"
	appindex "github.com/sentzunhat/hawp/librarian/src/internal/application/index"
	appsearch "github.com/sentzunhat/hawp/librarian/src/internal/application/search"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	"os"
	"path/filepath"
	"strings"
)

func Run(cwd string) error {
	root, err := repo.FindBacklogRepoRoot(cwd)
	if err != nil {
		fmt.Println("Not in a HAWP repo; nothing to index.")
		return nil
	}

	home, _ := os.UserHomeDir()
	hawpHome := ""
	if home != "" {
		hawpHome = filepath.Join(home, ".hawp")
	}
	searchCfg, err := appsearch.LoadSearchConfig(hawpHome, root)
	if err != nil {
		return fmt.Errorf("load search config: %w", err)
	}

	fmt.Printf("Building enriched document index from: %s\n", strings.Join(searchCfg.Index.Paths, ", "))

	// Build corpus from actual files
	corpus, err := buildCorpusFromRepo(root, searchCfg.Index.Paths)
	if err != nil {
		return fmt.Errorf("build corpus: %w", err)
	}

	// Path to the index DB
	dbPath := filepath.Join(root, ".hawp", "db", "index.sqlite")

	service := appindex.NewIngestService(dbPath)
	result, err := service.Execute(corpus)
	if err != nil {
		return err
	}

	fmt.Print(result.String())
	fmt.Printf("Index ready at: %s\n", dbPath)
	fmt.Println("Try: hawp search vector")
	return nil
}
