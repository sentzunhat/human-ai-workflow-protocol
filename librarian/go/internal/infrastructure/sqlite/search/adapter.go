package search

import (
	"context"
	"fmt"

	searchindex "github.com/sentzunhat/hawp/librarian/go/internal/domain/search/index"
	sqlite "github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/sqlite"
)

// Adapter implements the search index port over SQLite. It owns row-shape
// translation so application code never depends on driver values or maps.
type Adapter struct {
	db *sqlite.IndexDB
}

func NewAdapter(db *sqlite.IndexDB) *Adapter {
	return &Adapter{db: db}
}

func (a *Adapter) LexicalSearch(_ context.Context, query string, limit int) ([]searchindex.Candidate, error) {
	rows, err := a.db.QueryChunksLexical(query, limit)
	if err != nil {
		return nil, err
	}

	results := make([]searchindex.Candidate, 0, len(rows))
	for _, row := range rows {
		results = append(results, searchindex.Candidate{
			ID:            intValue(row["id"]),
			Content:       stringValue(row["text"]),
			Source:        stringValue(row["path"]),
			Title:         stringValue(row["folder_role"]),
			FolderContext: stringValue(row["folder_context"]),
			ChunkIndex:    int(intValue(row["chunk_idx"])),
			Type:          stringValue(row["type"]),
			Category:      stringValue(row["category"]),
			WorkUUID:      stringValue(row["work_uuid"]),
			Status:        stringValue(row["status"]),
		})
	}
	return results, nil
}

func (a *Adapter) HasVectors(_ context.Context) (bool, error) {
	return a.db.HasVectors()
}

func (a *Adapter) EmbeddingMetadata(_ context.Context) (searchindex.EmbeddingMetadata, bool, error) {
	meta, ok, err := a.db.GetEmbeddingMetadata()
	if err != nil {
		return searchindex.EmbeddingMetadata{}, false, err
	}
	return searchindex.EmbeddingMetadata{Backend: meta.Backend, Model: meta.Model, Dim: meta.Dim}, ok, nil
}

func (a *Adapter) ChunkVectors(_ context.Context, ids []int64) (map[int64][]float32, error) {
	return a.db.GetChunkVectors(ids)
}

func stringValue(value interface{}) string {
	if value == nil {
		return ""
	}
	if pointer, ok := value.(*string); ok {
		if pointer == nil {
			return ""
		}
		return *pointer
	}
	return fmt.Sprint(value)
}

func intValue(value interface{}) int64 {
	switch typed := value.(type) {
	case int64:
		return typed
	case int:
		return int64(typed)
	default:
		return 0
	}
}
