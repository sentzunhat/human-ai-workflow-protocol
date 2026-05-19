# TASK-047 Dry-Run Classification Evidence

Date: 2026-05-14
Mode: dry-run (no mutations)

## Command

- `npm --prefix librarian run validate:workflow`

Raw output evidence:

- `TASK-047-validate-workflow-output.txt`

## Classification Summary (Policy C)

### Fail (must remediate)

1. Active backlog consistency orphans

- `active/TASK-027.md`
- `active/TASK-040.md`
  Reason: Active Work table/file linkage must be deterministic for validator integrity.
  Policy class: parser/structure contract.

2. Untyped closed files post cutoff (`2026-05-10+`)

- `.hawp/work/closed/2026/05/12/0007.md`
- `.hawp/work/closed/2026/05/12/0008-CLARIFICATION-exact-paths.md`
- `.hawp/work/closed/2026/05/12/0008-install-update-distribution-review.md`
- `.hawp/work/closed/2026/05/12/HAWP-BACKLOG-VALIDATE-PLAN.md`
- `.hawp/work/closed/2026/05/13/0007.md`
- `.hawp/work/closed/2026/05/13/HAWP-BACKLOG-VALIDATE-PLAN.md`
  Reason: typed-ID contract is mandatory after cutoff.
  Policy class: fail.

3. Missing required sections in post-cutoff tracked files

- `.hawp/work/closed/2026/05/12/TASK-030-files.md` (missing Outcome, Close Checklist)
- `.hawp/work/closed/2026/05/13/TASK-030-files.md` (missing Outcome, Close Checklist)
- `.hawp/work/closed/2026/05/13/TASK-033.md` (missing Close Checklist)
- `.hawp/work/closed/2026/05/13/TASK-038.md` (missing Outcome, Close Checklist)
- `.hawp/work/closed/2026/05/14/TASK-044.md` (missing Outcome, Close Checklist)
  Reason: required closed-record completeness contract.
  Policy class: fail.

### Warn (visible, non-blocking by policy)

1. Legacy untyped files before cutoff

- `.hawp/work/closed/2026/04/26/2026-04-26-hawp-adr-template-review.md`

2. Legacy files missing sections before cutoff

- Multiple entries under `.hawp/work/closed/2026/04/30/` to `.hawp/work/closed/2026/05/03/`.
  Reason: historical tolerance band.

3. Verification clarity unproven note

- `TASK-036` unproven GitHub-hosted workflow verification.
  Reason: does not break parser/structure contracts.

## Dry-Run Safety Confirmation

- No files were modified by validator execution.
- No rename/move/delete actions were applied.
- This report is classification only.

## Next Action (Apply Phase Candidate)

- Use additive-only remediation in TASK-047 for fail-class findings.
- Preserve historical narrative content.
- Create backup snapshot and manifest before any apply-mode edits.
