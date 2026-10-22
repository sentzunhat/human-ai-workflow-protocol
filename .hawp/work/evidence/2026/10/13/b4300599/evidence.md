---

## Directly verified

- The user's signed-in Google Analytics Chrome session showed the existing account `sentzunhat` and existing property `sentzunhat-corp`.
- In the existing property, the date range Sep 15–Oct 12, 2026 showed the HAWP page hosted at `sentzunhat.com` with 5 views and 5 active users. This confirms that page is measured; it is not data for `hawp.online`.
- Created a separate property named **Human-AI Workflow Protocol** in the existing `sentzunhat` account. Reporting time zone is Winnipeg and currency is CAD.
- Created a web stream for `https://hawp.online`: stream ID `16121877274`, measurement ID `G-1FHPSYNPBS`.
- Enhanced measurement is on for page views, scrolls, and outbound clicks. Google Analytics reported that data collection may take up to 48 hours to start and showed no data received yet.
- Source now loads `gtag.js` only after an explicit allow choice, persists allow/decline locally, exposes an always-available settings control, and sets Google Consent Mode to denied when a visitor declines.
- The local desktop preview rendered the notice and both choices; selecting “No thanks” hid it and left the settings control available. No opt-in test was sent from the local preview, to avoid adding local development traffic to HAWP reporting.
- `npm run check` and `npm run build` passed with Node 26.10.0.
- PR #45 CI was read on remote head `a01546a037e39506d76059c6cbdab47fe97d4e51`: `build` passed, `deploy` was skipped, and `quality` failed because `website/docs/post-merge-discoverability-actions.md` had an extra blank line at EOF. The whitespace is fixed in the local working tree; no remote commit was created.
- Direct HTTPS and signed-in Chrome checks of `https://hawp.online/` did not load during this session.

## Still unverified

- Production deployment and accepted-consent traffic have not been verified. `hawp.online` did not respond to the bounded direct HTTPS check during this session.
- GA4 reporting for the new stream must be checked after deployment and visitor opt-in. The connected GSC Wizard reader is unavailable because its trial ended.
- Legal/privacy disclosure requirements have not been independently reviewed.
uuid: b4300599
title: "GA4 property and stream setup"
type: evidence
date: 2026-10-13
---
