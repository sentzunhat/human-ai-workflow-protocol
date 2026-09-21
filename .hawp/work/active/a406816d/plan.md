# Reject symlinked .hawp root in install/update shell scripts

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `a406816d-2458-4fbb-82ab-10fa1c2b594e`
**Type:** bug
**Reported:** 2026-09-20

---

## Input (verbatim)

> Install and update scripts follow a repository-controlled .hawp symlink during destructive refresh

## Intake Summary

The generated install and update scripts use `-d .hawp` for preflight checks,
which follows a symlink. Subsequent operations including `rm -rf .hawp/kit`,
`mkdir`, `cp`, `mktemp`, and `mv` then operate beneath the same unverified
path. A repository author who commits `.hawp` as a symlink can redirect all
of these to a chosen external directory.

## Current Context

The legacy migration step explicitly rejects a `hawp` (without dot) symlink,
but no equivalent check is applied to the primary `.hawp` tree. The
vulnerability is in the distribution source templates and propagates to every
generated provider guide.

## Initial Analysis

**Directly verified:**

- `distribution/sources/install/script-core.md:49-52` checks `[ -d .hawp ]`
  without a prior `[ -L .hawp ]` symlink rejection.
- `distribution/sources/install/script-core.md:198-214` executes
  `rm -rf .hawp/kit` and recreates it; this destructive operation follows a
  `.hawp` symlink.
- `distribution/sources/install/script-core.md:269-314` writes binary and
  scaffold files beneath `.hawp/bin`, `.hawp/config`, and `.hawp/work` via
  operations that traverse any symlink in the `.hawp` path component.
- Parallel sections are present in `script-core.md` for the update variant.
- The generated files replicate these sections verbatim.
- The existing legacy migration check uses `[ -L hawp ]` (no dot); it is
  not applied to `.hawp`.

**Inferred (not yet proven):**

- An attacker must be able to commit `.hawp` as a symlink in the repository;
  most VCSes permit this.
- Exploitation requires an operator to run install/update in the malicious
  repository.

**Likely scope:**

- `distribution/sources/install/script-core.md`
- `distribution/sources/update/script-core.md`
- All generated provider variants after regeneration.

## Root Cause

Preflight existence checks use `-d` (follows symlinks) rather than `-d` plus
`! -L` (reject symlinks). Once the directory-presence check passes, all
subsequent operations inherit the follow behavior.

## Options Considered

1. Add `[ -L .hawp ] && { echo 'error: .hawp must not be a symlink'; exit 1; }`
   before any `.hawp`-relative operation. Minimal, consistent with the existing
   legacy check. **Recommended.**
2. Resolve `.hawp` physically at the start of the script and operate on the
   canonical path throughout. More robust against nested symlinks but requires
   `realpath` availability.

## Recommended Fix

- At script startup (after repository root detection), assert `! -L .hawp` and
  abort with a descriptive error if the check fails.
- Apply the same check to managed ancestors such as `.hawp/kit`, `.hawp/bin`,
  and `.hawp/work` immediately before destructive operations on those paths.
- Apply to all distribution source templates and regenerate.

## Verification Plan

- Add a documented manual test: create `.hawp` as a symlink to a tmp directory,
  run the script, and assert it exits non-zero without modifying the external
  directory.
- Run `hawp distribution sync` and confirm generated output matches.
- Run `go run ./cmd/hawp check` and `go run ./cmd/hawp work validate`.

## Risk + Review Gate

**Risk:** medium — shell script distribution change; symlink check must be
portable across macOS and Linux target environments
**Gate:** explicit review before implementation and regeneration

## Backlog + Plan Link

**Status now:** plan-ready
**Plan file:** work/active/a406816d/plan.md

## Next Step

- [x] Investigation recorded above
- [x] Plan written with options, recommendation, and verification scope
- [x] Backlog moved to `plan-ready`
- [ ] Obtain explicit approval before implementation
