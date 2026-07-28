# verify hawp update end-to-end with real v0.0.1 -> v0.0.2 releases and changelogs

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `4c152ee3-09af-40c0-b372-cff849d063cd`
**Type:** improvement
**Reported:** 2026-07-21
**Risk Level:** low

---

### Input (what was reported)

> Test the auto update, planning to test using low versions like 0.0.1
> and 0.0.2, with changelogs and references, verbose and detailed
> changes for the changelog/release notes.

---

### Context

Closes the known gap noted at the close of `cdcf9f78-513a-47f4-9b44-9c24931ce231`
(self-update): the update mechanism was only verified against local
`httptest` fixtures, never a real GitHub Release round trip. `v0.0.1`
and `v0.0.2` are explicitly test/pre-release versions (marked
`prerelease` on GitHub) — not the first real product release.

---

### Analysis

**Root cause (or most likely cause):**
_Real release publishing (GitHub Actions building `make dist`, GitHub
computing asset digests, `hawp update` fetching from the live API) has
failure modes a local test server cannot reproduce: real network
latency, real digest computation timing, real cross-compiled binaries._

**Directly verified:**
_No releases exist yet for this repo (confirmed via `gh api .../releases`
returning `[]`); the release workflow
(`.github/workflows/release-librarian-go.yml`) has never fired._

**Inferred (not yet proven):**
_Pushing a `librarian-go-v0.0.1` tag triggers the workflow, produces a
GitHub Release with six platform binaries, and `hawp update` against
that real release succeeds end-to-end._

**Scope — what else is affected:**
_`librarian/go/CHANGELOG.md` (new), `.github/workflows/
release-librarian-go.yml` (changelog-driven release body), no source
code changes to the update mechanism itself unless real-world testing
finds a bug._

---

### Recommended Fix

- Add `librarian/go/CHANGELOG.md` (Keep a Changelog style) with detailed
  `0.0.1` and `0.0.2` entries.
- Update the release workflow to extract the matching changelog section
  into the release body (plus GitHub's auto-generated commit notes),
  and mark `0.0.x` tags as `prerelease`.
- Build and tag `v0.0.1`, push, wait for the release to publish.
- Give `0.0.2` real content (see `b4c8af81` — agent usage/commands
  surface) so the update is a meaningful upgrade, not a no-op.
- Build and tag `v0.0.2`, push, wait for the release.
- From a binary built as `v0.0.1`, run real `hawp update --check` (expect
  update available) then `hawp update` (expect success, binary now
  reports `v0.0.2`).

**What to verify after:**

- [x] `v0.0.1` release published with the changelog-derived body
      (Evidence: `gh release view librarian-go-v0.0.1` body matches the
      `## [0.0.1]` CHANGELOG.md section exactly)
- [x] `v0.0.2` release published with the changelog-derived body
      (Evidence: same, for the `## [0.0.2]` section, including the
      retroactive bug-fix note)
- [x] Real `hawp update --check` from `v0.0.1` reports `v0.0.2` available
      (Evidence: real run — `current: v0.0.1 / latest: v0.0.2 / Update
      available`)
- [x] Real `hawp update` from `v0.0.1` succeeds and the binary reports `v0.0.2`
      (Evidence: real run — downloaded checksum
      `b6358ff9af715795a327258df8b9311f3c7b111a09805c40fea24e3296afdeed`
      matched the published digest exactly; binary reported `v0.0.2`
      afterward and a repeat `--check` correctly said "Already up to
      date")
- [x] Both releases marked `prerelease` (test versions, not the real first release)
      (Evidence: `gh release list` shows both as `Pre-release`)

---

## Outcome (filled at close)

Closed 2026-07-21. Real end-to-end proof landed, and it found two
genuine bugs the local `httptest`-based unit tests could not have caught
(local fixtures don't reproduce GitHub's actual `/releases/latest`
semantics or its real tag format):

1. **`/releases/latest` excludes prereleases.** Since 0.0.x test
   releases are deliberately marked prerelease, the first cut of
   `hawp update` 404'd against the real API for both releases. Fixed by
   switching `Client.Latest` to the `/releases` list endpoint.
2. **Tag-prefix pollution.** Real tags are `librarian-go-vX.Y.Z`, not
   plain `vX.Y.Z`; `ParseVersion` only stripped `v`, silently zeroing the
   major-version slot. It happened to still compare correctly for 0.0.x
   by coincidence. Fixed with `CleanVersion`.

Both fixes are committed (`5cf1938`) with regression tests reproducing
the exact real-world failures. Because an already-published binary with
bug #1 can never discover *any* later release while every release stays
prerelease-only (it never gets past the 404), the only way to complete
the test was to delete and recreate both `v0.0.1`/`v0.0.2` releases at
the fixed commit — documented transparently in `CHANGELOG.md` under the
`0.0.2` entry, since fixing it after that binary was already published
would be a permanent dead end.

With the recreated releases, the full loop was proven for real: fresh
`v0.0.1` binary detects `v0.0.2`, downloads and checksum-verifies the
real release asset, atomically replaces itself, reports `v0.0.2`
afterward, and — critically — the *new* binary correctly reports
"Already up to date" on a repeat check, proving the fix is not
self-defeating.

## Verification (filled at close)

- Evidence: `gh release list` — both `librarian-go-v0.0.1` and
  `librarian-go-v0.0.2` published, `Pre-release`, 6 platform assets each.
- Evidence: `gh release view librarian-go-v0.0.2 --json body` matches
  the CHANGELOG.md `## [0.0.2]` section verbatim.
- Evidence: real `hawp update --check` and `hawp update` runs against
  the live repo (transcribed above); downloaded-asset SHA-256 matched
  GitHub's published digest exactly; `hawp commands --json` (17
  commands) confirmed working post-update.
- Evidence: `go test ./...` — new regression tests
  `TestLatestIncludesPrereleases`, `TestCleanVersion`,
  `TestParseVersionStripsTagPrefix`, plus updated `IsNewer` cases for
  prefixed tags, all pass.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled
- [x] Plan file saved under `closed/2026/07/21/4c152ee3-09af-40c0-b372-cff849d063cd.md`
- [x] BACKLOG.md updated
