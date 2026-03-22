```bash
# --- Provider overlay: Codex (core/providers/.codex/ -> AGENTS.md) ---
resolve_provider_pack() {
  if [ -d "$SRC/providers/.codex" ]; then
    echo "$SRC/providers/.codex"
    return 0
  fi
  echo "Error: Codex provider pack not found at core/providers/.codex/" >&2
  return 1
}
install_provider_overlay() {
  pack="$(resolve_provider_pack)" || return 1
  copy_file_no_clobber "$pack/AGENTS.md.seed" AGENTS.md
  echo "  installed: core/providers/.codex/ -> AGENTS.md (seed if missing)"
}
install_provider_overlay || exit 1
echo "Provider overlay: AGENTS.md (seed if missing)"
```

