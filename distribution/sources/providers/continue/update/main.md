# Update HAWP — Continue Provider (Main Branch)

Stable update of HAWP kit plus Continue overlays: refresh `.continue/rules/hawp-*.md`.

Requires an existing `.hawp/` install in the target repository.

## Update Steps

1. Open your target repository root in a terminal.
2. Run the **Update Command (Copy/Paste)** block below (`REF="main"`, `PROVIDER="continue"`).
3. Confirm `.hawp/kit/` and `.continue/rules/hawp-*.md` were refreshed.

Optional: `export HAWP_LOCAL_CORE="/absolute/path/to/human-ai-workflow-protocol/core"` for local testing.

## What Was Refreshed

- `.hawp/kit/**`, `.hawp/LICENSE`
- `.continue/rules/hawp-*.md` from `core/providers/.continue/rules/`

## What Was NOT Changed

- Non-HAWP rules under `.continue/rules/`, `.hawp/work/**` project records.

## Other guides

- Dev branch: `distribution/generated/continue/update/dev.md`
