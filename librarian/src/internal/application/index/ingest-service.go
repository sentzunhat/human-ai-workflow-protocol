// Package index provides application services for document indexing.
package index

import (
	"encoding/json"
	"fmt"
	"time"

	domainindex "github.com/sentzunhat/hawp/librarian/src/internal/domain/index"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/sqlite"
)

// IngestResult summarizes an ingest run.
type IngestResult struct {
	DocumentsProcessed int
	ChunksCreated      int
	MetadataRows       int
	ElapsedSeconds     float64
	BytesIndexed       int64
}

func (r IngestResult) String() string {
	return fmt.Sprintf(`
Ingest complete:
  Documents:  %d
  Chunks:     %d
  Metadata:   %d (work items only)
  Time:       %.2f seconds
  Content:    %d bytes
`,
		r.DocumentsProcessed, r.ChunksCreated, r.MetadataRows,
		r.ElapsedSeconds, r.BytesIndexed,
	)
}

// EnrichedDocument represents a document from `hawp index build` output.
type EnrichedDocument struct {
	Path       string                 `json:"path"`
	Type       string                 `json:"type"`
	Category   string                 `json:"category"` // "kit" or "work"
	FolderRole string                 `json:"folder_role"`
	Content    string                 `json:"content"`
	Status     *string                `json:"status,omitempty"`
	WorkUUID   *string                `json:"work_uuid,omitempty"`
	Owner      *string                `json:"owner,omitempty"`
	RiskLevel  *string                `json:"risk_level,omitempty"`
	ReportedAt *string                `json:"reported_at,omitempty"`
	ClosedAt   *string                `json:"closed_at,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// EnrichedCorpus is the input from `hawp index build`.
type EnrichedCorpus struct {
	Documents []EnrichedDocument `json:"documents"`
}

// IngestService orchestrates the ingest of enriched documents into the index.
type IngestService struct {
	dbPath string
}

// NewIngestService creates an ingest service.
func NewIngestService(dbPath string) *IngestService {
	return &IngestService{dbPath: dbPath}
}

// Execute ingests documents from a corpus into the index database.
func (s *IngestService) Execute(corpus *EnrichedCorpus) (IngestResult, error) {
	result := IngestResult{}
	start := time.Now()

	// Open or create the index database
	db, err := sqlite.Open(s.dbPath)
	if err != nil {
		return result, fmt.Errorf("open index db: %w", err)
	}
	defer db.Close()

	// Initialize schema
	if err := db.InitSchema(); err != nil {
		return result, fmt.Errorf("init schema: %w", err)
	}

	// Ingest each document
	for _, enriched := range corpus.Documents {
		// Insert document
		docID, err := db.InsertDocument(
			enriched.Category, enriched.Type, enriched.Path, enriched.FolderRole,
		)
		if err != nil {
			return result, fmt.Errorf("insert document %s: %w", enriched.Path, err)
		}
		result.DocumentsProcessed++

		// If it's a work item, insert metadata
		if enriched.Category == "work" && enriched.WorkUUID != nil {
			metadata := sqlite.DocumentMetadata{
				DocumentID: docID,
				WorkUUID:   *enriched.WorkUUID,
				Status:     *enriched.Status,
				Owner:      enriched.Owner,
				RiskLevel:  enriched.RiskLevel,
				ReportedAt: enriched.ReportedAt,
				ClosedAt:   enriched.ClosedAt,
			}
			if err := db.InsertMetadata(metadata); err != nil {
				return result, fmt.Errorf("insert metadata for %s: %w", enriched.Path, err)
			}
			result.MetadataRows++
		}

		// Create chunks
		doc := domainindex.Document{
			ID:         docID,
			Category:   enriched.Category,
			Type:       enriched.Type,
			Path:       enriched.Path,
			FolderRole: enriched.FolderRole,
			Content:    enriched.Content,
		}

		var metadata *domainindex.DocumentMetadata
		if enriched.Category == "work" && enriched.WorkUUID != nil {
			metadata = &domainindex.DocumentMetadata{
				DocumentID: docID,
				WorkUUID:   *enriched.WorkUUID,
				Status:     *enriched.Status,
				Owner:      enriched.Owner,
				RiskLevel:  enriched.RiskLevel,
				ReportedAt: enriched.ReportedAt,
				ClosedAt:   enriched.ClosedAt,
			}
		}

		folderContext := domainindex.BuildFolderContext(doc, metadata)
		metadataJSONBytes, _ := json.Marshal(enriched.Metadata)
		metadataJSONStr := string(metadataJSONBytes)

		// Re-ingesting the same document must be safe to repeat: clear its
		// prior chunks first so re-runs neither violate the
		// UNIQUE(document_id, chunk_idx) constraint nor leave stale trailing
		// chunks behind when the new content chunks to fewer pieces.
		if err := db.DeleteChunksForDocument(docID); err != nil {
			return result, fmt.Errorf("clear existing chunks for %s: %w", enriched.Path, err)
		}

		chunks := domainindex.ChunkBySection(enriched.Content)
		for i, chunkText := range chunks {
			chunk := sqlite.Chunk{
				DocumentID:    docID,
				ChunkIdx:      i,
				Text:          chunkText,
				FolderContext: &folderContext,
				MetadataJSON:  &metadataJSONStr,
			}
			if err := db.InsertChunk(chunk); err != nil {
				return result, fmt.Errorf("insert chunk %s#%d: %w", enriched.Path, i, err)
			}
			result.ChunksCreated++
			result.BytesIndexed += int64(len(chunkText))
		}
	}

	result.ElapsedSeconds = time.Since(start).Seconds()
	return result, nil
}
