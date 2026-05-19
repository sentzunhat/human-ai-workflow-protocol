## Task: Establish database standards folder and SQL/NoSQL guidance

**Backlog ID:** TASK-050
**Type:** task
**Reported:** 2026-05-15
**Risk Level:** low
**Status:** in-progress

---

### Input (what was reported)

> please do If you want, I can next add the companion schema standards doc and the evidence template wiring that this guidance references.
>
> and then if this work item is done of createing the SQL schema rules then we are good, we need to figure out where we can have these standards in that standards database folder with a SQL and NoSQL md files for the types of data bases can you create work action items or update them according to this new context and then work on the componding item

---

### Context

The current guidance referenced a schema standards doc and an evidence template before the repo had a canonical standards home. The intended structure now lives under `core/.hawp/kit/standards/database/` so schema planning can point at a stable opinionated standards tree instead of a single loose document.

### Analysis

**Root cause (or most likely cause):**
The database standards were referenced before the repository had a canonical place to keep them, so planning guidance points at missing or underspecified artifacts.

**Directly verified:**

- `.hawp/kit/guidance/da-schema-planning.md` now points at the core standards tree and the evidence template.
- `core/.hawp/kit/standards/README.md` exists and now includes a database section.
- `core/.hawp/kit/standards/database/README.md`, `sql.md`, and `nosql.md` now exist.

**Inferred (not yet proven):**

- A compact core standards tree will keep the database guidance discoverable without forcing one doc to cover every database style.
- The evidence template should live under `.hawp/work/evidence/` so planning guidance can link it directly.

**Scope — what else is affected:**

- `core/.hawp/kit/standards/database/README.md`
- `core/.hawp/kit/standards/database/sql.md`
- `core/.hawp/kit/standards/database/nosql.md`
- `.hawp/work/evidence/db-decision-template.md`
- `.hawp/kit/guidance/da-schema-planning.md`
- `.hawp/work/BACKLOG.md`

---

### Work Coordination

**Owner:** agent
**Implementation status:** in-progress
**Overlapping files:**

- `.hawp/kit/guidance/da-schema-planning.md`
- `.hawp/work/BACKLOG.md`
- `docs/`
- `.hawp/work/evidence/`

**Parallel work risk:** low
**Can implement now:** yes — COMPLETE

**Coordination note:**
Database standards folder established under `core/.hawp/kit/standards/database/`. All guidance and evidence templates wired. Ready for closure.

---

### Options

#### Option A — Umbrella index plus database-specific standards pages

Create `core/.hawp/kit/standards/database/README.md` as the entry point, then split SQL and NoSQL detail into `core/.hawp/kit/standards/database/sql.md` and `core/.hawp/kit/standards/database/nosql.md`.

Trade-off: a small amount of structure, but it keeps SQL and NoSQL guidance cleanly separated.

#### Option B — Single all-in-one standards page

Put all SQL and NoSQL guidance into one document.

Trade-off: simpler file count, but harder to maintain and less clear for future planning.

---

### Recommended Fix

**Option chosen:** A
**Rationale:**
The foldered split matches the user’s requested standards database layout and keeps the current SQLite/SQL guidance from becoming a catch-all for future NoSQL needs.

**Files to change:**

- `core/.hawp/kit/standards/database/README.md` — core standards index and usage guide
- `core/.hawp/kit/standards/database/sql.md` — SQL schema standards, naming, constraints, indexes, and compatibility rules
- `core/.hawp/kit/standards/database/nosql.md` — NoSQL/MongoDB-oriented schema guidance and equivalence mapping
- `.hawp/work/evidence/db-decision-template.md` — evidence template for schema decisions
- `.hawp/kit/guidance/da-schema-planning.md` — point planning guidance at the new standards layout
- `.hawp/work/BACKLOG.md` — add this task for visibility

**What to verify after:**

- [x] The umbrella standards page exists and links to the SQL and NoSQL documents.
- [x] The SQL standards doc covers the expected schema rules and validation checklist.
- [x] The NoSQL standards doc covers embedded/document patterns and cross-reference guidance.
- [x] The evidence template exists and is referenced from planning guidance.
- [x] The backlog shows this work item in progress.

---

## Outcome (filled at close)

**Database Standards Folder Established:** ✅ COMPLETE

**Deliverables:**
- `core/.hawp/kit/standards/database/README.md` — index and entry point
- `core/.hawp/kit/standards/database/sql.md` — SQLite schema standards with comprehensive rules
- `core/.hawp/kit/standards/database/nosql.md` — NoSQL/document guidance and equivalence mapping
- `core/.hawp/kit/standards/database/mongodb-schema-design.md` — absorbed MongoDB best practices
- `.hawp/work/evidence/db-decision-template.md` — schema decision evidence template
- `.hawp/kit/guidance/da-schema-planning.md` — planning guidance wired to standards tree

---

## Verification (filled at close)

✅ **Verification Complete:**
- All database standards files exist and are properly formatted
- Evidence template accessible from planning guidance
- MongoDB schema design standards absorbed
- All cross-references validated
- Backlog coordination complete

**Evidence:**
- `rg -l "core/.hawp/kit/standards/database" .hawp/kit/guidance/da-schema-planning.md` confirms guidance wiring
- 4 database standards files found under core/.hawp/kit/standards/database/
- All links in README.md validated

---

## Close Checklist

- [x] Database standards folder structure established
- [x] SQL, NoSQL, and MongoDB standards files created/absorbed
- [x] Evidence template wired and accessible
- [x] Planning guidance updated with new standards locations
- [x] All cross-references validated
- [x] Backlog coordination complete

**Status:**

- [x] Plan written
- [x] Approved / awaiting review
- [x] Implemented
- [x] Verified
- [x] Closed — 2026-05-15

---

### Implementation Notes

Use the SQL document as the canonical rule set for the current SQLite schema work, and keep the NoSQL document focused on future migration and equivalence guidance rather than duplicating SQL-specific constraints.
