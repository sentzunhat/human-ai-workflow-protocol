# Create a hand-drawn HAWP LinkedIn Featured card

**UUID:** `098c8a5c-a185-4f29-9bce-56aee7b9e39a`  
**Type:** improvement  
**Reported:** 2026-10-21  
**Risk:** medium  
**Status:** in-progress

## Input

> Make the HAWP LinkedIn Featured visual feel human and authored: build an original SVG drawing and a small deterministic repo-native script; do not use generated bitmaps, publish, or upload.

## Investigation

- The public HAWP README describes HAWP as an open workflow layer for humans and AI agents. The core public claim is about shaping and carrying intent across agents, sessions, and handoffs; HAWP is explicitly not a hosted platform, runtime, or orchestration engine.
- The current public website social-card source is `website/static/social-card.svg`, a hand-authored 1200×630 SVG using a dark background, warm sand/coral accents, a restrained grid, and a single pentagon brand mark. It is a visual reference only; this task does not edit or regenerate it.
- Existing parked item `96c45719` owns the HAWP website pentagon mark and social-card consistency. This task is separate: it creates an additional private LinkedIn Featured concept using the established visual vocabulary, without modifying the site or existing social card.
- `unDraw Handcrafts` offers downloadable customizable line-art doodles and states a free commercial/personal-use license, while prohibiting artwork packs/resale, competing integrations, and use of its assets/site to train automated systems. This concept uses original SVG geometry and does not copy, scrape, or bundle unDraw artwork. Source: https://handcrafts.undraw.co/license
- The owner's supplied LinkedIn screenshot showed a HAWP Featured card with a generic placeholder at capture time; current card state and LinkedIn's image requirements remain unverified.

## Goal

Create a human-looking, original vector illustration draft for HAWP's LinkedIn Featured link, generated reproducibly by a small script in this repository. Keep the source and output private in this work item until reviewed.

## Plan

1. Use the HAWP public README and `website/static/social-card.svg` only for supported wording and existing palette/brand reference.
2. Draw the concept directly in SVG with hand-drawn-style, slightly irregular strokes and simple shapes. Show a human-authored intent card moving through a compact HAWP structure to work artifacts and a handoff. Avoid robot/brain imagery and any implication that HAWP executes or orchestrates agents.
3. Add a small dependency-free Node script that writes the SVG deterministically. Keep the default output under this plan's private asset folder rather than `website/static/`.
4. Verify the generated SVG parses and inspect its rendered appearance at preview size. Keep the SVG source, generator, and visual iteration easy to edit.
5. Check current LinkedIn Featured image requirements and crop behavior before making any final export or upload recommendation. Founder review is required before any public use.

## Implementation scope

- `scripts/artwork/generate-hawp-featured-card.mjs` — dependency-free SVG generator.
- `.hawp/work/active/098c8a5c/hawp-featured-card.svg` — private generated source artwork for review.
- `.hawp/work/active/098c8a5c/plan.md` — plan and evidence.

No site files, package scripts, social profiles, accounts, or public assets are changed by this item.

## Acceptance criteria

- [x] SVG is generated deterministically by the checked-in script with no external runtime dependencies.
- [x] Visual uses HAWP's established public palette and custom drawing. No AI bitmap or third-party artwork is included; founder review remains required before public use of the HAWP mark.
- [ ] Text and imagery accurately describe HAWP as an open workflow protocol and remain readable at card size.
- [x] SVG source is valid and rendered at 1200×630 for visual review, matching the existing HAWP social-card source size; this does not assert a LinkedIn requirement.
- [ ] Verify LinkedIn-specific export dimensions and crop behavior from current authoritative guidance before final export or upload.
- [ ] Founder reviews the concept before it is treated as a final brand asset.
- [ ] No public site deployment, LinkedIn page/profile edit, image upload, or post publication occurs.

## Dependencies and decisions

- Current LinkedIn Featured link, image requirements, and crop behavior need an owner-authorized live check.
- HAWP public website branding provides a palette reference, but use of the official mark in a new public surface still requires founder review.
- Current scope is local drafting only. Public use needs separate owner approval.

## Verification

- [x] Inspected the public HAWP README, existing `website/static/social-card.svg`, and parked brand item `96c45719`.
- [x] Reviewed the linked unDraw Handcrafts page and license; no library artwork was copied or bundled.
- [x] Ran `node scripts/artwork/generate-hawp-featured-card.mjs`; it deterministically wrote `.hawp/work/active/098c8a5c/hawp-featured-card.svg`.
- [x] Parsed the generated SVG as XML and rendered it at 1200×630 for visual inspection.
- [ ] Preview and review at intended LinkedIn crop after current requirements are checked.

## Owner approval

The user's request authorizes creating this work item and a private SVG/script draft. Founder approval is still required before adding the artwork to the website, using it on LinkedIn, or treating it as final brand art.

## Next step

Implement the SVG generator and private draft, then present the concept for review. Do not publish or upload.
