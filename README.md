# human-ai-workflow-protocol

A minimal protocol for reliable human-AI collaboration. HAWP reduces task drift by giving work a small, durable shape before execution begins.

## What it is

HAWP is a compact task-shaping protocol with five required fields and one optional checkpoint field:

```ts
type Shape = {
  input: string; // the request as received
  context: string; // minimal background needed to interpret it
  mission: string; // the concrete objective for this run
  checkpoint?: string; // optional lightweight pause marker
  constraints: string; // hard boundaries and quality bars
  output: string; // what done looks like
};
```

It is not a runtime engine, memory system, validator, or orchestration framework. It is a lean shaping protocol.
Public examples in this repository are intended to stay fictional or clearly generic so the core remains portable across projects.

## Repository layout

```
benchmark/
  README.md               — benchmark purpose and install note
  benchmark-prompt.md     — guide for running a HAWP vs no-HAWP comparison
core/
  .hawp/
    LICENSE                 — Apache 2.0 license text copied with the kit
    README.md               — what HAWP is and why it exists
    SPEC.md                 — v0.1 field semantics and pipeline draft
    AUTHORING_PATTERNS.md   — compact authoring guidance for recurring task types
    templates/              — optional starter templates for common HAWP artifacts
    patterns/               — optional evidence/reporting patterns
    reviews/                — optional project review checklists
    types/shape.ts          — locked TypeScript type for v0.1
    examples/               — concrete filled-shape examples
    usage/
      INIT.md               — repo-local operating guide
      STATUS_REPORT.md      — status report format and rules
      GUARDRAIL_ADR.md      — guardrail ADR for improvement without schema expansion
      status/               — saved status reports
  .github/
    install.md              — instructions for installing the GitHub Copilot overlay
    update.md               — instructions for updating HAWP from GitHub main
    copilot-instructions.md — Copilot instructions template (uses repo-root .hawp/ paths)
    instructions/           — scoped intake instructions for the overlay
    prompts/                — curated prompt pack for handoff, status, and docs alignment
```

## Using HAWP

1. Read `core/.hawp/README.md` to understand the protocol.
2. Use `core/.hawp/AUTHORING_PATTERNS.md` to fill the shape for your task type.
3. See `core/.hawp/examples/` for concrete filled-shape examples.
4. Read `core/.hawp/usage/INIT.md` for how this repo uses HAWP in practice.

Templates and patterns are optional usage aids. They do not expand the HAWP core protocol.
Publication-safety guidance lives in `core/.hawp/usage/PUBLICATION_SAFETY_ADR.md` and `core/.hawp/reviews/public-safety-checklist.md`.

## Installing HAWP into another repository

To install HAWP into a target repository, copy only these two folders into `core/`:

- `.github/` — GitHub Copilot overlay files
- `.hawp/` — protocol docs, authoring patterns, examples, usage guidance, and `.hawp/LICENSE`

The root-level `benchmark/` folder is reference material for evaluating HAWP's practical value. Copy it into a target repository only if you intend to use it there. It is not required for HAWP to function.

For GitHub Copilot integration setup, follow `core/install.md`.

To refresh an already-installed setup from upstream `main`, use `core/update.md`.

## Quick install + usage (copy/paste)

If you want a single copy/paste flow in a target repo, use this:

> Replace `OWNER`, `REPO`, and `REF` if installing from a fork, internal mirror, or pinned release branch.

```bash
# From the target repository root
OWNER="sentzunhat"
REPO="human-ai-workflow-protocol"
REF="main"
KIT_DIR="core"  # optional: set to any folder name you prefer

TMP_DIR="$(mktemp -d)"
curl -fsSL "https://github.com/${OWNER}/${REPO}/archive/refs/heads/${REF}.tar.gz" \
  | tar -xz -C "$TMP_DIR"

SRC="$TMP_DIR/${REPO}-${REF}/core"

mkdir -p "$KIT_DIR" .github/instructions .github/prompts
cp -R "$SRC/.hawp" "$KIT_DIR"/
cp "$SRC/.github/instructions/"*.instructions.md \
  .github/instructions/
cp "$SRC/.github/prompts/"*.prompt.md \
  .github/prompts/

# If .github/copilot-instructions.md does not exist, seed it.
# If it already exists, merge HAWP guidance instead of overwriting.
if [ ! -f .github/copilot-instructions.md ]; then
  cp "$SRC/.github/copilot-instructions.md" .github/copilot-instructions.md
fi

rm -rf "$TMP_DIR"
```

That copy brings along `$KIT_DIR/.hawp/LICENSE`, so the installed HAWP folder carries the Apache 2.0 license text with it.

Then do this minimal usage wiring:

1. Ensure `.github/copilot-instructions.md` references `$KIT_DIR/.hawp/usage/INIT.md` and `$KIT_DIR/.hawp/usage/STATUS_REPORT.md`.
2. Use `$KIT_DIR/.hawp/START_HERE.md` to shape a task.
3. Use `$KIT_DIR/.hawp/usage/STATUS_REPORT.md` when you need a checkpoint/context-transfer artifact.

If your target repo uses repo-root `.hawp/` instead of `$KIT_DIR/.hawp/`, adapt paths in the installed `.github/` files accordingly.

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](./LICENSE).
