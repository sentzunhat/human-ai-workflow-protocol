// Package store defines persistence contracts for the index capability.
package store

// DocumentMetadata is optional work-item metadata stored with a document.
type DocumentMetadata struct {
	WorkUUID   string
	Status     string
	Owner      *string
	RiskLevel  *string
	ReportedAt *string
	ClosedAt   *string
}

// Chunk is a document chunk ready to persist.
type Chunk struct {
	Index         int
	Text          string
	FolderContext string
	MetadataJSON  string
}

// DocumentReplacement is the complete persisted state for one document.
// A CorpusWriter must replace it atomically.
type DocumentReplacement struct {
	Category   string
	Type       string
	Path       string
	FolderRole string
	Metadata   *DocumentMetadata
	Chunks     []Chunk
}

// CorpusWriter persists complete indexed documents.
type CorpusWriter interface {
	Initialize() error
	ReplaceDocument(DocumentReplacement) (int64, error)
	Close() error
}

// PendingChunk is a stored chunk that still needs an embedding.
type PendingChunk struct {
	ID   int64
	Text string
}

// EmbeddingMetadata identifies the vector space used by this index.
type EmbeddingMetadata struct {
	Backend string
	Model   string
	Dim     int
}

// EmbeddingStore owns pending embeddings and their persistence lifecycle.
type EmbeddingStore interface {
	Initialize() error
	PendingChunks() ([]PendingChunk, error)
	EmbeddingMetadata() (EmbeddingMetadata, bool, error)
	Begin() error
	UpdateEmbedding(chunkID int64, embedding []byte) error
	Commit() error
	Rollback() error
	SetEmbeddingMetadata(EmbeddingMetadata) error
	Close() error
}
