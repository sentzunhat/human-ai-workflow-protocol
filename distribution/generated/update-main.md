> **Agent instruction:** Update the local HAWP installation in this repository. Run the bash command block in the **Update Command (Copy/Paste)** section in a terminal from the repository root. Do not report "already present" based on content inspection alone — execute the script and report terminal output with `Source:` and `Source mode:` lines as proof.

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

# Repository Boundaries: What Gets Installed and Updated

The install and update operations respect clear boundaries to protect your project-specific work.

## HAWP-Managed Files (Refreshed on Every Install or Update)

These files are owned by HAWP and are written or replaced on every run:

- `.hawp/LICENSE` — Apache 2.0 license for the HAWP kit content
- `.hawp/kit/**` — Reusable HAWP protocol docs, templates, patterns, examples, references (replaced in full)
- `.github/instructions/*.instructions.md` — HAWP Copilot instructions
- `.github/prompts/*.prompt.md` — HAWP Copilot prompt templates

## Project-Owned Files (Never Overwritten)

These files are owned by your project and are never overwritten by install or update:

- `.hawp/work/**` — Your BACKLOG, active work, parked work, closed work, decisions, evidence, notes
- `.github/copilot-instructions.md` — Your custom Copilot instructions (seeded once when missing; never overwritten after)

## Scaffold Files (Seeded Once; Never Overwritten After)

These are created on first install if they do not already exist:

- `.hawp/work/README.md`, `.hawp/work/STATUS.md`, `.hawp/work/BACKLOG.md`
- `.hawp/work/active/README.md`, `parked/README.md`, `closed/README.md`
- `.hawp/work/decisions/README.md`, `evidence/README.md`, `status/README.md`, `notes/README.md`
- `.github/copilot-instructions.md`

## What Never Gets Installed Downstream

These files from the HAWP source repository are **never** copied to your target:

- The HAWP source repo's own `.hawp/work/BACKLOG.md`
- The HAWP source repo's own `.hawp/work/active/`, `closed/`, `decisions/`, `evidence/` operating state
- Any content from `benchmark/` (optional reference material, not installed)

## The Boundary Model

```
HAWP Source Repo (core/)
├─ .hawp/kit/             → Reusable HAWP assets (refreshed every install/update)
├─ .hawp/work/            → HAWP source repo's own operating state (NEVER installed)
└─ .github/               → HAWP overlay templates (refreshed every install/update)

Target Repo (your project)
├─ .hawp/kit/             ← Refreshed from HAWP source on every install/update
├─ .hawp/work/            ← YOUR project work (never overwritten; scaffold seeded once)
├─ .github/instructions/  ← HAWP overlays refreshed; copilot-instructions.md preserved
└─ [your files]           ← Your code, docs, and config (never touched)
```

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

# Update HAWP To Main Branch

Use this when you have HAWP already installed and want to upgrade to the latest stable improvements from the main branch.

## Prerequisites

- HAWP already installed in your repository (you have `.hawp/kit/` and `.hawp/work/`).
- `curl` and `tar` (standard Unix tools).
- Ability to run shell scripts or copy/paste a shell command block.

## Before You Update

1. **Check your `.hawp/work/BACKLOG.md`** to see if there are any active or parked work items.
   - The update script will automatically reconcile closed work items based on your backlog.
   - All your work will be preserved.

2. **Review any GitHub Copilot customizations** in `.github/copilot-instructions.md`.
   - If you have customized it, the update will preserve your customizations.
   - HAWP-provided instructions in `.github/instructions/` and `.github/prompts/` will be updated.

## Update Steps

1. **Open your repository in a terminal and move to its root directory.**

2. **Copy and run the command** from the "Update Command (Copy/Paste)" section below.
   - No edits are needed.
   - This guide already pins `REF="main"`.
   - The script output will show `Source: sentzunhat/human-ai-workflow-protocol@main`.
   - The command handles migration + update in one run.

   Optional local-source mode for branch testing:
   - If you need to test local or unpushed branch changes, run `export HAWP_LOCAL_CORE="/absolute/path/to/human-ai-workflow-protocol/core"` before running the command.
   - The script will print `Source mode: local core (...)` when that override is active.

3. **Review the output** to see what was updated.
   - Check that `.hawp/kit/` was refreshed with new protocol docs.
   - Check that stale legacy folders were removed (`.hawp/templates/`, `.hawp/patterns/`, etc. if present).
   - If the script printed `reconciled (link): <src> -> <dest>`, `reconciled (id-fallback): <src> -> <dest>`, or `retired (orphan): <src> -> <dest>` lines, those are expected — done work items were moved from `active/` to `closed/`.
   - The final output will show: `Refreshed: .hawp/LICENSE, .hawp/kit/**, .github/instructions/*, .github/prompts/*`

4. **Verify your work is intact**:
   - Check `.hawp/work/BACKLOG.md` to see if your items are still there.
   - Look for any automatically reconciled items moved to `closed/`.

## What Gets Updated

- `.hawp/LICENSE` — Apache 2.0 license (updated to latest).
- `.hawp/kit/**` — All HAWP protocol docs, templates, patterns, examples.
- `.github/instructions/*.instructions.md` — HAWP Copilot instructions.
- `.github/prompts/*.prompt.md` — HAWP Copilot prompt templates.
- Missing `.hawp/work/` scaffold files are seeded if they don't exist.

## What Is Preserved

- All of `.hawp/work/` — your backlog, active items, parked items, closed work, decisions, and evidence.
- `.github/copilot-instructions.md` — your custom instructions (not overwritten).
- Your project code and configuration.

## Automatic Work Reconciliation

On every run, the update script reads `.hawp/work/BACKLOG.md` and:

- Moves `.hawp/work/active/` files matching Done rows or `done`/`wont-fix` Active Work rows to `.hawp/work/closed/YYYY/MM/DD/`. Prints each move as `reconciled (link): <src> -> <dest>` or `reconciled (id-fallback): <src> -> <dest>`.
- Retires `.hawp/work/active/` files with no matching BACKLOG entry to `closed/YYYY/MM/DD/`. Prints each as `retired (orphan): <src> -> <dest>`.
- Both passes only run when the backlog has at least one data row.

All your work remains safe — no files are deleted, only moved.

## Troubleshooting

If something looks wrong after update:

- All files under `.hawp/work/` are unchanged, so your work is safe.
- If `.github/copilot-instructions.md` was not updated (expected if you customized it), you can manually merge in new guidance from the source.
- If stale folders were removed and you want them back, you can restore them from the HAWP source repository manually.
- Run the update command again; it is safe to re-run.

## Update Command (Copy/Paste)

Run this from the root of your target repository. No edits are required; branch selection is already configured in the command. Each run fetches the latest commit from that branch.

```bash
OWNER="sentzunhat"
REPO="human-ai-workflow-protocol"
REF="main"

echo "Source: ${OWNER}/${REPO}@${REF}"
echo "Archive: https://github.com/${OWNER}/${REPO}/archive/refs/heads/${REF}.tar.gz"

TMP_DIR=""
if [ -n "${HAWP_LOCAL_CORE:-}" ]; then
  SRC="${HAWP_LOCAL_CORE}"
  if [ ! -d "$SRC/.hawp/kit" ] || [ ! -d "$SRC/.github" ]; then
    echo "Error: HAWP_LOCAL_CORE must point to a core directory containing .hawp/kit and .github"
    exit 1
  fi
  echo "Source mode: local core (${SRC})"
else
  TMP_DIR="$(mktemp -d)"
  curl -fsSL "https://github.com/${OWNER}/${REPO}/archive/refs/heads/${REF}.tar.gz" \
    | tar -xz -C "$TMP_DIR"
  SRC="$TMP_DIR/${REPO}-${REF}/core"
  echo "Source mode: remote archive"
fi

if [ ! -d ".hawp" ]; then
  echo "Preflight: .hawp/ not found in this repository."
  echo "Run install-${REF}.md first, then run update-${REF}.md."
  if [ -n "$TMP_DIR" ] && [ -d "$TMP_DIR" ]; then
    rm -rf "$TMP_DIR"
  fi
  exit 1
fi
MDATE="$(date +%Y/%m/%d)"

# --- Helpers (no-clobber copy; never overwrite repo-owned files) ---
copy_dir_no_clobber() {
  src="$1"; dest="$2"
  if [ -d "$src" ]; then
    mkdir -p "$dest"
    cp -Rn "$src"/. "$dest"/ 2>/dev/null || true
  fi
}
copy_file_no_clobber() {
  src="$1"; dest="$2"
  if [ -f "$src" ] && [ ! -f "$dest" ]; then
    mkdir -p "$(dirname "$dest")"
    cp "$src" "$dest"
  fi
}
reconcile_closed_plans_from_backlog() {
  backlog=".hawp/work/BACKLOG.md"
  [ -f "$backlog" ] || return 0

  # Scans both ## Done rows (col5=Closed date) and ## Active Work rows where
  # status signals completion (done/wont-fix), so reconciliation works even
  # when items have not yet been reformatted into the Done section.
  awk -F'|' '
    function trim(s) { gsub(/^[[:space:]]+|[[:space:]]+$/, "", s); return s }
    /^## Done/        { section="done";   next }
    /^## Active/      { section="active"; next }
    /^## /            { section="";       next }
    section != "" && /^\|/ {
      id   = trim($2)
      col5 = trim($5)
      col6 = trim($6)
      if (id == "" || id ~ /^-+$/ || id == "ID") next
      if (section == "done") {
        closed = col5
        plan   = col6
      } else {
        # Active Work: col5=Status — only reconcile when status signals done
        if (col5 !~ /^(done|wont-fix)$/ && col5 !~ /^[0-9]{4}-[0-9]{2}-[0-9]{2}$/) next
        closed = (col5 ~ /^[0-9]{4}-[0-9]{2}-[0-9]{2}$/) ? col5 : ""
        plan   = col6
      }
      link = ""
      if (match(plan, /\(([^)]+)\)/)) {
        link = substr(plan, RSTART + 1, RLENGTH - 2)
      } else {
        link = plan
      }
      print id "\t" closed "\t" link
    }
  ' "$backlog" | while IFS=$'\t' read -r id closed link_path; do
    [ -n "$id" ] || continue

    closed_path=""
    case "$link_path" in
      .hawp/work/closed/*) closed_path="${link_path#.hawp/work/}" ;;
      work/closed/*) closed_path="${link_path#work/}" ;;
      closed/*) closed_path="$link_path" ;;
    esac
    closed_path="${closed_path%%#*}"

    if [ -n "$closed_path" ]; then
      plan_name="$(basename "$closed_path")"
      src=".hawp/work/active/$plan_name"
      dest=".hawp/work/$closed_path"
      if [ -f "$src" ] && [ ! -e "$dest" ]; then
        mkdir -p "$(dirname "$dest")"
        mv "$src" "$dest"
        echo "  reconciled (link): $src -> $dest"
      fi
      continue
    fi

    # Derive closed_dir: prefer explicit date from backlog, then from filename
    # prefix, then fall back to today so nothing is silently skipped.
    # Uses case-glob matching to stay portable (avoids grep -E / rg aliases).
    case "$closed" in
      [0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]) closed_dir="${closed//-//}" ;;
      *) closed_dir="" ;;
    esac

    find .hawp/work/active -maxdepth 1 -type f -name "*-${id}-*.md" | while IFS= read -r src; do
      [ -n "$src" ] || continue
      local_dir="$closed_dir"
      if [ -z "$local_dir" ]; then
        fname="$(basename "$src")"
        case "$fname" in
          [0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]-*)
            file_date="${fname:0:10}"
            local_dir="${file_date//-//}"
            ;;
          *)
            local_dir="$(date +%Y/%m/%d)"
            ;;
        esac
      fi
      dest=".hawp/work/closed/$local_dir/$(basename "$src")"
      if [ ! -e "$dest" ]; then
        mkdir -p "$(dirname "$dest")"
        mv "$src" "$dest"
        echo "  reconciled (id-fallback): $src -> $dest"
      fi
    done
  done
}

retire_orphan_active_items() {
  # Moves .hawp/work/active/ files that have no matching row in
  # .hawp/work/BACKLOG.md (no ID pattern match and no plan-link match) to
  # .hawp/work/closed/YYYY/MM/DD/. Skips when the backlog is absent or has
  # no data rows so fresh installs are never affected.
  backlog=".hawp/work/BACKLOG.md"
  [ -f "$backlog" ] || return 0

  known_ids="$(awk -F'|' '
    function trim(s) { gsub(/^[[:space:]]+|[[:space:]]+$/, "", s); return s }
    /^\|/ {
      id = trim($2)
      if (id == "" || id ~ /^-+$/ || id == "ID") next
      print id
    }
  ' "$backlog")"
  [ -n "$known_ids" ] || return 0

  known_links="$(awk -F'|' '
    function trim(s) { gsub(/^[[:space:]]+|[[:space:]]+$/, "", s); return s }
    /^\|/ {
      for (i = 2; i <= NF; i++) {
        cell = trim($i)
        if (match(cell, /\(([^)]+\.md[^)]*)\)/)) {
          link = substr(cell, RSTART + 1, RLENGTH - 2)
          n = split(link, parts, "/")
          gsub(/#.*$/, "", parts[n])
          if (parts[n] != "") print parts[n]
        }
      }
    }
  ' "$backlog" 2>/dev/null)"

  find .hawp/work/active -maxdepth 1 -type f -name '*.md' ! -name 'README.md' | while IFS= read -r src; do
    fname="$(basename "$src")"
    matched=0
    while IFS= read -r id; do
      [ -n "$id" ] || continue
      case "$fname" in
        *"-${id}-"*|*"-${id}.md") matched=1; break ;;
      esac
    done <<< "$known_ids"
    if [ "$matched" -eq 0 ] && [ -n "$known_links" ]; then
      while IFS= read -r link_name; do
        [ -n "$link_name" ] || continue
        [ "$fname" = "$link_name" ] && { matched=1; break; }
      done <<< "$known_links"
    fi
    [ "$matched" -eq 1 ] && continue
    case "$fname" in
      [0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]-*)
        file_date="${fname:0:10}"
        closed_dir="${file_date//-//}"
        ;;
      *)
        closed_dir="$(date +%Y/%m/%d)"
        ;;
    esac
    dest=".hawp/work/closed/$closed_dir/$fname"
    if [ ! -e "$dest" ]; then
      mkdir -p "$(dirname "$dest")"
      mv "$src" "$dest"
      echo "  retired (orphan): $src -> $dest"
    fi
  done
}

# --- 1. Migration: legacy hawp/ -> .hawp/ (preserves hawp/work/, hawp/usage/, hawp/status/) ---
# Only acts on a real directory; symlinks named hawp are left untouched.
if [ -d "hawp" ] && [ ! -L "hawp" ]; then
  mkdir -p .hawp/work
  copy_dir_no_clobber "hawp/work" ".hawp/work"
  if [ -d "hawp/usage" ]; then
    copy_file_no_clobber "hawp/usage/BACKLOG.md" ".hawp/work/BACKLOG.md"
    copy_dir_no_clobber  "hawp/usage/status"     ".hawp/work/active"
    mkdir -p ".hawp/work/decisions/$MDATE"
    for f in hawp/usage/*_ADR.md; do
      [ -f "$f" ] && copy_file_no_clobber "$f" ".hawp/work/decisions/$MDATE/$(basename "$f")"
    done
  fi
  mkdir -p ".hawp/work/notes/$MDATE"
  copy_dir_no_clobber "hawp/status" ".hawp/work/notes/$MDATE"
  copy_file_no_clobber "hawp/LICENSE" ".hawp/LICENSE"
  rm -rf hawp
fi

# --- 2. Migration: legacy .hawp/usage/ -> .hawp/work/ ---
if [ -d ".hawp/usage" ]; then
  copy_file_no_clobber ".hawp/usage/BACKLOG.md" ".hawp/work/BACKLOG.md"
  copy_dir_no_clobber  ".hawp/usage/status"     ".hawp/work/active"
  mkdir -p ".hawp/work/decisions/$MDATE"
  for f in .hawp/usage/*_ADR.md; do
    [ -f "$f" ] && copy_file_no_clobber "$f" ".hawp/work/decisions/$MDATE/$(basename "$f")"
  done
fi

# --- 3. Migration: pre-work/ layout (.hawp/status/) -> .hawp/work/notes/YYYY/MM/DD/ ---
if [ -d ".hawp/status" ]; then
  mkdir -p ".hawp/work/notes/$MDATE"
  copy_file_no_clobber ".hawp/status/STATUS.md" ".hawp/work/STATUS.md"
  copy_dir_no_clobber ".hawp/status" ".hawp/work/notes/$MDATE"
fi

# --- 3b. Migration: current .hawp/work/adrs/ -> .hawp/work/decisions/YYYY/MM/DD/ ---
if [ -d ".hawp/work/adrs" ]; then
  mkdir -p ".hawp/work/decisions/$MDATE"
  if cp -Rn .hawp/work/adrs/. ".hawp/work/decisions/$MDATE"/ 2>/dev/null; then
    rm -rf .hawp/work/adrs
  fi
fi

# --- 3d. Reconcile closed plans from backlog (active -> closed/...) ---
reconcile_closed_plans_from_backlog

# --- 3e. Retire orphan active items (no matching BACKLOG entry) ---
retire_orphan_active_items

# --- 4. Refresh .hawp/LICENSE and .hawp/kit/** (preserves .hawp/work/) ---
rm -rf .hawp/kit
mkdir -p .hawp/kit
cp "$SRC/.hawp/LICENSE" .hawp/
cp "$SRC/.hawp/kit/README.md"             .hawp/kit/
cp "$SRC/.hawp/kit/start-here.md"         .hawp/kit/
cp "$SRC/.hawp/kit/spec.md"               .hawp/kit/
cp "$SRC/.hawp/kit/authoring-patterns.md" .hawp/kit/
cp -R "$SRC/.hawp/kit/instructions" .hawp/kit/
cp -R "$SRC/.hawp/kit/templates" .hawp/kit/
cp -R "$SRC/.hawp/kit/patterns"  .hawp/kit/
cp -R "$SRC/.hawp/kit/reviews"   .hawp/kit/
cp -R "$SRC/.hawp/kit/examples"  .hawp/kit/
cp -R "$SRC/.hawp/kit/types"      .hawp/kit/
cp -R "$SRC/.hawp/kit/usage"      .hawp/kit/
cp -R "$SRC/.hawp/kit/references" .hawp/kit/
cp -R "$SRC/.hawp/kit/standards"  .hawp/kit/

# --- 5. Cleanup: remove legacy root-level kit folders and stray docs (now under .hawp/kit/) ---
# Safe because their reusable content has been rewritten under .hawp/kit/ and
# any repo-local items have already been migrated into .hawp/work/.
rm -rf .hawp/templates .hawp/patterns .hawp/reviews .hawp/examples .hawp/types .hawp/usage
rm -f .hawp/README.md .hawp/spec.md .hawp/start-here.md .hawp/authoring-patterns.md

# --- 5b. Cleanup: remove any .gitkeep files under .hawp/ (kit no longer ships them) ---
find .hawp -name .gitkeep -type f -delete 2>/dev/null || true

# --- 6. Seed .hawp/work/ scaffold (only when missing; never overwrites) ---
mkdir -p .hawp/work/active .hawp/work/parked .hawp/work/closed .hawp/work/decisions .hawp/work/evidence .hawp/work/status .hawp/work/notes
copy_file_no_clobber "$SRC/.hawp/work/README.md"               ".hawp/work/README.md"
copy_file_no_clobber "$SRC/.hawp/work/STATUS.md"               ".hawp/work/STATUS.md"
copy_file_no_clobber "$SRC/.hawp/work/BACKLOG.md"              ".hawp/work/BACKLOG.md"
copy_file_no_clobber "$SRC/.hawp/work/active/README.md"        ".hawp/work/active/README.md"
copy_file_no_clobber "$SRC/.hawp/work/parked/README.md"        ".hawp/work/parked/README.md"
copy_file_no_clobber "$SRC/.hawp/work/closed/README.md"        ".hawp/work/closed/README.md"
copy_file_no_clobber "$SRC/.hawp/work/decisions/README.md"     ".hawp/work/decisions/README.md"
copy_file_no_clobber "$SRC/.hawp/work/evidence/README.md"      ".hawp/work/evidence/README.md"
copy_file_no_clobber "$SRC/.hawp/work/status/README.md"        ".hawp/work/status/README.md"
copy_file_no_clobber "$SRC/.hawp/work/notes/README.md"         ".hawp/work/notes/README.md"

# --- 7. Refresh overlay files (HAWP-managed; safe to overwrite) ---
mkdir -p .github/instructions .github/prompts
cp "$SRC/.github/instructions/"*.instructions.md .github/instructions/
cp "$SRC/.github/prompts/"*.prompt.md            .github/prompts/

# --- 8. Cleanup: remove stale legacy-named overlay files ---
find .github/instructions -maxdepth 1 -type f -name 'human-ai-workflow-protocol-*.instructions.md' -delete 2>/dev/null || true
find .github/prompts -maxdepth 1 -type f -name 'human-ai-workflow-protocol-*.prompt.md' -delete 2>/dev/null || true

# --- 9. Refresh copilot-instructions (HAWP-managed) ---
cp "$SRC/.github/copilot-instructions.md" .github/copilot-instructions.md

if [ -n "$TMP_DIR" ] && [ -d "$TMP_DIR" ]; then
  rm -rf "$TMP_DIR"
fi

echo "HAWP update complete."
echo "Refreshed: .hawp/LICENSE, .hawp/kit/**, .github/instructions/*, .github/prompts/*, .github/copilot-instructions.md"
echo "Preserved: .hawp/work/** (no-overwrite)"
echo "Reconciled: Done rows + Active-Work 'done'/'wont-fix' rows moved from .hawp/work/active/ when eligible (see 'reconciled (link):' and 'reconciled (id-fallback):' lines above)"
```

## Source Reference

This file is generated. Do not edit it directly.

- Workflow gate: pushes and pull requests on `main` or `dev` fail when generated guides drift from source.
- Local sync: run `npm --prefix librarian run distribution:sync` after editing `distribution/sources/` or the distribution composition code.

Generated output file:

- `distribution/generated/update-main.md`

This generated guide is built from:

- `distribution/sources/update/script.md` — shell script (authoritative source)
