package index

import (
	"context"
	"errors"
	"testing"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/embeddings"
	domainstore "github.com/sentzunhat/hawp/librarian/go/internal/domain/index/store"
)

type fakeCorpusWriter struct {
	initialized  bool
	replacements []domainstore.DocumentReplacement
}

func (s *fakeCorpusWriter) Initialize() error { s.initialized = true; return nil }

func (s *fakeCorpusWriter) ReplaceDocument(replacement domainstore.DocumentReplacement) (int64, error) {
	s.replacements = append(s.replacements, replacement)
	return int64(len(s.replacements)), nil
}

func (s *fakeCorpusWriter) Close() error { return nil }

func TestIngestServiceUsesCorpusWriterPort(t *testing.T) {
	store := &fakeCorpusWriter{}
	service := NewIngestServiceWithStoreOpener("ignored", func(string) (domainstore.CorpusWriter, error) {
		return store, nil
	})

	status, uuid := "in-progress", "11111111-1111-1111-1111-111111111111"
	result, err := service.Execute(&EnrichedCorpus{Documents: []EnrichedDocument{{
		Path: ".hawp/work/active/example.md", Type: "task", Category: "work", FolderRole: "active",
		Content: "## Plan\n\nWrite the capability-local port.", Status: &status, WorkUUID: &uuid,
	}}})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !store.initialized {
		t.Fatal("corpus writer was not initialized")
	}
	if len(store.replacements) != 1 {
		t.Fatalf("replacements = %d, want 1", len(store.replacements))
	}
	replacement := store.replacements[0]
	if replacement.Metadata == nil || replacement.Metadata.WorkUUID != uuid {
		t.Fatalf("replacement metadata = %#v, want work UUID %q", replacement.Metadata, uuid)
	}
	if len(replacement.Chunks) != result.ChunksCreated || replacement.Chunks[0].FolderContext == "" {
		t.Fatalf("replacement chunks = %#v, want prepared chunks with folder context", replacement.Chunks)
	}
}

type fakeEmbeddingStore struct {
	metadataErr error
}

func (s *fakeEmbeddingStore) Initialize() error { return nil }

func (s *fakeEmbeddingStore) PendingChunks() ([]domainstore.PendingChunk, error) {
	return []domainstore.PendingChunk{{ID: 1, Text: "pending"}}, nil
}

func (s *fakeEmbeddingStore) EmbeddingMetadata() (domainstore.EmbeddingMetadata, bool, error) {
	return domainstore.EmbeddingMetadata{}, false, s.metadataErr
}

func (s *fakeEmbeddingStore) Begin() error { return nil }

func (s *fakeEmbeddingStore) UpdateEmbedding(int64, []byte) error { return nil }

func (s *fakeEmbeddingStore) Commit() error { return nil }

func (s *fakeEmbeddingStore) Rollback() error { return nil }

func (s *fakeEmbeddingStore) SetEmbeddingMetadata(domainstore.EmbeddingMetadata) error { return nil }

func (s *fakeEmbeddingStore) Close() error { return nil }

func TestEmbedServicePropagatesMetadataReadErrorBeforeOpeningModel(t *testing.T) {
	store := &fakeEmbeddingStore{metadataErr: errors.New("metadata unavailable")}
	newEmbedderCalled := false
	service := NewEmbedServiceWithDependencies("ignored",
		func(string) (domainstore.EmbeddingStore, error) { return store, nil },
		func(string, string) (embeddings.Embedder, error) {
			newEmbedderCalled = true
			return nil, nil
		},
	)

	_, err := service.Execute(context.Background(), "none", "none")
	if err == nil || !errors.Is(err, store.metadataErr) {
		t.Fatalf("Execute() error = %v, want metadata read error", err)
	}
	if newEmbedderCalled {
		t.Fatal("embedder was created despite the metadata read error")
	}
}
