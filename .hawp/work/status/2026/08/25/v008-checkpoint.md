# v0.0.8 Complete — Checkpoint 2026-08-25

**Branch:** `feature/v0.0.8` (18 commits ahead of main)  
**Status:** ALL WORK DONE — ready to ship  
**Validation:** 3/3 pass (0 issues, 1 warning)

## What Changed

### Shipped in this branch (vs v0.0.7 baseline)

| Area | Change |
|------|--------|
| Token reduction | Replaced silent no-op `DeduplicateResults` with `ContentJaccardDedup` (Jaccard > 0.70); ~30% savings on duplicate-heavy queries |
| Dynamic chunk cap | Greedy selection stops when next chunk exceeds `--max-tokens` |
| `--hybrid-ratio <f>` | Tune lexical/semantic blend: 0.3 (default) → 0.5 → 0.7; out-of-range exits 1 |
| `--verbose \| -v` | Prints `context: N chunks, ~M tokens (saved ~K tokens via dedup)` to stderr |
| hawp-first-workflow.md | New kit doc: MCP-search-first session pattern, worktree cleanup, token budget lesson |
| Docs cleanup | `context-reshaping.md`, `troubleshooting.md` deleted (entirely about removed `--llm-reshape`); 3 stale Go doc comments cleaned |
| Backlog audit | Recently Closed 15 → 10; `v010-3-3a` reason corrected; all plan links verified |
| benchmarks-v008.md | ONNX vs Ollama comparative benchmark; hybrid ratio sensitivity (0.3/0.5/0.7 all ~66-88ms) |

### Git log summary (newest first)
```
b618cb5  chore: add missing Outcome/Verification/Close Checklist to c98518bb and 1c743447
5cde028  chore: add Outcome/Verification/Close Checklist to 3 closed plans (work validate fix)
b92651b  feat: benchmarks-v008 + search.md update for --hybrid-ratio and --verbose
c8f8648  chore: close 1c743447 in BACKLOG
0485978  feat: 1c743447 — configurable --hybrid-ratio flag (squash)
c663bde  chore: session-close snapshot for v0.0.8 handoff
e37ff98  chore: close 0ca7cf49 in BACKLOG
c07fd65  chore: 0ca7cf49 — stale llm-reshape doc traces (squash)
787ae7b  chore: 8672216a — backlog audit (squash)
5a5d827  chore: CHANGELOG v0.0.8 full; close 4c88f451; token-budget lesson
b0991b8  feat: 4c88f451 — Jaccard dedup + dynamic chunk cap (squash)
9a51a42  chore: add work items 1c743447, 0ca7cf49, 8672216a
dec41c8  chore: worktree-cleanup lesson + close c98518bb
f75a93f  feat: c98518bb — hawp-first session workflow kit doc (squash)
12bcba8  chore: open active plans c98518bb and 4c88f451
9617e5d  docs: v0.0.8 cleanup — remove llm-reshape docs, fix CHANGELOG
b489a22  chore: remove --llm-reshape from CLI; engine key canonical
12fdb41  fix: v0.0.8 — engine as canonical JSON key for context config
```

## Closed Work Items

| ID | Title | Branch |
|----|-------|--------|
| c98518bb | hawp-first session workflow kit doc | feature/v008-hawp-workflow-docs |
| 4c88f451 | Token reduction: Jaccard dedup + dynamic chunk cap | feature/v008-token-reduction |
| 1c743447 | `--hybrid-ratio` flag | feature/v008-hybrid-ratio |
| 0ca7cf49 | Stale llm-reshape doc traces cleanup | feature/v008-llm-reshape-cleanup |
| 8672216a | Backlog audit: compact Recently Closed | feature/v008-backlog-audit |
| benchmarks | ONNX vs Ollama benchmark | inline on manager branch |

## Session Lessons (now in memory + kit docs)

1. **Worktree cleanup** — `git worktree remove --force` immediately after squash-merge; auto-cleanup only fires when agent made no changes.
2. **Token budget** — batch all parallel work before context limit; set ScheduleWakeup; commit to manager branch; save memory first.
3. **Flat branch naming** — `feature/v008-*` not `feature/v0.0.8/*`; git prevents nested refs when parent branch exists.

## Next Action

Ship v0.0.8: open PRs `feature/v0.0.8` → `development` → `main`.  
Release pipeline fires automatically: `tag-on-merge.yml` reads `version.go = "0.0.8"`, triggers `release-librarian-go.yml` (6-platform binaries + kit bundle).

Optionally before shipping: re-index + re-embed to pick up `librarian/docs/benchmarks-v008.md` in MCP search.

## v0.1.0 Gate Status

Partial: Jaccard dedup delivers ~30% savings on duplicate-heavy queries.  
Remaining: smart context sizing + dynamic budget-aware chunk cap demonstrably beating the pre-0.0.8 baseline.
