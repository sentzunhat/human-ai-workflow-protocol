// Package usage provides the SQLite adapter for the usage-log port.
package usage

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	domainusage "github.com/sentzunhat/hawp/librarian/src/internal/domain/usage"
	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS usage_log (
  id INTEGER PRIMARY KEY AUTOINCREMENT, ts TEXT NOT NULL, tool TEXT NOT NULL,
  query_hash TEXT NOT NULL, query_text TEXT, tokens_in INTEGER NOT NULL,
  tokens_out INTEGER NOT NULL, input_body TEXT, output_body TEXT
);
CREATE INDEX IF NOT EXISTS idx_usage_log_ts ON usage_log(ts DESC);`

const migrateQueryText = `ALTER TABLE usage_log ADD COLUMN query_text TEXT;`

type sqliteStore struct{ db *sql.DB }

func Open(path string) (domainusage.Store, error) {
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
	_, _ = db.Exec(migrateQueryText)
	return &sqliteStore{db: db}, nil
}

func (s *sqliteStore) Close() { _ = s.db.Close() }

func extractQueryText(inputJSON []byte) *string {
	var m map[string]any
	if json.Unmarshal(inputJSON, &m) != nil {
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

func (s *sqliteStore) Write(tool string, inputJSON, outputJSON []byte, logBodies bool) error {
	h := sha256.Sum256(inputJSON)
	queryHash := fmt.Sprintf("%x", h[:8])
	var inBody, outBody *string
	if logBodies {
		in, out := string(inputJSON), string(outputJSON)
		inBody, outBody = &in, &out
	}
	_, err := s.db.Exec(`INSERT INTO usage_log (ts, tool, query_hash, query_text, tokens_in, tokens_out, input_body, output_body) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, time.Now().UTC().Format(time.RFC3339), tool, queryHash, extractQueryText(inputJSON), (len(inputJSON)+3)/4, (len(outputJSON)+3)/4, inBody, outBody)
	return err
}

func (s *sqliteStore) Recent(n int) ([]domainusage.Entry, error) {
	rows, err := s.db.Query(`SELECT id, ts, tool, query_hash, query_text, tokens_in, tokens_out, input_body, output_body FROM usage_log ORDER BY id DESC LIMIT ?`, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entries []domainusage.Entry
	for rows.Next() {
		var e domainusage.Entry
		var ts string
		if err := rows.Scan(&e.ID, &ts, &e.Tool, &e.QueryHash, &e.QueryText, &e.TokensIn, &e.TokensOut, &e.InputBody, &e.OutputBody); err != nil {
			return nil, err
		}
		e.TS, _ = time.Parse(time.RFC3339, ts)
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (s *sqliteStore) GetTotals() (domainusage.Totals, error) {
	var t domainusage.Totals
	err := s.db.QueryRow(`SELECT COUNT(*), COALESCE(SUM(tokens_in),0), COALESCE(SUM(tokens_out),0) FROM usage_log`).Scan(&t.Calls, &t.TokensIn, &t.TokensOut)
	return t, err
}

func (s *sqliteStore) GetReport() (domainusage.Report, error) {
	var r domainusage.Report
	if err := s.db.QueryRow(`SELECT COUNT(*), COALESCE(SUM(tokens_in),0), COALESCE(SUM(tokens_out),0) FROM usage_log`).Scan(&r.Calls, &r.TokensIn, &r.TokensOut); err != nil {
		return r, err
	}
	var first, last string
	if err := s.db.QueryRow(`SELECT MIN(ts), MAX(ts) FROM usage_log`).Scan(&first, &last); err == nil && first != "" {
		t1, _ := time.Parse(time.RFC3339, first)
		t2, _ := time.Parse(time.RFC3339, last)
		r.Since, r.Until = &t1, &t2
	}
	rows, err := s.db.Query(`SELECT tool, COUNT(*), COALESCE(SUM(tokens_in),0), COALESCE(SUM(tokens_out),0) FROM usage_log GROUP BY tool ORDER BY COUNT(*) DESC`)
	if err != nil {
		return r, err
	}
	for rows.Next() {
		var stat domainusage.ToolStat
		if err := rows.Scan(&stat.Tool, &stat.Calls, &stat.TokensIn, &stat.TokensOut); err != nil {
			rows.Close()
			return r, err
		}
		r.ByTool = append(r.ByTool, stat)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return r, err
	}
	rows.Close()
	r.TopEntries, err = s.Recent(20)
	return r, err
}

func (s *sqliteStore) Clear() error { _, err := s.db.Exec(`DELETE FROM usage_log`); return err }
