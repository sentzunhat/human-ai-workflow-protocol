# Update HAWP: Shared Concepts

## Execution Preflight (Run First)

- Open a terminal at the target repository root.
- Confirm `.hawp/` exists; if not, run install for this provider first.
- Use the matching provider update guide (same provider as install).
- Verify output includes `Source:`, `Provider:`, and `Source mode:`.

## Update Work Item Contract (shared)

- Refresh `.hawp/LICENSE`, `.hawp/kit/**`, and **this provider's overlay folder only**.
- The provider folder refresh depends on the guide:
  - Claude Code: `.claude/rules/` (`CLAUDE.md` is preserved on update)
  - Codex: `AGENTS.md`
  - GitHub: `.github/instructions/`, `.github/prompts/`, `.github/copilot-instructions.md`
  - Cursor: `.cursor/rules/`, `AGENTS.md`
  - Continue: `.continue/rules/`
- Preserve `.hawp/work/**`.
- Reconcile closed work from backlog when eligible.

Provider-specific proof is in the **Update Contract** section in this guide.

## Explicit dispatch

See the provider **Install Contract** **Guide fetch (review-first)** block — it selects update when `.hawp/` exists and writes a script to `/tmp` for review before execution.

## When to Update

When HAWP is installed and you want the latest kit and provider overlay for your agent.

## What Update Does (all providers)

1. Refreshes kit from source branch.
2. Refreshes the provider overlay folder documented in this guide's boundaries.
3. Seeds missing work scaffold only.
4. Runs migrations and backlog reconciliation when eligible.

## What Update Does NOT Do

- Does not overwrite `.hawp/work/**`.
- Does not refresh other providers' folders.

## Update Is Safe to Re-Run

Safe to run multiple times.

## Agent Execution Contract (shared minimum)

- Required: run **Update Command (Copy/Paste)** in a terminal after reviewing the script block.
- Required: `Source:`, `Provider:`, `Source mode:` in output.
- Required: provider file proof from Update Contract.
- Optional: **Guide fetch (review-first)** writes a script to `/tmp` for inspection — run it explicitly; never pipe remote content directly to `bash`.

## Implementation Reference

Composed from `distribution/sources/update/script-core.md` + `providers/<provider>/script-update.md` + `script-footer.md`.
