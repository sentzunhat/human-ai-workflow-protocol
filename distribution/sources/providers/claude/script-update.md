```bash
# --- Provider overlay: Claude Code (refresh rules only; CLAUDE.md preserved) ---
resolve_provider_pack() {
  if [ -d "$SRC/providers/.claude" ]; then
    echo "$SRC/providers/.claude"
    return 0
  fi
  echo "Error: Claude Code provider pack not found at core/providers/.claude/" >&2
  return 1
}
update_provider_overlay() {
  pack="$(resolve_provider_pack)" || return 1
  mkdir -p .claude/rules
  if [ -d "$pack/rules" ]; then
    cp "$pack/rules/"hawp-*.md .claude/rules/ 2>/dev/null || true
  fi
  echo "  refreshed: core/providers/.claude/ -> .claude/rules/ (CLAUDE.md preserved)"
}
update_provider_overlay || exit 1
echo "Provider overlay: .claude/rules/* (CLAUDE.md not touched)"
```
