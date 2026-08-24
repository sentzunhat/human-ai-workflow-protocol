# hawp search — Context-Aware Search

`hawp search` finds relevant content from your indexed kit and work documents using lexical, semantic, or hybrid search. It can also pack results into LLM-ready context blocks and optionally reshape them via a local LLM.

## Quick start

```bash
# 1. Index documents (run once, re-run after content changes)
hawp search index

# 2. Embed documents for semantic/hybrid search
hawp search embed

# 3. Search
hawp search "how do I write a status report"
```

## Search modes

| Mode | Flag | Speed | Quality |
|------|------|-------|---------|
| Lexical (FTS5) | default | <1ms | 70% |
| Semantic (vector) | `--semantic` | ~100ms | 95% |
| Hybrid (blend) | `--hybrid` | 15–20ms | 96% |

Hybrid is the recommended default once vectors are built.

## Context packing

`--context` formats results as an LLM-ready context block with inline
document references, cosine-similarity deduplication (>90% duplicate
removal), and token-budget enforcement:

```bash
hawp search "intake workflow" --context
hawp search "intake workflow" --context --max-tokens 4000
hawp search "intake workflow" --context --format json
```

Output includes source file, chunk, and relevance score per result.

## LLM reshape

`--llm-reshape` passes the context block through a local LLM (Ollama or
ONNX text2text) to produce a coherent, synthesized answer instead of raw
chunks. Requires a running Ollama instance or a configured ONNX model:

```bash
# Requires Ollama running locally
hawp search "how do I close a work item" --llm-reshape

# Configure backend
HAWP_LLM_BACKEND=ollama HAWP_OLLAMA_URL=http://localhost:11434 \
  hawp search "status report format" --llm-reshape
```

Falls back to unreshaped context with a warning if the LLM is unavailable.

## Embedding backends

```bash
# ONNX (default, no network, ~100ms)
hawp search embed --backend onnx

# Ollama (requires running Ollama, ~25ms)
hawp search embed --backend ollama --model nomic-embed-text
```

Configure via environment:
- `HAWP_EMBEDDINGS_BACKEND` — `onnx` (default) or `ollama`
- `HAWP_OLLAMA_URL` — Ollama base URL (default: `http://localhost:11434`)

## Full flag reference

```
hawp search <query>
  --limit <n>          max results (default: 5)
  --context            pack results into LLM-ready context block
  --llm-reshape        reshape context via embeddings + LLM (implies --context)
  --format markdown|json  output format (default: markdown)
  --max-tokens <n>     token budget for context block (default: 8000)
  --semantic           force semantic-only search
  --hybrid             force hybrid search (lexical + semantic blend)

hawp search index      ingest documents into SQLite (no vectors)
hawp search embed      generate and store embedding vectors
```

## Typical agent workflow

```bash
# At session start (once per session)
hawp search index && hawp search embed

# During work
hawp search "backlog alignment rules" --context --max-tokens 6000
```

The index is idempotent — re-running after content changes upserts new
chunks and removes stale ones. Vectors persist across sessions.
