# Task: Update and align SQLite schema standards to new canonical version

**Backlog ID:** TASK-055
**Type:** standards-update
**Reported:** 2026-05-15
**Risk Level:** medium
**Status:** analyzing

---

### Input (what was reported)

> New update on the SQLite schema design standards was provided. The standards in `core/.hawp/kit/standards/database/sql.md` and related references must be updated to match the new canonical version. This should be planned, reviewed, and implemented before absorbing further shared standards.

---

### Context

- The user provided a detailed, updated SQLite schema standards document.
- The canonical standards file (`core/.hawp/kit/standards/database/sql.md`) and any references must be updated to match this new version.
- This update should be planned and reviewed before implementation, in line with HAWP workflow.
- Absorption of shared standards (TASK-054) should be paused until this update is complete.

---

### Analysis

**Root cause:**

- The standards file is out of date with the latest canonical guidance provided by the user.

**Directly verified:**

- User provided a new, detailed SQLite schema standards document.
- Existing standards file is less detailed and missing several new rules and examples.

**Inferred (not yet proven):**

- Updating the standards file first will prevent drift and ensure all future absorption work is based on the correct baseline.

**Scope — what else is affected:**

- `core/.hawp/kit/standards/database/sql.md`
- Any guidance or evidence templates referencing the standards
- The review/absorption process for shared standards (TASK-054)

---

### Work Coordination

**Owner:** agent
**Implementation status:** not-started
**Overlapping files:**

- `core/.hawp/kit/standards/database/sql.md`
- `.hawp/kit/guidance/da-schema-planning.md`
- `.hawp/work/evidence/db-decision-template.md`

**Parallel work risk:** medium
**Can implement now:** only after plan review

---

### Recommended Fix

Phase 1: Plan and review the update to the SQLite standards file.
Phase 2: Implement the update, ensuring all new rules, examples, and validation checklists are included.
Phase 3: Validate all references and update related guidance/templates as needed.
Phase 4: Resume shared standards absorption (TASK-054) once the baseline is current.

---

### Progress Update (2026-05-15)

- Intake plan created. Pending review and approval before implementation.
- Shared standards absorption (TASK-054) is paused until this update is complete.

---

## Outcome (filled at close)

SQLite schema standards were aligned to the new canonical baseline.

Completed deliverables:
- Updated `core/.hawp/kit/standards/database/sql.md` with comprehensive naming rules, FK conventions, canonical entity relationship model, index guidance, and validation checklist.
- Confirmed standards index alignment in `core/.hawp/kit/standards/database/README.md`.
- Confirmed evidence workflow alignment in `.hawp/work/evidence/db-decision-template.md` and planning references.
- Unblocked downstream shared-standards absorption work.

## Verification (filled at close)

Directly verified:
- `core/.hawp/kit/standards/database/sql.md` now contains the expected canonical sections: naming conventions, PK/FK rules, timestamps, enum/check guidance, canonical table model, index guidelines, and validation checklist.
- `core/.hawp/kit/standards/database/README.md` references `sql.md`, `nosql.md`, and absorbed MongoDB guidance.
- `.hawp/work/evidence/db-decision-template.md` references the canonical database standards location.

Confidence: high. The standards baseline is current and usable for follow-on review/absorption work.

## Close Checklist

- [x] Canonical SQL standards file updated to new baseline
- [x] Related standards references validated
- [x] Evidence template wiring validated
- [x] TASK-054 dependency removed
- [x] Work item moved to closed

**Status:**

- [x] Plan written
- [x] Approved / awaiting review
- [x] Implemented
- [x] Verified
- [x] Closed
