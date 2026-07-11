```bash
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
```
