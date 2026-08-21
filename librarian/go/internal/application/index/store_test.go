package index

import (
	"context"
	"errors"
	"testing"

	domainindex "github.com/sentzunhat/hawp/librarian/go/internal/domain/index"
	"github.com/sentzunhat/hawp/librarian/go/internal/domain/index/store"
)

// --- Mock DocumentStore ---

type mockDocumentStore struct {
	replaceErr error
	calls      int
	lastDocID  int64
}

func (m *mockDocumentStore) ReplaceDocument(doc domainindex.Document, meta *domainindex.DocumentMetadata, chunks []domainindex.Chunk) (int64, error) {
	if m.replaceErr != nil {
		return 0, m.replaceErr
	}
	m.calls++
	m.lastDocID++
	return m.lastDocID, nil
}

// --- Mock EmbeddingStore ---

type mockEmbeddingStore struct {
	chunks      []store.ChunkData
	metaReadErr error
	meta        store.EmbeddingMetadata
	metaOK      bool
	storedMeta  store.EmbeddingMetadata
	updatedIDs  []int64
}

func (m *mockEmbeddingStore) GetChunksNeedingEmbedding() ([]store.ChunkData, error) {
	return m.chunks, nil
}
func (m *mockEmbeddingStore) GetEmbeddingMetadata() (store.EmbeddingMetadata, bool, error) {
	return m.meta, m.metaOK, m.metaReadErr
}
func (m *mockEmbeddingStore) SetEmbeddingMetadata(meta store.EmbeddingMetadata) error {
	m.storedMeta = meta
	return nil
}
func (m *mockEmbeddingStore) UpdateChunkEmbedding(chunkID int64, _ []byte) error {
	m.updatedIDs = append(m.updatedIDs, chunkID)
	return nil
}
func (m *mockEmbeddingStore) BeginTx() error  { return nil }
func (m *mockEmbeddingStore) Commit() error   { return nil }
func (m *mockEmbeddingStore) Rollback() error { return nil }

// --- IngestService store tests ---

// TestIngestServiceFromStoreCallsReplaceDocument proves IngestService delegates
// to the DocumentStore contract rather than opening SQLite directly.
func TestIngestServiceFromStoreCallsReplaceDocument(t *testing.T) {
	ms := &mockDocumentStore{}
	svc := NewIngestServiceFromStore(ms)

	status := "in-progress"
	uuid := "abc12345"
	corpus := &EnrichedCorpus{
		Documents: []EnrichedDocument{
			{
				Path: "/repo/.hawp/work/active/abc12345-thing.md",
				Type: "task", Category: "work", FolderRole: "active",
				Content:  "## Section\n\nContent here.",
				Status:   &status,
				WorkUUID: &uuid,
			},
			{
				Path: "/repo/.hawp/kit/start-here.md",
				Type: "guide", Category: "kit", FolderRole: "start-here",
				Content: "## Intro\n\nKit content.",
			},
		},
	}

	result, err := svc.Execute(corpus)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result.DocumentsProcessed != 2 {
		t.Errorf("DocumentsProcessed = %d, want 2", result.DocumentsProcessed)
	}
	if ms.calls != 2 {
		t.Errorf("ReplaceDocument called %d times, want 2", ms.calls)
	}
	if result.MetadataRows != 1 {
		t.Errorf("MetadataRows = %d, want 1 (only work documents)", result.MetadataRows)
	}
}

// TestIngestServiceFromStorePropagatesReplaceError proves an error from
// DocumentStore.ReplaceDocument is returned to the caller.
func TestIngestServiceFromStorePropagatesReplaceError(t *testing.T) {
	ms := &mockDocumentStore{replaceErr: errors.New("storage unavailable")}
	svc := NewIngestServiceFromStore(ms)

	corpus := &EnrichedCorpus{
		Documents: []EnrichedDocument{
			{Path: "/doc.md", Type: "guide", Category: "kit", FolderRole: "start-here", Content: "x"},
		},
	}

	_, err := svc.Execute(corpus)
	if err == nil {
		t.Fatal("Execute() should return error when DocumentStore.ReplaceDocument fails")
	}
	if !errors.Is(err, ms.replaceErr) {
		t.Errorf("Execute() error = %v, want to wrap %v", err, ms.replaceErr)
	}
}

// --- EmbedService store tests ---

// TestEmbedServiceFromStorePropagatesMetadataReadError proves that a real
// GetEmbeddingMetadata failure is surfaced, not silently ignored. Previously
// the mixed-model guard was `if metaErr == nil && ok`, which swallowed errors.
func TestEmbedServiceFromStorePropagatesMetadataReadError(t *testing.T) {
	readErr := errors.New("disk read failure")
	ms := &mockEmbeddingStore{
		chunks:      []store.ChunkData{{ID: 1, Text: "hello"}},
		metaReadErr: readErr,
	}
	svc := NewEmbedServiceFromStore(ms)

	_, err := svc.Execute(context.Background(), "none", "none")
	if err == nil {
		t.Fatal("Execute() should return error when GetEmbeddingMetadata fails")
	}
	if !errors.Is(err, readErr) {
		t.Errorf("Execute() error = %v, want to wrap %v", err, readErr)
	}
}

// TestEmbedServiceFromStoreRejectsMixingModels proves the mixed-model guard
// still fires when the store reports a prior embedding run with different
// backend/model — no SQLite required.
func TestEmbedServiceFromStoreRejectsMixingModels(t *testing.T) {
	ms := &mockEmbeddingStore{
		chunks: []store.ChunkData{{ID: 1, Text: "hello"}},
		meta:   store.EmbeddingMetadata{Backend: "onnx", Model: "bge-base-en-v1.5", Dim: 768},
		metaOK: true,
	}
	svc := NewEmbedServiceFromStore(ms)

	_, err := svc.Execute(context.Background(), "ollama", "nomic-embed-text")
	if err == nil {
		t.Fatal("Execute() should reject embedding with a different backend/model")
	}
}

// TestEmbedServiceFromStoreNoChunks proves Execute returns a clear "no chunks"
// error when the store has nothing to embed — no SQLite required.
func TestEmbedServiceFromStoreNoChunks(t *testing.T) {
	ms := &mockEmbeddingStore{chunks: nil}
	svc := NewEmbedServiceFromStore(ms)

	_, err := svc.Execute(context.Background(), "none", "none")
	if err == nil {
		t.Fatal("Execute() should return error when no chunks need embedding")
	}
}

// --- Interface conformance checks ---

// These compile-time checks prove the mock types implement the interfaces,
// and that the constructors accept those interfaces correctly.
var _ store.DocumentStore = (*mockDocumentStore)(nil)
var _ store.EmbeddingStore = (*mockEmbeddingStore)(nil)
