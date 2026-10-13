# Give the homepage protocol cables electric, physics-led motion

**UUID:** `1607037f-4e55-46f5-91da-6405a6b08f17`
**Type:** improvement
**Reported:** 2026-10-13
**Status:** `in-progress`
**Risk Level:** medium
**Target Version:** website launch

## Investigation

`website/src/lib/components/Hero.svelte` already lays out five draggable fields and draws one SVG cubic path from the center to each field. `website/src/app.css` animates a dashed connector's opacity and dash offset, while the nodes float with CSS. This is a useful foundation, but the wires currently read as static decorative curves rather than responsive cables. Theme colors and `prefers-reduced-motion` support already exist.

**Directly inspected:** `website/src/lib/components/Hero.svelte`, `website/src/app.css`, and `website/src/lib/theme.ts`.

## Mission

Make the five homepage connections feel like softly energized cables: a restrained themed glow, visible energy traveling along each line, and elastic wobble when a field is dragged or released. Preserve field dragging, responsive layout, light/dark themes, and reduced-motion behavior.

## Plan

1. Build on the existing SVG paths and theme tokens; do not add a physics dependency or canvas renderer.
2. Add layered cable rendering with a subtle glow and a controlled energy pulse in both themes.
3. Add lightweight spring-like path response so dragged nodes tug their cables and the lines settle after release.
4. Keep animation bounded and respect reduced-motion preferences; avoid continuous layout reads and unnecessary state churn.
5. Review pointer, touch, keyboard-focus, small-screen, and reduced-motion behavior.

## Options

- **Recommended:** SVG path layers plus a small bounded spring state in the existing Svelte component. This keeps the visual crisp, interactive, and dependency-free.
- Use CSS-only pulse and wobble. This is simpler but cannot respond naturally to pointer movement.
- Add a general physics package or canvas renderer. This is unnecessary for five cables and adds bundle and maintenance cost.

## Work Coordination

**Owner:** Codex
**Implementation status:** in-progress
**Dependency:** 96c45719 establishes the pentagon mark first; this item then updates the shared hero artwork and cable styling.
**Expected overlapping / touched paths:**

- `website/src/lib/components/Hero.svelte`
- `website/src/app.css`

**Parallel work risk:** medium
**Can implement now:** yes, after 96c45719's edits to the overlapping hero/style files are complete.

## Risk

Medium because this changes an interactive hero animation. Motion must remain subtle, accessible, and performant across pointer and touch input.

## Verification

- Confirm cables visibly respond to dragging and settle without persistent jitter.
- Confirm theme colors work in light and dark mode and reduced-motion disables continuous movement.
- Review mouse and touch interactions at desktop and mobile widths.
- Added arrow-key movement for each focused field; Shift+Arrow moves in larger increments. Pointer drag was exercised in a local Chrome preview and visibly moved the node and its cable.
- Local desktop preview shows the pentagonal hub and themed cable glow in the dark theme. The iPhone 16 Pro Max emulation at 440 pixels shows the compact header and stacked hero; the initial consent notice overlays part of the hero artwork until a choice is made.
- The `prefers-reduced-motion` CSS disables the node float and cable pulse; the spring loop also clears wobble and velocity when reduced motion is active.
- `npm run check` and `npm run build` passed on Node 26.10.0.
- Direct touch-device behavior and owner visual acceptance remain open.

## Close Checklist

- [ ] Outcome and verification recorded.
- [ ] Plan moved to dated closed archive.
- [ ] Backlog updated.
