package index

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestIngestRealKitDocuments tests ingesting real kit files with timing.
func TestIngestRealKitDocuments(t *testing.T) {
	// This test reads the actual .hawp/kit/ directory and ingests it.
	repoRoot := findRepoRoot()
	if repoRoot == "" {
		t.Skip("not in a HAWP repo")
	}

	kitPath := filepath.Join(repoRoot, ".hawp", "kit")
	entries, err := os.ReadDir(kitPath)
	if err != nil {
		t.Fatalf("could not read kit directory: %v", err)
	}

	if len(entries) == 0 {
		t.Skip("no kit files found")
	}

	// Build a simple corpus from kit files
	corpus := &EnrichedCorpus{}
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".md" {
			path := filepath.Join(kitPath, entry.Name())
			content, _ := os.ReadFile(path)

			corpus.Documents = append(corpus.Documents, EnrichedDocument{
				Path:       filepath.Join(".hawp/kit", entry.Name()),
				Type:       "guide",
				Category:   "kit",
				FolderRole: "kit",
				Content:    string(content),
				Metadata:   map[string]interface{}{"file": entry.Name()},
			})
		}
	}

	t.Logf("Testing ingest of %d kit files\n", len(corpus.Documents))

	// Create a temp database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	service := NewIngestService(dbPath)

	// Measure ingest time
	start := time.Now()
	result, err := service.Execute(corpus)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("ingest failed: %v", err)
	}

	t.Logf("\nIngest Performance (Real Data):\n")
	t.Logf("  Documents:   %d\n", result.DocumentsProcessed)
	t.Logf("  Chunks:      %d\n", result.ChunksCreated)
	t.Logf("  Bytes:       %d\n", result.BytesIndexed)
	t.Logf("  Time:        %.3f seconds\n", elapsed.Seconds())
	t.Logf("  Throughput:  %.1f KB/s\n", float64(result.BytesIndexed)/(1024*elapsed.Seconds()))
	t.Logf("  Chunks/sec:  %.1f\n", float64(result.ChunksCreated)/elapsed.Seconds())
}

// TestIngestWithVectorEmbedding simulates embedding time (future work).
func TestIngestWithVectorEmbedding(t *testing.T) {
	// This test simulates the future embedding phase.
	// For now, it just measures the ingest overhead.

	corpus := &EnrichedCorpus{
		Documents: []EnrichedDocument{
			{
				Path:       ".hawp/kit/test-doc.md",
				Type:       "guide",
				Category:   "kit",
				FolderRole: "kit",
				Content:    generateLargeContent(1000), // ~1000 words
				Metadata:   map[string]interface{}{},
			},
		},
	}

	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	service := NewIngestService(dbPath)

	start := time.Now()
	result, err := service.Execute(corpus)
	ingestTime := time.Since(start)

	if err != nil {
		t.Fatalf("ingest failed: %v", err)
	}

	// Estimate embedding time (rough: 100 chunks @ ~10ms per chunk with model loaded)
	estimatedEmbeddingTime := time.Duration(result.ChunksCreated) * 10 * time.Millisecond

	t.Logf("\nEmbedding Performance (Estimate):\n")
	t.Logf("  Ingest time:            %.3f seconds\n", ingestTime.Seconds())
	t.Logf("  Chunks to embed:        %d\n", result.ChunksCreated)
	t.Logf("  Est. embedding time:    %.2f seconds (at 10ms/chunk)\n", estimatedEmbeddingTime.Seconds())
	t.Logf("  Total (ingest+embed):   %.2f seconds\n", (ingestTime + estimatedEmbeddingTime).Seconds())
	t.Logf("  \n  (Actual embedding time TBD when ONNX integration lands)\n")
}

func generateLargeContent(wordCount int) string {
	// Generate dummy content with approximately wordCount words
	sample := `This is a sample paragraph with multiple sentences. Each sentence adds to the content length. We use this to simulate realistic document sizes. The content generation process is simple and repeatable for testing purposes. `
	content := ""
	for i := 0; i < wordCount/30; i++ { // ~30 words per sample
		content += sample
	}
	return content
}

func findRepoRoot() string {
	// Try to find the repo root by looking for .hawp/kit/
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	// Walk up the directory tree
	for wd != "/" {
		if _, err := os.Stat(filepath.Join(wd, ".hawp", "kit")); err == nil {
			return wd
		}
		wd = filepath.Dir(wd)
	}
	return ""
}
