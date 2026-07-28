package sqlite

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGetChunkVector(t *testing.T) {
	// Create temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer db.Close()

	// Initialize schema
	if err := db.InitSchema(); err != nil {
		t.Fatalf("InitSchema() error = %v", err)
	}

	// Insert a document
	docID, err := db.InsertDocument("kit", "guide", "/test/doc.md", "start-here")
	if err != nil {
		t.Fatalf("InsertDocument() error = %v", err)
	}

	// Insert a chunk
	chunk := Chunk{
		DocumentID:   docID,
		ChunkIdx:     0,
		Text:         "sample text",
		FolderContext: nil,
		MetadataJSON: nil,
	}
	if err := db.InsertChunk(chunk); err != nil {
		t.Fatalf("InsertChunk() error = %v", err)
	}

	// Create test vector
	testVector := []float32{0.1, 0.2, 0.3, 0.4}
	vectorJSON, _ := json.Marshal(testVector)

	// Store embedding
	if err := db.UpdateChunkEmbedding(1, vectorJSON); err != nil {
		t.Fatalf("UpdateChunkEmbedding() error = %v", err)
	}

	// Retrieve vector
	got, err := db.GetChunkVector(1)
	if err != nil {
		t.Fatalf("GetChunkVector() error = %v", err)
	}

	if len(got) != len(testVector) {
		t.Errorf("GetChunkVector() length = %d, want %d", len(got), len(testVector))
	}

	for i, v := range got {
		if v != testVector[i] {
			t.Errorf("GetChunkVector()[%d] = %v, want %v", i, v, testVector[i])
		}
	}
}

func TestGetChunkVectorNotFound(t *testing.T) {
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

	// Try to get non-existent vector
	_, err = db.GetChunkVector(999)
	if err == nil {
		t.Error("GetChunkVector() should error for non-existent chunk")
	}
}

func TestGetAllChunkVectors(t *testing.T) {
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

	// Insert a document
	docID, err := db.InsertDocument("kit", "guide", "/test/doc.md", "start-here")
	if err != nil {
		t.Fatalf("InsertDocument() error = %v", err)
	}

	// Insert multiple chunks with embeddings
	vectors := [][]float32{
		{0.1, 0.2, 0.3},
		{0.4, 0.5, 0.6},
		{0.7, 0.8, 0.9},
	}

	for i, vec := range vectors {
		chunk := Chunk{
			DocumentID:   docID,
			ChunkIdx:     i,
			Text:         "sample",
			FolderContext: nil,
			MetadataJSON: nil,
		}
		if err := db.InsertChunk(chunk); err != nil {
			t.Fatalf("InsertChunk() error = %v", err)
		}

		vectorJSON, _ := json.Marshal(vec)
		if err := db.UpdateChunkEmbedding(int64(i+1), vectorJSON); err != nil {
			t.Fatalf("UpdateChunkEmbedding() error = %v", err)
		}
	}

	// Retrieve all vectors
	got, err := db.GetAllChunkVectors()
	if err != nil {
		t.Fatalf("GetAllChunkVectors() error = %v", err)
	}

	if len(got) != len(vectors) {
		t.Errorf("GetAllChunkVectors() returned %d vectors, want %d", len(got), len(vectors))
	}

	// Verify each vector
	for i, expectedVec := range vectors {
		chunkID := int64(i + 1)
		if gotVec, ok := got[chunkID]; ok {
			for j, v := range gotVec {
				if v != expectedVec[j] {
					t.Errorf("vector[%d][%d] = %v, want %v", i, j, v, expectedVec[j])
				}
			}
		} else {
			t.Errorf("missing vector for chunk %d", chunkID)
		}
	}
}

func TestIntegrationEmbedding(t *testing.T) {
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

	// Simulate ingest process
	docID, _ := db.InsertDocument("kit", "guide", "/test/doc.md", "start-here")
	for i := 0; i < 5; i++ {
		chunk := Chunk{
			DocumentID:   docID,
			ChunkIdx:     i,
			Text:         "sample text for embedding",
			FolderContext: nil,
			MetadataJSON: nil,
		}
		db.InsertChunk(chunk)
	}

	// Check chunks needing embedding
	needEmbed, err := db.ChunksNeedEmbedding()
	if err != nil {
		t.Fatalf("ChunksNeedEmbedding() error = %v", err)
	}

	if needEmbed != 5 {
		t.Errorf("ChunksNeedEmbedding() = %d, want 5", needEmbed)
	}

	// Simulate embedding process
	chunks, err := db.GetChunksNeedingEmbedding()
	if err != nil {
		t.Fatalf("GetChunksNeedingEmbedding() error = %v", err)
	}

	if len(chunks) != 5 {
		t.Errorf("GetChunksNeedingEmbedding() returned %d chunks, want 5", len(chunks))
	}

	// Add embeddings
	for _, chunk := range chunks {
		vec := make([]float32, 384)
		for i := range vec {
			vec[i] = float32(i) / 384.0
		}
		vecJSON, _ := json.Marshal(vec)
		db.UpdateChunkEmbedding(chunk.ID, vecJSON)
	}

	// Verify all embeddings stored
	afterEmbed, err := db.ChunksNeedEmbedding()
	if err != nil {
		t.Fatalf("ChunksNeedEmbedding() error = %v", err)
	}

	if afterEmbed != 0 {
		t.Errorf("ChunksNeedEmbedding() after embedding = %d, want 0", afterEmbed)
	}
}

// TestInsertDocumentIsIdempotent proves re-ingesting the same path twice is
// safe: previously, InsertDocument was a plain INSERT and `path` is UNIQUE,
// so a second run over an already-indexed document failed with a
// constraint violation instead of updating it in place.
func TestInsertDocumentIsIdempotent(t *testing.T) {
	tmpDir := t.TempDir()
	db, err := Open(filepath.Join(tmpDir, "test.db"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer db.Close()
	if err := db.InitSchema(); err != nil {
		t.Fatalf("InitSchema() error = %v", err)
	}

	firstID, err := db.InsertDocument("kit", "guide", "/test/doc.md", "start-here")
	if err != nil {
		t.Fatalf("first InsertDocument() error = %v", err)
	}

	secondID, err := db.InsertDocument("kit", "guide-updated", "/test/doc.md", "start-here")
	if err != nil {
		t.Fatalf("second InsertDocument() (re-ingest) error = %v, want no error", err)
	}

	if firstID != secondID {
		t.Errorf("re-ingesting the same path returned a different id: first=%d second=%d", firstID, secondID)
	}

	var typ string
	if err := db.db.QueryRow(`SELECT type FROM documents WHERE id = ?`, secondID).Scan(&typ); err != nil {
		t.Fatalf("query updated document: %v", err)
	}
	if typ != "guide-updated" {
		t.Errorf("re-ingest did not update fields: type = %q, want %q", typ, "guide-updated")
	}

	var count int
	if err := db.db.QueryRow(`SELECT COUNT(*) FROM documents WHERE path = ?`, "/test/doc.md").Scan(&count); err != nil {
		t.Fatalf("count documents: %v", err)
	}
	if count != 1 {
		t.Errorf("re-ingest created a duplicate row: count = %d, want 1", count)
	}
}

// TestDeleteChunksForDocumentClearsFTS proves the FTS5 delete trigger keeps
// chunks_fts in sync when chunks are cleared for re-ingest.
func TestDeleteChunksForDocumentClearsFTS(t *testing.T) {
	tmpDir := t.TempDir()
	db, err := Open(filepath.Join(tmpDir, "test.db"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer db.Close()
	if err := db.InitSchema(); err != nil {
		t.Fatalf("InitSchema() error = %v", err)
	}

	docID, err := db.InsertDocument("kit", "guide", "/test/doc.md", "start-here")
	if err != nil {
		t.Fatalf("InsertDocument() error = %v", err)
	}
	if err := db.InsertChunk(Chunk{DocumentID: docID, ChunkIdx: 0, Text: "uniquemarkertext"}); err != nil {
		t.Fatalf("InsertChunk() error = %v", err)
	}

	if err := db.DeleteChunksForDocument(docID); err != nil {
		t.Fatalf("DeleteChunksForDocument() error = %v", err)
	}

	var chunkCount int
	if err := db.db.QueryRow(`SELECT COUNT(*) FROM chunks WHERE document_id = ?`, docID).Scan(&chunkCount); err != nil {
		t.Fatalf("count chunks: %v", err)
	}
	if chunkCount != 0 {
		t.Errorf("chunks not deleted: count = %d, want 0", chunkCount)
	}

	var ftsCount int
	if err := db.db.QueryRow(`SELECT COUNT(*) FROM chunks_fts WHERE chunks_fts MATCH 'uniquemarkertext'`).Scan(&ftsCount); err != nil {
		t.Fatalf("query chunks_fts: %v", err)
	}
	if ftsCount != 0 {
		t.Errorf("chunks_fts still has stale entry after delete: count = %d, want 0", ftsCount)
	}
}
