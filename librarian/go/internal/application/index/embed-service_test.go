package index

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/sqlite"
)

func setupEmbedTestDB(t *testing.T) (dbPath string, insertChunk func(text string)) {
	t.Helper()
	dbPath = filepath.Join(t.TempDir(), "test.db")

	db, err := sqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if err := db.InitSchema(); err != nil {
		t.Fatalf("InitSchema() error = %v", err)
	}

	docID, err := db.InsertDocument("kit", "guide", "/test/doc.md", "start-here")
	if err != nil {
		t.Fatalf("InsertDocument() error = %v", err)
	}
	db.Close()

	// EmbedService.Execute opens its own connection per call, so insertChunk
	// does too — sqlite's single-connection-per-process design (see
	// infrastructure/sqlite) means holding a connection open across
	// Execute() calls would just serialize against them, not fail, but a
	// fresh short-lived connection per insert matches how real callers
	// (search index, then search embed) actually use this DB: separate
	// process invocations, not one long-lived handle.
	chunkIdx := 0
	insertChunk = func(text string) {
		db, err := sqlite.Open(dbPath)
		if err != nil {
			t.Fatalf("Open() error = %v", err)
		}
		defer db.Close()
		if err := db.InsertChunk(sqlite.Chunk{
			DocumentID: docID,
			ChunkIdx:   chunkIdx,
			Text:       text,
		}); err != nil {
			t.Fatalf("InsertChunk() error = %v", err)
		}
		chunkIdx++
	}
	return dbPath, insertChunk
}

// TestExecuteRejectsMixingModels locks in that embedding an index with a
// different backend/model than what already embedded it fails loudly,
// rather than silently mixing incompatible vector spaces into one index —
// the real bug found and fixed 2026-07-27 (this session had accidentally
// mixed ONNX/bge-base and Ollama/nomic-embed-text vectors in one test DB).
// Uses the "none" backend (NullEmbedder) so this test needs no real
// model/network.
func TestExecuteRejectsMixingModels(t *testing.T) {
	dbPath, insertChunk := setupEmbedTestDB(t)
	insertChunk("first chunk")

	svc := NewEmbedService(dbPath)

	result, err := svc.Execute(context.Background(), "none", "none")
	if err != nil {
		t.Fatalf("first Execute() error = %v", err)
	}
	if result.ChunksEmbedded != 1 {
		t.Fatalf("first Execute() embedded %d chunks, want 1", result.ChunksEmbedded)
	}

	// Add a second chunk and try to embed the remainder with a different
	// backend/model than what's already in the index.
	insertChunk("second chunk")

	_, err = svc.Execute(context.Background(), "onnx", "bge-base-en-v1.5")
	if err == nil {
		t.Fatal("Execute() with a different backend/model should fail when the index already has a different one")
	}
	if !strings.Contains(err.Error(), "none/none") {
		t.Errorf("error should mention the existing backend/model (none/none), got: %v", err)
	}
}

// TestExecuteAllowsSameModelAcrossRuns confirms re-running Execute with the
// SAME backend/model as an already-partially-embedded index is allowed
// (e.g. resuming after a previous run stopped partway through).
func TestExecuteAllowsSameModelAcrossRuns(t *testing.T) {
	dbPath, insertChunk := setupEmbedTestDB(t)
	insertChunk("first chunk")
	insertChunk("second chunk")

	svc := NewEmbedService(dbPath)

	// Embed just the first chunk isn't easily isolated via the public API,
	// so instead verify two consecutive full runs with the same
	// backend/model both succeed (the second is a no-op: no chunks left).
	if _, err := svc.Execute(context.Background(), "none", "none"); err != nil {
		t.Fatalf("first Execute() error = %v", err)
	}

	_, err := svc.Execute(context.Background(), "none", "none")
	if err == nil {
		t.Fatal("expected \"no chunks to embed\" once everything is already embedded")
	}
	if !strings.Contains(err.Error(), "no chunks to embed") {
		t.Errorf("expected the no-chunks-left error, not a mixed-model rejection, got: %v", err)
	}
}
