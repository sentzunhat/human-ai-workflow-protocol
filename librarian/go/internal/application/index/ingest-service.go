// Package index provides application services for document indexing.
package index

import (
	"encoding/json"
	"fmt"
	"time"

	domainindex "github.com/sentzunhat/hawp/librarian/go/internal/domain/index"
	domainstore "github.com/sentzunhat/hawp/librarian/go/internal/domain/index/store"
	sqliteindex "github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/sqlite/index"
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
	dbPath    string
	openStore func(string) (domainstore.CorpusWriter, error)
}

// NewIngestService creates an ingest service.
func NewIngestService(dbPath string) *IngestService {
	return NewIngestServiceWithStoreOpener(dbPath, func(path string) (domainstore.CorpusWriter, error) {
		return sqliteindex.Open(path)
	})
}

// NewIngestServiceWithStoreOpener composes ingestion with its persistence port.
func NewIngestServiceWithStoreOpener(dbPath string, opener func(string) (domainstore.CorpusWriter, error)) *IngestService {
	return &IngestService{dbPath: dbPath, openStore: opener}
}

// Execute ingests documents from a corpus into the index database.
func (s *IngestService) Execute(corpus *EnrichedCorpus) (IngestResult, error) {
	result := IngestResult{}
	start := time.Now()

	// Open or create the index database
	db, err := s.openStore(s.dbPath)
	if err != nil {
		return result, fmt.Errorf("open index db: %w", err)
	}
	defer db.Close()

	// Initialize schema
	if err := db.Initialize(); err != nil {
		return result, fmt.Errorf("init schema: %w", err)
	}

	// Ingest each document
	for _, enriched := range corpus.Documents {
		// Build optional work metadata before the atomic document replacement.
		var metadata *domainstore.DocumentMetadata
		if enriched.Category == "work" && enriched.WorkUUID != nil {
			metadata = &domainstore.DocumentMetadata{
				WorkUUID:   *enriched.WorkUUID,
				Status:     *enriched.Status,
				Owner:      enriched.Owner,
				RiskLevel:  enriched.RiskLevel,
				ReportedAt: enriched.ReportedAt,
				ClosedAt:   enriched.ClosedAt,
			}
		}

		// Create chunks
		doc := domainindex.Document{
			Category:   enriched.Category,
			Type:       enriched.Type,
			Path:       enriched.Path,
			FolderRole: enriched.FolderRole,
			Content:    enriched.Content,
		}

		var domainMetadata *domainindex.DocumentMetadata
		if metadata != nil {
			domainMetadata = &domainindex.DocumentMetadata{
				WorkUUID:   *enriched.WorkUUID,
				Status:     *enriched.Status,
				Owner:      enriched.Owner,
				RiskLevel:  enriched.RiskLevel,
				ReportedAt: enriched.ReportedAt,
				ClosedAt:   enriched.ClosedAt,
			}
		}

		folderContext := domainindex.BuildFolderContext(doc, domainMetadata)
		metadataJSONBytes, _ := json.Marshal(enriched.Metadata)
		metadataJSONStr := string(metadataJSONBytes)

		chunkTexts := domainindex.ChunkBySection(enriched.Content)
		replacement := domainstore.DocumentReplacement{
			Category: enriched.Category, Type: enriched.Type, Path: enriched.Path,
			FolderRole: enriched.FolderRole, Metadata: metadata,
			Chunks: make([]domainstore.Chunk, 0, len(chunkTexts)),
		}
		for i, chunkText := range chunkTexts {
			replacement.Chunks = append(replacement.Chunks, domainstore.Chunk{
				Index:         i,
				Text:          chunkText,
				FolderContext: folderContext,
				MetadataJSON:  metadataJSONStr,
			})
		}
		if _, err := db.ReplaceDocument(replacement); err != nil {
			return result, fmt.Errorf("replace document %s: %w", enriched.Path, err)
		}
		result.DocumentsProcessed++
		result.ChunksCreated += len(chunkTexts)
		if metadata != nil {
			result.MetadataRows++
		}
		for _, chunkText := range chunkTexts {
			result.BytesIndexed += int64(len(chunkText))
		}
	}

	result.ElapsedSeconds = time.Since(start).Seconds()
	return result, nil
}
