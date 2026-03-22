# Raw Output — HAWP (Audit pattern, Layer 1 + Layer 2)

The agent received the filled HAWP shape below (Audit Tasks pattern, now carrying the Layer-1 tiering + verified-correct requirements **and** the Layer-2 surface sweep), followed by the captured output. This re-runs the 2026-06-15 install-docs truth audit with the upgraded authoring pattern, on the same (unchanged) `distribution/` subject.

## The filled shape

**input:** "Check whether our install docs match what the scripts actually do."

**context:** "This is the HAWP source repository. Install/update behavior is defined by composed bash in `distribution/sources/install/script-core.md` plus per-provider fragments, rendered into `distribution/generated/<provider>/install|update/<branch>.md`. README Get Started links to those generated guides. Read-only doc-vs-script audit; make no changes."

**mission:** "Produce a bounded truth audit comparing install documentation against the actual install scripts, to support a decision on whether install docs are trustworthy enough to ship or need correction first."

**constraints:** "Evidence bar: separate observed from inferred for every claim; label confidence (Confirmed/Likely/Unclear). Scope: install paths only. Every mismatch must cite the doc section and the script lines/behavior on both sides. Categories: truth-risk drift, validation drift, maintainability drift, standard mismatch. No redesign. Bounded-absence wording. Run `distribution:validate` as a sync check. **Tier the findings (Layer 1):** cap *primary* findings at 7, but include an uncapped 'Minor / deferred' one-liner list. **State what is correct (Layer 1):** include a 'Verified correct' list, not only non-findings. **Sweep the scoped surface (Layer 2):** enumerate every in-scope install unit — shared docs (install.md, safety.md, repo-boundaries-kit.md), the install script fragments, each provider's install-contract/boundaries/safety/script-install/preamble, the generated install guides, and the README Get Started links — and confirm each was inspected, as a visible checklist; bounded-absence wording for anything not inspected."

**output:** "A decision artifact: (1) scope + method + visible surface-sweep checklist; (2) up to 7 primary findings (category, confidence, Observed doc + Observed script, Inference, Significance); (3) uncapped 'Minor / deferred' one-liners; (4) 'Verified correct' list; (5) verdict (trustworthy / partially / not) with rationale; (6) Fix first / Verify next / Defer."

---

## Captured output

# Install Docs vs Scripts — Truth Audit

**Scope:** Install guides for `cursor`, `github`, `continue` on `main` + `dev`, plus README Get Started mapping. **Method:** read source fragments and generated guides; cross-checked install contracts, boundary tables, safety prose, and embedded bash; ran `npm --prefix librarian run distribution:validate` (exit 0 — "generated outputs are current").

### Surface sweep (coverage checklist)

| In-scope install unit | Inspected | Result |
| --- | --- | --- |
| `shared/install.md` | yes | F5 |
| `shared/safety.md` | yes | F2, F4 |
| `shared/repo-boundaries-kit.md` | yes | verified correct |
| `install/script-core.md` + `script-footer.md` | yes | F2 (destructive kit refresh) |
| `providers/cursor/install-contract.md` | yes | F1 |
| `providers/cursor/{boundaries,safety,script-install,preamble-install}.md` | yes | F1 (cross-ref), verified correct |
| `providers/github/script-install.md` | yes | F3 |
| `providers/github/{install-contract,boundaries,safety,preamble-install}.md` | yes | verified correct |
| `providers/continue/{install-contract,boundaries,safety,script-install,preamble-install}.md` | yes | verified correct |
| `generated/*/install/{main,dev}.md` | yes | match sources (validate passed) |
| README Get Started links | yes | minor (branch-mapping not deep-checked) |

---

## Top findings (max 7)

### F1 — Cursor AGENTS.md install contract contradicts the script
**Category:** truth-risk drift · **Confidence:** Confirmed

- **Observed (doc):** `providers/cursor/install-contract.md` line 17: "Refresh all `hawp-*.mdc` rules from the provider pack and **AGENTS.md** on every install." `cursor/boundaries.md` + `safety.md` instead say AGENTS.md is **seed-if-missing** on install.
- **Observed (script):** `cursor/script-install.md` line 17: `copy_file_no_clobber "$pack/AGENTS.md.seed" AGENTS.md` — no refresh when the file exists.
- **Inference:** A user re-running install expecting AGENTS.md to update will not get that behavior (impact Likely).
- **Significance:** High — the install contract is the agent-facing execution section; read in isolation it overrides the accurate boundary table.

### F2 — Global "no-clobber (`cp -Rn`)" safety claim overgeneralizes the script
**Category:** truth-risk drift · **Confidence:** Confirmed

- **Observed (doc):** `shared/safety.md` line 20: "They use no-clobber copy semantics (`cp -Rn`) to avoid overwriting existing files."
- **Observed (script):** `script-core.md` does `rm -rf .hawp/kit` then full copy (destructive refresh); provider overlays use plain `cp` for rules/instructions (overwrite). Only `.hawp/work/**` scaffold and seed files use no-clobber helpers — and the helper is `copy_file_no_clobber`, not literally `cp -Rn`.
- **Inference:** A reader stopping at the safety section may believe re-running install never overwrites anything HAWP touches.
- **Significance:** Moderate — the top-level safety principle is wrong as stated; mitigated only if the reader also reads the boundary tables.

### F3 — GitHub legacy-filename deletion is undocumented
**Category:** maintainability drift · **Confidence:** Confirmed

- **Observed (script):** `github/script-install.md` lines 16–17 `find ... -name 'human-ai-workflow-protocol-*.instructions.md' -delete` and the matching `*.prompt.md` delete.
- **Observed (doc):** no evidence found in inspected guide prose (safety, boundaries, install steps) describing this deletion.
- **Inference:** Intentional cleanup of an older naming scheme, but invisible to readers; Unclear how many downstream repos still carry legacy-named files.
- **Significance:** Moderate for migrated repos, low for greenfield.

### F4 — Orphan-retirement trigger misdescribed
**Category:** truth-risk drift · **Confidence:** Confirmed (doc wrong); Likely (no fresh-install harm)

- **Observed (doc):** `shared/safety.md` line 34: orphan retirement "only runs when the backlog has at least one **data row**, so fresh installs are not affected."
- **Observed (script):** `retire_orphan_active_items` counts any `|`-row with a non-empty column 2; the seeded BACKLOG Status Key table rows satisfy that gate, so the pass executes on fresh installs (likely no-ops).
- **Inference:** Doc intent (skip until real work items exist) is not what the script implements.
- **Significance:** Low immediate impact; high trust impact for anyone debugging unexpected retirement.

### F5 — "Safe to re-run install" omits seed-vs-refresh nuance
**Category:** truth-risk drift · **Confidence:** Confirmed

- **Observed (doc):** `shared/install.md`: if `.hawp/` exists, run **update** instead; and "Kit and provider-managed paths refresh each run" under "Install Is Safe to Re-Run."
- **Observed (script):** install prints a preflight tip when `.hawp/` exists but continues with the **install** overlay (`copy_file_no_clobber` for AGENTS.md / copilot-instructions.md). Full refresh of those seed files requires the update script.
- **Inference:** Re-running install is only partially idempotent — kit and rules refresh, seed files do not.
- **Significance:** Moderate — users trusting "safe to re-run install" without switching to update may keep stale AGENTS.md / copilot-instructions.md.

---

## Minor / deferred (uncapped one-liners)

- `shared/safety.md` cites the literal command `cp -Rn`, but the script's no-clobber path is the `copy_file_no_clobber` helper — the example command is itself inaccurate (sub-point of F2).
- README Get Started branch mapping (main→`/main/`, dev→`/dev/` generated guides) was not deep-checked against each generated path — verify the four links resolve to the correct provider/branch files.
- `cursor/install-contract.md` "Guide fetch (review-first)" sets `REF="dev"` as the default in the copy block; confirm the stable-install guidance ("set to main for stable") is equally prominent so users don't default to dev unintentionally.
- Provider `preamble-install.md` wording was inspected for contract claims only, not prose tone — no mismatch found, flagged as shallow-checked.
- `repo-boundaries-kit.md` lists `benchmark/` as "Never Installed Downstream"; confirm the install script actually excludes it (consistent with intent, not separately re-verified here).

## Verified correct (checked, found sound)

- `npm --prefix librarian run distribution:validate` passes — generated install guides match composed sources; embedded bash blocks are not drifting.
- Kit boundary row (`core/.hawp/kit/** → .hawp/kit/**`, "Full refresh every install/update") matches `script-core.md` behavior.
- `.hawp/work/**` preservation is stated consistently across `safety.md`, `repo-boundaries-kit.md`, and the script's no-overwrite rule.
- GitHub and Continue overlay refresh contracts match their scripts (plain `cp` overwrite for instructions/rules) — the cross-provider parity the Cursor contract (F1) breaks.
- The guide-fetch "review-first / write to `/tmp`, never pipe to bash" security pattern is consistent across all three provider contracts.

## Verdict

**Partially trustworthy.** The copy/paste bash blocks and most boundary tables are accurate and validate-clean. Two prose sections actively contradict the scripts (Cursor install contract F1; global no-clobber F2), one script behavior is undocumented (F3), and two re-run/seed nuances are misstated (F4, F5). Safe to run the generated command blocks; do not trust every prose summary without checking the boundary tables.

---

## Fix first
- F1: correct `cursor/install-contract.md` line 17 to "seed AGENTS.md if missing on install; re-run users use update for refresh," matching boundaries + script.

## Verify next
- F2: reword `shared/safety.md` no-clobber claim to distinguish destructive kit refresh / overlay overwrite / seed-only files, and drop the inaccurate `cp -Rn` example.
- F5: align the "safe to re-run install" prose with the seed-vs-refresh reality.
- After edits, run `distribution:sync` and confirm generated guides pick up the changes.

## Defer
- F3: document the GitHub legacy-filename cleanup.
- F4: fix the orphan-retirement doc/script description.
- Minor-list verification items (README branch links, benchmark exclusion).
