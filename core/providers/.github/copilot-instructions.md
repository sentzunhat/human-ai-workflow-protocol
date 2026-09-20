# Rule Precedence (HAWP instructions)

When rules conflict, this order applies (highest → lowest):

1. **Current tool's own instructions** — `.github/copilot-instructions.md` for Copilot, `CLAUDE.md` for Claude Code, etc.
2. **Editor-agnostic entry point** — `AGENTS.md`
3. **Shared HAWP behavior sources** — `core/providers/shared/behaviors/`, mirrored to repo-local instructions in `.github/instructions/` and `.claude/rules/`
4. **Per-tool rule overlays** — `.claude/rules/*.md`, provider-generated files in `distribution/generated/`

HAWP kit templates (`.hawp/kit/`) are always the canonical policy source; instruction files adapt them for specific tools without contradicting them.

---

This repository uses HAWP as a lightweight workflow method.

Follow the repo-local HAWP guidance in:

- .hawp/kit/start-here.md
- .hawp/kit/usage/status-report.md

Use .hawp/kit/start-here.md as the operating guide for how this repo applies HAWP in practice.

Use .hawp/kit/usage/status-report.md when the user asks for a:

- status report
- checkpoint summary
- context transfer summary
- second-brain review artifact

Use .github/prompts/hawp-change-review-and-reference-sync.prompt.md when the user asks to review recent changes and synchronize HAWP/docs references.

Saved status reports belong in:

- .hawp/work/status/YYYY/MM/DD/{uuid}/status.md

For bugs/tasks, track in .hawp/work/BACKLOG.md. Active plan files go in .hawp/work/active/. Deferred items can live in .hawp/work/parked/. Close by moving to .hawp/work/closed/YYYY/MM/DD/.
Do not append completed work endlessly to BACKLOG.md. Move closed work to .hawp/work/closed/YYYY/MM/DD/ and keep BACKLOG.md compact.
When BACKLOG.md has many Done rows, create a work item titled "Compact BACKLOG.md and archive closed work." before adding more Done entries.
Follow backlog compaction guardrails in .github/instructions/hawp-backlog-alignment.instructions.md.
Use .github/prompts/hawp-backlog-alignment.prompt.md when asked to review or compact a backlog.

For changes that affect installed files or work in a downstream consumer repo, consult `.github/skills/multi-repo-context.md` for the source → downstream → consumer propagation model.

Keep the repo-local HAWP layer lean.

HAWP in this repository is a practical workflow layer, not a runtime engine, compiler, validator, orchestrator, or memory system.

Do not:

- invent per-field folders such as input/, context/, mission/, constraints/, output/, or checkpoint/
- imply a runtime engine, compiler, validator, orchestrator, persistence layer, or memory system
- overstate certainty; keep direct evidence separate from inference

Prefer compact, decision-useful outputs.

---

## MCP Tool Invocation Protocol

Default: try MCP first; report the exact tool name used.
Fallback: if MCP is unavailable (tool not found, server error, conversation doesn't expose tools), fall back to terminal commands (`go run ./cmd/hawp <command>`) with explicit disclosure ("MCP tool unavailable, using CLI fallback").
Never retry indefinitely — one fallback attempt per requested operation.

---

## Path and Import Handling Protocol

**Anchor to repo-root:** Before any file operation, capture `pwd`, `git rev-parse --show-prefix`. All references use exact repo-relative paths from root (e.g., `librarian/src/internal/platform/cli/model_commands.go`), never basenames alone.

**Import surgery rule:** Go imports are scoped to the file that declares them. Check usage in that file before removing an import; when moving code, inspect both source and destination imports. Compile/test the affected package afterward.

**Staging proof:** Before commit, run `git diff --name-status`, `git diff --check`, `git diff --cached --name-status`. If any change lacks a repo-root prefix or basename-only path, halt and correct.

---

## Context Gathering Before Changes

**Read before write:** Before any multi-file edit or import deletion, run `grep_search` across the relevant package/directory to find all usage sites. Only then proceed with edits.

**Batch reads:** When planning changes to N files, read all N files first (in parallel via multiple tool calls). Don't interleave reads and writes within a single change set.

**Dependency scan:** Use `rg` across the relevant package for symbol or signature changes. For import removal, check the declaring file; another file's imports do not satisfy its dependencies.

**Verify after batch:** After applying edits to N files, verify all N files compile/test before considering the change complete. Don't submit partial results and fix later.

## Validation Commands (canonical)

Source of truth: `CLAUDE.md`. Always use these commands when validating HAWP state.

```bash
cd librarian/src && go test ./... && go run ./cmd/hawp distribution sync && go run ./cmd/hawp check
```

Grouped validation:

```bash
go run ./cmd/hawp links check
go run ./cmd/hawp providers sync
go run ./cmd/hawp kit validate
go run ./cmd/hawp kit normalize --apply
go run ./cmd/hawp work validate
go run ./cmd/hawp work normalize --dry-run --validate
go run ./cmd/hawp check
```

## Agent Discipline

**Regex quoting:** Always wrap rg/grep patterns containing `|`, `-`, or shell metacharacters in single quotes. Unquoted alternations split into positional arguments or pipe operators, causing silent failures and wrong matches.

**Session checkpoints before compaction:** Before any `/compact` or context-reset, write or append a structured checkpoint to `.hawp/work/status/YYYY/MM/DD/{uuid}/status.md` using `hawp work status --title "checkpoint" --work-item <uuid>`; capture current plan IDs, completed commits (SHA), pending files, and verification status. After compaction, read the latest checkpoint entry and resume without re-investigation.

**Batch reads before writes:** When planning changes to N files, read all N files first (in parallel). Don't interleave reads and writes within a single change set.

**Verify after batch:** After applying edits to N files, verify all N files compile/test before considering the change complete. Don't submit partial results and fix later.
