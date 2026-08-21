// Package index provides application services for document indexing.
package index

import (
	"encoding/json"
	"fmt"
	"time"

	domainindex "github.com/sentzunhat/hawp/librarian/go/internal/domain/index"
	"github.com/sentzunhat/hawp/librarian/go/internal/domain/index/store"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/sqlite"
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
// Use NewIngestService to construct from a database path (backward compatible),
// or NewIngestServiceFromStore to inject an explicit DocumentStore.
type IngestService struct {
	dbPath string
	store  store.DocumentStore
}

// NewIngestService creates an ingest service that opens a SQLite database at
// dbPath on each Execute call (backward compatible constructor).
func NewIngestService(dbPath string) *IngestService {
	return &IngestService{dbPath: dbPath}
}

// NewIngestServiceFromStore creates an ingest service backed by the given
// DocumentStore. InitSchema is the caller's responsibility.
func NewIngestServiceFromStore(s store.DocumentStore) *IngestService {
	return &IngestService{store: s}
}

// Execute ingests documents from a corpus into the index.
func (s *IngestService) Execute(corpus *EnrichedCorpus) (IngestResult, error) {
	if s.store != nil {
		return ingest(corpus, s.store)
	}

	// Backward-compatible path: open SQLite, create adapter, ingest, close.
	db, err := sqlite.Open(s.dbPath)
	if err != nil {
		return IngestResult{}, fmt.Errorf("open index db: %w", err)
	}
	defer db.Close()
	if err := db.InitSchema(); err != nil {
		return IngestResult{}, fmt.Errorf("init schema: %w", err)
	}
	return ingest(corpus, sqliteindex.NewAdapter(db))
}

// ingest runs the core ingestion logic against any DocumentStore implementation.
func ingest(corpus *EnrichedCorpus, ds store.DocumentStore) (IngestResult, error) {
	result := IngestResult{}
	start := time.Now()

	for _, enriched := range corpus.Documents {
		doc := domainindex.Document{
			Category:   enriched.Category,
			Type:       enriched.Type,
			Path:       enriched.Path,
			FolderRole: enriched.FolderRole,
			Content:    enriched.Content,
		}

		var meta *domainindex.DocumentMetadata
		if enriched.Category == "work" && enriched.WorkUUID != nil {
			meta = &domainindex.DocumentMetadata{
				WorkUUID:   *enriched.WorkUUID,
				Status:     *enriched.Status,
				Owner:      enriched.Owner,
				RiskLevel:  enriched.RiskLevel,
				ReportedAt: enriched.ReportedAt,
				ClosedAt:   enriched.ClosedAt,
			}
		}

		folderContext := domainindex.BuildFolderContext(doc, meta)
		metadataJSONBytes, _ := json.Marshal(enriched.Metadata)
		metadataJSONStr := string(metadataJSONBytes)

		rawChunks := domainindex.ChunkBySection(enriched.Content)
		domainChunks := make([]domainindex.Chunk, len(rawChunks))
		for i, text := range rawChunks {
			fc := folderContext
			mj := metadataJSONStr
			domainChunks[i] = domainindex.Chunk{
				ChunkIdx:      i,
				Text:          text,
				FolderContext: fc,
				MetadataJSON:  &mj,
			}
			result.BytesIndexed += int64(len(text))
		}

		if _, err := ds.ReplaceDocument(doc, meta, domainChunks); err != nil {
			return result, fmt.Errorf("replace document %s: %w", enriched.Path, err)
		}

		result.DocumentsProcessed++
		result.ChunksCreated += len(rawChunks)
		if meta != nil {
			result.MetadataRows++
		}
	}

	result.ElapsedSeconds = time.Since(start).Seconds()
	return result, nil
}
