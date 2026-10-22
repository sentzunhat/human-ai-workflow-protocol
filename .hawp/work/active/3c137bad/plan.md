# Propagate MCP tool groups and resource/prompt guidance through provider packs and distributions

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `3c137bad-e686-4bab-a070-26fa6b821962`
**Type:** improvement
**Reported:** 2026-10-16
**Status:** in-progress
**Target Version:** Next release candidate; currently proposed as `v0.0.26` (not locked)

---

## Input (verbatim)

> After the HAWP MCP server exposes and verifies curated resources/prompts, update shared and provider-specific guidance so existing tool groups (Discover, Shape, Create, Verify) are consistent across Codex, Claude Code, Cursor, Continue, and GitHub Copilot surfaces. Regenerate and validate distribution guides from canonical provider sources, document client-specific feature limits, and verify installed-file propagation paths. This is a later follow-up; do not implement it as part of the server feature work.

## Release Direction (2026-10-17)

The maintainer directed that provider prompts and related distribution
propagation ship together in the next release. Keep this work in this existing
item; do not create a duplicate provider-prompt task. The latest local release
map proposes `v0.0.26` for MCP client proof and catalog closeout, so treat that
as the current candidate lane pending live backlog/release review. No release
date, lock, or claim of client support is implied.

## Intake Summary

Maintainer resumed provider and distribution propagation on 2026-10-16 after
the server implemented resources/prompts and protocol-level tests passed.
Client-specific smoke results remain absent, so the guidance explicitly
separates vendor-documented client UX from HAWP integration proof.

## Current Context

- Canonical HAWP MCP tool groups are documented in
  `.hawp/kit/usage/mcp/README.md`: Discover, Shape, Create, Verify.
- Provider packs live under `core/providers/`; generated public guides are
  authored in `distribution/sources/` and refreshed through the HAWP sync
  commands.
- HAWP server supports fixed resources and one prompt over stdio; focused
  protocol tests pass, but no supported client UI has been smoke-tested.
- Repository propagation guidance: `.github/skills/multi-repo-context.md`.

## Initial Analysis

**Directly verified:**

- The five provider overlays are Claude Code, Codex, GitHub Copilot, Cursor,
  and Continue, as listed in `core/providers/manifest.yaml`.
- `distribution/sources/providers/` owns generated provider install/update
  guides; `hawp providers sync` and `hawp distribution sync` are the expected
  regeneration steps for their respective sources.
- This item must follow server capability work and must not claim client
  support beyond documented or directly tested surfaces.

**Inferred (not yet proven):**

- Shared grouping can likely be introduced through canonical shared behaviors,
  with provider-specific notes for prompt/resource UX and unsupported surfaces.
- Generated distributions may need source updates if provider capabilities are
  described in installation/update workflows.

**Likely scope:**

- Sync canonical MCP tool grouping and reusable prompt guidance across relevant provider guidance.
- Document prompt discovery/use and resource discovery/use only on client
  surfaces where verified; state fallback behavior elsewhere.
- Regenerate provider and distribution outputs; inspect generated diffs.
- Verify install/update propagation boundaries without modifying a downstream
  consumer repository.
- Do not change server behavior or add new MCP capabilities in this item.

## Risk + Review Gate

**Risk:** medium
**Gate:** review generated provider/distribution changes and client-specific
claims before publishing.

## Backlog + Plan Link

**Status now:** in-progress
**Plan file:** work/active/3c137bad/plan.md

## Plan

1. [x] Reconcile the resource and prompt contracts from the server plans and
   capability catalog.
2. [x] Review current vendor docs for Claude Code, Copilot in VS Code, Cursor,
   Continue, and Codex; distinguish documented client UX from HAWP runtime
   evidence.
3. [x] Add consistent tool-group and conditional resource/prompt guidance to
   shared provider behavior and Codex/GitHub source instructions.
4. [x] Update the provider capability matrix and shared install/update sources;
   regenerate provider and distribution outputs and inspect diffs.
5. [x] Validate HAWP kit/work/distribution consistency and inspect propagation
   paths. Do not modify a downstream consumer repository.

## Verification

- `providers sync` passed and confirmed 15 materialized provider files current.
- `distribution sync` updated 20 of 40 generated guides; generated distribution
  validation passed.
- `kit validate`, `work validate` (6 active / 13 parked, 0 issues),
  `distribution validate`, `links check` (149 Markdown files), `hawp check`,
  and `git diff --check` passed.
- No HAWP live client smoke has been run; guidance says to use resources and
  prompts only when the connected client lists them.

## Backlog / Parking

- **Status:** `in-progress`; client verification remains an explicit follow-up.
- **Owner:** maintainer
- **Remaining proof:** run a bounded live HAWP discovery/read and prompt smoke
  in a supported client. Current guidance is conservative and points clients
  without those visible capabilities to the manual guide and available tools.
