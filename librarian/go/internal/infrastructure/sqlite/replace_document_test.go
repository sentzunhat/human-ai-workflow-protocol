package sqlite

import (
	"path/filepath"
	"testing"
)

// TestReplaceDocumentIsAtomic proves ReplaceDocument upserts the document,
// clears previous chunks, and inserts replacement chunks in one transaction.
func TestReplaceDocumentIsAtomic(t *testing.T) {
	db := openTestDB(t)

	fc := "ctx"
	chunks1 := []Chunk{
		{ChunkIdx: 0, Text: "first chunk A", FolderContext: &fc},
		{ChunkIdx: 1, Text: "first chunk B", FolderContext: &fc},
	}

	id1, err := db.ReplaceDocument(
		DocumentRow{Category: "kit", Type: "guide", Path: "/doc.md", FolderRole: "start-here"},
		nil,
		chunks1,
	)
	if err != nil {
		t.Fatalf("ReplaceDocument() first call error = %v", err)
	}
	if id1 == 0 {
		t.Fatal("ReplaceDocument() returned id == 0")
	}

	var chunkCount int
	if err := db.db.QueryRow(`SELECT COUNT(*) FROM chunks WHERE document_id = ?`, id1).Scan(&chunkCount); err != nil {
		t.Fatalf("count chunks: %v", err)
	}
	if chunkCount != 2 {
		t.Errorf("after first ReplaceDocument: chunk count = %d, want 2", chunkCount)
	}

	// Re-ingest with fewer chunks — stale ones must be cleared.
	chunks2 := []Chunk{
		{ChunkIdx: 0, Text: "replacement chunk", FolderContext: &fc},
	}
	id2, err := db.ReplaceDocument(
		DocumentRow{Category: "kit", Type: "guide", Path: "/doc.md", FolderRole: "start-here"},
		nil,
		chunks2,
	)
	if err != nil {
		t.Fatalf("ReplaceDocument() second call error = %v", err)
	}
	if id2 != id1 {
		t.Errorf("ReplaceDocument() returned different ID on re-ingest: first=%d second=%d", id1, id2)
	}

	if err := db.db.QueryRow(`SELECT COUNT(*) FROM chunks WHERE document_id = ?`, id1).Scan(&chunkCount); err != nil {
		t.Fatalf("count chunks after replace: %v", err)
	}
	if chunkCount != 1 {
		t.Errorf("after second ReplaceDocument: chunk count = %d, want 1 (stale chunk cleared)", chunkCount)
	}

	var docCount int
	if err := db.db.QueryRow(`SELECT COUNT(*) FROM documents WHERE path = ?`, "/doc.md").Scan(&docCount); err != nil {
		t.Fatalf("count documents: %v", err)
	}
	if docCount != 1 {
		t.Errorf("documents table has %d rows after re-ingest, want 1 (no duplicates)", docCount)
	}
}

// TestReplaceDocumentRollsBackOnFailure proves that a failure mid-operation
// (duplicate chunk_idx in the replacement set) leaves the document and its
// previous chunks untouched — no partial state is written.
func TestReplaceDocumentRollsBackOnFailure(t *testing.T) {
	db := openTestDB(t)

	fc := "ctx"
	original := []Chunk{
		{ChunkIdx: 0, Text: "original chunk", FolderContext: &fc},
	}

	id, err := db.ReplaceDocument(
		DocumentRow{Category: "kit", Type: "guide", Path: "/doc.md", FolderRole: "start-here"},
		nil,
		original,
	)
	if err != nil {
		t.Fatalf("first ReplaceDocument() error = %v", err)
	}

	// Attempt a replace with duplicate chunk_idx, which violates UNIQUE(document_id, chunk_idx).
	bad := []Chunk{
		{ChunkIdx: 0, Text: "new chunk A", FolderContext: &fc},
		{ChunkIdx: 0, Text: "duplicate idx — should fail", FolderContext: &fc}, // conflict
	}
	if _, err := db.ReplaceDocument(
		DocumentRow{Category: "kit", Type: "guide", Path: "/doc.md", FolderRole: "start-here"},
		nil,
		bad,
	); err == nil {
		t.Fatal("ReplaceDocument() with duplicate chunk_idx should return an error")
	}

	// The previous state must be intact.
	var chunkCount int
	if err := db.db.QueryRow(`SELECT COUNT(*) FROM chunks WHERE document_id = ?`, id).Scan(&chunkCount); err != nil {
		t.Fatalf("count chunks after failed replace: %v", err)
	}
	if chunkCount != 1 {
		t.Errorf("after failed ReplaceDocument: chunk count = %d, want 1 (rollback should preserve original chunk)", chunkCount)
	}

	var text string
	if err := db.db.QueryRow(`SELECT text FROM chunks WHERE document_id = ?`, id).Scan(&text); err != nil {
		t.Fatalf("read chunk text after failed replace: %v", err)
	}
	if text != "original chunk" {
		t.Errorf("chunk text after failed replace = %q, want %q (rollback should restore original text)", text, "original chunk")
	}
}

// TestReplaceDocumentWithMetadata proves that work-item metadata is upserted
// alongside the document in the same atomic operation.
func TestReplaceDocumentWithMetadata(t *testing.T) {
	db := openTestDB(t)

	status := "in-progress"
	uuid := "abc123"
	meta := &DocumentMetadata{
		WorkUUID: uuid,
		Status:   status,
	}

	id, err := db.ReplaceDocument(
		DocumentRow{Category: "work", Type: "task", Path: "/work/active/abc123.md", FolderRole: "active"},
		meta,
		nil,
	)
	if err != nil {
		t.Fatalf("ReplaceDocument() with metadata error = %v", err)
	}

	var gotUUID, gotStatus string
	if err := db.db.QueryRow(`SELECT work_uuid, status FROM documents_metadata WHERE document_id = ?`, id).Scan(&gotUUID, &gotStatus); err != nil {
		t.Fatalf("read metadata: %v", err)
	}
	if gotUUID != uuid || gotStatus != status {
		t.Errorf("metadata = {%q, %q}, want {%q, %q}", gotUUID, gotStatus, uuid, status)
	}

	// Re-ingest with updated status — must upsert, not duplicate.
	meta.Status = "done"
	if _, err := db.ReplaceDocument(
		DocumentRow{Category: "work", Type: "task", Path: "/work/active/abc123.md", FolderRole: "active"},
		meta,
		nil,
	); err != nil {
		t.Fatalf("ReplaceDocument() re-ingest with updated metadata error = %v", err)
	}

	var metaCount int
	if err := db.db.QueryRow(`SELECT COUNT(*) FROM documents_metadata WHERE document_id = ?`, id).Scan(&metaCount); err != nil {
		t.Fatalf("count metadata: %v", err)
	}
	if metaCount != 1 {
		t.Errorf("documents_metadata has %d rows, want 1 (upsert, not insert)", metaCount)
	}

	if err := db.db.QueryRow(`SELECT status FROM documents_metadata WHERE document_id = ?`, id).Scan(&gotStatus); err != nil {
		t.Fatalf("read updated status: %v", err)
	}
	if gotStatus != "done" {
		t.Errorf("status after re-ingest = %q, want %q", gotStatus, "done")
	}
}

// TestReplaceDocumentFTSSyncOnReplace proves the FTS5 triggers keep
// chunks_fts in sync when ReplaceDocument clears and re-inserts chunks.
func TestReplaceDocumentFTSSyncOnReplace(t *testing.T) {
	db := openTestDB(t)

	fc := "ctx"
	_, err := db.ReplaceDocument(
		DocumentRow{Category: "kit", Type: "guide", Path: filepath.Join(t.TempDir(), "doc.md"), FolderRole: "start-here"},
		nil,
		[]Chunk{{ChunkIdx: 0, Text: "uniqueftsmarker", FolderContext: &fc}},
	)
	if err != nil {
		t.Fatalf("ReplaceDocument() error = %v", err)
	}

	results, err := db.QueryChunksLexical("uniqueftsmarker", 5)
	if err != nil {
		t.Fatalf("QueryChunksLexical() error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("lexical search after ReplaceDocument returned %d results, want 1", len(results))
	}
}
