package ingest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	appindex "github.com/sentzunhat/hawp/librarian/src/internal/application/index"
	appsearch "github.com/sentzunhat/hawp/librarian/src/internal/application/search"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
)

func Run(cwd string) error {
	root, err := repo.FindBacklogRepoRoot(cwd)
	if err != nil {
		fmt.Println("Not in a HAWP repo; nothing to index.")
		return nil
	}

	home, err := os.UserHomeDir()
	hawpHome := ""
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not resolve home directory: %v\n", err)
	} else if home != "" {
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

	// Path to the index DB. Validate repository containment before the
	// ingest service can create the SQLite file, WAL, or runtime directory.
	dbPath, err := filesystem.ResolveSafeSearchIndexPath(root)
	if err != nil {
		return err
	}

	service := appindex.NewIngestService(dbPath)
	result, err := service.Execute(corpus)
	if err != nil {
		return fmt.Errorf("ingest search index: %w", err)
	}

	fmt.Print(result.String())
	fmt.Printf("Index ready at: %s\n", dbPath)
	fmt.Println(`Try: hawp search query "<query>"`)
	return nil
}
