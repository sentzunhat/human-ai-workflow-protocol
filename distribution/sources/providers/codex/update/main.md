# Update HAWP — Codex Provider (Main Branch)

Stable update of HAWP kit plus Codex `AGENTS.md` instructions.

## Update Steps

1. Open your target repository root in a terminal.
2. Run the **Update Command (Copy/Paste)** block below (`REF="main"`, `PROVIDER="codex"`).
3. Confirm `.hawp/kit/` reflects the selected branch and that `AGENTS.md` exists if you wanted the HAWP seed installed.
4. If your repo already had `AGENTS.md`, keep it and manually blend in any HAWP wording you want from the provider seed.

## What Was Updated

- `.hawp/kit/**` — refreshed from HAWP core.
- `AGENTS.md` — seeded from `core/providers/.codex/AGENTS.md.seed` only when missing.

## What Was NOT Changed

- Existing `AGENTS.md` content.
- `.github/**`
- `.cursor/**`
- `.continue/**`
- `.claude/**`
- Runtime CLI participant adapters.
- `.hawp/work/**` project records.

## Other guides

- Development branch: `distribution/generated/codex/update/development.md`
- GitHub/Copilot: `distribution/generated/github/update/main.md`
- Cursor: `distribution/generated/cursor/update/main.md`
- Continue: `distribution/generated/continue/update/main.md`
