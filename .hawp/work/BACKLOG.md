# Backlog

Active coordination index for open work. Closed history is archived under `.hawp/work/closed/`.

---

## Status Key

| Status        | Meaning                             |
| ------------- | ----------------------------------- |
| `inbox`       | Received, not yet analyzed          |
| `analyzing`   | Under investigation                 |
| `plan-ready`  | Plan written, awaiting review       |
| `approved`    | Plan approved, ready to implement   |
| `in-progress` | Being implemented                   |
| `parked`      | Deferred without closing            |
| `done`        | Implemented and verified            |
| `blocked`     | Blocked — reason noted in plan file |
| `wont-fix`    | Decided not to fix — reason noted   |

---

## Active Work

All v0.1.0 gates passed; feature branch
`feature/v0.0.24-work-folder-normalization` is near PR-ready after the selected
`v0.0.24` MCP intake cutoff is verified.

Release-scope note:
[`v0.0.24` cutoff and `v0.0.25` MCP lane](status/2026/09/13/6059dd4b/status.md).

Recommended `v0.0.24` cutoff: finish verifying `a8797f44`, including the small
structured-result subset of `c8d9e1cb` that stayed coupled to
`hawp_work_intake`.

Recommended `v0.0.25` lane: `2eea565c`, `429e075e`, and the remaining
elicitation/guidance work from `c8d9e1cb`.

| UUID | Type | Title | Status | Owner | Plan File | Updated |
| ---- | ---- | ----- | ------ | ----- | --------- | ------- |
| `c8d9e1cb` | feature | Interactive MCP intake refinement with structured results (`v0.0.24` subset, remainder `v0.0.25`) | `plan-ready` | Codex | [plan](active/c8d9e1cb/plan.md) | 2026-09-13 |
| `2eea565c` | improvement | MCP capability catalog for HAWP resources and prompts (`v0.0.25`) | `plan-ready` | Codex | [plan](active/2eea565c/plan.md) | 2026-09-13 |
| `429e075e` | improvement | Remote MCP transport and authorization readiness audit (`v0.0.25`) | `plan-ready` | Codex | [plan](active/429e075e/plan.md) | 2026-09-13 |
| `0c8d6391` | bug | Constrain project search index paths to the repository root | `plan-ready` | Codex | [plan](active/0c8d6391/plan.md) | 2026-09-20 |
| `225708ae` | bug | Reject external Markdown symlinks during corpus indexing and export | `plan-ready` | Codex | [plan](active/225708ae/plan.md) | 2026-09-20 |
| `cd4691d7` | bug | Reject symlink ancestors in kitsync repository writes | `plan-ready` | Codex | [plan](active/cd4691d7/plan.md) | 2026-09-20 |
| `a1014eef` | bug | Reject path traversal in install/update backlog reconciliation | `plan-ready` | Claude | [plan](active/a1014eef/plan.md) | 2026-09-20 |
| `a406816d` | bug | Reject symlinked `.hawp` root in install/update shell scripts | `plan-ready` | Claude | [plan](active/a406816d/plan.md) | 2026-09-20 |
| `be327af5` | bug | Reject symlinked `.gitignore` in MCP provider configuration | `plan-ready` | Claude | [plan](active/be327af5/plan.md) | 2026-09-20 |
| `984f038a` | bug | Reject Markdown symlinks in `links clean --apply` | `plan-ready` | Claude | [plan](active/984f038a/plan.md) | 2026-09-20 |

## Blocked / Parked

| ID                           | Type        | Title                                                                  | Reason                                                                                     | Detail                                            | Updated    |
| ---------------------------- | ----------- | ---------------------------------------------------------------------- | ------------------------------------------------------------------------------------------ | ------------------------------------------------- | ---------- |
| `89cf7a85`                   | improvement | v0.0.23 legacy work-item UUID canonicalization follow-up               | Older-repo safety and archive-policy follow-up; not a live `v0.0.23` blocker               | [plan](parked/89cf7a85/plan.md)                   | 2026-08-31 |
| `release-benchmark-backfill` | planning    | Benchmark evidence gate for `v0.0.11` through `v0.0.13`                | Older patch-train backfill, intentionally kept out of the `v0.0.23` release lane           | [plan](parked/release-benchmark-backfill/plan.md) | 2026-08-31 |
| `usage-tracking`             | feature     | `hawp usage` CLI + `hawp_usage` MCP tool for token/wall-clock metering | Scope too large for v0.0.11; belongs at v0.1.0 token-reduction gate                        | [plan](parked/usage-tracking/plan.md)             | 2026-08-25 |
| `bee15107`                   | improvement | Defer CLI participant adapters for Codex, Claude, and GitHub           | Not needed yet; provider packs are enough                                                  | [plan](parked/bee15107/plan.md)                   | 2026-07-06 |
| `v010-3-3a`                  | feature     | ONNX LLM Text2Text models (FLAN-T5-small)                              | Plan ready (FLAN-T5-small feasible via ONNX); deferred per incremental-only patch strategy | [detail](parked/v010-3-3a/plan.md)                | 2026-08-24 |
| `v010-3-2c`                  | feature     | OpenAI Embeddings (text-embedding-3-small/large)                       | Cloud-API tracks parked; see incremental-only strategy                                     | [detail](parked/v010-3-2c/plan.md)                | 2026-07-27 |
| `v010-3-3c`                  | feature     | OpenAI LLM (gpt-3.5-turbo, gpt-4-turbo)                                | Cloud-API tracks parked; see incremental-only strategy                                     | [detail](parked/v010-3-3c/plan.md)                | 2026-07-27 |
| `v010-3-2d`                  | feature     | Anthropic Embeddings Stub                                              | Cloud-API tracks parked; see incremental-only strategy                                     | [detail](parked/v010-3-2d/plan.md)                | 2026-07-27 |
| `v010-3-3d`                  | feature     | Anthropic LLM (claude-3-sonnet, claude-3-opus)                         | Cloud-API tracks parked; see incremental-only strategy                                     | [detail](parked/v010-3-3d/plan.md)                | 2026-07-27 |
| `v010-cost`                  | feature     | Cost Tracking + Rate Limiting (cross-cutting for cloud backends)       | Cloud-API tracks parked; see incremental-only strategy                                     | [detail](parked/v010-cost/plan.md)                | 2026-07-27 |

---

## Recently Closed

Limited to the last 5–10 items.

| ID                            | Type        | Title                                                                   | Closed     | Detail                                                        |
| ----------------------------- | ----------- | ----------------------------------------------------------------------- | ---------- | ------------------------------------------------------------- |
| `887ac089`                    | bug         | Fix model pull trailing `--onnx-file` parsing                          | 2026-09-20 | [plan](closed/2026/09/20/887ac089/plan.md)                   |
| `c5a48d9a`                    | bug         | Harden tar kit-bundle extraction against path traversal                 | 2026-09-20 | [plan](closed/2026/09/20/c5a48d9a/plan.md)                   |
| `b30f7a9b`                    | bug         | Validate kitsync manifest paths stay within repository and bundle roots  | 2026-09-20 | [plan](closed/2026/09/20/b30f7a9b/plan.md)                   |
| `24edc80c`                    | fix         | Wrap SQLite "no such table" as IndexNotFoundError in search infra       | 2026-09-16 | [plan](closed/2026/09/16/24edc80c/plan.md)                    |
| `a818b914`                    | fix         | Fix non-atomic writes in ApplyDuplicateLinks plan files                 | 2026-09-16 | [plan](closed/2026/09/16/a818b914/plan.md)                    |
| `f6f2818f`                    | fix         | Document TOCTOU in writeMCPJSON between Lstat and WriteFile             | 2026-09-16 | [plan](closed/2026/09/16/f6f2818f/plan.md)                    |
| `e61573f6`                    | fix         | Remove double expandConfigProviders call in Configure()                 | 2026-09-16 | [plan](closed/2026/09/16/e61573f6/plan.md)                    |
| `8c695bc5`                    | fix         | Improve WriteProviderConfigs partial-failure error message              | 2026-09-16 | [plan](closed/2026/09/16/8c695bc5/plan.md)                    |
| `4d6d8d2e`                    | fix         | Deduplicate idSet in duplicate_links.go                                 | 2026-09-16 | [plan](closed/2026/09/16/4d6d8d2e/plan.md)                    |
| `cd5c215b`                    | fix         | **[High]** doc: quote YAML front-matter title (newline injection)       | 2026-09-15 | [plan](closed/2026/09/15/cd5c215b/plan.md)                    |
| `025de149`                    | fix         | **[High]** intake: symlink guard for BACKLOG.md and active plan dir     | 2026-09-15 | [plan](closed/2026/09/15/025de149/plan.md)                    |
| `29b18cef`                    | fix         | **[High]** source-layout verify: write candidate tree under sourceRoot  | 2026-09-15 | [plan](closed/2026/09/15/29b18cef/plan.md)                    |
| `cop-path-traversal`          | fix         | **[High]** Security: path traversal in work doc folderID                | 2026-09-14 | [plan](closed/2026/09/14/cop-path-traversal/plan.md)          |
| `620ba34c`                    | infrastructure | Split domain work into cohesive subpackages                    | 2026-09-19 | [plan](closed/2026/09/19/620ba34c/plan.md)                    |
| `9c660e32`                    | infrastructure | Remove filesystem operations from domain work                  | 2026-09-19 | [plan](closed/2026/09/19/9c660e32/plan.md)                    |
| `a8797f44`                   | feature        | `hawp_work_intake` compound search + reshape                    | 2026-09-19 | [plan](closed/2026/09/19/a8797f44/plan.md)                    |
---

## Archive

- Closed work: `closed/`
- Status reports: `status/`
- Evidence: `evidence/`
- Decisions: `decisions/`

### Archived Recently-Closed (compacted 2026-09-16)

| ID                            | Type        | Title                                                                   | Closed     | Detail                                                        |
| ----------------------------- | ----------- | ----------------------------------------------------------------------- | ---------- | ------------------------------------------------------------- |
| `cop-mcp-limit-overflow`      | fix         | **[High]** Security: MCP search limit overflow before multiply-by-3     | 2026-09-14 | [plan](closed/2026/09/14/cop-mcp-limit-overflow/plan.md)      |
| `cop-root-flag-bypass`        | fix         | **[High]** Security: `--root` flag bypass in source-layout run.sh       | 2026-09-14 | [plan](closed/2026/09/14/cop-root-flag-bypass/plan.md)        |
| `cop-doc-overwrite`           | fix         | **[Medium]** Bug: silent overwrite of existing work docs                | 2026-09-14 | [plan](closed/2026/09/14/cop-doc-overwrite/plan.md)           |
| `cop-mcp-contract-mismatch`   | fix         | **[Medium]** Bug: MCP contract mismatches (state, JSON tags, draft)     | 2026-09-14 | [plan](closed/2026/09/14/cop-mcp-contract-mismatch/plan.md)   |
| `cop-reshape-backend-default` | fix         | **[Medium]** Bug: `--reshape-backend` defaults to onnx                  | 2026-09-14 | [plan](closed/2026/09/14/cop-reshape-backend-default/plan.md) |
| `cop-onnx-model-default`      | improvement | **[Low]** ONNX default model updated to Phi-3-mini                      | 2026-09-14 | [plan](closed/2026/09/14/cop-onnx-model-default/plan.md)      |

### Archived Recently-Closed (compacted 2026-09-15)

| ID                            | Type        | Title                                                                   | Closed     | Detail                                                        |
| ----------------------------- | ----------- | ----------------------------------------------------------------------- | ---------- | ------------------------------------------------------------- |
| `288d543c`                    | improvement | Token-reduction benchmark harness for reshape and search context        | 2026-09-12 | [plan](closed/2026/09/12/288d543c/plan.md)                    |
| `a3df8a9c`                    | improvement | Investigate reshape support for HAWP work intake                        | 2026-09-12 | [plan](closed/2026/09/12/a3df8a9c/plan.md)                    |
| `5b6d4e21`                    | improvement | Provider parity for shared HAWP agent guidance                          | 2026-09-11 | [plan](closed/2026/09/11/5b6d4e21/plan.md)                    |
| `d1fa0b72`                    | improvement | Install/update contract hardening for all HAWP providers                | 2026-09-11 | [plan](closed/2026/09/11/d1fa0b72/plan.md)                    |
| `e5fca9c7`                    | release     | HAWP v0.0.24 work-folder normalization and README positioning           | 2026-09-11 | [plan](closed/2026/09/11/e5fca9c7/plan.md)                    |
| `47c793d6`                    | improvement | v0.0.24 CLI decomposition and architecture audit continuation           | 2026-09-11 | [plan](closed/2026/09/11/47c793d6/plan.md)                    |
| `multi-repo-context-d9b2f3a1` | improvement | Multi-repo context skill for agent                                      | 2026-09-10 | [plan](closed/2026/09/10/multi-repo-context-d9b2f3a1/plan.md) |
| `8ddb06ea`                    | improvement | Move filesystem ops out of domain/kit, kitsync, provision, distribution | 2026-09-10 | [plan](closed/2026/09/10/8ddb06ea/plan.md)                    |

### Archived Recently-Closed (compacted 2026-09-12)

| ID                        | Type        | Title                                                          | Closed     | Detail                                                    |
| ------------------------- | ----------- | -------------------------------------------------------------- | ---------- | --------------------------------------------------------- |
| `74aaa332`                | improvement | Strengthen HAWP slice harnesses — providers, tools, UUID artifacts | 2026-09-10 | [plan](closed/2026/09/10/74aaa332/plan.md)            |
| `742aa60b`                | improvement | Separate work domain rules, use cases, and filesystem adapters | 2026-09-10 | [plan](closed/2026/09/10/742aa60b/plan.md)                |
| `b61770a6`                | improvement | Move model and usage adapters out of domain packages           | 2026-09-09 | [plan](closed/2026/09/09/b61770a6/plan.md)                |
| `e7a5e294`                | improvement | Group CLI commands into nested family/operation packages        | 2026-09-08 | [plan](closed/2026/09/08/e7a5e294/plan.md)                |
| `0cb0f9b0`                | improvement | Single-executable installation and configuration migration     | 2026-09-07 | [plan](closed/2026/09/07/0cb0f9b0/plan.md)                |
| `03299078`                | improvement | v0.0.23 TypeScript script deprecation plan and Go parity audit | 2026-08-31 | [plan](closed/2026/08/31/03299078/plan.md)                |
| `5957aaf4`                | task        | v0.0.23 follow-up architecture audit and simplification queue  | 2026-08-31 | [plan](closed/2026/08/31/5957aaf4/plan.md)                |
| `c804eec0`                | test        | v0.0.23 benchmark evidence and release artifact refresh        | 2026-08-31 | [plan](closed/2026/08/31/c804eec0/plan.md)                |

### Archived Recently-Closed (compacted 2026-09-10)

| ID                        | Type        | Title                                                   | Closed     | Detail                                                    |
| ------------------------- | ----------- | ------------------------------------------------------- | ---------- | --------------------------------------------------------- |
| `0bd32051`                | refactor    | v0.0.23 root tests layout and port adapter segregation  | 2026-08-31 | [plan](closed/2026/08/31/0bd32051/plan.md)                |
| `387a37b7`                | improvement | v0.0.23 refactor search pipeline and layout unification | 2026-08-31 | [plan](closed/2026/08/31/387a37b7/plan.md)                |
| `v0014-token-speed-bench` | feature     | Token-savings and speed benchmark gate                  | 2026-08-30 | [plan](closed/2026/08/30/v0014-token-speed-bench/plan.md) |

### Archived Recently-Closed (compacted 2026-08-24)

| ID                        | Type        | Title                                                                         | Closed     | Detail                                                    |
| ------------------------- | ----------- | ----------------------------------------------------------------------------- | ---------- | --------------------------------------------------------- |
| `b7e2a4f9`                | refactor    | Rename `librarian/go/` → `src/` + retire TS validators                        | 2026-08-22 | [plan](closed/2026/08/22/b7e2a4f9/plan.md)                |
| `b8d3e1f0`                | batch-close | All v0.0.1/0.0.2/0.0.3 planned work — shipped in 0.0.1 release (2026-08-21)   | 2026-08-22 | [index](closed/2026/08/22/b8d3e1f0/plan.md)               |
| `cursor-mcp-type-wrapper` | fix         | Cursor MCP: add `type: stdio` and `hawp-mcp` wrapper (relative command fails) | 2026-08-25 | [plan](closed/2026/08/25/cursor-mcp-type-wrapper/plan.md) |
| `w4x6y8z0`                | fix         | `--llm-reshape` wired into the CLI                                            | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md)         |
| `n4o6p8q0`                | fix         | `ReshapingConfig` now honors configured Ollama URL                            | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md)         |
| `j2k4l6m8`                | fix         | `DefaultConfig()` no longer defaults to nonexistent ONNX LLM model            | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md)         |
| `a1b3c5d7`                | fix         | 8 failing Ollama LLM tests repaired                                           | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md)         |
| `e9f1g3h5`                | fix         | Token budget `1995/0` display bug fixed                                       | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md)         |
| `i7j9k1l3`                | fix         | `aes256` silent downgrade now rejected at config validation                   | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md)         |
| `r5s7t9u1`                | fix         | `Reshape()` empty-block guard fixed                                           | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md)         |

## Notes

- Check this file before starting any new item.
- Active plan files go in `active/{uuid}/plan.md`. Close by moving to `closed/YYYY/MM/DD/{uuid}/plan.md`.
- Deferred items live in `parked/{id}/plan.md` without being closed.
- ADRs and decisions go in `decisions/YYYY/MM/DD/`.
- Each item gets one plan folder — no two agents on the same ID.
- Keep `Recently Closed` capped at 10; archive history lives in `closed/`.

## Future Improvements

- **Evidence recovery for legacy unproven claims** (optional, low priority): 85 legacy claims carry explicit unproven annotations (`a7ebe68a`, 2026-07-20).
