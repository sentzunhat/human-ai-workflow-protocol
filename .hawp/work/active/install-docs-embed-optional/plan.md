# install-docs-embed-optional — Clarify: index first, embed is optional and slow

**Type:** docs  
**Status:** plan-ready  
**Opened:** 2026-08-25  
**Target:** v0.0.11

## Input

Downstream install of hawp 0.0.9. Install docs and kit `usage/search.md` did
not clearly communicate:

1. `hawp search index` must be run before `hawp_search` returns results.
2. Lexical MCP search works without embeddings — `hawp search embed` is not required.
3. Embed can be extremely slow on CPU-only machines (~1 chunk/s measured; kit
   table claims ~110ms/chunk which is hardware-dependent).
4. Ollama local models save $0 API tokens only when the agent is routed to Ollama;
   installing Ollama does not automatically cut Cursor cloud token spend.
5. `--max-tokens` / `--verbose` counts packed context size (chars/4 estimate),
   not Cursor subscription billing.

Evidence source: downstream install evidence 2026-08-25.

## Goal

After this fix, install docs and `usage/search.md` give clear, accurate guidance:
- Step sequence: index → (optional embed) → search
- Lexical is the default; embed is optional
- Embed timing is hardware-dependent; do not block on it
- Local Ollama ≠ automatic Cursor token savings

## Constraints

- Do not change CLI behavior — docs only.
- Keep it concise; don't pad with disclaimers.
- Separate direct evidence from inference in any timing claims.

## Plan

### Step 1 — Audit `usage/search.md` and install guide

Identify every place that:
- Suggests or implies embed is required for MCP search
- Shows timing numbers without a hardware qualifier
- Conflates local Ollama with Cursor token savings
- Omits `hawp search index` as a prerequisite step

### Step 2 — Update `usage/search.md`

Add or rewrite the "getting started" sequence:

```
1. hawp search index          # required — builds FTS5 lexical index
2. hawp search embed ...      # optional — needed only for --semantic or --hybrid
3. hawp_search (MCP)          # works after step 1, no embed needed for lexical
```

Add a warning box or callout:

> **Embed timing is hardware-dependent.** The embedded timing in this doc was
> measured on warm GPU. CPU-only machines may take ~1s/chunk (hours for a full
> corpus). Run embed in a background session and use lexical search in the meantime.

Add a note on Ollama and tokens:

> **Ollama saves API tokens only when the agent is routed to Ollama.** Cursor
> agents still use the cloud model unless you explicitly route them to a local
> Ollama endpoint. `--max-tokens` is a context-packing budget (chars/4 estimate),
> not a Cursor billing meter.

### Step 3 — Review install guide for the same gaps

Same fixes in `distribution/sources/shared/install.md` or any generated
install/quickstart docs.

### Step 4 — Validate with `npm run check:markdown-links`

Run link check to confirm no broken references after edits.

## Verification

- A user following the updated docs can run `hawp_search` after only `hawp search index`
- Embed section is clearly optional with timing caveats
- No misleading Ollama/token claims remain
