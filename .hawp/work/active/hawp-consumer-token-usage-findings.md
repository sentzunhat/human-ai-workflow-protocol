# Consumer findings: Cursor token usage and HAWP install-time context

**Audience:** HAWP maintainers (this repo), not a Cursor billing memo.
**From:** a consumer infrastructure repository (HAWP installed).
**Type:** task (feedback / investigation record)
**Reported:** 2026-08-21
**Status:** in-progress (capture complete; no kit or rule changes requested as the "fix")
**Risk:** n/a (context only)

---

## Input

> Capture 30-day Cursor usage findings for HAWP upstream. Do not change kit or rules. Give HAWP maintainers evidence on whether install-time rules cause the harness to gather more context.

## Mission

Hand HAWP a compact, evidence-split report from a consumer infrastructure repository and one operator’s Cursor export so they can review always-on vs on-demand kit loading. This item is feedback. It is not a request to patch the consumer.

## Constraints

- Do not modify `.hawp/kit/`, provider overlays, or generated alwaysApply rules as a "fix" for this consumer bill.
- Do not copy usage export files into HAWP work.
- Separate direct evidence from inference.
- Do not treat HAWP as the whole bill.

## Output

This file. Follow-on review items in this repo: `hawp-upstream-alwaysapply-imperative-reads.md` and `hawp-upstream-always-on-token-budget.md`.

---

## Usage snapshot (consumer, 30 days)

Window: 2026-07-23 through 2026-08-21. Source: one operator’s Cursor export for a consumer infrastructure repository.

| Measure | Value |
| --- | --- |
| On-demand spend | $1,018.90 |
| Events | 2,209 |
| Tokens | 1.356B |
| Cache reads | 87.5% (1.187B cache / 169.7M fresh) |

HAWP was installed in that consumer during the window (2026-08). Spend continued after install while the same people ran infrastructure work in long Cursor Agent chats.

**Question HAWP should care about:** do install-time always-on rules cause the harness to gather more context than a skinny consumer install needs?

**Short answer:** bucket A (always-injected HAWP text) does not explain the bill. Bucket B (obedient kit reads) is a real but bounded slice. Bucket C (install-session history plus later long chats) is the HAWP-shaped cost that actually matters.

---

## Evidence A / B / C

### A) Auto-injected on every Cursor request (Confirmed)

These files are `alwaysApply: true` in the consumer:

- `.cursor/rules/hawp-core.mdc` (~230 tokens; measured ~232 from file size)
- `.cursor/rules/hawp-backlog-alignment.mdc` (~235 tokens)
- `CLAUDE.md` from the HAWP install PR (~305 tokens)

Together that is about **800 HAWP tokens** on every request, plus whatever else Cursor injects that is not HAWP.

**Confirmed:** Cursor does **not** auto-inject `.hawp/kit/start-here.md`, `.hawp/work/BACKLOG.md`, or active plans. Those are not `alwaysApply` files.

Also `alwaysApply: false` with globs (not on every request; only when matching files are in play):

- `.cursor/rules/hawp-intake.mdc` (`**/.hawp/**`)
- `.cursor/rules/hawp-docs-alignment.mdc` (`.hawp/kit/**`, `.hawp/work/**`)

### B) Instruction-following kit reads (Confirmed text; Likely cost)

Consumer `hawp-core.mdc` is imperative. It tells the agent to read `.hawp/kit/start-here.md` for task shaping, `.hawp/kit/usage/status-report.md` for handoffs, and `.hawp/kit/usage/workflow-loop.md` for multi-iteration work, and to track `.hawp/work/BACKLOG.md` / `active/` / `closed/`.

`CLAUDE.md` repeats the start-here and status-report reads.

Typical obedient first turn: on the order of **10k to 25k tokens** of kit/work files, then those reads cache in that thread.

Whole kit on disk in that consumer: about **419 KB** (~105k tokens) if an agent reads all of it. `start-here.md` lists a large tree (templates, standards, reviews, usage guides). A careless agent can cascade through that tree. That cascade is not auto-injection; it is following links.

### C) Session history (Confirmed size; Likely dominant HAWP-shaped cost)

One HAWP-install transcript ran **2026-07-28 through 2026-08-11**: about **6.1 MB**, **1,024 user turns**. That history sits in the conversation cache. It is not the 800-token alwaysApply payload.

A later coordinator-lineage chat from 2026-08-12 is about 1.7 MB. Median cache read in the export is ~210k through early August, then ~788k in mid-August. **A cannot explain that jump.** B is a slice of first-turn and occasional re-reads. C plus long chats plus parallel in-IDE agents plus high thinking is the bill.

---

## Cost drivers (ranked)

Compared Aug 3-6 vs Aug 17-20 (same consumer export, not a HAWP-only split):

1. **Context size (~69% of the intensity jump).** Tokens/request 495k to 1.25M. Median cache ~210k to ~788k. Request volume fell (757 to 169). Intensity rose.
2. **High-thinking mix (~31%).** High-thinking share 1.5% to ~98%. Claude 4.6 medium and high share the same list rates ($3 input / $3.75 cache write / $0.30 cache read / $15 output per 1M). High effort writes more cache and bills more thinking output.

HAWP did not set thinking level. Users and the IDE did. HAWP can still reduce how much *stable* text lands in every new thread (A+B) and how install chats are run (C).

---

## What we are not asking HAWP to do

- Do not rewrite the kit or rules in the consumer as a "fix."
- Do not treat this as a demand to remove HAWP.
- Do not re-run the 30-day usage audit from this file. Numbers above are from the already completed audit.

This item is **context and feedback** for the HAWP project to review.

---

## Ideas for HAWP-side review (not consumer patches)

1. **Always-on rules should not imperatively say "read start-here on every task."** Point at start-here as the on-ramp when the agent is shaping work, not as a mandatory first tool call. Likely lives in generated `core/providers/shared/behaviors` that become consumer `hawp-core.mdc`.
2. **Make kit reads clearly on-demand.** Distinguish "injected every request" from "open if you are doing X."
3. **Warn that long install chats accumulate B-reads and then become C.** Install sessions that last two weeks will dwarf A forever after.
4. **Document an expected always-on token budget** for consumer installs (this consumer's A is ~800 HAWP tokens; say what you intend).
5. **Optional skinny vs full kit** for consumer installs: protocol + backlog loop vs the full standards/templates tree.
6. **Avoid listing the entire kit tree in start-here as something agents will open.** A link index is useful for humans; models treat it as a reading list.

Follow-ons in this repo: `hawp-upstream-alwaysapply-imperative-reads.md`, `hawp-upstream-always-on-token-budget.md`.

---

## Direct evidence vs inference

### Direct evidence

- Export-backed totals in the consumer (not copied here): $1,018.90, 2,209 events, 1.356B tokens, 87.5% cache.
- `alwaysApply: true` on consumer `.cursor/rules/hawp-core.mdc` and `.cursor/rules/hawp-backlog-alignment.mdc`; `CLAUDE.md` present from HAWP install.
- `hawp-core.mdc` text: "Read `.hawp/kit/start-here.md` ... `status-report.md` ... `workflow-loop.md`."
- Kit tree size ~419 KB on disk in that consumer.

### Inference (not proven from the export alone)

- That ~800 always-on HAWP tokens are a rounding error next to 210k-788k median cache.
- That a typical obedient first turn is 10k-25k (order of magnitude from file sizes + common agent behavior, not a per-event export field labeled "HAWP").
- That the 6.1 MB install transcript is the main *HAWP-shaped* contributor to later cache, vs other long infra chats in the same month.
- That changing alwaysApply wording would move the monthly bill a lot. Unclear. It would shrink B on *new* threads. It would not shrink C on threads that already exist.

### Unknowns

- Cursor's exact injection packing besides the HAWP files named above.
- How much of median cache is HAWP vs infrastructure plans, ticketing, and parallel Agent sessions.
- Whether other HAWP consumers see the same install-chat pattern.

---

## Non-findings

- No evidence that HAWP auto-injects the full kit on every request.
- No evidence that optional high-cost Cursor runtimes drove this 30-day on-demand total.
- This consumer is not asking to compact or rewrite `.hawp/kit/` locally.

---

## Sources (do not duplicate secrets)

- One operator’s Cursor usage export for a consumer infrastructure repository (gitignored in the consumer; not copied here)
- Consumer git: `.cursor/rules/hawp-core.mdc`, `.cursor/rules/hawp-backlog-alignment.mdc`, `CLAUDE.md` (from the HAWP install)
- Local analysis canvas in the consumer project (not in this repo; not required to review this item)

## Next step

Review A/B/C and the idea list. Do not treat filing this item as approval to change kit or generated provider overlays until a maintainer explicitly scopes that work.
