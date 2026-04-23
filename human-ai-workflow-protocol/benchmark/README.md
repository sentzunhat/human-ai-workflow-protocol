# HAWP Benchmark

This folder contains a simple human-runnable benchmark for evaluating whether HAWP-guided work is better than unguided work.

## Purpose

HAWP is a workflow protocol, not a runtime system. As such, its value cannot be measured by running automated tests. It can only be measured by comparing real work outputs — and asking whether the HAWP-guided output was more bounded, more trustworthy, and more useful than the unguided one.

This benchmark is a human evaluation aid. It is not a scoring framework, a scientific instrument, or an automated suite.

## What this folder contains

- `benchmark-prompt.md` — a reusable prompt and guide for running a HAWP vs no-HAWP comparison across practical task types

## Installation note

This folder is part of this repository's content and is not part of the minimal installation pair.

When installing HAWP into a target repository, copy only:

- `github/`
- `hawp/` (includes `hawp/LICENSE` with the Apache 2.0 text)

into `human-ai-workflow-protocol/`.

The `benchmark/` folder is a reference resource for teams who want to evaluate HAWP's practical value. Copy it into a target repository only if you intend to use it there.

## Related docs

- [hawp/README.md](../hawp/README.md) — what HAWP is
- [hawp/AUTHORING_PATTERNS.md](../hawp/AUTHORING_PATTERNS.md) — guidance for specific task types
- [hawp/examples/](../hawp/examples/) — concrete filled-shape examples
