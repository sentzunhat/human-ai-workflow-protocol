// Package index adapts the SQLite index database to domain index store ports.
package index

import (
	domainstore "github.com/sentzunhat/hawp/librarian/go/internal/domain/index/store"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/sqlite"
)

// Adapter provides index persistence through capability-local store ports.
type Adapter struct {
	db *sqlite.IndexDB
}

// Open opens the SQLite-backed index adapter.
func Open(path string) (*Adapter, error) {
	db, err := sqlite.Open(path)
	if err != nil {
		return nil, err
	}
	return &Adapter{db: db}, nil
}

func (a *Adapter) Initialize() error { return a.db.InitSchema() }

func (a *Adapter) Close() error { return a.db.Close() }

func (a *Adapter) ReplaceDocument(replacement domainstore.DocumentReplacement) (int64, error) {
	var metadata *sqlite.DocumentMetadata
	if replacement.Metadata != nil {
		metadata = &sqlite.DocumentMetadata{
			WorkUUID: replacement.Metadata.WorkUUID, Status: replacement.Metadata.Status,
			Owner: replacement.Metadata.Owner, RiskLevel: replacement.Metadata.RiskLevel,
			ReportedAt: replacement.Metadata.ReportedAt, ClosedAt: replacement.Metadata.ClosedAt,
		}
	}

	chunks := make([]sqlite.Chunk, 0, len(replacement.Chunks))
	for _, chunk := range replacement.Chunks {
		folderContext, metadataJSON := chunk.FolderContext, chunk.MetadataJSON
		chunks = append(chunks, sqlite.Chunk{
			ChunkIdx: chunk.Index, Text: chunk.Text,
			FolderContext: &folderContext, MetadataJSON: &metadataJSON,
		})
	}
	return a.db.ReplaceDocument(sqlite.DocumentReplacement{
		Category: replacement.Category, Type: replacement.Type, Path: replacement.Path,
		FolderRole: replacement.FolderRole, Metadata: metadata, Chunks: chunks,
	})
}

func (a *Adapter) PendingChunks() ([]domainstore.PendingChunk, error) {
	chunks, err := a.db.GetChunksNeedingEmbedding()
	if err != nil {
		return nil, err
	}
	pending := make([]domainstore.PendingChunk, len(chunks))
	for i, chunk := range chunks {
		pending[i] = domainstore.PendingChunk{ID: chunk.ID, Text: chunk.Text}
	}
	return pending, nil
}

func (a *Adapter) EmbeddingMetadata() (domainstore.EmbeddingMetadata, bool, error) {
	metadata, ok, err := a.db.GetEmbeddingMetadata()
	return domainstore.EmbeddingMetadata{Backend: metadata.Backend, Model: metadata.Model, Dim: metadata.Dim}, ok, err
}

func (a *Adapter) Begin() error { return a.db.BeginTx() }

func (a *Adapter) UpdateEmbedding(chunkID int64, embedding []byte) error {
	return a.db.UpdateChunkEmbedding(chunkID, embedding)
}

func (a *Adapter) Commit() error { return a.db.Commit() }

func (a *Adapter) Rollback() error { return a.db.Rollback() }

func (a *Adapter) SetEmbeddingMetadata(metadata domainstore.EmbeddingMetadata) error {
	return a.db.SetEmbeddingMetadata(sqlite.EmbeddingMetadata{
		Backend: metadata.Backend, Model: metadata.Model, Dim: metadata.Dim,
	})
}

var _ domainstore.CorpusWriter = (*Adapter)(nil)
var _ domainstore.EmbeddingStore = (*Adapter)(nil)
