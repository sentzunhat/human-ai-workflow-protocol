# GitHub / Copilot Install Contract

## Work item goal

Install HAWP kit plus **GitHub Copilot overlays only**. This refreshes `core/providers/.github/` into `.github/instructions/`, `.github/prompts/`, and `.github/copilot-instructions.md`.

## Agent execution

- Run the **Install Command** bash block in a terminal from repo root.
- Report proof lines: `Source:`, `Provider: github`, `Source mode:`.
- File proof:
  - `git status --short .hawp/LICENSE .hawp/kit .github/instructions .github/prompts`
  - `find .hawp/kit -maxdepth 2 -type f | head -n 20`

## Provider-specific rules

- Seed `.github/copilot-instructions.md` on install only when missing.
- Refresh `.github/instructions/` and `.github/prompts/` every install.
- Do **not** expect `.cursor/` or `AGENTS.md` changes from this guide.

## Guide fetch (review-first)

**Recommended:** use **Install Command (Copy/Paste)** below — the full script is visible in this guide for review before you run it.

**Optional convenience:** downloads the remote guide and writes the extracted command block to a local script file. **Does not execute automatically** — review the file, then run it explicitly.

> Security: do not pipe remote content directly to `bash`. This helper writes to `/tmp` so you can inspect the script first.

````bash
OWNER="sentzunhat"
REPO="human-ai-workflow-protocol"
PROVIDER="github"
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
