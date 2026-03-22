# Quick Start: Implementation Ordering

## Do This First (Sequential)

```
1. TASK-029: Data Models
   └─ Creates TypeScript types (DetectionReport, BlockedItem, etc)
   └─ Takes ~1 day
   └─ Unblocks everything else
```

## Then Do These In Parallel

```
2a. TASK-027: CLI Entry Point
    └─ Creates ./hawp command + argument parsing
    └─ Takes ~1-1.5 days
    └─ Produces: ./hawp backlog upgrade --help works

2b. TASK-028: Detection + Dry-Run
    └─ Implements scanning, rules, reporting
    └─ Takes ~3-4 days
    └─ Produces: ./hawp backlog upgrade --dry-run works
    └─ (depends on TASK-029 models, can start once TASK-029 done)
```

## Then Integrate

```
3. Verify combination of TASK-027 + TASK-028
   └─ Test all flag combinations
   └─ Confirm no file modifications
   └─ Verify text + JSON output formats
   └─ Takes ~0.5-1 day
```

## Then GATE & PAUSE

```
❌ STOP HERE before starting TASK-030 (apply mode)
   Ensure Slice 1 is stable, tested, verified
   Get approval before moving to destructive operations
```

## What's Gated (Later)

- **TASK-030+**: Apply mode (destructive, writes files, creates backups)
- **TASK-031+**: Validation integration
- **TASK-032+**: Evidence reports with hashes
- **V2+**: AI-assisted synthesis with governance gates

---

## Files to Create (Slice 1 scope)

```
./.hawp/bin/hawp                                    (entry point script)
librarian/scripts/backlog-upgrade/
  ├─ index.ts                                       (main entry)
  ├─ cli.ts                                         (arg parser)
  ├─ models/
  │  ├─ index.ts
  │  ├─ types.ts
  │  ├─ detection-report.ts
  │  ├─ backlog-fix-plan.ts
  │  ├─ blocked-item.ts
  │  └─ evidence-report.ts
  └─ detection/
     ├─ backlog-parser.ts
     ├─ plan-scanner.ts
     ├─ detector.ts
     ├─ rules/
     │  ├─ auto-fix-rules.ts          (A1-A7)
     │  └─ blocked-rules.ts            (B1-B6)
     └─ output/
        ├─ formatters.ts              (text + JSON)
        └─ report-generator.ts
```

---

## Key Dependencies

```
TASK-029 (models)
    ↓
    ├─→ TASK-027 (CLI uses models)
    └─→ TASK-028 (detection uses models)
         ├─ detection/rules use BlockedItem model
         └─ output/formatters render from models
```

**No circular dependencies. Linear critical path: TASK-029 → (TASK-027 + TASK-028) → verify → gate.**

---

## Success for Each Task

**TASK-029:**

- [ ] TypeScript compiles cleanly
- [ ] All models export from index.ts
- [ ] No circular dependencies
- [ ] BlockedItem has rule/confidence/candidates/reason fields

**TASK-027:**

- [ ] `./hawp backlog upgrade --help` outputs correctly
- [ ] `--dry-run` and `--apply` are mutually exclusive
- [ ] `--validate` can combine with both
- [ ] Exit codes correct (0/1/2)

**TASK-028:**

- [ ] Detection completes in < 2s for test backlog
- [ ] All A1-A7 rules trigger on expected inputs
- [ ] All B1-B6 blocks trigger on expected inputs
- [ ] Text and JSON output formats both valid
- [ ] No files modified (dry-run confirmed)

---

## Current Status

✅ Design approved
✅ Work items created (TASK-027, TASK-028, TASK-029)
✅ Implementation plan documented
🟡 **READY TO START: Begin with TASK-029**
