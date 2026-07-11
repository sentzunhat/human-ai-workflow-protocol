# Librarian CLI Architecture Using Zacatl Patterns

Related work: `.hawp/work/active/8dedc4e2-69c5-42ca-aaad-93b62d7fb899.md`

## Purpose

Shape the installable librarian CLI architecture before implementation. This
note uses Zacatl as the architecture reference while preserving the current
HAWP librarian constraints:

- HAWP is a workflow layer, not a runtime engine.
- `librarian/` is currently repo-local maintenance tooling.
- `.hawp/bin/hawp` is a repo-local wrapper, not the published CLI contract.
- Runtime/model and single-binary packaging questions stay in their own spikes.

## Direct Evidence

Confirmed from this repository:

- `librarian/scripts/README.md` already defines the command-domain boundary:
  `index.ts` owns executable process concerns, `cli.ts` owns argument parsing
  and help text, and `script.ts` owns reusable logic.
- `librarian/scripts/README.md` separates HAWP workflow tooling under
  `scripts/hawp/`, build/distribution pipelines under `scripts/librarian/`, and
  shared helpers under `scripts/lib/`.
- `librarian/package.json` is private and has no `bin` field today.
- `README.md` and `librarian/README.md` describe `librarian/` as repo-local and
  `.hawp/bin/hawp` as a thin wrapper.
- A narrow `@oclif/core` PoC can route `hawp-poc kit validate` through
  application/domain/infrastructure/platform folders while reusing the existing
  kit validation logic.
- esbuild can produce a bundled CJS executable artifact at
  `librarian/build-bin/hawp-poc.cjs`; the observed artifact is 10.2 MB.
- Node 26 SEA supports direct single-executable generation with `--build-sea`.
  Cross-platform SEA builds should keep `useSnapshot` and `useCodeCache` false.
- The PoC target matrix is Linux x64/ARM64, macOS x64/ARM64, and Windows
  x64/ARM64. Local builds are host-target only unless target Node 26 executables
  are supplied; CI builds use matching hosted runners.

Confirmed from Zacatl:

- Zacatl documents a layered architecture: application, domain,
  infrastructure, and platform boundaries.
- Zacatl's `ServiceType.CLI` is declared but not runnable yet; its CLI platform
  throws explicit "not implemented" errors.
- Zacatl recommends standalone dependency-injection helpers for non-HTTP
  scripts and jobs today, rather than depending on the unfinished CLI runtime.
- Zacatl publish tooling already treats `bin` as a package-surface concern that
  needs explicit build/publish handling.

## Inference

Likely architecture direction:

- Use Zacatl's layering vocabulary and boundary discipline, not Zacatl's current
  CLI runtime implementation.
- Make `@hawp/librarian` an installable package only after the command contract
  is explicit and current repo-local wrapper behavior can be preserved.
- Introduce a publishable CLI entrypoint that routes to command modules instead
  of growing `.hawp/bin/hawp` into the product architecture.
- Keep command implementations library-first so future binary packaging and
  model/runtime commands can compose the same modules without process coupling.
- Prefer a CJS bundled executable for the near-term Node CLI PoC. Raw `tsc` ESM
  output needs import-extension handling, and the ESM esbuild bundle hit dynamic
  CommonJS require behavior inside `@oclif/core`.

## Proposed Layers

```text
librarian/
  src/
    application/
      commands/
        kit/
        work/
        providers/
        distribution/
      command-registry.ts
    domain/
      kit/
      work/
      providers/
      distribution/
    infrastructure/
      filesystem/
      markdown/
      provider-packs/
      distribution-writer/
    platform/
      cli/
        index.ts
        parse.ts
        render-help.ts
        run.ts
    lib/
      index.ts
  bin/
    hawp.ts
  scripts/
    ...
```

This is a target shape, not a required immediate folder move. The first
implementation pass can route existing `scripts/hawp/*` domains through a small
platform CLI adapter while leaving build/distribution scripts in place.

PoC mapping:

```text
librarian/scripts/hawp-cli-poc/
  bin/hawp-poc.ts
  platform/cli/
  application/commands/kit/validate.ts
  domain/kit/validate-kit.ts
  infrastructure/filesystem/hawp-repo.ts
```

## Command Surface

Installable CLI name:

- Primary: `hawp`
- Package: keep `@hawp/librarian` as the likely package identity unless a
  publish decision changes it.
- Repo-local wrapper: keep `.hawp/bin/hawp` as a compatibility shim that invokes
  the local workspace.

Initial stable command groups:

| Command | Status | Notes |
| --- | --- | --- |
| `hawp kit validate` | publishable | Existing domain; read-only validation. |
| `hawp kit normalize` | publishable with guardrails | Mutating; must preserve clean-tree guard and dry-run/apply semantics. |
| `hawp work validate` | publishable | Existing active-work integrity command. |
| `hawp work normalize` | publishable with guardrails | Existing raw CLI behavior should remain available. |
| `hawp backlog validate` | compatibility alias | Alias to work validation/check behavior; avoid making backlog a separate domain. |
| `hawp backlog upgrade` | compatibility alias | Alias to `work normalize`; deprecate wording later if needed. |
| `hawp uuid` | publishable | Small utility, no repo dependency. |

Repo-local or deferred:

- `providers:*` and `distribution:*` stay repo-maintenance scripts until the
  package has a clear maintainer-only/admin command story.
- AI/model commands stay out of this CLI architecture slice and belong to the
  Node/WASM shared-cache spike.
- Single-binary packaging stays out of this slice and belongs to the packaging
  strategy spike.

## Bin Contract

The published package should eventually expose:

```json
{
  "bin": {
    "hawp": "build/bin/hawp.js"
  }
}
```

Contract requirements:

- `bin/hawp.ts` must be a thin executable boundary: argv in, exit code out.
- Argument parsing and help rendering stay in `platform/cli/`.
- Command handlers return structured results and exit codes.
- Filesystem access belongs below application/domain boundaries, not in parsing.
- Repo discovery must be explicit and testable; commands that require a HAWP
  repo should say so clearly.
- Mutating commands must fail closed when the worktree is dirty unless the
  existing command-specific guard says otherwise.

## Migration Sequence

1. Document the command taxonomy and repo-local vs publishable boundary.
2. Add a package/bin implementation plan without changing runtime behavior.
3. Introduce a central CLI router that delegates to existing command domains.
4. Add tests for routing, help output, aliases, unknown commands, and exit codes.
5. Switch `.hawp/bin/hawp` to call the same router locally.
6. Only after parity is proven, consider `private: false`, `files`, `exports`,
   build output, and publish preparation.

PoC packaging lessons:

- `npm --prefix librarian run cli:poc` runs the source-mode proof through `tsx`.
- `npm --prefix librarian run cli:poc:bundle` builds a bundled CJS executable
  artifact with esbuild.
- `npm --prefix librarian run cli:poc:node -- kit validate` runs the bundled
  artifact with Node.
- `npm --prefix librarian run cli:poc:binary` attempts a Node 26 SEA binary for
  the current host target.
- `npm --prefix librarian run cli:poc:binaries -- --target all` writes a target
  manifest and attempts all configured platform targets.
- `.github/workflows/hawp-cli-poc-binaries.yml` is the path for real
  all-platform binary artifacts.
- Node 26 proof still needs a real Node 26 runtime in the shell. Current proof
  was executed under Node v24.14.0, where oclif emitted a Node >=26 warning.

## Verification Plan

For the implementation item that follows this architecture note:

- `npm --prefix librarian run typecheck`
- `npm --prefix librarian test`
- `./.hawp/bin/hawp kit validate`
- `./.hawp/bin/hawp work validate`
- focused CLI route tests for stable commands and aliases

## Open Decisions

- Whether `providers` and `distribution` ever become public CLI commands or
  stay maintainer-only npm scripts.
- Whether the package stays named `@hawp/librarian` or moves to a user-facing
  package name.
- Whether `backlog` aliases are documented as permanent compatibility aliases
  or a transition path toward `work`.
