# ADR: Keep HAWP Generic, Portable, and Privacy-Safe for Public Use

Date: 2026-04-27
Status: Proposed
Project: Human AI Workflow Protocol

## Context

HAWP is intended to be reusable across repositories, teams, and AI-assisted workflows.
For that to hold in public, the repository must read like a neutral protocol toolkit rather than a record of one person's projects, plans, or private operating context.

The main risk is not protocol complexity.
The main risk is leakage or unnecessary anchoring through examples, prompts, status artifacts, install docs, or repo-local language that only makes sense inside the original author's work.

## Decision

Treat the public HAWP repository as a generic core.
Keep private or repo-specific usage in private overlays, private notes, or downstream repo-local implementations.

Public HAWP core:

- protocol docs
- reusable templates
- fictional or clearly generic examples
- prompts and review aids that do not depend on private context
- checklists for publication safety and maintenance

Private or repo-local overlays:

- real product audits
- business strategy or roadmap details
- internal status reports
- personal workflow notes
- client-specific materials
- logs, secrets, or environment-specific paths

## Rationale

This keeps HAWP teachable without exposing private context.
It also makes adoption easier because external users do not need to understand the original author's projects to understand the protocol.

## Review Rules

When reviewing content for the public core:

1. If content teaches HAWP generically, keep it.
2. If content is useful but too specific, rewrite it as a fictional or generic example.
3. If content only describes private context, remove it from the public core or move it to a private overlay.

## Audit Focus

Review these surfaces before publishing or broad sharing:

- README files
- examples
- templates
- prompts
- status artifacts intended for public inclusion
- install and update docs
- file names and metadata

Look for:

- personal names
- private project or product names
- company, client, or customer references
- internal strategy or launch details
- local file paths or infrastructure specifics not needed for teaching
- credentials, tokens, emails, phone numbers, or logs
- personal context unrelated to the protocol

## Consequences

Advantages:

- stronger privacy protection
- clearer open-source positioning
- easier adoption across projects
- lower long-term cleanup risk

Tradeoffs:

- examples may need periodic rewriting
- some real-world detail will stay outside the public repo
- maintainers must enforce the boundary over time

## Implementation Notes

- Keep public examples fictional or clearly generic.
- Use the publication-safety review docs before release.
- Prefer small, targeted rewrites over broad abstraction that makes examples unusable.
- Do not expand the HAWP schema to solve publication-safety concerns.
