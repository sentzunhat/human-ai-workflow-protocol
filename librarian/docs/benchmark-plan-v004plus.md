# Comprehensive Benchmarking Plan for v0.0.4+

**Document Version:** 2026-07-26  
**Target Releases:** v0.0.4, v0.0.5, v0.1.0+ (incremental provider additions)  
**Scope:** All embedding and LLM backend combinations, quality measurement, cost tracking  

---

## Executive Summary

This plan establishes a repeatable benchmarking methodology for evaluating all embedding and LLM provider combinations as they ship in v0.0.4+. It covers:

1. **Performance metrics** — latency, throughput, memory footprint
2. **Quality metrics** — semantic similarity, LLM output coherence, factual accuracy
3. **Cost tracking** — API pricing per provider
4. **Matrix testing** — all embedding × LLM combinations
5. **Test infrastructure** — reproducible harness, dataset management, result storage

---

## Current Baseline (v0.0.3)

### Embeddings Already Measured

| Backend | Model | Dim | Single Embed | Batch (5x) | Notes |
|---|---|---|---|---|---|
| ONNX | all-MiniLM-L6-v2 | 384 | 104ms | ~100ms | Fast, lower quality |
| ONNX | bge-base-en-v1.5 | 768 | 605ms | ~950ms | Higher quality, slower |
| Ollama | all-minilm | 384 | 34ms | ~161ms | HTTP overhead low |
| Ollama | nomic-embed-text | 768 | 35ms | ~164ms | Same latency, higher dim |
| Ollama | mxbai-embed-large | 1024 | 31ms | ~175ms | Largest dim, no slowdown |

### LLM Already Measured

| Backend | Model | Params | 100-Token Reshape | CPU Cost |
|---|---|---|---|---|
| Ollama | qwen3.5:4b | 4B | ~40s | Mac M1 CPU-only |
| Ollama | mistral | 7B (Q4) | ~12.4s (measured on 100 tokens) | Mac M1 CPU-only |

**Note:** Times are on **Mac M1 Max, CPU-only** (no GPU). GPU results would be 2-5x faster. Add to test plan for future.

---

## Part 1: Dimensions to Measure

### 1.1 Embedding Model Dimensions

#### Performance (Required)

- **Single embed latency** (ms)
  - Measure: Time for `Embed()` call on a single text
  - Report: min/median/max/stdev across 10 runs
  - Text sample: 100-character, 500-character, 2000-character samples (different sizes affect tokenization)

- **Batch embed latency** (ms)
  - Measure: Time for `EmbedBatch()` on N texts (5, 10, 50)
  - Report: Total time + per-text average
  - Key question: Does batching improve efficiency? By how much?

- **Throughput** (texts/sec)
  - Calculated: texts_per_second = batch_size / batch_latency
  - Compare: ONNX (in-process) vs. Ollama (HTTP) vs. cloud APIs

- **Memory footprint** (MB)
  - Measure: Process resident memory after model load
  - Report: Model memory only (subtract baseline)
  - Context: Important for deployment constraints

- **Model download size** (MB)
  - Report: First-run download size from source
  - Compare: ONNX (Hugot) vs. Ollama (GGUF) vs. OpenAI (API) vs. Anthropic (API)

#### Quality (Required)

- **Vector dimension** (int)
  - Report: Output dimension of each model
  - Question: Does higher dimension always improve quality, or is there a sweet spot?

- **Semantic coherence** (score 0-1)
  - Measure: Cosine similarity between semantically similar pairs (e.g., "car" vs. "automobile")
  - Test set: 20 semantic pairs (see Test Data section)
  - Acceptance: Similarity > 0.7 for related terms, < 0.3 for unrelated
  - Report: Median similarity for related pairs, median dissimilarity for unrelated

- **Embedding stability** (score 0-1)
  - Measure: Cosine similarity between two embeddings of the same text
  - Repeat 5 times per model
  - Question: Are results deterministic across calls?
  - Acceptance: Correlation > 0.99 (drift indicates numerical instability or randomization)

#### Context-Specific (When Available)

- **Language coverage** — does the model handle non-English equally well?
- **Domain bias** — does it favor certain domains (e.g., software terminology)?
- **Multilingual quality** (for multilingual models)

---

### 1.2 LLM Model Dimensions

#### Performance (Required)

- **Latency by token budget** (ms)
  - Measure: Time to complete `Reshape()` with varying `max_tokens`
  - Token budgets: 50, 100, 200, 500, 1000
  - Report: Latency, tokens/sec, total tokens generated
  - Question: Is latency linear or superlinear with token budget?

- **Cold start vs. warm** (ms)
  - Measure: First call after model load vs. subsequent calls
  - Context: Some models cache; others don't
  - Report: First call / warm call + ratio

- **Throughput** (tokens/sec)
  - Calculated: tokens_generated / latency
  - Compare across models and backends

- **Memory footprint** (MB)
  - Measure: Resident memory with model loaded
  - Report: Model memory + context window size
  - Question: Can we quantize to fit mobile/edge?

- **Model size / download** (MB)
  - Report: Quantization level (if applicable)
  - For Ollama: Report quantization (Q4_K_M, Q8_0, etc.)
  - For ONNX: Report actual model file size

#### Quality (Required — Objective Metrics)

- **Output coherence** (0-5 Likert scale)
  - Measure: Is the reshaped text grammatically correct and readable?
  - Criteria: No truncation mid-word, no obvious token corruption, flows naturally
  - Scorer: Automated check for length + basic NLP metrics (no human scoring initially)
  - Acceptance: Score >= 4 for production use

- **Factual accuracy retention** (0-1 percentage)
  - Measure: Does reshaping preserve factual claims from the original?
  - Method: Given a block with N facts, how many are preserved/present in the output?
  - Test cases: 10 domain-specific passages with 3-5 facts each (see Test Data)
  - Acceptance: >= 80% facts retained

- **Context relevance** (0-1 cosine similarity)
  - Measure: Similarity between embedded context + reshaped output
  - Rationale: Output should stay semantically close to input
  - Test: Embed input + output, compute cosine similarity
  - Acceptance: >= 0.6

- **Conciseness trade-off** (ratio)
  - Measure: Original word count vs. reshaped word count
  - Expected: 30-50% of original (summarization effect)
  - Report: Mean compression ratio
  - Question: Do some models over-truncate?

#### Context-Specific (When Available)

- **Hallucination rate** (0-1 percentage)
  - Measure: Does the model add information not in the original?
  - Manual verification on small sample
  - Acceptance: < 5% for production use

- **Latency under load** (ms)
  - Measure: Latency with N concurrent requests (5, 10, 20)
  - Context: For cloud APIs and server-based backends
  - Question: Do they scale linearly or degrade?

- **Error rate under load** (percentage)
  - Measure: Timeouts, rate limit hits, failures
  - Context: For cloud APIs especially

---

## Part 2: Complete Combination Matrix

### 2.1 Planned Provider Roadmap

#### v0.0.3 (Current — ✅ Shipped)

| Embeddings | LLM | Status |
|---|---|---|
| ONNX (MiniLM, BGE) | Ollama (mistral, qwen) | ✅ Live |
| Ollama (all-minilm, nomic, mxbai) | Ollama (mistral, qwen) | ✅ Live |

#### v0.1.0 Planned

| Embeddings | LLM | Status |
|---|---|---|
| ONNX (MiniLM, BGE) | ONNX (text2text: FLAN-T5, etc.) | 🔮 ONNX LLM |
| Ollama (all-minilm, nomic, mxbai) | ONNX (text2text: FLAN-T5, etc.) | 🔮 Hybrid |
| OpenAI (text-embedding-3-small/large) | OpenAI (gpt-3.5, gpt-4) | 🔮 Cloud |
| Anthropic (future API) | Anthropic (claude-3-sonnet, opus) | 🔮 Cloud |

#### v0.2.0+ (Future)

| Embeddings | LLM | Status |
|---|---|---|
| Azure OpenAI | Azure OpenAI | 🔮 Enterprise |
| HuggingFace Inference API | HuggingFace Inference API | 🔮 Alt Cloud |
| Any | Any | 🔮 Full matrix |

### 2.2 Benchmark Matrix Heatmap (Quality)

```
                    ONNX LLM  Ollama LLM  OpenAI LLM  Anthropic LLM
ONNX Embeddings     TBD       ✅ TESTED   TBD         TBD
Ollama Embeddings   TBD       ✅ TESTED   TBD         TBD
OpenAI Embeddings   TBD       TBD         TBD         TBD
Anthropic Embed.    TBD       TBD         TBD         TBD
```

**Key Question:** Which combinations are "stable" (good quality + predictable performance)?

Expected answer (hypothesis):
- **High quality:** OpenAI/Anthropic embeddings + their own LLM
- **Balanced:** Ollama embeddings + Ollama LLM (good enough, local only)
- **Niche:** ONNX embeddings + ONNX LLM (when available — full offline)
- **Risky:** Mixed cloud + local (API dependencies, debugging harder)

### 2.3 Pairing Strategy

For each new backend added, benchmark these pairings:

1. **Self-pairing** (e.g., OpenAI embeddings + OpenAI LLM)
   - Expectation: Vendor optimizes for their own models
   - Measure: Quality, latency, cost

2. **Cross-pairing with Ollama** (e.g., OpenAI embeddings + Ollama LLM)
   - Expectation: Useful for cost control (cheaper local LLM, premium embeddings)
   - Measure: Quality still good? Latency acceptable?

3. **Cross-pairing with existing** (e.g., OpenAI embeddings + existing Ollama baseline)
   - Expectation: Shows impact of better embeddings on existing LLM
   - Measure: Does semantic quality improve output?

---

## Part 3: Testing Methodology

### 3.1 Test Data Repository

Store test data in `.hawp/benchmarks/test-data/`:

```
.hawp/benchmarks/test-data/
├── README.md                              # How to use this data
├── semantic-pairs.json                    # 20 related/unrelated word pairs
├── passages-with-facts.json              # 10 passages, each tagged with 3-5 facts
├── domain-queries/
│   ├── kubernetes.txt                     # Domain-specific sample passages
│   ├── go-concurrency.txt
│   ├── ml-fundamentals.txt
│   └── database-optimization.txt
└── evaluation/
    ├── human-scores.json                 # (Future) Ground truth for quality
    └── reference-embeddings.json         # (Future) Pre-computed for regression
```

#### Semantic Pairs (semantic-pairs.json)

```json
[
  {"pair": ["cat", "feline"], "expected": "high", "domain": "animals"},
  {"pair": ["car", "automobile"], "expected": "high", "domain": "vehicles"},
  {"pair": ["cat", "algorithm"], "expected": "low", "domain": "unrelated"},
  ...
]
```

#### Passages with Facts (passages-with-facts.json)

```json
[
  {
    "domain": "kubernetes",
    "text": "Kubernetes uses pods as the smallest deployable unit. Each pod runs one or more containers. Pods are ephemeral and can be replaced dynamically.",
    "facts": ["pods are smallest deployable unit", "pods contain containers", "pods are ephemeral"]
  },
  ...
]
```

### 3.2 Benchmark Harness Structure

Create `librarian/src/benchmarks/` directory with:

```
librarian/src/benchmarks/
├── harness.go                  # Main benchmarking engine
├── metrics.go                  # Metric collection / aggregation
├── reporters.go                # Output formatting (CSV, JSON, markdown)
├── quality_test.go            # Objective quality checks
├── suite_test.go              # Full matrix runner
├── testdata/                  # Link to .hawp/benchmarks/test-data/
└── results/
    └── 2026-07-26-v003-baseline.json  # Saved runs for regression
```

#### Harness Example (harness.go pseudo-code)

```go
type BenchmarkSuite struct {
    Embedders map[string]Embedder
    LLMs map[string]LLMClient
    TestData *TestDataSet
    Results *BenchmarkResults
}

func (bs *BenchmarkSuite) RunMatrix() {
    for embName, emb := range bs.Embedders {
        for llmName, llm := range bs.LLMs {
            bs.RunPairing(embName, llmName, emb, llm)
        }
    }
}

func (bs *BenchmarkSuite) RunPairing(embName, llmName string, emb, llm ...) {
    // Single embed latency
    // Batch embed latency
    // Semantic coherence
    // Reshape latency (multiple token budgets)
    // Output quality checks
    // Record results
}

func (bs *BenchmarkSuite) GenerateReport() markdown {
    // Tables: performance, quality, cost
    // Heatmaps: quality matrix
    // Recommendations: which pairings are production-ready
}
```

### 3.3 Running Benchmarks

#### Full Suite

```bash
# Run all combinations (requires all backends available)
go test -v -run TestBenchmarkSuite_Matrix ./benchmarks/... \
    -timeout 3600s \
    -tags benchmark

# Saves results to: librarian/src/benchmarks/results/YYYY-MM-DD-suite.json
```

#### Per-Backend

```bash
# Only ONNX embeddings
go test -v -run TestBenchmarkSuite_ONNX_Embed ./benchmarks/...

# Only Ollama (requires: ollama serve)
go test -v -run TestBenchmarkSuite_Ollama ./benchmarks/... \
    -tags requireollama

# Only cloud APIs (requires: API keys in env)
go test -v -run TestBenchmarkSuite_Cloud ./benchmarks/... \
    -tags requirecloud \
    -env OPENAI_API_KEY,ANTHROPIC_API_KEY
```

#### Continuous Regression

```bash
# Compare against v0.0.3 baseline
go test -v ./benchmarks/... \
    -tags benchmark,regression \
    -baseline librarian/src/benchmarks/results/2026-07-26-v003-baseline.json
```

---

## Part 4: Success Criteria & Thresholds

### 4.1 Embedding Model Acceptance

A new embedding model is **production-ready** if:

- [x] Single embed latency < 1000ms (allows local + cloud models)
- [x] Semantic coherence >= 0.7 for related terms
- [x] Embedding stability > 0.99 (deterministic)
- [x] Download size < 2GB (practical for distribution)
- [x] No deadlocks, timeouts, or crashes on test data

### 4.2 LLM Model Acceptance

A new LLM is **production-ready** if:

- [x] 100-token reshape latency < 30 seconds (user-acceptable for CLI)
- [x] Output coherence >= 4/5 (readable, no corruption)
- [x] Fact retention >= 80% (preserves original content)
- [x] Context relevance >= 0.6 (semantically similar to input)
- [x] No hallucinations on clean inputs (< 5% rate)

### 4.3 Combination Acceptance

A **pairing is viable** if:

- [x] Both endpoints work independently (backend acceptance)
- [x] Combined latency < 60 seconds (total pipeline usable)
- [x] Quality metrics both pass (embeddings + LLM)
- [x] No degradation when used together

**Recommended pairing**: All metrics pass + cost is reasonable for use case

**Caution pairing**: Metrics pass but quality slightly lower (e.g., some hallucinations)

**Unstable pairing**: High latency, quality issues, or frequent errors — requires investigation

---

## Part 5: Storage & Reporting

### 5.1 Results Storage

Save each benchmark run to:

```
.hawp/benchmarks/results/
└── YYYY-MM-DD_release-version/
    ├── README.md                         # Run metadata, hardware, environment
    ├── summary.json                      # Aggregate metrics
    ├── detailed_embeddings.json          # Per-model metrics
    ├── detailed_llm.json                 # Per-model metrics
    ├── detailed_pairings.json            # Combined results
    ├── quality_checks.json               # Coherence, accuracy, etc.
    ├── cost_analysis.json                # API costs (for cloud models)
    └── regression_vs_baseline.json       # Comparison to v0.0.3
```

### 5.2 Report Template (README.md)

```markdown
# Benchmark Run: {date} {version}

## Hardware & Environment
- Machine: {CPU/GPU details}
- OS: {OS version}
- Go version: {version}
- Backend service versions: {Ollama version, OpenAI API date, etc.}

## Summary
- Embeddings tested: N models across M backends
- LLMs tested: N models across M backends
- Total pairings: N
- Execution time: XXm

## Performance Highlights
[Table: fastest models per category]

## Quality Highlights
[Table: highest semantic coherence, best fact retention, etc.]

## Recommendations
- **Production-ready**: [list of stable pairings]
- **Caution**: [list of working but lower-quality pairings]
- **Unstable**: [list requiring investigation]

## Regressions
[Any metrics worse than v0.0.3 baseline?]

## Notes
[Any anomalies, hardware constraints, skipped tests?]
```

### 5.3 Update Documentation

After each benchmark run, update:

- `librarian/docs/backends.md` — performance characteristics table
- `librarian/docs/benchmarks-vXXX.md` — latest results
- `librarian/src/README.md` — quick performance reference

---

## Part 6: Quality Measurement Details

### 6.1 Semantic Coherence Test

```go
func TestEmbeddingCoherence(t *testing.T, embedder Embedder, testPairs []SemanticPair) {
    related := 0
    unrelated := 0
    for _, pair := range testPairs {
        vec1, _ := embedder.Embed(pair.Term1)
        vec2, _ := embedder.Embed(pair.Term2)
        sim := CosineSimilarity(vec1, vec2)
        
        if pair.Expected == "high" && sim > 0.7 {
            related++
        } else if pair.Expected == "low" && sim < 0.3 {
            unrelated++
        }
    }
    score := (related + unrelated) / len(testPairs)
    t.Logf("Coherence: %.2f (related: %d/%d, unrelated: %d/%d)", 
        score, related, relatedCount, unrelated, unrelatedCount)
}
```

### 6.2 Fact Retention Test

```go
func TestFactRetention(t *testing.T, llm LLMClient, passages []PassageWithFacts) {
    for _, passage := range passages {
        reshaped, _ := llm.Reshape(passage.Text, 200)
        
        retained := 0
        for _, fact := range passage.Facts {
            if strings.Contains(reshaped, fact) {
                retained++
            }
        }
        score := float64(retained) / float64(len(passage.Facts))
        t.Logf("Passage %s: %d/%d facts retained (%.1f%%)", 
            passage.Domain, retained, len(passage.Facts), score*100)
    }
}
```

### 6.3 Context Relevance Test

```go
func TestContextRelevance(t *testing.T, embedder Embedder, llm LLMClient, passages []string) {
    for _, passage := range passages {
        originalVec, _ := embedder.Embed(passage)
        reshaped, _ := llm.Reshape(passage, 200)
        reshapedVec, _ := embedder.Embed(reshaped)
        
        sim := CosineSimilarity(originalVec, reshapedVec)
        t.Logf("Context relevance: %.3f", sim)
    }
}
```

---

## Part 7: Cost Tracking (For Cloud APIs)

### 7.1 Cost Model

For each API-based backend, track:

```go
type CostMetrics struct {
    BackendName    string    // "openai", "anthropic"
    Model          string    // "gpt-4", "claude-3-opus"
    InputTokens    int       // Tokens sent
    OutputTokens   int       // Tokens generated
    APICallCount   int       // Number of API calls
    TotalCost      float64   // USD
    CostPerCall    float64   // USD/call
}
```

### 7.2 Benchmark Cost Report

```markdown
## Cost Analysis (v0.1.0+ with Cloud APIs)

| Backend | Model | Calls | Input Tokens | Output Tokens | Total Cost | $/Call |
|---|---|---|---|---|---|---|
| OpenAI | gpt-3.5-turbo | 50 | 25,000 | 5,000 | $0.08 | $0.0016 |
| Anthropic | claude-3-sonnet | 50 | 25,000 | 5,000 | $0.15 | $0.003 |
| Ollama | mistral | 50 | N/A | N/A | $0.00 | $0.00 |

**Recommendation:** For budget-sensitive workloads, Ollama is free but slower. For speed, consider OpenAI.
```

---

## Part 8: Integration with CI/CD

### 8.1 Benchmark CI Job

Add to `.github/workflows/benchmarks.yml`:

```yaml
name: Benchmark Suite

on:
  schedule:
    - cron: '0 0 * * *'  # Daily
  workflow_dispatch:

jobs:
  benchmark:
    runs-on: [self-hosted, macos, arm64]  # Stable hardware
    steps:
      - uses: actions/checkout@v4
      
      - name: Start Ollama
        run: ollama serve > /tmp/ollama.log 2>&1 &
      
      - name: Pull Models
        run: |
          ollama pull nomic-embed-text
          ollama pull mistral
      
      - name: Run Benchmarks
        run: |
          go test -v ./benchmarks/... \
            -tags benchmark,requireollama \
            -timeout 3600s \
            -json > bench_result.json
      
      - name: Generate Report
        run: |
          go run ./benchmarks/cmd/reporter.go \
            -input bench_result.json \
            -baseline librarian/src/benchmarks/results/v003-baseline.json \
            -output librarian/docs/benchmarks-latest.md
      
      - name: Commit Results
        run: |
          git add librarian/src/benchmarks/results/
          git add librarian/docs/benchmarks-latest.md
          git commit -m "ci: benchmark run $(date +%Y-%m-%d)" || true
          git push
```

### 8.2 Regression Detection

If any metric regresses by > 10%, file an issue:

```go
func checkRegressions(baseline, current *BenchmarkResults) []Issue {
    var issues []Issue
    for model, metrics := range current.Embeddings {
        baseMetrics := baseline.Embeddings[model]
        latencyChange := (metrics.LatencyMs - baseMetrics.LatencyMs) / baseMetrics.LatencyMs
        
        if latencyChange > 0.1 {  // 10% slower
            issues = append(issues, Issue{
                Title: fmt.Sprintf("Embedding %s regressed %.1f%%", model, latencyChange*100),
                Severity: "warn",
            })
        }
    }
    return issues
}
```

---

## Part 9: Research Questions to Explore

As benchmarks accumulate, investigate:

1. **Embedding dimension vs. quality trade-off**
   - Do higher-dimension embeddings always help LLM quality?
   - Or is there an optimal sweet spot (e.g., 768-dim) beyond which gains plateau?

2. **Latency vs. quality trade-off**
   - Can we quantize large models (e.g., Q4 vs. Q8) without major quality loss?
   - Example: Q4 mistral 30% faster than Q8 mistral — is quality acceptable?

3. **Batch efficiency**
   - ONNX: Does batching reduce per-item latency? By how much?
   - Ollama: HTTP overhead fixed? Does batch size matter?
   - Cloud APIs: Always cheaper per token in batches?

4. **Cold start impact**
   - For Ollama: Does model load time dominate first request?
   - For cloud APIs: Is first call slower?
   - Practical advice for caching strategies?

5. **Domain performance**
   - Are some embeddings better for code than text?
   - Do LLMs preserve domain terminology better than others?
   - Test: Kubernetes vs. Go concurrency vs. databases

6. **GPU vs. CPU performance**
   - v0.0.3 baseline is CPU-only (Mac M1)
   - Add GPU benchmarks when available (NVIDIA, Metal)
   - Expected: 2-5x speedup for both ONNX and Ollama

---

## Part 10: Timeline & Milestones

| Phase | Version | Work | Expected Date |
|---|---|---|---|
| v0.0.4 | Minor providers | Add ONNX LLM scaffolding → v0.0.4 | 2026-08-06 |
| v0.0.5 | More ONNX | Benchmark ONNX LLM (if models available) | 2026-08-13 |
| v0.1.0 | Cloud APIs | OpenAI/Anthropic backends + benchmarks | 2026-09-15 |
| v0.1.0+ | Continuous | Monthly regression checks, new model testing | Ongoing |

---

## Part 11: Roles & Responsibilities

| Role | Responsibility |
|---|---|
| **Benchmark Lead** | Maintain harness, run suite, generate reports |
| **QA** | Manual quality spot-checks, edge case testing |
| **DevOps** | CI/CD setup, stable hardware provisioning, secrets management (API keys) |
| **Platform Lead** | Prioritize new backends, interpret trade-offs, make pairing recommendations |

---

## Part 12: Example: First v0.1.0 Benchmark Run

### Scenario: OpenAI Embeddings Added

**Steps:**

1. **Implement OpenAI embeddings backend** (`openai_embedder.go`)
   - Factory integration
   - Error handling for API failures
   - Rate limiting (1000 req/min)

2. **Add to test data:**
   ```go
   var allBackends = map[string]EmbedderFactory{
       "onnx": NewONNXEmbedder,
       "ollama": NewOllamaEmbedder,
       "openai": NewOpenAIEmbedder,  // NEW
   }
   ```

3. **Run benchmarks:**
   ```bash
   OPENAI_API_KEY=sk-xxx go test -v ./benchmarks/... \
       -run TestBenchmarkSuite_Cloud -tags requirecloud
   ```

4. **Results (hypothetical):**

   | Backend | Model | Latency | Cost/1M Tokens | Quality Score |
   |---|---|---|---|---|
   | ONNX | all-MiniLM | 104ms | $0 | 0.78 |
   | Ollama | nomic | 35ms | $0 | 0.82 |
   | OpenAI | text-embedding-3-small | 450ms | $0.02 | 0.91 |
   | OpenAI | text-embedding-3-large | 650ms | $0.13 | 0.93 |

5. **Recommendations:**
   - OpenAI 3-small: Good quality, affordable, but slower than local options
   - OpenAI 3-large: Best quality, highest cost, use for precision tasks only
   - Default remains ONNX (free, offline)

6. **Update docs:**
   - Add OpenAI to backends.md
   - Update benchmarks-v010.md with full comparison
   - Add pairing recommendations to README

---

## Part 13: Checklist for Future Maintainers

When adding a new backend (Anthropic, HF, Azure, etc.):

- [ ] Implement provider factory interface
- [ ] Add to `allBackends` map in harness
- [ ] Write unit tests (mocks)
- [ ] Write integration test (real service, requires flag/env var)
- [ ] Add to benchmark suite
- [ ] Run full matrix (all existing backends + new one)
- [ ] Verify no regressions in existing benchmarks
- [ ] Spot-check quality manually (5-10 examples)
- [ ] Update cost model if API-based
- [ ] Document results in benchmarks-vXXX.md
- [ ] Update backends.md pairing recommendations
- [ ] File follow-up issues for deep dives (e.g., "Investigate why OpenAI+Ollama pairing has low relevance score")

---

## Appendix A: Test Data Reference

### Semantic Pairs
- Animals: cat/feline, dog/canine, bear/ursine
- Vehicles: car/automobile, bike/bicycle, truck/lorry
- Technology: code/program, error/bug, cache/memory
- Unrelated: cat/algorithm, computer/banana, dog/database

### Sample Passages (Diverse Domains)

**Kubernetes (Infrastructure):**
```
Kubernetes uses pods as the smallest deployable unit. Each pod runs one or more containers. Pods are ephemeral and can be replaced dynamically. Services provide stable endpoints for pod access across cluster restarts.
```

**Go Concurrency (Programming):**
```
Go concurrency is built on goroutines and channels. Goroutines are lightweight threads managed by the Go runtime. Channels enable safe communication between goroutines. The select statement multiplexes channel operations.
```

**ML Fundamentals (Data Science):**
```
Machine learning models learn patterns from labeled data. Training adjusts model parameters to minimize loss. Overfitting occurs when a model memorizes training data rather than generalizing. Regularization techniques prevent overfitting.
```

---

## Appendix B: Useful Metrics Formulas

### Cosine Similarity
```
similarity(u, v) = (u · v) / (||u|| * ||v||)
Range: [-1, 1], where 1 = identical, 0 = orthogonal, -1 = opposite
```

### Throughput
```
throughput = items / time_seconds
Units: texts/sec, tokens/sec, etc.
```

### Compression Ratio
```
ratio = output_tokens / input_tokens
Example: 0.35 = 35% of original size
```

### Quality Score (0-1)
```
score = (met_criteria) / (total_criteria)
Example: 4/5 coherence checks pass = 0.8 score
```

### Cost per Call (USD)
```
cost = (input_tokens * cost_per_input_token) + (output_tokens * cost_per_output_token)
Example: GPT-3.5: $0.0015/1k input, $0.002/1k output
```

---

## Appendix C: Glossary

| Term | Definition |
|---|---|
| **Embedder** | A model that converts text to a vector (dense representation) |
| **LLM Client** | A model that takes text input and produces reshaped/summarized text |
| **Backend** | The infrastructure running the model (ONNX in-process, Ollama HTTP, OpenAI API, etc.) |
| **Pairing** | A specific (Embedder, LLM) combination tested together |
| **Coherence** | Quality of semantic similarity between related terms |
| **Fact Retention** | Percentage of factual claims preserved after reshaping |
| **Relevance** | Semantic similarity between original and reshaped text |
| **Cold Start** | Latency of first call after model load |
| **Warm Call** | Latency of subsequent calls (after model is cached) |
| **Quantization** | Reducing model precision (Q4_K_M vs. Q8_0) to reduce size/latency |
| **Hallucination** | Model generating information not present in the input |

---

## Appendix D: Tools & Commands

### Install Dependencies

```bash
# Go testing
go install gotest.tools/gotestsum@latest

# Ollama (if testing locally)
brew install ollama
ollama serve

# OpenAI CLI (optional, for manual testing)
pip install openai

# Result analysis
pip install pandas matplotlib
```

### Useful Commands

```bash
# Run benchmarks with JSON output
go test -v ./benchmarks/... -json > results.json

# Parse JSON results
jq '.[] | select(.Action=="output") | .Output' results.json

# Compare two benchmark runs
diff <(jq '.results | sort_by(.model)' v003.json) \
     <(jq '.results | sort_by(.model)' v010.json)

# Watch Ollama logs
tail -f ~/.ollama/logs/server.log

# Test API key setup
export OPENAI_API_KEY=sk-xxx
go run ./benchmarks/cmd/test-connection.go
```

---

## Appendix E: Known Issues & Caveats

| Issue | Context | Workaround |
|---|---|---|
| ONNX BGE-large deadlock | Hugot dependency bug | Reverted; use bge-base-en-v1.5 instead |
| Ollama cold start | Model load on first call | Measure warm latency separately; document in report |
| Cloud API rate limits | OpenAI/Anthropic throttle requests | Add backoff to benchmark harness; note in results |
| Determinism | LLMs may use sampling (temperature > 0) | Set temperature=0 for reproducible benchmarks |
| GPU variance | M1 Metal accelerated ops | Document hardware; note GPU results separately |

---

**Last Updated:** 2026-07-26  
**Next Review:** 2026-08-31 (after v0.1.0 planning)
