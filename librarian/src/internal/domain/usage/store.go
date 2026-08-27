// Package usage records MCP tool calls to a local SQLite log.
// Logging is opt-in; body capture requires a second explicit opt-in.
package usage

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS usage_log (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  ts          TEXT    NOT NULL,
  tool        TEXT    NOT NULL,
  query_hash  TEXT    NOT NULL,
  query_text  TEXT,
  tokens_in   INTEGER NOT NULL,
  tokens_out  INTEGER NOT NULL,
  input_body  TEXT,
  output_body TEXT
);
CREATE INDEX IF NOT EXISTS idx_usage_log_ts ON usage_log(ts DESC);
`

const migrateQueryText = `
ALTER TABLE usage_log ADD COLUMN query_text TEXT;
`

// Config is the user preference file at ~/.hawp/config/usage.json.
type Config struct {
	Enabled   bool `json:"enabled"`
	LogBodies bool `json:"log_bodies"`
}

// Entry is a single recorded tool call.
type Entry struct {
	ID         int64
	TS         time.Time
	Tool       string
	QueryHash  string
	QueryText  *string // first 256 chars of the query field; always stored when present
	TokensIn   int
	TokensOut  int
	InputBody  *string
	OutputBody *string
}

// Totals summarises the whole log.
type Totals struct {
	Calls     int
	TokensIn  int
	TokensOut int
}

// Store wraps the usage SQLite DB.
type Store struct {
	db *sql.DB
}

// Open opens (or creates) the usage DB at path. Creates parent dirs as needed.
func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("usage schema: %w", err)
	}
	// Idempotent migration: add query_text column to existing databases.
	// SQLite returns "duplicate column name" when the column already exists;
	// that error is expected and safe to ignore.
	_, _ = db.Exec(migrateQueryText)
	return &Store{db: db}, nil
}

// Close releases the DB handle.
func (s *Store) Close() { s.db.Close() }

// extractQueryText pulls the first 256 chars of the "query" (or "title")
// field from an input JSON blob. Returns nil when neither field is present.
func extractQueryText(inputJSON []byte) *string {
	var m map[string]any
	if err := json.Unmarshal(inputJSON, &m); err != nil {
		return nil
	}
	for _, key := range []string{"query", "title"} {
		if v, ok := m[key].(string); ok && v != "" {
			if len(v) > 256 {
				v = v[:256]
			}
			return &v
		}
	}
	return nil
}

// Write inserts one log entry. inputJSON and outputJSON are the raw JSON
// bytes of the request arguments and response result. When logBodies is
// false the bodies are stored as NULL.
func (s *Store) Write(tool string, inputJSON, outputJSON []byte, logBodies bool) error {
	h := sha256.Sum256(inputJSON)
	queryHash := fmt.Sprintf("%x", h[:8]) // 16 hex chars — enough for dedup, not full fingerprint
	tokensIn := (len(inputJSON) + 3) / 4
	tokensOut := (len(outputJSON) + 3) / 4

	queryText := extractQueryText(inputJSON)

	var inBody, outBody *string
	if logBodies {
		in := string(inputJSON)
		out := string(outputJSON)
		inBody = &in
		outBody = &out
	}

	_, err := s.db.Exec(
		`INSERT INTO usage_log (ts, tool, query_hash, query_text, tokens_in, tokens_out, input_body, output_body)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		time.Now().UTC().Format(time.RFC3339),
		tool, queryHash, queryText, tokensIn, tokensOut, inBody, outBody,
	)
	return err
}

// Recent returns up to n most recent entries, newest first.
func (s *Store) Recent(n int) ([]Entry, error) {
	rows, err := s.db.Query(
		`SELECT id, ts, tool, query_hash, query_text, tokens_in, tokens_out, input_body, output_body
         FROM usage_log ORDER BY id DESC LIMIT ?`, n,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entries []Entry
	for rows.Next() {
		var e Entry
		var tsStr string
		if err := rows.Scan(&e.ID, &tsStr, &e.Tool, &e.QueryHash, &e.QueryText,
			&e.TokensIn, &e.TokensOut, &e.InputBody, &e.OutputBody); err != nil {
			return nil, err
		}
		e.TS, _ = time.Parse(time.RFC3339, tsStr)
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

// GetTotals returns aggregate counts across the whole log.
func (s *Store) GetTotals() (Totals, error) {
	var t Totals
	err := s.db.QueryRow(
		`SELECT COUNT(*), COALESCE(SUM(tokens_in),0), COALESCE(SUM(tokens_out),0) FROM usage_log`,
	).Scan(&t.Calls, &t.TokensIn, &t.TokensOut)
	return t, err
}

// Clear truncates the log. Irreversible.
func (s *Store) Clear() error {
	_, err := s.db.Exec(`DELETE FROM usage_log`)
	return err
}

// LoadConfig reads ~/.hawp/config/usage.json.
// Returns defaults (disabled) when the file is absent or unreadable.
func LoadConfig(configFile string) Config {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return Config{}
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return Config{}
	}
	return c
}

// SaveConfig writes the config to configFile, creating parent dirs as needed.
func SaveConfig(configFile string, c Config) error {
	if err := os.MkdirAll(filepath.Dir(configFile), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0o644)
}

// QuerySummary extracts the first "query" field value from a JSON args blob.
// Used for display; falls back to the raw hash when no query field is present.
func QuerySummary(inputJSON []byte, queryHash string) string {
	var m map[string]any
	if err := json.Unmarshal(inputJSON, &m); err != nil {
		return queryHash
	}
	if q, ok := m["query"].(string); ok && q != "" {
		if len(q) > 60 {
			q = q[:57] + "..."
		}
		return q
	}
	if t, ok := m["title"].(string); ok && t != "" {
		if len(t) > 60 {
			t = t[:57] + "..."
		}
		return t
	}
	return queryHash
}

// EntrySummary returns the best display string for an entry's query.
// Prefers the always-stored query_text field; falls back to QuerySummary
// against the input body (when bodies are enabled); finally falls back to
// the hash. This means `hawp usage log` shows a human-readable query even
// when body capture is disabled (the default).
func EntrySummary(e Entry) string {
	if e.QueryText != nil && *e.QueryText != "" {
		s := *e.QueryText
		if len(s) > 60 {
			s = s[:57] + "..."
		}
		return s
	}
	if e.InputBody != nil {
		return QuerySummary([]byte(*e.InputBody), e.QueryHash)
	}
	return e.QueryHash
}

// FormatTotals returns a human-readable totals summary.
func FormatTotals(t Totals) string {
	if t.Calls == 0 {
		return "No calls recorded. Run `hawp usage enable` then make some MCP tool calls."
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Calls:       %d\n", t.Calls)
	fmt.Fprintf(&sb, "Tokens in:   ~%d\n", t.TokensIn)
	fmt.Fprintf(&sb, "Tokens out:  ~%d\n", t.TokensOut)
	saved := t.TokensIn - t.TokensOut
	if saved > 0 {
		pct := float64(saved) / float64(t.TokensIn) * 100
		fmt.Fprintf(&sb, "Saved:       ~%d (~%.0f%% reduction via context shaping)\n", saved, pct)
	}
	return sb.String()
}
