# Claude Code Install Contract

## Work item goal

Install HAWP kit plus **Claude Code overlays only**. This copies `core/providers/.claude/rules/hawp-*.md` into `.claude/rules/` and seeds `CLAUDE.md` if absent.

## Agent execution

- Run the **Install Command** bash block in a terminal from repo root.
- Report proof lines: `Source:`, `Provider: claude`, `Source mode:`.
- File proof:
  - `git status --short .hawp/LICENSE .hawp/kit .claude/rules CLAUDE.md`
  - `for rule in .claude/rules/hawp-*.md; do [ -f "$rule" ] && printf '%s\n' "$rule"; done | sort`

## Provider-specific rules

- Refresh all `hawp-*.md` rules from the provider pack on every install.
- Seed `CLAUDE.md` from `CLAUDE.md.seed` only when `CLAUDE.md` is absent.
- Do **not** expect `.github/`, `.cursor/`, or `.continue/` changes from this guide.

## Guide fetch (review-first)

**Recommended:** use **Install Command (Copy/Paste)** below — the full script is visible in this guide for review before you run it.

**Optional convenience:** downloads the remote guide and writes the extracted command block to a local script file. **Does not execute automatically** — review the file, then run it explicitly.

> Security: do not pipe remote content directly to `bash`. This helper writes to `/tmp` so you can inspect the script first.

````bash
OWNER="sentzunhat"
REPO="human-ai-workflow-protocol"
PROVIDER="claude"
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
