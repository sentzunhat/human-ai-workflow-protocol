package usage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func strPtr(s string) *string { return &s }

func openTemp(t *testing.T) *Store {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "usage.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(s.Close)
	return s
}

func TestOpenCreatesSchema(t *testing.T) {
	s := openTemp(t)
	totals, err := s.GetTotals()
	if err != nil {
		t.Fatalf("GetTotals on empty store: %v", err)
	}
	if totals.Calls != 0 {
		t.Errorf("expected 0 calls, got %d", totals.Calls)
	}
}

func TestWriteAndRecent(t *testing.T) {
	s := openTemp(t)

	input := []byte(`{"query":"kubernetes deployments"}`)
	output := []byte(`{"content":"some context block here"}`)

	if err := s.Write("hawp_search", input, output, false); err != nil {
		t.Fatalf("Write: %v", err)
	}

	entries, err := s.Recent(10)
	if err != nil {
		t.Fatalf("Recent: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Tool != "hawp_search" {
		t.Errorf("tool: got %q want %q", e.Tool, "hawp_search")
	}
	if e.TokensIn <= 0 {
		t.Error("tokens_in should be > 0")
	}
	if e.InputBody != nil {
		t.Error("input_body should be nil when logBodies=false")
	}
	// query_text should be stored even when bodies are off
	if e.QueryText == nil || *e.QueryText != "kubernetes deployments" {
		t.Errorf("query_text: got %v, want %q", e.QueryText, "kubernetes deployments")
	}
}

func TestExtractQueryText(t *testing.T) {
	cases := []struct {
		input []byte
		want  *string
	}{
		{[]byte(`{"query":"find something"}`), strPtr("find something")},
		{[]byte(`{"title":"My work item"}`), strPtr("My work item")},
		{[]byte(`{"other":"field"}`), nil},
		{[]byte(`not json`), nil},
	}
	for _, c := range cases {
		got := extractQueryText(c.input)
		if c.want == nil && got != nil {
			t.Errorf("extractQueryText(%s) = %q, want nil", c.input, *got)
		} else if c.want != nil && (got == nil || *got != *c.want) {
			gotStr := "<nil>"
			if got != nil {
				gotStr = *got
			}
			t.Errorf("extractQueryText(%s) = %q, want %q", c.input, gotStr, *c.want)
		}
	}
}

func TestEntrySummary(t *testing.T) {
	qt := "kubernetes pods"
	body := `{"query":"from body"}`

	// query_text wins over body
	e1 := Entry{QueryText: &qt, QueryHash: "abc"}
	if EntrySummary(e1) != "kubernetes pods" {
		t.Errorf("expected query_text to win: %q", EntrySummary(e1))
	}

	// no query_text → body fallback
	e2 := Entry{InputBody: &body, QueryHash: "abc"}
	if EntrySummary(e2) != "from body" {
		t.Errorf("expected body fallback: %q", EntrySummary(e2))
	}

	// no query_text, no body → hash
	e3 := Entry{QueryHash: "abc123"}
	if EntrySummary(e3) != "abc123" {
		t.Errorf("expected hash fallback: %q", EntrySummary(e3))
	}
}

func TestWriteWithBodies(t *testing.T) {
	s := openTemp(t)

	input := []byte(`{"query":"search term"}`)
	output := []byte(`{"content":"result"}`)

	if err := s.Write("hawp_search", input, output, true); err != nil {
		t.Fatalf("Write: %v", err)
	}

	entries, err := s.Recent(1)
	if err != nil || len(entries) == 0 {
		t.Fatalf("Recent: %v / %d entries", err, len(entries))
	}
	if entries[0].InputBody == nil {
		t.Error("input_body should be set when logBodies=true")
	}
	if entries[0].OutputBody == nil {
		t.Error("output_body should be set when logBodies=true")
	}
}

func TestGetTotals(t *testing.T) {
	s := openTemp(t)

	for i := 0; i < 3; i++ {
		_ = s.Write("hawp_search", []byte(`{"query":"test"}`), []byte(`{"content":"abc"}`), false)
	}

	totals, err := s.GetTotals()
	if err != nil {
		t.Fatalf("GetTotals: %v", err)
	}
	if totals.Calls != 3 {
		t.Errorf("expected 3 calls, got %d", totals.Calls)
	}
	if totals.TokensIn <= 0 {
		t.Error("expected tokens_in > 0")
	}
}

func TestClear(t *testing.T) {
	s := openTemp(t)
	_ = s.Write("hawp_search", []byte(`{"query":"x"}`), []byte(`{}`), false)

	if err := s.Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	totals, _ := s.GetTotals()
	if totals.Calls != 0 {
		t.Errorf("expected 0 after clear, got %d", totals.Calls)
	}
}

func TestIdempotentSchema(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "usage.db")

	s1, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	s1.Close()

	s2, err := Open(path)
	if err != nil {
		t.Fatalf("second Open failed: %v", err)
	}
	s2.Close()
}

func TestLoadSaveConfig(t *testing.T) {
	dir := t.TempDir()
	configFile := filepath.Join(dir, "config", "usage.json")

	// missing file → defaults
	c := LoadConfig(configFile)
	if c.Enabled || c.LogBodies {
		t.Error("defaults should have enabled=false, log_bodies=false")
	}

	// save then reload
	c.Enabled = true
	c.LogBodies = true
	if err := SaveConfig(configFile, c); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}
	c2 := LoadConfig(configFile)
	if !c2.Enabled || !c2.LogBodies {
		t.Error("reloaded config should have enabled=true, log_bodies=true")
	}
}

func TestGetReport(t *testing.T) {
	s := openTemp(t)

	tools := []string{"hawp_search", "hawp_search", "hawp_work_new"}
	// large input, small output → guaranteed savings
	longInput := []byte(`{"query":"find all kubernetes deployment configs and service mesh patterns that apply to this workflow"}`)
	shortOutput := []byte(`{"text":"ok"}`)
	for _, tool := range tools {
		if err := s.Write(tool, longInput, shortOutput, false); err != nil {
			t.Fatalf("Write: %v", err)
		}
	}

	rep, err := s.GetReport()
	if err != nil {
		t.Fatalf("GetReport: %v", err)
	}
	if rep.Calls != 3 {
		t.Errorf("calls: got %d want 3", rep.Calls)
	}
	if rep.TokensIn <= 0 {
		t.Error("tokens_in should be > 0")
	}
	if len(rep.ByTool) != 2 {
		t.Errorf("by_tool: got %d distinct tools, want 2", len(rep.ByTool))
	}
	if rep.Since == nil || rep.Until == nil {
		t.Error("since/until should be set when log has entries")
	}
	if len(rep.TopEntries) != 3 {
		t.Errorf("top_entries: got %d, want 3", len(rep.TopEntries))
	}

	// FormatReport smoke test
	out := FormatReport(rep)
	for _, want := range []string{"# HAWP Usage Report", "hawp_search", "hawp_work_new", "Tokens saved"} {
		if !strings.Contains(out, want) {
			t.Errorf("FormatReport missing %q", want)
		}
	}
}

func TestFormatReportEmpty(t *testing.T) {
	out := FormatReport(Report{})
	if out == "" {
		t.Error("FormatReport(empty) should return non-empty string")
	}
}

func TestQuerySummary(t *testing.T) {
	cases := []struct {
		input []byte
		hash  string
		want  string
	}{
		{[]byte(`{"query":"kubernetes pods"}`), "abc123", "kubernetes pods"},
		{[]byte(`{"title":"Fix the thing"}`), "abc123", "Fix the thing"},
		{[]byte(`{"other":"field"}`), "abc123", "abc123"},
		{[]byte(`not json`), "abc123", "abc123"},
	}
	for _, c := range cases {
		got := QuerySummary(c.input, c.hash)
		if got != c.want {
			t.Errorf("QuerySummary(%s) = %q, want %q", c.input, got, c.want)
		}
	}
}

func TestFormatTotals(t *testing.T) {
	out := FormatTotals(Totals{})
	if out == "" {
		t.Error("FormatTotals(empty) should return non-empty string")
	}
	out2 := FormatTotals(Totals{Calls: 5, TokensIn: 1000, TokensOut: 400})
	if out2 == "" {
		t.Error("FormatTotals(non-empty) should return non-empty string")
	}
	_ = os.Stderr // silence unused import
}
