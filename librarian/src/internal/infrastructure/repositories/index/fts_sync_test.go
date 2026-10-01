package sqlite

// TestFTS5SyncsOnInsert locks in a real, severe bug found 2026-07-27:
// chunks_fts is an "external content" FTS5 table (content=chunks), which
// SQLite does NOT automatically keep in sync — that requires triggers,
// which this schema previously had none of. A fresh InitSchema +
// InsertChunk produced a chunks_fts table whose row COUNT(*) matched
// chunks exactly (misleadingly suggesting it worked), yet
// `chunks_fts.text MATCH '...'` returned zero rows for any chunk inserted
// after table creation — meaning a brand-new `hawp search index` run on a
// repo indexed for the first time would have zero working lexical search
// results, the most foundational feature this whole index exists for.
//
// The repo's long-committed .hawp/db/index.sqlite happened to have working
// FTS from however it was originally populated, which masked this for the
// current InitSchema/InsertChunk code path — nothing in this session's
// testing had actually built a fully fresh index from scratch until this
// test did.
import (
	"path/filepath"
	"testing"
)

func TestFTS5SyncsOnInsert(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := Open(dbPath)
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
	if err := db.InsertChunk(Chunk{
		DocumentID: docID,
		ChunkIdx:   0,
		Text:       "Kubernetes manages containerized applications across distributed clusters.",
	}); err != nil {
		t.Fatalf("InsertChunk() error = %v", err)
	}

	results, err := db.QueryChunksLexical("kubernetes", 10)
	if err != nil {
		t.Fatalf("QueryChunksLexical() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("QueryChunksLexical(\"kubernetes\") returned %d results, want 1 — "+
			"a freshly inserted chunk must be findable by lexical search immediately", len(results))
	}
}

// TestFTS5TriggersExist confirms the update/delete sync triggers are part
// of the schema too, not just insert — a schema-level check (rather than
// exercising update/delete through IndexDB's public API, which doesn't
// currently expose a way to update chunk text or delete a chunk directly).
func TestFTS5TriggersExist(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer db.Close()
	if err := db.InitSchema(); err != nil {
		t.Fatalf("InitSchema() error = %v", err)
	}

	for _, name := range []string{"chunks_fts_ai", "chunks_fts_ad", "chunks_fts_au"} {
		var got string
		err := db.db.QueryRow(`SELECT name FROM sqlite_master WHERE type = 'trigger' AND name = ?`, name).Scan(&got)
		if err != nil {
			t.Errorf("trigger %q not found in schema: %v", name, err)
		}
	}
}
