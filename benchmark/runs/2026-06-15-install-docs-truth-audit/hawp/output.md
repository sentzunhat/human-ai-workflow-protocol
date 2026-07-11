# Raw Output — HAWP

The agent received the filled HAWP shape below as its prompt (authored using the Audit Tasks pattern in `.hawp/kit/references/authoring-patterns.md`), followed by the captured output.

## The filled shape

**input:** "Check whether our install docs match what the scripts actually do."

**context:** "This is the HAWP source repository. Install/update behavior is defined by composed bash in `distribution/sources/install/script-core.md` plus per-provider fragments (`script-install.md`, `script-update.md`, `script-footer.md`), rendered into `distribution/generated/<provider>/install|update/<branch>.md`. Top-level README Get Started links to those generated guides. This is a read-only doc-vs-script audit; make no changes."

**mission:** "Produce a bounded truth audit comparing install documentation (generated guides + README install summary) against the actual install scripts, to support a decision on whether install docs are trustworthy enough to ship or need correction first."

**constraints:** "Evidence bar: separate directly observed evidence from inference for every substantive claim; label confidence Confirmed, Likely, or Unclear. Audit scope: install paths only (not update guides unless needed to explain install seed-vs-refresh semantics). Standards for 'mismatch' must cite the specific doc section and the specific script lines/behavior being compared — do not call something a violation without naming both sides. Cap findings at 7. Cap non-findings at 2. Use only these categories: truth-risk drift, validation drift, maintainability drift, standard mismatch. Do not propose architecture redesign. When claiming absence in docs, use bounded wording ('no evidence found in inspected guide prose'). Run `npm --prefix librarian run distribution:validate` as a generation sync check."

**output:** "A decision artifact: (1) audit scope and method; (2) per-finding entries with category, confidence, Observed (doc claim + script behavior), Inference/uncertainty, Significance; (3) up to 2 non-findings; (4) verdict — docs trustworthy / partially trustworthy / not trustworthy for install, with one-sentence rationale; (5) closing sequence: Fix first, Verify next, Defer."

---

## Captured output

# Install Docs vs Scripts — Truth Audit

**Scope:** Install guides for all three providers (`cursor`, `github`, `continue`) on `main` and `dev`, plus README Get Started install mapping. Script sources: `distribution/sources/install/script-core.md` and `distribution/sources/providers/*/script-install.md`. **Method:** Read source fragments and generated guides; cross-check boundary tables, safety prose, install contracts, and embedded bash blocks; ran `npm --prefix librarian run distribution:validate` (exit 0 — "generated outputs are current").

## Verdict

**Partially trustworthy.** Generated install command blocks and most boundary tables match the scripts. Two doc sections actively contradict the scripts and each other (Cursor install contract; global no-clobber wording). One script behavior is undocumented (GitHub legacy filename cleanup). Safe to use the copy/paste bash blocks; unsafe to trust every prose summary without checking boundaries.

---

### Finding 1 — Cursor AGENTS.md install contract contradicts script

**Category:** Truth-risk drift  
**Confidence:** Confirmed

**Observed (doc):** `distribution/sources/providers/cursor/install-contract.md` states: "Refresh all `hawp-*.mdc` rules from the provider pack and **AGENTS.md on every install**."

**Observed (doc, elsewhere):** `cursor/boundaries.md` and `cursor/safety.md` say `AGENTS.md` is **seed if missing** on install, **refresh** on update.

**Observed (script):** `cursor/script-install.md` uses `copy_file_no_clobber "$pack/AGENTS.md.seed" AGENTS.md` — no refresh when file exists.

**Inference / uncertainty:** None for the contradiction itself. Impact on users: Likely — someone re-running install expecting AGENTS.md to update will not get that behavior.

**Significance:** High. The install contract is the agent-facing execution section; it overrides the more accurate boundaries table if read in isolation.

---

### Finding 2 — Global no-clobber safety claim overgeneralizes script behavior

**Category:** Truth-risk drift  
**Confidence:** Confirmed

**Observed (doc):** `distribution/sources/shared/safety.md`: "They use no-clobber copy semantics (`cp -Rn`) to avoid overwriting existing files."

**Observed (script):** `script-core.md` runs `rm -rf .hawp/kit` then full copy (destructive refresh). Provider overlays use plain `cp` for rules/instructions (overwrite). Only `.hawp/work/**` scaffold and seed files use no-clobber helpers.

**Observed (doc, partial correction):** `repo-boundaries-kit.md` kit row says "Full refresh every install/update" — accurate for kit.

**Inference / uncertainty:** Readers who stop at the safety section may believe re-run install never overwrites anything HAWP touches.

**Significance:** Moderate. Mitigated if readers also read boundary tables, but the top-level safety principle is wrong as stated.

---

### Finding 3 — GitHub legacy filename deletion undocumented

**Category:** Maintainability drift  
**Confidence:** Confirmed

**Observed (script):** `github/script-install.md` deletes `human-ai-workflow-protocol-*.instructions.md` and matching prompt files before copy.

**Observed (doc):** No evidence found in inspected guide prose (safety, boundaries, install steps, "What Was Added/NOT Changed") describing this deletion.

**Inference / uncertainty:** Unclear how many downstream repos still have legacy-named files; deletion is intentional cleanup but invisible to readers.

**Significance:** Moderate for repos migrated from an older naming scheme; low for greenfield installs.

---

### Finding 4 — Orphan retirement trigger misdescribed

**Category:** Truth-risk drift  
**Confidence:** Confirmed (doc wrong); Likely (no user harm on fresh install)

**Observed (doc):** `safety.md`: orphan retirement "only runs when the backlog has at least one **data row**, so fresh installs are not affected."

**Observed (script):** `retire_orphan_active_items` builds `known_ids` from **all** `|`-rows in BACKLOG.md where column 2 is non-empty, skipping only header/separator/ID rows. Seeded `core/.hawp/work/BACKLOG.md` includes a Status Key table (`| inbox |`, etc.), which satisfies the gate.

**Observed (behavior):** Fresh install seeds BACKLOG template and `active/README.md` only; orphan pass likely no-ops but still executes.

**Inference / uncertainty:** Doc intent (skip until real work items exist) is not what the script implements.

**Significance:** Low immediate impact; high trust impact if someone debugs unexpected retirement behavior.

---

### Finding 5 — Re-run install vs update semantics for seed files

**Category:** Truth-risk drift  
**Confidence:** Confirmed

**Observed (doc):** `shared/install.md`: if `.hawp/` exists, run **update** instead of install. Also: "Kit and provider-managed paths refresh each run" under "Install Is Safe to Re-Run."

**Observed (script):** Install script prints a preflight tip when `.hawp/` exists but continues with **install** overlay (`copy_file_no_clobber` for AGENTS.md / copilot-instructions.md). Full refresh of those files requires the update script.

**Inference / uncertainty:** Re-running install is partially idempotent — kit and rules refresh, seed files do not.

**Significance:** Moderate. Users following "safe to re-run install" without switching to update may have stale AGENTS.md or copilot-instructions.md.

---

## Non-findings

### Non-Finding 1 — Provider overlay refresh behavior (GitHub, Continue)

**Type:** Verified alignment  
**Confidence:** Confirmed

**Observed:** GitHub install contract says refresh instructions/prompts every install; script uses plain `cp` (overwrite) — matches. Continue install contract says refresh `hawp-*.md` every install; script copies with overwrite — matches. Cursor `hawp-*.mdc` refresh on install — matches.

**Why not a finding:** Doc and script agree for these paths.

---

### Non-Finding 2 — Generated outputs match sources

**Type:** Validation check passed  
**Confidence:** Confirmed

**Observed:** `npm --prefix librarian run distribution:validate` → "generated outputs are current." Embedded bash in generated guides matches source composition.

**Why not a finding:** The build pipeline enforces doc-script parity for generated blocks; drift would fail CI, not silently ship.

---

## Fix first

Correct `cursor/install-contract.md` line 17 to match boundaries + script: seed `AGENTS.md` if missing on install; point re-run users to update for refresh.

## Verify next

After editing sources, run `npm --prefix librarian run distribution:sync` and confirm generated cursor install guides pick up the contract fix. Spot-check one install re-run vs update on a scratch repo with a pre-existing `AGENTS.md`.

## Defer

Rewording global no-clobber in `safety.md`, documenting GitHub legacy deletes, and fixing orphan-retirement doc/script alignment — real issues but lower urgency than the Cursor contract contradiction.
