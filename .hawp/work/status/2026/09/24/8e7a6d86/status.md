---
uuid: 8e7a6d86-0529-48af-9479-16ecee516f82
title: "Website discoverability and post-merge actions"
type: status
date: 2026-09-24
---

### Status Report

#### Intent

Continue the `feature/website` lane by documenting repository and domain
discoverability actions that require the website PR to merge, while improving
the static site's SEO and AI-retrieval signals within the branch.

#### Current State

The website branch now contains a post-merge discoverability checklist, a
machine-readable `llms.txt` summary, expanded JSON-LD, and cleaned website
README formatting. The separate PR #41 hardening lane was not changed.

#### What Was Inspected

- `website/README.md`
- `website/src/app.html`
- `website/src/routes/+page.svelte`
- `website/src/lib/content/site.ts`
- `website/static/robots.txt`
- `website/static/sitemap.xml`
- `website/docs/discoverability-baseline.md`
- `.github/workflows/website.yml`
- live branch and worktree status

#### What Changed

- Added `website/docs/post-merge-discoverability-actions.md` covering the
  proposed GitHub description, homepage, topics, Pages/custom-domain steps,
  Community Standards items, search-engine submission, and +7/+14-day evidence.
- Added `website/static/llms.txt` with canonical links, plain-language project
  scope, supported concepts, non-goals, and retrieval terms.
- Expanded static JSON-LD in `website/src/app.html` with publisher, keywords,
  `WebPage`, and `FAQPage` data.
- Documented the new discovery surfaces in `website/README.md` and fixed its
  literal escaped-newline section break.
- Regenerated the tracked static build output through the normal website build.

#### What Was Directly Verified

- `mise exec node@26.5.0 -- npm run check` passed with 0 errors and 0 warnings.
- `mise exec node@26.5.0 -- npm run build` passed and wrote `website/build`.
- `website/build/llms.txt` exists.
- Generated `website/build/index.html` contains one each of `FAQPage`,
  `WebPage`, and `SoftwareSourceCode`, plus the website graph entries.
- Generated output contains no `@html` template residue.
- `git diff --check` passed.

#### What Remains Unproven

- The PR has not been merged and GitHub repository settings have not been
  changed.
- GitHub Pages deployment, custom-domain validation, DNS, and enforced HTTPS
  require post-merge external verification.
- Search indexing, rankings, referrals, AI answer inclusion, and traffic lift
  require live post-launch measurements.
- No claim is made that the proposed GitHub description or topics are applied.

#### Constraints

- Preserve the separate PR #41 hardening lane and existing website artifacts.
- Do not change GitHub settings, DNS, search-console accounts, or merge/push
  state as part of this turn.
- Treat user-provided GitHub screenshots as baseline evidence only.

#### Help Wanted

After merge, verify the final GitHub metadata, Pages deployment, custom domain,
HTTPS, and the +7/+14-day discoverability snapshots against the checklist.

#### Suggested Next Step

Review the generated website preview on desktop and mobile, then use the
post-merge checklist once the PR reaches `main`.

### Interaction enhancement follow-up

The hero's five-field work packet was extended with a dependency-free
pointer interaction. Nodes are semantic buttons that can be dragged with a
mouse or touch pointer; the SVG connectors recalculate as nodes move. CSS
floating motion, hover lift, theme-aware shadows, elastic Bézier curves, and
reduced-motion fallbacks reuse the existing theme tokens rather than adding a
visual library.

The local preview was visually inspected in dark/system and light themes. A
live drag of the Input node moved the node and its connector together. The
remaining visual QA boundary is wider responsive/mobile inspection and a
production deployment check.

### Lighthouse follow-up

The supplied Lighthouse run reported 79 Performance, 95 Accessibility, 100
Best Practices, 100 SEO, and 2/3 Agentic Browsing. The report's 3,011 KiB
payload and 2,719 KiB minification estimate do not match the current static
production build, which is approximately 240 KiB including assets. The test
URL used port `5174`, so the next comparison should use the production preview
server rather than a Vite development server.

Applied follow-up changes:

- `llms.txt` now uses Markdown hyperlink lists and an `Optional` section in the
  documented specification shape.
- `app.html` now exposes `llms.txt` through a `describedby` link relation.
- The animated protocol node/detail text now uses a dedicated higher-contrast
  token for Light and Dark themes.

The supplied report's remaining contrast finding was not independently
reproduced with a fresh Lighthouse run in this session; the color adjustment
is therefore a targeted remediation, not a claim that the full audit is now
perfect.

The attached Lighthouse JSON later confirmed the exact contrast failure and
also identified a separate `label-content-name-mismatch`: each draggable
button used an action-only `aria-label` such as `Drag Input field` while its
visible text was `Input request`. The implementation now lets the visible
node text provide the accessible name and associates the shared drag
instruction through `aria-describedby`. `svelte-check` remains clean after
this change.

### Second Lighthouse JSON comparison

The second attached run was captured at `2026-09-24T19:20:41Z` against the
same Vite development URL, `http://localhost:5174/#benchmarks`. Compared with
the first attached run:

- Performance improved from `0.66` to `0.81`.
- Accessibility improved from `0.95` to `1.00`.
- Best Practices remained `1.00`.
- SEO remained `1.00`.
- Agentic Browsing remained `0.67`.
- LCP improved from `20.2 s` to `3.5 s`.
- Contrast passed in the second report.

The second report still contains the old `aria-label="Drag Input field"` and
still says `llms.txt` has no links. Those findings predate the current source
fixes, which use `aria-describedby` and Markdown hyperlink lists. A fresh
server restart is required before treating those audits as current.

### Release lookup loading state

The HAWP release-download section now shows a compact shimmer skeleton beside
each platform while the GitHub latest-release request is pending. The static
`currentVersion` and `/releases/latest/download/...` fallback behavior remain
available when the request is slow or unavailable. The loading status is also
announced with `aria-live`, and the shimmer stops under
`prefers-reduced-motion`.

Direct verification after this change:

- `mise exec node@26.5.0 -- npm run check` passed with 0 errors and 0 warnings.
- `mise exec node@26.5.0 -- npm run build` passed.
- `git diff --check` passed.

### Unified dynamic version labels

The page previously updated the release version in the download section but
left the hero badge, agent description, and benchmark label on their initial
static value while the GitHub request was pending. These references now share
the same reactive release state and reusable `ReleaseVersion` component. Every
visible version therefore follows the same sequence: readable `v0.0.0`
fallback, gradient shimmer while the latest-release request is pending, then
the fetched GitHub tag when available. The pending fallback is intentionally
provisional, and a failed request leaves it readable rather than leaving an
empty label.

Direct verification after this refinement:

- `mise exec node@26.5.0 -- npm run check` passed with 0 errors and 0 warnings.
- `mise exec node@26.5.0 -- npm run build` passed.
- `git diff --check` passed.

### Current HAWP Lighthouse JSON

The new attached HAWP report was captured at `2026-09-24T19:57:45Z` against
`http://localhost:5174/#benchmarks`. It reports Performance `0.79`,
Accessibility `1.00`, Best Practices `1.00`, SEO `1.00`, and Agentic Browsing
`1.00`. It also reports zero CLS, 10 ms TBT, and 3.9 s LCP. The network table
shows two successful requests to the same GitHub releases API endpoint, one
from the homepage release-version check and one from the download cards.

The shared `loadLatestRelease()` helper now caches the in-flight promise, so
those callers share one request and one result. The release skeleton remains
the visual pending state while that single request resolves. During that state,
each platform card keeps the readable fallback label `v0.0.24` visible and
layers a small theme-aware gradient sheen over the pill. Once GitHub responds,
the fetched release version replaces `v0.0.24`; the sheen is removed. The text
remains above the sheen for contrast, and the animation is disabled for users
who prefer reduced motion.

Direct verification after the version-label refinement:

- `mise exec node@26.5.0 -- npm run check` passed with 0 errors and 0 warnings.
- `mise exec node@26.5.0 -- npm run build` passed.
- `git diff --check` passed.
