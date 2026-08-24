// Package sqlite provides the local index database for kit/work documents.
package sqlite

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

// IndexDB wraps the documents index database.
type IndexDB struct {
	db *sql.DB
	tx *sql.Tx
}

// Open opens or creates the index database at dbPath.
func Open(dbPath string) (*IndexDB, error) {
	if err := ensureDir(dbPath); err != nil {
		return nil, err
	}
	sqliteDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Optimize for SQLite concurrency: single connection, WAL mode
	sqliteDB.SetMaxOpenConns(1)
	sqliteDB.SetMaxIdleConns(1)
	sqliteDB.SetConnMaxLifetime(0)

	if err := sqliteDB.Ping(); err != nil {
		return nil, err
	}
	return &IndexDB{db: sqliteDB}, nil
}

// InitSchema initializes the documents, metadata, and chunks tables.
func (ix *IndexDB) InitSchema() error {
	// Enable pragmas for concurrency
	pragmas := `
PRAGMA journal_mode=WAL;
PRAGMA synchronous=NORMAL;
PRAGMA foreign_keys=ON;
PRAGMA busy_timeout=30000;
`
	if _, err := ix.db.Exec(pragmas); err != nil {
		return err
	}

	schema := `
CREATE TABLE IF NOT EXISTS documents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category TEXT NOT NULL CHECK(category IN ('kit', 'work')),
    type TEXT NOT NULL,
    path TEXT NOT NULL UNIQUE,
    folder_role TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS documents_metadata (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    document_id INTEGER NOT NULL UNIQUE,
    work_uuid TEXT NOT NULL,
    status TEXT NOT NULL,
    owner TEXT,
    risk_level TEXT,
    reported_at DATE,
    closed_at DATE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (document_id) REFERENCES documents(id)
);

CREATE TABLE IF NOT EXISTS chunks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    document_id INTEGER NOT NULL,
    chunk_idx INT NOT NULL,
    text TEXT NOT NULL,
    folder_context TEXT,
    metadata_json TEXT,
    embedding BLOB,
    line_start INT NOT NULL DEFAULT 0,
    line_end INT NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (document_id, chunk_idx),
    FOREIGN KEY (document_id) REFERENCES documents(id)
);

CREATE VIRTUAL TABLE IF NOT EXISTS chunks_fts USING fts5(
    text,
    content=chunks,
    content_rowid=id
);

-- chunks_fts is an "external content" FTS5 table: SQLite does NOT
-- automatically keep it in sync with chunks — that's the application's
-- job, via these triggers. Without them, InsertChunk populates chunks but
-- chunks_fts.text MATCH queries silently return nothing for any row
-- inserted after the table was created (confirmed 2026-07-27: a fresh
-- InitSchema + InsertChunk with no triggers produced a chunks_fts table
-- whose row COUNT(*) matched chunks exactly, yet MATCH found zero results —
-- COUNT(*) reflects the content table's rowids regardless of whether the
-- FTS index segments were ever populated for search).
CREATE TRIGGER IF NOT EXISTS chunks_fts_ai AFTER INSERT ON chunks BEGIN
    INSERT INTO chunks_fts(rowid, text) VALUES (new.id, new.text);
END;
CREATE TRIGGER IF NOT EXISTS chunks_fts_ad AFTER DELETE ON chunks BEGIN
    INSERT INTO chunks_fts(chunks_fts, rowid, text) VALUES ('delete', old.id, old.text);
END;
CREATE TRIGGER IF NOT EXISTS chunks_fts_au AFTER UPDATE ON chunks BEGIN
    INSERT INTO chunks_fts(chunks_fts, rowid, text) VALUES ('delete', old.id, old.text);
    INSERT INTO chunks_fts(rowid, text) VALUES (new.id, new.text);
END;

-- Single-row table (id always 1) tracking which backend/model embedded the
-- chunks currently in this index. Query-time embedding (hybrid search) must
-- use this same backend/model — comparing a query vector from one model
-- against document vectors from another is comparing incompatible vector
-- spaces, even when dimensions happen to match. Also guards against
-- accidentally mixing models into one index across separate embed runs.
CREATE TABLE IF NOT EXISTS index_metadata (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    embedding_backend TEXT NOT NULL,
    embedding_model TEXT NOT NULL,
    embedding_dim INTEGER NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_documents_category ON documents(category);
CREATE INDEX IF NOT EXISTS idx_documents_folder_role ON documents(folder_role);
CREATE INDEX IF NOT EXISTS idx_documents_metadata_work_uuid ON documents_metadata(work_uuid);
CREATE INDEX IF NOT EXISTS idx_documents_metadata_status ON documents_metadata(status);
CREATE INDEX IF NOT EXISTS idx_chunks_document_id ON chunks(document_id);
`
	if _, err := ix.db.Exec(schema); err != nil {
		return err
	}
	// Migrate existing DBs: add line columns if absent (ignore "duplicate column" error).
	ix.db.Exec(`ALTER TABLE chunks ADD COLUMN line_start INT NOT NULL DEFAULT 0`)
	ix.db.Exec(`ALTER TABLE chunks ADD COLUMN line_end INT NOT NULL DEFAULT 0`)
	return nil
}

// InsertDocument inserts a document row, or updates it in place if `path`
// already exists, and returns its ID either way. Re-ingesting the same
// corpus (e.g. re-running `hawp search index`) must be safe to repeat, since
// `path` is UNIQUE — a plain INSERT would fail with a constraint violation
// on the second run over any document already indexed.
func (ix *IndexDB) InsertDocument(category, typ, path, folderRole string) (int64, error) {
	var id int64
	err := ix.db.QueryRow(
		`INSERT INTO documents (category, type, path, folder_role)
		 VALUES (?, ?, ?, ?)
		 ON CONFLICT(path) DO UPDATE SET
		   category=excluded.category, type=excluded.type,
		   folder_role=excluded.folder_role, updated_at=CURRENT_TIMESTAMP
		 RETURNING id`,
		category, typ, path, folderRole,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// DocumentMetadata is the data to insert into documents_metadata.
type DocumentMetadata struct {
	DocumentID int64
	WorkUUID   string
	Status     string
	Owner      *string
	RiskLevel  *string
	ReportedAt *string
	ClosedAt   *string
}

// InsertMetadata inserts or updates a document's metadata.
func (ix *IndexDB) InsertMetadata(m DocumentMetadata) error {
	_, err := ix.db.Exec(
		`INSERT INTO documents_metadata (document_id, work_uuid, status, owner, risk_level, reported_at, closed_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(document_id) DO UPDATE SET
		   work_uuid=excluded.work_uuid, status=excluded.status, owner=excluded.owner,
		   risk_level=excluded.risk_level, reported_at=excluded.reported_at,
		   closed_at=excluded.closed_at, updated_at=CURRENT_TIMESTAMP`,
		m.DocumentID, m.WorkUUID, m.Status, m.Owner, m.RiskLevel, m.ReportedAt, m.ClosedAt,
	)
	return err
}

// Chunk is a searchable unit of a document.
type Chunk struct {
	DocumentID    int64
	ChunkIdx      int
	Text          string
	FolderContext *string
	MetadataJSON  *string
	LineStart     int
	LineEnd       int
}

// InsertChunk inserts a chunk and updates the FTS5 index.
func (ix *IndexDB) InsertChunk(c Chunk) error {
	_, err := ix.db.Exec(
		`INSERT INTO chunks (document_id, chunk_idx, text, folder_context, metadata_json, line_start, line_end)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		c.DocumentID, c.ChunkIdx, c.Text, c.FolderContext, c.MetadataJSON, c.LineStart, c.LineEnd,
	)
	if err != nil {
		return err
	}
	// FTS5 automatically indexes on insert via triggers
	return nil
}

// DeleteChunksForDocument removes every chunk belonging to a document (the
// AFTER DELETE trigger keeps chunks_fts in sync). Call this before
// re-inserting a document's chunks on re-ingest: re-running the same corpus
// must be safe to repeat, and a plain re-insert would either violate the
// UNIQUE(document_id, chunk_idx) constraint or, if content shrank, leave
// stale trailing chunks from the previous run still searchable.
func (ix *IndexDB) DeleteChunksForDocument(documentID int64) error {
	_, err := ix.db.Exec(`DELETE FROM chunks WHERE document_id = ?`, documentID)
	return err
}

// Close closes the database connection.
func (ix *IndexDB) Close() error {
	return ix.db.Close()
}

// BeginTx starts a transaction and stores it for use with Tx methods.
func (ix *IndexDB) BeginTx() error {
	tx, err := ix.db.Begin()
	if err != nil {
		return err
	}
	ix.tx = tx
	return nil
}

// Commit commits the current transaction.
func (ix *IndexDB) Commit() error {
	if ix.tx == nil {
		return nil
	}
	err := ix.tx.Commit()
	ix.tx = nil
	return err
}

// Rollback rolls back the current transaction.
func (ix *IndexDB) Rollback() error {
	if ix.tx == nil {
		return nil
	}
	err := ix.tx.Rollback()
	ix.tx = nil
	return err
}

func ensureDir(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "" {
		return os.MkdirAll(dir, 0o755)
	}
	return nil
}

// sanitizeFTSQuery strips chars that FTS5 treats as syntax (dots, dashes,
// parens, quotes, etc.) so raw user input never causes a parse error.
// Each non-word run is replaced by a space; resulting tokens are joined by space
// (implicit AND in FTS5). An empty result means no searchable terms remain.
func sanitizeFTSQuery(q string) string {
	var b strings.Builder
	for _, r := range q {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			b.WriteRune(r)
		} else {
			b.WriteRune(' ')
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

// QueryChunksLexical returns chunks matching a lexical FTS5 query.
func (ix *IndexDB) QueryChunksLexical(query string, limit int) ([]map[string]interface{}, error) {
	query = sanitizeFTSQuery(query)
	if query == "" {
		return nil, nil
	}
	rows, err := ix.db.Query(`
		SELECT c.id, c.document_id, c.chunk_idx, c.text, c.folder_context,
		       d.category, d.type, d.path, d.folder_role,
		       dm.work_uuid, dm.status, dm.closed_at,
		       c.line_start, c.line_end
		FROM chunks_fts
		JOIN chunks c ON chunks_fts.rowid = c.id
		JOIN documents d ON c.document_id = d.id
		LEFT JOIN documents_metadata dm ON d.id = dm.document_id
		WHERE chunks_fts.text MATCH ?
		LIMIT ?
	`, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id, docID, chunkIdx int64
		var text, folderContext *string
		var category, typ, path, folderRole *string
		var workUUID, status, closedAt *string
		var lineStart, lineEnd int64

		if err := rows.Scan(&id, &docID, &chunkIdx, &text, &folderContext,
			&category, &typ, &path, &folderRole, &workUUID, &status, &closedAt,
			&lineStart, &lineEnd); err != nil {
			continue
		}

		result := map[string]interface{}{
			"id":             id,
			"text":           text,
			"folder_context": folderContext,
			"path":           path,
			"type":           typ,
			"category":       category,
			"folder_role":    folderRole,
			"chunk_idx":      chunkIdx,
			"line_start":     lineStart,
			"line_end":       lineEnd,
		}
		if workUUID != nil {
			result["work_uuid"] = workUUID
			result["status"] = status
			result["closed_at"] = closedAt
		}
		results = append(results, result)
	}
	return results, rows.Err()
}

// DocumentCount returns the total number of documents in the index.
func (ix *IndexDB) DocumentCount() (int, error) {
	var count int
	err := ix.db.QueryRow(`SELECT COUNT(*) FROM documents`).Scan(&count)
	return count, err
}

// ChunkCount returns the total number of chunks in the index.
func (ix *IndexDB) ChunkCount() (int, error) {
	var count int
	err := ix.db.QueryRow(`SELECT COUNT(*) FROM chunks`).Scan(&count)
	return count, err
}

// ChunksNeedEmbedding returns the count of chunks with NULL embeddings.
func (ix *IndexDB) ChunksNeedEmbedding() (int, error) {
	var count int
	err := ix.db.QueryRow(`SELECT COUNT(*) FROM chunks WHERE embedding IS NULL`).Scan(&count)
	return count, err
}

// UpdateChunkEmbedding stores a vector embedding for a chunk.
// Uses active transaction if available, otherwise executes on db directly.
func (ix *IndexDB) UpdateChunkEmbedding(chunkID int64, embedding []byte) error {
	var execer interface {
		Exec(string, ...interface{}) (sql.Result, error)
	}

	if ix.tx != nil {
		execer = ix.tx
	} else {
		execer = ix.db
	}

	_, err := execer.Exec(
		`UPDATE chunks SET embedding = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		embedding, chunkID,
	)
	return err
}

// GetChunksNeedingEmbedding returns all chunks with NULL embeddings.
type ChunkData struct {
	ID   int64
	Text string
}

func (ix *IndexDB) GetChunksNeedingEmbedding() ([]ChunkData, error) {
	rows, err := ix.db.Query(`SELECT id, text FROM chunks WHERE embedding IS NULL ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []ChunkData
	for rows.Next() {
		var c ChunkData
		if err := rows.Scan(&c.ID, &c.Text); err != nil {
			continue
		}
		chunks = append(chunks, c)
	}
	return chunks, rows.Err()
}

// GetChunkVector retrieves the embedding vector for a chunk, if available.
func (ix *IndexDB) GetChunkVector(chunkID int64) ([]float32, error) {
	var embeddingJSON []byte
	err := ix.db.QueryRow(`SELECT embedding FROM chunks WHERE id = ? AND embedding IS NOT NULL`, chunkID).Scan(&embeddingJSON)
	if err != nil {
		return nil, err
	}

	var vector []float32
	if err := json.Unmarshal(embeddingJSON, &vector); err != nil {
		return nil, err
	}
	return vector, nil
}

// GetAllChunkVectors retrieves all chunks with their embedding vectors.
func (ix *IndexDB) GetAllChunkVectors() (map[int64][]float32, error) {
	rows, err := ix.db.Query(`SELECT id, embedding FROM chunks WHERE embedding IS NOT NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vectors := make(map[int64][]float32)
	for rows.Next() {
		var chunkID int64
		var embeddingJSON []byte

		if err := rows.Scan(&chunkID, &embeddingJSON); err != nil {
			continue
		}

		var vector []float32
		if err := json.Unmarshal(embeddingJSON, &vector); err != nil {
			continue
		}
		vectors[chunkID] = vector
	}
	return vectors, rows.Err()
}

// HasVectors checks if the database contains any embeddings.
func (ix *IndexDB) HasVectors() (bool, error) {
	var count int
	err := ix.db.QueryRow(`SELECT COUNT(*) FROM chunks WHERE embedding IS NOT NULL`).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetChunkVectors retrieves embeddings for multiple chunk IDs.
func (ix *IndexDB) GetChunkVectors(chunkIDs []int64) (map[int64][]float32, error) {
	if len(chunkIDs) == 0 {
		return make(map[int64][]float32), nil
	}

	// Build placeholder list for SQL IN clause
	placeholders := make([]string, len(chunkIDs))
	args := make([]interface{}, len(chunkIDs))
	for i, id := range chunkIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := `SELECT id, embedding FROM chunks WHERE id IN (` + strings.Join(placeholders, ",") + `) AND embedding IS NOT NULL`
	rows, err := ix.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vectors := make(map[int64][]float32)
	for rows.Next() {
		var chunkID int64
		var embeddingJSON []byte

		if err := rows.Scan(&chunkID, &embeddingJSON); err != nil {
			continue
		}

		var vector []float32
		if err := json.Unmarshal(embeddingJSON, &vector); err != nil {
			continue
		}
		vectors[chunkID] = vector
	}
	return vectors, rows.Err()
}

// EmbeddingMetadata records which backend/model embedded the chunks
// currently in this index.
type EmbeddingMetadata struct {
	Backend string
	Model   string
	Dim     int
}

// GetEmbeddingMetadata returns which backend/model embedded this index's
// chunks, and ok=false if nothing has been embedded yet (no row present).
func (ix *IndexDB) GetEmbeddingMetadata() (meta EmbeddingMetadata, ok bool, err error) {
	row := ix.db.QueryRow(`SELECT embedding_backend, embedding_model, embedding_dim FROM index_metadata WHERE id = 1`)
	err = row.Scan(&meta.Backend, &meta.Model, &meta.Dim)
	if err == sql.ErrNoRows {
		return EmbeddingMetadata{}, false, nil
	}
	if err != nil {
		return EmbeddingMetadata{}, false, err
	}
	return meta, true, nil
}

// SetEmbeddingMetadata records the backend/model used for this embed run.
// Upserts the single row (id always 1).
func (ix *IndexDB) SetEmbeddingMetadata(meta EmbeddingMetadata) error {
	_, err := ix.db.Exec(`
		INSERT INTO index_metadata (id, embedding_backend, embedding_model, embedding_dim, updated_at)
		VALUES (1, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(id) DO UPDATE SET
			embedding_backend = excluded.embedding_backend,
			embedding_model = excluded.embedding_model,
			embedding_dim = excluded.embedding_dim,
			updated_at = CURRENT_TIMESTAMP
	`, meta.Backend, meta.Model, meta.Dim)
	return err
}
