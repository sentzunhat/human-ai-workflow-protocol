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
| `b6c4e8a2` | — | task | Audit and align librarian architecture with domain port/adapter boundaries | in-progress | unassigned | [plan](active/b6c4e8a2-architecture-port-adapter-audit.md) | 2026-08-21 |
| `c1d2e402` | — | fix | Extract work normalization scan and mutation boundary | plan-ready | unassigned | [plan](active/c1d2e402-fix-work-normalization-boundary.md) | 2026-08-13 |
| `392313e4-eded-402b-9d5e-20350c86b856` | — | migration | Normalize dated UUID work-item folders and repair workflow references | plan-ready | unassigned | [plan](active/392313e4-eded-402b-9d5e-20350c86b856-workspace-identity-migration.md) | 2026-08-13 |
| `c1d2e3fe` | — | audit | Recursive audit: SQLite infrastructure capabilities | plan-ready | unassigned | [plan](active/c1d2e3fe-audit-sqlite-capabilities.md) | 2026-08-10 |
| `c1d2e3ff` | — | fix | Group SQLite persistence by capability and remove raw result maps | plan-ready | unassigned | [plan](active/c1d2e3ff-fix-sqlite-capabilities.md) | 2026-08-10 |
| `c1d2e400` | — | audit | Recursive audit: CLI command capabilities | plan-ready | unassigned | [plan](active/c1d2e400-audit-cli-capabilities.md) | 2026-08-10 |
| `c1d2e401` | — | fix | Split CLI routing by command capability | plan-ready | unassigned | [plan](active/c1d2e401-fix-cli-capability-splits.md) | 2026-08-10 |
| `c42d1443-fbd0-4a98-8379-6058ace02dc4` | — | fix | Unify search index through shared build and ingest flow | plan-ready | unassigned | [plan](active/c42d1443-fbd0-4a98-8379-6058ace02dc4-unify-search-index-flow.md) | 2026-08-15 |
| `fcb5c8a2-eb3e-4554-9889-2e9118f4ad02` | — | feature | Optional HAWP parallel coordination execution space | plan-ready | unassigned | [plan](active/fcb5c8a2-eb3e-4554-9889-2e9118f4ad02-parallel-coordination-space.md) | 2026-08-15 |
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

Limited to the last 30 days.

| ID       | Type             | Title                                                                     | Closed     | Detail                                |
| -------- | ---------------- | ------------------------------------------------------------------------- | ---------- | ------------------------------------- |
| `c1d2e3fd` | fix | Separate context retrieval, formatting, and reshape seams | 2026-08-21 | [plan](closed/2026/08/21/c1d2e3fd-fix-context-seams.md) |
| `c1d2e3fb` | fix | Introduce typed index persistence capability contracts | 2026-08-21 | [plan](closed/2026/08/21/c1d2e3fb-fix-index-persistence-boundary.md) |
| `c1d2e3f9` | fix | Extract work source and link-resolution boundaries | 2026-08-21 | [plan](closed/2026/08/21/c1d2e3f9-fix-domain-work-boundary.md) |
| `d7e5a1b9` | infrastructure/docs | Propagate release-readiness guidance across all provider packs | 2026-08-21 | [plan](closed/2026/08/21/d7e5a1b9-release-readiness-provider-alignment.md) |
| `c1d2e3fc` | audit | Recursive audit: application context capability | 2026-08-13 | [plan](closed/2026/08/13/c1d2e3fc-audit-application-context.md) |
| `c1d2e3fa` | audit | Recursive audit: application index capability | 2026-08-13 | [plan](closed/2026/08/13/c1d2e3fa-audit-application-index.md) |
| `c1d2e3f8` | audit | Recursive audit: domain work capability | 2026-08-13 | [plan](closed/2026/08/13/c1d2e3f8-audit-domain-work.md) |
| `c1d2e3f7` | fix | Isolate kit content input from normalization and validation | 2026-08-13 | [plan](closed/2026/08/13/c1d2e3f7-fix-domain-kit-boundary.md) |
| `c1d2e3f6` | audit | Recursive audit: domain kit capability | 2026-08-11 | [plan](closed/2026/08/11/c1d2e3f6-audit-domain-kit.md) |
| `c1d2e3f5` | fix | Extract domain context corpus/source boundary | 2026-08-11 | [plan](closed/2026/08/11/c1d2e3f5-fix-domain-context-boundary.md) |
| `c1d2e3f4` | audit | Recursive audit: domain context capability | 2026-08-10 | [plan](closed/2026/08/10/c1d2e3f4-audit-domain-context.md) |
| `u7d4c8a1` | task | Audit docs drift and checkpoint alignment | 2026-08-10 | [plan](closed/2026/08/10/u7d4c8a1-audit-docs-drift.md) |
| `c9a7f2e1` | infrastructure | Cross-platform GitHub Actions CI/CD pipeline (6-binary matrix) | 2026-07-23 | [plan](closed/2026/07/23/c9a7f2e1-github-actions-pipeline.md) |
| `d1b3e8f4` | infrastructure | Repo audit & cleanup: verified clean, all systems ready for v0.0.1 | 2026-07-23 | [plan](closed/2026/07/23/d1b3e8f4-repo-audit-cleanup.md) |
Older closed work (`w4x6y8z0`, `r5s7t9u1`, `e2c4f9g5`, `4b114c0c`, `273f3e4b`, `748609a8`, `b95436f2`, `f93bee55` and earlier) is archived under `closed/YYYY/MM/DD/`.

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
