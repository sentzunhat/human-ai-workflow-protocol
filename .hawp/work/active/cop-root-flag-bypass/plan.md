# Security: --root flag bypass in source-layout script

**UUID:** `cop-root-flag-bypass`
**Type:** fix
**Severity:** High
**Reported:** 2026-09-13
**Status:** `inbox`
**Source:** Copilot PR #40 review — `scripts/source-layout/run.sh`

---

## Finding

The safety wrapper enforces `--root "$repo_root"` but places it before `"$@"`:

```sh
exec go -C "$script_dir" run ./cmd/source-layout --root "$repo_root" "$@"
```

Go's `flag` package uses last-wins for repeated flags. A caller can pass their
own `--root /some/other/path` at the end of `"$@"` and override the enforced
root entirely.

## Fix Plan

1. Move the enforced `--root` flag to the end, after `"$@"`:
   ```sh
   exec go -C "$script_dir" run ./cmd/source-layout "$@" --root "$repo_root"
   ```
2. Confirm that `flag.Parse` in `cmd/source-layout/main.go` uses standard
   `flag` (last-wins) so the appended value wins.
3. If the tool uses a library that does first-wins, document the decision.

## Files

- `scripts/source-layout/run.sh` — one-line fix
