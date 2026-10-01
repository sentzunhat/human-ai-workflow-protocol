# CLI regex pattern quoting discipline (rg/grep)

**UUID:** `cli-rg-quoting-b9f3e1d2`
**Type:** improvement
**Reported:** 2026-09-10
**Status:** `done`

---

## Input (what was reported)

Ripgrep CLI error: `error parsing flag -E: grep config error: unknown encoding: FAIL|---`.
Root cause — regex patterns containing `|`, `--`, or other special characters are
not always wrapped in single quotes, causing the shell to pass unquoted alternations
and flags to ripgrep instead of treating them as a single pattern argument.

## Context

Every agent/session that runs grep/rg commands with multi-alternation regex patterns
needs consistent quoting discipline. Without it, terminal errors are confusing and
the search silently fails (wrong matches).

## Analysis

**Root cause:** Shell parses `|` before the command sees it when not quoted. Patterns
like `gateway|interface|expire` without quotes split into positional arguments or
pipe operators.

**Directly verified:**

- The error message `unknown encoding: FAIL|---` confirms unquoted `FAIL|---` was
  parsed as a flag (`-E`) argument by rg, not a regex pattern.

**Likely fix:** Document and enforce single-quote wrapping for any regex pattern
containing `|`, `-`, or shell metacharacters. Add to `.github/copilot-instructions.md`
and agent mode docs.

**Scope — what else is affected:**

- `.github/copilot-instructions.md` — add quoting safety note
- `CLAUDE.md` — add quoting safety note
- `.claude/rules/` — update rules if needed
- `AGENTS.md` — update rules if needed

## Work Coordination

**Owner:** Codex
**Implementation status:** done (2026-09-10)
**Can implement now:** yes

---

## Completion Record (2026-09-10)

The regex quoting rule was folded into `.github/copilot-instructions.md` as part
of the consolidated Copilot agent-instruction update. The same rule has now been
propagated to `core/providers/.github/copilot-instructions.md` so GitHub provider
updates can carry the guidance to downstream repositories.

## Outcome

Copilot-facing HAWP instructions now tell agents to wrap `rg`/`grep` patterns
containing `|`, `-`, or shell metacharacters in single quotes.

## Verification

Directly checked `.github/copilot-instructions.md` and
`core/providers/.github/copilot-instructions.md` for the regex quoting rule.

## Close Checklist

- [x] Rule added to Copilot instructions.
- [x] Provider source pack aligned.
- [x] Plan moved to `closed/2026/09/10/cli-rg-quoting-b9f3e1d2/`.
