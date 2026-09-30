# Reject path traversal in install/update backlog reconciliation

## Outcome

The PR-review finding was implemented and reconciled with the current branch.

## Verification

The focused regression coverage and repository-wide Go, HAWP, distribution,
formatting, and diff-hygiene checks passed before close.

## Close Checklist

- [x] Outcome recorded.
- [x] Verification evidence recorded or referenced.
- [x] Backlog row removed from active coordination.
- [x] Plan archived under the close date.

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `a1014eef-d325-49a9-948d-79e6590b9d23`
**Type:** bug
**Reported:** 2026-09-20

---

## Input (verbatim)

> Backlog reconciliation in install/update shell scripts accepts traversal components in closed-plan destinations

## Intake Summary

The generated install and update shell scripts parse BACKLOG.md to move
completed active plans into `.hawp/work/closed/`. The destination path is
built from the plan link by accepting any `closed/`-prefixed token, then
passing it via `mkdir -p` and `mv` without normalization. A malicious repository
with a BACKLOG.md link containing `../../../.config/...` can redirect a plan
move to a target outside the repository.

## Current Context

This is a shell-script boundary, not a Go boundary. The vulnerability lives in
`distribution/sources/install/script-core.md` and
`distribution/sources/update/script-core.md`; the generated files
`distribution/generated/codex/install/main.sh` and
`distribution/generated/codex/update/main.sh` are the runnable surfaces.
Source documents must be fixed and `hawp distribution sync` run to regenerate.

## Initial Analysis

**Directly verified:**

- `distribution/generated/codex/install/main.sh:69-117` matches backlog links
  beginning with `closed/`, `work/closed/`, or `.hawp/work/closed/` and
  constructs `dest=".hawp/work/$closed_path"` from the matched suffix.
- `mkdir -p "$dest_dir"` and `mv "$plan_file" "$dest_dir/"` are called without
  a further normalization or containment check.
- Variables are double-quoted (no shell metacharacter injection), and `mv`
  uses a non-existent destination, but neither prevents a new external file
  from being created through `../` traversal.
- The same pattern is present in the update script counterpart.

**Inferred (not yet proven):**

- The attacker controls BACKLOG.md content; the operator runs the script on
  a malicious or compromised checkout.
- The source Markdown templates share the same logic; regeneration without
  source fix would reproduce the issue.

**Likely scope:**

- `distribution/sources/install/script-core.md`
- `distribution/sources/update/script-core.md`
- `distribution/generated/codex/install/main.sh` (regenerated)
- `distribution/generated/codex/update/main.sh` (regenerated)
- Any other generated provider variants.

## Root Cause

The reconciliation loop trusts a BACKLOG.md-controlled path component as the
closed destination suffix without canonicalization. The prefix check verifies
only the leading token, not the full resolved path.

## Options Considered

1. Normalize the destination and verify it remains under `.hawp/work/closed/`
   before creating directories or moving files. Low complexity and consistent
   with existing shell safety patterns. **Recommended.**
2. Remove backlog-driven reconciliation from shell scripts and rely solely on
   the Go `hawp work normalize` command. Simpler and more robust, but a larger
   behavioral change.

## Recommended Fix

- After extracting `$closed_path`, canonicalize with `realpath --no-symlinks`
  or a POSIX-compatible path normalization, then check that the result begins
  with `<abs repo>/.hawp/work/closed/`.
- Reject and log entries that fail containment.
- Apply the same fix to all generated provider variants via source template
  changes plus `hawp distribution sync`.

## Verification Plan

- Add shell test or documented manual test with a BACKLOG.md closed link
  containing `../` components; assert no file is created outside
  `.hawp/work/closed/`.
- Run `hawp distribution sync` and confirm generated output matches fixed
  sources.
- Run `go run ./cmd/hawp check` and `go run ./cmd/hawp work validate` from
  `librarian/src`.

## Risk + Review Gate

**Risk:** medium — shell script mutation in distribution sources
**Gate:** explicit review before implementation; regeneration must follow any
source template change

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/active/a1014eef/plan.md

## Next Step

- [x] Investigation recorded above
- [x] Plan written with options, recommendation, and verification scope
- [x] Backlog moved to `plan-ready`
- [x] Obtain explicit approval before implementation

## Outcome

Done. Install/update reconciliation reject closed-plan link paths containing
`.` or `..` components before deriving destinations. Both source templates
contain the guard, generated outputs are current, and reconciliation/distribution
tests pass.
