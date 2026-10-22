# Unify HAWP's pentagon brand mark across the website

**UUID:** `96c45719-8758-4bd2-a210-a26ffebab262`
**Type:** improvement
**Reported:** 2026-10-13
**Status:** `parked`
**Risk Level:** low
**Target Version:** website launch

## Investigation

The header currently renders an outlined rotated square and a smaller rotated square. The hero's central packet repeats the same diamond motif through CSS pseudo-elements. The SVG social card also uses diamonds, while the favicon uses a circular mark. A shared five-sided mark should replace these mixed geometries with a regular, upright pentagon while keeping the existing HAWP palette and wordmark.

**Directly inspected:** `website/src/lib/components/Header.svelte`, `website/src/lib/components/Hero.svelte`, `website/src/app.css`, `website/static/favicon.svg`, and `website/static/social-card.svg`.

## Mission

Use a consistent, optically balanced five-sided HAWP mark in the site header, hero center, favicon, and social card. User review of the first implementation showed that the nested filled pentagons and surrounding squares were too heavy. The mark should be a single clear pentagon, and the hero hub should itself be pentagonal.

## Plan

1. Establish one simple regular pentagon outline in the website palette.
2. Remove the header's surrounding square and make the hero packet hub a pentagon with the WORK label inside.
3. Align the transparent SVG favicon and social card with the same silhouette; regenerate the social-card PNG from its SVG source.
4. Inspect the result at compact header/favicon and larger hero/social-card sizes; refine proportions if needed.

## Options

- **Recommended:** one SVG path for the brand mark, with CSS/theme colors supplied where the component is live. This gives the same five-sided geometry at every size without a nested icon or square frame.
- Keep separate CSS and artwork approximations. This avoids shared geometry but risks the current mismatch returning.

## Work Coordination

**Owner:** Codex
**Implementation status:** in-progress
**Expected overlapping / touched paths:**

- `website/src/lib/components/Header.svelte`
- `website/src/lib/components/Hero.svelte`
- `website/src/app.css`
- `website/static/favicon.svg`
- `website/static/social-card.svg`

**Parallel work risk:** low
**Can implement now:** yes; item 1607037f follows this mark pass and may then touch `Hero.svelte` and `app.css` for cable animation.

## Risk

Low. This is a visual identity refinement confined to existing website artwork and styles.

## Verification

- Source inspection confirms one five-sided path in `website/src/lib/components/Header.svelte`, `website/src/lib/components/Hero.svelte`, `website/static/favicon.svg`, and `website/static/social-card.svg`. The social-card PNG was regenerated from the SVG.
- `npm run check` completed with zero errors and zero warnings using the website's Node 26.10.0 toolchain. `npm run build` completed successfully after the final header breakpoint adjustment.
- Local browser review at `127.0.0.1:5174` showed the header and hero on desktop and at a 390-pixel viewport in dark and light themes. The 775-pixel preview exposed cramped navigation; moving the mobile-menu breakpoint to `lg` resolved it in the rendered preview.
- Local desktop preview now shows the header and hero using the same single outline pentagon. The regenerated 1200×630 PNG social card was also inspected and uses that silhouette without the rejected nested mark or square frame.
- `npm run check` and `npm run build` passed on Node 26.10.0 after the current local edits.
- Owner aesthetic acceptance remains open before this item can close.

## Close Checklist

- [ ] Outcome and verification recorded.
- [ ] Plan moved to dated closed archive.
- [ ] Backlog updated.


## Parking / Resume

- **Reason:** Temporarily paused while c147f2ae is prioritized for the HAWP tool-first usage audit. Resume after its plan is reviewed; website work returns first in dependency order: 96c45719, then 1607037f, with b4300599 also returning to active work.
- **Resume priority:** Resume after `c147f2ae` is reviewed and ready for execution. Website sequencing: `96c45719` before `1607037f`; `b4300599` remains a website priority as well.
