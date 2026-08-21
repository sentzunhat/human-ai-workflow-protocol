// Package sqliteindex provides the SQLite adapter implementing the
// domain/index/store capability contracts. It bridges domain types to the
// sqlite.IndexDB low-level operations so application services depend only on
// typed port interfaces.
package sqliteindex

import (
	domainindex "github.com/sentzunhat/hawp/librarian/go/internal/domain/index"
	"github.com/sentzunhat/hawp/librarian/go/internal/domain/index/store"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/sqlite"
)

// Adapter wraps *sqlite.IndexDB and implements store.DocumentStore and
// store.EmbeddingStore. It converts domain types to the sqlite-local types
// expected by IndexDB.
type Adapter struct {
	db *sqlite.IndexDB
}

// NewAdapter creates an Adapter backed by an already-opened IndexDB.
// The caller retains ownership of db (Open/InitSchema/Close).
func NewAdapter(db *sqlite.IndexDB) *Adapter {
	return &Adapter{db: db}
}

// ReplaceDocument implements store.DocumentStore. It maps domain types to
// sqlite-local types and delegates to IndexDB.ReplaceDocument for the atomic
// upsert + chunk replacement.
func (a *Adapter) ReplaceDocument(
	doc domainindex.Document,
	meta *domainindex.DocumentMetadata,
	chunks []domainindex.Chunk,
) (int64, error) {
	var sqliteMeta *sqlite.DocumentMetadata
	if meta != nil {
		sqliteMeta = &sqlite.DocumentMetadata{
			WorkUUID:   meta.WorkUUID,
			Status:     meta.Status,
			Owner:      meta.Owner,
			RiskLevel:  meta.RiskLevel,
			ReportedAt: meta.ReportedAt,
			ClosedAt:   meta.ClosedAt,
		}
	}

	sqliteChunks := make([]sqlite.Chunk, len(chunks))
	for i, c := range chunks {
		var fc *string
		if c.FolderContext != "" {
			v := c.FolderContext
			fc = &v
		}
		sqliteChunks[i] = sqlite.Chunk{
			ChunkIdx:     c.ChunkIdx,
			Text:         c.Text,
			FolderContext: fc,
			MetadataJSON: c.MetadataJSON,
		}
	}

	return a.db.ReplaceDocument(
		sqlite.DocumentRow{
			Category:   doc.Category,
			Type:       doc.Type,
			Path:       doc.Path,
			FolderRole: doc.FolderRole,
		},
		sqliteMeta,
		sqliteChunks,
	)
}

// GetChunksNeedingEmbedding implements store.EmbeddingStore.
func (a *Adapter) GetChunksNeedingEmbedding() ([]store.ChunkData, error) {
	raw, err := a.db.GetChunksNeedingEmbedding()
	if err != nil {
		return nil, err
	}
	out := make([]store.ChunkData, len(raw))
	for i, c := range raw {
		out[i] = store.ChunkData{ID: c.ID, Text: c.Text}
	}
	return out, nil
}

// GetEmbeddingMetadata implements store.EmbeddingStore.
func (a *Adapter) GetEmbeddingMetadata() (store.EmbeddingMetadata, bool, error) {
	raw, ok, err := a.db.GetEmbeddingMetadata()
	if err != nil {
		return store.EmbeddingMetadata{}, false, err
	}
	if !ok {
		return store.EmbeddingMetadata{}, false, nil
	}
	return store.EmbeddingMetadata{
		Backend: raw.Backend,
		Model:   raw.Model,
		Dim:     raw.Dim,
	}, true, nil
}

// SetEmbeddingMetadata implements store.EmbeddingStore.
func (a *Adapter) SetEmbeddingMetadata(meta store.EmbeddingMetadata) error {
	return a.db.SetEmbeddingMetadata(sqlite.EmbeddingMetadata{
		Backend: meta.Backend,
		Model:   meta.Model,
		Dim:     meta.Dim,
	})
}

// UpdateChunkEmbedding implements store.EmbeddingStore.
func (a *Adapter) UpdateChunkEmbedding(chunkID int64, embedding []byte) error {
	return a.db.UpdateChunkEmbedding(chunkID, embedding)
}

// BeginTx implements store.EmbeddingStore.
func (a *Adapter) BeginTx() error { return a.db.BeginTx() }

// Commit implements store.EmbeddingStore.
func (a *Adapter) Commit() error { return a.db.Commit() }

// Rollback implements store.EmbeddingStore.
func (a *Adapter) Rollback() error { return a.db.Rollback() }
