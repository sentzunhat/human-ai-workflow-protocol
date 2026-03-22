---
work-item: t5l7m9n1
type: planning
title: "Benchmarking Infrastructure for v0.0.4+ — Comprehensive Provider Matrix"
status: plan-ready
owner: unassigned
created: 2026-07-26
updated: 2026-07-26
---

# Benchmarking Plan: v0.0.4+ Provider Matrix

## Mission

Establish a repeatable, comprehensive benchmarking methodology for all embedding and LLM provider combinations as they ship in v0.0.4+. This plan enables data-driven decisions about provider pairings, quality trade-offs, and cost optimization.

---

## Context

v0.0.3 collected initial performance baselines for ONNX and Ollama backends:
- ONNX embeddings: 104–605ms (varies by model dimension)
- Ollama embeddings: 24–35ms (all models similar, despite dimension)
- Ollama LLM: 12–120s per reshape (CPU-bound)

v0.0.4+ will add incrementally:
- ONNX LLM (text2text models like FLAN-T5)
- OpenAI embeddings + LLM
- Anthropic embeddings + LLM

**Problem:** Without structured benchmarking, we lack:
- Systematic quality measurement (not just latency)
- Cost analysis for cloud APIs
- Pairing recommendations ("which embedding works best with which LLM?")
- Regression detection (did a new provider slow down existing backends?)

---

## Scope

### What This Plan Covers

1. **Metrics definition** — what to measure for embeddings and LLMs
2. **Test data** — reproducible, domain-specific test passages
3. **Matrix testing** — all combinations of backends
4. **Quality measurement** — objective metrics: semantic coherence, fact retention, context relevance
5. **Cost tracking** — per-API cost accounting
6. **CI/CD integration** — daily/weekly benchmark runs, regression detection
7. **Reporting** — markdown summaries, heatmaps, recommendations

### What This Plan Doesn't Cover

- **GPU benchmarking** — document as future work (implies stable hardware provisioning)
- **Distributed/batch processing optimization** — future v0.2.0
- **Human evaluation** — future if quality metrics plateau
- **Streaming output** — future if needed

---

## Key Decisions (Locked)

1. **Test hardware:** Stable Mac (M1 Max) — GPU runs separate from CPU baseline
2. **Quality = objective only:** Automated checks (coherence, fact retention), no manual scoring initially
3. **Storage:** `.hawp/benchmarks/results/YYYY-MM-DD_version/` with detailed JSON
4. **CI cadence:** Daily for local backends (ONNX, Ollama), weekly for cloud APIs (cost control)
5. **Pairing strategy:** Test self-pairing + cross-pairing with Ollama to isolate provider impact

---

## Deliverables

| Item | Location | Status |
|---|---|---|
| Complete benchmark plan | `librarian/docs/BENCHMARK_PLAN_v004plus.md` | ✅ Done |
| Test data repository | `.hawp/benchmarks/test-data/` | Stub only |
| Benchmark harness (`go test` suite) | `librarian/src/benchmarks/` | Stub only |
| CI/CD workflow | `.github/workflows/benchmarks.yml` | Not started |
| First v0.1.0 run (OpenAI added) | `.hawp/benchmarks/results/2026-09-15_v010/` | Future |

---

## Quality Thresholds (Success Criteria)

### Embedding Model Acceptance

- Single embed latency < 1000ms (local + cloud)
- Semantic coherence >= 0.7 for related terms
- Embedding stability > 0.99 (deterministic)
- No deadlocks/crashes on test data

### LLM Model Acceptance

- 100-token reshape < 30 seconds (user-acceptable)
- Output coherence >= 4/5 (readable, no token corruption)
- Fact retention >= 80% (preserves original)
- Context relevance >= 0.6 (stays semantically close)

### Combination Acceptance

- Both backends pass independently
- Combined latency < 60 seconds
- Quality metrics both pass

---

## Implementation Timeline

### Phase 1: Infrastructure (Before v0.0.4)
- [x] Write comprehensive benchmark plan (this document)
- [ ] Create test data repository (semantic pairs, domain passages)
- [ ] Scaffold `librarian/src/benchmarks/` harness structure
- [ ] Add to `package.json` scripts for easy invocation

### Phase 2: v0.0.4 Cycle (During v0.0.4 development)
- [ ] Implement basic harness (single backend runner, metric collection)
- [ ] Add ONNX LLM benchmarks when models available
- [ ] Run regression tests vs v0.0.3 baseline
- [ ] Update BACKENDS.md with new metrics

### Phase 3: v0.1.0 Preparation (Before v0.1.0)
- [ ] Set up CI/CD workflow (`.github/workflows/benchmarks.yml`)
- [ ] Add cloud API support to harness (OpenAI, Anthropic)
- [ ] Implement cost tracking module
- [ ] Full matrix benchmark (all combinations)
- [ ] Generate comprehensive report

### Phase 4: Continuous (v0.1.0+)
- [ ] Daily regression checks (ONNX, Ollama)
- [ ] Weekly cloud API benchmarks (cost control)
- [ ] Monthly reports with recommendations
- [ ] Investigate quality regressions or anomalies

---

## Known Dependencies & Risks

| Dependency | Status | Risk Level |
|---|---|---|
| ONNX LLM models available | Waiting for T5/BART/FLAN models | Medium — may not ship in v0.0.4 |
| Stable CI hardware | Need M1 Mac, not random cloud runners | Medium — currently using dev machine |
| OpenAI/Anthropic API keys | Need for CI environment secrets | Low — standard practice |
| Test data quality | Semantic pairs + passages must be curated | Low — can iterate |

---

## Next Steps

1. **Immediate (before v0.0.4):**
   - Create `.hawp/benchmarks/test-data/` structure
   - Populate with semantic pairs and domain passages
   - Scaffold `librarian/src/benchmarks/harness.go`

2. **During v0.0.4:**
   - Implement basic latency + throughput benchmarks
   - Add ONNX LLM tests when models available
   - Run regression tests vs v0.0.3

3. **Before v0.1.0:**
   - Full CI/CD setup
   - Cloud API integration
   - Complete matrix run
   - Generate final report

---

## Reference

**Full Plan:** `librarian/docs/BENCHMARK_PLAN_v004plus.md`  
**v0.0.3 Baseline:** `librarian/docs/BENCHMARKS_v003.md`  
**Backend Status:** `librarian/docs/BACKENDS.md`
