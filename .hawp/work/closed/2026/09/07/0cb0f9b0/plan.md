# v0.0.24 single-executable installation and configuration migration

**UUID:** `0cb0f9b0-1df6-4250-8274-1c3bfffbbf38`
**Type:** improvement
**Reported:** 2026-09-06
**Owner:** Codex
**Status now:** done
**Closed:** 2026-09-07

## Input

> Coordinate single-executable installers, updates, guides, and configuration migration. Validate completion before proceeding to remaining CLI argument validation, UUID correctness, storage boundaries, preservation-first unknown work handling, and broader indexing. Prefer available HAWP MCP tools and maintain existing work context. Do not merge PRs or modify downstream repositories.

## Current Context

Native destination alignment and source-wrapper retirement are implemented.
The earlier blanket closure claim is superseded by this scoped closeout; full prior notes remain in
[progress history](progress-archive.md). Parent release:
[e5fca9c7](../../../../../active/e5fca9c7/plan.md); architecture:
[47c793d6](../../../../../active/47c793d6/plan.md).

## Investigation

Confirmed from current source:

- Canonical install/update binary functions target `.hawp/bin/hawp` (or
  `hawp.exe` on Windows); generated guides are checked against those sources.
- `librarian/src/Makefile` installs the native build by temporary file and rename.
- MCP configs launch native `hawp mcp --repo-root`; JSON/TOML preservation and
  refusal tests exist under `librarian/src/internal/platform/mcp/`.
- The unused `core/.hawp/bin/hawp` source launcher is retired. Checked-in
  release packaging, kit sync, provider mappings, and canonical installers
  do not consume it. Existing installed compatibility files remain untouched.
- Earlier download tests covered installed hawp/hawp-bin files only. The new
  matrix covers fresh/legacy directories, old hawp-mcp preservation, and repeats.

## Plan And Coordination

Keep the native path, retain existing legacy files, and prove the downloader
boundary without changing installed tools or provider configurations. Clean up
current status and guidance while preserving historical evidence. This is a
low-risk test/documentation continuation, authorized by the user's cleanup
request. Can implement now: yes; preserve prior architecture/reshape dirty work.
No new dependencies, runtime changes, or downstream writes.

## Outcome

Current source uses the canonical native destination. Configuration migrations
preserve supported custom settings and refuse unsupported layouts. Fresh,
legacy, failure, and repeat-run fixtures are now explicit. This outcome covers
extracted downloader functions, not a transactional whole installer or every
platform's live execution.

## Verification

Verified 2026-09-07: 40 download scenarios pass (2 modes × 2 layouts × 10
outcomes), with four successful cases repeated, executable-permission checks,
legacy hawp-bin/hawp-mcp preservation, and staging cleanup. Full Go tests and
vet pass, as do distribution validation, root/core kit parity, diff checks, and
connected `hawp_work_validate` (kit/work/links; zero issues/warnings).

The tests execute extracted canonical
shell functions against temporary files with network and Linux asset selection
mocked; payloads are fixture bytes, not executable model/provider binaries.
Historical six-target builds and temporary consumer-copy evidence are retained
in progress history, not claimed as freshly rerun proof.

## Close Checklist

- [x] Canonical native path aligned across sources and local install
- [x] Preserving provider configuration implementation exists
- [x] Fresh/legacy/failure/repeat downloader regression matrix added
- [x] Current plan reconciled; conflicting historical claims preserved separately
- [x] Decide the core launcher policy after inspecting maintained packaging consumers
- [x] Record final acceptance boundaries and close the backlog entry consistently

## Next Step

Migration item complete within the local source/configuration compatibility
scope. Release publication remains with the parent release item; architecture
and intake reshaping continue in their own work items. No global legacy-file
removal or additional provider execution is implied by closure.

## Final launcher policy investigation

2026-09-07 confirmed: `.github/workflows/release.yml` packages only kit and
provider overlays. `librarian/src/internal/application/kitsync/kitsync.go`
consumes those roots, and `core/providers/manifest.yaml` has no binary mapping.
Canonical install/update scripts copy kit/provider files and download native
hawp; no core/bin consumer was found in the inspected maintained paths.

Decision: retire only the unused `core/.hawp/bin/hawp` source launcher. Keep
existing installed hawp-bin/hawp-mcp compatibility files; tested installers
preserve them. This does not require a global inventory because no external
files are being removed. Historical release sources remain in git history.
Can implement now: yes, within the authorized migration and cleanup scope.
Verify current release-bundle assembly, fresh/legacy downloader tests, kit sync,
MCP configuration tests, full suite/vet, generated distribution, and HAWP checks.
Then close the item and repair current references to its archived location.

## Final Verification And Acceptance Boundaries

After source-launcher removal, full `go test ./...` and `go vet ./...` pass;
distribution validation passes. Release bundle assembly from the checked-in
workflow produced 195 kit/provider members with no launcher (temporary output
path; macOS copy metadata disabled to avoid AppleDouble entries). All temporary
bundle/staging files were removed. This is local packaging proof, not CI execution.

The existing downloader suite verifies 40 fresh/legacy cases and repeat success.
No installed binary was rebuilt; no external wrapper/configuration was changed.
Historical cross-target builds and user-reported provider checks remain labeled
in the preserved history. Live downloads and every-platform runtime are outside
this local compatibility closeout; release approval/publication remains separate.
