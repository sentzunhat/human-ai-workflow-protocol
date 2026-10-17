---
uuid: 82e3c31f-86d0-4832-85fd-9ab15236fd53
title: "HAWP v0.0.24 website and PR #41 lane checkpoint"
type: status
date: 2026-09-24
---

# Status Report

## Checkpoint Title

**2026-09-24 — HAWP v0.0.24 public website lane added while PR #41 hardening remains separately gated**

## Intent

Carry forward the existing HAWP v0.0.24 repository timeline without duplicating the completed PR #41 hardening checkpoint or unrelated project timelines.

## Current State

The checked-out `feature/website` lane is at `0db6b673`, also published as `origin/feature/website`. The PR #41 hardening lane remains separate at `6dface57` on `feature/v0.0.24` / `origin/feature/v0.0.24`; it remains unmerged and still requires explicit GitHub-side review-state readback before any merge decision. The working tree has untracked `website/.svelte-kit/` output and `website/package-lock.json`; these were preserved and not reconciled in this checkpoint.

## What Changed

**Before:** the latest checkpoint described a pushed but externally incomplete PR #41 security-hardening slice: nested corpus symlink rejection and schema-aware escaped-pipe reconciliation were locally verified, but fresh Copilot review, thread resolution, and GitHub check readback remained pending.

**After:** the repository also has a committed public website slice for v0.0.24: SvelteKit page structure, themed components, discoverability documentation, static metadata/assets, and a website CI workflow were added in `0db6b673`. This is a public-facing product/distribution lane, not evidence that PR #41 is reviewed or merge-ready. The two lanes remain intentionally separate.

## What Was Directly Verified

- `feature/website` and `origin/feature/website` both resolve to `0db6b673`.
- `feature/v0.0.24` and `origin/feature/v0.0.24` both resolve to `6dface57`.
- `0db6b673` adds 32 website/workflow files, including the SvelteKit app, CI workflow, discoverability baseline, static metadata, and social-card assets.
- The working tree is dirty only with untracked `website/.svelte-kit/` and `website/package-lock.json` in the inspected status output.
- The 2026-09-23 PR #41 local validation evidence remains recorded in `.hawp/work/status/2026/09/23/20bc9b7a/status.md`; it was not re-run as part of this documentation-only checkpoint.

## What Remains Unproven

- Website build, test, accessibility, deployment, and live rendering were not re-verified in this checkpoint.
- The untracked website artifacts have not been classified as intended generated output versus work needing cleanup.
- PR #41 GitHub inline-thread resolution, fresh Copilot review, and current CI/readback remain unverified.

## Strategic Impact

The repository now has two active strategic lanes: make the v0.0.24 public surface coherent and verifiable, while preserving the security-hardening lane’s narrow review gate. Priority is to validate the website and classify its untracked artifacts, then perform only the bounded GitHub review actions needed for PR #41. Do not merge PR #41 or treat the website commit as release proof without those checks.

## Constraints

No merge, external review request, cleanup, or generated-file deletion was authorized by this checkpoint. Preserve unrelated dirty state and keep the PR #41 and website lanes distinct unless new evidence materially connects them.

## Suggested Next Step

Run the website’s declared validation/build path and inspect its CI/deployment result; separately, when explicitly authorized, resolve only addressed PR #41 threads, request one fresh Lite review, and read back GitHub checks/review state.

## Memory Delta

Updated the existing repository HAWP v0.0.24 / PR #41 timeline with the material 2026-09-24 change: committed website lane `0db6b673`, its separation from hardening commit `6dface57`, current untracked artifacts, and the revised validation priority. No unrelated project checkpoint was duplicated or changed.

## 2026-10-13 continuation — website launch polish and MCP adoption

**Event and archive date:** 2026-10-13. This section extends the website timeline above; its current state supersedes the 2026-09-24 branch, verification, and next-step snapshot. This is a documentation checkpoint, not a work-item closure or website release.

### Before

- The public SvelteKit website existed on `feature/website`, with static prerendering, theme controls, release-label loading, and website CI. The 2026-09-24 checkpoint had not proved a deployed site or completed launch.
- A read-only active-work review found four non-website plans still `plan-ready`: `7394f90c` (parallel coordination), `c8d9e1cb` (remaining MCP intake elicitation and client guidance), `2eea565c` (MCP resource and prompt catalog), and `429e075e` (remote transport and authorization audit). The structured-result subset of `c8d9e1cb` had already landed for v0.0.24, while its remaining scope belongs to the proposed v0.0.25 lane.
- The user reported that HAWP MCP tools were used inconsistently during ordinary intake. This is an adoption observation, not proof that the tools are unavailable or defective.

### During: decisions and local work

- Priority shifted to finishing the nearly complete website before expanding MCP scope. The near-term sequence is website visual and interaction polish, website PR/release gates, then practical MCP intake guidance and discoverability.
- Two compoundable website plans were added locally to `.hawp/work/BACKLOG.md` and `.hawp/work/active/`: `96c45719` for one consistent pentagon mark, followed by `1607037f` for responsive, electrically styled hero cables. Both are `in-progress`, not `done`.
- The first nested, filled pentagon and square-framed treatment was rejected after visual review. The current local revision uses a single upright pentagon outline for the header and favicon, a pentagonal `WORK` hub in the hero, and matching social-card artwork. A responsive header review also moved the full navigation breakpoint to `lg` after the intermediate-width layout looked crowded.
- The five existing SVG connections gained themed base/current layers and bounded spring-like movement on drag/release. This reuses the existing hero, avoids a physics dependency, and retains reduced-motion handling. The cable plan still calls for touch, keyboard-focus, small-screen, and reduced-motion review before closure.

### After: verified state and limits

- At this checkpoint, local `feature/website` and `origin/feature/website` both resolve to `b12a1960785d002aabf13b52d9feba62ab0b139f` **before this documentation commit**. The revised pentagon, cable, backlog, and two plan changes are uncommitted local work and must be preserved separately from this checkpoint commit.
- The prior implementation pass reported `npm run check` with zero errors/warnings, `npm run build`, HAWP work validation, and `git diff --check` passing. Local browser review covered desktop and 390-pixel light/dark layouts; it led to the `lg` breakpoint correction. Those checks are prior-pass evidence, not reruns performed for this archival section.
- Live GitHub readback on 2026-10-13 found public PR #45 (`feature/website` to `development`) open and non-draft at `b12a1960`; `quality` had failed and merge state was `UNSTABLE`. A website build check had succeeded, while deployment was skipped/cancelled in the inspected check rollup. The v0.0.24 GitHub release is labeled **Pre-release**. Recheck these mutable states before acting.
- No production deployment, final aesthetic acceptance, complete interaction/accessibility review, or website launch was established here. The latest pentagon revision still awaits user visual review; the cable behavior has not yet met its plan's full verification checklist.

### Strategic impact and continuation

The website is the immediate product surface. Complete its two active visual items in order, then diagnose PR #45's failing `quality` check and verify the published workflow, deployment, metadata, and live page before claiming launch. MCP adoption is the next product/workflow slice: start with the existing `c8d9e1cb` plan's no-guessing intake/client guidance, then the `2eea565c` catalog. Keep `429e075e` as a remote-security planning lane and `7394f90c` as coordination work unless a concrete dependency changes that order.

**Resume from:** Inspect the local uncommitted pentagon and cable changes plus their existing plans; get visual acceptance for `96c45719`, finish the interaction and reduced-motion checks for `1607037f`, and update/close each existing work item only after its evidence is complete. Preserve unrelated working-tree changes. Then address PR #45's current `quality` failure and re-read GitHub checks and merge state. No merge or deployment is implied by this checkpoint.

## 2026-10-13 continuation — Analytics, local interaction review, and PR #45 quality diagnosis

### Current State

Website work remains local and uncommitted on `feature/website`. The pentagon and cable items remain `in-progress`, and new feature item `b4300599` tracks dedicated, opt-in GA4 for `hawp.online`.

### Directly Verified

- Signed-in GA4 showed the existing `sentzunhat-corp` property recorded 5 views and 5 active users for the HAWP project page on `sentzunhat.com`, Sep 15–Oct 12, 2026. That is not traffic for `hawp.online`.
- Created property `Human-AI Workflow Protocol` inside the existing `sentzunhat` Analytics account, using Winnipeg reporting time and CAD. Its `https://hawp.online` stream is ID `16121877274` with measurement ID `G-1FHPSYNPBS`; page views, scrolls, and outbound clicks are enabled.
- The property stream reported no data received yet and may take up to 48 hours to begin collecting. Website source now loads the Google tag only after explicit opt-in, records allow/decline locally, and retains a settings control.
- Local desktop and iPhone 16 Pro Max previews showed the pentagon, responsive hero, cable response on pointer drag, and consent UI. Selecting “No thanks” dismissed the prompt; the filtered Network panel showed no `googletagmanager` request. The initial notice overlays some hero artwork until a choice is made.
- Node 26.10.0 `npm run check` passed with zero errors and warnings; `npm run build` passed. `git diff --check` and HAWP work validation passed after the local whitespace cleanup.
- PR #45 is open on remote head `a01546a037e39506d76059c6cbdab47fe97d4e51`. Its `build` check passed, `deploy` was skipped, and `quality` failed because `website/docs/post-merge-discoverability-actions.md` had an extra blank line at EOF. That whitespace is corrected locally only; the remote check remains a failure on the unchanged remote head until a later commit updates it.

### What Remains Unproven

- Both `curl` and the signed-in browser failed to load the public `hawp.online` site during this session. Production deployment and current live analytics collection remain unverified.
- No post-deployment opt-in event has been observed in GA4. The connected GSC Wizard reader is unavailable because its trial ended.
- Direct touch-device behavior, final owner visual acceptance, and legal/privacy disclosure review remain open.

### Suggested Next Steps

1. Review and accept the pentagon/cable treatment; finish touch and reduced-motion interaction checks.
2. Commit/push the corrected local website changes only after review, then read PR #45 checks on the new head.
3. Verify Pages/custom-domain availability; after deployment, opt in on `hawp.online` and confirm the stream in GA4 Realtime.
4. Keep MCP adoption next: resume `c8d9e1cb`, then `2eea565c`.

No commit, push, merge, or deployment occurred in this continuation.

## 2026-10-13 continuation — Node 24 deployment, analytics, and dev-domain verification

### Checkpoint Title

**2026-10-13 — HAWP website deployed to the DigitalOcean app and `dev.hawp.online`; apex and GA4 collection remain open**

### Intent

Continue the existing website timeline after the owner approved Node 24 for the static-site build, requested the deployment be pushed, and supplied the DigitalOcean and custom-domain URLs. This addendum supersedes the earlier same-day local-only and DNS-pending snapshots below.

### Before → During → After

- **Before:** the DigitalOcean buildpack rejected the website's Node 26 engine range. Website visuals, opt-in analytics, and their work plans were local changes. The custom domain had not been proven reachable.
- **During:** the website engine was set to Node `24.x`; `.nvmrc` and GitHub Actions use `24.15.0`. Local Node `24.17.0` check/build passed. Commit `2d02ab94d2804e6a8d5bdb4eed92af48b087807c` was pushed to `feature/website`. The first DigitalOcean attempt built the static site successfully but failed asset upload while looking for `/workspace/website/website/build` instead of the emitted `build` directory. Later direct checks of the app URL and custom domain returned the deployed site successfully. The exact platform-side adjustment that resolved the initial upload-path failure was not captured.
- **After:** the site is reachable over HTTPS at both the generated DigitalOcean app URL and `https://dev.hawp.online/`. The apex `hawp.online` is not configured yet. A dedicated GA4 property/stream and a private HAWP traffic dashboard exist, but no post-deployment opt-in event has been verified.

### What Was Directly Verified

- Local `npm run check` with Node `24.17.0` reported zero errors and warnings; `npm run build` completed and emitted the static site.
- DigitalOcean's build log for the pushed revision showed SvelteKit's static adapter writing `build` successfully. That deployment attempt then failed during asset upload due to the path mismatch noted above.
- On 2026-10-13, `https://hawp-website-pg5fj.ondigitalocean.app/` returned HTTP 200 with the HAWP page title. Its `favicon.svg` also returned HTTP 200.
- DNS for `dev.hawp.online` resolves by CNAME to the generated DigitalOcean app hostname. On 2026-10-13, the custom domain and its favicon both returned HTTP 200 over HTTPS.
- GitHub PR #45 remains open from `feature/website` to `development`, at head `2d02ab94d2804e6a8d5bdb4eed92af48b087807c`. The Quality and Website build checks succeeded; the GitHub deploy job was skipped. Review decision is empty and merge state is `UNSTABLE`.
- The feature branch and `origin/feature/website` both read back at `2d02ab94d2804e6a8d5bdb4eed92af48b087807c`; the tree was clean before this checkpoint update.
- A private Looker Studio report now uses a HAWP GA4 data source and omits the source report's HAWP-inapplicable custom section chart. It currently has no traffic to display.

### What Remains Unproven

- No A/AAAA answer currently exists for apex `hawp.online`; a successful dev subdomain does not prove the apex is configured.
- GA4 has not shown a production event following an explicit visitor opt-in. Do not treat property/stream setup or a page load without consent as collection proof.
- The DigitalOcean UI's final deployment record after the later successful live responses was not read back; those responses prove current content delivery but not the precise platform-side change or deployment history.
- Final owner visual acceptance and the planned full touch, keyboard, responsive, and reduced-motion review of the pentagon/cable work remain open.

### Strategic Impact

The Node runtime blocker is resolved and work item `4b3f1581` is archived as done. The near-term website gates are now apex DNS/custom-domain setup, live consent-to-GA4 verification, and final visual/interaction acceptance. Keep the existing `b4300599`, `96c45719`, and `1607037f` items in progress until their own evidence is complete. Resume the planned MCP lane only after website launch gates are resolved; keep the remote transport and coordination plans separate.

### Constraints

This public repository checkpoint omits private analytics report links, account identifiers, and personal account details. No PR merge, review approval, or apex DNS change is implied. Zacatl is unaffected because Node 24 applies to the separate static website package.

### Continuation State

1. Configure the apex `hawp.online` in DigitalOcean App Platform and its DNS zone; verify HTTPS after the certificate becomes active.
2. With explicit visitor opt-in on the deployed site, verify the expected GA4 event in Realtime; then confirm the private dashboard populates.
3. Review and accept the pentagon treatment; finish cable touch/keyboard/reduced-motion checks and close the two visual work items only after evidence.
4. Re-read PR #45 checks, review state, and merge status before any merge decision.
5. Continue with HAWP intake guidance `c8d9e1cb`, then resource/prompt discovery `2eea565c`.

### Suggested Resume Point

**Resume from:** `feature/website` at pushed commit `2d02ab94`, with the static site live at `dev.hawp.online`; Node item `4b3f1581` is archived as done, while analytics and visual work remain active.

**Next objective:** configure and verify the apex domain, then confirm one consented production GA4 event without exposing private analytics links in public repository artifacts.

### Memory Delta

The existing local note `2026-10-13-hawp-website-launch-checkpoint.md` is extended by a follow-up memory update for this deployment milestone. The memory keeps the prior visual/analytics context, records the Node 24 push and live dev-domain proof, notes the initial upload-path failure without guessing how it was corrected, and points to this canonical HAWP timeline. Private Analytics IDs, report links, and unrelated project context are not repeated.
