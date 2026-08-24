//go:build integration

package index

// Four live embed+search runs:
//
//  Run 1 — ONNX single (all-MiniLM-L6-v2)
//  Run 2 — Ollama single (nomic-embed-text @ localhost:11434)
//  Run 3 — ONNX full pipeline (ingest → embed → search)
//  Run 4 — Ollama full pipeline (ingest → embed → search)
//
// Run with:
//   go test -v -tags integration -run TestEmbed ./internal/application/index/...

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	appsearch "github.com/sentzunhat/hawp/librarian/src/internal/application/search"
	"github.com/sentzunhat/hawp/librarian/src/internal/domain/embeddings"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/sqlite"
)

const (
	ollamaURL       = "http://localhost:11434"
	ollamaModel     = "nomic-embed-text"
	onnxModel       = "all-MiniLM-L6-v2"
	testCorpusChunk = "HAWP backlog alignment keeps the backlog concise and actionable. Active work is short and current. Recently closed items are capped at 10 entries."
	searchQuery     = "backlog alignment concise"
)

// checkOllama skips the test if Ollama is not reachable.
func checkOllama(t *testing.T) {
	t.Helper()
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(ollamaURL + "/api/tags")
	if err != nil || resp.StatusCode != 200 {
		t.Skipf("Ollama not reachable at %s (start with: ollama serve)", ollamaURL)
	}
	resp.Body.Close()
}

// buildTestDB creates a temp DB, ingest one document, return its path.
func buildTestDB(t *testing.T) string {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "index.sqlite")
	db, err := sqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer db.Close()
	if err := db.InitSchema(); err != nil {
		t.Fatalf("InitSchema: %v", err)
	}
	docID, err := db.InsertDocument("kit", "guide", ".hawp/kit/backlog-alignment.md", "kit/references")
	if err != nil {
		t.Fatalf("InsertDocument: %v", err)
	}
	if err := db.InsertChunk(sqlite.Chunk{
		DocumentID: docID,
		ChunkIdx:   0,
		Text:       testCorpusChunk,
		LineStart:  1,
		LineEnd:    3,
	}); err != nil {
		t.Fatalf("InsertChunk: %v", err)
	}
	return dbPath
}

// ── Run 1: ONNX single model ────────────────────────────────────────────────

func TestEmbedRun1_ONNXSingle(t *testing.T) {
	home, _ := os.UserHomeDir()
	modelPath := filepath.Join(home, ".hawp", "models", "embedding", "sentence-transformers_all-MiniLM-L6-v2")
	if _, err := os.Stat(modelPath); err != nil {
		t.Skipf("ONNX model not downloaded: %s (run: hawp search embed --backend onnx --model %s)", modelPath, onnxModel)
	}

	t.Logf("Run 1 — ONNX single: embedding one text with %s", onnxModel)
	start := time.Now()

	embedder, err := embeddings.NewONNXEmbedder(onnxModel)
	if err != nil {
		t.Fatalf("NewONNXEmbedder: %v", err)
	}
	defer embedder.Close()

	vec, err := embedder.Embed(context.Background(), testCorpusChunk)
	if err != nil {
		t.Fatalf("Embed: %v", err)
	}

	elapsed := time.Since(start)
	if len(vec) != embedder.Dimension() {
		t.Errorf("vector dim = %d, want %d", len(vec), embedder.Dimension())
	}
	nonZero := countNonZero(vec)
	t.Logf("  backend:   onnx/%s", onnxModel)
	t.Logf("  dimension: %d", len(vec))
	t.Logf("  non-zero:  %d/%d", nonZero, len(vec))
	t.Logf("  latency:   %s", elapsed.Round(time.Millisecond))
	if nonZero == 0 {
		t.Error("embedding is all zeros")
	}
}

// ── Run 2: Ollama single model ──────────────────────────────────────────────

func TestEmbedRun2_OllamaSingle(t *testing.T) {
	checkOllama(t)

	t.Logf("Run 2 — Ollama single: embedding one text with %s", ollamaModel)
	start := time.Now()

	embedder, err := embeddings.NewOllamaEmbedder(ollamaURL, ollamaModel)
	if err != nil {
		t.Fatalf("NewOllamaEmbedder: %v", err)
	}
	defer embedder.Close()

	vec, err := embedder.Embed(context.Background(), testCorpusChunk)
	if err != nil {
		t.Fatalf("Embed: %v", err)
	}

	elapsed := time.Since(start)
	if len(vec) != embedder.Dimension() {
		t.Errorf("vector dim = %d, want %d", len(vec), embedder.Dimension())
	}
	nonZero := countNonZero(vec)
	t.Logf("  backend:   ollama/%s", ollamaModel)
	t.Logf("  dimension: %d", len(vec))
	t.Logf("  non-zero:  %d/%d", nonZero, len(vec))
	t.Logf("  latency:   %s", elapsed.Round(time.Millisecond))
	if nonZero == 0 {
		t.Error("embedding is all zeros")
	}
}

// ── Run 3: ONNX full pipeline (ingest → embed → search) ────────────────────

func TestEmbedRun3_ONNXPipeline(t *testing.T) {
	home, _ := os.UserHomeDir()
	modelPath := filepath.Join(home, ".hawp", "models", "embedding", "sentence-transformers_all-MiniLM-L6-v2")
	if _, err := os.Stat(modelPath); err != nil {
		t.Skipf("ONNX model not downloaded: %s", modelPath)
	}

	dbPath := buildTestDB(t)
	svc := NewEmbedService(dbPath)

	t.Logf("Run 3 — ONNX pipeline: embed → search with %s", onnxModel)
	embedStart := time.Now()
	result, err := svc.Execute(context.Background(), "onnx", onnxModel)
	embedElapsed := time.Since(embedStart)
	if err != nil {
		t.Fatalf("Execute (onnx): %v", err)
	}
	if result.ChunksEmbedded != 1 {
		t.Errorf("embedded %d chunks, want 1", result.ChunksEmbedded)
	}
	t.Logf("  embed:     %d chunks in %s", result.ChunksEmbedded, embedElapsed.Round(time.Millisecond))

	// Search phase
	db, err := sqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open db for search: %v", err)
	}
	defer db.Close()

	rows, err := db.QueryChunksLexical(searchQuery, 5)
	if err != nil {
		t.Fatalf("QueryChunksLexical: %v", err)
	}
	t.Logf("  lexical:   %d result(s)", len(rows))

	searchStart := time.Now()
	ranked := appsearch.HybridRank(rows, searchQuery, db, 3)
	searchElapsed := time.Since(searchStart)
	t.Logf("  hybrid:    %d result(s) in %s", len(ranked), searchElapsed.Round(time.Millisecond))

	if len(ranked) == 0 {
		t.Error("search returned no results after embedding")
	} else {
		path := derefPath(ranked[0]["path"])
		t.Logf("  top result: %s (score %.4f)", path, ranked[0]["_hybrid_score"])
		if !strings.Contains(path, "backlog-alignment") {
			t.Errorf("top result path %q does not match inserted document", path)
		}
	}
}

// ── Run 4: Ollama full pipeline (ingest → embed → search) ──────────────────

func TestEmbedRun4_OllamaPipeline(t *testing.T) {
	checkOllama(t)

	dbPath := buildTestDB(t)
	svc := NewEmbedService(dbPath)

	t.Logf("Run 4 — Ollama pipeline: embed → search with %s", ollamaModel)
	embedStart := time.Now()
	result, err := svc.Execute(context.Background(), "ollama", ollamaModel)
	embedElapsed := time.Since(embedStart)
	if err != nil {
		t.Fatalf("Execute (ollama): %v", err)
	}
	if result.ChunksEmbedded != 1 {
		t.Errorf("embedded %d chunks, want 1", result.ChunksEmbedded)
	}
	t.Logf("  embed:     %d chunks in %s", result.ChunksEmbedded, embedElapsed.Round(time.Millisecond))

	// Search phase
	db, err := sqlite.Open(dbPath)
	if err != nil {
		t.Fatalf("open db for search: %v", err)
	}
	defer db.Close()

	rows, err := db.QueryChunksLexical(searchQuery, 5)
	if err != nil {
		t.Fatalf("QueryChunksLexical: %v", err)
	}
	t.Logf("  lexical:   %d result(s)", len(rows))

	searchStart := time.Now()
	ranked := appsearch.HybridRank(rows, searchQuery, db, 3)
	searchElapsed := time.Since(searchStart)
	t.Logf("  hybrid:    %d result(s) in %s", len(ranked), searchElapsed.Round(time.Millisecond))

	if len(ranked) == 0 {
		t.Error("search returned no results after embedding")
	} else {
		path := derefPath(ranked[0]["path"])
		t.Logf("  top result: %s (score %.4f)", path, ranked[0]["_hybrid_score"])
		if !strings.Contains(path, "backlog-alignment") {
			t.Errorf("top result path %q does not match inserted document", path)
		}
	}
}

// derefPath extracts a string from a map value that may be *string or string.
func derefPath(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(*string); ok && s != nil {
		return *s
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

// countNonZero counts non-zero float32 values.
func countNonZero(v []float32) int {
	n := 0
	for _, f := range v {
		if f != 0 {
			n++
		}
	}
	return n
}
