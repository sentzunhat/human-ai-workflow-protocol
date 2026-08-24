package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/sqlite"
)

// setupTestRepo creates a minimal HAWP repo with an indexed document and
// returns its root path. The indexed document has known content and line
// positions so assertions can be precise.
func setupTestRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	// Minimal HAWP structure
	must(t, os.MkdirAll(filepath.Join(root, ".hawp", "db"), 0o755))
	must(t, os.MkdirAll(filepath.Join(root, ".hawp", "kit"), 0o755))
	must(t, os.MkdirAll(filepath.Join(root, ".hawp", "work"), 0o755))

	// Write a minimal BACKLOG.md so toolWorkValidate doesn't error on missing file.
	backlog := "# HAWP Work Backlog\n\n## Active Work\n\n| ID | Title | Status |\n|---|---|---|\n\n## Recently Closed\n\n| ID | Title | Closed |\n|---|---|---|\n"
	must(t, os.WriteFile(filepath.Join(root, ".hawp", "work", "BACKLOG.md"), []byte(backlog), 0o644))

	// Write a kit start-here stub
	must(t, os.WriteFile(filepath.Join(root, ".hawp", "kit", "start-here.md"), []byte("# HAWP start-here\n\nThis is the start-here guide.\n"), 0o644))

	// Seed the index with a document that has known line positions.
	dbPath := filepath.Join(root, ".hawp", "db", "index.sqlite")
	db, err := sqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	defer db.Close()
	if err := db.InitSchema(); err != nil {
		t.Fatalf("InitSchema: %v", err)
	}

	docID, err := db.InsertDocument("kit", "guide", ".hawp/kit/backlog-alignment.md", "kit/references")
	if err != nil {
		t.Fatalf("InsertDocument: %v", err)
	}

	content := "line one\nline two backlog alignment policy\nline three\nline four\nline five"
	err = db.InsertChunk(sqlite.Chunk{
		DocumentID: docID,
		ChunkIdx:   0,
		Text:       content,
		LineStart:  1,
		LineEnd:    5,
	})
	if err != nil {
		t.Fatalf("InsertChunk: %v", err)
	}

	return root
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

// TestToolSearchReturnsStructuredJSON verifies that hawp_search returns the
// structured SearchResponse schema with lines and context.window fields.
func TestToolSearchReturnsStructuredJSON(t *testing.T) {
	root := setupTestRepo(t)

	args, _ := json.Marshal(map[string]any{"query": "backlog alignment", "limit": 3})
	resp := toolSearch(args, root)

	if resp.Result == nil {
		t.Fatal("toolSearch returned nil Result")
	}
	res := resp.Result.(toolResult)
	if res.IsError {
		t.Fatalf("toolSearch returned error: %s", res.Content[0].Text)
	}
	if len(res.Content) == 0 {
		t.Fatal("toolSearch returned empty content")
	}

	var sr SearchResponse
	if err := json.Unmarshal([]byte(res.Content[0].Text), &sr); err != nil {
		t.Fatalf("response is not valid JSON: %v\ntext: %s", err, res.Content[0].Text)
	}
	if len(sr.Results) == 0 {
		t.Fatal("expected at least one result")
	}

	r := sr.Results[0]

	// source path must be set
	if r.Source == "" {
		t.Error("result.source is empty")
	}
	// content must be non-empty
	if r.Content == "" {
		t.Error("result.content is empty")
	}
	// lines.range must span at least one line
	if r.Lines.Range.Start < 1 {
		t.Errorf("lines.range.start = %d, want >= 1", r.Lines.Range.Start)
	}
	if r.Lines.Range.End < r.Lines.Range.Start {
		t.Errorf("lines.range.end %d < start %d", r.Lines.Range.End, r.Lines.Range.Start)
	}
	// lines.source must be within the chunk range
	if r.Lines.Source < r.Lines.Range.Start || r.Lines.Source > r.Lines.Range.End+1 {
		t.Errorf("lines.source %d outside range [%d, %d]", r.Lines.Source, r.Lines.Range.Start, r.Lines.Range.End)
	}
	// context.window must be non-zero
	if r.Context.Window.Start < 1 {
		t.Errorf("context.window.start = %d, want >= 1", r.Context.Window.Start)
	}
	if r.Context.Window.End < r.Lines.Source {
		t.Errorf("context.window.end %d < lines.source %d", r.Context.Window.End, r.Lines.Source)
	}
}

// TestToolSearchNoIndex returns a tool error when DB is missing.
func TestToolSearchNoIndex(t *testing.T) {
	root := t.TempDir()
	args, _ := json.Marshal(map[string]any{"query": "anything"})
	resp := toolSearch(args, root)
	res := resp.Result.(toolResult)
	if !res.IsError {
		t.Error("expected IsError=true when no index exists")
	}
}

// TestToolSearchEmptyQueryErrors returns a tool error for blank query.
func TestToolSearchEmptyQueryErrors(t *testing.T) {
	root := t.TempDir()
	args, _ := json.Marshal(map[string]any{"query": ""})
	resp := toolSearch(args, root)
	res := resp.Result.(toolResult)
	if !res.IsError {
		t.Error("expected IsError=true for empty query")
	}
}

// TestToolSearchNoResults returns empty results array (not error) when no match.
func TestToolSearchNoResults(t *testing.T) {
	root := setupTestRepo(t)
	args, _ := json.Marshal(map[string]any{"query": "xyzzy_no_match_at_all"})
	resp := toolSearch(args, root)
	res := resp.Result.(toolResult)
	if res.IsError {
		t.Fatalf("unexpected error for no-match query: %s", res.Content[0].Text)
	}
	var sr SearchResponse
	if err := json.Unmarshal([]byte(res.Content[0].Text), &sr); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if sr.Results == nil {
		t.Error("results should be an empty array, not null")
	}
}

// TestToolWorkNewCreatesItem verifies hawp_work_new returns a UUID and plan path.
func TestToolWorkNewCreatesItem(t *testing.T) {
	root := setupTestRepo(t)
	args, _ := json.Marshal(map[string]any{"title": "e2e test item", "type": "task"})
	resp := toolWorkNew(args, root)
	res := resp.Result.(toolResult)
	if res.IsError {
		t.Fatalf("toolWorkNew error: %s", res.Content[0].Text)
	}
	text := res.Content[0].Text
	if !strings.Contains(text, ".hawp/work/active/") {
		t.Errorf("expected work item path in response, got: %s", text)
	}
	if !strings.Contains(text, "Created work item") {
		t.Errorf("expected 'Created work item' in response, got: %s", text)
	}
}

// TestToolWorkNewMissingTitle returns a tool error.
func TestToolWorkNewMissingTitle(t *testing.T) {
	root := t.TempDir()
	args, _ := json.Marshal(map[string]any{})
	resp := toolWorkNew(args, root)
	res := resp.Result.(toolResult)
	if !res.IsError {
		t.Error("expected error for missing title")
	}
}
