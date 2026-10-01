# Prepare v0.0.24 release without a tracked binary

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `d605db54-0c6e-41fe-b787-f3332e028d8b`
**Type:** release
**Reported:** 2026-09-30

---

## Input (verbatim)

> Untrack the repository-local hawp binary, use release assets for platform installs, and finish the v0.0.24 changelog and release instructions.

## Intake Summary

Prepare the existing `0.0.24` release lane without committing a machine-specific
executable. Keep `.hawp/bin/hawp` as the local install path and publish platform
assets only through the release workflow after the merge to `main`.

## Current Context

The repository tracks `.hawp/bin/hawp`, while `librarian/src/bin/` is ignored.
`make install` builds and installs the host binary. Generated install/update
scripts download `hawp-<os>-<arch>` plus `checksums.txt` from a GitHub release,
then replace `.hawp/bin/hawp` after verification. The release workflow builds
all six platform assets and reads release notes from `librarian/CHANGELOG.md`.
The version constant is already `0.0.24`; the latest published release is
`0.0.23` as of this investigation.

## Initial Analysis

**Directly verified:**

- `git ls-files .hawp/bin` lists only the native `hawp` binary.
- `distribution/sources/install/script-core.md` and `update/script-core.md`
  download a platform asset and verify its SHA256 before replacement.
- `RELEASE.md` still asks for a checked-in wrapper; the `0.0.24` changelog was
  dated September 13 although release publication is still pending.

**Inferred (not yet proven):**

- Existing local MCP processes may continue running an earlier executable
  image until their client restarts them; removing the path from Git does not
  stop those processes.

**Likely scope:**

- Git tracking/ignore rules, release playbook, changelog, canonical kit wording,
  and release selection in the install/update source scripts.

## Risk + Review Gate

**Risk:** medium. Install selection affects new consumers; preserve checksum
validation and the existing executable on any failed download.
**Gate:** User explicitly requested implementation and release preparation.

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/active/d605db54/plan.md

## Plan

1. Ignore `.hawp/bin/` and remove the tracked binary from Git while preserving
   the local executable used by this checkout.
2. Make install/update select the newest published release, including a
   prerelease, with a stable-release endpoint fallback. Preserve checksum and
   atomic replacement behavior.
3. Generate distribution files, align version and release documentation, and
   check that release notes extract correctly for `0.0.24`.
4. Verify the focused installer behavior and required release checks. Commit
   and push the coherent release-preparation change to `development` for PR #42.

## Outcome

- `.hawp/bin/hawp` was removed from Git tracking and `.hawp/bin/` was added to
  `.gitignore`; the local executable remains available for this checkout.
- Both install and update sources now resolve the newest published release from
  the releases list first, including prereleases, then use `/releases/latest`
  as a stable-release fallback. All 40 distribution guides were regenerated.
- `librarian/CHANGELOG.md` has a dated `0.0.24` section and release highlights.
  `RELEASE.md` and the canonical kit explain the binary's source and install
  path. `librarian/src/internal/domain/update/version.go` already contained
  `0.0.24`, so no version-code edit was needed.

## Verification

- Focused binary download test passed for both install/update sources, fresh
  and legacy layouts, release-list fallback, checksum failures, and preserved
  installed files.
- `make check VERSION=0.0.24` passed (`go vet`, full Go test suite, build).
- `make dist VERSION=0.0.24` produced all six standard platform binaries;
  the host `librarian/src/bin/dist/hawp-darwin-arm64 version` reports `0.0.24`.
- `go run ./cmd/hawp distribution validate`, `providers validate`,
  `kit validate`, `work validate`, and `links check` passed. Work validation
  retained one existing repository-wide warning for 25 ambiguous historical
  verification records and reported zero issues.
- Root and core kit copies match; `git diff --check` and the release workflow's
  `0.0.24` changelog extraction passed.
- The GitHub Actions checks on the pushed PR head remain to be read back.
  Publication, release downloads of `0.0.24`, and a fresh downstream install
  remain unproven until merge to `main` and release completion.

## Close Checklist

- [x] Outcome and direct verification recorded.
- [x] Local binary preserved and ignored after removal from Git.
- [x] Release tag deferred until merge to `main`.
- [x] Status report linked below.

Status: [v0.0.24 release preparation](../../../../../status/2026/09/30/d605db54/status.md).
