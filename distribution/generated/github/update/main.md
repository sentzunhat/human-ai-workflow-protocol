> **Agent instruction:** Update the local HAWP **GitHub/Copilot provider** installation in this repository. Run the bash command block in the **Update Command (Copy/Paste)** section in a terminal from the repository root. Report terminal output with `Source:`, `Provider:`, and `Source mode:` lines as proof.

# Safety: Install and Update Principles

These are foundational safety rules that apply to both install and update operations.

## Project Work Is Never Overwritten

- `.hawp/work/` is **project-owned and must not be overwritten**.
- This includes your BACKLOG, active work, parked work, closed work, decisions, and evidence files.
- Install and update operations preserve all existing `.hawp/work/**` files, always.

## Provider Overlay Behavior

- Each guide installs **one** provider overlay from `core/providers/.<provider>/`.
- Refresh vs seed rules and paths not touched by this guide are in **Provider Boundaries** below.

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
- Only items explicitly marked done or wont-fix in the backlog are moved. Unlinked active items are left alone.
- This means an update can legitimately change `.hawp/work/active/` and `.hawp/work/closed/` in the target repo when the backlog says a plan is finished.
- All moves use the no-overwrite rule: if the destination already exists, the source file is not moved.

## Verification Before and After

- Review what the script will do before running it (read the **Install Command** or **Update Command** block first).
- Do not pipe remote guide content directly to `bash`. Optional guide-fetch helpers write a script to `/tmp` for review first.
- After running install or update, check `.hawp/work/BACKLOG.md` and `.hawp/kit/` to confirm changes.
- If something looks wrong, `.hawp/work/` is already preserved and safe.

## Privacy-Safe Path Logging

- Do not persist machine-local absolute paths in plans, evidence, status reports, or prompts.
- Avoid storing host-local prefixes such as `<user-home>/...`, `<linux-home>/...`, or `<windows-user-home>\\...` in repository artifacts.
- If command output includes absolute host paths, redact only the machine-local prefix with a placeholder (for example `<repo-root-abs>`) while preserving command identity and repo-relative path evidence.

# GitHub / Copilot Overlay Safety

This guide installs the GitHub provider pack only.

- Refreshes `.github/instructions/` and `.github/prompts/` from `core/providers/.github/` on every install and update.
- Seeds `.github/copilot-instructions.md` on first install if missing; refreshes it on update.
- Does not modify `.cursor/`, `AGENTS.md`, or `.continue/`.

# Kit and Work Boundaries (All Providers)

Every provider install/update guide refreshes the agent-neutral kit and preserves project work.

## Kit (always installed)

| Source | Target | Behavior |
|--------|--------|----------|
| `core/.hawp/kit/**` | `.hawp/kit/**` | Full refresh every install/update |
| `core/.hawp/LICENSE` | `.hawp/LICENSE` | Refreshed |
| `core/.hawp/work/` scaffold | `.hawp/work/` READMEs, `BACKLOG.md` seed | Seed only when missing |

## Project-Owned (never overwritten)

- `.hawp/work/**` — backlog, active/parked/closed work, decisions, evidence, notes

## Never Installed Downstream

- HAWP source repo's `.hawp/work/**` operating state
- `benchmark/` reference material

Provider-specific overlay boundaries are documented in the next section for this guide.

# GitHub / Copilot Provider Boundaries

This guide sets `PROVIDER=github`. Only the GitHub overlay is installed — not Cursor or Continue paths.

## Source pack

`core/providers/.github/`

## Install mapping

| Source | Installs to | Install | Update |
|--------|-------------|---------|--------|
| `instructions/*.instructions.md` | `.github/instructions/` | refresh | refresh |
| `prompts/*.prompt.md` | `.github/prompts/` | refresh | refresh |
| `copilot-instructions.md` | `.github/copilot-instructions.md` | seed if missing | refresh |

## Not touched by this guide

- `.cursor/**`, `AGENTS.md`, `.continue/**`
- Custom files under `.github/` outside the HAWP-managed paths above

## Boundary model

```
core/providers/.github/  →  .github/instructions/
                         →  .github/prompts/
                         →  .github/copilot-instructions.md
```

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

# GitHub / Copilot Update Contract

## Work item goal

Refresh HAWP kit plus **GitHub Copilot overlays** from the selected branch. This refreshes `core/providers/.github/` into `.github/instructions/`, `.github/prompts/`, and `.github/copilot-instructions.md`.

## Agent execution

- Requires existing `.hawp/`.
- Report proof: `Source:`, `Provider: github`, `Source mode:`.
- File proof: `git status --short .hawp/kit .github/instructions .github/prompts .github/copilot-instructions.md`

## Provider-specific rules

- Refreshes `.github/copilot-instructions.md` on every update.
- Does not modify `.cursor/` or `AGENTS.md`.

## Auto-dispatch

Use the install-contract **Guide fetch (review-first)** block with `PROVIDER="github"`; it selects update when `.hawp/` exists.

# Update HAWP — GitHub/Copilot Provider (Main Branch)

Upgrade an existing HAWP installation to the latest stable kit and GitHub Copilot overlays from `main`.

## Prerequisites

- HAWP already installed (`.hawp/kit/` and `.hawp/work/` present).
- `curl` and `tar`.

## Before You Update

1. **Check `.hawp/work/BACKLOG.md`** — reconciliation may move done items from `active/` to `closed/`.
2. **Review `.github/copilot-instructions.md`** — update refreshes this file from the HAWP GitHub provider pack.

## Update Steps

1. **Open your repository root in a terminal.**

2. **Copy and run the command** from the "Update Command (Copy/Paste)" section below (`REF="main"`).

   Optional: `export HAWP_LOCAL_CORE="/absolute/path/to/human-ai-workflow-protocol/core"`.

3. **Review output** — `.hawp/kit/` refreshed; reconciliation lines if any; final `Refreshed:` summary.

4. **Verify `.hawp/work/`** is intact.

## What Gets Updated

- `.hawp/LICENSE`, `.hawp/kit/**`
- `.github/instructions/*.instructions.md`, `.github/prompts/*.prompt.md`
- `.github/copilot-instructions.md` (refreshed on update)
- Missing `.hawp/work/` scaffold files (seeded only when absent)

## What Is Preserved

- All of `.hawp/work/**`
- Your project code and configuration

## Troubleshooting

- `.hawp/work/` is never overwritten.
- Re-run update safely if needed.
- Source pack: `core/providers/.github/` in the HAWP repository.

## Other branches

- Development update: `distribution/generated/github/update/development.md`

## Update Command (Copy/Paste)

Run this from the root of your target repository. No edits are required; branch and provider are already configured in the command. Each run fetches the latest commit from that branch.

```bash
set -euo pipefail

OWNER="sentzunhat"
REPO="human-ai-workflow-protocol"
REF="main"
PROVIDER="github"

echo "Source: ${OWNER}/${REPO}@${REF}"
echo "Provider: ${PROVIDER}"
echo "Archive: https://github.com/${OWNER}/${REPO}/archive/refs/heads/${REF}.tar.gz"

TMP_DIR=""
if [ -n "${HAWP_LOCAL_CORE:-}" ]; then
  SRC="${HAWP_LOCAL_CORE}"
  if [ ! -d "$SRC/.hawp/kit" ]; then
    echo "Error: HAWP_LOCAL_CORE must point to a core directory containing .hawp/kit"
    exit 1
  fi
  echo "Source mode: local core (${SRC})"
else
  TMP_DIR="$(mktemp -d)"
  curl -fsSL "https://github.com/${OWNER}/${REPO}/archive/refs/heads/${REF}.tar.gz" \
    | tar -xz -C "$TMP_DIR"
  SRC="$(find "$TMP_DIR" -maxdepth 2 -type d -path "*/core" | head -n 1)"
  if [ -z "$SRC" ]; then
    echo "Error: downloaded archive did not contain core/"
    exit 1
  fi
  echo "Source mode: remote archive"
fi

if [ ! -d ".hawp" ]; then
  echo "Preflight: .hawp/ not found in this repository."
  echo "Run the matching install guide first, then rerun this update guide."
  if [ -n "$TMP_DIR" ] && [ -d "$TMP_DIR" ]; then
    rm -rf "$TMP_DIR"
  fi
  exit 1
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
MDATE="$(date +%Y/%m/%d)"

# --- 1. Migration: legacy hawp/ -> .hawp/ ---
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

# --- 3. Migration: .hawp/status/ -> .hawp/work/notes/ ---
if [ -d ".hawp/status" ]; then
  mkdir -p ".hawp/work/notes/$MDATE"
  copy_file_no_clobber ".hawp/status/STATUS.md" ".hawp/work/STATUS.md"
  copy_dir_no_clobber ".hawp/status" ".hawp/work/notes/$MDATE"
fi

# --- 3b. Migration: .hawp/work/adrs/ -> decisions/ ---
if [ -d ".hawp/work/adrs" ]; then
  mkdir -p ".hawp/work/decisions/$MDATE"
  if cp -Rn .hawp/work/adrs/. ".hawp/work/decisions/$MDATE"/ 2>/dev/null; then
    rm -rf .hawp/work/adrs
  fi
fi

reconcile_closed_plans_from_backlog

# --- 4. Refresh .hawp/LICENSE and .hawp/kit/** ---
rm -rf .hawp/kit
mkdir -p .hawp/kit
cp "$SRC/.hawp/LICENSE" .hawp/
cp "$SRC/.hawp/kit/README.md"             .hawp/kit/
cp "$SRC/.hawp/kit/start-here.md"         .hawp/kit/
cp -R "$SRC/.hawp/kit/instructions" .hawp/kit/
cp -R "$SRC/.hawp/kit/templates" .hawp/kit/
cp -R "$SRC/.hawp/kit/patterns"  .hawp/kit/
cp -R "$SRC/.hawp/kit/reviews"   .hawp/kit/
cp -R "$SRC/.hawp/kit/examples"  .hawp/kit/
cp -R "$SRC/.hawp/kit/types"      .hawp/kit/
cp -R "$SRC/.hawp/kit/usage"      .hawp/kit/
cp -R "$SRC/.hawp/kit/references" .hawp/kit/
cp -R "$SRC/.hawp/kit/standards"  .hawp/kit/

# --- 4b. Update hawp CLI binary (platform-detected from GitHub release) ---
# Downloads the latest Go binary to .hawp/bin/hawp-bin (beside the shell
# wrapper at .hawp/bin/hawp). Never copies the shell wrapper over the binary.
update_hawp_binary() {
  local _os _arch _asset _ext _dest _url _checksum_url _expected _actual
  local _tag

  _os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  _arch="$(uname -m)"
  _ext=""

  case "$_os" in
    linux)  _os="linux"  ;;
    darwin) _os="darwin" ;;
    mingw*|msys*|cygwin*|windows_nt*) _os="windows"; _ext=".exe" ;;
    *)
      echo "hawp update: unsupported OS '$_os' — skipping binary update."
      return 0
      ;;
  esac

  case "$_arch" in
    x86_64|amd64)  _arch="amd64" ;;
    aarch64|arm64) _arch="arm64" ;;
    *)
      echo "hawp update: unsupported arch '$_arch' — skipping binary update."
      return 0
      ;;
  esac

  _asset="hawp-${_os}-${_arch}${_ext}"
  _dest=".hawp/bin/hawp-bin${_ext}"

  _tag="$(curl -fsSL "https://api.github.com/repos/${OWNER}/${REPO}/releases/latest" 2>/dev/null \
    | awk -F'"' '/"tag_name"/ { print $4; exit }' || true)"
  if [ -z "$_tag" ]; then
    _tag="$(curl -fsSL "https://api.github.com/repos/${OWNER}/${REPO}/releases?per_page=1" 2>/dev/null \
      | awk -F'"' '/"tag_name"/ { print $4; exit }' || true)"
    [ -n "$_tag" ] && echo "hawp update: /releases/latest unavailable; using releases list fallback."
  fi

  if [ -z "$_tag" ]; then
    echo "hawp update: could not resolve latest release tag — skipping binary update."
    return 0
  fi

  _url="https://github.com/${OWNER}/${REPO}/releases/download/${_tag}/${_asset}"
  _checksum_url="https://github.com/${OWNER}/${REPO}/releases/download/${_tag}/checksums.txt"

  echo "hawp binary: ${_asset} (release ${_tag})"
  mkdir -p .hawp/bin
  curl -fsSL -o "${_dest}.tmp" "$_url" || {
    echo "hawp update: download failed — skipping binary update."
    return 0
  }

  if curl -fsSL -o /tmp/hawp-checksums.txt "$_checksum_url" 2>/dev/null; then
    _expected="$(grep " ${_asset}$" /tmp/hawp-checksums.txt | awk '{print $1}')"
    if [ -n "$_expected" ]; then
      if command -v sha256sum >/dev/null 2>&1; then
        _actual="$(sha256sum "${_dest}.tmp" | awk '{print $1}')"
      elif command -v shasum >/dev/null 2>&1; then
        _actual="$(shasum -a 256 "${_dest}.tmp" | awk '{print $1}')"
      else
        _actual=""
      fi
      if [ -n "$_actual" ] && [ "$_actual" != "$_expected" ]; then
        echo "hawp update: SHA256 mismatch — aborting binary update."
        rm -f "${_dest}.tmp"
        return 1
      fi
      [ -n "$_actual" ] && echo "hawp binary: SHA256 verified."
    fi
    rm -f /tmp/hawp-checksums.txt
  fi

  mv "${_dest}.tmp" "$_dest"
  chmod +x "$_dest"

  # Install the shell wrapper alongside the binary (idempotent).
  if [ -f "$SRC/.hawp/bin/hawp" ]; then
    cp "$SRC/.hawp/bin/hawp" .hawp/bin/hawp
    chmod +x .hawp/bin/hawp
  fi
  if [ -f "$SRC/.hawp/bin/hawp-mcp" ]; then
    cp "$SRC/.hawp/bin/hawp-mcp" .hawp/bin/hawp-mcp
    chmod +x .hawp/bin/hawp-mcp
  fi

  echo "hawp binary: updated to ${_tag} at ${_dest}"
}
update_hawp_binary

rm -rf .hawp/templates .hawp/patterns .hawp/reviews .hawp/examples .hawp/types .hawp/usage
rm -f .hawp/README.md .hawp/spec.md .hawp/start-here.md .hawp/authoring-patterns.md
find .hawp -name .gitkeep -type f -delete 2>/dev/null || true

# --- 5. Seed .hawp/work/ scaffold (only when missing) ---
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

# --- Provider overlay: GitHub/Copilot (refresh) ---
resolve_provider_pack() {
  if [ -d "$SRC/providers/.github" ]; then
    echo "$SRC/providers/.github"
    return 0
  fi
  echo "Error: GitHub provider pack not found at core/providers/.github/" >&2
  return 1
}
update_provider_overlay() {
  pack="$(resolve_provider_pack)" || return 1
  mkdir -p .github/instructions .github/prompts
  cp "$pack/instructions/"*.instructions.md .github/instructions/
  cp "$pack/prompts/"*.prompt.md            .github/prompts/
  find .github/instructions -maxdepth 1 -type f -name 'human-ai-workflow-protocol-*.instructions.md' -delete 2>/dev/null || true
  find .github/prompts -maxdepth 1 -type f -name 'human-ai-workflow-protocol-*.prompt.md' -delete 2>/dev/null || true
  cp "$pack/copilot-instructions.md" .github/copilot-instructions.md
  echo "  refreshed: core/providers/.github/ -> .github/"
}
update_provider_overlay || exit 1
echo "Provider overlay: .github/instructions/*, .github/prompts/*, .github/copilot-instructions.md"

if [ -n "$TMP_DIR" ] && [ -d "$TMP_DIR" ]; then
  rm -rf "$TMP_DIR"
fi

echo "HAWP update complete (provider: ${PROVIDER})."
echo "Refreshed: .hawp/LICENSE, .hawp/kit/**"
echo "Preserved: .hawp/work/** (no-overwrite)"
echo "Reconciled: Done rows + Active-Work 'done'/'wont-fix' rows moved from .hawp/work/active/ when eligible (see 'reconciled (link):' and 'reconciled (id-fallback):' lines above)"
```

## Source Reference

This file is generated. Do not edit it directly.

- Workflow gate: pushes and pull requests on `main` or `development` fail when generated guides drift from source.
- Local sync: run `hawp distribution sync` after editing `distribution/sources/` or the distribution composition code.

Generated output file:

- `distribution/generated/github/update/main.md`

Provider: `github` · Operation: `update` · Branch: `main`

Install mapping: `core/providers/.github/` -> downstream paths in this guide.

This generated guide is built from:

- `distribution/sources/providers/github/preamble-update.md`
- `distribution/sources/shared/safety.md`
- `distribution/sources/providers/github/safety.md`
- `distribution/sources/shared/repo-boundaries-kit.md`
- `distribution/sources/providers/github/boundaries.md`
- `distribution/sources/shared/update.md`
- `distribution/sources/providers/github/update-contract.md`
- `distribution/sources/providers/github/update/main.md`

Composed shell script (core + provider overlay + footer):

- `distribution/sources/update/script-core.md`
- `distribution/sources/providers/github/script-update.md`
- `distribution/sources/update/script-footer.md`
