# HAWP Librarian Design Note - Fix-up and Upgrade Commands

## Date

2026-05-11

## Purpose

Define a safe command model for future HAWP cleanup and upgrade flows, using validator findings as the source input and preserving file-based truth.

## Scope

This note covers:

- command naming and intent
- safety model and trust boundaries
- dry-run and apply behavior
- mechanical fix eligibility
- AI drafting boundaries
- hard no-auto-fix boundaries
- dry-run output format
- recommended V1 implementation slice

## Design Inputs

- Existing validator findings for local and external `.hawp` roots
- Existing work files (`BACKLOG.md`, `work/active/`, `work/closed/`, `work/status/`, `work/evidence/`)
- Optional git evidence supplied at runtime (for example, `git status`, `git diff`, or commit metadata)

## Source-of-Truth Rules

1. Files are canonical. No database is required to apply or verify changes.
2. Any index is optional and derived from files.
3. No command may discard source text. Original context must remain recoverable.
4. Verification claims are evidence-bound. If evidence is missing, output must mark claim as unproven.

## Proposed Commands

### 1) `hawp validate`

Primary drift scanner and preflight gate.

- Reads one or more HAWP roots.
- Emits findings with stable finding IDs.
- Supports output for downstream fix-up planning (`json` preferred for machine use).

### 2) `hawp work-items fix-up --dry-run`

Proposes safe, mechanical edits that normalize structure and references without rewriting intent.

- Input: validator findings (directly or via findings artifact).
- Output: a patch plan with explicit per-file hunks and rationale.
- Default mode for `fix-up` family.

### 3) `hawp work-items fix-up --apply`

Applies a previously generated, reviewed fix-up plan.

- Requires a plan artifact hash match to current files, or refuses to apply.
- Writes backup snapshots for touched files (or stores reversible patch artifacts).
- Refuses if preconditions changed since dry-run.

### 4) `hawp work-items upgrade --dry-run`

Proposes schema/layout upgrades for legacy HAWP layouts.

- Input: findings + detected layout generation.
- Output: migration plan, compatibility notes, and data-preservation map.
- Never deletes legacy context; uses move/copy with traceability.

### 5) `hawp backlog refine --dry-run`

Proposes safe backlog normalization:

- compact stale Done rows per policy
- repair plan links
- align table fields/status vocabulary
- keep Recently Closed cap and archive references

### 6) `hawp backlog work-items-upgrade --dry-run`

Specialized backlog-to-work-item upgrade planning.

- Maps backlog rows to missing plan/status/evidence scaffolding.
- Proposes file and link updates while preserving original prose.
- Explicitly flags ambiguous mappings for manual review.

## Safety Model

## Safety Levels

- Level 0 (`validate`): read-only; no file modifications.
- Level 1 (`--dry-run`): generate proposal only (patch/plan/report).
- Level 2 (`--apply`): apply reviewed proposal with guardrails.

## Guardrails

1. Dry-run first by default.
2. Deterministic proposal artifact:
   - includes command, version, root(s), timestamp, finding IDs, and file hashes.
3. Preconditions required for apply:
   - file hash match
   - finding set match
   - no unknown-risk operations in plan
4. Explicit refusal paths:
   - missing evidence for verification claim
   - ambiguous record mapping
   - destructive edits without reversible snapshot
5. Human review required before apply for plans containing:
   - inferred content generation
   - cross-file semantic rewrites
   - archive/close transitions

## Dry-run and Apply Behavior

## `--dry-run`

- Produces:
  - proposed file patch hunks
  - operation classification (`mechanical`, `ai-draft`, `manual-required`)
  - risk tags (`low`, `medium`, `high`)
  - unproven markers where evidence is missing
- Exit codes:
  - `0` no changes proposed
  - `2` changes proposed
  - `3` blocked by ambiguity or safety gate

## `--apply`

- Allowed only for operations classified as apply-safe by policy.
- Requires either:
  - explicit `--plan <path>` from a previous dry-run artifact, or
  - direct confirmation flow for exact generated plan hash.
- Emits:
  - applied changes summary
  - skipped/blocked operations
  - rollback artifact location

## What Can Be Fixed Mechanically

1. Path/link repairs where target can be uniquely resolved.
2. Folder placement corrections for known HAWP structures.
3. Table status normalization to allowed status vocabulary.
4. Missing but derivable relative links (plan/evidence/status) when unambiguous.
5. Date-folder normalization where date is explicit in source context.
6. Consistent renames with traceability notes (old -> new path).

Mechanical fixes must not change semantic meaning of user-authored narratives.

## What Can Be AI-Drafted

AI assist is optional and off for apply unless explicitly approved.

Potential model path:

- local `HuggingFaceTB/SmolLM2-135M-Instruct`
- via WASM/Transformers.js inference

AI may draft only these sections in work item plans:

1. Outcome
2. Verification
3. Close Checklist prefill suggestions

AI input constraints:

- existing file content only
- evidence files and status reports only
- backlog rows only
- optional git evidence if user supplied

AI output constraints:

- must cite source snippets/paths used
- must mark verification lines as `UNPROVEN` unless directly backed by explicit evidence
- must label inferred claims explicitly (`INFERRED`)
- must produce draft blocks only; never auto-commit to source files without review

Human/agent review is mandatory before any apply step.

## What Must Never Be Auto-Fixed

1. Invented verification evidence.
2. Claims of completion not supported by files/evidence.
3. Deletion of historical context or original narrative text.
4. Ambiguous row-to-plan mapping when multiple candidates exist.
5. Risky status transitions (`in-progress` -> `done`) without verification artifacts.
6. Semantic rewrites of Root Cause, Recommended Fix, or decision rationale.
7. Any operation that reduces traceability of when/why a change happened.

## Example Dry-run Output

```text
$ hawp work-items fix-up --dry-run --hawp-root ./project/.hawp --format text

Command: hawp work-items fix-up --dry-run
Root: ./project/.hawp
Findings input: findings-2026-05-11T142201Z.json
Plan ID: plan_01jv0wz7q4
Policy: v1-safe-fixup

Summary:
  findings scanned: 18
  proposed operations: 7
  mechanical: 5
  ai-draft: 1
  manual-required: 1
  blocked: 2

Proposed operations:
  [M-001] mechanical low
    file: .hawp/work/BACKLOG.md
    change: normalize status token "inprogress" -> "in-progress"
    evidence: direct (status vocabulary policy)

  [M-002] mechanical low
    file: .hawp/work/BACKLOG.md
    change: repair broken plan link for TASK-044
    evidence: direct (unique target exists)

  [M-003] mechanical low
    file: .hawp/work/closed/2026/05/09/TASK-041.md
    change: add missing evidence filename reference
    evidence: direct (file exists)

  [A-001] ai-draft medium REVIEW REQUIRED
    file: .hawp/work/active/TASK-052.md
    section: Verification
    draft basis: evidence/2026/05/11/TASK-052-test-output.md
    note: all generated checks tagged UNPROVEN unless explicit pass evidence found

  [X-001] blocked high MANUAL REQUIRED
    file: .hawp/work/BACKLOG.md
    issue: ambiguous mapping for TASK-053 plan link (2 candidate files)

Safety checks:
  file-hash preconditions: recorded
  destructive operations: none
  verification invention: prevented by policy

Result:
  no files changed (dry-run)
  export plan: .hawp/work/status/2026/05/11/plan_01jv0wz7q4.md
  next step: review plan, then run --apply --plan ...
  exit code: 2
```

## Recommended V1 Implementation Slice

1. Keep `hawp validate` as the authoritative finding generator.
2. Add `hawp work-items fix-up --dry-run` with JSON and text output.
3. Support only low-risk mechanical operations in V1 apply:
   - status token normalization
   - unambiguous link/path repairs
   - folder placement corrections with reversible patch artifacts
4. Add `hawp work-items fix-up --apply --plan <artifact>` using hash preconditions.
5. Implement refusal engine for ambiguity, missing evidence, and risky semantic edits.
6. Add optional AI draft mode behind explicit flag (for example `--ai-draft`) and keep it dry-run only in V1.
7. Defer `upgrade` and backlog-specialized commands to V1.1 after validating safety telemetry and reviewer workflow.

## Open Questions

1. Should plan artifacts be stored under `.hawp/work/status/YYYY/MM/DD/` or a dedicated `.hawp/work/plans/` path?
2. Should AI drafts be emitted as sidecar files only (`*.draft.md`) rather than inline patch hunks?
3. What minimum git evidence set should be accepted when users provide no CI/test logs?

## Notes for Implementation Hand-off

- Reuse validator finding IDs as stable operation anchors.
- Keep operation handlers pure and deterministic from findings + files.
- Treat AI drafting as a post-processing stage, never as a source parser.
- Maintain strict separation between direct evidence and inference in all generated output.
