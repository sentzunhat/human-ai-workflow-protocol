# releases-prerelease-fallback — Verify/fix: `/releases/latest` 404 when all releases are prerelease

**Type:** investigation  
**Status:** analyzing  
**Opened:** 2026-08-25  
**Target:** v0.0.11

## Input

Downstream install of hawp 0.0.9 hit a 404 from `GET /repos/.../releases/latest`.
This endpoint returns 404 when every GitHub release on the repo is marked
as a prerelease, even if tagged releases exist. The install/update script was
able to fall back to listing tags and downloading the SHA256-verified binary.

Hawp 0.0.9 CHANGELOG notes a `--check` fix (`runUpdateVerify`) but does not
explicitly call out a releases/latest fallback to list-releases. Whether the
current update client already handles this needs verification.

Evidence source: downstream install evidence 2026-08-25.

## Goal

Confirm whether the current `hawp update` and install script handle the
`/releases/latest` 404 gracefully. If not, implement a fallback to
`GET /repos/.../releases` (list all), pick the newest non-draft tag, and
proceed. SHA256 verification must apply regardless of path.

## Constraints

- The fix must not change the happy path (when `/releases/latest` returns 200).
- The fallback must still verify the SHA256 checksum of the downloaded binary.
- Do not remove the prerelease flag from GitHub releases as a workaround.

## Plan

### Step 1 — Read current `githubrelease` infrastructure

Read `librarian/src/internal/infrastructure/githubrelease/`. Determine:
- Does it hit `/releases/latest` or `/releases`?
- Does it already fall back on 404?
- What error does it surface to the user?

### Step 2 — Reproduce the 404 scenario

In a test or with a mocked HTTP client, confirm whether a 404 on
`/releases/latest` causes an error or is already handled.

### Step 3 — Implement fallback (if needed)

If the fallback is missing:
- On a 404 from `/releases/latest`, call `GET /repos/.../releases`
- Filter out drafts; prefer the highest semver tag
- Continue with SHA256 verification and binary replacement as normal
- Log a note: `latest release is prerelease; using /releases list fallback`

### Step 4 — Test

- Unit test: mock `/releases/latest` → 404; mock `/releases` → list with one prerelease entry; confirm update succeeds
- Ensure SHA256 check still runs in the fallback path

## Verification

- `hawp update` on a repo where all GitHub releases are prerelease successfully finds and installs the latest binary
- SHA256 mismatch still aborts the update
