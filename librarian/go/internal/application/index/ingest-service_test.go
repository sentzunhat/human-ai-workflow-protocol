package index

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

// TestExecuteIsIdempotent proves that running the same corpus through
// IngestService twice succeeds both times and leaves the index in the same
// state — previously the second run failed outright (InsertDocument was a
// plain INSERT against a UNIQUE(path) column), and even after that's fixed,
// naively re-inserting chunks without clearing old ones first would either
// violate UNIQUE(document_id, chunk_idx) or leave stale chunks behind when
// content shrinks between runs.
func TestExecuteIsIdempotent(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "index.sqlite")
	svc := NewIngestService(dbPath)

	status := "in-progress"
	uuid := "abc12345"
	corpus := &EnrichedCorpus{
		Documents: []EnrichedDocument{
			{
				Path:       "/repo/.hawp/work/active/abc12345-thing.md",
				Type:       "task",
				Category:   "work",
				FolderRole: "active",
				Content:    "## Section One\n\nSome content here.\n\n## Section Two\n\nMore content.",
				Status:     &status,
				WorkUUID:   &uuid,
			},
		},
	}

	firstResult, err := svc.Execute(corpus)
	if err != nil {
		t.Fatalf("first Execute() error = %v", err)
	}
	if firstResult.DocumentsProcessed != 1 {
		t.Fatalf("first run: DocumentsProcessed = %d, want 1", firstResult.DocumentsProcessed)
	}

	secondResult, err := svc.Execute(corpus)
	if err != nil {
		t.Fatalf("second Execute() (re-ingest) error = %v, want no error", err)
	}
	if secondResult.DocumentsProcessed != 1 {
		t.Errorf("second run: DocumentsProcessed = %d, want 1", secondResult.DocumentsProcessed)
	}
	if secondResult.ChunksCreated != firstResult.ChunksCreated {
		t.Errorf("second run: ChunksCreated = %d, want %d (same as first run)",
			secondResult.ChunksCreated, firstResult.ChunksCreated)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open db for verification: %v", err)
	}
	defer db.Close()

	var docCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM documents`).Scan(&docCount); err != nil {
		t.Fatalf("count documents: %v", err)
	}
	if docCount != 1 {
		t.Errorf("documents table has %d rows after re-ingest, want 1 (no duplicates)", docCount)
	}

	var chunkCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM chunks`).Scan(&chunkCount); err != nil {
		t.Fatalf("count chunks: %v", err)
	}
	if chunkCount != firstResult.ChunksCreated {
		t.Errorf("chunks table has %d rows after re-ingest, want %d (no stale duplicates)",
			chunkCount, firstResult.ChunksCreated)
	}
}

// TestExecuteReingestWithFewerChunksRemovesStaleOnes proves that when a
// document's content shrinks between runs (fewer chunks), the previous
// run's extra trailing chunks are removed, not left orphaned and stale.
func TestExecuteReingestWithFewerChunksRemovesStaleOnes(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "index.sqlite")
	svc := NewIngestService(dbPath)

	makeCorpus := func(content string) *EnrichedCorpus {
		return &EnrichedCorpus{
			Documents: []EnrichedDocument{
				{
					Path:       "/repo/doc.md",
					Type:       "guide",
					Category:   "kit",
					FolderRole: "start-here",
					Content:    content,
				},
			},
		}
	}

	longResult, err := svc.Execute(makeCorpus("## A\n\ntext\n\n## B\n\ntext\n\n## C\n\ntext"))
	if err != nil {
		t.Fatalf("first Execute() error = %v", err)
	}
	if longResult.ChunksCreated < 2 {
		t.Fatalf("first run created %d chunks, want at least 2 to make this test meaningful", longResult.ChunksCreated)
	}

	shortResult, err := svc.Execute(makeCorpus("## A\n\ntext"))
	if err != nil {
		t.Fatalf("second Execute() (shrinking content) error = %v", err)
	}
	if shortResult.ChunksCreated != 1 {
		t.Fatalf("second run created %d chunks, want 1", shortResult.ChunksCreated)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open db for verification: %v", err)
	}
	defer db.Close()

	var chunkCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM chunks`).Scan(&chunkCount); err != nil {
		t.Fatalf("count chunks: %v", err)
	}
	if chunkCount != 1 {
		t.Errorf("chunks table has %d rows after shrinking re-ingest, want 1 (stale chunks should be removed)", chunkCount)
	}
}
