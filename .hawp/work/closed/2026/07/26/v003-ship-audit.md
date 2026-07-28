---
type: audit
title: "v0.0.3 ship-readiness audit — 7 findings fixed and re-verified"
status: done
closed: 2026-07-26
---

# v0.0.3 ship-readiness audit — closed record

This is a **shared closed record**: one audit session fixed and re-verified
7 backlog rows at once. Rather than 7 near-identical files, this single
record names every row it closes (`work-validate` matches a closed record
by content, not only by filename — see `findClosedFile` /
`recordListsId` in `librarian/scripts/hawp/work-validate/validations/backlog-consistency.ts`).

**Closes:** `w4x6y8z0`, `a1b3c5d7`, `e9f1g3h5`, `i7j9k1l3`, `j2k4l6m8`,
`n4o6p8q0`, `r5s7t9u1`

Full findings, root causes, and executed proof for each fix:
[SHIP_READY_VERIFICATION.md](../../../../../librarian/docs/SHIP_READY_VERIFICATION.md)

## Outcome (filled at close)

All 7 findings from the 2026-07-26 audit were fixed and re-verified against
real running backends (Ollama + ONNX Runtime, no mocks):

- `w4x6y8z0` — `--llm-reshape` wired into the CLI (was unreachable)
- `n4o6p8q0` — `ReshapingConfig` now honors the configured Ollama URL
- `j2k4l6m8` — `DefaultConfig()` no longer defaults to a nonexistent ONNX LLM model
- `a1b3c5d7` — 8 of 10 failing Ollama LLM tests repaired (mock endpoint mismatch)
- `e9f1g3h5` — token budget display (`1995/0`) fixed
- `i7j9k1l3` — `aes256` silent downgrade now rejected at config validation
- `r5s7t9u1` — `Reshape()`'s empty-block guard now checks `len(block.Results)`, not `String() == ""`

## Verification (filled at close)

Each fix in [SHIP_READY_VERIFICATION.md](../../../../../librarian/docs/SHIP_READY_VERIFICATION.md)
carries **Confirmed** (directly executed/observed) or **Likely**
(evidence-backed but interpretive) labels per claim — see that document's
"Method" section. Executed proof includes real `hawp search --llm-reshape`
runs against live Ollama/ONNX backends, not just source re-reading.

## Close Checklist

- [x] Outcome recorded (all 7 findings, one per bullet above)
- [x] Verification recorded (see linked audit doc for full command output)
- [x] Follow-on work identified — none; audit found no further blockers
