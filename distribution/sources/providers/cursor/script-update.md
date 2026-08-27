```bash
# --- Provider overlay: Cursor (refresh) ---
resolve_provider_pack() {
  if [ -d "$SRC/providers/.cursor" ]; then
    echo "$SRC/providers/.cursor"
    return 0
  fi
  echo "Error: Cursor provider pack not found at core/providers/.cursor/" >&2
  return 1
}
update_provider_overlay() {
  pack="$(resolve_provider_pack)" || return 1
  mkdir -p .cursor/rules
  if [ -d "$pack/rules" ]; then
    cp "$pack/rules/"*.mdc .cursor/rules/ 2>/dev/null || true
  fi
  copy_file_no_clobber "$pack/AGENTS.md.seed" AGENTS.md
  echo "  refreshed: core/providers/.cursor/ -> .cursor/rules/"
  echo "  seeded if missing: core/providers/.cursor/ -> AGENTS.md"
}
update_provider_overlay || exit 1
echo "Provider overlay: .cursor/rules/*, AGENTS.md (seed if missing)"
```
