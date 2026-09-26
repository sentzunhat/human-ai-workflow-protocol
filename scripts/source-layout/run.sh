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

# Keep the enforced root after caller flags so it wins over a caller-supplied
# --root, but before the caller's -- terminator so Go's flag parser still sees
# it as a flag. Arguments after -- remain positional and are rejected by the
# source-layout command as intended.
forwarded=()
root_injected=false
for arg in "$@"; do
  if [[ "$arg" == "--" && "$root_injected" == false ]]; then
    forwarded+=(--root "$repo_root")
    root_injected=true
  fi
  forwarded+=("$arg")
done
if [[ "$root_injected" == false ]]; then
  forwarded+=(--root "$repo_root")
fi

exec go -C "$script_dir" run ./cmd/source-layout "${forwarded[@]}"
