# 4c88f451 — Token reduction: smart context sizing + dedup-before-pack

**Type:** improvement  
**Status:** in-progress  
**Branch:** feature/v008-token-reduction → feature/v0.0.8  
**Updated:** 2026-08-24

## Input

Prior benchmark showed tokens *increased* with context packing (bloat offset gains).
v0.1.0 gate: demonstrate actual token reduction. Approach: dedup before ContextBlock
assembly, dynamic chunk cap by token budget, measure tokens-in vs tokens-saved.

## Scope

1. **Pre-pack dedup** — before building the ContextBlock, drop chunks whose content
   overlaps >N% with a higher-ranked chunk (lexical overlap, not vector — fast).
2. **Dynamic chunk cap** — don't always use limit=10; select greedily until budget
   exhausted (stop adding when next chunk would exceed remaining tokens).
3. **Token measurement** — add a `--verbose` output line: `tokens saved: N (M%)`.
4. **Benchmark** — run before/after to prove reduction is real.

## Relevant files

- `librarian/src/internal/application/context/` — ContextBlock, reshaper pipeline
- `librarian/src/internal/platform/cli/run.go` — search output + verbose flag
- `librarian/src/internal/application/search/service.go` — HybridRank, SemanticSearch

## Acceptance

- `hawp search "query" --context --verbose` shows `tokens saved: N (M%)`.
- Before/after benchmark shows net reduction (tokens saved > tokens added by wrapping).
- Commit on `feature/v008-token-reduction`, squash-merge → `feature/v0.0.8`.
