## Task: Audit and consolidate opinionated standards folder layout

**Backlog ID:** TASK-053
**Type:** task
**Reported:** 2026-05-15
**Risk Level:** medium
**Status:** in-progress

---

### Input (what was reported)

> an action item to review and plan an implementation of getting all the opinionates standards folders in the stadards folder which are opinionated in code and arch and anything else related to standards

---

### Context

New database standards now live under `core/.hawp/kit/standards/database/`, but related standards-like guidance still spans the repo. The next planning step is to classify what belongs in the core standards tree, what stays as supporting guidance, and what should be split into smaller follow-up items.

---

### Analysis

**Root cause (or most likely cause):**
Standards content has evolved incrementally and is currently distributed across multiple folders, making discovery and governance harder.

**Directly verified:**

- `core/.hawp/kit/standards/database/README.md`, `sql.md`, and `nosql.md` now exist.
- `.hawp/kit/guidance/da-schema-planning.md` now references the core standards tree.
- TASK-054 now exists to review shared public-safe standards before any larger absorption pass.

**Inferred (not yet proven):**

- Consolidating standards navigation under a coherent `docs/standards/` map will reduce drift and duplicate guidance.

**Scope — what else is affected:**

- `core/.hawp/kit/standards/**`
- `.hawp/kit/guidance/**`
- `README.md` links to standards entry points
- `shared_standards/**` once provided

---

### Work Coordination

**Owner:** agent
**Implementation status:** not-started
**Overlapping files:**

- `docs/`
- `.hawp/kit/guidance/`
- `README.md`

**Parallel work risk:** medium
**Can implement now:** only after plan review

---

### Recommended Fix

Phase 1: inventory all standards-like docs and classify canonical vs supporting.
Phase 2: propose target folder map and redirect/link migration plan.
Phase 3: implement in small safe moves with reference validation after each pass.

---

### Progress Update (2026-05-15)

Kickoff inventory completed for standards-related docs/guidance paths:

- `core/.hawp/kit/standards/database/README.md`
- `core/.hawp/kit/standards/database/sql.md`
- `core/.hawp/kit/standards/database/nosql.md`
- `.hawp/kit/guidance/da-schema-planning.md`
- `.hawp/kit/references/backlog-alignment.md`
- `.hawp/kit/references/docs-alignment.md`
- `.hawp/kit/references/install-update-safety.md`
- `.hawp/kit/references/work-item-file-tracking.md`
- `.hawp/work/evidence/db-decision-template.md`

Next move for TASK-053: classify these into canonical standards vs supporting references, then produce a no-break migration/link plan.

### Progress Update (2026-05-15, phase-2)

Classification pass completed and codified in `core/.hawp/kit/standards/README.md`:

- Canonical standards (normative):
  - `core/.hawp/kit/standards/guidelines/*.md`
  - `core/.hawp/kit/standards/nodejs/*.md`
  - `core/.hawp/kit/standards/database/sql.md`
  - `core/.hawp/kit/standards/database/nosql.md`
- Supporting (non-normative) references:
  - `core/.hawp/kit/standards/database/README.md`
  - `.hawp/kit/guidance/da-schema-planning.md`
  - `.hawp/work/evidence/db-decision-template.md`

No-break migration/link plan status:

1. Keep standards docs in `core/.hawp/kit/standards/**` as canonical.
2. Keep execution/planning guidance in `.hawp/kit/guidance/**`.
3. Keep evidence templates in `.hawp/work/evidence/**`.
4. Enforce reference safety with markdown-link audits after each standards update.

Remaining for TASK-053:

- Coordinate with TASK-054 once shared standards source content is available.
- Decide whether additional standards domains (for example, architecture overlays beyond current guideline set) require new folders or remain under `guidelines/`.

---

## Outcome (filled at close)

Standards layout audit and consolidation was completed.

Completed outcomes:
- Inventory and classification pass finished for standards-like assets.
- Canonical vs supporting boundaries codified in `core/.hawp/kit/standards/README.md`.
- No-break migration/link strategy established and applied.
- Follow-on shared standards work completed through TASK-054/TASK-063 with boundaries preserved.

## Verification (filled at close)

Directly verified:
- Canonical standards sections now include `guidelines/`, `nodejs/`, `database/`, `patterns/`, and `service-design/` in `core/.hawp/kit/standards/README.md`.
- Supporting guidance separation remains explicit (`.hawp/kit/guidance/**`, `.hawp/work/evidence/**`).
- Backlog/workflow validation passes with no FAIL checks.

Confidence: high. Consolidation goals for this item are met.

## Close Checklist

- [x] Standards-like docs inventoried
- [x] Canonical vs supporting classification documented
- [x] No-break migration/link plan applied
- [x] Follow-on decomposition and coordination completed
- [x] Validation pass confirmed

**Status:**

- [x] Plan written
- [x] Approved / awaiting review
- [x] Implemented
- [x] Verified
- [x] Closed
