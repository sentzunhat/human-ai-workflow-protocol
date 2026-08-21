# Release-readiness guidance across provider packs

**UUID:** `d7e5a1b9` · **Type:** infrastructure/docs · **Priority:** P1
**Reported:** 2026-08-15 · **Status:** in-progress · **Owner:** unassigned

## Intent

Propagate provider-neutral guidance for dependency consolidation, CI/Docker
verification, release gates, and secure OIDC publication to every supported
HAWP provider without changing HAWP's runtime or project-owned work records.

## Scope

- Add one canonical shared release-readiness behavior.
- Materialize it for GitHub, Claude Code, Cursor, and Continue.
- Update the hand-maintained Codex and GitHub global instruction overlays.
- Regenerate and validate provider/distribution artifacts.

## Acceptance

- All provider materialized files validate as current.
- Distribution-generated guides validate as current.
- Guidance keeps merge, tag, publish, and credential-revocation actions
  owner-approved and distinguishes direct evidence from inference.
- Existing worktree changes remain preserved.

## Verification

- `mise exec node@26.5.0 -- npm --prefix librarian run distribution:sync`
  passed: provider materialization and validation passed; distribution build
  and validation passed.
- Pre-existing dirty files were preserved; no merge, push, or publish ran.
