# Install HAWP: Shared Concepts

## Execution Preflight (Run First)

Treat this run as a new execution work item for the current repository.

- Open a terminal at the target repository root.
- Decide mode by checking `.hawp/`:
  - If `.hawp/` does not exist, run install.
  - If `.hawp/` already exists, run update for the same branch (`update-main.md` or `update-dev.md`).
- Use the generated branch command block exactly as written.
- Run once and verify output lines for `Source:` and `Source mode:`.

## Install Work Item Contract

Treat install as an execution work item.

- If the request is "install HAWP", run the generated install command block in a terminal from the target repo root.
- The goal is to add or refresh `.hawp/LICENSE`, `.hawp/kit/**`, `.github/instructions/**`, and `.github/prompts/**`.
- Seed missing `.hawp/work/**` scaffold files if they are absent.
- Preserve `.hawp/work/**` project records and any existing `.github/copilot-instructions.md`.
- If the user asks to update the GitHub root workflow protocol during install, that means refreshing the HAWP-managed files under `.github/instructions/` and `.github/prompts/`.

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

## When to Install

Install HAWP when you have a repository where you want to:

- Adopt the HAWP task-shaping protocol for planning and context transfer.
- Use HAWP's work tracking structure (`.hawp/work/`, backlog, active items, decisions, evidence).
- Integrate HAWP Copilot instructions and prompt templates for your team.
- Ensure consistent patterns across your project for human–AI collaboration.

Install mode is intended to be copy/paste ready:

- Choose the branch guide you want (`install-main.md` or `install-dev.md`).
- Copy the command block exactly as shown.
- Paste and run in your terminal from your repository root (not in Copilot chat).
- No manual branch edits are required.

## What Install Does

The install operation:

1. Creates `.hawp/` at your repository root if it doesn't exist.
2. Copies the HAWP kit (protocol docs, templates, patterns, examples) into `.hawp/kit/`.
3. Seeds missing `.hawp/work/` scaffold files (starter `BACKLOG.md`, folder structure, README files).
4. Adds HAWP Copilot instructions under `.github/instructions/` and `.github/prompts/`.
5. Seeds `.github/copilot-instructions.md` only if it doesn't already exist.
6. Handles migration from old `hawp/` (no dot), `.hawp/usage/`, `.hawp/status/`, and `.hawp/work/adrs/` layouts if they are detected.

## What Install Does NOT Do

- Does not overwrite your existing `.hawp/work/` files (your project work is preserved).
- Does not delete project work records; it may move eligible closed or orphan items from `active/` to `closed/` during reconciliation.
- Does not overwrite `.github/copilot-instructions.md` if you already have one.
- Does not copy this repository's operating state (our BACKLOG, decisions, evidence) into your repository.

## Install Is Safe to Re-Run

- You can run install multiple times on the same repository.
- It will not overwrite any of your existing project work or customizations.
- HAWP-managed kit and overlay files are refreshed to the latest source on every run.

## Copy/Paste Behavior

- The generated command already includes `OWNER`, `REPO`, and `REF` values.
- `REF` is branch-specific per generated file.
  - `install-main.md` uses `REF="main"`.
  - `install-dev.md` uses `REF="dev"`.
- The command output now prints the exact source ref and archive URL before applying changes.
- For local branch testing (including unpushed changes), set `HAWP_LOCAL_CORE` to a local HAWP `core/` path before running the command. In that mode the script prints `Source mode: local core (...)` and skips the archive download.
- If you paste docs into Copilot chat instead of running the bash block in a terminal, Copilot may only analyze content and report "already present" without applying any file changes.
- The final output prints `Refreshed: .hawp/LICENSE, .hawp/kit/**, .github/instructions/*, .github/prompts/*` so you can quickly confirm install finished.
- Once pasted, the command completes install end-to-end without interactive prompts.

## Agent Execution Contract

- Required: run the **Install Command (Copy/Paste)** bash block in a terminal from repo root.
- Required: report execution proof, including command output lines with `Source:` and `Source mode:`.
- Required: report file-level result proof with:
  - `git status --short .hawp/LICENSE .hawp/kit .github/instructions .github/prompts`
  - `find .hawp/kit -maxdepth 2 -type f | head -n 20`
- Not allowed as final answer: "already present" without command execution proof.

If no files changed after execution, state that explicitly as an execution result (already up to date), not as a content comparison result.

## Implementation Reference

The install script is maintained in `distribution/sources/install/script.md` and composed into the generated branch-specific guides. The generated guide you are reading is the user-facing install entry point for this branch.

## The Install Process

Install works by:

1. Downloading the HAWP source repository (`core/` directory) from GitHub.
2. Using conservative copy semantics (`cp -Rn`, no-clobber) to avoid overwriting existing files.
3. Running automatic migrations if old HAWP layouts are detected.
4. Printing reconciliation activity when done/wont-fix work items are moved from `active/` to `closed/` (each line now shows exact repo-relative `source -> destination` paths).
5. Completing with `.hawp/work/` fully preserved.
