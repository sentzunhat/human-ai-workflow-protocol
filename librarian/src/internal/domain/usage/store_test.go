package usage

import (
	"os"
	"path/filepath"
	"testing"
)

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
