# Work Intake - Plan

## Bug / Task: Validate downstream HAWP install compatibility

**Backlog ID:** TASK-019  
**Type:** task  
**Reported:** 2026-05-10  
**Risk Level:** medium

---

### Input (what was reported)

> Run the HAWP workflow validator against real downstream HAWP installs for compatibility verification.
>
> Context:
> The validator currently lives in the HAWP repo under:
>
> librarian/scripts/validate-hawp-workflow/
>
> It was developed against the HAWP repo's own `.hawp/work` folder. Now we need to verify it against real downstream project histories with older HAWP layouts.
>
> Target downstream HAWP folders:
>
> tekit/.hawp
> mictlan/.hawp
>
> Mission:
> Run the validator from the HAWP repo against Tekit and Mictlan, without modifying either downstream repo.
>
> Goals:
>
> 1. Confirm the validator can read external `.hawp` roots.
> 2. Confirm backwards compatibility with older layouts:
>    - date-prefixed closed files
>    - title-suffixed BUG/TASK files
>    - active files directly under active/
>    - active files nested under active/YYYY/MM/DD/
>    - status/evidence/summary/archive support files
>    - legacy files missing modern close checklist sections
> 3. Report whether validation passes, warns, or fails for each target.
> 4. Save the results as evidence in the HAWP repo.
>
> Scope:
>
> - Add a CLI argument if needed, such as:
>   - `--hawp-root /path/to/.hawp`
>     or
>   - `--work-root /path/to/.hawp/work`
> - Keep current default behavior working for the HAWP repo itself.
> - Do not modify Tekit or Mictlan files.
> - Do not migrate, rename, move, or rewrite any downstream work files.
> - Do not start UUID migration.
> - Do not add SQLite, indexing, search, or queueing.

---

### Context

Validator compatibility needs to be verified against real downstream `.hawp` histories from other repos while running from this repository's `librarian` package. The validator must remain read-only for downstream projects.

### Analysis

**Root cause (or most likely cause):**  
The validator resolves `.hawp/work` by walking up from `cwd` and currently has no explicit external-root CLI option. Some checks may also be too strict for legacy support-file naming patterns.

**Directly verified:**

- Current CLI uses `findWorkDirectory()` with ancestor search and no arg parsing.
- Closed-file completeness check already has a legacy cutoff and warning mode before `2026-05-10`.
- Active-file lookup supports flat and date-nested active layouts.

**Inferred (not yet proven):**

- None.

**Scope - what else is affected:**

- `librarian/scripts/validate-hawp-workflow/*` only.
- HAWP evidence and backlog files for task tracking.

### Work Coordination

**Owner:** agent  
**Implementation status:** done  
**Overlapping files:**

- `.hawp/work/BACKLOG.md`
- `librarian/scripts/validate-hawp-workflow/index.ts`
- `librarian/scripts/validate-hawp-workflow/validations/closed-task-completeness.ts`

**Parallel work risk:** low  
**Can implement now:** yes

**Coordination note:**  
No active work items are currently listed in backlog Active Work beyond this task.

### Options

#### Option A - Add explicit external root CLI resolution

Parse `--hawp-root` and `--work-root` in the entrypoint, preserve current fallback discovery, and keep validators unchanged except compatibility-safe support-file classification improvements.

#### Option B - Keep discovery logic and invoke from target directories

Run validator with `cwd` in each downstream repo and absolute script path, without adding CLI options.

### Recommended Fix

**Option chosen:** A  
**Rationale:** Explicit root arguments provide deterministic external validation from this repo and reduce operator error while preserving backward-compatible defaults.

**Files to change:**

- `librarian/scripts/validate-hawp-workflow/index.ts` - add CLI args for external `.hawp` or work root.
- `librarian/scripts/validate-hawp-workflow/validations/closed-task-completeness.ts` - broaden clearly-supporting legacy suffix patterns if needed.
- `.hawp/work/evidence/2026/05/10/TASK-019-*.md` - store Tekit/Mictlan outputs.

**What to verify after:**

- [x] `npm run typecheck` passes in `librarian`.
- [x] validator runs against Tekit `.hawp` without file modifications.
- [x] validator runs against Mictlan `.hawp` without file modifications.
- [x] evidence files saved in this repo under requested path.

### Implementation Notes

Keep all downstream paths read-only and avoid any writes outside this repository.

---

## Outcome (filled at close)

Implemented compatibility-focused validator changes and ran downstream checks from this repo:

- Added external root CLI support:
  - `--hawp-root /path/to/.hawp`
  - `--work-root /path/to/.hawp/work`
  - preserved existing default auto-discovery for local runs
- Added `## Done` section support when parsing closed backlog rows.
- Fixed case-insensitive ID matching for date-prefixed/title-suffixed file discovery.
- Expanded support-file suffix handling for clear non-plan files (`-status`, `-status-report`, `-summary`, `-checkpoint`, `-evidence`, including suffix-ending variants).
- Executed validator against Tekit and Mictlan external roots.
- Saved required evidence files under `.hawp/work/evidence/2026/05/10/`.

## Verification (filled at close)

- [x] `npm run typecheck` succeeded in `librarian`. **Evidence:** command completed with no TypeScript errors.
- [x] Tekit external validation executed successfully (read-only). **Evidence:** `.hawp/work/evidence/2026/05/10/TASK-019-tekit-validator-compatibility.md`.
- [x] Mictlan external validation executed successfully (read-only). **Evidence:** `.hawp/work/evidence/2026/05/10/TASK-019-mictlan-validator-compatibility.md`.
- [x] Default no-arg behavior still works. **Evidence:** `npm run validate:workflow` exited `0` and validated local `.hawp/work`.

## Close Checklist

**Before marking this task done in BACKLOG.md, verify ALL:**

- [x] Outcome section filled (what was actually implemented)
- [x] Verification section filled (all checks listed, each with direct evidence or "unproven" tag)
- [x] All evidence files referenced exist in `../evidence/YYYY/MM/DD/` or are noted as inline
- [x] Plan file moved to `../closed/YYYY/MM/DD/<ID>.md`
- [x] BACKLOG.md row moved from Active to Recently Closed (or marked done)
- [x] Status report written (optional: only if non-trivial OR if something remains unproven OR if a decision/pattern emerged)
- [ ] Decision file created if applicable (only if this task resolves a design question)

**Status:**

- [x] Plan written
- [x] Approved / user-request approved (medium risk)
- [x] Implemented
- [x] Verified
- [x] Closed
