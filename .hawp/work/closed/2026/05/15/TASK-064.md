## Task: Compact BACKLOG.md and archive closed work

**Backlog ID:** TASK-064
**Type:** task
**Reported:** 2026-05-15
**Risk Level:** low
**Status:** inbox

---

### Input (what was reported)

Derived from backlog alignment guardrails: recently closed items are active and should stay capped while detailed history is preserved in closed archives.

---

### Context

Multiple work items were closed in this cycle and the active coordination index should remain compact and decision-useful.

---

### Analysis

**Root cause (or most likely cause):**
BACKLOG.md is functioning as both coordination index and rolling history, which increases noise if not compacted regularly.

**Directly verified:**

- `Recently Closed` is currently capped to 10 items.
- Closed records already exist under `.hawp/work/closed/YYYY/MM/DD/`.
- Validation currently passes.

**Inferred (not yet proven):**

- A periodic compaction pass will keep intake/navigation faster while preserving traceability in archives.

**Scope — what else is affected:**

- `.hawp/work/BACKLOG.md`
- `.hawp/work/closed/**`
- Optional status/evidence references if index cleanup needs cross-links

---

### Work Coordination

**Owner:** agent
**Implementation status:** not-started
**Overlapping files:**

- `.hawp/work/BACKLOG.md`

**Parallel work risk:** low
**Can implement now:** yes

---

### Recommended Fix

1. Keep `Active Work` concise and current.
2. Keep `Recently Closed` bounded (last 5-10 rows or recent window).
3. Ensure older detail remains discoverable in `.hawp/work/closed/YYYY/MM/DD/`.
4. Re-run backlog validation.

---

## Outcome (filled at close)

Backlog compaction pass completed.

Completed outcomes:
- Active lane reduced to only currently open work during execution; stale duplicate plan drift reconciled.
- `Recently Closed` remains compact and bounded to the last 10 items.
- Closed-history discoverability preserved under `.hawp/work/closed/YYYY/MM/DD/`.
- Backlog validation re-run with PASS result.

## Verification (filled at close)

Directly verified:
- Duplicate active/closed drift (`TASK-054`) was reconciled by preserving the newer active content into the closed path and removing the stale active duplicate.
- `BACKLOG.md` remains compact and decision-useful.
- `./.hawp/bin/hawp backlog validate` returns PASS with no FAIL checks.

Confidence: high. Compaction goal met without losing archive traceability.

## Close Checklist

- [x] Active Work kept concise and current
- [x] Recently Closed bounded
- [x] Archive paths retained as source of full history
- [x] Validation re-run and passing

**Status:**

- [x] Plan written
- [x] Approved / awaiting review
- [x] Implemented
- [x] Verified
- [x] Closed
