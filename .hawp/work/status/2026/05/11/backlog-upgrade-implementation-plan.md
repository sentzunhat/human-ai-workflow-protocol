# Implementation Status: Backlog Upgrade Command

**Date:** 2026-05-11
**Status:** TASK-029 complete and closed; TASK-027 parse-only slice complete
**Next Phase:** Start TASK-028 detect/report-only slice after approval

---

## Design Approval Summary

Design document: [hawp-backlog-upgrade-command-design.md](../../../../notes/2026/05/11/hawp-backlog-upgrade-command-design.md)

**Approved scope:**

- Command: `./.hawp/bin/hawp backlog upgrade`
- Default mode: `--dry-run` (safe, non-destructive)
- Modes: `--dry-run`, `--apply`, `--validate`, `--export-plan`
- Output: text and JSON formats
- V1 architecture: deterministic, mechanical, non-AI
- Safety: 9 explicit safety rules + 3 governance principles
- Auditability: immutable SHA256 hashes for all operations
- Data model: JSON-first (objects as source-of-truth, text/json as renderers)
- Operational guarantee: idempotency (apply twice = stable state)
- Governance: validator remains authoritative

---

## Implementation Plan (V1)

### Slice 1: Foundation + Detection + Dry-Run

**Timeline estimate:** ~5-7 days for serial implementation

#### Phase 1.1: Data Models (TASK-029) — **START HERE**

**Effort:** Low
**Duration:** 1 day

Creates TypeScript type definitions:

- `DetectionReport` — scan results
- `BacklogFixPlan` — complete plan with operations
- `BacklogFixOperation` — individual fixes
- `BlockedItem` — blocked items with rule/confidence/candidates/reason
- `EvidenceReport` — post-apply summary with immutable hashes
- `RuleId`, `OperationType` — shared enums

**Unblocks:** TASK-027, TASK-028 (both depend on these models)

#### Phase 1.2: CLI Entry Point (TASK-027) — **PARSE ONLY**

**Effort:** Low
**Duration:** 1-1.5 days

Creates scaffolding only:

- `./.hawp/bin/hawp` — Bash entry point
- `librarian/scripts/backlog-upgrade/cli.ts` — Argument parser
- Validates mutually exclusive flags (--dry-run vs --apply)
- Validates combined flags (--validate with both modes allowed)

Not implemented in TASK-027:

- no detection engine
- no backlog scanning
- no dry-run report generation
- no apply/write behavior
- no validator execution
- no evidence report writing

**Deliverable:** `./hawp backlog upgrade --help` works

#### Phase 1.3: Detection + Dry-Run Report (TASK-028) — **DETECT/REPORT ONLY**

**Effort:** Medium
**Duration:** 3-4 days

Implements detect/report logic only:

- Backlog parser — parse `.hawp/work/BACKLOG.md`
- Plan scanner — locate and read all plan files
- Rule evaluator — apply A1-A7 auto-fixes, detect B1-B6 blocks
- Report generator — compile findings into structured report
- Formatters — text and JSON output renderers

Not implemented in TASK-028:

- no file writes
- no apply mode execution
- no validator execution
- no evidence report writing

**Deliverable:** `./hawp backlog upgrade --dry-run` detects and reports only

#### Phase 1.4: Verification + Integration

**Effort:** Low
**Duration:** 0.5-1 day

- Test all combinations of flags
- Verify no files modified in --dry-run mode
- Validate report structure (text and JSON)
- Confirm exit codes (0 = success, 1 = error, 2 = usage error)
- Verify report includes all rule matches with confidence + candidates

---

### Slice 2: Apply Mode (TASK-030+) — **GATED, AFTER SLICE 1 STABILIZES**

**Prerequisite:** Slice 1 complete and verified

**Estimated effort:** 5-7 days

Implements:

- Plan executor — apply all auto-fix operations
- Backup mechanism — save original files before modifications
- Git helper — optional git staging of changes
- Idempotency checker — verify state after apply
- Evidence report generator — immutable hashes for audit trail

**Will be gated:** Requires separate approval before starting

---

### Slice 3: Validation Integration (TASK-031+) — **GATED, AFTER SLICE 2**

**Prerequisite:** Slice 2 complete

**Estimated effort:** 2-3 days

Implements:

- Validator rerun after apply
- Validator state hash collection (before/after)
- Report comparison (issues reduced/eliminated)
- CI/CD integration readiness

---

## Work Items & Sequencing

### Recommended sequence:

1. **Start:** TASK-029 (data models) — foundational, unblocks others
2. **Parallel:** TASK-027 (CLI) + TASK-028 (detection) — independent once TASK-029 exists
3. **Then:** Integrate and verify (merge, test combinations)
4. **Gate:** Pause before apply mode until dry-run is stable
5. **Later:** TASK-030+ (apply mode — separate approval cycle)

### Blocking relationships:

```
TASK-029 (models) ──┬──→ TASK-027 (CLI)
                    └──→ TASK-028 (detection)
                         └──→ Verify + integrate
                         └──→ GATE: Pause before TASK-030

TASK-030 (apply) ──→ TASK-031 (validation)
```

---

## Design-to-Implementation Mapping

### Core design elements covered in first slice:

| Design Element             | Implemented In      | Status    |
| -------------------------- | ------------------- | --------- |
| Command shape              | TASK-027            | ✓ Planned |
| Argument parsing           | TASK-027            | ✓ Planned |
| Auto-fix types A1-A7       | TASK-028            | ✓ Planned |
| Blocked types B1-B6        | TASK-028            | ✓ Planned |
| Rule/confidence/candidates | TASK-029 models     | ✓ Planned |
| JSON-first architecture    | TASK-029 + TASK-028 | ✓ Planned |
| Text/JSON output           | TASK-028            | ✓ Planned |
| Dry-run safety             | TASK-027 + TASK-028 | ✓ Planned |
| No semantic mutation       | TASK-028 rules      | ✓ Planned |
| Path boundary (`.hawp/**`) | TASK-028 rules      | ✓ Planned |

### Deferred to later slices:

| Design Element        | Deferred To         | Reason                         |
| --------------------- | ------------------- | ------------------------------ |
| Apply mode            | TASK-030            | Gated until dry-run stable     |
| Immutable hashes      | TASK-030 + TASK-032 | Needed only for apply/evidence |
| Idempotency guarantee | TASK-030            | Applies to apply mode          |
| Validator integration | TASK-031            | Separate gate                  |
| Evidence reports      | TASK-032            | Post-apply only                |
| V2+ AI synthesis      | V2+                 | Governance gates required      |

---

## Quality Gates

Each slice must pass before proceeding:

### Slice 1 Gate (before starting TASK-030):

- [ ] `./hawp backlog upgrade --help` works
- [ ] `./hawp backlog upgrade --dry-run` completes in < 2s
- [ ] Detects all A1-A7 auto-fixes with confidence > threshold
- [ ] Detects all B1-B6 blocked items with rule/confidence/candidates
- [ ] Text output is human-readable with clear sections
- [ ] JSON output is valid, parseable, includes all fields
- [ ] No files modified (--dry-run confirmed)
- [ ] Exit code 0 for success, 1 for errors
- [ ] Tested on clean backlog (no drift) — should report "No modifications needed"
- [ ] Tested on known drift scenarios — should match design examples

### Slice 2 Gate (before starting TASK-031):

- [ ] `./hawp backlog upgrade --apply` executes auto-fixes
- [ ] Modified files have correct content
- [ ] Blocked items are skipped (not modified)
- [ ] Backup files created (if implemented)
- [ ] Running `--apply` twice produces zero changes on second run
- [ ] Evidence report includes immutable hashes
- [ ] Evidence report is immutable and persistent

### Slice 3 Gate (before production):

- [ ] Validator rerun produces expected drift reduction
- [ ] Validator state hashes collected correctly
- [ ] CI/CD integration works
- [ ] Full integration testing across HAWP repo

---

## Risk & Mitigation

| Risk                          | Probability | Impact | Mitigation                                           |
| ----------------------------- | ----------- | ------ | ---------------------------------------------------- |
| Detection rules incomplete    | Medium      | High   | Reference design specs, iterative rule refinement    |
| Idempotency edge cases        | Medium      | Medium | Explicit testing of repeated applies                 |
| Performance on large backlogs | Low         | Medium | Benchmark on real data, optimize if needed           |
| Path boundary escapes         | Low         | High   | Explicit guard assertions in all write operations    |
| Validator incompatibility     | Low         | High   | Coordinate with validator (TASK-020) before TASK-031 |

---

## Success Criteria

**Slice 1 complete when:**

- All three tasks (TASK-027, TASK-028, TASK-029) are implemented and verified
- Dry-run report matches design examples
- No modifications made by dry-run (safety confirmed)
- Exit codes correct
- Text and JSON formats both valid and complete

**Overall V1 complete when:**

- Apply mode works and is idempotent
- Validation integration works
- All quality gates passed
- Evidence reports are immutable and correct

---

## Next Immediate Steps

1. ✅ Approve implementation plan (this document)
2. → Review TASK-029 plan file (data models)
3. → Begin implementation: start with TASK-029
4. → Proceed to TASK-027 and TASK-028 in parallel once TASK-029 models exist
5. → Verify and gate before TASK-030 (apply mode)
