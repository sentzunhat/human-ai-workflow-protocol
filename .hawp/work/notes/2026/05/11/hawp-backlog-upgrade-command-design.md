# HAWP Backlog Upgrade Command Design

## Date

2026-05-11 (Revised)

## Purpose

Design the `.hawp/bin/hawp backlog upgrade` command for automatically detecting and applying safe mechanical fixes to old backlog/workflow record formats. This ensures legacy records automatically conform to current HAWP templates while blocking on ambiguous cases and preventing evidence invention.

---

## Part 1: Command Shape & Execution Model

### Primary Command

```bash
# Development invocation (local HAWP kit)
./.hawp/bin/hawp backlog upgrade [--dry-run | --apply] [options]

# Future packaged form (installed globally or via npm)
hawp backlog upgrade [--dry-run | --apply] [options]
```

### Default Behavior

- **Default mode:** `--dry-run` (propose fixes, do not modify files)
- **Entry point:** `./.hawp/bin/hawp` CLI script (executable)
- **Scope:** scan `.hawp/work/` for legacy formats and propose/apply fixes

### Modes

#### `--dry-run` (default)

```bash
./.hawp/bin/hawp backlog upgrade --dry-run [--format text|json]
```

**Behavior:**

- Scan backlog and work item files for legacy formats
- Generate upgrade plan with safe mechanical fixes listed
- Report blocked/ambiguous items requiring manual review
- Produce upgrade report (text or JSON)
- **No files modified**

**Output:**

- List of safe mechanical fixes (auto-applicable)
- List of blocked items (manual-required)
- Rationale and evidence for each category

**Exit codes:**

- `0` → no changes needed
- `2` → safe fixes available (ready for `--apply`)
- `3` → blocked by ambiguity or error

#### `--apply`

```bash
./.hawp/bin/hawp backlog upgrade --apply [--force-dirty]
```

**Preconditions:**

1. Git working tree must be clean (or `--force-dirty` given)
2. Backlog and work files unchanged since last `--dry-run` scan

**Behavior:**

- Apply all safe mechanical fixes from plan
- Skip blocked items (requires manual resolution)
- Write evidence report with summary
- Rerun validator to confirm results
- Attempt git commit (if in git repo) or write backup artifacts

**Output:**

- Applied fixes summary (before/after snippets)
- Blocked/skipped items list
- Evidence report file path
- Validator results

**Exit codes:**

- `0` → success (all safe fixes applied)
- `1` → precondition failed or error
- `2` → success but blocked items remain (manual follow-up needed)

---

## Part 2: V1 Automatic Fixes (Safe, Always Mechanical)

### Fix Type A1: Status Token Normalization

**Trigger:** Backlog table contains non-standard status vocabulary.

**Old → New mappings:**

- `complete` → `done`
- `finished` → `done`
- `hold` → `parked`
- `inprogress` → `in-progress`
- `pending-review` → `plan-ready`

**Safety:** Direct policy mapping; no semantic change.
**Classification:** Mechanical (low risk)
**Evidence basis:** Direct (status vocabulary policy)

### Fix Type A2: Type Column Addition

**Trigger:** Backlog table missing explicit `type` column.

**Inference sources (in priority order):**

1. Plan file name pattern (TASK-NN, BUG-NN, IMP-NN)
2. Plan file template type (bug-plan.md vs intake-plan.md)
3. Git history (commit message containing type hint)
4. Title heuristic ("fix", "bug" → bug; "add", "update", "implement" → task; "refactor", "improve" → improvement)

**Auto-apply only if:**

- Confidence ≥ 90% (sources 1 or 2 match)
- Exactly one plan file found for the task ID

**Block if:**

- Confidence < 90% (heuristic-only inference)
- Multiple plan file candidates exist
- No plan file found and title is ambiguous

**Safety:** Mechanical when source is direct (filename/template); blocked when inferred.
**Classification:** Mechanical (high confidence) or blocked (manual-required if ambiguous)
**Evidence basis:** Direct (plan filename/path) or inferred (title heuristic with low confidence flag)

### Fix Type A3: Owner & Updated Columns Addition

**Trigger:** Backlog table missing `Owner` or `Updated` columns.

**Owner inference:**

- Extract from plan file YAML frontmatter if present
- If not present, use "unassigned" (never invent author)

**Updated inference:**

- Use file modification time (mtime) if available
- Or extract from backlog row date context (e.g., closed folder YYYY/MM/DD)
- Add comment: "auto-inferred from file mtime" for traceability

**Safety:** Mechanical when sourced directly; auto-filled with trace note otherwise.
**Classification:** Mechanical (low risk)
**Evidence basis:** Direct (file metadata) or inferred (mtime)

### Fix Type A4: Unambiguous Plan File Link Repair

**Trigger:** Backlog row references broken or outdated plan file path.

**Repair logic:**

1. Parse backlog row plan link target
2. Search `.hawp/work/` for matching plan file (by task ID)
3. If exactly one match found, propose link update
4. If no match or multiple matches, block (manual-required)

**Auto-apply only if:**

- Exactly one valid target found
- Target file is readable and parseable

**Block if:**

- No valid target found (orphaned record)
- Multiple candidates exist (ambiguous)

**Safety:** Direct file existence check; no guessing.
**Classification:** Mechanical if unambiguous; blocked if ambiguous
**Evidence basis:** Direct (file system scan)

### Fix Type A5: Missing Section Heading Scaffolding

**Trigger:** Closed task plan file lacks required modern template sections.

**Sections to scaffold (for closed tasks):**

1. `Outcome (filled at close)` — if missing
2. `Verification (filled at close)` — if missing
3. `Close Checklist` — if missing

**Scaffold content:**

```markdown
## Outcome (filled at close)

<!-- TODO: Add outcome summary. See ../evidence/YYYY/MM/DD/ for supporting files. -->

## Verification (filled at close)

<!-- TODO: Add verification checklist or summary of tests/reviews performed. -->

## Close Checklist

- [ ] Outcome section filled
- [ ] Verification section filled
- [ ] All evidence files documented
```

**Safety:** Empty scaffolding with TODO comments; no content invention.
**Classification:** Mechanical (low risk)
**Evidence basis:** Policy (required section structure)
**Trace:** Add comment linking to evidence folder path

### Fix Type A6: Legacy Folder Path Migration

**Trigger:** Detection of flat or pre-date-structured legacy folders.

**Migration paths:**

- `adrs/<file>` → `decisions/YYYY/MM/DD/<file>`
- `status/<file>` → `status/YYYY/MM/DD/<file>`
- `evidence/<file>` → `evidence/YYYY/MM/DD/<file>`
- `closed/<ID>.md` (flat) → `closed/YYYY/MM/DD/<ID>.md` (dated)

**Date inference (in priority order):**

1. Explicit date in backlog row (if available)
2. File modification time (mtime)
3. Backlog context (latest related row date)

**Auto-apply only if:**

- Date can be inferred with confidence ≥ 90%
- Target folder path is writable

**Block if:**

- Date cannot be inferred reliably
- Target path already occupied

**Safety:** Files moved, not deleted; traceability comments added.
**Classification:** Mechanical if date determinable; blocked if date unclear
**Evidence basis:** Direct (backlog date) or inferred (mtime)
**Trace:** Add comment to affected files: "Migrated from flat structure to YYYY/MM/DD on [date]."

### Fix Type A7: Backlog Row Table Schema Normalization

**Trigger:** Backlog table columns don't match current schema.

**Current schema columns:**

- ID | Type | Title | Status | Owner | Plan File | Updated

**Auto-apply fixes:**

- Insert missing columns in correct position
- Preserve existing data
- Fill new columns with values inferred from plan files/metadata

**Safety:** Purely structural; no row content changed.
**Classification:** Mechanical (low risk)
**Evidence basis:** Direct (file metadata)

---

## Part 3: V1 Blocked/Non-Automatic Fixes (Manual Required)

### Blocked Type B1: Ambiguous Type Inference

**Trigger:** Task type cannot be determined with ≥90% confidence.

**Example:** Backlog row title is "Update system" (could be task, improvement, or maintenance).

**Structured Explanation (always included in reports):**

- `rule`: AMBIGUOUS_TYPE_INFERENCE
- `confidence`: confidence score (0.0-1.0) for best candidate
- `candidates`: list of possible types inferred
- `reason`: why automation stopped (human-readable)
- `evidence`: supporting data (title, plan hints, git patterns)

**Action:** Block, report ambiguity with structured reason, require manual type assignment.

**Recovery:** User edits backlog row, adds explicit type, reruns `--dry-run`.

### Blocked Type B2: Orphaned Records

**Trigger:** Backlog row references a plan file that doesn't exist AND no evidence files exist AND no ADR exists.

**Example:** TASK-055 row in backlog, but no TASK-055.md file anywhere.

**Action:** Block, report orphaned status, require manual resolution (delete row, create plan, or archive).

**Recovery:** User decides: delete row, create minimal plan file, or mark as historical archive.

### Blocked Type B3: Multiple Plan File Candidates

**Trigger:** Multiple plan files found for same task ID in different date folders.

**Example:** BUG-003.md found in `closed/2026/05/08/` AND `closed/2026/05/10/`.

**Structured Explanation (always included in reports):**

- `rule`: MULTIPLE_PLAN_CANDIDATES
- `confidence`: 0.0 (no winner, ambiguous)
- `candidates`: list of all plan file paths with metadata (size, mtime)
- `reason`: why automation stopped (cannot determine canonical)
- `evidence`: backlog row reference date, file sizes, timestamps, content hints

**Action:** Block, report ambiguity with candidate list and metadata, require manual consolidation.

**Recovery:** User inspects evidence folder for context, determines canonical version, consolidates files. Rerun after consolidation.

### Blocked Type B4: Outcome/Verification/Evidence Content Synthesis

**Trigger:** V1 encounters missing Outcome, Verification, or evidence content.

**Action:** Scaffold empty sections only (see Fix Type A5); never invent content.

**Recovery:** User/DA fills sections manually with actual outcomes and verification evidence.

### Blocked Type B5: Root Cause or Rationale Rewrites

**Trigger:** Plan file has incomplete Analysis, Recommended Fix, or decision narrative.

**Action:** Block semantic rewrites; never invent or infer details.

**Recovery:** User/DA manually fills or edits narrative sections based on actual work performed.

### Blocked Type B6: Cross-file Semantic Edits

**Trigger:** Fix would require changing meaning across multiple files or backlog rows.

**Action:** Block; preserve manual review gate.

**Recovery:** User reviews and manually applies semantic changes.

---

## Part 4: Safety Rules (Never Auto-Fix)

### Absolute Constraints

1. **✗ Invented verification evidence**
   - Never auto-generate "tests passed" or "verified by" claims
   - Scaffold empty/TODO sections only

2. **✗ Semantic rewrites of intent** (EXPLICIT NO-SEMANTIC-MUTATION RULE)
   - Never rewrite Root Cause Analysis without explicit source
   - Never change Recommended Fix rationale
   - Never infer implementation details or decision logic
   - **DO NOT rewrite task titles** — preserve original human language
   - **DO NOT reinterpret status meanings** — only normalize vocabulary
   - **DO NOT rewrite human-authored analysis** — preserve narrative as-is
   - Tool may normalize structure, but MUST NOT change semantic meaning of user content

3. **✗ Destructive deletion**
   - Legacy files are moved or preserved, never deleted
   - Historical narratives are linked-to, not overwritten
   - Orphaned records are reported, not silently removed

4. **✗ Ambiguous mappings**
   - Multiple plan file candidates → block
   - Ambiguous type inference → block
   - Uncertain date mapping → block

5. **✗ Out-of-band modifications**
   - No implicit sidecars or temporary files
   - All file operations explicit in upgrade plan
   - All changes traceable and reversible

6. **✗ Out-of-repo evidence invention**
   - Evidence files must exist in `.hawp/work/evidence/`
   - Never synthesize, download, or generate evidence
   - Empty scaffolding only; content must be human-provided

7. **✗ Writes outside `.hawp/**` boundary\*\* (PATH-BOUNDARY PROTECTION)
   - Tool refuses all writes outside `.hawp/` directory tree
   - Prevents accidental repo-wide mutations or side effects
   - All file operations scoped to `.hawp/work/`, `.hawp/kit/`, backlog metadata
   - Future extensions must be explicitly approved in design
   - Boundary is strict: no escape hatches or relative path traversal

8. **✓ Idempotency guarantee** (EXPLICIT OPERATIONAL GUARANTEE)
   - Running `--apply` twice produces zero additional changes
   - Second run returns `0` with identical validation results
   - Critical property: file state is recognized as already-upgraded
   - Implementation: verify no diff between current and post-upgrade state

9. **✓ Validator remains authoritative** (CRITICAL GOVERNANCE RULE)
   - Upgrade must adapt to validator rules, never vice versa
   - If validator detects drift upgrade cannot fix, block
   - Validator rules are source-of-truth; upgrade implements them

---

## Part 5: Preconditions & Safety Gates

### For `--dry-run`

- No preconditions required
- Read-only access to `.hawp/work/` and git history

### For `--apply`

1. **Git working tree must be clean**
   - `git status --porcelain` returns empty (or `--force-dirty` given)
   - Ensures rollback via `git reset --hard` is possible

2. **Backlog and work files unchanged since dry-run**
   - File hash preconditions checked
   - Blocks if concurrent edits detected

3. **No unknown-risk operations in plan**
   - Only low-risk mechanical operations auto-applied
   - Blocked items skipped (not failed)

4. **Sufficient disk space & permissions**
   - Write access to `.hawp/work/` verified
   - No disk space errors during apply

### Git Commit Behavior

If in git repository:

- Commit applied fixes with message: `chore: upgrade legacy backlog/workflow records to v1 templates`
- Add evidence report file to commit
- Return commit hash for rollback reference

If not in git repository:

- Write backup artifacts to `.hawp/work/status/YYYY/MM/DD/backlog-upgrade-<timestamp>.backup`
- Document rollback procedure in evidence report

---

## Part 5.5: Evidence Integrity & Immutability

### Immutable Hashes for Auditability

All upgrade operations produce immutable hash artifacts for:

- Rollback verification
- Reproducibility validation
- CI/CD integration
- Agent auditability
- Forensic analysis

### Required Hash Artifacts

**In all upgrade plans:**

- Plan ID (unique, deterministic)
- Plan hash (SHA256 of complete plan)
- File hashes before apply (SHA256 of each affected file)
- Timestamp (ISO 8601)

**In evidence reports (post-apply):**

- Report ID and hash
- Before/after SHA256 for each modified file
- Plan hash (links evidence to executed plan)
- Validator state hash (before and after)
- All immutable and queryable

### Hash Usage Patterns

**Rollback verification:**

- Compare current SHA256 with evidence report after-hash
- If mismatch, file has been edited; rollback may conflict
- Operator decides: force rollback or investigate

**Reproducibility:**

- Rerun same plan by hash; verify idempotency
- Hashes must match for safety gate approval
- Prove: two runs produce identical results

**CI/CD integration:**

- Validate: post-apply validator state better than pre-apply
- Fail job if upgrade worsened compliance
- Pass only if upgrade improved or maintained validator status

**Agent auditability:**

- Agents verify plan hash matches claimed mutations
- Evidence hashes prove actual changes
- Validator hashes prove upgrade improved compliance

---

## Part 6: Example Outputs

### Dry-run Output (Text Format)

```
$ ./.hawp/bin/hawp backlog upgrade --dry-run

=== SCAN PHASE ===

Backlog file: .hawp/work/BACKLOG.md
  Status: Format B + F detected
  - Missing Type column (6 rows)
  - Missing Owner column (all rows)
  - Missing Updated column (all rows)

Work item files: 26 found
  - Current v1 format: 18 files
  - Missing Outcome/Verification: 5 files
  - Untyped legacy: 3 files
  - Broken references: 2 rows

Legacy folders:
  - ./closed/ (flat): 8 files
  - ./adrs/: 4 files
  - ./status/ (flat): 0 files

=== UPGRADE PLAN ===

Plan ID: upgrade_20260511_143201
Safe mechanical fixes: 23 operations
Blocked/ambiguous: 3 items
Effort: Medium (1-2 hours manual follow-up for blocked items)

--- SAFE MECHANICAL FIXES (23 ops, auto-applicable) ---

[FIX-001] Normalize status tokens
  files: .hawp/work/BACKLOG.md
  changes: complete→done (3 rows), inprogress→in-progress (2 rows)
  evidence: policy

[FIX-002] Add Type column from plan filenames
  files: .hawp/work/BACKLOG.md
  inferred: 6 values (confidence ≥90%), source = plan filename
  affected rows: 6

[FIX-003] Add Owner + Updated columns
  files: .hawp/work/BACKLOG.md
  source: plan file metadata (owner), file mtime (updated)
  affected rows: 32

[FIX-004] Repair unambiguous plan links
  files: .hawp/work/BACKLOG.md
  repairs:
    - TASK-042: active/TASK-042.md.old → active/TASK-042.md (unique match)
    - TASK-039: closed/TASK-039.md → closed/2026/04/22/TASK-039.md (unique match)

[FIX-005] Scaffold missing required sections
  files: .hawp/work/closed/2026/05/09/TASK-041.md (+ 4 more)
  sections: Outcome, Verification, Close Checklist
  action: add empty/TODO scaffolds
  affected: 5 files

[FIX-006] Migrate legacy folder paths
  files:
    - 8 items: closed/* → closed/YYYY/MM/DD/
    - 4 items: adrs/* → decisions/YYYY/MM/DD/
  date source: file mtime + backlog context
  traceability: comment added to each file

--- BLOCKED ITEMS (3 items, manual required) ---

[BLOCKED-001] Ambiguous type inference
  row: TASK-053
  title: "Update system configuration"
  issue: Could be task, improvement, or maintenance
  resolution: User must assign type explicitly in backlog

[BLOCKED-002] Orphaned record (no plan file, no evidence)
  row: TASK-055
  title: "Refactor legacy storage layer"
  issue: No plan file found, no evidence files, no ADR
  resolution: User decides: delete row, create plan, or archive to decisions/legacy/

[BLOCKED-003] Multiple plan file candidates (ambiguous)
  row: BUG-003
  issue: BUG-003.md found in 2 different date folders
  locations:
    - .hawp/work/closed/2026/05/08/BUG-003.md
    - .hawp/work/closed/2026/05/10/BUG-003.md
  resolution: User inspects evidence, determines canonical version, consolidates

=== NEXT STEPS ===

1. Review automatic fixes above (all safe mechanical operations)
2. Resolve 3 blocked items (manual follow-up):
   - Assign type for TASK-053
   - Decide fate of TASK-055 (delete/create/archive)
   - Consolidate BUG-003 plan files
3. Run: ./.hawp/bin/hawp backlog upgrade --apply
   (blocked items will be skipped)

No files modified (dry-run mode).
Exit code: 2 (safe fixes available)
```

### Dry-run Output (JSON Format)

```json
{
  "command": "hawp backlog upgrade --dry-run",
  "timestamp": "2026-05-11T14:32:01Z",
  "planId": "upgrade_20260511_143201",
  "mode": "dry-run",
  "detections": {
    "backlogFormat": ["untyped-rows", "old-backlog-schema"],
    "workFilesCount": 26,
    "legacyFoldersFound": ["closed-flat", "adrs-folder"],
    "brokenReferences": 2
  },
  "safeMechanicalFixes": [
    {
      "fixId": "FIX-001",
      "type": "status-normalization",
      "file": ".hawp/work/BACKLOG.md",
      "operations": [
        { "old": "complete", "new": "done", "count": 3 },
        { "old": "inprogress", "new": "in-progress", "count": 2 }
      ]
    },
    {
      "fixId": "FIX-002",
      "type": "type-column-add",
      "file": ".hawp/work/BACKLOG.md",
      "inferred": 6,
      "confidence": "high",
      "source": "plan-filename"
    },
    {
      "fixId": "FIX-005",
      "type": "scaffold-sections",
      "files": [".hawp/work/closed/2026/05/09/TASK-041.md"],
      "sections": ["Outcome", "Verification", "Close Checklist"],
      "action": "add-empty-with-todo"
    }
  ],
  "blockedItems": [
    {
      "blockId": "BLOCKED-001",
      "reason": "ambiguous-type-inference",
      "taskId": "TASK-053",
      "details": "Title is too generic; cannot infer type with confidence ≥90%"
    },
    {
      "blockId": "BLOCKED-002",
      "reason": "orphaned-record",
      "taskId": "TASK-055",
      "details": "No plan file, no evidence, no ADR found"
    }
  ],
  "summary": {
    "totalOperations": 26,
    "safeMechanical": 23,
    "blocked": 3,
    "estimatedApplyTime": "< 5 seconds",
    "estimatedManualFollowUp": "1-2 hours"
  },
  "nextSteps": [
    "Review safeMechanicalFixes above",
    "Resolve blockedItems (manual step)",
    "Run: hawp backlog upgrade --apply"
  ],
  "exitCode": 2
}
```

### Apply Output

```
$ ./.hawp/bin/hawp backlog upgrade --apply

=== PRECONDITIONS CHECK ===

Git working tree: CLEAN ✓
Backlog files: UNCHANGED (hash match) ✓
Validator state: STABLE ✓

All preconditions met. Proceeding with apply.

=== APPLYING SAFE MECHANICAL FIXES ===

[FIX-001] Normalize status tokens
  .hawp/work/BACKLOG.md: 5 tokens normalized
  ✓ Applied

[FIX-002] Add Type column
  .hawp/work/BACKLOG.md: 32 rows updated with Type values
  ✓ Applied

[FIX-003] Add Owner + Updated columns
  .hawp/work/BACKLOG.md: all rows populated
  ✓ Applied

[FIX-004] Repair plan links
  .hawp/work/BACKLOG.md: 2 links repaired
  ✓ Applied

[FIX-005] Scaffold missing sections
  .hawp/work/closed/2026/05/09/TASK-041.md: Outcome, Verification, Close Checklist added
  .hawp/work/closed/2026/05/09/TASK-040.md: Outcome, Verification added
  (+ 3 more files)
  ✓ Applied (5 files modified)

[FIX-006] Migrate legacy folder paths
  closed/TASK-033.md → closed/2026/04/15/TASK-033.md
  closed/TASK-034.md → closed/2026/04/18/TASK-034.md
  (+ 6 more files)
  adrs/architecture-001.md → decisions/2026/03/22/architecture-001.md
  (+ 3 more files)
  ✓ Applied (12 files moved)

=== SUMMARY OF APPLIED FIXES ===

Total fixes applied: 23 operations
Files modified: 33 (5 backlog changes, 5 section scaffolds, 12 folder migrations, 11 link repairs)
Blocked items skipped: 3 (await manual resolution)

Evidence report: .hawp/work/evidence/2026/05/11/backlog-upgrade-applied-20260511-143215.md

=== BLOCKED ITEMS (Skipped) ===

[BLOCKED-001] TASK-053 — ambiguous type (manual assignment required)
[BLOCKED-002] TASK-055 — orphaned record (manual resolution required)
[BLOCKED-003] BUG-003 — multiple candidates (manual consolidation required)

=== VALIDATOR RERUN ===

Running validator to confirm upgrade success...

Backlog consistency: PASS
Closed task completeness: PASS
Evidence integrity: PASS
Verification clarity: PASS

Overall: PASS ✓

=== GIT COMMIT ===

Committing applied fixes...
  commit: a1b2c3d4e5f6... "chore: upgrade legacy backlog/workflow records to v1 templates"
  files: 28 changed, 45 insertions(+), 12 deletions(-)
  rollback: git reset --hard <parent_commit>

=== NEXT STEPS ===

1. Manual follow-up for 3 blocked items:
   - TASK-053: assign type explicitly
   - TASK-055: decide on plan creation vs archival
   - BUG-003: consolidate duplicate plan files

2. After resolving blocked items, rerun:
   ./.hawp/bin/hawp backlog upgrade --dry-run
   (should return 0 — no changes needed)

Upgrade complete. Exit code: 0 (success, blocked items remain)
```

---

## Part 7: Command Invocation Reference

### Arguments & Flags

```bash
./.hawp/bin/hawp backlog upgrade \
  [--dry-run | --apply] \           # mode (default: --dry-run)
  [--export-plan <path>] \          # export plan as JSON (governance workflows)
  [--format text | json] \          # output format (default: text)
  [--output <path>] \               # save output to file (optional)
  [--validate] \                    # rerun validator after apply (apply mode only)
  [--force-dirty] \                 # allow apply with dirty git tree (risky)
  [--verbose]                       # include detailed reasoning (optional)
```

### Usage Examples

```bash
# Dry-run, human-readable output (default)
./.hawp/bin/hawp backlog upgrade

# Dry-run, structured JSON output (for tooling)
./.hawp/bin/hawp backlog upgrade --format json --output /tmp/plan.json

# Export plan for governance review before apply
./.hawp/bin/hawp backlog upgrade --export-plan plan.json

# Apply fixes after manual review of dry-run
./.hawp/bin/hawp backlog upgrade --apply

# Apply fixes and automatically rerun validator (closes loop)
./.hawp/bin/hawp backlog upgrade --apply --validate

# Force apply in dirty git tree (not recommended)
./.hawp/bin/hawp backlog upgrade --apply --force-dirty

# Verbose mode with detailed detection reasoning
./.hawp/bin/hawp backlog upgrade --dry-run --verbose
```

### Exit Codes

| Code | Meaning                                                                                 |
| ---- | --------------------------------------------------------------------------------------- |
| `0`  | Success; no changes needed (dry-run) or success with no blocked items (apply)           |
| `1`  | Error: precondition failed, invalid input, or I/O error                                 |
| `2`  | Success with changes proposed (dry-run) or success with blocked items remaining (apply) |
| `3`  | Dry-run blocked by ambiguity or safety gate (no safe fixes available)                   |

---

## Part 8: Implementation Workstreams (V1 Scope)

### Workstream 1: Command Structure & Entry Point

**Title:** Set up `.hawp/bin/hawp` CLI executable and command routing

**Tasks:**

- Create `.hawp/bin/hawp` as TypeScript/Node executable
- Implement argument parsing for `backlog upgrade [--dry-run|--apply]`
- Wire command to implementation module

**Effort:** 1 day | **Risk:** Low

### Workstream 2: Format Detection Engine

**Title:** Implement backlog format detection for all 6 legacy formats

**Tasks:**

- Create `detection/format-detector.ts` with detection rules
- Create `detection/backlog-parser.ts` to parse BACKLOG.md table
- Create `detection/task-scanner.ts` to scan work item files
- Unit tests for detection accuracy

**Effort:** 3 days | **Risk:** Medium (false positives must be rare)

### Workstream 3: Safe Mechanical Fix Plan Generator

**Title:** Build upgrade plan generator with automatic fix classifications

**Tasks:**

- Implement fix generators for all 7 safe fix types (A1–A7)
- Implement blocker detection for 6 blocked types (B1–B6)
- Create plan serialization (artifact format)
- Hash calculation for precondition matching

**Effort:** 4 days | **Risk:** Medium (must ensure all fixes are truly safe)

### Workstream 4: CLI Output Formatters

**Title:** Implement human-readable text and structured JSON outputs

**Tasks:**

- Create `output/text-formatter.ts` for human-readable reports
- Create `output/json-formatter.ts` for structured output
- Implement `--format` flag routing
- Pretty-print with proper alignment and emphasis

**Effort:** 2 days | **Risk:** Low

### Workstream 5: Apply Engine with Safety Gates

**Title:** Implement file modification with preconditions and rollback support

**Tasks:**

- Git precondition checks (clean working tree)
- File hash precondition validation
- Safe file write operations (atomic or with backup)
- Git commit or backup artifact generation
- Rollback documentation

**Effort:** 5 days | **Risk:** High (must avoid partial application, data loss)

### Workstream 6: Evidence Report Generation

**Title:** Generate detailed post-apply evidence report

**Tasks:**

- Record all applied fixes with before/after snippets
- List blocked/skipped items with reasons
- Include validator rerun results
- Store in `.hawp/work/evidence/YYYY/MM/DD/`

**Effort:** 2 days | **Risk:** Low

### Workstream 7: Integration Tests & Safety Validation

**Title:** Comprehensive test fixtures and safety regression tests

**Tasks:**

- Create fixtures for all 6 legacy formats (Format A–F)
- Create fixtures for edge cases (ambiguous, orphaned, multiple candidates)
- Integration tests for dry-run (verify detection accuracy)
- Integration tests for apply (verify file modifications, rollback)
- Safety gate regression tests (no evidence invention, no semantic rewrites)

**Effort:** 5 days | **Risk:** Medium (test fixtures must be realistic)

### Workstream 8: User Documentation & Runbook

**Title:** Write user-facing docs and troubleshooting guide

**Tasks:**

- Command reference with examples
- Step-by-step upgrade guide
- Troubleshooting section for common errors
- Manual recovery procedure for blocked items
- FAQ for common questions

**Effort:** 2 days | **Risk:** Low

---

## Part 9: V1 Implementation Scope Summary

### Includes ✓

- Format detection for 6 legacy formats (A–F)
- 7 safe automatic fix types (A1–A7)
- 6 blocked/non-automatic fix types (B1–B6)
- Dry-run mode with text/JSON output
- Apply mode with preconditions and git support
- **Validator integration:** `--validate` flag for post-apply validation rerun
- **Plan export mode:** `--export-plan` for governance and review workflows
- **Explicit path-boundary protections:** refuses writes outside `.hawp/**` tree
- **Explicit no-semantic-mutation rule:** structure normalization only, content preservation
- Evidence report generation with optional validation summary
- Full test coverage
- User documentation

### Excludes ✗

- **AI-assisted evidence synthesis (deferred to V2 — V1 stays purely mechanical and deterministic)**
- Interactive approval workflow (future enhancement)
- Bulk migration for multiple HAWP roots (future enhancement)
- Custom repair rules/plugins (future enhancement)

### V1 Strategic Decision: Mechanical-Only, No AI

**Keeping V1 purely mechanical and deterministic is intentional and strategic.**

Rationale:

- **Predictability:** Operator knows exactly what changes and why; no surprises
- **Auditability:** All mutations are policy-driven, traceable, and reversible
- **Trust:** No unpredictable inference or synthesis; rules are explicit
- **Historical integrity:** Legacy records preserved without model rewriting
- **Foundation:** Establishes deterministic baseline before AI assistance layers on top

**Clean operational architecture:**

- **Validator:** Detects drift and inconsistencies (truth checker)
- **Upgrade:** Applies safe mechanical fixes (structural normalizer)
- **Librarian:** Discovers, indexes, surfaces patterns (future: knowledge layer)
- **HAWP:** Governs workflows, enforces policy (orchestrator)

AI-assisted synthesis (V2+) will layer on top with explicit governance gates, evidence tracking, and review workflows.

---

## Part 10: Architecture & Extensibility

### Data Model Architecture (JSON-First)

**Critical design principle: Generate structured objects internally; render text/json from shared model.**

NOT: Text output is source-of-truth (with JSON secondary)
YES: Internal objects are source-of-truth (with text/json rendered)

Why JSON-first internally:

- Decouples UI rendering from business logic
- Enables future extensions: web dashboards, agents, indexing, queues
- Ensures consistency across CLI, API, web outputs
- Future-proofs: CLI, web APIs, agent orchestration share same model
- Semantic analysis and tooling integration become straightforward

Internal data flow:

```
Backlog files → Detection objects → Plan object → Formatters
                                    ↓
                                    ├→ text (CLI display)
                                    ├→ json (API/tooling)
                                    └→ evidence report
```

Key types always generated as internal objects:

- `DetectionReport` (scan results)
- `BacklogFixPlan` (complete plan)
- `BacklogFixOperation[]` (individual operations)
- `BlockedItem` with explicit rule/confidence/candidates
- `EvidenceReport` (post-apply summary with hashes)

Formatters are pure rendering:

- `textFormatter(plan: BacklogFixPlan): string`
- `jsonFormatter(plan: BacklogFixPlan): string`
- `reportFormatter(evidence: EvidenceReport): string`

### Folder Structure

```
.hawp/bin/
  hawp                              # Main CLI executable (TypeScript)

.hawp/kit/lib/backlog-upgrade/      # Shared upgrade module
  cli.ts                            # CLI entry point
  types.ts                          # Type definitions

  detection/
    format-detector.ts              # Format detection engine
    backlog-parser.ts               # BACKLOG.md parser
    task-scanner.ts                 # Work file scanner

  planner/
    fix-generator.ts                # Safe fix generators (A1–A7)
    blocker-detector.ts             # Blocked item detector (B1–B6)
    plan-generator.ts               # Plan serialization

  applicator/
    apply-engine.ts                 # Safe apply with preconditions
    git-helper.ts                   # Git operations
    backup-helper.ts                # Backup artifact generation

  output/
    text-formatter.ts               # Human-readable output
    json-formatter.ts               # Structured JSON output
    report-generator.ts             # Evidence report

  __tests__/
    fixtures/                       # Test fixtures (6 legacy formats + edges)
    *.test.ts                       # Integration tests
```

### Future Extension Points

1. **Custom repair rules:** Plugin system for domain-specific fixes
2. **AI-assisted drafting:** Add `--ai-draft` flag for V2 (Outcome/Verification synthesis)
3. **Interactive approval:** Interactive CLI mode for manual case-by-case approval
4. **Bulk migration:** `--for-each-hawp-root` flag to upgrade multiple HAWP installations
5. **Policy customization:** Pluggable policy engine for status vocabulary, folder naming, etc.

---

## Part 11: Open Questions

1. **Command location:** Should `.hawp/bin/hawp` be shell script or Node executable? (Recommend: Node shebang for portability)

2. **Backup strategy:** For non-git repos, should we create `.git`-like snapshot or separate backup files? (Recommend: separate backup files under `evidence/`)

3. **Blocked item handling:** If `--apply` encounters blocked items, should it fail (exit 1) or succeed partially (exit 2)? (Recommend: exit 2 with clear message about manual follow-up)

4. **Concurrent edits detection:** Should file hash precondition also check git status, or rely on mtime? (Recommend: check both for safety)

5. **Evidence report storage:** Should it go to `evidence/YYYY/MM/DD/` or `status/YYYY/MM/DD/`? (Recommend: `evidence/` for consistency with validator finding storage)

---

## Part 11.5: Governance Principles & Strategic Architecture

### Validator Remains Authoritative

**Critical governance rule (Safety Rule #9):**

- Upgrade tool must adapt to validator rules, never vice versa
- If validator detects drift that upgrade cannot fix, block
- Validator rules are source-of-truth; upgrade implements them
- Protects long-term integrity: tool serves protocol, not convenience

**Why this matters:**

- Prevents tool from weakening validator standards for expedience
- Ensures validator remains independent truth checker
- Enables policy evolution without corrupting upgrade logic
- Maintains clear separation of concerns

### Idempotency Guarantee

**Critical operational property (Safety Rule #8):**

- Running `--apply` twice produces identical outcome both times
- Second run returns exit code `0` (no changes needed)
- State after first apply is stable and recognized as valid
- Implementation: verify no diff between current state and post-apply expectation

**Why this matters:**

- Enables safe automated/scripted apply workflows
- Simplifies CI/CD integration (idempotent operations are composable)
- Proves determinism and reproducibility
- Makes tool suitable for orchestration and queueing

### Strategic Architecture Convergence

**HAWP protocol evolution path:**

```
HAWP protocol (workflow policy)
    ↓
Validator (detect drift, truth checker)
    ↓
Upgrade tool (normalize structure, implement policy)
    ↓
Librarian/Indexer (discover patterns, future: knowledge layer)
    ↓
Multi-agent orchestration (future: autonomous coordination)
```

This path is cleaner and safer than:

- Jumping directly to "AI agents modifying repos autonomously"
- Trying to make tools intelligent without governance foundations
- Mixing truth-checking with automated fixing

**V1 strategic position:**

- Establish deterministic, policy-driven baseline
- Prove upgrade tool can be trusted (all mechanical, all traceable)
- Build immutability/hash infrastructure for future auditability
- Prepare for V2+ where AI-assisted synthesis layers on with explicit governance

---

## Part 12: Sign-off

**Design Date:** 2026-05-11 (Revised)
**Domain:** HAWP Backlog Tooling
**Status:** Ready for implementation review
**Next Step:** User approval → implement Workstreams 1–8 → deliver V1

---

## Quick Reference: Safe Fix Types vs Blocked Types

| Category             | Fix ID | Type                  | Auto-Apply? | Evidence Basis                            |
| -------------------- | ------ | --------------------- | ----------- | ----------------------------------------- |
| **Safe (V1 auto)**   | A1     | Status normalization  | YES         | Direct (policy)                           |
|                      | A2     | Type column add       | YES\*       | Direct (filename) if confidence ≥90%      |
|                      | A3     | Owner/Updated add     | YES         | Direct (metadata) or inferred (mtime)     |
|                      | A4     | Link repair           | YES\*       | Direct (file scan) if unambiguous         |
|                      | A5     | Section scaffold      | YES         | Policy (template requirement)             |
|                      | A6     | Legacy path migration | YES\*       | Direct (backlog date) or inferred (mtime) |
|                      | A7     | Table schema norm     | YES         | Policy (schema definition)                |
| **Blocked (manual)** | B1     | Ambiguous type        | NO          | Ambiguity blocks                          |
|                      | B2     | Orphaned record       | NO          | No plan/evidence blocks                   |
|                      | B3     | Multiple candidates   | NO          | Ambiguity blocks                          |
|                      | B4     | Evidence synthesis    | NO          | Policy (no invention)                     |
|                      | B5     | Rationale rewrites    | NO          | Policy (no semantic changes)              |
|                      | B6     | Cross-file edits      | NO          | Policy (preserve complexity)              |

\* With confidence/unambiguity checks; blocked if ambiguous
