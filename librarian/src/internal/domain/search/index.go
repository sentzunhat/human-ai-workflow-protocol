package search

// EmbeddingMetadata identifies the vector space used by an index.
type EmbeddingMetadata struct {
	Backend string
	Model   string
	Dim     int
}

// Index is the repository contract required by lexical and semantic retrieval.
// Implementations own persistence; callers close an opened index after use.
type Index interface {
	Close() error
	QueryChunksLexical(query string, limit int) ([]map[string]interface{}, error)
	HasVectors() (bool, error)
	GetEmbeddingMetadata() (EmbeddingMetadata, bool, error)
	GetAllChunkVectors() (map[int64][]float32, error)
	QueryChunksByIDs(ids []int64) ([]map[string]interface{}, error)
	GetChunkVectors(chunkIDs []int64) (map[int64][]float32, error)
}
