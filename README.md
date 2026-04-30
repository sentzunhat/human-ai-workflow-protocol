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
    LICENSE                       — Apache 2.0 license text copied with the kit
    kit/                          — reusable HAWP material (installed downstream)
      README.md                   — what HAWP is and why it exists
      START_HERE.md               — fastest on-ramp / minimal template
      SPEC.md                     — v0.1 field semantics and pipeline draft
      AUTHORING_PATTERNS.md       — compact authoring guidance for recurring task types
      templates/                  — optional starter templates (includes empty backlog template)
      patterns/                   — optional evidence/reporting patterns
      reviews/                    — optional project review checklists
      examples/                   — concrete filled-shape examples
      types/shape.ts              — locked TypeScript type for v0.1
      usage/                      — reusable workflow guides
        INIT.md                   — operating guide for HAWP in a repo
        INTAKE_WORKFLOW.md        — intake loop for bugs, tasks, and improvements
        STATUS_REPORT.md          — status report guidance
    work/                         — repo-local operating state (not installed downstream)
      README.md                   — explains BACKLOG/adrs/status/evidence
      BACKLOG.md                  — single source of truth for open and completed work
      adrs/                       — repo-local ADRs (e.g. guardrail, install boundary)
      status/                     — plan files and status reports per backlog item
      evidence/                   — verification artifacts referenced from plans
  .github/
    install.md              — instructions for installing the GitHub Copilot overlay
    update.md               — instructions for updating HAWP from GitHub main
    copilot-instructions.md — Copilot instructions template (uses repo-root .hawp/ paths)
    instructions/           — scoped intake and docs-alignment instructions
      intake.instructions.md                              — ambient intake trigger (auto-activates on bug/task reports)
      human-ai-workflow-protocol-intake.instructions.md  — HAWP modular intake overlay
      human-ai-workflow-protocol-docs-alignment.instructions.md — docs alignment overlay
    prompts/                — curated prompt pack for intake, handoff, status, and docs alignment
      intake.prompt.md                                           — explicit intake loop trigger
      human-ai-workflow-protocol-status-report.prompt.md
      human-ai-workflow-protocol-intent-first-handoff.prompt.md
      human-ai-workflow-protocol-docs-alignment-deterministic.prompt.md
      human-ai-workflow-protocol-docs-alignment-simplicity.prompt.md
      human-ai-workflow-protocol-conservative-docs-drift-cleanup.prompt.md
```

## Using HAWP

1. Read `core/.hawp/kit/README.md` to understand the protocol.
2. Use `core/.hawp/kit/AUTHORING_PATTERNS.md` to fill the shape for your task type.
3. See `core/.hawp/kit/examples/` for concrete filled-shape examples.
4. Start from `core/.hawp/kit/START_HERE.md` for task shaping and `core/.hawp/kit/templates/status-report.md` for continuity reports.

Templates and patterns are optional usage aids. They do not expand the HAWP core protocol.
Publication-safety guidance lives in `core/.hawp/kit/reviews/public-safety-checklist.md` and `core/.hawp/kit/reviews/publication-safety-guidelines.md`.

## Installing HAWP into another repository

To install HAWP into a target repository, install the reusable `.hawp/` kit files and the `.github/` Copilot overlay. The default install flow uses an allowlist and does not copy `core/.hawp/work/`.

Installed by default:

- `.github/` — GitHub Copilot overlay files
- `.hawp/LICENSE` and `.hawp/kit/` — protocol docs, authoring patterns, templates, reviews, patterns, examples, types, and reusable workflow guides under `kit/usage/`

Not installed by default:

- `core/.hawp/work/` — source-repo backlog, ADRs, status reports, and evidence

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
mkdir -p "$KIT_DIR/.hawp/kit"

cp "$SRC/.hawp/LICENSE" "$KIT_DIR/.hawp/"
cp "$SRC/.hawp/kit/README.md" "$KIT_DIR/.hawp/kit/"
cp "$SRC/.hawp/kit/START_HERE.md" "$KIT_DIR/.hawp/kit/"
cp "$SRC/.hawp/kit/SPEC.md" "$KIT_DIR/.hawp/kit/"
cp "$SRC/.hawp/kit/AUTHORING_PATTERNS.md" "$KIT_DIR/.hawp/kit/"

cp -R "$SRC/.hawp/kit/templates" "$KIT_DIR/.hawp/kit/"
cp -R "$SRC/.hawp/kit/patterns" "$KIT_DIR/.hawp/kit/"
cp -R "$SRC/.hawp/kit/reviews" "$KIT_DIR/.hawp/kit/"
cp -R "$SRC/.hawp/kit/examples" "$KIT_DIR/.hawp/kit/"
cp -R "$SRC/.hawp/kit/types" "$KIT_DIR/.hawp/kit/"
cp -R "$SRC/.hawp/kit/usage" "$KIT_DIR/.hawp/kit/"
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

1. Ensure `.github/copilot-instructions.md` references `$KIT_DIR/.hawp/kit/START_HERE.md` and `$KIT_DIR/.hawp/kit/templates/status-report.md`.
2. Use `$KIT_DIR/.hawp/kit/START_HERE.md` to shape a task.
3. Use `$KIT_DIR/.hawp/kit/templates/status-report.md` when you need a checkpoint/context-transfer artifact.

If your target repo uses repo-root `.hawp/` instead of `$KIT_DIR/.hawp/`, adapt paths in the installed `.github/` files accordingly.

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](./LICENSE).
