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

// Store is the port (interface) for reading and writing usage log entries.
// The concrete SQLite implementation is in sqliteStore; callers depend only
// on this interface so the storage backend can be swapped or stubbed in tests.
type Store interface {
	Write(tool string, inputJSON, outputJSON []byte, logBodies bool) error
	Recent(n int) ([]Entry, error)
	GetTotals() (Totals, error)
	GetReport() (Report, error)
	Clear() error
	Close()
}

// sqliteStore is the SQLite-backed implementation of Store.
// It directly imports database/sql and modernc.org/sqlite; callers that want
// a pure-domain boundary should depend on the Store interface, not this type.
// Moving Open() to infrastructure/sqlite is tracked in the arch audit.
type sqliteStore struct {
	db *sql.DB
}

// Open opens (or creates) the usage DB at path. Creates parent dirs as needed.
// Returns a Store interface so callers don't depend on the concrete sqliteStore type.
func Open(path string) (Store, error) {
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
	return &sqliteStore{db: db}, nil
}

// Close releases the DB handle.
func (s *sqliteStore) Close() { s.db.Close() }

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
func (s *sqliteStore) Write(tool string, inputJSON, outputJSON []byte, logBodies bool) error {
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
func (s *sqliteStore) Recent(n int) ([]Entry, error) {
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
func (s *sqliteStore) GetTotals() (Totals, error) {
	var t Totals
	err := s.db.QueryRow(
		`SELECT COUNT(*), COALESCE(SUM(tokens_in),0), COALESCE(SUM(tokens_out),0) FROM usage_log`,
	).Scan(&t.Calls, &t.TokensIn, &t.TokensOut)
	return t, err
}

// ToolStat holds per-tool aggregate counts.
type ToolStat struct {
	Tool      string
	Calls     int
	TokensIn  int
	TokensOut int
}

// Report is the full usage summary for the report subcommand.
type Report struct {
	Calls      int
	TokensIn   int
	TokensOut  int
	ByTool     []ToolStat
	TopEntries []Entry // up to 20 most recent, for query listing
	Since      *time.Time
	Until      *time.Time
}

// GetReport returns aggregate and per-tool stats, plus the top recent entries.
func (s *sqliteStore) GetReport() (Report, error) {
	var r Report

	// Overall totals
	if err := s.db.QueryRow(
		`SELECT COUNT(*), COALESCE(SUM(tokens_in),0), COALESCE(SUM(tokens_out),0) FROM usage_log`,
	).Scan(&r.Calls, &r.TokensIn, &r.TokensOut); err != nil {
		return r, err
	}

	// Time range
	var first, last string
	if err := s.db.QueryRow(`SELECT MIN(ts), MAX(ts) FROM usage_log`).Scan(&first, &last); err == nil && first != "" {
		t1, _ := time.Parse(time.RFC3339, first)
		t2, _ := time.Parse(time.RFC3339, last)
		r.Since = &t1
		r.Until = &t2
	}

	// Per-tool breakdown
	rows, err := s.db.Query(
		`SELECT tool, COUNT(*), COALESCE(SUM(tokens_in),0), COALESCE(SUM(tokens_out),0)
         FROM usage_log GROUP BY tool ORDER BY COUNT(*) DESC`,
	)
	if err != nil {
		return r, err
	}
	defer rows.Close()
	for rows.Next() {
		var ts ToolStat
		if err := rows.Scan(&ts.Tool, &ts.Calls, &ts.TokensIn, &ts.TokensOut); err != nil {
			return r, err
		}
		r.ByTool = append(r.ByTool, ts)
	}
	if err := rows.Err(); err != nil {
		return r, err
	}

	// Top 20 recent entries for query listing
	r.TopEntries, err = s.Recent(20)
	return r, err
}

// FormatReport renders a Report as human-readable Markdown. The result is
// suitable for printing to stdout or committing as an evidence artifact.
func FormatReport(rep Report) string {
	if rep.Calls == 0 {
		return "No calls recorded. Run `hawp usage enable` then make some MCP tool calls.\n"
	}

	var sb strings.Builder
	fmt.Fprintln(&sb, "# HAWP Usage Report")
	fmt.Fprintln(&sb)
	if rep.Since != nil && rep.Until != nil {
		fmt.Fprintf(&sb, "Period: %s → %s\n", rep.Since.Format("2006-01-02 15:04 UTC"), rep.Until.Format("2006-01-02 15:04 UTC"))
	}
	fmt.Fprintln(&sb)
	fmt.Fprintln(&sb, "## Totals")
	fmt.Fprintln(&sb)
	fmt.Fprintf(&sb, "| Metric | Value |\n|--------|-------|\n")
	fmt.Fprintf(&sb, "| Calls | %d |\n", rep.Calls)
	fmt.Fprintf(&sb, "| Tokens in (est.) | ~%d |\n", rep.TokensIn)
	fmt.Fprintf(&sb, "| Tokens out (est.) | ~%d |\n", rep.TokensOut)
	saved := rep.TokensIn - rep.TokensOut
	if saved > 0 {
		pct := float64(saved) / float64(rep.TokensIn) * 100
		fmt.Fprintf(&sb, "| Tokens saved (est.) | ~%d (~%.0f%%) |\n", saved, pct)
	}
	fmt.Fprintln(&sb)

	if len(rep.ByTool) > 0 {
		fmt.Fprintln(&sb, "## By Tool")
		fmt.Fprintln(&sb)
		fmt.Fprintln(&sb, "| Tool | Calls | Tokens In | Tokens Out | Saved |")
		fmt.Fprintln(&sb, "|------|-------|-----------|------------|-------|")
		for _, ts := range rep.ByTool {
			s := ts.TokensIn - ts.TokensOut
			pct := 0.0
			if ts.TokensIn > 0 {
				pct = float64(s) / float64(ts.TokensIn) * 100
			}
			fmt.Fprintf(&sb, "| %s | %d | ~%d | ~%d | ~%d (~%.0f%%) |\n",
				ts.Tool, ts.Calls, ts.TokensIn, ts.TokensOut, s, pct)
		}
		fmt.Fprintln(&sb)
	}

	if len(rep.TopEntries) > 0 {
		fmt.Fprintln(&sb, "## Recent Queries")
		fmt.Fprintln(&sb)
		fmt.Fprintln(&sb, "| # | Time | Tool | Tokens In | Tokens Out | Query |")
		fmt.Fprintln(&sb, "|---|------|------|-----------|------------|-------|")
		for i, e := range rep.TopEntries {
			fmt.Fprintf(&sb, "| %d | %s | %s | %d | %d | %s |\n",
				i+1, e.TS.Format("01-02 15:04"), e.Tool, e.TokensIn, e.TokensOut, EntrySummary(e))
		}
		fmt.Fprintln(&sb)
	}

	fmt.Fprintln(&sb, "_Token estimates: chars/4 (MCP request/response JSON byte length)_")
	return sb.String()
}

// Clear truncates the log. Irreversible.
func (s *sqliteStore) Clear() error {
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
