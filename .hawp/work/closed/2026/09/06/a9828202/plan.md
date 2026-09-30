# Batch File Reads Before Multi-File Changes

**UUID:** `d4e5f6a7-b8c9-0123-defa-234567890123`
**Type:** improvement
**Reported:** 2026-09-06
**Risk Level:** low

---

## Input (what was reported)

> Add batch-read-before-edit rules to `.github/copilot-instructions.md` so agents gather full context across related files before applying changes, preventing destructive patterns like the import cleanup that broke imports in `model_commands.go`.

---

## Context

The CLI decomposition session (`47c793d6`) demonstrated a damage pattern: an overly aggressive cleanup removed required imports from `model_commands.go` without first verifying which other CLI handler files (if any) shared those dependencies. The agent should have read the full dependency graph across all relevant files before executing import deletions via `multi_replace_string_in_file`.

Similarly, in the `local-print-farm` session (`70e4f88e`), 117 turns suggest iterative read-execute-retry patterns rather than upfront context gathering — the agent appeared to attempt commands without fully understanding the network state first.

**Directly verified:**

- Session `47c793d6` — `multi_replace_string_in_file` deleted `sqlite` and `application/index` imports from `model_commands.go` because the agent didn't scan for usage sites across all CLI handler files first.
- The fix required: reading lines 1-20 to see current (broken) imports, then lines 80-100 to find `appindex.NewEmbedService()` usage — i.e., the agent had to be told "read more context" instead of doing it proactively.

**Inferred (not yet proven):**

- Agents tend to optimize for quick edits over thorough reads because they're not given an explicit rule requiring batch reads before any edit.
- This pattern likely affects non-GO projects too — any multi-file refactor without upfront context gathering risks missing dependencies, shared abstractions, or cross-cutting concerns.

**Scope — what else is affected:**

- `.github/copilot-instructions.md` (add batch-read protocol)
- All agent workflows performing any file edits beyond trivial single-line changes

---

## Analysis

**Root cause (or most likely cause):**
No explicit "read before write" rule exists in Copilot instructions. The agent defaults to surgical edits when it could be batching reads first. Without a required pre-flight read phase, the agent misses cross-references, shared imports, and dependencies that aren't visible in the file being directly edited.

**Directly verified:**

- Import deletion happened via `multi_replace_string_in_file` without a preceding `grep_search` across the package to find all usage sites.
- The fix pattern was: broken → read more files → identify missing deps → restore imports → verify. This is reactive, not proactive.

**Inferred (not yet proven):**

- A batch-read rule would also improve non-code tasks like documentation updates, config changes, and workflow adjustments where context from multiple sources matters.

**Scope — what else is affected:**

- Only `.github/copilot-instructions.md` needs changes.

---

## Verification

Implemented 2026-09-06: Added "Context Gathering Before Changes" section to `.github/copilot-instructions.md` with batch-read-before-edit, dependency scan, and verify-after-batch rules. Prevents destructive patterns like the `model_commands.go` import cleanup from session `47c793d6`. Verified `go run ./cmd/hawp check --no-update-check` passes all 3 validations (kit validate, work validate, links check). Plan closed via commit `92cb6ac`.

---

## Work Coordination

**Owner:** unassigned
**Implementation status:** done
**Parallel work risk:** low
**Can implement now:** yes
**Coordination note:** This is the 4th of 4 related instruction improvements, all targeting `.github/copilot-instructions.md`. Can be batched into a single implementation.

---

## Next Step

1. Read current `.github/copilot-instructions.md`
2. Add a section: "Context Gathering Before Changes" with these rules:
   - **Read before write:** Before any multi-file edit or import deletion, run `grep_search` across the relevant package/directory to find all usage sites. Only then proceed with edits.
   - **Batch reads:** When planning changes to N files, read all N files first (in parallel via multiple tool calls). Don't interleave reads and writes within a single change set.
   - **Dependency scan:** For import/scope changes in Go, run `grep_search` across `**/*.go` in the same package before removing any import. The removed import must have zero remaining references.
   - **Verify after batch:** After applying edits to N files, verify all N files compile/test before considering the change complete. Don't submit partial results and fix later.
3. Validate: run `go vet ./...` after any test implementation to confirm no regressions.

---

## Outcome

- `.github/copilot-instructions.md` updated with "Context Gathering Before Changes" section (read before write, batch reads, dependency scan, verify after batch rules)
- Prevents destructive patterns like the `model_commands.go` import cleanup from session `47c793d6`
- `go run ./cmd/hawp check --no-update-check` passes all 3 validations (kit validate, work validate, links check)

---

## Close Checklist

- [x] Plan written and reviewed
- [x] Implementation completed
- [x] Verification evidence recorded
- [x] Changes committed (`92cb6ac` on `feature/v0.0.24-work-folder-normalization`)
