#!/usr/bin/env bash
# Prepare a temporary copy of this repository with agent-loaded HAWP guidance removed.
# Use the copy for the no-HAWP benchmark arm. Never modifies the source repository.
#
# Usage:
#   ./benchmark/prepare-clean-workspace.sh
#   ./benchmark/prepare-clean-workspace.sh --open
#   ./benchmark/prepare-clean-workspace.sh --dest /tmp/my-benchmark --task "Your bare prompt here"
#   ./benchmark/prepare-clean-workspace.sh --cleanup

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

OPEN_CURSOR=0
CLEANUP=0
DEST=""
TASK=""

usage() {
  cat <<'EOF'
Prepare a clean workspace for the no-HAWP benchmark arm.

Copies this repository to a temporary folder and removes files that Cursor,
Continue, and GitHub Copilot load automatically (AGENTS.md, rules, instructions).
The original repository is never modified.

Options:
  --open              Open the clean copy in a new Cursor window
  --dest PATH         Write the copy to PATH (default: /tmp/hawp-benchmark-clean-<timestamp>)
  --task "TEXT"       Write the bare prompt to BENCHMARK-TASK.txt in the copy
  --cleanup           Remove /tmp/hawp-benchmark-clean-* folders
  -h, --help          Show this help

After the script runs:
  1. Open the printed path as its own workspace (or use --open).
  2. In that window, follow benchmark/instructions/run.md through the STOP step.
  3. Save no-hawp/output.md in the real repo, then close the clean window.
  4. In the original repo, follow benchmark/instructions/cleanup.md for the HAWP arm + cleanup.

See benchmark/instructions/setup.md and benchmark/benchmark-prompt.md.
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --open) OPEN_CURSOR=1; shift ;;
    --dest)
      DEST="${2:-}"
      if [[ -z "$DEST" ]]; then
        echo "error: --dest requires a path" >&2
        exit 1
      fi
      shift 2
      ;;
    --task)
      TASK="${2:-}"
      if [[ -z "$TASK" ]]; then
        echo "error: --task requires text" >&2
        exit 1
      fi
      shift 2
      ;;
    --cleanup)
      CLEANUP=1
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "error: unknown option: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

if [[ "$CLEANUP" -eq 1 ]]; then
  shopt -s nullglob
  removed=0
  for dir in /tmp/hawp-benchmark-clean-*; do
    rm -rf "$dir"
    echo "removed: $dir"
    removed=$((removed + 1))
  done
  if [[ "$removed" -eq 0 ]]; then
    echo "nothing to remove under /tmp/hawp-benchmark-clean-*"
  fi
  exit 0
fi

if [[ -z "$DEST" ]]; then
  DEST="/tmp/hawp-benchmark-clean-$(date +%Y%m%d-%H%M%S)"
fi

if [[ -e "$DEST" ]]; then
  echo "error: destination already exists: $DEST" >&2
  echo "Use --dest with a new path, or remove it first." >&2
  exit 1
fi

mkdir -p "$DEST"

echo "Copying repository to: $DEST"
echo "(excluding node_modules for speed; .git is included so git status still works)"

if command -v rsync >/dev/null 2>&1; then
  rsync -a \
    --exclude 'node_modules' \
    --exclude '.git/objects/pack' \
    "${REPO_ROOT}/" "${DEST}/"
else
  cp -a "${REPO_ROOT}/." "${DEST}/"
  rm -rf "${DEST}/node_modules" 2>/dev/null || true
fi

# Paths that editors load automatically. Removed in the copy only.
STRIP_PATHS=(
  "AGENTS.md"
  ".cursor/rules"
  ".continue/rules"
  ".github/instructions"
  "core/providers/.cursor/rules"
  "core/providers/.continue/rules"
  "core/providers/.github/instructions"
)

echo
echo "Removing agent-loaded guidance from the copy:"
for rel in "${STRIP_PATHS[@]}"; do
  target="${DEST}/${rel}"
  if [[ -e "$target" ]]; then
    rm -rf "$target"
    echo "  removed: ${rel}"
  fi
done

cat > "${DEST}/.benchmark-clean-workspace" <<EOF
HAWP benchmark clean workspace
Created: $(date -u +"%Y-%m-%dT%H:%M:%SZ")
Source: ${REPO_ROOT}

This folder is a temporary copy for the no-HAWP benchmark arm.
Agent-loaded guidance was removed from the copy only; the source repo was not modified.

Stripped paths:
$(printf '  - %s\n' "${STRIP_PATHS[@]}")

Note: .hawp/kit/ and other docs remain readable if the agent explores the tree.
Only always-on injected rules were removed.
EOF

if [[ -n "$TASK" ]]; then
  cat > "${DEST}/BENCHMARK-TASK.txt" <<EOF
Paste this into a fresh agent chat in this workspace only (no-HAWP arm):

${TASK}

Return your complete answer as your final response.
EOF
  echo
  echo "Wrote bare prompt to: ${DEST}/BENCHMARK-TASK.txt"
fi

echo
echo "Clean workspace ready."
echo
cat > "${DEST}/BENCHMARK-NEXT-STEPS.md" <<'NEXT'
# No-HAWP arm — next steps

You are in the clean benchmark workspace. Follow these instructions:

  benchmark/instructions/run.md

Do NOT run the HAWP arm here. After saving no-hawp/output.md in the real repo,
close this window and continue with cleanup.md in the original repository
(it runs the HAWP arm, scoring, recording, then cleanup).
NEXT

echo "Next steps (clean workspace window):"
echo "  1. Open this folder as its own workspace:"
echo "     ${DEST}"
if [[ "$OPEN_CURSOR" -eq 1 ]]; then
  if command -v cursor >/dev/null 2>&1; then
    cursor -n "$DEST"
    echo "  (opened in a new Cursor window)"
  else
    echo "  warning: cursor CLI not found; open the path manually"
  fi
else
  echo "     or run: ./benchmark/prepare-clean-workspace.sh --open --dest \"${DEST}\""
fi
echo "  2. Follow: benchmark/instructions/run.md"
echo "  3. Save no-hawp/output.md in the real repo, then CLOSE this window."
echo
echo "Then (original repo window):"
echo "  4. Follow: benchmark/instructions/cleanup.md (HAWP arm + score + cleanup)"
echo "     Source repo: ${REPO_ROOT}"
echo
echo "Pointer file in this copy: BENCHMARK-NEXT-STEPS.md"
echo "When finished: ./benchmark/prepare-clean-workspace.sh --cleanup"
