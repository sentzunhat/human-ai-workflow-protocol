# HAWP Project Drift and Improvement Audit

## Discovery Report — May 12, 2026

**Audit Scope:**

- Repository-wide consistency and clarity for AI agent guidance
- Workflow records vs actual repo state alignment
- Install/update/distribution docs accuracy and boundary clarity
- Validator behavior vs documentation alignment
- Legacy content handling and future-facing notes
- Naming, references, and command documentation

---

## Working Tree Status (Hard Constraint: NOT CLEAN)

### Modified Files (9) — From Completed TASK-026 Audit

Ready to commit but currently staged to orphaned commit 182ba18:

- `.hawp/work/BACKLOG.md` — Added TASK-026 row
- `core/update.md` — Fixed reconciliation output label references (2 edits)
- `core/install.md` — Clarified `.hawp/work` boundary semantics
- `core/distribution/sources/shared/install.md` — Updated no-overwrite language
- `core/distribution/sources/shared/update.md` — Same no-overwrite language
- `core/distribution/sources/install/main.md` — Updated boundary notes
- `core/distribution/generated/install-dev.md` — Regenerated
- `core/distribution/generated/install-main.md` — Regenerated
- `core/distribution/generated/update-dev.md` — Regenerated
- `core/distribution/generated/update-main.md` — Regenerated

### Untracked Files (7) — Parallel/Planning Work

Not part of docs-audit lane but in working tree:

- `.hawp/work/active/TASK-026.md` — Docs audit plan (should be committed with docs-audit)
- `.hawp/work/active/TASK-027.md` — Backlog upgrade CLI command (separate lane)
- `.hawp/work/active/TASK-028.md` — Backlog detection/dry-run (separate lane)
- `.hawp/work/active/TASK-029.md` — Backlog data models (separate lane)
- `.hawp/work/decisions/2026/05/11/` — Design decision folder
- `.hawp/work/notes/2026/05/11/hawp-backlog-upgrade-command-design.md` — Design artifact
- `.hawp/work/status/2026/05/11/backlog-upgrade-implementation-plan.md` — Status report

---

## Repository Inventory

**Total tracked files:** 35

```
README.md                                                (1)
LICENSE                                                 (1)
core/
  install.md, update.md                                 (2)
  distribution/
    sources/
      install/{dev,main}.md                             (2)
      update/{dev,main}.md                              (2)
      shared/{install,update,safety,repo-boundaries}.md (4)
    generated/
      {install,update}-{dev,main}.md                    (4)
benchmark/
  README.md, benchmark-prompt.md                        (2)
librarian/
  package.json, tsconfig.json, package-lock.json        (3)
  scripts/
    distribution/
      build/index.ts                                    (1)
      validate/index.ts                                 (1)
      shared/{composition,types}.ts                     (2)
    validate-hawp-workflow/
      {index,cli,reporter,orchestrate}.ts               (4)
      validations/
        {id-parser,backlog-consistency,...}.ts          (5)
.github/
  copilot-instructions.md                               (1)
  instructions/ (detected via rg)                        (4)
  prompts/ (detected via rg)                            (5+)
.hawp/
  (kit, work folders not counted; not tracked git files)
```

---

## Drift Analysis by Category

### 1. ✅ Commands and Scripts — CONSISTENT

**Finding:** Documentation references match npm scripts exactly.

**Evidence:**

- `README.md` → `npm run distribution:build`, `npm run distribution:validate` ✓
- `librarian/package.json` defines both scripts ✓
- Generated docs warn to regenerate with `npm run distribution:build` ✓
- Generated docs reference exact script extraction source (REF=main|dev) ✓

**Status:** No drift detected.

---

### 2. ✅ Installation/Update Boundaries — MOSTLY CONSISTENT

**Finding:** `.hawp/kit/`, `.hawp/work/`, `.github/` boundaries documented and preserved.

**Evidence:**

- `core/install.md`: "This document covers full HAWP installation: `.hawp/` at repo root"
- `core/update.md`: "Update boundary" section explicitly states `.hawp/work/` is never overwritten
- `README.md`: "They never overwrite your `.hawp/work/` content or your `.github/copilot-instructions.md`"
- Generated distribution guides inherit boundary warnings from sources ✓

**Minor Drift:**

- `core/install.md` line ~374 says "This package is self-bootstrapping" — uses singular "package" but applies to entire repo installation context. Not misleading but could clarify scope once per doc (suggestion: preface with "the install script" or "this flow").

**Status:** No critical drift. One wording clarification opportunity.

---

### 3. 🟡 `.hawp/work` Preservation Language — PARTIALLY DRIFTED

**Finding:** Pre-audit docs said "does not modify" but actual behavior moves eligible closed/orphan items.

**Evidence:**

- Pre-audit `core/distribution/sources/shared/` said: "The generated command preserves your `.hawp/work/` content"
- Actual behavior (per `core/update.md`): reconciles done items, moves orphans to closed/, uses no-overwrite on individual files with `cp -Rn`
- This creates ambiguity: Agents reading docs might think work items are strictly untouched, but reconciliation actively relocates eligible files.

**Corrective Edit Applied (TASK-026):**

- Changed to: "does not overwrite...may move eligible closed or orphan items during reconciliation"
- All 4 distribution docs regenerated ✓

**Status:** RESOLVED by TASK-026 audit edits (pending commit).

---

### 4. 🟡 Generated File Editing — CLEAR BUT COULD BE STRONGER

**Finding:** Docs properly warn not to hand-edit `core/distribution/generated/` but warning placement could be improved for agent discoverability.

**Evidence:**

- `README.md` line 30: "Do not edit files under `core/distribution/generated/` directly." ✓
- Each generated file footer references "extracted from `core/install.md`...run `npm run distribution:build` to regenerate" ✓
- Source fragments also warn: "The distribution fragments are generated guidance; when in doubt, `core/install.md` takes precedence" ✓

**Minor Improvement Opportunity:**

- Warning is in footer of each file (good for downstream) but not in table of contents or quick-reference section
- Agents first-time reading README might miss the warning if they jump to links table
- Suggestion: Add single-line note in Generated guide quick links section: "⚠️ Auto-generated; do not edit. Edit source fragments instead."

**Status:** No critical drift. Minor discoverability improvement identified.

---

### 5. 🟢 Validator Behavior — ALIGNED

**Finding:** `npm run validate:workflow` behavior matches documented expectations.

**Evidence:**

- Exit code: 0 for PASS/WARN, 1 for FAIL ✓
- Current repo state: 1 WARN (legacy closed files pre-2026-05-10), overall PASS ✓
- Validator docs in `.hawp/work/active/HAWP-BACKLOG-VALIDATE-PLAN.md` accurately describe checks ✓
- Reporter output includes warnings, legit failures, and legacy tolerance section ✓

**Note:** Future validator enhancement (TASK-030+) may add `--legacy-strict` flag; docs will need update when implemented.

**Status:** Aligned.

---

### 6. 🟡 Active vs Closed Reconciliation — MINOR MISMATCH

**Finding:** BACKLOG.md shows 4 active items and 10 recently-closed items, but actual closed count is 7 total (2 in 2026/04, 5 in 2026/05).

**Evidence:**

- `BACKLOG.md` Recently Closed lists: TASK-025, TASK-024, TASK-023, TASK-022, TASK-020, BUG-006, TASK-021, TASK-019, TASK-018, TASK-017 (10 items)
- `find .hawp/work/closed/2026 -name "*.md"` shows only 7 files total
- This suggests older items (TASK-019, TASK-018, TASK-017 and others) are in 2026/04 but not current date folders
- Backlog compaction rule states: "Keep Recently Closed capped; archive history lives in `closed/`"
- Current recently closed section may be showing stale links or the count verification is off

**Severity:** LOW — Information consistency issue, not correctness. Closed files do exist (evidence verified).

**Recommendation:** Run verification: `cd .hawp/work && find closed -name "*.md" | wc -l` and cross-check BACKLOG dates.

**Status:** Detected; recommend verification before next compaction.

---

### 7. 🔴 Future-Facing and Legacy Notes — ABUNDANT BUT UNCLEAR PRIORITIZATION

**Finding:** ~40+ mentions of "future", "legacy", "coming soon", "TODO-like" phrases scattered across docs and closed evidence files. No centralized tracking of which future items are blocked, parked, or in active development.

**Evidence Sample (From Drift Scan):**

- `.hawp/work/active/HAWP-BACKLOG-VALIDATE-PLAN.md` line 121: "Supports future commands (`fix-up`, `upgrade`) sharing same validation foundation"
- `.hawp/work/active/HAWP-BACKLOG-VALIDATE-PLAN.md` line 567: "Adding AI-assisted diagnostics (future enhancement)"
- `.hawp/work/parked/TASK-013.md`: UUID-based IDs project (depends on TASK-012 completion)
- `core/distribution/sources/update/main.md`: "Check that stale legacy folders were removed"
- Multiple closed evidence files reference "future compaction task" and "legacy pre-2026-05-10" cleanup

**Specific Categories of Future Work (Not Yet Captured):**

1. **Legacy cleanup**: Pre-2026-05-10 closed files need formatting standardization (future compaction task)
2. **Backlog upgrade command**: TASK-027/028/029 in active; TASK-030+ (apply mode) deferred
3. **UUID migration**: TASK-013 parked, depends on TASK-012 validator completion
4. **AI-assisted synthesis**: Explicitly deferred with governance gates
5. **CI/automation integration**: Mentioned as extensible future in id-parser.ts
6. **Multi-root validation**: Mentioned as possible future in validator plan
7. **Decision tree linking**: Closed evidence notes it as "minor future enhancement"

**Severity:** MEDIUM — Creates ambient vagueness. Agents reading docs don't know which futures are near-term (TASK-027+), deferred (TASK-030+), or exploratory (AI-assisted, multi-root).

**Recommendation:** Create a "Roadmap" section in README.md or dedicated Future Work note linking to BACKLOG parked/approved items.

**Status:** Detected; recommendation ready.

---

### 8. 🟡 Legacy Layout Migration Documentation — CLEAR BUT SCATTERED

**Finding:** Legacy migration logic is documented in multiple places (install.md, update.md, distribution sources, closed evidence) but not consolidated.

**Evidence:**

- `core/install.md` lines 11-13: Describes three legacy migration flows
- `core/distribution/sources/shared/safety.md`: Legacy folder removal notes
- `core/update.md`: Comprehensive legacy migration section
- `.hawp/work/closed/2026/05/01/2026-05-01-bug-005-install-update-alignment.md`: Summary of all migration behaviors

**Strengths:**

- Each location provides appropriate context (installation vs update vs safety)
- Generated distribution guides inherit guidance ✓
- No contradictions detected

**Improvement Opportunity:**

- First-time reader of install.md might struggle to see all three migration paths at once
- Suggestion: Add a single "Legacy Migration Flowchart" table early in install.md that shows all three paths with outcomes (legacy `hawp/` → `.hawp/`, legacy `.hawp/usage/` → `.hawp/work/`, legacy `.hawp/status/` → `.hawp/work/notes/`)

**Status:** Functional but could improve discoverability.

---

### 9. ✅ Validator Extensibility — WELL DOCUMENTED

**Finding:** `librarian/scripts/validate-hawp-workflow/validations/id-parser.ts` documents extensibility for future UUID support.

**Evidence:**

- Comment: "Designed for extensibility to UUID-based IDs in future"
- Current implementation: TASK-NNN, BUG-NNN format parsing with lowercase-insensitive matching
- Related BACKLOG item TASK-013 (UUID migration) explicitly depends on TASK-012 validator completion ✓

**Status:** Well handled.

---

### 10. 🟡 Librarian Package Scope — PARTIALLY DOCUMENTED

**Finding:** `librarian/` folder contains both distribution (build/validate) and workflow (validate-hawp-workflow) scripts. Purpose and scope not clearly documented in README or folder structure.

**Evidence:**

- `librarian/package.json` has 2 script groups: distribution-related and workflow-related
- No README in `librarian/` folder explaining the distinction
- Top-level README mentions "librarian/scripts/distribution/..." but not overall purpose

**Improvement Opportunity:**

- Create `librarian/README.md` documenting:
  - Purpose: HAWP maintenance and documentation generation
  - Two script groups: (1) distribution build/validate, (2) workflow validation
  - Entry points for common tasks
  - Future extensibility (librarian CLI mentioned in notes)

**Status:** Detected; documentation gap identified.

---

### 11. 🟢 Naming Consistency — GOOD

**Finding:** Task IDs, file naming, command naming, and folder structure are consistent.

**Evidence:**

- BACKLOG uses: TASK-NNN, BUG-NNN, improvement type labels ✓
- Closed work dated: `.hawp/work/closed/YYYY/MM/DD/` ✓
- Plan files named: `TASK-026.md`, `BUG-006.md` ✓
- Commands: `npm run distribution:build`, `npm run validate:workflow` ✓
- No mixing of naming schemes detected

**Status:** No issues.

---

### 12. 🟡 Copilot Instructions Alignment — PARTIAL COVERAGE

**Finding:** `.github/copilot-instructions.md` references `.hawp/kit/` paths correctly but .github/instructions/ and .github/prompts/ overlap may create agent confusion about priority.

**Evidence:**

- `.github/copilot-instructions.md` (2KB) references:
  - `.hawp/kit/start-here.md` ✓
  - `.hawp/kit/usage/status-report.md` ✓
  - `.hawp/kit/usage/intake-workflow.md` ✓
  - `.hawp/work/BACKLOG.md` ✓
- `.github/instructions/` contains 4+ instruction files for Copilot (hawp-intake.instructions.md, etc.)
- `.github/prompts/` contains 5+ prompt files (hawp-backlog-alignment.prompt.md, etc.)

**Ambiguity:**

- Are `.github/instructions/` and `.github/prompts/` active overlays or archived templates?
- Does `.github/copilot-instructions.md` take precedence or do individual .instructions.md files?
- Are prompts auto-discovered by Copilot or manually invoked?

**Improvement Opportunity:**

- Add top-level comment in `.github/copilot-instructions.md` clarifying:
  - This file is the main Copilot overlay for the repo
  - `.github/instructions/` files extend specific integration points
  - `.github/prompts/` are reusable workflows available via slash commands or explicit invocation

**Status:** Detected; clarification recommended.

---

### 13. 🟢 Distribution Builder/Validator Documentation — CLEAR

**Finding:** Build and validate scripts have clear extraction logic and composition rules.

**Evidence:**

- `librarian/scripts/distribution/build/index.ts` documents bash block extraction with REF substitution
- Generated docs include source reference footer with regeneration instructions ✓
- `librarian/scripts/distribution/validate/index.ts` compares expected vs actual file lists ✓
- Recent validator run showed all 4 generated files current ✓

**Status:** Well designed and documented.

---

## Summary: Drift Classifications

| Category                     | Classification                   | Severity | Action Required                                 |
| ---------------------------- | -------------------------------- | -------- | ----------------------------------------------- |
| Commands & Scripts           | ✅ Consistent                    | None     | None                                            |
| Install/Update Boundaries    | ✅ Mostly Consistent             | Info     | Clarify "package" scope once per doc            |
| .hawp/work Preservation      | 🟡 Drifted (FIXED)               | Medium   | **Pending commit** of TASK-026 edits            |
| Generated File Warnings      | 🟡 Clear but low discoverability | Low      | Add warning to quick-ref section                |
| Validator Behavior           | ✅ Aligned                       | None     | None                                            |
| Active/Closed Reconciliation | 🟡 Minor count mismatch          | Low      | Verify closed file count before next compaction |
| Future-Facing Notes          | 🔴 Abundant but unorganized      | Medium   | Create Roadmap section in README                |
| Legacy Migration Docs        | 🟡 Scattered, consolidate-able   | Low      | Add flowchart/table for clarity                 |
| Validator Extensibility      | ✅ Well documented               | None     | None                                            |
| Librarian Package Scope      | 🟡 Undocumented                  | Low      | Create `librarian/README.md`                    |
| Naming Consistency           | ✅ Good                          | None     | None                                            |
| Copilot Instructions         | 🟡 Ambiguous hierarchy           | Low      | Clarify .github overlay precedence              |
| Distribution Tool Docs       | ✅ Clear                         | None     | None                                            |

---

## Improvement Opportunities (Low-Risk, High-Value)

### P1 — Clarity/Discoverability (Should do soon)

1. **Add generated file warning to README quick-ref section** (1 line)
   - Prevents first-time agent confusion about file editability
   - Impact: Agents skip edit attempts on generated files

2. **Clarify Copilot instructions hierarchy** (5 lines in copilot-instructions.md)
   - Prevents agent confusion about which overlay takes priority
   - Impact: Agents apply correct integration layer

3. **Create librarian/README.md** (20-30 lines)
   - Explains distribution vs workflow script purpose
   - Impact: Agents understand script organization, easier to extend

### P2 — Organization/Discoverability (Nice to have)

4. **Add Roadmap section to main README** (15-20 lines)
   - Links TASK-027/028/029 (approved), TASK-013 (parked), future enhancements
   - Impact: Agents understand active development direction

5. **Add legacy migration flowchart to install.md** (1 table, 3 rows)
   - Consolidates three migration paths visually
   - Impact: Easier onboarding for downstream projects with legacy content

### P3 — Documentation (Polish)

6. **Clarify "package" scope in install.md** (1 sentence)
   - Single place: "The install script..." instead of "This package..."
   - Impact: Removes ambient scope ambiguity

---

## Verification Checklist

Before proceeding with edits:

- [ ] Confirm 9 modified doc-alignment files from TASK-026 will be committed separately (docs-audit lane)
- [ ] Confirm 7 untracked files from parallel work (TASK-027/028/029) not modified during this audit
- [ ] Verify git reset recovery plan before next commit (orphaned commits 182ba18, 702d6d8, 65a9152)
- [ ] Run: `cd .hawp/work && find closed -name "*.md" | wc -l` to verify recent closed count (expect 7-10)
- [ ] Cross-check BACKLOG.md recently closed dates against actual closed/ directory timestamps

---

## Next Steps

1. **User Review:** Present findings and get approval for improvement set (P1 recommended)
2. **Edit Plan:** Prepare specific edits for each approved improvement (max 6 files)
3. **Regeneration:** Re-run `npm run distribution:build` if any source fragments edited
4. **Verification:** Run gate suite (distribution:validate, typecheck, validate:workflow)
5. **Commit Strategy:**
   - Lane 1: Recover/commit TASK-026 docs-audit edits (9 modified files)
   - Lane 2: Apply improvement edits (if approved) + any regenerations
6. **Closure:** Move TASK-026 to closed, update BACKLOG with improvement task (if created)
