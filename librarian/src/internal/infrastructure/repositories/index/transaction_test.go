package sqlite

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestTransactionPersistence(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration test; set INTEGRATION=1")
	}

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer db.Close()

	if err := db.InitSchema(); err != nil {
		t.Fatalf("InitSchema() error = %v", err)
	}

	// Insert test document
	docID, err := db.InsertDocument("kit", "guide", "/test/doc.md", "start-here")
	if err != nil {
		t.Fatalf("InsertDocument() error = %v", err)
	}

	// Insert test chunks
	const numChunks = 100
	for i := 0; i < numChunks; i++ {
		chunk := Chunk{
			DocumentID:    docID,
			ChunkIdx:      i,
			Text:          "test content",
			FolderContext: nil,
			MetadataJSON:  nil,
		}
		if err := db.InsertChunk(chunk); err != nil {
			t.Fatalf("InsertChunk() error = %v", err)
		}
	}

	// Verify chunks needing embedding
	beforeCount, err := db.ChunksNeedEmbedding()
	if err != nil {
		t.Fatalf("ChunksNeedEmbedding() error = %v", err)
	}
	if beforeCount != numChunks {
		t.Errorf("ChunksNeedEmbedding() = %d, want %d", beforeCount, numChunks)
	}

	// Start transaction
	if err := db.BeginTx(); err != nil {
		t.Fatalf("BeginTx() error = %v", err)
	}

	// Write embeddings within transaction with periodic commits
	const commitFreq = 32
	totalEmbedded := 0

	for i := 0; i < numChunks; i++ {
		vec := make([]float32, 384)
		for j := range vec {
			vec[j] = float32(i) / float32(numChunks)
		}
		vectorJSON, _ := json.Marshal(vec)

		if err := db.UpdateChunkEmbedding(int64(i+1), vectorJSON); err != nil {
			db.Rollback()
			t.Fatalf("UpdateChunkEmbedding() error = %v", err)
		}
		totalEmbedded++

		// Commit periodically
		if totalEmbedded%commitFreq == 0 {
			if err := db.Commit(); err != nil {
				t.Fatalf("Commit() error = %v", err)
			}
			if err := db.BeginTx(); err != nil {
				t.Fatalf("BeginTx() restart error = %v", err)
			}
		}
	}

	// Final commit
	if err := db.Commit(); err != nil {
		t.Fatalf("Final Commit() error = %v", err)
	}

	// Verify persistence
	afterCount, err := db.ChunksNeedEmbedding()
	if err != nil {
		t.Fatalf("ChunksNeedEmbedding() after embedding error = %v", err)
	}

	if afterCount != 0 {
		t.Errorf("ChunksNeedEmbedding() after transaction = %d, want 0 (all vectors should be stored)", afterCount)
	}

	// Verify vectors are actually readable
	allVectors, err := db.GetAllChunkVectors()
	if err != nil {
		t.Fatalf("GetAllChunkVectors() error = %v", err)
	}

	if len(allVectors) != numChunks {
		t.Errorf("GetAllChunkVectors() returned %d vectors, want %d", len(allVectors), numChunks)
	}
}
