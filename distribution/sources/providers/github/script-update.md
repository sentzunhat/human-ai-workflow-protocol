```bash
# --- Provider overlay: GitHub/Copilot (refresh) ---
resolve_provider_pack() {
  if [ -d "$SRC/providers/.github" ]; then
    echo "$SRC/providers/.github"
    return 0
  fi
  echo "Error: GitHub provider pack not found at core/providers/.github/" >&2
  return 1
}
update_provider_overlay() {
  pack="$(resolve_provider_pack)" || return 1
  mkdir -p .github/instructions .github/prompts
  cp "$pack/instructions/"*.instructions.md .github/instructions/
  cp "$pack/prompts/"*.prompt.md            .github/prompts/
  for existing in .github/instructions/human-ai-workflow-protocol-*.instructions.md; do
    [ -e "$existing" ] || continue
    rm -f "$existing"
  done
  for existing in .github/prompts/human-ai-workflow-protocol-*.prompt.md; do
    [ -e "$existing" ] || continue
    rm -f "$existing"
  done
  cp "$pack/copilot-instructions.md" .github/copilot-instructions.md
  echo "  refreshed: core/providers/.github/ -> .github/"
}
update_provider_overlay || exit 1
echo "Provider overlay: .github/instructions/*, .github/prompts/*, .github/copilot-instructions.md"
```
