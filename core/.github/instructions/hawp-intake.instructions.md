---
applyTo: ".hawp/**,**/.hawp/**"
---

# HAWP Modular Intake

Use this as a scoped, drop-in instruction when integrating HAWP into repos that already have Copilot instructions and prompts.

## Integration Intent

- Layer HAWP behavior on top of existing repository rules.
- Do not replace or duplicate repository baseline instructions.
- Keep HAWP guidance scoped to HAWP-related folders and workflows.

## Source Resolution Order

When a task asks for status, handoff, bounded review, audit, planning, or comparison output, resolve HAWP guidance in this order:

1. `.hawp/kit/START_HERE.md` and `.hawp/kit/templates/status-report.md` (installable operating layer)
2. `.hawp/kit/README.md`, `.hawp/kit/SPEC.md`, and `.hawp/kit/AUTHORING_PATTERNS.md` (protocol reference)
3. `.hawp/kit/usage/INIT.md` and `.hawp/kit/usage/STATUS_REPORT.md` only when present as repo-local extensions

If a higher-priority path is missing, continue to the next source without failing.

## Behavior Rules

- Keep outputs compact and decision-useful.
- Separate direct evidence from inference for substantive findings.
- Avoid inventing runtime engines, validators, orchestrators, or memory systems.
- Do not introduce per-field runtime folders such as `input/`, `context/`, `mission/`, `constraints/`, or `output/` unless explicitly required by repo policy.

## Non-Conflict Rule

- If existing repository instructions conflict with this file, prefer the more specific path-scoped rule.
- Treat this file as a modular overlay, not a global override.
