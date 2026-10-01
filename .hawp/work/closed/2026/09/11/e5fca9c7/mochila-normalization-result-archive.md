# Mochila normalization follow-up

2026-09-05: user supplied a successful normalization report from
`mochila-archive-viewer`, branch `chore/hawp-work-normalization-2026-09-05`:
five checks passed, zero issues/warnings, no commits, existing dirty work preserved.
The report lists stale completed rows, missing closed-record sections, and
flat-to-folder migration. Thirty-five unidentified legacy records were skipped.

Independent read-only check after the report, using the current HAWP source:
`work validate --hawp-root <mochila-repo>/.hawp --no-update-check` passed all
five checks, with 10 active rows, 10 recent closed rows, 40 complete closed
plans, zero issues, and zero warnings. No downstream files were changed here.

## Command Boundaries

- `work normalize --dry-run --validate` detects record/backlog drift.
- `work normalize --apply --validate` performs supported mechanical cleanup.
- `work normalize --dry-run --migrate-folders --validate` previews folder moves.
- `work normalize --apply --migrate-folders --validate` applies those moves.
- Dirty repositories require the explicit `--force-dirty` override for apply.
- Numeric legacy IDs remain valid. Folder migration does not silently reassign
  their identity to UUIDs; newly created items use UUID folders.
- Skipped records without inferable identity and missing real verification
  evidence require review. Added sections are scaffolding, not invented proof.

## Safety Follow-Up

The user explicitly requires preservation of unreferenced copies. Duplicate
normalization now adds reciprocal links when there is one archived counterpart,
preserving original text, all files, and supporting artifacts. Ambiguous matches
are reported for review. Referenced active/parked items remain untouched.
The duplicate portion of JSON dry-run output is available as `duplicateLinks`.
Cross-references do not change lifecycle status or suppress duplicate warnings;
that requires an explicit decision about which record remains active.

Linking verification: focused preservation/ambiguity/idempotency tests, full Go
suite, vet, six standard platform builds, local 0.0.24 rebuild, and HAWP check
all passed. These linking changes were tested in temporary fixtures only; no
linking or removal was performed against Mochila or Tekit in this continuation.

Earlier Tekit temp-copy success predates this preservation policy; do not
interpret that earlier result as proof that all current duplicate cases auto-fix.
