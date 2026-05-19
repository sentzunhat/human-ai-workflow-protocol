## Task: Adapt repo-specific shared standards overlaps (docs and templates)

**Backlog ID:** TASK-062
**Type:** standards-update
**Reported:** 2026-05-15
**Risk Level:** low
**Status:** done

---

### Input (what was reported)

Derived from TASK-058 triage: several public shared files are useful but overlap repo-local references/templates and should be reviewed as adaptation candidates, not direct absorbs.

---

### Context

TASK-058 marked 25 files as `repo-specific` with `adapt` action, including:

- `shared_standards/public/standards/docs/hawp-install-update-safety.md`
- `shared_standards/public/templates/ADR.template.md`
- public governance/index/export artifacts

TASK-061 narrowed the Zacatl lane into adaptation candidates (not direct absorb):

- `shared_standards/public/standards/zacatl/service-boundaries.md`
- `shared_standards/public/standards/zacatl/handler-responsibilities.md`
- `shared_standards/public/standards/zacatl/dependency-registration.md`
- `shared_standards/public/standards/zacatl/layered-composition.md`
- `shared_standards/public/standards/zacatl/contract-testing.md`
- `shared_standards/public/standards/zacatl/evidence-linked-documentation.md`

---

### Analysis

**Root cause:**

Public shared artifacts are not uniformly absorbable because this repo already has local HAWP references and templates.

**Directly verified:**

- Existing local overlap paths:
  - `.hawp/kit/references/install-update-safety.md`
  - `.hawp/kit/templates/adr-template.md`
- Repo-specific adaptation bucket is documented in `.hawp/work/evidence/2026/05/15/TASK-058-shared-standards-triage.md`.

**Inferred (not yet proven):**

- Selective line-level adaptation can improve consistency without replacing repo-local conventions.

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `.hawp/kit/references/install-update-safety.md`
- `.hawp/kit/templates/adr-template.md`
- shared standards docs/template sources

**Parallel work risk:** low
**Can implement now:** yes (approved by continuation request)

---

### Recommended Fix

1. Diff shared and local versions for docs/template overlap candidates.
2. Propose minimal adaptation patch list (no large rewrites).
3. Apply only approved deltas and validate links/references.

**What to verify after:**

- [x] Overlap candidates are diff-reviewed.
- [x] Adapted content preserves repo-local boundaries and conventions.
- [x] Validation passes with no FAIL checks.

---

## Outcome

- Completed direct overlap diff review for:
  - `.hawp/kit/references/install-update-safety.md` vs `shared_standards/public/standards/docs/hawp-install-update-safety.md`
  - `.hawp/kit/templates/adr-template.md` vs `shared_standards/public/templates/ADR.template.md`
- Applied adaptation outcome as **reviewed, no-content-change** for both local files.
- Recorded evidence and rationale in `.hawp/work/evidence/2026/05/15/TASK-062-overlap-adaptation-review.md`.
- Split Zacatl extraction into a dedicated follow-up: `TASK-063`.

## Verification

Direct evidence:

- Repo-root proof captured pre-edit:
  - `pwd` => `<repo-root-abs>`
  - `git rev-parse --show-toplevel` => `<repo-root-abs>`
  - `git rev-parse --show-prefix` => `(empty)`
- Diff commands executed and reviewed:
  - `diff -u .hawp/kit/references/install-update-safety.md shared_standards/public/standards/docs/hawp-install-update-safety.md`
  - `diff -u .hawp/kit/templates/adr-template.md shared_standards/public/templates/ADR.template.md`
- Evidence file created:
  - `.hawp/work/evidence/2026/05/15/TASK-062-overlap-adaptation-review.md`

Unproven:

- Generalized extraction quality for Zacatl candidate docs remains pending `TASK-063`.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled (all claims have direct evidence or "unproven" tag)
- [x] Evidence files created if large/complex
- [x] Plan file moved to closed/YYYY/MM/DD/
- [x] BACKLOG.md updated
- [x] Status report written (not required)
- [x] Decision file created (not required)
