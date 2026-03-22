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
```

## Using HAWP

1. Read `human-ai-workflow-protocol/hawp/README.md` to understand the protocol.
2. Use `human-ai-workflow-protocol/hawp/AUTHORING_PATTERNS.md` to fill the shape for your task type.
3. See `human-ai-workflow-protocol/hawp/examples/` for concrete filled-shape examples.
4. Read `human-ai-workflow-protocol/hawp/usage/INIT.md` for how this repo uses HAWP in practice.

## GitHub Copilot integration

To add HAWP-aware Copilot behavior to a target repository that already has a `hawp/` folder, follow `human-ai-workflow-protocol/github/install.md`.
