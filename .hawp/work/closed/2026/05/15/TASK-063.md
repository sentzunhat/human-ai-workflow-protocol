## Task: Extract generalized standards from Zacatl adaptation candidates

**Backlog ID:** TASK-063
**Type:** standards-update
**Reported:** 2026-05-15
**Risk Level:** medium
**Status:** done

---

### Input (what was reported)

Derived from TASK-061/TASK-062: six Zacatl principle-level docs are adaptation candidates but should not be absorbed directly under Zacatl labeling.

---

### Context

TASK-061 classified six Zacatl files as `repo-specific` adaptation candidates and one README as private-lane. TASK-062 kept docs/template overlap unchanged for boundary safety and split Zacatl extraction as a dedicated follow-up.

Candidate sources:

- `shared_standards/public/standards/zacatl/service-boundaries.md`
- `shared_standards/public/standards/zacatl/handler-responsibilities.md`
- `shared_standards/public/standards/zacatl/dependency-registration.md`
- `shared_standards/public/standards/zacatl/layered-composition.md`
- `shared_standards/public/standards/zacatl/contract-testing.md`
- `shared_standards/public/standards/zacatl/evidence-linked-documentation.md`

---

### Analysis

**Root cause:**

Useful principle-level guidance is currently tied to framework/domain naming and should be extracted into neutral standards language.

**Directly verified:**

- Six files contain principle-level content and explicit "does not include" boundaries.
- Direct absorb is not approved in current form.

**Inferred (not yet proven):**

- Generalized extraction can improve reuse without leaking internal domain references.

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `core/.hawp/kit/standards/**`
- `.hawp/work/active/TASK-054.md`
- `.hawp/work/active/TASK-062.md` (source decision context)

**Parallel work risk:** medium
**Can implement now:** yes (approved by continuation request)

---

### Recommended Fix

1. Draft neutralized standards pages (remove framework-specific labels and internal-domain references).
2. Preserve only principle-level guidance and boundary-safe sections.
3. Add references in `core/.hawp/kit/standards/README.md` only after privacy scan and link validation.

**What to verify after:**

- [x] No internal-domain references remain in extracted text.
- [x] New pages are principle-level and reusable.
- [x] Validation passes with no FAIL checks.

---

## Outcome

- Extracted six neutral service-design standards from Zacatl adaptation candidates into:
	- `core/.hawp/kit/standards/service-design/service-boundaries.md`
	- `core/.hawp/kit/standards/service-design/handler-responsibilities.md`
	- `core/.hawp/kit/standards/service-design/dependency-composition.md`
	- `core/.hawp/kit/standards/service-design/layered-composition.md`
	- `core/.hawp/kit/standards/service-design/contract-testing.md`
	- `core/.hawp/kit/standards/service-design/evidence-linked-documentation.md`
- Added index page: `core/.hawp/kit/standards/service-design/README.md`.
- Updated `core/.hawp/kit/standards/README.md` to include service-design section and canonical map entries.
- Recorded extraction evidence at `.hawp/work/evidence/2026/05/15/TASK-063-generalized-extraction-review.md`.

## Verification

Direct evidence:

- Source and destination paths are documented in `.hawp/work/evidence/2026/05/15/TASK-063-generalized-extraction-review.md`.
- Privacy scan terms checked against extracted files:
	- `tekit|mictlan|micltan|zacatl|internal-only`
- Workflow validator run after updates reports PASS with no FAIL checks.

Unproven:

- Real-world adoption clarity of these generalized pages remains to be evaluated in future standards-review cycles.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled (all claims have direct evidence or "unproven" tag)
- [x] Evidence files created if large/complex
- [x] Plan file moved to closed/YYYY/MM/DD/
- [x] BACKLOG.md updated
- [x] Status report written (not required)
- [x] Decision file created (not required)
