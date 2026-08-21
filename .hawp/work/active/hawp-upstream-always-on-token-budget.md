# Always-on token budget and install-session context growth

**From:** a consumer infrastructure repository (follow-on review)
**Type:** improvement (upstream review; not a consumer patch)
**Reported:** 2026-08-21
**Status:** inbox
**Risk:** low if docs-only; higher if it becomes a kit-tree or install-product change

## Input

> HAWP upstream: document always-on token budget and install-session context growth for consumer repos.

## Context

Consumer findings: `.hawp/work/active/hawp-consumer-token-usage-findings.md`. Evidence A is ~800 HAWP tokens always on (hawp-core, hawp-backlog-alignment, CLAUDE.md). Evidence C is a 6.1 MB / 1,024-turn install chat (Jul 28 to Aug 11). Median cache ~210k to ~788k is not explained by A.

## Mission

Give consumer installers a written expectation: what is always injected, what is on-demand, and that long install threads dwarf alwaysApply files.

## Constraints

Docs and product notes only until scoped. Do not rewrite the kit as part of filing this item. Optional skinny vs full kit is a later product decision.

## Next action

Draft maintainer-facing install notes from the consumer findings. Keep consumer kit installs unchanged until that draft is approved.
