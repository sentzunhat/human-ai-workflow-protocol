// Package store defines capability-local port interfaces for index persistence.
// Application services depend on these contracts; infrastructure adapters
// implement them. This keeps application logic free of SQLite construction.
package store

import (
	"github.com/sentzunhat/hawp/librarian/go/internal/domain/index"
)

// DocumentStore is the write port for atomic enriched-document persistence.
// A single call to ReplaceDocument is the only write operation the ingest
// service needs — it atomically upserts the document with optional work
// metadata, clears previous chunks, and inserts replacement chunks.
type DocumentStore interface {
	// ReplaceDocument atomically upserts a document with optional work metadata,
	// clears its previous chunks, and inserts replacement chunks.
	// Returns the database-assigned document ID.
	ReplaceDocument(doc index.Document, meta *index.DocumentMetadata, chunks []index.Chunk) (int64, error)
}

// ChunkData is a chunk record returned when chunks need embedding.
type ChunkData struct {
	ID   int64
	Text string
}

// EmbeddingMetadata records which backend/model embedded an index's chunks.
// Mixed-model protection: comparing a query vector from one model against
// document vectors from another is comparing incompatible vector spaces even
// when dimensions coincidentally match.
type EmbeddingMetadata struct {
	Backend string
	Model   string
	Dim     int
}

// EmbeddingStore is the read/write port for chunk embedding persistence.
// The embed service depends on this contract rather than on SQLite directly.
type EmbeddingStore interface {
	// GetChunksNeedingEmbedding returns all chunks with a NULL embedding vector.
	GetChunksNeedingEmbedding() ([]ChunkData, error)

	// GetEmbeddingMetadata returns the backend/model that embedded this index's
	// chunks. ok is false (and err is nil) when no embedding run has happened
	// yet. Implementations must return a non-nil error for real read failures
	// so the embed service can propagate them rather than silently skipping the
	// mixed-model guard.
	GetEmbeddingMetadata() (EmbeddingMetadata, bool, error)

	// SetEmbeddingMetadata records the backend/model for this embed run.
	// Upserts the single metadata row.
	SetEmbeddingMetadata(meta EmbeddingMetadata) error

	// UpdateChunkEmbedding stores a serialised vector for one chunk.
	// Uses the active transaction when one is open.
	UpdateChunkEmbedding(chunkID int64, embedding []byte) error

	// BeginTx starts a write transaction.
	BeginTx() error
	// Commit commits the current transaction.
	Commit() error
	// Rollback rolls back the current transaction.
	Rollback() error
}
