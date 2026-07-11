# No-HAWP arm — reused same-state run (not re-run)

The Layer-2 change touched only the **HAWP authoring pattern** (the new surface-sweep habit in `authoring-patterns.md`) and regenerated dot-folder artifacts. The reviewed subject — `librarian/` — is byte-for-byte unchanged from the Layer-1 run, so the no-HAWP arm is carried over rather than re-executed:

- Output: [`../../2026-06-15-layer1-coverage-balance/no-hawp/output.md`](../../2026-06-15-layer1-coverage-balance/no-hawp/output.md)
- Scorecard: [`../../2026-06-15-layer1-coverage-balance/no-hawp/scorecard.md`](../../2026-06-15-layer1-coverage-balance/no-hawp/scorecard.md)
- Score: **45 / 60 → 75%**

This is a legitimate same-state baseline: that no-HAWP arm was run in a clean `/tmp/` workspace against the current Node 26 `librarian/` tree, which the Layer-2 work did not alter. The HAWP authoring change cannot affect the unguided arm by construction.
