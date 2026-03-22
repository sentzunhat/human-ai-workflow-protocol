# Install HAWP: Shared Concepts

## Execution Preflight (Run First)

Treat this run as a new execution work item for the current repository.

- Open a terminal at the target repository root.
- Use the **provider-specific** guide for your agent (this file's provider section and boundaries apply).
- If `.hawp/` already exists, run **update** for the same provider and branch instead of install.
- Run the generated command block exactly as written.
- Verify output includes `Source:`, `Provider:`, and `Source mode:`.

## Install Work Item Contract (shared)

- Install always refreshes `.hawp/LICENSE`, `.hawp/kit/**`, and seeds missing `.hawp/work/` scaffold files.
- Install also runs **one provider overlay** for this guide's provider. That means the matching provider folder is refreshed alongside the kit:
  - Claude Code: `core/providers/.claude/` → `.claude/rules/`, `CLAUDE.md`
  - Codex: `core/providers/.codex/` → `AGENTS.md`
  - GitHub: `core/providers/.github/` → `.github/instructions/`, `.github/prompts/`, `.github/copilot-instructions.md`
  - Cursor: `core/providers/.cursor/` → `.cursor/rules/`, `AGENTS.md`
  - Continue: `core/providers/.continue/` → `.continue/rules/`
- Preserve `.hawp/work/**` project records.

Provider-specific execution proof and optional guide-fetch helpers are in the **Install Contract** section above the branch steps.

## When to Install

Install when you want HAWP kit plus agent overlays in a repository that does not have `.hawp/` yet (or you are re-running install intentionally).

## What Install Does (all providers)

1. Downloads HAWP `core/` from the selected branch (or uses `HAWP_LOCAL_CORE`).
2. Runs legacy layout migrations when detected.
3. Refreshes `.hawp/kit/**` and seeds `.hawp/work/` scaffold when missing.
4. Installs **only this guide's provider overlay folder** (see provider boundaries + contract).

## What Install Does NOT Do

- Does not overwrite `.hawp/work/**`.
- Does not install other providers' overlays (e.g. a Cursor guide does not write `.github/`).

## Install Is Safe to Re-Run

Idempotent for project-owned files. Kit and provider-managed paths refresh each run.

## Copy/Paste Behavior

- `REF` and `PROVIDER` are pre-set in the command block for this guide.
- For local testing: `export HAWP_LOCAL_CORE="/path/to/human-ai-workflow-protocol/core"`.
- Run in a terminal from repo root — not doc inspection alone.

## Agent Execution Contract (shared minimum)

- Required: run **Install Command (Copy/Paste)** in a terminal after reviewing the script block.
- Required: report `Source:`, `Provider:`, `Source mode:` from output.
- Required: provider-specific file proof from the Install Contract section.
- Optional: **Guide fetch (review-first)** writes a script to `/tmp` for inspection — run it explicitly; never pipe remote content directly to `bash`.
- Not allowed: "already present" without execution proof.

## Implementation Reference

Composed from `distribution/sources/install/script-core.md` + `providers/<provider>/script-install.md` + `script-footer.md`.
