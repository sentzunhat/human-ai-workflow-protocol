# Safety: Install and Update Principles

These are foundational safety rules that apply to both install and update operations.

## Project Work Is Never Overwritten

- `.hawp/work/` is **project-owned and must not be overwritten**.
- This includes your BACKLOG, active work, parked work, closed work, decisions, and evidence files.
- Install and update operations preserve all existing `.hawp/work/**` files, always.

## Existing Copilot Instructions Are Preserved

- `.github/copilot-instructions.md` is only seeded on first install if it does not exist.
- If you have customized `.github/copilot-instructions.md`, update operations will not overwrite it.

## Install and Update Are Safe to Re-Run

- Both operations are idempotent.
- Running them multiple times is safe and supported.
- They use no-clobber copy semantics (`cp -Rn`) to avoid overwriting existing files.

## Legacy Layout Migration Is Automatic

If your repository has an older HAWP layout, migration runs automatically:

- `hawp/` (no leading dot, real directory only — symlinks are skipped) → migrated to `.hawp/`, preserving `hawp/work/`, `hawp/usage/`, and `hawp/status/` content.
- `.hawp/usage/` → migrated to `.hawp/work/` (`BACKLOG.md` → `work/BACKLOG.md`, `status/*` → `work/active/`, `*_ADR.md` → `work/decisions/YYYY/MM/DD/`).
- `.hawp/status/` → migrated to `.hawp/work/notes/YYYY/MM/DD/`; `STATUS.md` promoted to `.hawp/work/STATUS.md`.
- `.hawp/work/adrs/` → migrated to `.hawp/work/decisions/YYYY/MM/DD/`, then the legacy folder is removed.

## Active Work Reconciliation Runs Automatically

- After migration, the script reads `.hawp/work/BACKLOG.md` Done rows and Active Work rows with `done` or `wont-fix` status, and moves matching `.hawp/work/active/*.md` files to `.hawp/work/closed/...`. Each moved file is printed in the output.
- After reconciliation, `.hawp/work/active/` files with no matching BACKLOG entry are retired to `closed/YYYY/MM/DD/` using the filename date prefix. This pass only runs when the backlog has at least one data row, so fresh installs are not affected.
- All moves use the no-overwrite rule: if the destination already exists, the source file is not moved.

## Verification Before and After

- Review what the script will do before running it (read the script block first).
- After running install or update, check `.hawp/work/BACKLOG.md` and `.hawp/kit/` to confirm changes.
- If something looks wrong, `.hawp/work/` is already preserved and safe.

## Privacy-Safe Path Logging

- Do not persist machine-local absolute paths in plans, evidence, status reports, or prompts.
- Avoid storing host-local prefixes such as `<user-home>/...`, `<linux-home>/...`, or `<windows-user-home>\\...` in repository artifacts.
- If command output includes absolute host paths, redact only the machine-local prefix with a placeholder (for example `<repo-root-abs>`) while preserving command identity and repo-relative path evidence.
