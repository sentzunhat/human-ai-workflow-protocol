# 2026-09-05 — README positioning and open-source vision checkpoint

**Conversation/event date:** 2026-09-05  
**Archive date:** 2026-10-12  
**Owning work item:** [`e5fca9c7`](../../../closed/2026/09/11/e5fca9c7/plan.md)  
**Checkpoint type:** retrospective archival status  
**Repository state at archive:** `main` at `66bc64be87ca918eace4601dc0a62b7dd6a631c1` before this archival branch

This checkpoint preserves the September 5 README/positioning discussion without
creating a second work item or treating the October archive as a new product
milestone.

## Before

Before the September 5 discussion, the root README was technically useful but
presented HAWP mainly through the problem of context loss and agent drift. It
opened with "Shape the work before execution begins" and described the core as
a five-field protocol.

The repository already contained substantially more than that framing exposed:
a small repo-native protocol, local CLI/search tooling, MCP integration,
provider overlays, work/status/evidence conventions, local model paths, and
benchmark evidence.

Two positioning gaps were identified:

- the README did not lead with the larger open-source idea of project-owned,
  portable intent across AI agents and sessions;
- the README's "five fields" wording did not reflect the canonical v0.1 shape,
  which has five required fields plus optional `checkpoint`.

At the time of review, repository-level GitHub discoverability was also weak:
the repository had no description, homepage, or topics configured.

## During

The September 5 review inspected the README together with the canonical HAWP
shape/specification, start-here guidance, provider layout, CLI capabilities,
benchmarks, backlog, and recent release history.

The resulting product/communication direction was:

- position HAWP as an **open workflow layer for humans and AI agents** rather
  than only as a context-drift remedy;
- make **project-owned intent** the central idea: intent belongs with the
  project, not inside one vendor's chat/session history;
- emphasize portability across providers and the multi-agent development
  model;
- keep the protocol small and distinguish it from the optional tooling around
  it;
- foreground repo-native/local operation, provider independence, measurable
  evidence, and open-source participation;
- restructure the README toward a fast builder path: hero -> short explanation
  -> minimal shape/demo -> install -> architecture/layers -> proof -> vision;
- correct the protocol description to five required fields plus optional
  `checkpoint`.

Two memorable lines were proposed as public positioning:

> The open workflow layer for humans and AI agents.

> Your model can change. Your intent shouldn't.

The recommendation also included adding a concise GitHub repository description
and relevant topics for discoverability.

## After

The positioning direction was later absorbed into the existing v0.0.24 cleanup
lane owned by `e5fca9c7`; that work item closed on 2026-09-11.

The current README on `main` now preserves the main September 5 decisions:

- it opens with "The open workflow layer for humans and AI agents";
- it uses "Shape intent once. Carry it across agents, sessions, and handoffs";
- it includes "Your model can change. Your intent shouldn't";
- it describes the canonical shape as five required fields plus one optional
  checkpoint;
- it adds explicit multi-agent, local-by-design, evidence, and open-source
  vision sections.

This means the September 5 conversation should be treated as the **positioning
decision and design input**, while the later `e5fca9c7` work records the
implementation and verification.

Subsequent HAWP development continued beyond this conversation. As of the
archive point, `main` includes the later local AI search/MCP intake,
work-record, provider-setup, and filesystem-hardening work merged on
2026-10-01. That later engineering work is not retroactively part of the
September 5 event.

## Strategic impact

The README shifted from describing a useful workflow technique to defining a
clearer open-source category:

**HAWP is a portable, project-owned contract for human intent across changing AI
agents.**

That framing aligns the public story with the actual architecture:

1. a deliberately small protocol;
2. repository-owned kit/work artifacts;
3. optional local CLI/search/MCP tooling;
4. provider adapters rather than provider lock-in.

The positioning also establishes a useful boundary for future product decisions:
new tooling should strengthen portability, retrieval, evidence, and continuity
without turning the core protocol into a proprietary runtime or ever-growing
schema.

## Continuation state

### Current state

- README positioning from the September 5 review is implemented.
- The canonical field-count correction is implemented.
- The owning `e5fca9c7` work item is closed.
- The repository has advanced beyond that work through later HAWP releases and
  hardening.

### Next logical milestone

Treat README positioning as stable unless new user/adoption evidence justifies a
change. Future public-facing work should improve discoverability and demonstrate
the workflow with concrete, credible examples rather than repeatedly rewriting
the category statement.

### Next actions

- Add a concise GitHub repository description aligned with the README.
- Add focused repository topics such as AI agents, MCP, developer tools,
  human-AI workflows, and supported coding-agent ecosystems.
- Prefer real examples, benchmark updates, and adoption evidence over additional
  slogan changes.
- Keep README claims synchronized with canonical spec and measured evidence.

### Blockers

No blocker remains for the September 5 README work itself.

Repository metadata is still incomplete at archive time: description, homepage,
and topics remain unset.

### Decisions not to revisit without new evidence

- HAWP remains protocol-first, not an agent runtime or hosted memory platform.
- Intent and durable work context should remain project-owned and portable.
- Provider independence is a core product property.
- The protocol should stay deliberately small; capabilities should grow in the
  surrounding kit/tooling rather than by casually expanding the core shape.

## Resume summary

**Resume from:** the implemented multi-agent/open-workflow README on current
`main`; do not redo the September 5 positioning work.

**Next objective:** improve repository discoverability and adoption proof while
keeping the current public category statement stable.
