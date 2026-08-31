package search_test

import (
	"os"
	"path/filepath"
	"testing"

	appsearch "github.com/sentzunhat/hawp/librarian/src/internal/application/search"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/sqlite"
)

func TestQueryErrorsWithoutIndex(t *testing.T) {
	_, err := appsearch.Query(t.TempDir(), "anything", 5)
	if err == nil {
		t.Error("Query should error when the index doesn't exist yet")
	}
}

func TestQueryLexicalOnly(t *testing.T) {
	repoRoot := t.TempDir()
	dbPath := filepath.Join(repoRoot, ".hawp", "db", "index.sqlite")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		t.Fatal(err)
	}

	db, err := sqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if err := db.InitSchema(); err != nil {
		t.Fatalf("InitSchema() error = %v", err)
	}
	docID, err := db.InsertDocument("kit", "guide", "/test/kubernetes.md", "start-here")
	if err != nil {
		t.Fatalf("InsertDocument() error = %v", err)
	}
	if err := db.InsertChunk(sqlite.Chunk{
		DocumentID: docID,
		ChunkIdx:   0,
		Text:       "Kubernetes manages containerized applications across clusters.",
	}); err != nil {
		t.Fatalf("InsertChunk() error = %v", err)
	}
	db.Close()

	results, err := appsearch.Query(repoRoot, "kubernetes", 5)
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Query() returned %d results, want 1", len(results))
	}
	if results[0].Source != "/test/kubernetes.md" {
		t.Errorf("Source = %q, want /test/kubernetes.md", results[0].Source)
	}
	if results[0].ChunkID == "" {
		t.Error("ChunkID should not be empty")
	}
	if results[0].Priority != 0 {
		t.Errorf("Priority = %d, want 0", results[0].Priority)
	}
}

func TestQueryNoMatches(t *testing.T) {
	repoRoot := t.TempDir()
	dbPath := filepath.Join(repoRoot, ".hawp", "db", "index.sqlite")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		t.Fatal(err)
	}
	db, err := sqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if err := db.InitSchema(); err != nil {
		t.Fatalf("InitSchema() error = %v", err)
	}
	db.Close()

	results, err := appsearch.Query(repoRoot, "nonexistent query terms", 5)
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Query() on an empty index should return no results, got %d", len(results))
	}
}
