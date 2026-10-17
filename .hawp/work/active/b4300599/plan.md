# Add privacy-first Google Analytics to the HAWP website

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `b4300599-8938-4919-a32d-4fd3b0c89dc4`
**Type:** feature
**Reported:** 2026-10-13
**Status:** `in-progress`
**Risk Level:** medium
**Target Version:** website launch

## Investigation

The HAWP website source had no Google Analytics tag or visitor choice. The existing GA4 `sentzunhat-corp` property reports activity on the HAWP project page hosted at `sentzunhat.com`, but it does not cover the separate `hawp.online` domain.

**Directly inspected:** Google Analytics in the user's signed-in Chrome session, the HAWP site source, and the local browser preview.

## Mission

Measure page views, scrolls, and outbound link clicks on `hawp.online` in a dedicated GA4 property. Keep the Google tag unloaded until a visitor explicitly allows analytics, provide a decline option and a persistent settings control, and make revoking an existing choice update Google Consent Mode to denied.

## Plan

1. Create `Human-AI Workflow Protocol` as a property in the existing `sentzunhat` Analytics account; use Winnipeg time and CAD.
2. Add a web stream for `https://hawp.online` with enhanced measurement for page views, scrolls, and outbound clicks.
3. Add a visitor-facing allow/decline interface and a persistent settings control. Persist the choice locally; never load `gtag.js` before allow.
4. Verify consent UI, website build, and production collection after deployment. Do not count local preview or property creation as production data proof.
5. Provide a private Looker Studio dashboard based on the Sentzunhat Web Traffic report, mapped to the HAWP GA4 source and limited to HAWP reporting.

## Options

- **Recommended:** a dedicated property under the existing Sentzunhat Analytics account. This keeps HAWP reporting separate while using the current account organization.
- Add `hawp.online` to the existing Sentzunhat property. This mixes HAWP-only web reporting with the company website and makes clean attribution harder.
- Create a second Analytics account. This adds unnecessary account-level separation for one public project site.

## Work Coordination

**Owner:** Codex
**Implementation status:** in-progress
**Expected touched paths:**

- `website/src/lib/components/AnalyticsConsent.svelte`
- `website/src/routes/+layout.svelte`
- `.hawp/work/BACKLOG.md`

## Risk

Medium because the integration sends visitor data to Google Analytics. The tag must remain absent until explicit opt-in, reject must be available, and a visitor must be able to revisit settings. The site must not send personally identifiable information to Analytics.

## Verification

- Confirmed in Google Analytics that property `Human-AI Workflow Protocol` and web stream `https://hawp.online` exist.
- Stream ID: `16121877274`; Measurement ID: `G-1FHPSYNPBS`.
- Enhanced measurement is on for page views, scrolls, and outbound clicks.
- In the Analytics property, the stream reports no data received in the past 48 hours and the stream details warn that website data collection is not active. This verifies the property/stream configuration only, not production tracking.
- Created a private Looker Studio copy titled `Human-AI Workflow Protocol Web Traffic`, using a reusable HAWP GA4 data source. Removed the Sentzunhat-specific custom homepage-section chart. The standard traffic panels are present and show no data yet.
- DNS now points to DigitalOcean nameservers, but no apex A/AAAA or `www` CNAME answers are present; `https://hawp.online` does not resolve. Therefore neither production consent flow nor accepted-consent events can yet be verified.
- `npm run check` using Node 26.10.0: zero errors and zero warnings.
- `npm run build` using Node 26.10.0: passed; static output generated.
- Local desktop preview showed the pentagon, cable layers, and consent controls. No Google Analytics tag was exercised in the local preview to avoid contaminating production reporting with local traffic.
- Continuation update (2026-10-13): `https://dev.hawp.online/` and its favicon respond over HTTPS with HTTP 200 and serve the current HAWP static site. Production measurement remains unverified until a visitor opts in and the event appears in GA4 Realtime. The GSC Wizard Analytics reader is unavailable because its trial ended.

## Close Checklist

- [ ] Production deployment and accepted-consent event verified in GA4 Realtime.
- [x] Remove the inherited Sentzunhat-only homepage-section chart from the Looker dashboard.
- [ ] Consent UI and revocation reviewed at desktop and mobile widths.
- [ ] Outcome and verification recorded; plan archived; backlog updated.
