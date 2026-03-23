# human-ai-workflow-protocol

A minimal protocol for reliable human-AI collaboration. HAWP reduces task drift by giving work a small, durable shape before execution begins.

## What it is

HAWP is a six-field carrier for intent and task framing:

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

## Repository layout

```
human-ai-workflow-protocol/
  hawp/
    README.md               — what HAWP is and why it exists
    SPEC.md                 — v0.1 field semantics and pipeline draft
    AUTHORING_PATTERNS.md   — compact authoring guidance for recurring task types
    src/shape.ts            — locked TypeScript type for v0.1
    examples/               — concrete filled-shape examples
    usage/
      INIT.md               — repo-local operating guide
      STATUS_REPORT.md      — status report format and rules
      status/               — saved status reports
  github/
    install.md              — instructions for installing the GitHub Copilot overlay
    copilot-instructions.md — Copilot instructions template (uses repo-root hawp/ paths)
    instructions/           — scoped intake instructions for the overlay
    prompts/                — status report prompt for the overlay
  benchmark/
    README.md               — benchmark purpose and install note
    benchmark-prompt.md     — guide for running a HAWP vs no-HAWP comparison
```

## Using HAWP

1. Read `human-ai-workflow-protocol/hawp/README.md` to understand the protocol.
2. Use `human-ai-workflow-protocol/hawp/AUTHORING_PATTERNS.md` to fill the shape for your task type.
3. See `human-ai-workflow-protocol/hawp/examples/` for concrete filled-shape examples.
4. Read `human-ai-workflow-protocol/hawp/usage/INIT.md` for how this repo uses HAWP in practice.

## Installing HAWP into another repository

To install HAWP into a target repository, copy only these two folders into `human-ai-workflow-protocol/`:

- `github/` — GitHub Copilot overlay files
- `hawp/` — protocol docs, authoring patterns, examples, and usage guidance

The `benchmark/` folder is reference material for evaluating HAWP's practical value. Copy it into a target repository only if you intend to use it there. It is not required for HAWP to function.

For GitHub Copilot integration setup, follow `human-ai-workflow-protocol/github/install.md`.
