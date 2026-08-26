# releases-prerelease-fallback — Verify/fix: `/releases/latest` 404 when all releases are prerelease

**Type:** investigation  
**Status:** done  
**Opened:** 2026-08-25  
**Closed:** 2026-08-25
**Target:** v0.0.11

## Goal

Confirm whether the current `hawp update` and install script handle the
`/releases/latest` 404 gracefully. If not, implement a fallback to
`GET /repos/.../releases` (list all), pick the newest non-draft tag, and
proceed. SHA256 verification must apply regardless of path.

## Outcome

Within inspected scope, the Go release client was already correct: it uses
`GET /repos/.../releases?per_page=1`, not `/releases/latest`, specifically so
prereleases remain visible. Existing unit tests in
`librarian/src/internal/infrastructure/githubrelease/` already covered that
behavior.

The real gap was in the generated install/update shell guides, whose binary
download blocks still resolved tags from `/releases/latest`. Those source
templates now fall back to `/releases?per_page=1` when `/releases/latest` is
unavailable, while preserving the existing SHA256 verification flow.

## Verification

- [x] Confirmed `githubrelease.Client.Latest()` uses `/releases?per_page=1`
- [x] Confirmed prerelease coverage tests already exist in `githubrelease_test.go`
- [x] Updated `distribution/sources/install/script-core.md` and `distribution/sources/update/script-core.md` with releases-list fallback
- [x] `go test ./internal/infrastructure/githubrelease/... ./internal/domain/kitsync/... ./internal/application/kitsync/... ./internal/platform/cli/... ./internal/platform/mcp/...`
- [x] `mise exec node@26.5.0 -- npm run distribution:sync`

## What was done

- Verified the CLI update client was already prerelease-safe
- Added prerelease-safe tag resolution to the install/update shell templates
- Regenerated all provider install/update guides so downstream copy-paste flows inherit the fix
