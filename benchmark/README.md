# HAWP Benchmark

This folder contains a simple human-runnable benchmark for evaluating whether HAWP-guided work is better than unguided work.

## Quick start

Follow these three steps in order (two Cursor windows; stop between arms):

1. **[instructions/setup.md](instructions/setup.md)** — original repo: checks, choose task, launch clean copy
2. **[instructions/run.md](instructions/run.md)** — clean `/tmp/` workspace: no-HAWP arm only; close window when done
3. **[instructions/cleanup.md](instructions/cleanup.md)** — original repo: HAWP arm, score, record, then remove `/tmp/` copies

Reference while running: **[benchmark-prompt.md](benchmark-prompt.md)** (task types, scoring dimensions, interpretation).

Completed examples: **[runs/](runs/)**.

## Purpose

HAWP is a workflow protocol, not a runtime system. As such, its value cannot be measured by running automated tests. It can only be measured by comparing real work outputs — and asking whether the HAWP-guided output was more bounded, more trustworthy, and more useful than the unguided one.

This benchmark is a human evaluation aid. It is not a scoring framework, a scientific instrument, or an automated suite.

## What this folder contains

| File / folder | Purpose |
| --- | --- |
| [instructions/setup.md](instructions/setup.md) | Setup: choose task, launch clean workspace (original repo) |
| [instructions/run.md](instructions/run.md) | No-HAWP arm in clean `/tmp/` workspace — stop before HAWP |
| [instructions/cleanup.md](instructions/cleanup.md) | HAWP arm, scoring, recording, then cleanup (original repo) |
| [benchmark-prompt.md](benchmark-prompt.md) | Task types, scoring dimensions, interpretation |
| [prepare-clean-workspace.sh](prepare-clean-workspace.sh) | Fair no-HAWP workspace: copy repo to `/tmp/`, strip agent rules from copy only |
| [runs/](runs/) | Saved comparison runs (outputs, scorecards, comparisons) |

## Installation note

This folder is part of this repository's content and is not part of the minimal installation pair.

When installing HAWP into a target repository, copy only:

- `.github/`
- `.hawp/` (includes `.hawp/LICENSE` with the Apache 2.0 text)

into `core/`.

The `benchmark/` folder is a reference resource for teams who want to evaluate HAWP's practical value. Copy it into a target repository only if you intend to use it there.

## Related docs

- [.hawp/kit/README.md](../core/.hawp/kit/README.md) — what HAWP is
- [.hawp/kit/references/authoring-patterns.md](../core/.hawp/kit/references/authoring-patterns.md) — guidance for specific task types
- [.hawp/kit/examples/](../core/.hawp/kit/examples/) — concrete filled-shape examples
