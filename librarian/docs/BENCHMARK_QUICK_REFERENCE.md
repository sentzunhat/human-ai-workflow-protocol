# Benchmark Quick Reference

**Full Plan:** [`BENCHMARK_PLAN_v004plus.md`](BENCHMARK_PLAN_v004plus.md)  
**Work Item:** `.hawp/work/active/t5l7m9n1-benchmark-plan-v004.md`  
**Current Baseline:** [`BENCHMARKS_v003.md`](BENCHMARKS_v003.md)  
**Backend Status:** [`BACKENDS.md`](BACKENDS.md)

---

## At a Glance: What to Measure

### Embeddings (per model)

| Metric | Purpose | Acceptance |
|---|---|---|
| Single embed latency | Speed | < 1000ms |
| Batch embed latency (5, 10, 50 texts) | Throughput | TBD |
| Semantic coherence (test pairs) | Quality | >= 0.7 |
| Embedding stability | Determinism | > 0.99 |
| Memory footprint | Deployment constraints | Varies |
| Download size | Distribution | < 2GB |

### LLMs (per model)

| Metric | Purpose | Acceptance |
|---|---|---|
| Latency by token budget (50–1000) | Speed curve | 100-tok < 30s |
| Cold start vs. warm latency | Caching behavior | Report both |
| Output coherence | Quality | >= 4/5 |
| Fact retention | Accuracy | >= 80% |
| Context relevance | Semantic fidelity | >= 0.6 |
| Memory footprint | Deployment | Varies |

---

## Combination Matrix (Future Roadmap)

```
                    ONNX LLM  Ollama LLM  OpenAI LLM  Anthropic LLM
ONNX Embeddings     🔮        ✅          🔮         🔮
Ollama Embeddings   🔮        ✅          🔮         🔮
OpenAI Embeddings   🔮        🔮          🔮         🔮
Anthropic Embed.    🔮        🔮          🔮         🔮
```

**✅ = Tested in v0.0.3**  
**🔮 = Planned for v0.1.0+**

---

## Quick Start: Running Benchmarks

### Full Suite (All Combinations)

```bash
# Requires: Ollama running (ollama serve)
# Optional: OpenAI/Anthropic API keys in env

go test -v ./benchmarks/... \
    -tags benchmark \
    -timeout 3600s \
    -json > results.json
```

### ONNX Only (No Dependencies)

```bash
go test -v -run ".*ONNX.*" ./benchmarks/... \
    -tags benchmark \
    -timeout 600s
```

### Ollama Only (Requires: ollama serve)

```bash
go test -v -run ".*Ollama.*" ./benchmarks/... \
    -tags benchmark,requireollama \
    -timeout 1800s
```

### Cloud APIs (Requires: API Keys)

```bash
OPENAI_API_KEY=sk-xxx \
ANTHROPIC_API_KEY=sk-xxx \
go test -v -run ".*Cloud.*" ./benchmarks/... \
    -tags benchmark,requirecloud \
    -timeout 900s
```

### Regression vs. v0.0.3 Baseline

```bash
go test -v ./benchmarks/... \
    -tags benchmark,regression \
    -baseline librarian/go/benchmarks/results/v003-baseline.json
```

---

## Test Data Location

```
.hawp/benchmarks/test-data/
├── semantic-pairs.json              # 20 word pairs (high/low similarity)
├── passages-with-facts.json         # 10 domain passages with fact tags
└── domain-queries/
    ├── kubernetes.txt
    ├── go-concurrency.txt
    ├── ml-fundamentals.txt
    └── database-optimization.txt
```

---

## Result Storage

Each run saves to:

```
.hawp/benchmarks/results/YYYY-MM-DD_version/
├── README.md                    # Metadata, hardware, environment
├── summary.json                 # Aggregate metrics
├── detailed_embeddings.json     # Per-model results
├── detailed_llm.json            # Per-model results
├── detailed_pairings.json       # Combined backend results
├── quality_checks.json          # Coherence, fact retention, etc.
├── cost_analysis.json           # API costs (cloud models)
└── regression_vs_baseline.json  # Comparison to previous run
```

---

## Key Thresholds (Production Readiness)

### Embedding Acceptance

- ✅ Single embed < 1000ms
- ✅ Semantic coherence >= 0.7
- ✅ Stability > 0.99
- ✅ No crashes/deadlocks

### LLM Acceptance

- ✅ 100-token reshape < 30s
- ✅ Coherence >= 4/5
- ✅ Fact retention >= 80%
- ✅ Relevance >= 0.6

### Pairing Acceptance

- ✅ Both backends pass independently
- ✅ Combined latency < 60s
- ✅ Quality metrics pass

---

## Timeline

| Phase | When | Work |
|---|---|---|
| **1. Infrastructure** | Before v0.0.4 | Create test data, scaffold harness |
| **2. v0.0.4 Development** | During v0.0.4 | Add ONNX LLM benchmarks, regression tests |
| **3. v0.1.0 Preparation** | Before v0.1.0 | CI/CD setup, cloud API support, full matrix |
| **4. Continuous** | v0.1.0+ | Daily/weekly benchmarks, regression detection |

---

## Useful Commands

### Parse JSON Results

```bash
# Show all models and latencies
jq '.results[] | {model: .model, latency: .latency_ms}' results.json

# Export to CSV
jq -r '.results[] | [.backend, .model, .latency_ms] | @csv' results.json > results.csv

# Find regressions (> 10% slower than baseline)
jq '.regressions[] | select(.change > 0.1)' regression_vs_baseline.json
```

### Monitor Ollama

```bash
# Check available models
curl http://localhost:11434/api/tags | jq '.models[].name'

# Watch inference logs
tail -f ~/.ollama/logs/server.log

# Pull a model
ollama pull nomic-embed-text
```

### Compare Two Runs

```bash
diff <(jq '.results | sort_by(.model)' run1.json) \
     <(jq '.results | sort_by(.model)' run2.json)
```

---

## Roles

| Role | Responsibility |
|---|---|
| **Benchmark Lead** | Maintain harness, run suite, generate reports |
| **QA** | Manual quality spot-checks, edge cases |
| **DevOps** | CI/CD setup, hardware provisioning, secrets |
| **Platform Lead** | Interpret results, prioritize backends, pairing recommendations |

---

## Known Issues & Workarounds

| Issue | Workaround |
|---|---|
| ONNX BGE-large deadlock | Use bge-base-en-v1.5 instead |
| Ollama model cold start | Measure warm latency separately |
| Cloud API rate limits | Add exponential backoff to harness |
| LLM non-determinism | Set temperature=0 for benchmarks |

---

## Next Steps for v0.0.4

1. [ ] Create `.hawp/benchmarks/test-data/` with semantic pairs + passages
2. [ ] Scaffold `librarian/go/benchmarks/harness.go`
3. [ ] Add ONNX LLM benchmarks (when models available)
4. [ ] Run regression tests vs v0.0.3
5. [ ] Update BACKENDS.md with results

---

**See full plan for details:** [`BENCHMARK_PLAN_v004plus.md`](BENCHMARK_PLAN_v004plus.md)
