# cursor-docs-ui-path — Fix stale Cursor MCP UI path in kit docs

**Type:** docs  
**Status:** done  
**Opened:** 2026-08-25  
**Closed:** 2026-08-25
**Target:** v0.0.11

## Goal

Kit docs show the correct Cursor UI path and include the reload note.
No stale `Settings → Tools & MCP` references remain.

## Outcome

The kit search guide now tells users to enable the workspace server in
sidebar `Customize → MCPs` and to open a new chat or reload Cursor if tools
still do not appear. Within inspected scope, no stale `Settings → Tools & MCP`
references remain under `.hawp/kit/` or `distribution/`.

## Verification

- [x] Updated `.hawp/kit/usage/search.md` with the confirmed current Cursor path
- [x] `rg -n "Tools & MCP|Settings.*Tools.*MCP|Tools.*MCP.*Settings" .hawp/kit distribution` returned no matches
- [x] Distribution outputs were regenerated after the docs/source update

## What was done

- Added the current Cursor MCP enable path and reload guidance to `.hawp/kit/usage/search.md`
- Confirmed stale Cursor UI wording is absent in the checked kit and distribution docs
