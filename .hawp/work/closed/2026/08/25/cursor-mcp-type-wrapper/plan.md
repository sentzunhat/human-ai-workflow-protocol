# cursor-mcp-type-wrapper — Cursor MCP: add `type: stdio` and `hawp-mcp` wrapper

**Type:** fix  
**Status:** done  
**Opened:** 2026-08-25  
**Closed:** 2026-08-25
**Target:** v0.0.11

## Goal

`hawp init --provider cursor` generates a `.cursor/mcp.json` that Cursor
actually spawns. Also: ship a small `hawp-mcp` wrapper script alongside the
hawp binary that cds to the repo root before executing `hawp mcp`.

## Outcome

Cursor MCP config generation is now provider-specific. `hawp init --provider cursor`
writes `.cursor/mcp.json` with `type: stdio` and an absolute command path to a
repo-local `hawp-mcp` wrapper. The wrapper changes to repo root before execing
`hawp mcp`, which preserves the repo-root assumption inside `runMCP()`. Claude
Code keeps the existing `.mcp.json` shape, and Codex remains unchanged.

The install/update distribution scripts now copy `hawp-mcp` alongside the main
`hawp` binary, and the kit search docs now explain the current Cursor enable path:
`Customize → MCPs`, workspace toggle on, then open a new chat or reload if tools
do not appear immediately.

## Verification

- [x] `go test ./internal/platform/mcp/...`
- [x] `go test ./internal/platform/cli/...`
- [x] `mise exec node@26.5.0 -- npm run distribution:sync`
- [x] Cursor stale UI-path grep under `.hawp/kit/` and `distribution/` returned no matches after doc updates
- [x] Generated distribution outputs validated current after source changes

## What was done

- Added `core/.hawp/bin/hawp-mcp` wrapper that `cd`s to repo root and execs `hawp mcp`
- Split JSON MCP config generation into provider-specific entries
- Kept Claude Code on relative `.hawp/bin/hawp` + args, while Cursor now gets `type: stdio` and absolute wrapper path
- Updated Cursor-oriented notes in `.hawp/kit/usage/search.md`
- Updated install/update source templates so downstream repos receive `hawp-mcp`
- Regenerated checked-in distribution outputs from source templates
