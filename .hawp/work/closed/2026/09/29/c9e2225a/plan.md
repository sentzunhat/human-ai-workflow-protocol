# Close root ancestry gaps across MCP and filesystem guards

**UUID:** `c9e2225a-1242-47e1-a02b-2c1ed750f0ba`
**Type:** bug
**Reported:** 2026-09-29

## Investigation

User requested the new MCP root findings and related boundaries be reviewed and fixed with low usage. Both explicit MCP roots use os.Stat; the shared RejectSymlinkAncestors guard starts at the supplied root, overlooking ancestors above it. Direct provider configuration and search-index resolution use this same guard. Usage SQLite sidecar checks omit symlink validation. Independent read-only investigation confirmed the root gap and the shared enforcement point.

## Coordination

Clean feature/v0.0.24 at 4d39675c before edits. pwd and git top-level matched <repo-root-abs>; prefix and short status were empty. Existing active items concern deferred MCP features, not these guard files. Can implement now: yes; user authorized fixes and review resolution.

## Plan

Strengthen the shared ancestry guard using a private volume-rooted walker; retain containment, missing-root errors, missing descendants, relative paths and macOS system aliases. Reject unsafe explicit MCP roots before serving/configuring. Apply full ancestry checks to usage SQLite sidecars. This closes equivalent callers; only adding two CLI checks would leave direct callers exposed.

Risk: medium, shared filesystem helper. Changes: librarian/src/internal/infrastructure/filesystem/symlink_guard.go; librarian/src/internal/platform/cli/mcp/commands.go; librarian/src/internal/platform/cli/mcp/configure/command.go; librarian/src/internal/infrastructure/repositories/usage/store.go and focused tests.

Verify original and alternate symlink cases, legitimate paths, full make check, source-layout tests, work validation and final diff. Review candidate independently once. Push and resolve the two threads after green CI. No merge is authorized.

## Outcome

Fixed both explicit MCP root resolvers and the shared filesystem guard, protecting direct provider configuration, indexing, archive, normalization and atomic-write callers from symlinked ancestors above the supplied root. Also fixed the confirmed usage-database sidecar symlink omission.

## Verification

- Before patch: focused regressions failed for both MCP resolvers (root link, parent link, relative path), shared guard/index resolver/atomic write (external sentinel changed), and usage sidecars (database created through accepted sidecar links).
- After patch: the same regressions pass, external sentinels remain unchanged and SQLite is not opened for rejected sidecars.
- Focused MCP CLI, provider configuration, filesystem and usage tests passed. Existing legitimate absolute/relative/space-containing/default-root controls, missing-root/descendant behavior and normal database operations pass.
- Full make check from librarian/src passed (vet, tests, static build); source-layout go test ./... passed.
- Independent investigation checked the root-boundary call sites; independent candidate review found no concrete remaining bypass or compatibility regression in scope.
- Additional static sweep covered explicit root flags, database opens, token multiplication and direct write sites. This is not a new manual line-by-line review of every PR document or binary. Concurrent path substitution after preflight remains outside these static path checks.
- Remote CI and review-thread closure will be checked after this commit is pushed; no merge or human approval is claimed.

## Close Checklist

- [x] Outcome and local verification recorded.
- [x] Work record moved to closed and backlog reconciled.
- [x] Implementation independently reviewed.

Status: local implementation and verification complete; publication readback recorded in the accompanying response.
