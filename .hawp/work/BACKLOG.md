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

| UUID | Legacy ID | Type | Title | Status | Owner | Plan File | Updated |
| ---- | --------- | ---- | ----- | ------ | ----- | --------- | ------- |
| `c9a7f2e1` | — | infrastructure | Cross-platform GitHub Actions CI/CD pipeline (6-binary matrix) | done | — | [plan](active/c9a7f2e1-github-actions-pipeline.md) | 2026-07-23 |
| `d1b3e8f4` | — | infrastructure | Repo audit & cleanup: verified clean, all systems ready for v0.0.1 | done | — | [evidence](evidence/2026/07/23/d1b3e8f4-repo-clean.md) | 2026-07-23 |
| `h5f7c2j8` | — | fix | Retry v0.0.1 release: tag fixed workflow + push (GitHub Actions corrected) | plan-ready | unassigned | [plan](active/h5f7c2j8-retry-v001-release.md) | 2026-07-23 |
| `f3d5a0h6` | — | test | Release verification: test all 6 binaries cross-platform | plan-ready | unassigned | [plan](active/f3d5a0h6-release-verification.md) | 2026-07-23 |
| `i6g8d3k9` | — | feature | Context Packing (Slice 4): v0.0.2 - Phases 1-5 (NO APIs) | done | unassigned | [plan](active/i6g8d3k9-context-packing-slice4.md) | 2026-07-24 |
|   | ✓ | — | ├─ Phase 1: Deduplication & Grouping (9 tests ✅) | done | — | — | 2026-07-23 |
|   | ✓ | — | ├─ Phase 2: Formatting & Truncation (6 tests ✅) | done | — | — | 2026-07-23 |
|   | ⏹️ | — | ├─ Phase 3: Context Reshaping via LLM Inference (ONNX + embeddings + API backends) ↓ DEFERRED TO v0.0.3+ | parked | — | — | 2026-07-24 |
|   | ✓ | — | ├─ Phase 4: CLI Integration (`hawp search --context`) | done | — | — | 2026-07-24 |
|   | ✓ | — | └─ Phase 5: Testing & Documentation | done | — | — | 2026-07-24 |
| `m2o4p6q8` | — | release | Tag v0.0.2: Release with Context Packing (Phases 1-2, 4-5) | plan-ready | unassigned | [plan](active/m2o4p6q8-tag-v002.md) | 2026-07-27 |
| `j7h9e4l0` | — | test | Auto-update testing: v0.0.1 → v0.0.2 (test update command on all 6 platforms) | plan-ready | unassigned | — | 2026-07-24 |
| `k8i9f5m1` | — | feature | v0.0.3: Phase 3 - Context Reshaping (Configurable LLM + Embeddings backends) | done | unassigned | [plan](active/i6g8d3k9-context-packing-slice4.md) | 2026-07-25 |
|   | ✓ | — | ├─ Phase 3.1: Config System & Secure Key Storage (DONE + README auto-gen) | done | — | [detail](active/k8i9f5m1-phase3-1-config-system.md) | 2026-07-24 |
|   | ✓ | — | ├─ Phase 3.2a: ONNX Embeddings (Hugot integration DONE, 10 tests) | done | — | [detail](active/k8i9f5m1-phase3-2-embeddings.md) | 2026-07-24 |
|   | ✓ | — | ├─ Phase 3.2b: Ollama Embeddings (HTTP API DONE, 10 tests) | done | — | [detail](active/k8i9f5m1-phase3-2-embeddings.md) | 2026-07-25 |
|   | ⏹️ | — | ├─ Phase 3.2c/d: OpenAI/Anthropic Embeddings → deferred to v0.1.0 | parked | — | [detail](active/k8i9f5m1-phase3-2-embeddings.md) | 2026-07-25 |
|   | ✓ | — | ├─ Phase 3.3a: ONNX LLM (Scaffolding DONE, 8 tests, awaiting ONNX LLM models) | done | — | [detail](active/k8i9f5m1-phase3-3-llm.md) | 2026-07-24 |
|   | ✓ | — | ├─ Phase 3.3b: Ollama LLM (HTTP API DONE, 10 tests) | done | — | [detail](active/k8i9f5m1-phase3-3-llm.md) | 2026-07-25 |
|   | ⏹️ | — | ├─ Phase 3.3c/d: OpenAI/Anthropic LLM → deferred to v0.1.0 | parked | — | [detail](active/k8i9f5m1-phase3-3-llm.md) | 2026-07-25 |
|   | ✓ | — | ├─ Phase 3.4: Context Reshaping Pipeline (DONE, 100+ tests) | done | — | [detail](active/k8i9f5m1-phase3-4-5-pipeline.md) | 2026-07-25 |
|   | ✓ | — | └─ Phase 3.5: Testing & Documentation (DONE — 4 docs, version bumped) | done | — | [detail](active/k8i9f5m1-phase3-4-5-pipeline.md) | 2026-07-25 |
| `p8r0t2u4` | — | test | Integration + benchmarks: verify ONNX + Ollama live (embeddings + LLM) | done | — | [benchmarks](../../librarian/docs/BENCHMARKS_v003.md) | 2026-07-26 |
| `n3p5r7s9` | — | release | Tag v0.0.3: ship context reshaping (ONNX + Ollama backends) | plan-ready | unassigned | [plan](active/n3p5r7s9-tag-v003.md) | 2026-07-27 |
| `t5l7m9n1` | — | planning | Benchmarking Infrastructure for v0.0.4+: comprehensive provider matrix plan + test harness design | plan-ready | unassigned | [plan](active/t5l7m9n1-benchmark-plan-v004.md) | 2026-07-26 |
| `g4e6b1i7` | — | task | Final: Update version.go + CHANGELOG.md, tag v0.0.1, push | done | — | [evidence](evidence/2026/07/23/g4e6b1i7-release-executed.md) | 2026-07-23 |
| `fbf12a93` | — | feature | Vector search + local context building on downloaded ONNX models (epic; slices below) | analyzing | unassigned | [plan](active/fbf12a93-7700-4236-b6bb-e7b55b63a622.md) | 2026-07-22 |
| `77a6879a` | — | feature | ├─ Slice 2: vector embedding via ONNX | analyzing | unassigned | [plan](active/77a6879a-vector-embedding-onnx.md) | 2026-07-22 |

---

## Blocked / Parked

| ID  | Type | Title | Reason | Detail | Updated |
| --- | ---- | ----- | ------ | ------ | ------- |
| `bee15107` | improvement | defer CLI participant adapters for Codex, Claude, and GitHub | Not needed yet; provider packs are enough for current use | [plan](parked/bee15107-e4d1-41b3-b67e-04c87328fef3.md) | 2026-07-06 |
| `v010-3-3a` | feature | ONNX LLM Text2Text models (FLAN-T5-small) | No v0.1.0 milestone; all future work ships as incremental v0.0.4+ patches | [detail](parked/v010-phase3-3a-onnx-llm-feasible.md) | 2026-07-27 |
| `v010-3-2c` | feature | OpenAI Embeddings (text-embedding-3-small/large) | Cloud-API tracks parked indefinitely; see incremental-only strategy | [detail](parked/v010-phase3-2c-openai-embeddings.md) | 2026-07-27 |
| `v010-3-3c` | feature | OpenAI LLM (gpt-3.5-turbo, gpt-4-turbo) | Cloud-API tracks parked indefinitely; see incremental-only strategy | [detail](parked/v010-phase3-3c-openai-llm.md) | 2026-07-27 |
| `v010-3-2d` | feature | Anthropic Embeddings Stub | Cloud-API tracks parked indefinitely; see incremental-only strategy | [detail](parked/v010-phase3-2d-anthropic-embeddings.md) | 2026-07-27 |
| `v010-3-3d` | feature | Anthropic LLM (claude-3-sonnet, claude-3-opus) | Cloud-API tracks parked indefinitely; see incremental-only strategy | [detail](parked/v010-phase3-3d-anthropic-llm.md) | 2026-07-27 |
| `v010-cost` | feature | Cost Tracking + Rate Limiting (cross-cutting for cloud backends) | Cloud-API tracks parked indefinitely; see incremental-only strategy | [detail](parked/v010-cost-tracking-rate-limiting.md) | 2026-07-27 |

---

## Recently Closed

Limited to the last 10 items.

| ID       | Type             | Title                                                                     | Closed     | Detail                                |
| -------- | ---------------- | ------------------------------------------------------------------------- | ---------- | ------------------------------------- |
| `w4x6y8z0` | fix | `--llm-reshape` wired into the CLI (shared record, see `v003-ship-audit`) | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit.md) |
| `n4o6p8q0` | fix | `ReshapingConfig` now honors configured Ollama URL (shared record) | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit.md) |
| `j2k4l6m8` | fix | `DefaultConfig()` no longer defaults to nonexistent ONNX LLM model (shared record) | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit.md) |
| `a1b3c5d7` | fix | 8 failing Ollama LLM tests repaired (shared record) | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit.md) |
| `e9f1g3h5` | fix | token budget `1995/0` display bug fixed (shared record) | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit.md) |
| `i7j9k1l3` | fix | `aes256` silent downgrade now rejected at config validation (shared record) | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit.md) |
| `r5s7t9u1` | fix | `Reshape()` empty-block guard fixed (shared record) | 2026-07-26 | [plan](closed/2026/07/26/v003-ship-audit.md) |
| `e2c4f9g5` | infrastructure | update command v2 (platform detect, download, SHA256 verify) was already complete | 2026-07-23 | [plan](closed/2026/07/23/e2c4f9g5-update-already-complete.md) |
| `4b114c0c` | feature | document index ingest pipeline (fbf12a93 Slice 1: SQLite + FTS5 lexical search) | 2026-07-22 | [plan](closed/2026/07/22/4b114c0c-document-index-ingest.md) |
| `273f3e4b` | feature | hawp update also syncs .hawp/kit/ and the installed provider overlay (unified version) | 2026-07-21 | [plan](closed/2026/07/21/273f3e4b-4978-4fc1-ae8c-2f8ecc99402a.md) |
| `748609a8` | feature | adopt hugot as a transformers.js-style Go wrapper for hawp (model pull + embed, pure-Go, no cgo) | 2026-07-21 | [plan](closed/2026/07/21/748609a8-ac1a-4e27-a0dc-7a67d3b9fd14.md) |
| `b95436f2` | feature | draft-and-approve release workflow: workflow_dispatch cuts the tag, builds, auto-publishes or drafts | 2026-07-21 | [plan](closed/2026/07/21/b95436f2-e9e2-4259-ac65-53a28068ccab.md) |
Older closed work (`f93bee55` / 2026-07-21 and earlier) is archived under `closed/YYYY/MM/DD/`.

---

## Archive

- Closed work: `closed/`
- Status reports: `status/`
- Evidence: `evidence/`
- Decisions: `decisions/`

## Notes

- Check this file before starting any new item.
- Active plan files go in `active/`. Close by moving to `closed/YYYY/MM/DD/`.
- Deferred items can live in `parked/` without being closed.
- ADRs and decisions go in `decisions/YYYY/MM/DD/`.
- Each item gets one plan file — no two agents on the same ID.
- Work started outside this loop should still get a row added for visibility.
- Keep `Recently Closed` capped; archive history lives in `closed/`.
- Provider rollout (TASK-072–075), standards adaptation (TASK-070), audit remediation (TASK-076), and scripts structure alignment (TASK-077) are closed.
- TASK-080 and TASK-081 closed 2026-06-29 after a paired-lane workflow-loop trial; see closed plans under `closed/2026/06/29/`.

## Future Improvements

Sequenced 2026-07-20 — none of these become work items until their gate is met:

- **UUID migration Phase 3 (optional)**: retroactive migration of closed records to UUID paths. Gate: only after Go indexing lands (`fbf12a93`) and only if indexing needs it. Phases 1–2 closed 2026-07-03; see `notes/2026/07/03/migration-sequential-to-uuid.md`.
- **Optional MCP server overlay for librarian tooling** (deferred, much later): thin MCP wrapper over kit/work validate + normalize + uuid. CLI stays canonical. Decision context: `closed/2026/07/03/2bcbf995-983a-4c23-ab2b-2703ef50d477.md`.
- **Evidence recovery for legacy unproven claims** (optional, low priority): 85 legacy claims now carry explicit unproven annotations (`a7ebe68a`, 2026-07-20); individual claims can be upgraded to `Evidence:` if ever re-verified. The clarity WARN is honest debt, not ambiguity.
- **Retire TS validators after CI switches to Go** (gated on CI): `.hawp/bin/hawp` now prefers the Go binary (2026-07-20); the TS kit/work validate + normalize scripts remain as the npm-pipeline fallback. Once CI builds and runs `librarian/go/bin/hawp`, the TS equivalents and their npm scripts can be removed.
- **Delete the `librarian-go-v0.0.1`/`v0.0.2` test releases and tags** (optional cleanup, low priority): both were real end-to-end proof of the release/update mechanism (`4c152ee3`, closed 2026-07-21) — safe to remove once no longer useful as a reference, or keep as a documented example.
