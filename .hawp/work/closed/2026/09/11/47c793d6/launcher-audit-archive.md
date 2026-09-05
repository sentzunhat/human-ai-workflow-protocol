# Native Launcher Boundary

## Observed Contract

- `.hawp/bin/hawp` is a Bash compatibility launcher, not the implementation.
- `.hawp/bin/hawp-bin` is the compiled Go implementation; Windows installs use
  `hawp-bin.exe`.
- `.hawp/bin/hawp-mcp` previously changed the working directory before starting
  MCP. Retired 2026-09-06 — superseded by `hawp mcp --repo-root`.
- Canonical install/update scripts still distribute the launchers and native
  executable together. Existing provider configurations may reference them.

## Implemented Step

`hawp mcp --repo-root <path>` selects the repository without changing process
working directory. New provider configurations launch the native executable
with that argument, including the Windows executable suffix. No new dependency
or CLI framework is introduced. Existing no-argument MCP invocation remains
supported. Configuration generation does not automatically migrate consumers.

## Remaining Work

The coordinated migration now has a dedicated work item:
[0cb0f9b0](../../closed/2026/09/07/0cb0f9b0/plan.md). JSON provider migration preserves custom settings
and rejects incompatible shapes. Codex now uses parser-backed edits of explicit
tables, preserving unrelated settings and surrounding text; unsupported layouts
are refused. Installer changes remain open; this is not a completed
single-executable migration.

A single executable is technically sufficient. Retiring the wrappers requires
coordinated canonical installer/update changes, regenerated distribution guides,
command references, and migration of existing provider configurations. Preserve
the wrappers until that compatibility work is tested. Cross-compilation alone
does not prove native execution or provider integration on every platform.
