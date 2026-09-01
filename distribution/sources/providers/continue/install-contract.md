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
REF="main"   # set to "development" to install from the dev branch

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
