# Improvement: Add legacy tolerance mode to HAWP validator

**Backlog ID:** TASK-016
**Type:** improvement
**Reported:** 2026-05-10
**Risk Level:** low
**Status:** in-progress

---

### Context

Real HAWP repos (Tekit, Mictlan, and others) already use multiple historical
file-naming and directory layouts. The validator must be **strict for new work,
tolerant for old work, and aware of historical layouts**. It must not become a
"new HAWP shape or fail" tool.

---

### Known Historical Layouts

**Closed file naming (Tekit examples):**

- `2026-04-29-BUG-001-landing-page-component-split.md` — date-prefixed + ID + title
- `BUG-063-google-drive-provider-integration-review-and-plan.md` — ID + title (plan file)
- `BUG-091-status-report.md` — ID + supporting-file keyword (skip, not a plan file)
- `2026-05-02-backlog-archive-BUG-001-to-BUG-042.md` — archive range, no single ID (skip)

**Active file layout (Mictlan examples):**

- `active/TASK-003.md` — flat (current canonical form)
- `active/2026/05/06/TASK-011.md` — date-nested (legacy/alternate)

**Supporting files inside `closed/` (Mictlan examples):**

- `TASK-017-summary.md` — task summary, not a plan file (skip)
- `BACKLOG.done.archive.md` — archive file, not a plan file (skip)

---

### Design Rule

The validator must be:

- **Strict for new work** — closed files on/after `2026-05-10` must have Outcome, Verification, Close Checklist → missing = FAIL
- **Tolerant for old work** — closed files before `2026-05-10` → missing sections = WARN, not FAIL
- **Aware of historical layouts** — ID parsing and scanning handle all documented layouts

Constraints (out of scope):

- Do not migrate, rename, or rewrite old files
- Do not start UUID migration
- Do not add SQLite, indexing, search, queueing

---

### Implementation Plan

#### 1. `id-parser.ts` — Extend ID extraction

Add date-prefixed filename support to `extractIdFromFilename`:

- `2026-04-29-BUG-001-title` → `BUG-001`
- `BUG-063-title` → already works (standard prefix regex)
- `2026-05-02-backlog-archive-BUG-001-to-BUG-042` → null (no single ID at start after date → treated as supporting)

#### 2. `closed-task-completeness.ts` — File classification + legacy cutoff

**File classification:**

- **Supporting file** (skip, no sections required):
  - No extractable single ID → skip
  - Suffix after ID is exactly `-summary`, `-status-report`, or `-checkpoint` → skip
  - Suffix after ID contains `-archive` → skip
  - Filename starts with `BACKLOG` → skip
- **Plan file** (check sections): has a single extractable ID AND not a supporting file

**Legacy cutoff** (`LEGACY_CUTOFF = "2026-05-10"`):

- Date extracted from `closed/YYYY/MM/DD/` path
- Date < cutoff → WARN (legacy tolerance)
- Date >= cutoff → FAIL (strict)
- Unknown date → treated as legacy (safe fallback)

**Result shape:** replace `missing[]` with `failing[]` (FAIL) + `warnings[]` (WARN); add `skipped: number`

#### 3. `types.ts` — Update `ClosedTaskCheck`

Replace `missing` with `failing` + `warnings`; add `skipped`.

#### 4. `backlog-consistency.ts` — Nested active/ path support

Active file lookup tries `active/<ID>.md` first, then scans `active/YYYY/MM/DD/<ID>.md`.
Orphan scan covers both flat and nested directories.

#### 5. `reporter.ts` — Separate FAIL / WARN / skipped output

```
2. CLOSED TASK COMPLETENESS
Checking N plan file(s)  (M supporting files skipped):
	Outcome: X/N | Verification: Y/N | Close Checklist: Z/N

	[FAIL] Missing sections (2026-05-10 or later — must fix):
		TASK-XXX: ✗ missing Outcome, Verification

	[WARN] Legacy files missing sections (before 2026-05-10 — tolerated):
		BUG-001: missing Outcome, Verification, Close Checklist  (2026-04-29)
```

#### 6. `index.ts` — Include completeness WARN in warning count

---

### Verification Plan

1. `npm run typecheck` — 0 errors
2. `npm run validate:workflow` — FAIL count drops to 0; 16 legacy files become WARN
3. Document Tekit/Mictlan file classification
4. Save evidence at `evidence/2026/05/10/TASK-016-validation-output.md`

---

### Status

- [x] Plan written
- [x] Approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed

---

## Outcome

All six implementation goals delivered:

1. **ID parsing** — `extractIdFromFilename` now handles date-prefixed filenames (`2026-04-29-BUG-001-title` → `BUG-001`).
2. **File classification** — `isSupportingFile()` identifies status reports, summaries, checkpoints, archives, and BACKLOG files; those are skipped.
3. **Legacy cutoff** — `LEGACY_CUTOFF = "2026-05-10"`; files before this date produce `WARN`, files on/after produce `FAIL`.
4. **Nested active/ layout** — `findActiveFile()` checks both `active/<ID>.md` and `active/YYYY/MM/DD/<ID>.md`; `collectOrphanedActive()` covers both.
5. **Validator output** — Section 2 now shows `[FAIL]` / `[WARN]` / skipped counts separately with dates.
6. **Warning count** — `countWarnings` includes `completeness.status === "WARN"` so summary reflects tolerant issues.

Result: `VALIDATION PASS` — 0 FAIL, 1 WARN (completeness), 15 legacy files tolerated.

---

## Verification

- [x] `npm run typecheck` — 0 errors — **Evidence:** TASK-016-validation-output.md
- [x] `npm run validate:workflow` exit 0 — **Evidence:** TASK-016-validation-output.md
- [x] 15 pre-2026-05-10 files → `[WARN]` not `[FAIL]` — **Evidence:** validator output section 2
- [x] 0 items in `[FAIL]` — **Evidence:** validator output section 2
- [x] 1 supporting file correctly skipped — **Evidence:** "(1 supporting file(s) skipped)"
- [x] Tekit/Mictlan classification documented — **Evidence:** TASK-016-validation-output.md

---

## Close Checklist

- [x] Legacy cutoff implemented (`LEGACY_CUTOFF = "2026-05-10"`)
- [x] File classification implemented (supporting files skipped)
- [x] Date-prefixed ID parsing implemented
- [x] Nested `active/YYYY/MM/DD/` layout supported
- [x] Reporter shows `[FAIL]` / `[WARN]` / skipped separately
- [x] `countWarnings` includes completeness WARN
- [x] `npm run typecheck` passes
- [x] `npm run validate:workflow` returns PASS
- [x] Evidence saved at `evidence/2026/05/10/TASK-016-validation-output.md`
- [x] TASK-016 closed and moved to `closed/2026/05/10/`
- [x] BACKLOG.md updated
