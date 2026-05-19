# START_HERE

This is the fastest place to start with HAWP. Copy the template below, fill it in, and hand it to an AI agent. Keep it simple. Do not add fields.

---

## Minimal Template

```text
input: |
  <the request as received>

context: |
  <minimal background needed to interpret the input>

mission: |
  <one concrete objective derived from input + context>

checkpoint: |
  <optional pause point or handoff marker>

constraints: |
  <scope limits, required quality bars, and non-goals>

output: |
  <what done looks like — the expected artifact or response form>
```

`checkpoint` is optional. Add it only when you need a handoff marker or pause point.

---

## Usage Notes

- **context** — minimal but sufficient. Stable facts only. No transcript dumps.
- **mission** — one concrete objective. Name the lens (review, audit, plan, handoff).
- **constraints** — state what must not happen as clearly as what must.
- **output** — describe the artifact form, not just "an answer."
- **checkpoint** — omit unless you need a handoff or pause anchor.
- For path-sensitive work, reference files using exact repo-relative paths from repository root.
- Basename-only file mentions are unsafe for path-sensitive work unless the file is truly at repository root and explicitly marked as such.
- Capture repo-root proof (`pwd`, `git rev-parse --show-toplevel`, `git rev-parse --show-prefix`, `git status --short`) before path-sensitive edits.
- Do not persist machine-local absolute paths in artifacts; if proof output includes host-local prefixes, redact only the machine-local prefix with a placeholder such as `<repo-root-abs>`.

---

## Good First Uses

- repo review
- implementation plan
- handoff / status report
- gap analysis
- standards drift audit

---

## Evidence Discipline Quick Tips

- Separate directly observed facts from inference in your findings.
- Use `within inspected scope` or `no evidence found in inspected material` when claiming absence.
- Only call a standard mismatch when supported by docs, ADRs, contracts, or strong repeated conventions.
- Label confidence: `Confirmed`, `Likely`, or `Unclear`. Match to the weakest part of the claim.

---

## Boundary Reminder

HAWP is a **shaping protocol**, not a runtime. Better results come from stronger authoring, not more fields. The shape is locked to v0.1; quality improvements are achieved through guidance and examples.

## Next Resources

### Learning & Examples

- **Authoring guidance**: [authoring-patterns.md](authoring-patterns.md) — guidance for recurring task types
- **Real examples**: [examples/](examples/) — seven concrete filled-in HAWP shapes
- **Core spec**: [spec.md](spec.md) — v0.1 semantics, pipeline, and non-goals

### Workflow Guides (Recommended Order)

1. **[usage/init.md](usage/init.md)** — one-time setup for a new project
2. **[usage/intake-workflow.md](usage/intake-workflow.md)** — intake loop for structured bug/task handling
3. **[usage/status-report.md](usage/status-report.md)** — context handoff and session continuity

### Templates by Task Type

**Start here for shaping your first task:**

- **Micro** (constraints, ~1 page): [templates/micro-hawp.md](templates/micro-hawp.md)
- **Standard** (full form, 2–3 pages): [templates/standard-hawp.md](templates/standard-hawp.md)
- **Work intake** (bug reports, features): [templates/work-intake.md](templates/work-intake.md)
- **Bug plan** (troubleshooting + fix): [templates/bug-plan.md](templates/bug-plan.md)
- **Handoff** (context transfer): [templates/intent-first-handoff.md](templates/intent-first-handoff.md)
- **Status report** (checkpoint/review): [templates/status-report.md](templates/status-report.md)
- **Audit** (review findings): [templates/audit-report.md](templates/audit-report.md)

### Reference & Patterns

**Quality disciplines:**

- [patterns/evidence-discipline.md](patterns/evidence-discipline.md) — separate fact from inference
- [patterns/non-findings.md](patterns/non-findings.md) — how to report absence of findings
- [patterns/parallel-work-guardrails.md](patterns/parallel-work-guardrails.md) — multi-agent coordination

**Operational references:**

- [references/backlog-alignment.md](references/backlog-alignment.md) — backlog format, scope, closures
- [references/docs-alignment.md](references/docs-alignment.md) — keeping docs sync'd with reality
- [references/work-item-file-tracking.md](references/work-item-file-tracking.md) — file ownership in multi-agent work
- [references/install-update-safety.md](references/install-update-safety.md) — kit distribution and maintenance

### Standards & Domain Guidance

**Code and architecture:**

- [standards/guidelines/](standards/guidelines/) — general architecture, code style, documentation, testing, security, git
- [standards/nodejs/](standards/nodejs/) — Node.js-specific structure and patterns
- [standards/service-design/](standards/service-design/) — layered composition, handler responsibilities, dependency patterns

**Database & persistence:**

- [standards/database/](standards/database/) — SQL, NoSQL, and schema design patterns

**Patterns library:**

- [standards/patterns/](standards/patterns/) — cross-cutting concerns (evidence discipline, non-findings, parallel work)

### Quality & Safety

- [reviews/project-review-checklist.md](reviews/project-review-checklist.md) — HAWP maintenance quality gate
- [reviews/public-safety-checklist.md](reviews/public-safety-checklist.md) — safety gate for public distribution
- [reviews/publication-safety-guidelines.md](reviews/publication-safety-guidelines.md) — privacy and generic-ness guidance

---

**Note:** Templates and patterns are optional. Better shaping comes from clearer authoring, not more fields. The core shape is locked; improvements flow through guidance, examples, and discipline.
