# Work Status

Local checkpoint: 2026-09-10. The [backlog](BACKLOG.md) owns lifecycle status.

## Active

- [a3df8a9c](active/a3df8a9c/plan.md): request-to-intake reshaping. Existing-tool
  worker guidance and an internal draft contract are implemented; adapter
  fidelity and public wiring remain open.

## Recently Closed (2026-09-11)

- [e5fca9c7](closed/2026/09/11/e5fca9c7/plan.md): v0.0.24 work-folder normalization and README positioning — closed 2026-09-11.
- [47c793d6](closed/2026/09/11/47c793d6/plan.md): CLI decomposition and architecture audit continuation — closed 2026-09-11.
- [5b6d4e21](closed/2026/09/11/5b6d4e21/plan.md): provider parity for shared HAWP agent guidance — closed 2026-09-11.
- [d1fa0b72](closed/2026/09/11/d1fa0b72/plan.md): install/update contract hardening — closed 2026-09-11.

## Recently Completed

[0cb0f9b0](closed/2026/09/07/0cb0f9b0/plan.md): native path/configuration
migration and unused source-launcher retirement. Existing installed legacy
files are preserved. This is a local compatibility closeout, not a release.

## Current Preview Evidence

The [source-layout checkpoint](status/2026/09/10/source-layout-preview.md)
records the earlier passing preview and candidate checks. The subsequent scripts
cleanup adds deeper context folders and retains those artifacts as a prior
snapshot. A fresh preview follows the requested commit; no migration is applied.

## Earlier Evidence (2026-09-07)

- Full Go tests and vet pass, including model-pull and mutation-boundary cases.
  Conflicting kit modes and explicit empty path overrides are rejected;
  temporary-fixture snapshots confirm rejected commands preserve files.
- Connected HAWP MCP validation passes kit/work/links with zero issues/warnings.
- Distribution generated-output validation and root/core kit parity pass.
- `librarian/src/Makefile` installs the native build at `.hawp/bin/hawp`
  via temporary file plus rename. The older launcher/sibling description in
  this dashboard was stale.
- Changes remain uncommitted, including work preserved from the prior session.
  No binary refresh, model inference, cross-platform execution, or release
  validation was performed in this continuation.

## Next Queue

Topology now precedes new features, per the user's 2026-09-07 direction.
[Package map](closed/2026/09/11/47c793d6/package-boundaries-archive.md).

1. After the checkpoint commit, regenerate the
   [prior layout snapshot](../../scripts/source-layout/review/preview.md) using the
   [current instructions](../../scripts/source-layout/README.md), then review it.
   No broad migration is applied.
2. [742aa60b](closed/2026/09/10/742aa60b/plan.md): work rules/use cases/filesystem boundaries — closed 2026-09-10.

[CLI grouping](closed/2026/09/08/e7a5e294/plan.md) and
[model/usage adapters](closed/2026/09/09/b61770a6/plan.md) are closed.
Existing reshape contract and metadata
work remain tracked; resume feature changes after the relevant structural slice.
One small slice per turn, narrow reads, concise updates, focused tests, full checks
at completion. No automatic model switch or usage-reset redemption.
