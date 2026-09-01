# Codex Install Contract

## Work item goal

Install HAWP kit plus the **Codex overlay only**. This seeds `core/providers/.codex/AGENTS.md.seed` into `AGENTS.md` when the target repo does not already have one.

## Agent execution

- Run the **Install Command (Copy/Paste)** bash block in a terminal from repo root.
- Report proof lines: `Source:`, `Provider: codex`, `Source mode:`.
- File proof:
  - `git status --short .hawp/LICENSE .hawp/kit AGENTS.md`
  - `test -f AGENTS.md && sed -n '1,40p' AGENTS.md`

## Provider-specific rules

- Seed `AGENTS.md` on install when it is missing.
- If your repo already has `AGENTS.md`, HAWP preserves it. Manually blend in any desired HAWP guidance from the provider seed instead of overwriting your existing file.
- Do **not** expect `.github/`, `.cursor/`, `.continue/`, or `.claude/` changes from this guide.
- Do **not** create runtime CLI participant adapters.

## Guide fetch (review-first)

**Recommended:** use **Install Command (Copy/Paste)** below — the full script is visible in this guide for review before you run it.

**Optional convenience:** downloads the remote guide and writes the extracted command block to a local script file. **Does not execute automatically** — review the file, then run it explicitly.

> Security: do not pipe remote content directly to `bash`. This helper writes to `/tmp` so you can inspect the script first.

````bash
OWNER="sentzunhat"
REPO="human-ai-workflow-protocol"
PROVIDER="codex"
REF="main"   # set to "development" to install from the development branch

case "$REF" in
  main|development) ;;
  *) echo "Error: REF must be 'main' or 'development'"; exit 1 ;;
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
