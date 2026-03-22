---
work-item: k8i9f5m1-p3-4-5
type: feature
title: "v0.0.3 Phase 3.4-3.5: Context Reshaping Pipeline & Testing"
status: plan-ready
owner: unassigned
created: 2026-07-24
updated: 2026-07-24
---

# Phase 3.4-3.5: Context Reshaping Pipeline & Testing

## Mission

Integrate embeddings (Phase 3.2) and LLM backends (Phase 3.3) into the full context reshaping pipeline, then test and document for production use.

---

## Context

**What we have so far (after Phases 3.1-3.3):**
- ✅ Phase 3.1: Config system (select backend + model)
- ✅ Phase 3.2: Embedding backends (ONNX, Ollama, OpenAI, Anthropic)
- ✅ Phase 3.3: LLM backends (ONNX, Ollama, OpenAI, Anthropic)

**What Phase 3.4 does:**
- Combines embeddings + LLM into one reshaping pipeline
- Identifies key concepts via embeddings
- Sends to LLM for final structuring
- Returns improved context

**What Phase 3.5 does:**
- End-to-end testing of all backend combinations
- Performance benchmarking
- Documentation and usage examples

---

## Phase 3.4: Context Reshaping Pipeline

### Design

```go
// ContextReshaper orchestrates embeddings + LLM for context improvement
type ContextReshaper struct {
    embedder Embedder
    llm      LLMClient
    config   ContextConfig
}

func NewContextReshaper(config ContextConfig) (*ContextReshaper, error) {
    // 1. Initialize embedder from config
    embedder, err := NewEmbedder(config.Embeddings)
    if err != nil {
        return nil, err
    }
    
    // 2. Initialize LLM from config
    llm, err := NewLLMClient(config.LLM)
    if err != nil {
        embedder.Close()
        return nil, err
    }
    
    return &ContextReshaper{
        embedder: embedder,
        llm:      llm,
        config:   config,
    }, nil
}

// Reshape improves a context block via embeddings + LLM
func (r *ContextReshaper) Reshape(ctx context.Context, block *ContextBlock) (*ContextBlock, error) {
    // Step 1: Convert context block to text
    markdown := block.String()
    
    // Step 2: Identify key concepts via embeddings
    // - Extract main concepts from each result
    // - Rank by relevance
    concepts := r.identifyKeyConcepts(ctx, block)
    
    // Step 3: Prepare prompt with key concepts
    prompt := r.buildReshapingPrompt(markdown, concepts)
    
    // Step 4: Send to LLM for restructuring
    reshaped, err := r.llm.Reshape(ctx, prompt, r.config.LLM.MaxTokens)
    if err != nil {
        return nil, err
    }
    
    // Step 5: Parse reshaped output back into ContextBlock
    improved := r.parseReshapedOutput(reshaped, block)
    
    return improved, nil
}

// identifyKeyConcepts uses embeddings to find key themes
func (r *ContextReshaper) identifyKeyConcepts(ctx context.Context, block *ContextBlock) []string {
    // 1. Embed each result's content
    // 2. Find clusters of similar concepts
    // 3. Extract top-N representative concepts
    // 4. Return ranked list
}

// buildReshapingPrompt creates a rich prompt with context + concepts
func (r *ContextReshaper) buildReshapingPrompt(markdown string, concepts []string) string {
    // Include markdown + key concepts for LLM
    return fmt.Sprintf(`
System: You are a technical context reshaping expert.

Key concepts identified: %v

Original context:
%s

Reshaped context:
`, concepts, markdown)
}

// parseReshapedOutput converts LLM output back to ContextBlock
func (r *ContextReshaper) parseReshapedOutput(reshapedText string, original *ContextBlock) *ContextBlock {
    // 1. Parse reshapedText into structured format
    // 2. Preserve original result references
    // 3. Update Content with reshaped text
    // 4. Return new ContextBlock
}

// Close releases resources
func (r *ContextReshaper) Close() error {
    r.embedder.Close()
    r.llm.Close()
    return nil
}
```

### CLI Integration

Add to `hawp search --context`:

```bash
# Use default backends (ONNX embeddings + ONNX LLM)
hawp search "query" --context

# Reshape with OpenAI embeddings + Anthropic LLM
hawp search "query" --context \
    --embeddings-backend openai \
    --llm-backend anthropic

# Save reshaped context to file
hawp search "query" --context --output reshaped.md

# JSON format
hawp search "query" --context --format json --llm-reshape
```

### Acceptance Criteria for Phase 3.4

- [ ] ContextReshaper struct integrates embedder + LLM
- [ ] Key concept identification via embeddings working
- [ ] Prompt engineering tested (quality output)
- [ ] CLI flags for embeddings + LLM backend selection
- [ ] Context block parsing + restructuring working
- [ ] Tests: 10+ (pipeline, concept identification, prompt building, parsing)
- [ ] Documentation: usage examples + performance expectations

### Effort Estimate: Phase 3.4

| Task | Time |
|------|------|
| ContextReshaper struct | 2 hours |
| Key concept identification | 3 hours |
| Prompt engineering + testing | 3 hours |
| CLI integration | 2 hours |
| Tests (10+) | 3 hours |
| Documentation | 1 hour |
| **Total** | **~14 hours** |

---

## Phase 3.5: Testing & Documentation

### Testing Strategy

**Test Matrix: 16 Backend Combinations**

| Embeddings | LLM | Test Case |
|------------|-----|-----------|
| ONNX | ONNX | ✅ Local-only (fastest) |
| ONNX | Ollama | Hybrid (local embed + local LLM) |
| ONNX | OpenAI | Hybrid (local embed + API LLM) |
| ONNX | Anthropic | Hybrid (local embed + API LLM) |
| Ollama | Ollama | ✅ Full Ollama (local) |
| Ollama | ONNX | Hybrid |
| Ollama | OpenAI | Hybrid |
| Ollama | Anthropic | Hybrid |
| OpenAI | ONNX | Hybrid |
| OpenAI | Ollama | Hybrid |
| OpenAI | OpenAI | ✅ Full OpenAI (API-only) |
| OpenAI | Anthropic | Hybrid |
| Anthropic | ONNX | Hybrid |
| Anthropic | Ollama | Hybrid |
| Anthropic | OpenAI | Hybrid |
| Anthropic | Anthropic | ✅ Full Anthropic (API-only) |

**Priority:** Test the 4 primary paths (✅) first, then hybrid combinations.

### Test Implementation

```go
// Test structure
type ReshapingTest struct {
    name           string
    embeddingsBackend string
    llmBackend     string
    config         ContextConfig
    input          *ContextBlock
    expectedOutput *ContextBlock
}

// Run tests
func TestContextReshaping(t *testing.T) {
    cases := []ReshapingTest{
        // ONNX + ONNX (fastest, for CI)
        {name: "onnx-onnx", ...},
        // Ollama + Ollama (if available)
        {name: "ollama-ollama", ...},
        // OpenAI + OpenAI (real API test, optional in CI)
        {name: "openai-openai", ...},
        // Anthropic + Anthropic (real API test, optional in CI)
        {name: "anthropic-anthropic", ...},
        // Hybrids (mock backends)
        {name: "onnx-openai", ...},
        {name: "openai-onnx", ...},
        // ... etc
    }
    
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            r, err := NewContextReshaper(tc.config)
            if err != nil {
                t.Fatal(err)
            }
            defer r.Close()
            
            output, err := r.Reshape(context.Background(), tc.input)
            if err != nil {
                t.Fatal(err)
            }
            
            // Verify output has expected properties
            if output.TokenCount > tc.config.LLM.MaxTokens {
                t.Errorf("output exceeds token budget")
            }
            
            if len(output.Results) == 0 {
                t.Errorf("output has no results")
            }
        })
    }
}
```

### Benchmarking

**Measure for each backend combination:**
- Time to reshape 10 typical queries
- Memory usage peak
- Token consumption (API backends only)
- Cost estimation (API backends)

```go
func BenchmarkContextReshaping(b *testing.B) {
    configs := []struct {
        name   string
        config ContextConfig
    }{
        {name: "onnx-onnx", config: ...},
        {name: "ollama-ollama", config: ...},
        {name: "openai-openai", config: ...},
    }
    
    for _, cfg := range configs {
        b.Run(cfg.name, func(b *testing.B) {
            r, _ := NewContextReshaper(cfg.config)
            defer r.Close()
            
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                r.Reshape(context.Background(), sampleBlock)
            }
        })
    }
}
```

### Documentation

**Sections:**
1. **Quick Start:** ONNX defaults (no config needed)
2. **Configuration:** How to configure embeddings + LLM backends
3. **Usage Examples:** CLI + Go API
4. **Backend Comparison:** Performance/quality/cost table
5. **Troubleshooting:** Common issues per backend
6. **Architecture:** How reshaping works under the hood

**Example docs/context-reshaping.md:**

```markdown
# Context Reshaping in v0.0.3

## Quick Start (Local, No APIs)

```bash
hawp search "your query" --context --llm-reshape
```

Uses default: ONNX embeddings (BGE) + ONNX LLM (TinyLlama)

## Configure Backends

Create ~/.hawp/config/context.json:

```json
{
  "embeddings": {
    "backend": "onnx",
    "model": "bge-base-en-v1.5"
  },
  "llm": {
    "backend": "openai",
    "model": "gpt-4-turbo",
    "temperature": 0.3
  },
  "backends": {
    "openai": {
      "api_key": "[encrypted key]"
    }
  }
}
```

## Backend Comparison

| Backend | Speed | Quality | Cost | Setup |
|---------|-------|---------|------|-------|
| ONNX | ⚡ Fast | ⭐⭐⭐ Good | Free | auto-download models |
| Ollama | ⚡ Fast | ⭐⭐⭐ Good | Free | run `ollama serve` |
| OpenAI | 🐢 Slow | ⭐⭐⭐⭐ Great | $$ | API key |
| Anthropic | 🐢 Slow | ⭐⭐⭐⭐ Great | $$ | API key |
```

### Acceptance Criteria for Phase 3.5

- [ ] 16 backend combinations tested (4 primary required, 12 hybrid optional)
- [ ] Benchmarking: time/memory/tokens for each
- [ ] All tests passing (CI runs ONNX, manual runs for APIs)
- [ ] Documentation complete (5+ sections)
- [ ] Example configs for each backend
- [ ] Troubleshooting guide (common issues)
- [ ] Cost estimation guide (for API users)

### Effort Estimate: Phase 3.5

| Task | Time |
|------|------|
| Integration tests (16 combinations) | 4 hours |
| Benchmarking + data collection | 2 hours |
| Documentation (5+ sections) | 3 hours |
| Examples + troubleshooting | 2 hours |
| Code review + final polish | 1 hour |
| **Total** | **~12 hours** |

---

## Timeline

**Phase 3.4-3.5 Sequential:**
- Phase 3.4 (Pipeline): ~14 hours (1.5-2 days)
- Phase 3.5 (Testing & Docs): ~12 hours (1-1.5 days)
- **Total: ~26 hours (3-3.5 days)**

**Note:** Phases 3.2 and 3.3 run in parallel (2 days total), so while they're finishing, Phase 3.1 can be done first, then 3.4 starts.

---

## Dependencies

- ✅ Phase 3.1: Config system
- ✅ Phase 3.2: Embedding backends
- ✅ Phase 3.3: LLM backends

---

## Success Metrics

✅ **Pipeline works end-to-end:**
- Input: packed context block
- Process: embeddings + key concept identification + LLM reshaping
- Output: improved context block

✅ **All backend combinations functional:**
- ONNX + ONNX (primary, always works)
- Ollama + Ollama (if user runs Ollama)
- OpenAI + OpenAI (if user has API key)
- Anthropic + Anthropic (if user has API key)
- All hybrids work (any mix)

✅ **Well-documented and tested:**
- 30+ tests covering all paths
- Benchmarks showing speed/cost trade-offs
- Clear usage examples
- Troubleshooting guide

---

## v0.0.3 Complete

After Phases 3.1-3.5:
- ✅ Config system (foundation)
- ✅ Embedding backends (ONNX, Ollama, OpenAI, Anthropic)
- ✅ LLM backends (ONNX, Ollama, OpenAI, Anthropic)
- ✅ Context reshaping pipeline (integrated)
- ✅ Testing & documentation (production-ready)

**Users can:**
- Use free local ONNX (default)
- Swap to Ollama (local LLM server)
- Use OpenAI / Anthropic for better quality
- Mix and match embeddings + LLM backends
- Configure via file or env vars
- Securely store API keys

---

**Status:** Ready (after Phases 3.1-3.3 complete).

## Outcome

Shipped in the `0.0.1` release (tag `0.0.1`, 2026-08-21). Work complete.

## Verification

Release published at https://github.com/sentzunhat/human-ai-workflow-protocol/releases/tag/0.0.1 with all 7 assets.

## Close Checklist

- [x] Work shipped in 0.0.1 release
- [x] Archived to closed/2026/08/22/v001-shipped-cleanup/
