# Task: Apply repository audit remediation (quality, safety, drift, docs)

**Backlog ID:** TASK-076

## Context

A full-repository audit (2026-06-10) produced a prioritized improvement backlog (IMP-01..IMP-21). The user requested all needed changes and fixes be applied in one pass.

## Scope

- P0: regenerate stale intake overlays + root copies (CI drift gate).
- P1: fix dangerous dirty-tree test, fail-closed git guard, npm test script + CI quality workflow, install-script hardening.
- P2: rule ID collision (A4/A8), verification heading matching, evidence-link resolution bug, error surfacing, path containment, shared `librarian/scripts/lib/`, dead-code removal, import/async consistency, new unit tests, backlog/status hygiene, doc fixes.
- P3: `.nvmrc` + engines, orphan fragment cleanup, CLI documentation.

## Outcome

All planned remediation items applied. Librarian suite: 37/37 tests passing, typecheck clean, distribution and provider materialization validators green, workflow validation PASS. Generated guides (12) and materialized overlays (11) rebuilt from updated sources; repo root overlays synced with provider packs.

Deliberately not done (kept as future work): UUID work-item migration (tracked in Future Improvements), checksum pinning for remote archive downloads (would require release tagging strategy), removal of the auto-dispatch `curl | bash` pattern (kept as product direction, now annotated with a security note).

## Verification

- [x] `npm --prefix librarian run typecheck` exits 0. Evidence: ../evidence/2026/06/10/TASK-076-audit-fixes-verification.md
- [x] `npm --prefix librarian test` 37/37 pass. Evidence: ../evidence/2026/06/10/TASK-076-audit-fixes-verification.md
- [x] `distribution:sync` and `providers:sync` validators pass with regenerated outputs. Evidence: ../evidence/2026/06/10/TASK-076-audit-fixes-verification.md
- [x] `workflow:validate` reports VALIDATION PASS. Evidence: ../evidence/2026/06/10/TASK-076-audit-fixes-verification.md
- [ ] GitHub-hosted CI run for `librarian-quality.yml` — unproven until pushed (local-only session).

## Close Checklist

- [x] Outcome recorded with explicit not-done items.
- [x] Evidence file stored under `evidence/2026/06/10/`.
- [x] BACKLOG.md Recently Closed updated and capped at 10 rows.
- [x] STATUS.md refreshed.
