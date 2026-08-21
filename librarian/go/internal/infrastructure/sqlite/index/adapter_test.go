package sqliteindex_test

import (
	"path/filepath"
	"testing"

	domainindex "github.com/sentzunhat/hawp/librarian/go/internal/domain/index"
	"github.com/sentzunhat/hawp/librarian/go/internal/domain/index/store"
	sqliteindex "github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/sqlite/index"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/sqlite"
)

func openAdapter(t *testing.T) (*sqliteindex.Adapter, func()) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := sqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("sqlite.Open() error = %v", err)
	}
	if err := db.InitSchema(); err != nil {
		t.Fatalf("db.InitSchema() error = %v", err)
	}
	return sqliteindex.NewAdapter(db), func() { db.Close() }
}

// TestAdapterDocumentPersistence proves document persistence via the typed
// DocumentStore contract: a domain Document is persisted and its ID returned.
func TestAdapterDocumentPersistence(t *testing.T) {
	adapter, cleanup := openAdapter(t)
	defer cleanup()

	doc := domainindex.Document{
		Category:   "kit",
		Type:       "guide",
		Path:       "/hawp/kit/start-here.md",
		FolderRole: "start-here",
		Content:    "## Section\n\ncontent",
	}

	id, err := adapter.ReplaceDocument(doc, nil, nil)
	if err != nil {
		t.Fatalf("ReplaceDocument() error = %v", err)
	}
	if id == 0 {
		t.Error("ReplaceDocument() returned id == 0")
	}

	// Re-ingest same path must return same id (upsert, not insert).
	id2, err := adapter.ReplaceDocument(doc, nil, nil)
	if err != nil {
		t.Fatalf("ReplaceDocument() re-ingest error = %v", err)
	}
	if id2 != id {
		t.Errorf("re-ingest returned id=%d, want %d (same document)", id2, id)
	}
}

// TestAdapterChunkPersistence proves chunk persistence is separate from
// document persistence: chunks are cleared and replaced atomically.
func TestAdapterChunkPersistence(t *testing.T) {
	adapter, cleanup := openAdapter(t)
	defer cleanup()

	doc := domainindex.Document{
		Category:   "kit",
		Type:       "guide",
		Path:       "/hawp/kit/doc.md",
		FolderRole: "start-here",
	}
	chunks := []domainindex.Chunk{
		{ChunkIdx: 0, Text: "chunk zero", FolderContext: "ctx"},
		{ChunkIdx: 1, Text: "chunk one", FolderContext: "ctx"},
	}

	if _, err := adapter.ReplaceDocument(doc, nil, chunks); err != nil {
		t.Fatalf("ReplaceDocument() error = %v", err)
	}

	// Verify chunks are searchable immediately via the adapter's embedding path.
	raw, err := adapter.GetChunksNeedingEmbedding()
	if err != nil {
		t.Fatalf("GetChunksNeedingEmbedding() error = %v", err)
	}
	if len(raw) != 2 {
		t.Errorf("GetChunksNeedingEmbedding() returned %d chunks, want 2", len(raw))
	}
}

// TestAdapterEmbeddingMetadataPersistence proves the embedding-metadata
// read/write round-trip via the EmbeddingStore contract.
func TestAdapterEmbeddingMetadataPersistence(t *testing.T) {
	adapter, cleanup := openAdapter(t)
	defer cleanup()

	_, ok, err := adapter.GetEmbeddingMetadata()
	if err != nil {
		t.Fatalf("GetEmbeddingMetadata() before any embed error = %v", err)
	}
	if ok {
		t.Error("ok should be false before any embedding run")
	}

	want := store.EmbeddingMetadata{Backend: "onnx", Model: "bge-base-en-v1.5", Dim: 768}
	if err := adapter.SetEmbeddingMetadata(want); err != nil {
		t.Fatalf("SetEmbeddingMetadata() error = %v", err)
	}

	got, ok, err := adapter.GetEmbeddingMetadata()
	if err != nil {
		t.Fatalf("GetEmbeddingMetadata() after set error = %v", err)
	}
	if !ok {
		t.Fatal("ok should be true after SetEmbeddingMetadata")
	}
	if got != want {
		t.Errorf("GetEmbeddingMetadata() = %+v, want %+v", got, want)
	}
}

// TestAdapterImplementsDocumentStore is a compile-time check that Adapter
// satisfies the DocumentStore interface.
func TestAdapterImplementsDocumentStore(t *testing.T) {
	adapter, cleanup := openAdapter(t)
	defer cleanup()
	var _ store.DocumentStore = adapter
}

// TestAdapterImplementsEmbeddingStore is a compile-time check that Adapter
// satisfies the EmbeddingStore interface.
func TestAdapterImplementsEmbeddingStore(t *testing.T) {
	adapter, cleanup := openAdapter(t)
	defer cleanup()
	var _ store.EmbeddingStore = adapter
}
