## Bug / Task: Fix active main-branch HAWP alignment drift

**Backlog ID:** TASK-009
**Type:** task
**Reported:** 2026-05-02
**Risk Level:** low

---

### Input (what was reported)

> Fix active main-branch HAWP alignment drift.
>
> Required fixes:
> 1. Update README.md so repo-local work state is `.hawp/work/`, not `.work/`.
> 2. Remove or clearly mark root `.work/` as legacy/historical, or delete it if fully migrated.
> 3. Update core/install.md to copy `core/.hawp/kit/references/` into installed `.hawp/kit/references/`.
> 4. Update core/install.md boundary wording from `.work/` to `.hawp/work/`.
> 5. Update core/.hawp/work/BACKLOG.md scaffold so it uses compact active-index wording, not “single source of truth” + permanent Done.
> 6. Verify core/update.md and core/install.md match on kit references, work scaffold seeding, and `.hawp/work` preservation.
> 7. Search for stale active references:
>    - `.work/`
>    - `core/BACKLOG_ALIGNMENT.md`
>    - `BACKLOG_ALIGNMENT.md`
>    - `single source of truth`
>    - permanent `Done` sections
> 8. Do not add CLI, DB, runtime, or backlog compaction beyond scaffold/doc alignment.

---

### Context

This task is a documentation and scaffold alignment pass focused on current active guidance files, keeping source/downstream path boundaries consistent.

---

### Analysis

**Root cause (or most likely cause):**
Residual wording and drift checks are not fully normalized across README/install/update/scaffold docs.

**Directly verified:**
- Root `.work/` directory is absent.
- `README.md` already points repo-local work state to `.hawp/work/`.
- `core/install.md` script already copies `.hawp/kit/references`.
- `core/.hawp/work/BACKLOG.md` is compact-style but uses mixed path wording (`work/...`).
- Active references to `core/BACKLOG_ALIGNMENT.md` and `.work/` are largely in archived historical records.
- `single source of truth` appears in active `core/.hawp/kit/usage/intake-workflow.md` and docs-alignment prompts.

**Inferred (not yet proven):**
A few explicit wording updates will make branch alignment unambiguous and satisfy requested drift checks.

**Scope - what else is affected:**
- `README.md`
- `core/install.md`
- `core/update.md`
- `core/.hawp/work/BACKLOG.md`
- `core/.hawp/kit/usage/intake-workflow.md`
- `.hawp/work/BACKLOG.md` lifecycle row updates

---

### Work Coordination

**Owner:** agent
**Implementation status:** not-started
**Overlapping files:**

- `README.md`
- `core/install.md`
- `core/update.md`
- `core/.hawp/work/BACKLOG.md`
- `core/.hawp/kit/usage/intake-workflow.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
No current active plan overlaps these files.

---

### Options

#### Option A - Minimal alignment edits in active docs only

Apply small wording/path consistency fixes and verify stale references are limited to historical archives.
Trade-off: preserves historical records as-is.

#### Option B - Rewrite historical archives too

Normalize old references inside archived plans/evidence/status files.
Trade-off: rewrites history and expands scope beyond requested scaffold/doc alignment.

---

### Recommended Fix

**Option chosen:** A
**Rationale:**
Matches requested constraints and keeps scope tight.

**Files to change:**

- `README.md` - add explicit legacy-root `.work/` retirement note
- `core/install.md` - make references-copy and boundary wording explicit
- `core/update.md` - verify wording parity with install doc for references/scaffold/preservation
- `core/.hawp/work/BACKLOG.md` - ensure compact active-index wording and `.hawp/work` path style
- `core/.hawp/kit/usage/intake-workflow.md` - replace stale `single source of truth` phrase with active-index language

**What to verify after:**

- [ ] Active docs/scaffolds have no stale `.work/` or `core/BACKLOG_ALIGNMENT.md` references
- [ ] install/update wording matches on kit references, scaffold seeding, and `.hawp/work` preservation
- [ ] No CLI/DB/runtime additions and no backlog compaction work added

---

### Implementation Notes

Keep edits minimal and textual only. Treat historical archive files as immutable records.

### Outcome

- Updated active docs and scaffold wording to reinforce `.hawp/work/` as the repo-local work state.
- Added an explicit README note marking legacy root `.work/` as retired/historical.
- Confirmed `core/install.md` and `core/update.md` both include kit references copy and matching `.hawp/work` scaffold/preservation semantics.
- Updated scaffold and intake workflow wording to use active-index language.
- Searched for stale references and confirmed remaining `.work/` and `core/BACKLOG_ALIGNMENT.md` hits are in archived historical records.

### Verification

- `rg -n --hidden -uu "\.work/|core/BACKLOG_ALIGNMENT\.md|BACKLOG_ALIGNMENT\.md|single source of truth|## Done|permanent history" README.md core .github .hawp/work core/.hawp/work .hawp/kit`
- `rg -n "kit/references|copy_file_no_clobber \"\$SRC/\.hawp/work|cp -R \"\$SRC/\.hawp/kit/references\"|seed \.hawp/work/ scaffold|project-owned" core/install.md core/update.md`
- `ls -la` (root) confirmed no root `.work/` directory

---

### Status

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
