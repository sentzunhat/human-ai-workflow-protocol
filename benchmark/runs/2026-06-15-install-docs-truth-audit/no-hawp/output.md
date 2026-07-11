# Raw Output — No HAWP

Prompt given to the agent (verbatim, no framing added):

> Check whether our install docs match what the scripts actually do.
>
> Return your complete answer as your final response.

**Arm:** no-HAWP (clean workspace at `/tmp/hawp-benchmark-clean-20260615-211005`; agent-loaded rules stripped)

Output captured below.

---

I compared the generated install guides under `distribution/generated/*/install/{main,dev}.md`, their composed script sources under `distribution/sources/`, and the top-level README install summary. I also ran `npm --prefix librarian run distribution:validate` (passed — generated outputs match sources).

## Overall verdict

Most of the install documentation **does** match the scripts. The generated command blocks are faithful copies of `script-core.md` + provider `script-install.md` + `script-footer.md`. Provider path mappings in the README and per-guide boundary tables align with what the bash actually copies.

There are, however, several **internal contradictions within the guides** and a few **undocumented script behaviors** that a reader relying on prose alone would get wrong.

## Matches (docs and scripts agree)

**Kit install behavior.** Guides say `core/.hawp/kit/**` → `.hawp/kit/**` with a full refresh every install. The script does `rm -rf .hawp/kit` then copies README, start-here, instructions, templates, patterns, reviews, examples, types, usage, references, and standards — matching the "full refresh" claim.

**Work scaffold seeding.** Docs say `.hawp/work/**` scaffold files are seeded only when missing. The script uses `copy_file_no_clobber` for all scaffold READMEs, BACKLOG, and STATUS — correct.

**GitHub Copilot instructions on install.** Boundaries table and "What Was NOT Changed" say `copilot-instructions.md` is seeded only if missing on install. Script uses `copy_file_no_clobber` — correct. Update script uses plain `cp` — also documented as refresh on update.

**Cursor rules on install.** Safety section and boundaries table say `hawp-*.mdc` rules refresh on every install. Script uses `cp "$pack/rules/"*.mdc` (overwrite) — correct.

**Continue overlay.** Docs say only `hawp-*.md` rules are installed. Script copies `hawp-*.md` only — correct.

**Legacy migration steps.** The safety section's four migration bullets (hawp/ → .hawp/, .hawp/usage/, .hawp/status/, .hawp/work/adrs/) match the numbered migration blocks in `script-core.md` line-for-line in behavior.

**Remote vs local source.** Docs describe `HAWP_LOCAL_CORE` override; script checks the env var and sets `Source mode: local core` vs `remote archive` — correct.

**README Get Started table.** The three provider mapping one-liners (`core/providers/.github/` → `.github/...`, etc.) match the install scripts.

## Mismatches and doc problems

### 1. Cursor `AGENTS.md` — install contract contradicts boundaries and script

This is the clearest truth drift.

- **Provider boundaries table** (`distribution/sources/providers/cursor/boundaries.md`): `AGENTS.md.seed` → `AGENTS.md` is **seed if missing** on install, **refresh** on update.
- **Cursor overlay safety** (`cursor/safety.md`): same — seed on install when missing, refresh on update.
- **Install script** (`cursor/script-install.md`): `copy_file_no_clobber "$pack/AGENTS.md.seed" AGENTS.md` — seed only.
- **But Cursor install contract** (`cursor/install-contract.md` line 17): *"Refresh all `hawp-*.mdc` rules from the provider pack and **`AGENTS.md` on every install**."*

A reader following the install contract will expect `AGENTS.md` to update on re-run install. The script will not touch an existing file. The boundaries table and script agree with each other; the install contract does not.

### 2. Global "no-clobber" claim is overstated

`safety.md` says install/update use no-clobber copy semantics (`cp -Rn`) throughout. That is true for `.hawp/work/**` and seed files, but:

- Kit refresh **deletes and recreates** `.hawp/kit/` (`rm -rf .hawp/kit`).
- Provider rules use plain `cp` (overwrite), not `cp -Rn`.
- GitHub instructions/prompts use plain `cp` (overwrite).

The kit boundaries table correctly says "Full refresh every install/update," so the information exists — but the blanket safety bullet is misleading if read without the provider-specific tables.

### 3. GitHub legacy filename cleanup is undocumented in user-facing prose

Both GitHub install and update scripts run:

```bash
find .github/instructions -maxdepth 1 -type f -name 'human-ai-workflow-protocol-*.instructions.md' -delete
find .github/prompts -maxdepth 1 -type f -name 'human-ai-workflow-protocol-*.prompt.md' -delete
```

No install guide section describes this cleanup. Users with old naming convention files would see them vanish with no documented explanation.

### 4. Orphan active-item retirement trigger is misdescribed

Safety docs say orphan retirement *"only runs when the backlog has at least one data row, so fresh installs are not affected."*

The script's `retire_orphan_active_items` gates on `known_ids` from **any** pipe table row in BACKLOG.md whose second column is non-empty — including the **Status Key** reference table (`| inbox |`, `| analyzing |`, etc.). A freshly seeded BACKLOG from the template contains those rows, so the pass **does run** on first install. In practice it likely has no effect because `active/` only contains `README.md` (excluded), but the documented trigger condition is wrong.

### 5. Re-run install when `.hawp/` already exists — soft mismatch

Shared install concepts say: *"If `.hawp/` already exists, run **update** for the same provider and branch instead of install."*

The install script detects existing `.hawp/`, prints a preflight tip to use the update guide, but **continues running the install overlay path** (seed semantics for `AGENTS.md` / `copilot-instructions.md`). Docs elsewhere say install is "safe to re-run" and that "kit and provider-managed paths refresh each run" — which is only fully true if you use the **update** script for seed-vs-refresh files. Re-running install refreshes kit and rules but not `AGENTS.md` or `copilot-instructions.md`.

### 6. Minor: kit `instructions/` subfolder not named in prose summaries

"What Was Added" sections list templates, patterns, examples, references — the script also copies `instructions/` and `standards/`. Not wrong (`.hawp/kit/**` covers it), but easy to under-count what's installed.

## What looks fine

- Generated guides are in sync with sources (`distribution:validate` passed).
- Install vs update script split correctly implements seed-on-install / refresh-on-update for the files where boundaries tables document it — aside from the install-contract contradiction above.
- Auto-dispatch commands in install contracts use four-backtick fences so they are not confused with the copy/paste install block — intentional and consistent.

## Suggested fixes (if you want docs to match reality)

1. Fix `cursor/install-contract.md` provider-specific rules to say seed-if-missing for `AGENTS.md` on install (match boundaries + script).
2. Narrow the global no-clobber bullet in `safety.md` to project-owned files, and point to the kit/overlay tables for refresh behavior.
3. Document the GitHub legacy filename deletion in `github/safety.md` or boundaries.
4. Reword orphan-retirement trigger to match the script (or fix the script to ignore non-work-item tables).
5. Clarify in shared install docs that re-run install ≠ update for seed files.
