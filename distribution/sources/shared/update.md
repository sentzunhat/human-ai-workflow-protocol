# Update HAWP: Shared Concepts

## Execution Preflight (Run First)

Treat this run as a new execution work item for the current repository.

- Open a terminal at the target repository root.
- Confirm `.hawp/` already exists before running update.
- If `.hawp/` does not exist, run install first (`install-main.md` or `install-dev.md`), then run update.
- Use the generated branch command block exactly as written.
- Run once and verify output lines for `Source:` and `Source mode:`.

## Update Work Item Contract

Treat update as an execution work item.

- If the request is "update HAWP", "update the HAWP lib", or "refresh the workflow protocol", run the generated update command block in a terminal from the target repo root.
- The goal is to refresh `.hawp/LICENSE`, `.hawp/kit/**`, `.github/instructions/**`, and `.github/prompts/**` to the selected branch.
- Preserve `.hawp/work/**` and any existing `.github/copilot-instructions.md`.
- If the user asks to update the GitHub root folder with the new workflow protocol, that means refreshing the HAWP-managed files in `.github/instructions/` and `.github/prompts/`.

## Explicit Work Item Command (Branch-Aware)

Use this when you want one command that treats the run as a real execution work item and dispatches automatically:

````bash
OWNER="sentzunhat"
REPO="human-ai-workflow-protocol"
REF="dev"   # set to "main" for stable

case "$REF" in
  main|dev) ;;
  *)
    echo "Error: REF must be 'main' or 'dev'"
    exit 1
    ;;
esac

if [ -d ".hawp" ]; then
  MODE="update"
else
  MODE="install"
fi

GUIDE_URL="https://raw.githubusercontent.com/${OWNER}/${REPO}/${REF}/distribution/generated/${MODE}-${REF}.md"
echo "Execution work item: ${MODE} HAWP from ${OWNER}/${REPO}@${REF}"
echo "Guide: ${GUIDE_URL}"

curl -fsSL "$GUIDE_URL" | awk '
  /^```bash$/ { in_block = 1; next }
  /^```$/ && in_block { exit }
  in_block { print }
' | bash
````

Expected behavior:

- If `.hawp/` is missing, it runs install for the selected branch.
- If `.hawp/` exists, it runs update for the selected branch.
- Output prints `Source:` and `Source mode:` lines from the executed script.

## When to Update

Update HAWP when you already have it installed and you want to:

- Get the latest HAWP kit improvements, new templates, or new patterns.
- Refresh HAWP Copilot instructions and prompts with the latest guidance.
- Adopt improvements from the HAWP source repository while keeping your own work intact.
- Support new HAWP features or migration helpers without losing your project's decision records.

Update mode is intended to be copy/paste ready:

- Choose the branch guide you want (`update-main.md` or `update-dev.md`).
- Copy the command block exactly as shown.
- Paste and run in your terminal from your repository root (not in Copilot chat).
- No manual branch edits are required.

## What Update Does

The update operation:

1. Refreshes `.hawp/LICENSE` and all of `.hawp/kit/**` from the source repository.
2. Updates `.github/instructions/*.instructions.md` and `.github/prompts/*.prompt.md` with new HAWP guidance.
3. Seeds any missing `.hawp/work/` scaffold files (only if they don't already exist).
4. Automatically reconciles closed work: moves `.hawp/work/active/` files matching Done rows or `done`/`wont-fix` Active Work rows in `.hawp/work/BACKLOG.md` to `.hawp/work/closed/...`. Prints each move. Only runs when the backlog has at least one data row.
5. Retires orphan active items (files with no matching BACKLOG entry) to `closed/YYYY/MM/DD/`. Only runs when the backlog has at least one data row.
6. Handles migration from old `hawp/` (no dot), `.hawp/usage/`, `.hawp/status/`, and `.hawp/work/adrs/` layouts if detected.
7. Cleans up stale legacy root-level kit folders (`.hawp/templates`, `.hawp/patterns`, `.hawp/reviews`, `.hawp/examples`, `.hawp/types`, `.hawp/usage`) and any `.gitkeep` files under `.hawp/`.

## What Update Does NOT Do

- Does not overwrite files in `.hawp/work/` (your project work is preserved).
- Does not delete project work records; it may move eligible closed or orphan items from `active/` to `closed/` during reconciliation.
- Does not overwrite `.github/copilot-instructions.md` if you have customized it.
- Does not remove any of your existing project code, decisions, or evidence files.
- Does not change `.github/copilot-instructions.md` if it already exists (your customizations are safe).

## Update Is Safe to Re-Run

- You can run update multiple times on the same repository.
- It will not harm your existing project work or customizations.
- Running it again is idempotent with respect to project-owned files.

## Copy/Paste Behavior

- The generated command already includes `OWNER`, `REPO`, and `REF` values.
- `REF` is branch-specific per generated file.
  - `update-main.md` uses `REF="main"`.
  - `update-dev.md` uses `REF="dev"`.
- If you need dev updates, run `update-dev.md`; if you need stable main updates, run `update-main.md`.
- The command output prints `Source: <owner>/<repo>@<ref>` and the archive URL so you can verify which branch was fetched.
- For local branch testing (including unpushed changes), set `HAWP_LOCAL_CORE` to a local HAWP `core/` path before running the command. In that mode the script prints `Source mode: local core (...)` and skips the archive download.
- If you paste docs into Copilot chat instead of running the bash block in a terminal, Copilot may only analyze content and report "already present" without applying any file changes.
- Once pasted, the command completes update end-to-end without interactive prompts.

## Agent Execution Contract

- Required: run the **Update Command (Copy/Paste)** bash block in a terminal from repo root.
- Required: report execution proof, including command output lines with `Source:` and `Source mode:`.
- Required: report file-level result proof with:
  - `git status --short .hawp/LICENSE .hawp/kit .github/instructions .github/prompts`
  - `find .hawp/kit -maxdepth 2 -type f | head -n 20`
- Not allowed as final answer: "already present" without command execution proof.

If no files changed after execution, state that explicitly as an execution result (already up to date), not as a content comparison result.

## Privacy-Safe Evidence Logging

- Keep file evidence repo-relative from repository root.
- Do not write machine-local absolute paths into work artifacts.
- If terminal proof outputs include host-local prefixes (for example `<user-home>/...`, `<linux-home>/...`, `<windows-user-home>\\...`), redact only the machine-local prefix in persisted artifacts (for example `<repo-root-abs>`), while preserving command identity and repo-relative evidence.

## Implementation Reference

The update script is maintained in `distribution/sources/update/script.md` and composed into the generated branch-specific guides. The generated guide you are reading is the user-facing update entry point for this branch.

## The Update Process

Update works by:

1. Downloading the HAWP source repository from GitHub.
2. Using conservative copy semantics to avoid overwriting existing files.
3. Running automatic migrations if old HAWP layouts are detected.
4. Running backlog-based reconciliation to move done work from `active/` to `closed/` and retire orphan items, with each move logged as full repo-relative `source -> destination` paths.
5. Removing stale legacy kit folders and `.gitkeep` files.
6. Preserving all of your `.hawp/work/` and your GitHub customizations.
