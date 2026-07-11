# Evidence: TASK-076 — Repository audit remediation

Date: 2026-06-10
Verified by: agent (Cursor session), commands run locally from repo root.

## Direct evidence (commands and observed output)

- `npm --prefix librarian run typecheck` — exit 0, no diagnostics.
- `npm --prefix librarian test` — `# tests 37 / # pass 37 / # fail 0` (was 22 passing + 1 failing before remediation).
- `npm --prefix librarian run distribution:sync` — `distribution build complete: 12/12 file(s) updated`, then `distribution validation passed: generated outputs are current`.
- `npm --prefix librarian run providers:sync` — `provider validation passed: 11 materialized file(s) are current`.
- `npm --prefix librarian run workflow:validate` — `Result: VALIDATION PASS` (3 passed, 0 issues, 1 tolerated legacy warning). Verification clarity now scans 190 checklist items (was 144) after heading-matching fix.

## Key changes verified

- Stale intake overlays regenerated and root copies synced (`.cursor/rules/`, `.continue/rules/`, `.github/instructions/`) — CI drift gate unblocked.
- `backlog-upgrade` dirty-tree test rewritten against a temp git repo (no longer mutates the live repo); `hasDirtyWorkingTree` fails closed when git is unavailable.
- Rule ID collision resolved (structural section check now `A8`; `A4` remains the per-row Outcome heading rule).
- Evidence-link resolution bug fixed: `checkEvidenceIntegrity` now resolves `../evidence/` links against the work root (previously mis-resolved under `closed/YYYY/MM/`), with path-escape containment.
- Shared `librarian/scripts/lib/` extracted (repo-root finder, markdown walker, `toRepoRelative`, `normalizeForCompare`, legacy cutoff constant); duplicates removed from 7 call sites.
- Swallowed errors now surface stderr warnings in all validators; backlog plan paths are contained to `.hawp/work`.
- `set -euo pipefail` added to install/update script cores; legacy `core/.github` fallback removed; Cursor `AGENTS.md` is seed-if-missing on install (manifest updated); auto-dispatch blocks carry a security note. All 12 guides regenerated.
- `npm test` script added; `.github/workflows/librarian-quality.yml` runs typecheck + tests + workflow validation in CI; `.nvmrc` and `engines` pin Node.
- New unit tests for distribution composition, provider materialization render, and all four workflow validations.
- Docs corrected: root README (Continue overlay no longer "planned", stale roadmap, `.hawp/bin/hawp` documented), librarian README (CLI section, test docs, JSON claim removed), generated banner now uses `npm --prefix librarian run providers:sync`, wrong repo URL in CLI help fixed.
- Orphan fragments `distribution/sources/{install,update}/preamble.md` deleted (referenced only from a closed historical plan).
- `db-decision-template.md` moved from `evidence/` root to `notes/`.

## Addendum (2026-06-11, second-pass doc-drift audit)

- `core/providers/.cursor/README.md` and `core/providers/DOC-VERIFICATION.md` updated: Cursor `AGENTS.md` is seeded on install only when missing (was documented as "refreshed on install/update").
- Root README Validation section now lists both CI workflows (`Validate Distribution Generated`, `Librarian Quality`).
- `sync-distribution-generated.yml`: added `distribution/generated/**` and `librarian/scripts/lib/**` trigger paths; Node now pinned via `.nvmrc` (was `lts/*`).
- Re-ran full suite: typecheck PASS, tests 37/37, providers:sync current (11 files), distribution:sync 0/12 rebuilt + validation PASS, workflow:validate PASS, `hawp backlog validate` PASS, `hawp backlog upgrade --dry-run --validate` PASS (0 auto-fixable, 0 blocked), `workflow:normalize` correctly blocked by dirty-tree guard (exit 2, uncommitted changes present).

## Inference / not directly proven

- CI behavior on GitHub-hosted runners (workflow YAML added/changed but not executed remotely from this session).
- Downstream install/update script behavior with `set -euo pipefail` was reviewed against bash AND-OR-list semantics but not executed end-to-end against a live downstream repository.
