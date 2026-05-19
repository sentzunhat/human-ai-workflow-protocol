> **Agent instruction:** Install HAWP in this repository. Run the bash command block in the **Install Command (Copy/Paste)** section in a terminal from the repository root. Do not report "already present" based on content inspection alone — execute the script and report terminal output with `Source:` and `Source mode:` lines as proof.

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

# Install HAWP From Dev Branch

Use this installation path when you want to test the latest unreleased or in-progress HAWP changes and improvements.

## When to Use Dev Install

- You are actively contributing to HAWP development.
- You want early access to new features or bug fixes in progress.
- You are testing HAWP improvements before they land on `main`.
- You want to provide feedback on new HAWP guidance or tools.

## Prerequisites

- A repository where you want to test HAWP.
- `curl` and `tar` (standard Unix tools).
- Ability to run shell scripts or copy/paste a shell command block.
- Understanding that dev branch changes may not be fully stable.

## Installation Steps

1. **Open your target repository in a terminal and move to its root directory.**

2. **Copy and run the command** from the "Install Command (Copy/Paste)" section below.
   - No edits are needed.
   - This guide already pins `REF="dev"`.
   - The script output will show `Source: sentzunhat/human-ai-workflow-protocol@dev`.
   - The command handles migration + install in one run.

   Optional local-source mode for branch testing:
   - If you need to test local or unpushed branch changes, run `export HAWP_LOCAL_CORE="/absolute/path/to/human-ai-workflow-protocol/core"` before running the command.
   - The script will print `Source mode: local core (...)` when that override is active.

3. **Review the output** to confirm files were added from the dev branch.
   - Check that `.hawp/kit/` was created with dev-branch HAWP documentation.
   - Check that `.github/instructions/` and `.github/prompts/` have dev-branch files.
   - Check that `.hawp/work/BACKLOG.md` exists.

4. **Test the new features or changes**:
   - Read the HAWP kit to understand any new concepts or improvements.
   - Test the new features in your workflow.
   - Open issues or feedback on the HAWP repository if you find problems.

## What Was Added

Same as main branch install, but with dev-branch versions of:

- `.hawp/kit/**` (dev-branch protocol docs, templates, patterns).
- `.github/instructions/` and `.github/prompts/` (dev-branch Copilot overlays).
- Any other HAWP-managed files.

## Important: Dev Branch Is Unstable

- Dev branch changes are in progress and may not be fully tested.
- You may encounter bugs or incomplete features.
- Do not use dev branch for production unless you are actively testing and providing feedback.
- Switch back to `main` branch if you need stability.

## Switching Back to Main

To return to the stable `main` branch:

1. Use the main branch install guide (`install-main.md`) and run that script.
2. All HAWP-managed files will be updated to main branch versions.
3. Your `.hawp/work/` and `.github/copilot-instructions.md` will be preserved.

## Install Command (Copy/Paste)

Run this from the root of your target repository. No edits are required; branch selection is already configured in the command. Each run fetches the latest commit from that branch.

```bash
OWNER="sentzunhat"
REPO="human-ai-workflow-protocol"
REF="dev"

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

if [ -d ".hawp" ]; then
  echo "Preflight: detected existing .hawp/."
  echo "Switching to update-compatible refresh mode for this run."
  echo "Tip: use update-${REF}.md next time when .hawp/ already exists."
fi

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

MDATE="$(date +%Y/%m/%d)"

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

# --- 7. Install Copilot overlay into .github/ ---
mkdir -p .github/instructions .github/prompts
cp "$SRC/.github/instructions/"*.instructions.md .github/instructions/
cp "$SRC/.github/prompts/"*.prompt.md            .github/prompts/

# --- 8. Cleanup: remove stale legacy-named overlay files ---
find .github/instructions -maxdepth 1 -type f -name 'human-ai-workflow-protocol-*.instructions.md' -delete 2>/dev/null || true
find .github/prompts -maxdepth 1 -type f -name 'human-ai-workflow-protocol-*.prompt.md' -delete 2>/dev/null || true

# --- 9. Install copilot-instructions (HAWP-managed) ---
cp "$SRC/.github/copilot-instructions.md" .github/copilot-instructions.md

if [ -n "$TMP_DIR" ] && [ -d "$TMP_DIR" ]; then
  rm -rf "$TMP_DIR"
fi

echo "HAWP install complete."
echo "Refreshed: .hawp/LICENSE, .hawp/kit/**, .github/instructions/*, .github/prompts/*, .github/copilot-instructions.md"
echo "Preserved: .hawp/work/** (no-overwrite)"
echo "Reconciled: Done rows + Active-Work 'done'/'wont-fix' rows moved from .hawp/work/active/ when eligible (see 'reconciled (link):' and 'reconciled (id-fallback):' lines above)"
```

## Source Reference

This file is generated. Do not edit it directly.

- Workflow gate: pushes and pull requests on `main` or `dev` fail when generated guides drift from source.
- Local sync: run `npm --prefix librarian run distribution:sync` after editing `distribution/sources/` or the distribution composition code.

Generated output file:

- `distribution/generated/install-dev.md`

This generated guide is built from:

- `distribution/sources/install/script.md` — shell script (authoritative source)
