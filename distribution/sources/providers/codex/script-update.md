```bash
# --- Provider overlay: Codex (refresh AGENTS.md) ---
resolve_provider_pack() {
  if [ -d "$SRC/providers/.codex" ]; then
    echo "$SRC/providers/.codex"
    return 0
  fi
  echo "Error: Codex provider pack not found at core/providers/.codex/" >&2
  return 1
}
update_provider_overlay() {
  pack="$(resolve_provider_pack)" || return 1
  cp "$pack/AGENTS.md.seed" AGENTS.md
  echo "  refreshed: core/providers/.codex/ -> AGENTS.md"
}
update_provider_overlay || exit 1
echo "Provider overlay: AGENTS.md"
```

