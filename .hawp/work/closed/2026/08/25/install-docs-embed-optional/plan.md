# install-docs-embed-optional — Clarify: index first, embed is optional and slow

**Type:** docs  
**Status:** done  
**Opened:** 2026-08-25  
**Closed:** 2026-08-25
**Target:** v0.0.11

**Closes:** `install-docs-embed-optional`

## Goal

After this fix, install docs and `usage/search.md` give clear, accurate guidance:
- Step sequence: index → (optional embed) → search
- Lexical is the default; embed is optional
- Embed timing is hardware-dependent; do not block on it
- Local Ollama ≠ automatic Cursor token savings

## Outcome

The search guide now states that `hawp search index` is the required first step
for CLI and MCP search, and that `hawp search embed` is optional background prep
used only for semantic or hybrid retrieval. The timing table is now qualified as
machine-specific, `--max-tokens` / `--verbose` are described as context-packing
estimates rather than billing meters, and the Ollama note is explicit that local
embedding/search execution does not automatically change an editor agent's cloud
model routing or token spend.

The shared install guide now also warns that install does not build the search
index for the user, so immediate `hawp_search` expectations are anchored to a
post-install `hawp search index` step.

## Verification

- [x] Updated `.hawp/kit/usage/search.md` quick-start, MCP, context-budget, and workflow sections. Evidence: this plan's What was done section lists those guide changes.
- [x] Updated `distribution/sources/shared/install.md` post-install expectations. Evidence: this plan's What was done section lists the shared install-guide update.
- [x] `mise exec node@26.5.0 -- npm run distribution:sync`. Evidence: the command is recorded in this plan's Verification section.
- [x] Verified generated install guides were refreshed for all providers' install paths. Evidence: this plan's Outcome section states the generated guides were refreshed.

## What was done

- Reframed quick start as `index` required, `embed` optional, lexical search immediate after indexing
- Added MCP note that `hawp_search` needs indexing but not embeddings for lexical use
- Clarified hardware-dependent embed timings and that local Ollama does not automatically reduce editor-agent cloud-token spend
- Clarified `--max-tokens` / `--verbose` as context-packing estimates
- Added post-install shared guide note so generated install docs set the same expectation

## Close Checklist

- [x] Outcome recorded
- [x] Verification includes source and generated install-guide evidence
- [x] Lexical-first guidance is captured explicitly
- [x] Ready to stay in closed history
