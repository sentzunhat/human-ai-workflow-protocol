# 2026-09-29 — Mobile coding-agent workflow direction

Status: archived public-safe status checkpoint
Conversation/event date: 2026-09-29
Archive date: 2026-09-29
Repository: `sentzunhat/human-ai-workflow-protocol`
Branch: `checkpoint/2026-09-29-mobile-codex-workflow`
Path: `.hawp/work/status/2026/09/29/mobile-codex-workflow/status.md`

This checkpoint preserves a public-safe workflow requirement discovered during a discussion about open-source coding-agent alternatives that can use local Ollama-backed models. It avoids private project, business, and personal context.

## What changed

### Before

- Continue was already a strong VS Code extension for local model coding workflows, especially with Llama/Ollama-backed models.
- Aider, Cline, and Roo Code were known as possible coding-agent options with different strengths around terminal-driven changes, workspace edits, and tool use.
- Existing HAWP direction already included provider-independent handoff concepts, local-model support, provider guidance, and durable state outside a single agent conversation.
- The unresolved gap was not basic IDE integration. The missing capability was a mobile-friendly way to start and review work that runs in a trusted development environment.

### During

- The discussion narrowed the requirement from "open-source Codex alternative with Ollama" to a more specific mobile-to-development-host workflow:
  - start or guide coding work from a phone,
  - run the actual workspace operations in a trusted user-controlled development environment,
  - support repository edits, command results, and Git history review,
  - use local/open models where practical through Ollama,
  - and preserve durable session evidence and resume state.
- Continue remained valid for desktop/IDE local-model coding, but it does not fully satisfy the mobile control-plane requirement by itself.
- A plausible architecture emerged: secure access to a trusted development host plus a thin mobile-accessible orchestration surface above HAWP-style work state.
- No implementation was completed in this conversation. No repository code, CLI adapter, mobile UI, webhook, or remote-execution service was built.

### After

- Current direction: treat mobile coding-agent control as a workflow/orchestration problem, not as a replacement for Continue.
- Continue remains the preferred desktop local-model coding interface.
- Aider remains useful for deliberate terminal/Git patch workflows.
- HAWP should continue to model durable task state, handoff, evidence, review, and continuation separately from whichever coding agent or editor performs the actual edit.
- The next architectural question is whether to define a HAWP-compatible remote participant/control-plane pattern for mobile-initiated work.

## Completed work

- Clarified that the desired alternative is not just an IDE extension.
- Identified the core requirement: phone-initiated coordination of a trusted development environment that can produce inspectable repository changes and Git state.
- Preserved the distinction between local desktop coding tools and a mobile-friendly session/control-plane.
- Archived this direction in a public-safe repository checkpoint without including private surrounding context.

## Decisions

- Continue + Ollama remains the best desktop/local IDE path, but it is not sufficient for the mobile workflow by itself.
- The mobile workflow should not depend on one vendor-specific model or editor.
- The durable state layer should remain provider-neutral: work packet, command/output evidence, file/Git changes, review checkpoints, and resume state should be separable from the UI and model provider.
- A future solution should prefer a controlled, user-owned development environment with explicit review gates rather than broad uncontrolled access.

## Experiments / hypotheses

- Hypothesis: HAWP can serve as the durable state layer for a mobile-initiated coding-agent workflow, while Continue/Aider/Cline-like participants perform code changes inside a trusted development workspace.
- Hypothesis: the first useful prototype can be CLI-first before any native mobile UI is built.
- Hypothesis: a phone-friendly web or chat interface could submit work packets to a dev host, while the host records evidence and prepares diffs for review.
- Hypothesis: the existing CLI participant-loop direction and parked participant-adapter work may overlap with this requirement, but the mobile control-plane adds access, security, authorization, and UX concerns.

## Unresolved questions

- What is the minimum safe mobile control surface: chat, web UI, GitHub issue/PR comments, webhook, or dedicated app?
- Where should execution run: personal computer, home server, dev container, cloud VM, or CI-like ephemeral runner?
- How should authentication, authorization, audit logs, and command allowlists be handled?
- Which operations require explicit approval before execution or Git commit?
- Should Git changes be committed automatically, staged for review, or opened as a branch/PR?
- How should local Ollama models be made available without exposing the host broadly?
- How does this relate to the previously identified CLI-driven iterative participant loop direction?

## Strategic impact

- This reframes the Codex alternative search into a broader HAWP-adjacent workflow opportunity: a secure mobile-initiated coding session that can continue work on a real repository without being tied to one cloud vendor.
- The likely architecture should layer:
  1. a mobile-friendly request/control interface,
  2. a secure access boundary,
  3. a trusted dev-host runner,
  4. one or more participant adapters such as Continue/Aider/Cline-style tools,
  5. a durable HAWP work/evidence/resume record,
  6. and explicit review/approval gates for risky operations.
- This should be treated as an architecture proposal before implementation. The safety, repository-integrity, and Git-history implications require deliberate design.

## Continuation state

Current state:

- No implementation exists for the mobile coding-agent workflow.
- The requirement is now clarified as mobile-initiated coordination over a trusted repository-capable development environment.
- Continue remains part of the local/desktop toolchain rather than the whole solution.
- HAWP's provider-neutral handoff/evidence/review concepts appear relevant to the durable state layer.

Next milestone:

- Write a small architecture/work-intake proposal for a mobile-initiated coding-agent control plane.

Next actions:

1. Review the existing CLI-driven participant-loop checkpoint and any parked CLI participant-adapter work.
2. Define the minimum participant contract for a remote coding session: input packet, allowed tools, execution result, diff summary, command log, Git status, and resume token.
3. Define the security model: host location, network boundary, auth, allowlist, approval gates, audit log, and secret handling.
4. Prototype the smallest CLI-first flow before building a mobile UI.
5. Validate with one low-risk repository task: create a branch, make a small edit, show diff, ask for approval, then commit.

Blockers:

- No target architecture has been approved.
- No mobile control surface has been selected.
- The execution/security boundary is unresolved.
- The relationship to existing HAWP CLI participant-loop work needs review.

Dependencies:

- Existing HAWP work/handoff/evidence model.
- Local or remote Ollama runtime availability.
- A secure remote access path such as SSH, VPN, or equivalent controlled network boundary.
- A coding participant adapter capable of editing files and producing inspectable diffs.

## Memory delta

- Add compact durable context only: the user wants a Codex-like workflow that can be initiated from a phone and coordinate work in a trusted computer/remote dev environment, including repository changes, command evidence, Git state, and local/Ollama-compatible agents.
- Preserve that Continue remains good for local VS Code/Ollama use, but the new requirement is a secure mobile-to-dev-host control plane.
- Do not duplicate broader GPU/AI infrastructure, HAWP v0.0.24, CLI loop, or existing Continue/Ollama history.

## GitHub archival

- Related repository found: Yes
- Organization: `sentzunhat`
- Repository: `sentzunhat/human-ai-workflow-protocol`
- Source branch inspected: `feature/v0.0.24`
- Archival branch: `checkpoint/2026-09-29-mobile-codex-workflow`
- Checkpoint file: `.hawp/work/status/2026/09/29/mobile-codex-workflow/status.md`
- Fallback archive: No
- Note: archival was isolated on a checkpoint branch so the active v0.0.24 feature branch was not modified directly.

## Resume summary

**Resume from:** Treat the mobile coding-agent workflow as an unimplemented HAWP-adjacent architecture proposal, with Continue/Aider-style tools as possible participants rather than the whole solution.

**Next objective:** Draft a focused architecture/work-intake proposal for a secure phone-to-dev-host coding-agent control plane with local/Ollama-compatible execution, Git diff review, and explicit approval gates.
