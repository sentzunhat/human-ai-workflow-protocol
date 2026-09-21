#!/usr/bin/env bash
set -euo pipefail

# Default: preview. Relative artifact paths resolve in this source-layout
# directory; repository reads stay anchored to the Git root that owns this
# script.
script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(git -C "$script_dir" rev-parse --show-toplevel)"
expected_root="$(cd -- "$script_dir/../.." && pwd)"

if [[ "$repo_root" != "$expected_root" ]]; then
  printf 'source-layout: script path is not under expected repository root\n' >&2
  printf '  git root: %s\n' "$repo_root" >&2
  printf '  expected: %s\n' "$expected_root" >&2
  exit 1
fi

exec go -C "$script_dir" run ./cmd/source-layout "$@" --root "$repo_root"
