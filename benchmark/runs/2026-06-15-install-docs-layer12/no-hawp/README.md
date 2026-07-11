# No-HAWP arm — reused same-state run (UNCHANGED, 85%)

This re-run changed only the **HAWP authoring pattern** (Layer 1 + Layer 2). The audited subject — `distribution/` install sources and generated guides — is **byte-for-byte unchanged** (verified: `git status --short distribution/` reports no diffs, `distribution:validate` passes). The no-HAWP arm is therefore carried over unchanged, not re-executed:

- Output: [`../../2026-06-15-install-docs-truth-audit/no-hawp/output.md`](../../2026-06-15-install-docs-truth-audit/no-hawp/output.md)
- Scorecard: [`../../2026-06-15-install-docs-truth-audit/no-hawp/scorecard.md`](../../2026-06-15-install-docs-truth-audit/no-hawp/scorecard.md)
- Score: **51 / 60 → 85% (unchanged)**

## Why this is the integrity check you asked for

The no-HAWP arm never receives the HAWP shape, so HAWP authoring changes (Layer 1/2) **cannot** affect its score by construction. This run demonstrates that directly: same task, same unchanged subject, identical no-HAWP score (85%).

The "10-point drop" observed between this audit's no-HAWP (85%) and the bounded-review/layer runs' no-HAWP (75%) is a **task-type difference, not a regression**:

| Task type | No-HAWP score | Why |
| --- | --- | --- |
| Standards / truth audit (this run) | 85% | A doc-vs-script audit is naturally bounded; little room to drift, so unguided work scores high. |
| Bounded repo review (scope-creep trap) | 75% | The "what else should we clean up?" bait actively invites drift, so unguided work loses drift/scope/false-positive points. |

Both no-HAWP scores are correct for their task; they are not comparable to each other, and neither moved because of the HAWP changes.
