# hawp search — Context-Aware Search

`hawp search` finds relevant content from your indexed kit and work documents using lexical, semantic, or hybrid search. It can also pack results into LLM-ready context blocks.

## Quick start (CLI)

```bash
# 1. Index documents (run once, re-run after content changes)
hawp search index

# 2. Embed documents for semantic/hybrid search (--backend is required)
hawp search embed --backend ollama   # recommended: fast + high quality
hawp search embed --backend onnx     # offline fallback, no Ollama needed

# 3. Search
hawp search "how do I write a status report"
hawp search "backlog alignment rules"          # hybrid when vectors present
hawp search "HAWP shape mission constraints" --semantic  # pure-vector
```

## Embedding backends

| Backend | Model (default) | First embed | Warm embed (batch) | Quality |
|---------|----------------|-------------|-------------------|---------|
| `ollama` | `nomic-embed-text` (768d) | ~110ms/chunk | ~32ms/chunk | 0.83 |
| `onnx` | `all-MiniLM-L6-v2` (384d) | ~476ms/batch | ~8ms/chunk | 0.70 |

`--backend` is required. Ollama must be running for the `ollama` backend. ONNX runs offline (pure Go, no system deps).

```bash
# Override model
hawp search embed --backend ollama --model nomic-embed-text
hawp search embed --backend onnx   --model all-MiniLM-L6-v2
```

## Search modes

| Mode | Flag | Speed | Notes |
|------|------|-------|-------|
| Lexical (FTS5) | default | <1ms | No vectors needed |
| Semantic (pure vector) | `--semantic` | ~480ms | Requires embed step; no FTS5 |
| Hybrid (lexical + semantic) | auto / `--hybrid-ratio` | ~70ms | Requires embed step; ratio tunable |

Hybrid is the automatic default once vectors are built — `hawp search` upgrades to hybrid re-ranking when a vector index exists, no flag required. Use `--semantic` when a query has few keyword matches and you want full corpus coverage ranked by concept similarity.

## Context packing

```bash
hawp search "intake workflow" --context
hawp search "intake workflow" --context --max-tokens 4000 --format json
```

Output includes source file, chunk, and relevance score per result.

## Full flag reference

```
hawp search <query>
  --limit <n>             max results (default: 10)
  --semantic              pure-vector search — ranks all stored vectors by cosine similarity
                          (no FTS5; requires embed step; ~480ms warm Ollama)
  --hybrid-ratio <f>      lexical fraction for hybrid blend, range [0.0, 1.0] (default: 0.3)
                          0.3 = semantic-dominant (default); 0.5 = equal; 0.7 = lexical-heavy
  --context               pack results into LLM-ready context block
  --format markdown|json  output format (default: markdown)
  --max-tokens <n>        token budget for context block (default: 2000)
  --verbose | -v          print token accounting to stderr: chunks, ~tokens, saved via dedup

hawp search index         ingest configured paths into SQLite (reads .hawp/config/search.json)
hawp search embed         generate and store embedding vectors
  --backend onnx|ollama   required: embedding backend to use
  --model <name>          override default model for the backend
hawp search benchmark     3-way speed + quality benchmark (lexical / semantic / hybrid)
```

## Configuring index paths

By default `hawp search index` ingests `.hawp/kit/` and `.hawp/work/`. To include additional directories or files, create `.hawp/config/search.json` in your repo:

```json
{
  "index": {
    "paths": [
      ".hawp/kit",
      ".hawp/work",
      "librarian/docs",
      "README.md"
    ]
  }
}
```

Paths are relative to the repo root. Directories are walked recursively for `.md` files; individual files are indexed directly. Missing paths print a warning and are skipped. A home-level default can be set at `~/.hawp/config/search.json` — the project config overrides it.

## MCP tool: hawp_search

When using Claude Code with the HAWP MCP server (`hawp mcp`), the `hawp_search` tool provides structured results with precise line positions and context windows — suitable for automated code navigation and documentation lookup.

Run `hawp init --provider <name>` to write the correct config file for your agent:

| Agent | Config file written | Command |
|-------|--------------------|---------| 
| Claude Code | `.mcp.json` | `hawp init --provider claude` |
| Cursor | `.cursor/mcp.json` | `hawp init --provider cursor` |
| Codex | `.codex/config.toml` | `hawp init --provider codex` |

**Codex: project trust is required.** Codex only loads project-scoped MCP config
(`.codex/config.toml`) for projects it considers trusted. If `codex mcp list` does not
show `hawp` after writing the config, trust the project in Codex settings and start a
fresh task or session. The desktop UI does not hot-reload MCP changes mid-session.
Verify with the CLI: `codex mcp list` and `codex mcp get hawp`.

For Claude Code, the config is:

```json
{
  "mcpServers": {
    "hawp": { "command": ".hawp/bin/hawp", "args": ["mcp"] }
  }
}
```

### Raw results mode (default)

```json
{ "query": "backlog alignment", "limit": 5 }
```

Returns structured chunks with relevance scores and line positions:

```json
{
  "query": "backlog alignment",
  "results": [
    {
      "source": ".hawp/kit/references/backlog-alignment.md",
      "relevance": 0.91,
      "content": "Keep Active Work short and current…",
      "lines": {
        "range": { "start": 1, "end": 18 },
        "source": 3
      },
      "context": {
        "window": { "start": 1, "end": 58 }
      }
    }
  ]
}
```

- `lines.range` — full chunk span in the source file (1-indexed)
- `lines.source` — line containing the best match for the query term
- `context.window` — suggested read window (±40 lines around `lines.source`, clamped to file bounds)

### Context mode (pre-shaped for LLM injection)

Pass `context: true` to apply the same Jaccard dedup + greedy token cap + markdown formatting pipeline as the CLI `--context` flag. The result is a single pre-shaped block ready to inject directly into an LLM prompt — smaller and already deduplicated:

```json
{ "query": "backlog alignment", "limit": 5, "context": true, "max_tokens": 2000 }
```

Response:

```json
{
  "query": "backlog alignment",
  "content": "# Search Results: \"backlog alignment\"\n\n**Results:** 3 | **Tokens:** 847/2000\n\n…",
  "token_count": 847,
  "budget": 2000,
  "chunks_used": 3,
  "chunks_dropped": 1
}
```

- `content` — markdown context block, ready for LLM injection
- `token_count` — estimated tokens in the block (chars/4)
- `chunks_dropped` — chunks removed by Jaccard dedup (>70% word-set overlap)

`max_tokens` defaults to 2000. Use context mode when you want the search result to slot directly into a system prompt or retrieval step without post-processing.

Other MCP tools: `hawp_work_new` (create work item), `hawp_work_validate` (validate kit + work integrity).

## Typical agent workflow

```bash
# At session start (once per session, if content changed)
hawp search index && hawp search embed --backend ollama

# During work
hawp search "backlog alignment rules" --hybrid --context --max-tokens 6000
```

The index is idempotent — re-running after content changes upserts new chunks. Vectors persist across sessions.
