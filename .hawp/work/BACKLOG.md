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

_No active work items. `v0.0.23` release lanes are closed or intentionally
parked._

| UUID | Type | Title | Status | Owner | Plan File | Updated |
| ---- | ---- | ----- | ------ | ----- | --------- | ------- |

## Blocked / Parked

| ID  | Type | Title | Reason | Detail | Updated |
| --- | ---- | ----- | ------ | ------ | ------- |
| `89cf7a85` | improvement | v0.0.23 legacy work-item UUID canonicalization follow-up | Older-repo safety and archive-policy follow-up; not a live `v0.0.23` blocker | [plan](parked/89cf7a85/plan.md) | 2026-08-31 |
| `release-benchmark-backfill` | planning | Benchmark evidence gate for `v0.0.11` through `v0.0.13` | Older patch-train backfill, intentionally kept out of the `v0.0.23` release lane | [plan](parked/release-benchmark-backfill/plan.md) | 2026-08-31 |
| `usage-tracking` | feature | `hawp usage` CLI + `hawp_usage` MCP tool for token/wall-clock metering | Scope too large for v0.0.11; belongs at v0.1.0 token-reduction gate | [plan](parked/usage-tracking/plan.md) | 2026-08-25 |
| `bee15107` | improvement | Defer CLI participant adapters for Codex, Claude, and GitHub | Not needed yet; provider packs are enough | [plan](parked/bee15107/plan.md) | 2026-07-06 |
| `v010-3-3a` | feature | ONNX LLM Text2Text models (FLAN-T5-small) | Plan ready (FLAN-T5-small feasible via ONNX); deferred per incremental-only patch strategy | [detail](parked/v010-3-3a/plan.md) | 2026-08-24 |
| `v010-3-2c` | feature | OpenAI Embeddings (text-embedding-3-small/large) | Cloud-API tracks parked; see incremental-only strategy | [detail](parked/v010-3-2c/plan.md) | 2026-07-27 |
| `v010-3-3c` | feature | OpenAI LLM (gpt-3.5-turbo, gpt-4-turbo) | Cloud-API tracks parked; see incremental-only strategy | [detail](parked/v010-3-3c/plan.md) | 2026-07-27 |
| `v010-3-2d` | feature | Anthropic Embeddings Stub | Cloud-API tracks parked; see incremental-only strategy | [detail](parked/v010-3-2d/plan.md) | 2026-07-27 |
| `v010-3-3d` | feature | Anthropic LLM (claude-3-sonnet, claude-3-opus) | Cloud-API tracks parked; see incremental-only strategy | [detail](parked/v010-3-3d/plan.md) | 2026-07-27 |
| `v010-cost` | feature | Cost Tracking + Rate Limiting (cross-cutting for cloud backends) | Cloud-API tracks parked; see incremental-only strategy | [detail](parked/v010-cost/plan.md) | 2026-07-27 |

---

## Recently Closed

Limited to the last 10 items.

| ID | Type | Title | Closed | Detail |
| -- | ---- | ----- | ------ | ------ |
| `03299078` | improvement | v0.0.23 TypeScript script deprecation plan and Go parity audit | 2026-08-31 | [plan](closed/2026/08/31/03299078/plan.md) |
| `5957aaf4` | task | v0.0.23 follow-up architecture audit and simplification queue | 2026-08-31 | [plan](closed/2026/08/31/5957aaf4/plan.md) |
| `c804eec0` | test | v0.0.23 benchmark evidence and release artifact refresh | 2026-08-31 | [plan](closed/2026/08/31/c804eec0/plan.md) |
| `0bd32051` | refactor | v0.0.23 root tests layout and port adapter segregation | 2026-08-31 | [plan](closed/2026/08/31/0bd32051/plan.md) |
| `387a37b7` | improvement | v0.0.23 refactor search pipeline and layout unification | 2026-08-31 | [plan](closed/2026/08/31/387a37b7/plan.md) |
| `v0014-ort-release-fix` | fix | Repair ORT release lanes — shipped v0.0.14–v0.0.16 | 2026-08-27 | [plan](closed/2026/08/27/v0014-ort-release-fix/plan.md) |
| `onnx-llm-release-build` | feature | Ship ONNX LLM in release binary (ORT tarballs) — shipped v0.0.14 | 2026-08-27 | [plan](closed/2026/08/27/onnx-llm-release-build/plan.md) |
| `usage-log` | feature | Local MCP call log with token counts (`hawp usage`) — shipped v0.0.12 | 2026-08-27 | [plan](closed/2026/08/27/usage-log/plan.md) |
| `v0014-token-speed-bench` | feature | Token-savings and speed benchmark gate | 2026-08-30 | [plan](closed/2026/08/30/v0014-token-speed-bench/plan.md) |
| `archive-uuid-structure-audit` | investigation | Review mixed `closed/` and `evidence/` archive structure vs UUID-folder guidance | 2026-08-25 | [plan](closed/2026/08/25/archive-uuid-structure-audit/plan.md) |
| `releases-prerelease-fallback` | investigation | Verify/fix: `/releases/latest` 404 when all GitHub releases are prerelease | 2026-08-25 | [plan](closed/2026/08/25/releases-prerelease-fallback/plan.md) |
| `manager-branch-kit-pattern` | improvement | Document optional manager-branch / worktree operating pattern in kit | 2026-08-25 | [plan](closed/2026/08/25/manager-branch-kit-pattern/plan.md) |
| `agents-seed-if-missing` | fix | AGENTS.md must not be overwritten on `hawp init` update | 2026-08-25 | [plan](closed/2026/08/25/agents-seed-if-missing/plan.md) |
| `install-docs-embed-optional` | docs | Clarify: `hawp search index` first; embed is optional and slow | 2026-08-25 | [plan](closed/2026/08/25/install-docs-embed-optional/plan.md) |
| `cursor-docs-ui-path` | docs | Fix stale Cursor MCP UI path: Customize → MCPs (not Settings → Tools & MCP) | 2026-08-25 | [plan](closed/2026/08/25/cursor-docs-ui-path/plan.md) |

---

## Archive

- Closed work: `closed/`
- Status reports: `status/`
- Evidence: `evidence/`
- Decisions: `decisions/`

### Archived Recently-Closed (compacted 2026-08-24)

| ID | Type | Title | Closed | Detail |
| -- | ---- | ----- | ------ | ------ |
| `b7e2a4f9` | refactor | Rename `librarian/go/` → `src/` + retire TS validators | 2026-08-22 | [plan](closed/2026/08/22/b7e2a4f9/plan.md) |
| `b8d3e1f0` | batch-close | All v0.0.1/0.0.2/0.0.3 planned work — shipped in 0.0.1 release (2026-08-21) | 2026-08-22 | [index](closed/2026/08/22/b8d3e1f0/plan.md) |
| `cursor-mcp-type-wrapper` | fix | Cursor MCP: add `type: stdio` and `hawp-mcp` wrapper (relative command fails) | 2026-08-25 | [plan](closed/2026/08/25/cursor-mcp-type-wrapper/plan.md) |
| `w4x6y8z0` | fix | `--llm-reshape` wired into the CLI | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md) |
| `n4o6p8q0` | fix | `ReshapingConfig` now honors configured Ollama URL | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md) |
| `j2k4l6m8` | fix | `DefaultConfig()` no longer defaults to nonexistent ONNX LLM model | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md) |
| `a1b3c5d7` | fix | 8 failing Ollama LLM tests repaired | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md) |
| `e9f1g3h5` | fix | Token budget `1995/0` display bug fixed | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md) |
| `i7j9k1l3` | fix | `aes256` silent downgrade now rejected at config validation | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md) |
| `r5s7t9u1` | fix | `Reshape()` empty-block guard fixed | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit/plan.md) |

## Notes

- Check this file before starting any new item.
- Active plan files go in `active/{uuid}/plan.md`. Close by moving to `closed/YYYY/MM/DD/{uuid}/plan.md`.
- Deferred items live in `parked/{id}/plan.md` without being closed.
- ADRs and decisions go in `decisions/YYYY/MM/DD/`.
- Each item gets one plan folder — no two agents on the same ID.
- Keep `Recently Closed` capped at 10; archive history lives in `closed/`.

## Future Improvements

- **Evidence recovery for legacy unproven claims** (optional, low priority): 85 legacy claims carry explicit unproven annotations (`a7ebe68a`, 2026-07-20).
