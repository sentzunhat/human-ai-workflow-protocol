package index

import "context"

// Candidate is the typed boundary between an index adapter and search rules.
// Storage-specific row maps must not cross this boundary.
type Candidate struct {
	ID            int64
	Content       string
	Source        string
	Title         string
	FolderContext string
	ChunkIndex    int
	Type          string
	Category      string
	WorkUUID      string
	Status        string
	LexicalRank   float32
	SemanticScore float32
	Relevance     float32
}

type EmbeddingMetadata struct {
	Backend string
	Model   string
	Dim     int
}

// IndexPort is the search capability required by the application layer.
// Concrete storage stays behind infrastructure adapters.
type IndexPort interface {
	LexicalSearch(ctx context.Context, query string, limit int) ([]Candidate, error)
	HasVectors(ctx context.Context) (bool, error)
	EmbeddingMetadata(ctx context.Context) (EmbeddingMetadata, bool, error)
	ChunkVectors(ctx context.Context, ids []int64) (map[int64][]float32, error)
}
