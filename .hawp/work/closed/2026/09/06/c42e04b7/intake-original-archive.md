# Repo-Absolute Path Handling Conventions

**UUID:** `c42e04b7-c19a-451a-a76d-12ed2e31a1d3`
**Type:** improvement
**Reported:** 2026-09-06
**Risk Level:** low

---

## Input (what was reported)

> Add path-handling rules to `.github/copilot-instructions.md` so agents anchor file operations to repo-root and capture proof before edits, preventing build-breaking changes like the import cleanup disaster in session `47c793d6`.

---

## Context

During CLI decomposition wiring (backlog item `47c793d6`), an overly aggressive import cleanup removed required dependencies (`sqlite`, `application/index`) from `model_commands.go`, breaking the build. The agent operated without anchoring to repo-root paths or capturing staging proof before committing changes. This caused multiple sessions of debugging to fix scope errors and undefined references.

**Directly verified:**

- Session `47c793d6` — `multi_replace_string_in_file` deleted `"github.com/sentzunhat/hawp/librarian/src/internal/application/index"` import, causing `appindex.NewEmbedService()` and `appindex.DefaultEmbeddingModel` to be undefined.
- Import cleanup operated on `librarian/src/internal/platform/cli/model_commands.go` without verifying the full dependency graph across all CLI handler files.
- Build was broken (`go vet ./internal/platform/cli/...` failed with `undefined: appindex`).

**Inferred (not yet proven):**

- Similar import-scope errors may occur during any multi-file refactor or cleanup where agents modify imports without understanding cross-references.
- Terminal output showing `pwd`, `git rev-parse --show-prefix` isn't captured in plans/evidence as proof of correct working directory.

**Scope — what else is affected:**

- `.github/copilot-instructions.md` (add path anchoring rules)
- All agent workflows that modify Go imports or file paths

---

## Analysis

**Root cause (or most likely cause):**
Agents don't have explicit rules about capturing repo-relative proof before editing files. Without `pwd`, `git rev-parse --show-prefix`, and `git diff --name-status` checkpoints, the agent operated in path-ambiguous mode and removed imports that were actually required downstream.

**Directly verified:**

- The fix required reading `model_commands.go` at lines 1-20 to see current (broken) imports, then at lines 80-100 to confirm `sqlite.Open` and `appindex.NewEmbedService` usage sites.
- Import block was rewritten from scratch rather than surgically adding one missing import, because the agent didn't know the full import set that should exist.

**Inferred (not yet proven):**

- The same pattern could affect any language project where transitive dependencies get stripped during cleanup operations.

**Scope — what else is affected:**

- Only `.github/copilot-instructions.md` needs changes (add path protocol section).

---

## Work Coordination

**Owner:** unassigned
**Implementation status:** not-started
**Parallel work risk:** low
**Can implement now:** yes
**Coordination note:** This and the MCP fallback protocol item both edit `.github/copilot-instructions.md`. Can batch into a single implementation.

---

## Next Step

1. Read current `.github/copilot-instructions.md`
2. Add a section: "Path and Import Handling Protocol" with these rules:
   - **Anchor to repo-root:** Before any file operation, capture `pwd`, `git rev-parse --show-prefix`. All references use exact repo-relative paths from root (e.g., `librarian/src/internal/platform/cli/model_commands.go`), never basenames alone.
   - **Import surgery rule:** When removing unused imports, verify no downstream files in the same package reference them. Use `grep_search` across all `.go` files in the package before any import deletion.
   - **Staging proof:** Before commit, run `git diff --name-status`, `git diff --check`, `git diff --cached --name-status`. If any change lacks a repo-root prefix or basename-only path, halt and correct.
3. Update existing sections if conflicting rules exist.
4. Validate: `go vet ./...` passes with new rules as guidance (not enforced).
