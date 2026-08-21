# AlwaysApply imperative reads vs on-demand kit loading

**From:** a consumer infrastructure repository (follow-on review)
**Type:** improvement (upstream review; not a consumer patch)
**Reported:** 2026-08-21
**Status:** inbox
**Risk:** medium if implemented (touches generated provider behaviors / consumer alwaysApply text)

## Input

> HAWP upstream: review alwaysApply imperative reads (start-here / status-report) vs on-demand kit loading.

## Context

Consumer findings: `.hawp/work/active/hawp-consumer-token-usage-findings.md`. Evidence B.

Generated `hawp-core` (from `core/providers/shared/behaviors`) tells agents to read start-here, status-report, and workflow-loop. Cursor does not auto-inject those kit files; obedient agents still open them (~10k to 25k once per task; whole kit ~105k if cascaded).

## Mission

Decide whether always-on overlay text should stop ordering a kit read on every task, and load those files only when shaping, handoff, or loop work is in play.

## Constraints

Do not change kit or rules until a maintainer scopes an implementation item. This row is review intake only. Do not rewrite consumer kit as the "fix."

## Next action

Read the consumer findings. Inspect generated alwaysApply core behavior vs glob-scoped intake. Decide keep / reword / split skinny vs full.
