package search

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/sqlite"
)

func TestQueryErrorsWithoutIndex(t *testing.T) {
	// A fresh repo root with no .hawp/db/index.sqlite at all.
	_, err := Query(t.TempDir(), "anything", 5)
	if err == nil {
		t.Error("Query should error when the index doesn't exist yet")
	}
}

func TestQueryLexicalOnly(t *testing.T) {
	repoRoot := t.TempDir()
	dbPath := filepath.Join(repoRoot, ".hawp", "db", "index.sqlite")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
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

	// No vectors embedded — this should still return lexical results
	// (no index_metadata row, so HybridRank falls back gracefully).
	results, err := Query(repoRoot, "kubernetes", 5)
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Query() returned %d results, want 1", len(results))
	}
	if results[0].Source != "/test/kubernetes.md" {
		t.Errorf("Source = %q, want /test/kubernetes.md", results[0].Source)
	}
}

func TestQueryNoMatches(t *testing.T) {
	repoRoot := t.TempDir()
	dbPath := filepath.Join(repoRoot, ".hawp", "db", "index.sqlite")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
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

	results, err := Query(repoRoot, "nonexistent query terms", 5)
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Query() on an empty index should return no results, got %d", len(results))
	}
}
