package searchembed

import (
	"context"
	"fmt"
	"path/filepath"

	indexapp "github.com/sentzunhat/hawp/librarian/src/internal/application/index"
	"github.com/sentzunhat/hawp/librarian/src/internal/bootstrap"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	sqlite "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repositories/index"
)

func Run(args []string, cwd string) error {
	opts, err := parseSearchEmbedArgs(args)
	if err != nil {
		return err
	}

	root, err := repo.FindBacklogRepoRoot(cwd)
	if err != nil {
		fmt.Println("Not in a HAWP repo; nothing to embed.")
		return nil
	}

	dbPath := filepath.Join(root, ".hawp", "db", "index.sqlite")

	// Check if DB exists
	db, err := sqlite.Open(dbPath)
	if err != nil {
		fmt.Printf("Index not found at %s. Run `hawp search index` first.\n", dbPath)
		return nil
	}

	// Check how many chunks need embedding
	needEmbed, err := db.ChunksNeedEmbedding()
	if err != nil {
		db.Close()
		return fmt.Errorf("check embeddings: %w", err)
	}
	db.Close()

	if needEmbed == 0 {
		fmt.Println("All chunks already embedded.")
		return nil
	}

	modelID := opts.model
	if modelID == "" {
		switch opts.backend {
		case "onnx":
			modelID = indexapp.DefaultEmbeddingModel
		case "ollama":
			modelID = "nomic-embed-text"
		}
	}

	fmt.Printf("Embedding %d chunks with %s (%s)...\n", needEmbed, modelID, opts.backend)
	if opts.backend == "onnx" {
		fmt.Printf("Estimated time: %.0f seconds\n\n", float64(needEmbed)*0.008)
	}

	service := bootstrap.NewEmbedService(dbPath)
	result, err := service.Execute(context.Background(), opts.backend, modelID)
	if err != nil {
		return err
	}

	fmt.Print(result.String())
	fmt.Println("Vectors ready for hybrid search. Try: hawp search <query>")
	return nil
}
