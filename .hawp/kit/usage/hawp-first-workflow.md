# Hawp-First Session Workflow

Use `hawp_search` (MCP) as the default context strategy for kit and work files. Direct reads
load whole files and surface irrelevant content; `hawp_search` returns ranked chunks sized to
your token budget.

---

## Why hawp-first

- `hawp_search` ranks chunks by relevance and trims to `--max-tokens N`; a direct read of a
  large kit file can cost 3-10x more tokens with most of that content unused.
- MCP exposes the same index to any agent without requiring file-system access.
- Ranked results separate signal from noise; direct reads require you to skim the noise yourself.

---

## Session-start pattern

Run once per session if kit or work files changed since the last index:

```bash
hawp search index
hawp search embed --backend ollama   # or --backend onnx for offline
```

Skip if the index is current and no kit/work files changed this session.

---

## During-work pattern

Before reading any `.hawp/kit/` or `.hawp/work/` file, search first:

```
hawp_search  query="<topic>"  context=true  max_tokens=2000
```

Use the returned chunks to answer the question. Fall through to a direct read only if the
chunks are insufficient (see below).

Adjust `max_tokens` to your remaining budget. For quick lookups, 500-1000 is usually enough.
For broad context gathering, use 2000-4000.

---

## When to fall through to direct reads

Direct reads are appropriate when:

- The target is **implementation code** (Go source, TypeScript scripts, CI configs) — not kit
  or work docs.
- You already know the **exact file and section** and need the full text verbatim.
- The search index is stale or unavailable and re-indexing is not practical mid-session.
- A single small file (under ~100 lines) where reading costs less than formulating a query.

When falling through, prefer `Read` with `offset` and `limit` to read only the relevant section.

---

## Re-index trigger

Re-run `hawp search index && hawp search embed --backend ollama` any time:

- Kit files were edited this session.
- Work files (backlog, active plans, evidence) were added or modified.
- A new status report or evidence file was committed.

The index is fast (seconds for typical repo sizes); err on the side of re-indexing.

---

## Quick reference

| Situation | Action |
|-----------|--------|
| Session start, files changed | `hawp search index && hawp search embed --backend ollama` |
| Need context on a topic | `hawp_search` via MCP, `max_tokens` to budget |
| Broad context pull | `hawp_search context=true max_tokens=3000` |
| Implementation code | Direct `Read` — search index does not cover non-kit files |
| Exact known section | Direct `Read` with `offset`/`limit` |
| Index stale mid-session | Re-run index + embed, then search |
