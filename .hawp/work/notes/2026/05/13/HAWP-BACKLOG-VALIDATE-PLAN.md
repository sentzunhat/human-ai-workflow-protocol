# Implementation Plan: `hawp backlog validate` Command

**Date:** 2026-05-11
**Author:** Planning DA
**Status:** Plan (not implemented)
**Related:** `hawp backlog upgrade` (sibling command)

---

## 1. Current Validator Entry Points

### Existing Script

- **Location:** `librarian/scripts/validate-hawp-workflow/`
- **Entry point:** `npm run validate:workflow` (defined in `librarian/package.json`)
- **Command line:** `tsx scripts/validate-hawp-workflow/index.ts [options]`

### Module Export Points

| Module             | Exports                                                                                                                | Purpose                               |
| ------------------ | ---------------------------------------------------------------------------------------------------------------------- | ------------------------------------- |
| `cli.ts`           | `parseArgs()`, `getHelpText()`, `resolveWorkDirectory()`                                                               | CLI arg parsing + help                |
| `orchestrate.ts`   | `parseBacklog()`, `orchestrateValidation()`                                                                            | Core validation orchestration         |
| `reporter.ts`      | `formatReport()`                                                                                                       | Human-readable text report formatting |
| `types.ts`         | Type definitions (`ValidationReport`, `CheckResult`, etc.)                                                             | Shared type contracts                 |
| `validations/*.ts` | `checkBacklogConsistency()`, `checkClosedTaskCompleteness()`, `checkEvidenceIntegrity()`, `checkVerificationClarity()` | Individual validation checks          |

### Current Behavior

- Auto-discovers `.hawp/work` by searching up from `cwd` (max 10 levels)
- Supports explicit `--hawp-root` and `--work-root` override
- Exits with code 0 on PASS/WARN, code 1 on FAIL
- Outputs human-readable text to stdout

---

## 2. Proposed Command Shape

### Surface Command

```bash
.hawp/bin/hawp backlog validate [options]
```

### Future Packaged Form

```bash
hawp backlog validate [options]
```

### Supported Invocations

```bash
# Default: validate current repo's .hawp folder, output text
.hawp/bin/hawp backlog validate

# Explicit HAWP root
.hawp/bin/hawp backlog validate --hawp-root /path/to/.hawp

# Explicit work directory
.hawp/bin/hawp backlog validate --work-root /path/to/.hawp/work

# Human-readable text output (default)
.hawp/bin/hawp backlog validate --format text

# Machine-readable JSON output
.hawp/bin/hawp backlog validate --format json

# Combined options
.hawp/bin/hawp backlog validate --hawp-root ./project/.hawp --format json

# Debug mode (include closed-task diagnostics)
.hawp/bin/hawp backlog validate --debug
```

### Default Behavior

- Validates the current repo's `.hawp/work` directory
- Outputs human-readable text
- Exits with code 0 on PASS/WARN, code 1 on FAIL
- Never modifies any files (read-only)

---

## 3. Minimal Module Boundaries

### Proposed Directory Structure

```text
librarian/
  src/
    commands/
      backlog/
        validate.ts           # CLI command handler for 'backlog validate'
        types.ts              # Command-specific types (CLIOptions, CommandResult)
    cli/
      args.ts                 # Shared CLI arg parsing utilities
      formatter.ts            # Output formatters (text, json)
    validation/
      index.ts                # Re-export existing checks + orchestration
      types.ts                # Shared validation types (ValidationReport, etc.)
      reporter.ts             # Report formatting (text + json)
      orchestrate.ts          # Validation orchestration (moves from scripts/)
      backlog-consistency.ts  # Existing check (moves from scripts/)
      closed-task-completeness.ts
      evidence-integrity.ts
      verification-clarity.ts
      id-parser.ts
```

### Migration Path (Phase 2)

- **Phase 1 (MVP):** Wrap existing `validate-hawp-workflow/` scripts, leave files in place
- **Phase 2 (Refactor):** Move validation logic to `src/validation/`, keep scripts as thin wrappers

### Rationale

- Keep validation logic centralized in `src/validation/`
- Commands layer is thin and focused on CLI integration
- Reusable by both CLI and npm scripts
- Supports future commands (`fix-up`, `upgrade`) sharing same validation foundation

---

## 4. TypeScript Types and Interfaces Needed

### Command-Level Types

```typescript
// src/commands/backlog/types.ts

export interface BacklogValidateOptions {
  hawpRoot?: string; // Path to .hawp directory
  workRoot?: string; // Path to .hawp/work directory
  format: "text" | "json"; // Output format
  debug?: boolean; // Include debug diagnostics
  help?: boolean; // Show help
}

export interface CommandResult {
  success: boolean; // Whether validation passed (no FAILs)
  message: string; // Human-readable result message
  report: ValidationReport; // Full validation report
  exitCode: 0 | 1; // Process exit code
}
```

### CLI Argument Types

```typescript
// src/cli/args.ts

export interface ParsedArgs {
  hawpRoot?: string;
  workRoot?: string;
  format: "text" | "json";
  debug?: boolean;
  help?: boolean;
  _: string[]; // Remaining positional args
}

export interface ResolvedPaths {
  workDir: string; // Resolved .hawp/work directory
  hawpRoot: string; // Resolved .hawp directory
}
```

### Output Formatter Types

```typescript
// src/cli/formatter.ts

export interface FormattedOutput {
  text: string;             // Plain text output
  json: object;             // JSON-serializable object
}

export const formatValidationReport = (
  report: ValidationReport,
  format: 'text' | 'json'
): FormattedOutput => { ... }
```

### Reused Validation Types

```typescript
// From src/validation/types.ts (existing, no changes)

export type CheckStatus = "PASS" | "FAIL" | "WARN";
export type ReportStatus = "PASS" | "FAIL";

export interface ValidationReport {
  timestamp: string;
  checks: {
    backlogConsistency: BacklogCheck;
    closedTaskCompleteness: ClosedTaskCheck;
    evidenceIntegrity: EvidenceCheck;
    verificationClarity: VerificationCheck;
  };
  summary: {
    passed: number;
    failed: number;
    warnings: number;
    totalChecks: number;
  };
  overallStatus: ReportStatus;
}

// ... existing check types unchanged
```

---

## 5. CLI Argument Mapping

### Argument Parsing Rules

| Argument             | Type   | Default         | Mapping                                       |
| -------------------- | ------ | --------------- | --------------------------------------------- |
| `--hawp-root <path>` | string | (auto-discover) | `options.hawpRoot`                            |
| `--work-root <path>` | string | (auto-discover) | `options.workRoot`                            |
| `--format <fmt>`     | enum   | `text`          | `options.format` (validate: `text` \| `json`) |
| `--debug`            | flag   | false           | `options.debug = true`                        |
| `--help` / `-h`      | flag   | false           | `options.help = true`                         |

### Argument Resolution Order

1. **Explicit `--work-root`:** Use directly if it exists
2. **Explicit `--hawp-root`:** Use `<hawp-root>/work` if it exists
3. **Auto-discovery:** Search up from `cwd` for `.hawp/work` (max 10 levels)
4. **Not found:** Error and exit(1)

### Example Parsing Flow

```typescript
const options = parseBacklogValidateArgs(process.argv.slice(2));
// Input:  ['--hawp-root', '/project/.hawp', '--format', 'json']
// Output: { hawpRoot: '/project/.hawp', format: 'json', debug: false, help: false }

const paths = await resolveValidationPaths(options);
// Resolves: { workDir: '/project/.hawp/work', hawpRoot: '/project/.hawp' }

const report = await orchestrateValidation(paths.workDir, ...);
// Validates and returns ValidationReport

const output = formatValidationReport(report, options.format);
// Returns { text: '...', json: {...} }

console.log(output[options.format]);
process.exit(report.overallStatus === 'FAIL' ? 1 : 0);
```

---

## 6. Exit Code Behavior

### Exit Codes

| Code       | Meaning       | When                                                             |
| ---------- | ------------- | ---------------------------------------------------------------- |
| 0          | Success       | Validation completed with PASS or WARN (no FAILs)                |
| 1          | Failure       | At least one FAIL check OR missing `.hawp/work` OR parsing error |
| (implicit) | Process error | Uncaught exception                                               |

### Validation Status to Exit Code

```
ValidationReport.overallStatus
  ├─ 'PASS'  → exit(0)
  └─ 'FAIL'  → exit(1)
```

### Integration with CI/CD

```bash
# CI: fail on any validation failure
.hawp/bin/hawp backlog validate || { echo "Validation failed"; exit 1; }

# Pre-commit: warn but don't fail
.hawp/bin/hawp backlog validate --format json > /tmp/report.json
# Inspect report JSON for warnings/failures
```

---

## 7. Relationship to `hawp backlog upgrade`

### Command Family Hierarchy

```
hawp backlog
  ├── validate      ← This plan (Level 0: read-only, no changes)
  └── upgrade       ← Sibling command (Level 1/2: dry-run → apply, modifies files)
```

### Safety Model Alignment

| Level | Command             | Behavior        | Files Modified        |
| ----- | ------------------- | --------------- | --------------------- |
| 0     | `validate`          | Scan + report   | None (read-only)      |
| 1     | `upgrade --dry-run` | Propose changes | None (simulation)     |
| 2     | `upgrade --apply`   | Apply changes   | Yes (with safeguards) |

### Data Flow

```
User runs: hawp backlog validate
    ↓
Command finds `.hawp/work`
    ↓
Orchestrates 4 validation checks
    ↓
Generates ValidationReport
    ↓
Formats as text/json
    ↓
Outputs to stdout + sets exit code

Later, when upgrade command exists:
User runs: hawp backlog upgrade --dry-run
    ├─ Reads ValidationReport or invokes validate internally
    ├─ Proposes fixes based on findings
    └─ Outputs upgrade plan (text/json)

User reviews plan, then:
hawp backlog upgrade --apply --plan <path>
    ├─ Verifies plan preconditions
    ├─ Applies changes
    └─ Outputs summary + rollback info
```

### Design Principle: Clear Separation

- **`validate`:** "What's wrong?" (finding/diagnosis)
- **`upgrade`:** "How to fix it?" (planning/action)

They do not overlap:

- `validate` is purely read-only
- `upgrade` reads `validate` output but adds its own proposal logic
- No cross-talk; each command independently callable and testable

---

## 8. Implementation Work Items

### Item 1: Create CLI Wrapper Command Handler

**File:** `librarian/src/commands/backlog/validate.ts`
**Work:**

- [ ] Define `BacklogValidateOptions` interface
- [ ] Implement `handleBacklogValidate(options: BacklogValidateOptions): Promise<CommandResult>`
- [ ] Integrate with existing `orchestrateValidation()` function
- [ ] Call appropriate formatter based on `options.format`
- [ ] Return `CommandResult` with exit code

**Dependencies:** Existing `orchestrate.ts`, `reporter.ts`, `types.ts`

**Estimated effort:** 1-2 hours

---

### Item 2: Create Shared CLI Utilities

**File:** `librarian/src/cli/args.ts`
**Work:**

- [ ] Extract arg parsing logic from existing `validate-hawp-workflow/cli.ts`
- [ ] Create `parseBacklogValidateArgs()` function
- [ ] Create `resolveValidationPaths()` function (replaces `resolveWorkDirectory`)
- [ ] Support `--hawp-root`, `--work-root`, `--format`, `--debug`
- [ ] Return `ParsedArgs` interface

**Dependencies:** None (pure utility)

**Estimated effort:** 1 hour

---

### Item 3: Create JSON Output Formatter

**File:** `librarian/src/cli/formatter.ts`
**Work:**

- [ ] Reuse existing `formatReport()` for text output
- [ ] Implement `toJSON()` serializer for `ValidationReport`
- [ ] Create unified `formatValidationReport(report, format)` function
- [ ] Ensure JSON output is machine-parseable and contains all check details
- [ ] Document JSON schema (optional but helpful for downstream tools)

**Dependencies:** Existing `reporter.ts`, `types.ts`

**Estimated effort:** 1-2 hours

---

### Item 4: Create `.hawp/bin/hawp` Bootstrap Script

**File:** `.hawp/bin/hawp`
**Work:**

- [ ] Create shell/node bootstrap script that routes commands
- [ ] Detect `backlog validate` subcommand
- [ ] Invoke `tsx librarian/src/commands/backlog/validate.ts` with args
- [ ] Pass through exit code
- [ ] Eventually support other subcommands (`backlog upgrade`, etc.)

**Dependencies:** Items 1-3

**Estimated effort:** 1-2 hours

---

### Item 5: Create Unit Tests

**File:** `librarian/src/commands/backlog/validate.test.ts` + CLI tests
**Work:**

- [ ] Test arg parsing for all flag combinations
- [ ] Test path resolution (auto-discover, explicit root, errors)
- [ ] Test text output formatting
- [ ] Test JSON output formatting
- [ ] Test exit codes (0, 1)
- [ ] Mock `orchestrateValidation()` to verify integration

**Dependencies:** Items 1-3

**Estimated effort:** 2-3 hours

---

### Item 6: Update npm Scripts and Documentation

**File:** `librarian/package.json`, `README.md`
**Work:**

- [ ] Keep existing `npm run validate:workflow` working (backward compat)
- [ ] Add new `npm run backlog:validate` (optional)
- [ ] Document `.hawp/bin/hawp backlog validate` usage in README
- [ ] Add examples (default, --format json, --hawp-root)
- [ ] Link to design note and upgrade command when available

**Dependencies:** Items 1-4

**Estimated effort:** 1 hour

---

### Item 7: Integration Test with Real Repo Structure

**File:** `test/fixtures/sample-hawp-repo/` + integration test
**Work:**

- [ ] Create minimal `.hawp` fixture with test backlog
- [ ] Run `hawp backlog validate` against it
- [ ] Verify text and JSON output match expected schema
- [ ] Verify exit codes are correct
- [ ] Test with external repo (`--hawp-root` flag)

**Dependencies:** Items 1-4

**Estimated effort:** 1-2 hours

---

## 9. Risks and Open Questions

### Risk: Format Stability for JSON Output

**Issue:** JSON schema for `ValidationReport` may change when other features (AI drafting, finding IDs, patch hunks) are added later.

**Mitigation:**

- Add `schema_version: "1.0"` field to JSON output
- Document JSON structure clearly
- Use semantic versioning for schema changes
- Provide migration guide when schema evolves

---

### Risk: External Root Validation Scope Creep

**Issue:** Validating multiple external `.hawp` roots could be requested, but current design is single-root.

**Mitigation:**

- For now, support only one root per invocation
- If multi-root is needed later, create separate `hawp backlog validate-all` command or `--roots` flag
- Keep validate pure and composable so callers can loop over multiple roots

---

### Open Question: Should `validate` auto-invoke on `.hawp/bin/hawp backlog upgrade`?

**Design options:**

1. **Option A (Current Plan):** Each command is independent. `upgrade` can invoke `validate` if needed but doesn't require it.
2. **Option B:** `upgrade` always runs validation first and includes findings in its proposal.
3. **Option C:** `upgrade --auto-validate` flag to opt-in to inline validation.

**Recommendation:** Option A for now. Keep commands decoupled. Document that running `validate` before `upgrade` is recommended practice but optional.

---

### Open Question: Passthrough of Backlog Row Metadata in JSON

**Issue:** Should JSON output include full backlog row content (title, type, status, detail link) for downstream consumption?

**Recommendation:** Yes. Include full row data in `checks.backlogConsistency.activeWork.missing[]` items, not just IDs. Enables downstream tools to surface issue context without re-parsing backlog.

---

### Open Question: Dry-run Format for Future `upgrade` Command

**Issue:** When `upgrade --dry-run` is implemented, should it emit the same JSON schema as `validate`, or a different structure?

**Recommendation:** Different structure. `validate` reports **findings**; `upgrade` reports **proposals**. They are semantically different and should have distinct JSON schemas. Consider adding `reportType: "findings" | "proposals"` to distinguish them.

---

## 10. Implementation Sequencing

### Recommended Order

1. **Create CLI utilities** (Item 2) — foundation for everything
2. **Create command handler** (Item 1) — core logic
3. **Create formatters** (Item 3) — output layer
4. **Create bootstrap script** (Item 4) — user-facing entry point
5. **Write tests** (Item 5) — verify all pieces work together
6. **Integration testing** (Item 7) — real-world validation
7. **Documentation & npm scripts** (Item 6) — finalize

### Estimated Total Effort

| Phase         | Items   | Estimate       |
| ------------- | ------- | -------------- |
| Core          | 1, 2, 3 | 3-4 hours      |
| Integration   | 4       | 1-2 hours      |
| Testing       | 5, 7    | 3-4 hours      |
| Documentation | 6       | 1 hour         |
| **Total**     |         | **8-11 hours** |

---

## 11. Success Criteria

The `hawp backlog validate` command will be considered complete when:

- [ ] Command is callable as `.hawp/bin/hawp backlog validate`
- [ ] Supports all required flags (`--hawp-root`, `--work-root`, `--format`, `--debug`)
- [ ] Outputs correct exit codes (0 for PASS/WARN, 1 for FAIL)
- [ ] Text output matches existing validator formatting
- [ ] JSON output is valid, well-documented, and machine-parseable
- [ ] Validation logic is reused (no duplication from `validate-hawp-workflow/`)
- [ ] All unit tests pass (arg parsing, formatting, integration)
- [ ] Real-world integration test passes (fixture repo)
- [ ] Backward compatibility maintained (existing `npm run validate:workflow` still works)
- [ ] Documentation is complete (README, usage examples, exit code behavior)
- [ ] No file modifications (read-only guarantee held)

---

## 12. Non-Goals (Out of Scope for This Plan)

- Implementing `hawp backlog upgrade` command (separate plan)
- Adding AI-assisted diagnostics (future enhancement)
- Supporting multi-root validation in single invocation (may be future)
- Modifying existing validation check logic (only wrapping)
- Changing validator exit code semantics

---

## 13. Next Steps

1. **Review Plan:** User reviews and approves plan structure
2. **Identify Gaps:** Flag any missing pieces or unanswered questions
3. **Implementation Gate:** Once approved, proceed to implementation in sequence
4. **Parallel Track:** Consider designing `upgrade` command plan concurrently
5. **Archive Plan:** Move plan to `.hawp/work/closed/` once implementation begins

---

## Appendices

### A. Existing Validator Checks (for reference)

#### Check 1: Backlog Consistency

- Verifies active/closed/parked rows in BACKLOG.md have matching plan files
- Detects orphaned files with no backlog row
- Status: FAIL if files missing, PASS if complete

#### Check 2: Closed Task Completeness

- Verifies closed plan files have Outcome, Verification, Close Checklist sections
- Uses cutoff date (2026-05-10) for legacy tolerance (WARN) vs current (FAIL)
- Status: FAIL if recent files missing sections, WARN if legacy files incomplete

#### Check 3: Evidence Integrity

- Verifies referenced evidence files exist and are complete
- Checks for missing or broken evidence links
- Status: FAIL if critical evidence missing

#### Check 4: Verification Clarity

- Verifies verification statements are present and substantive
- Detects placeholder or incomplete verification text
- Status: WARN if verification unclear

### B. JSON Output Schema (Proposed)

```json
{
  "schema_version": "1.0",
  "timestamp": "2026-05-11T14:22:01Z",
  "command": "hawp backlog validate",
  "root": ".hawp/work",
  "checks": {
    "backlogConsistency": {
      "status": "PASS",
      "activeWork": { "total": 5, "found": 5, "missing": [] },
      "recentlyClosed": { "total": 10, "found": 10, "missing": [] },
      "parkedWork": { "total": 2, "found": 2, "missing": [] },
      "orphanedFiles": [],
      "orphanedParked": []
    },
    "closedTaskCompleteness": { ... },
    "evidenceIntegrity": { ... },
    "verificationClarity": { ... }
  },
  "summary": {
    "passed": 4,
    "failed": 0,
    "warnings": 0,
    "totalChecks": 4
  },
  "overallStatus": "PASS"
}
```

### C. Command Cheat Sheet

```bash
# Default (auto-discover, text output)
.hawp/bin/hawp backlog validate

# Explicit root with JSON
.hawp/bin/hawp backlog validate --hawp-root ./project/.hawp --format json

# Debug output
.hawp/bin/hawp backlog validate --debug

# Save JSON report for downstream tooling
.hawp/bin/hawp backlog validate --format json > /tmp/validation.json

# Use in CI/CD: fail if validation fails
.hawp/bin/hawp backlog validate || exit 1
```

---

**End of Plan**
