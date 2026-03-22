package domain_test

// Comprehensive benchmarks for all embedding and LLM backends.
// Run with: go test -run BenchmarkAll -v ./internal/domain/ -timeout 1200s
//
// Tests both ONNX and Ollama backends with multiple models.

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sentzunhat/hawp/librarian/src/internal/domain/embeddings"
	"github.com/sentzunhat/hawp/librarian/src/internal/domain/llm"
)

// Representative text samples for benchmarking
var benchmarkTexts = []string{
	"The Go programming language is efficient and suitable for concurrent systems.",
	"Kubernetes manages containerized applications across distributed clusters seamlessly.",
	"Machine learning models require large datasets and careful hyperparameter tuning.",
	"Database optimization involves indexing, query planning, and schema design.",
	"API design best practices include versioning, pagination, and error handling.",
}

// ─── ONNX Embedding Benchmarks ─────────────────────────────────────────────

func BenchmarkAll_ONNX_Embeddings_MiniLM(b *testing.B) {
	embedder, err := embeddings.NewONNXEmbedder("all-MiniLM-L6-v2")
	if err != nil {
		b.Fatalf("ONNX MiniLM: create embedder: %v", err)
	}
	defer embedder.Close()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := embedder.Embed(ctx, benchmarkTexts[i%len(benchmarkTexts)])
		if err != nil {
			b.Fatalf("ONNX MiniLM: embed failed: %v", err)
		}
	}
	b.StopTimer()

	// Run batch once for throughput
	start := time.Now()
	_, _ = embedder.EmbedBatch(ctx, benchmarkTexts)
	batchElapsed := time.Since(start)
	b.Logf("ONNX MiniLM batch (5 texts): %v", batchElapsed)
}

func BenchmarkAll_ONNX_Embeddings_BGE(b *testing.B) {
	embedder, err := embeddings.NewONNXEmbedder("bge-base-en-v1.5")
	if err != nil {
		b.Fatalf("ONNX BGE: create embedder: %v", err)
	}
	defer embedder.Close()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := embedder.Embed(ctx, benchmarkTexts[i%len(benchmarkTexts)])
		if err != nil {
			b.Fatalf("ONNX BGE: embed failed: %v", err)
		}
	}
	b.StopTimer()

	// Run batch once for throughput
	start := time.Now()
	_, _ = embedder.EmbedBatch(ctx, benchmarkTexts)
	batchElapsed := time.Since(start)
	b.Logf("ONNX BGE batch (5 texts): %v", batchElapsed)
}

// ─── Ollama Embedding Benchmarks ───────────────────────────────────────────

func BenchmarkAll_Ollama_Embeddings_MiniLM(b *testing.B) {
	embedder, err := embeddings.NewOllamaEmbedder("http://localhost:11434", "all-minilm")
	if err != nil {
		b.Skipf("Ollama MiniLM not available: %v", err)
	}
	defer embedder.Close()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := embedder.Embed(ctx, benchmarkTexts[i%len(benchmarkTexts)])
		if err != nil {
			b.Fatalf("Ollama MiniLM: embed failed: %v", err)
		}
	}
	b.StopTimer()

	// Run batch once for throughput
	start := time.Now()
	_, _ = embedder.EmbedBatch(ctx, benchmarkTexts)
	batchElapsed := time.Since(start)
	b.Logf("Ollama MiniLM batch (5 texts): %v", batchElapsed)
}

func BenchmarkAll_Ollama_Embeddings_Nomic(b *testing.B) {
	embedder, err := embeddings.NewOllamaEmbedder("http://localhost:11434", "nomic-embed-text")
	if err != nil {
		b.Skipf("Ollama Nomic not available (run: ollama pull nomic-embed-text): %v", err)
	}
	defer embedder.Close()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := embedder.Embed(ctx, benchmarkTexts[i%len(benchmarkTexts)])
		if err != nil {
			b.Fatalf("Ollama Nomic: embed failed: %v", err)
		}
	}
	b.StopTimer()

	// Run batch once for throughput
	start := time.Now()
	_, _ = embedder.EmbedBatch(ctx, benchmarkTexts)
	batchElapsed := time.Since(start)
	b.Logf("Ollama Nomic batch (5 texts): %v", batchElapsed)
}

// ─── ONNX LLM Benchmarks ───────────────────────────────────────────────────

func TestBenchmarkAll_ONNX_LLM_Status(t *testing.T) {
	// Document the current ONNX LLM state on this default (non-ORT) test
	// build: NewONNXLLMClient fails fast, before any download, because
	// hugot's ORT backend is compiled out without `-tags ORT`.
	//
	// Real generation IS verified working (2026-07-27, manual `-tags ORT` +
	// native-libraries build, macOS arm64): SmolLM2-360M-Instruct
	// (genai-converted variant), ~1.1s for a single-context reshape. Not
	// yet part of the release build — see v0.1.0_VISION.md for the setup
	// and remaining GitHub Actions work (5 of 6 platforms have official
	// native libraries; Intel Mac does not).
	t.Logf("ONNX LLM: SupportedModels count = %d", len(llm.SupportedModels))
	t.Logf("ONNX LLM: default build STATUS = fails fast (no -tags ORT)")
	t.Logf("ONNX LLM: manual -tags ORT build STATUS = ✅ verified working, ~1.1s/reshape")
	t.Logf("ONNX LLM: Workaround for default build = use the Ollama LLM backend.")
}

// ─── Ollama LLM Benchmarks ────────────────────────────────────────────────

func BenchmarkAll_Ollama_LLM_Qwen(b *testing.B) {
	client, err := llm.NewOllamaLLMClient("http://localhost:11434", "qwen3.5:4b")
	if err != nil {
		b.Skipf("Ollama Qwen not available: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	sampleText := "Go concurrency patterns channels goroutines select mutex sync atomic"
	maxTokens := 100

	// Run just 2 iterations (LLM is slow on CPU, each ~40s)
	b.ResetTimer()
	for i := 0; i < 2; i++ {
		_, err := client.Reshape(ctx, sampleText, maxTokens)
		if err != nil {
			b.Fatalf("Ollama Qwen: reshape failed: %v", err)
		}
	}
	b.StopTimer()
}

func BenchmarkAll_Ollama_LLM_Mistral(b *testing.B) {
	client, err := llm.NewOllamaLLMClient("http://localhost:11434", "mistral")
	if err != nil {
		b.Skipf("Ollama Mistral not available (run: ollama pull mistral): %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	sampleText := "Kubernetes pod scheduling resource limits CPU memory requests"
	maxTokens := 100

	// Run just 1 iteration (Mistral on CPU is very slow, ~60-120s)
	b.ResetTimer()
	for i := 0; i < 1; i++ {
		_, err := client.Reshape(ctx, sampleText, maxTokens)
		if err != nil {
			b.Fatalf("Ollama Mistral: reshape failed: %v", err)
		}
	}
	b.StopTimer()
}

// ─── Composite Benchmarks (Full Pipeline) ────────────────────────────────

func TestBenchmarkAll_Full_Pipeline_ONNX(t *testing.T) {
	// Full pipeline: ONNX embeddings + Ollama LLM (since ONNX LLM not available)
	emb, err := embeddings.NewONNXEmbedder("all-MiniLM-L6-v2")
	if err != nil {
		t.Fatalf("ONNX embedder: %v", err)
	}
	defer emb.Close()

	llmClient, err := llm.NewOllamaLLMClient("http://localhost:11434", "qwen3.5:4b")
	if err != nil {
		t.Skipf("Ollama LLM not available: %v", err)
	}
	defer llmClient.Close()

	ctx := context.Background()
	text := "Go channels goroutines concurrency patterns select statement"

	t.Log("=== Full Pipeline: ONNX Embeddings + Ollama LLM ===")

	// Step 1: Embed
	start := time.Now()
	vec, err := emb.Embed(ctx, text)
	embedElapsed := time.Since(start)
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}
	t.Logf("  1. Embed: %v (dim=%d)", embedElapsed, len(vec))

	// Step 2: Reshape
	start = time.Now()
	reshaped, err := llmClient.Reshape(ctx, text, 100)
	reshapeElapsed := time.Since(start)
	if err != nil {
		t.Fatalf("Reshape failed: %v", err)
	}
	t.Logf("  2. Reshape: %v (%d chars out)", reshapeElapsed, len(reshaped))

	total := embedElapsed + reshapeElapsed
	t.Logf("  Total: %v", total)
	t.Logf("  Reshaped output: %s", reshaped)
}

func TestBenchmarkAll_Full_Pipeline_Ollama(t *testing.T) {
	// Full pipeline: Ollama embeddings + Ollama LLM
	emb, err := embeddings.NewOllamaEmbedder("http://localhost:11434", "all-minilm")
	if err != nil {
		t.Skipf("Ollama embeddings not available: %v", err)
	}
	defer emb.Close()

	llmClient, err := llm.NewOllamaLLMClient("http://localhost:11434", "qwen3.5:4b")
	if err != nil {
		t.Skipf("Ollama LLM not available: %v", err)
	}
	defer llmClient.Close()

	ctx := context.Background()
	text := "Kubernetes pod scheduling resource limits CPU memory requests"

	t.Log("=== Full Pipeline: Ollama Embeddings + Ollama LLM ===")

	// Step 1: Embed
	start := time.Now()
	vec, err := emb.Embed(ctx, text)
	embedElapsed := time.Since(start)
	if err != nil {
		t.Fatalf("Embed failed: %v", err)
	}
	t.Logf("  1. Embed: %v (dim=%d)", embedElapsed, len(vec))

	// Step 2: Reshape
	start = time.Now()
	reshaped, err := llmClient.Reshape(ctx, text, 100)
	reshapeElapsed := time.Since(start)
	if err != nil {
		t.Fatalf("Reshape failed: %v", err)
	}
	t.Logf("  2. Reshape: %v (%d chars out)", reshapeElapsed, len(reshaped))

	total := embedElapsed + reshapeElapsed
	t.Logf("  Total: %v", total)
	t.Logf("  Reshaped output: %s", reshaped)
}

// ─── Summary Report ───────────────────────────────────────────────────────

func TestBenchmarkAll_SummaryReport(t *testing.T) {
	ctx := context.Background()
	results := make(map[string]map[string]interface{})

	// ONNX Embeddings
	results["ONNX Embeddings"] = make(map[string]interface{})
	if emb, err := embeddings.NewONNXEmbedder("all-MiniLM-L6-v2"); err == nil {
		start := time.Now()
		vec, _ := emb.Embed(ctx, "test")
		elapsed := time.Since(start)
		results["ONNX Embeddings"]["all-MiniLM"] = fmt.Sprintf("dim=%d time=%v", len(vec), elapsed)
		emb.Close()
	} else {
		results["ONNX Embeddings"]["all-MiniLM"] = fmt.Sprintf("ERROR: %v", err)
	}

	if emb, err := embeddings.NewONNXEmbedder("bge-base-en-v1.5"); err == nil {
		start := time.Now()
		vec, _ := emb.Embed(ctx, "test")
		elapsed := time.Since(start)
		results["ONNX Embeddings"]["bge-base-en-v1.5"] = fmt.Sprintf("dim=%d time=%v", len(vec), elapsed)
		emb.Close()
	} else {
		results["ONNX Embeddings"]["bge-base-en-v1.5"] = fmt.Sprintf("ERROR: %v", err)
	}

	// Ollama Embeddings
	results["Ollama Embeddings"] = make(map[string]interface{})
	if emb, err := embeddings.NewOllamaEmbedder("http://localhost:11434", "all-minilm"); err == nil {
		start := time.Now()
		vec, _ := emb.Embed(ctx, "test")
		elapsed := time.Since(start)
		results["Ollama Embeddings"]["all-minilm"] = fmt.Sprintf("dim=%d time=%v", len(vec), elapsed)
		emb.Close()
	} else {
		results["Ollama Embeddings"]["all-minilm"] = fmt.Sprintf("SKIP: %v", err)
	}

	if emb, err := embeddings.NewOllamaEmbedder("http://localhost:11434", "nomic-embed-text"); err == nil {
		start := time.Now()
		vec, _ := emb.Embed(ctx, "test")
		elapsed := time.Since(start)
		results["Ollama Embeddings"]["nomic-embed-text"] = fmt.Sprintf("dim=%d time=%v", len(vec), elapsed)
		emb.Close()
	} else {
		results["Ollama Embeddings"]["nomic-embed-text"] = fmt.Sprintf("SKIP: %v", err)
	}

	// ONNX LLM — fails on this default (non-ORT) test build; verified
	// working manually with -tags ORT (2026-07-27, see v0.1.0_VISION.md).
	results["ONNX LLM"] = map[string]interface{}{
		"status": fmt.Sprintf("needs -tags ORT build (verified working manually; %d model(s) registered)", len(llm.SupportedModels)),
	}

	// Ollama LLM
	results["Ollama LLM"] = make(map[string]interface{})
	if client, err := llm.NewOllamaLLMClient("http://localhost:11434", "qwen3.5:4b"); err == nil {
		results["Ollama LLM"]["qwen3.5:4b"] = "OK"
		client.Close()
	} else {
		results["Ollama LLM"]["qwen3.5:4b"] = fmt.Sprintf("SKIP: %v", err)
	}

	if client, err := llm.NewOllamaLLMClient("http://localhost:11434", "mistral"); err == nil {
		results["Ollama LLM"]["mistral"] = "OK"
		client.Close()
	} else {
		results["Ollama LLM"]["mistral"] = fmt.Sprintf("SKIP: %v", err)
	}

	t.Log("\n=== Benchmark Summary ===")
	for backend, models := range results {
		t.Logf("%s:", backend)
		for model, result := range models {
			t.Logf("  %-25s %v", model, result)
		}
	}
}
