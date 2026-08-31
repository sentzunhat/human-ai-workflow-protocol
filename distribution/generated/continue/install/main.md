> **Agent instruction:** Install HAWP with the **Continue provider** in this repository. Run the bash command block in the **Install Command (Copy/Paste)** section in a terminal from the repository root. Do not report "already present" based on content inspection alone — execute the script and report terminal output with `Source:`, `Provider:`, and `Source mode:` lines as proof.

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

# Continue Overlay Safety

This guide installs the Continue provider pack only.

- Refreshes `.continue/rules/hawp-*.md` from `core/providers/.continue/rules/` on every install and update.
- Does not modify `.github/`, `.cursor/`, or `AGENTS.md`.
- Non-HAWP rules already in `.continue/rules/` are left unchanged.

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

# Continue Provider Boundaries

This guide sets `PROVIDER=continue`. Only the Continue overlay is installed — not GitHub or Cursor paths.

## Source pack

`core/providers/.continue/`

## Install mapping

| Source | Installs to | Install | Update |
|--------|-------------|---------|--------|
| `rules/hawp-*.md` | `.continue/rules/` | refresh | refresh |

## Not touched by this guide

- `.github/**`
- `.cursor/**`, `AGENTS.md`
- Non-HAWP rules in `.continue/rules/` (only `hawp-*.md` from the provider pack are refreshed)

## Boundary model

```
core/providers/.continue/  →  .continue/rules/hawp-*.md
```

Local Continue rules load automatically; no Hub `config.yaml` is required for this pack.

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
- Does not build the search index for you. If you plan to use `hawp search` or
  `hawp_search` right away, run `hawp search index` after install. `hawp search embed`
  is optional and only needed for semantic or hybrid search.

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

# Continue Install Contract

## Work item goal

Install HAWP kit plus **Continue overlays only**. This refreshes `core/providers/.continue/` into `.continue/rules/`.

## Agent execution

- Run the **Install Command** bash block in a terminal from repo root.
- Report proof lines: `Source:`, `Provider: continue`, `Source mode:`.
- File proof:
  - `git status --short .hawp/LICENSE .hawp/kit .continue/rules`
  - `find .continue/rules -maxdepth 1 -name 'hawp-*.md' 2>/dev/null | sort`

## Provider-specific rules

- Refresh all `hawp-*.md` rules from the provider pack on every install.
- Do **not** expect `.github/` or `.cursor/` changes from this guide.

## Guide fetch (review-first)

**Recommended:** use **Install Command (Copy/Paste)** below — the full script is visible in this guide for review before you run it.

**Optional convenience:** downloads the remote guide and writes the extracted command block to a local script file. **Does not execute automatically** — review the file, then run it explicitly.

> Security: do not pipe remote content directly to `bash`. This helper writes to `/tmp` so you can inspect the script first.

````bash
OWNER="sentzunhat"
REPO="human-ai-workflow-protocol"
PROVIDER="continue"
REF="dev"   # set to "main" for stable

case "$REF" in
  main|dev) ;;
  *) echo "Error: REF must be 'main' or 'dev'"; exit 1 ;;
esac

if [ -d ".hawp" ]; then MODE="update"; else MODE="install"; fi

GUIDE_URL="https://raw.githubusercontent.com/${OWNER}/${REPO}/${REF}/distribution/generated/${PROVIDER}/${MODE}/${REF}.md"
SCRIPT="$(mktemp "/tmp/hawp-${PROVIDER}-${MODE}.XXXXXX.sh")"

echo "Execution work item: ${MODE} HAWP (${PROVIDER}) from ${OWNER}/${REPO}@${REF}"
echo "Guide: ${GUIDE_URL}"

curl -fsSL "$GUIDE_URL" | awk '
  /^```bash$/ { in_block = 1; next }
  /^```$/ && in_block { exit }
  in_block { print }
' > "$SCRIPT"

chmod 700 "$SCRIPT"
echo "Script written: $SCRIPT"
echo "Review it, then run:"
echo "  bash \"$SCRIPT\""
````

# Install HAWP — Continue Provider (Main Branch)

Stable install of HAWP kit plus Continue overlays: `.continue/rules/hawp-*.md`.

**Source → target mapping:**

| `core/providers/.continue/` | Your repo |
|-----------------------------|-----------|
| `rules/hawp-*.md` | `.continue/rules/` |

## Prerequisites

- A repository where you use Continue (VS Code extension or IDE integration).
- `curl` and `tar`.

## Installation Steps

1. Open your target repository root in a terminal.
2. Run the **Install Command (Copy/Paste)** block below (`REF="main"`, `PROVIDER="continue"`).
3. Confirm `.hawp/kit/` and `.continue/rules/hawp-*.md` exist.
4. Open Continue and verify HAWP rules appear in the rules toolbar.

Optional: `export HAWP_LOCAL_CORE="/absolute/path/to/human-ai-workflow-protocol/core"` for local testing.

## What Was Added

- `.hawp/kit/**` — agent-neutral HAWP kit (always installed).
- `.continue/rules/hawp-*.md` — Continue rules from `core/providers/.continue/rules/`.
- `.hawp/work/` scaffold — seeded once when missing.

## What Was NOT Changed

- Non-HAWP files under `.continue/rules/` (only `hawp-*.md` from the provider pack are copied).
- `.hawp/work/**` project records.

## Other guides

- Dev branch: `distribution/generated/continue/install/dev.md`
- GitHub/Copilot: `distribution/generated/github/install/main.md`
- Cursor: `distribution/generated/cursor/install/main.md`

## Install Command (Copy/Paste)

Run this from the root of your target repository. No edits are required; branch and provider are already configured in the command. Each run fetches the latest commit from that branch.

```bash
set -euo pipefail

OWNER="sentzunhat"
REPO="human-ai-workflow-protocol"
REF="main"
PROVIDER="continue"

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
  SRC="$TMP_DIR/${REPO}-${REF}/core"
  echo "Source mode: remote archive"
fi

if [ -d ".hawp" ]; then
  echo "Preflight: detected existing .hawp/."
  echo "Switching to update-compatible refresh mode for this run."
  echo "Tip: use distribution/generated/${PROVIDER}/update/${REF}.md next time when .hawp/ already exists."
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

# --- 4b. Install hawp CLI binary (platform-detected from GitHub release) ---
install_hawp_binary() {
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
      echo "hawp install: unsupported OS '$_os' — skipping binary install."
      return 0
      ;;
  esac

  case "$_arch" in
    x86_64|amd64)  _arch="amd64" ;;
    aarch64|arm64) _arch="arm64" ;;
    *)
      echo "hawp install: unsupported arch '$_arch' — skipping binary install."
      return 0
      ;;
  esac

  _asset="hawp-${_os}-${_arch}${_ext}"
  _dest=".hawp/bin/hawp-bin${_ext}"

  # Resolve latest release tag from GitHub API.
  _tag="$(curl -fsSL "https://api.github.com/repos/${OWNER}/${REPO}/releases/latest" 2>/dev/null \
    | awk -F'"' '/"tag_name"/ { print $4; exit }' || true)"
  if [ -z "$_tag" ]; then
    _tag="$(curl -fsSL "https://api.github.com/repos/${OWNER}/${REPO}/releases?per_page=1" 2>/dev/null \
      | awk -F'"' '/"tag_name"/ { print $4; exit }' || true)"
    [ -n "$_tag" ] && echo "hawp install: /releases/latest unavailable; using releases list fallback."
  fi

  if [ -z "$_tag" ]; then
    echo "hawp install: could not resolve latest release tag — skipping binary install."
    return 0
  fi

  _url="https://github.com/${OWNER}/${REPO}/releases/download/${_tag}/${_asset}"
  _checksum_url="https://github.com/${OWNER}/${REPO}/releases/download/${_tag}/checksums.txt"

  echo "hawp binary: ${_asset} (release ${_tag})"
  mkdir -p .hawp/bin
  curl -fsSL -o "${_dest}.tmp" "$_url" || {
    echo "hawp install: download failed — skipping binary install."
    return 0
  }

  # Verify SHA256 when checksums.txt is available
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
        echo "hawp install: SHA256 mismatch — aborting binary install."
        rm -f "${_dest}.tmp"
        return 1
      fi
      [ -n "$_actual" ] && echo "hawp binary: SHA256 verified."
    fi
    rm -f /tmp/hawp-checksums.txt
  fi

  mv "${_dest}.tmp" "$_dest"
  chmod +x "$_dest"

  # Install the shell wrapper at .hawp/bin/hawp so it delegates to hawp-bin.
  if [ -f "$SRC/.hawp/bin/hawp" ]; then
    cp "$SRC/.hawp/bin/hawp" .hawp/bin/hawp
    chmod +x .hawp/bin/hawp
  fi
  if [ -f "$SRC/.hawp/bin/hawp-mcp" ]; then
    cp "$SRC/.hawp/bin/hawp-mcp" .hawp/bin/hawp-mcp
    chmod +x .hawp/bin/hawp-mcp
  fi

  echo "hawp binary: installed to ${_dest}"
}
install_hawp_binary

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

# --- Provider overlay: Continue (core/providers/.continue/ -> .continue/rules/) ---
resolve_provider_pack() {
  if [ -d "$SRC/providers/.continue" ]; then
    echo "$SRC/providers/.continue"
    return 0
  fi
  echo "Error: Continue provider pack not found at core/providers/.continue/" >&2
  return 1
}
install_provider_overlay() {
  pack="$(resolve_provider_pack)" || return 1
  mkdir -p .continue/rules
  if [ -d "$pack/rules" ]; then
    cp "$pack/rules"/hawp-*.md .continue/rules/ 2>/dev/null || true
  fi
  echo "  installed: core/providers/.continue/ -> .continue/rules/"
}
install_provider_overlay || exit 1
echo "Provider overlay: .continue/rules/hawp-*.md"

if [ -n "$TMP_DIR" ] && [ -d "$TMP_DIR" ]; then
  rm -rf "$TMP_DIR"
fi

echo "HAWP install complete (provider: ${PROVIDER})."
echo "Refreshed: .hawp/LICENSE, .hawp/kit/**, .hawp/bin/hawp (platform binary)"
echo "Preserved: .hawp/work/** (no-overwrite)"
echo "Reconciled: Done rows + Active-Work 'done'/'wont-fix' rows moved from .hawp/work/active/ when eligible (see 'reconciled (link):' and 'reconciled (id-fallback):' lines above)"
```

## Source Reference

This file is generated. Do not edit it directly.

- Workflow gate: pushes and pull requests on `main` or `dev` fail when generated guides drift from source.
- Local sync: run `hawp distribution sync` after editing `distribution/sources/` or the distribution composition code.

Generated output file:

- `distribution/generated/continue/install/main.md`

Provider: `continue` · Operation: `install` · Branch: `main`

Install mapping: `core/providers/.continue/` -> downstream paths in this guide.

This generated guide is built from:

- `distribution/sources/providers/continue/preamble-install.md`
- `distribution/sources/shared/safety.md`
- `distribution/sources/providers/continue/safety.md`
- `distribution/sources/shared/repo-boundaries-kit.md`
- `distribution/sources/providers/continue/boundaries.md`
- `distribution/sources/shared/install.md`
- `distribution/sources/providers/continue/install-contract.md`
- `distribution/sources/providers/continue/install/main.md`

Composed shell script (core + provider overlay + footer):

- `distribution/sources/install/script-core.md`
- `distribution/sources/providers/continue/script-install.md`
- `distribution/sources/install/script-footer.md`
