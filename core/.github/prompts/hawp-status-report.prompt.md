---
name: status-report
description: Generate a HAWP-style status report for the current work
---

Follow `.hawp/kit/templates/status-report.md`.

If `.hawp/kit/usage/status-report.md` exists in this repository, treat it as an optional repo-local supplement.

Using the current chat, workspace context, and any relevant open files, produce a status report that:

- stays focused on one intent/work thread unless the user explicitly wants combined reporting
- separates direct evidence from inference
- is compact, portable, and decision-useful
- includes a clear Help Wanted section when possible

If the user did not specify the kind of help wanted, infer a reasonable one based on the work, such as:

- review reasoning
- challenge assumptions
- help prioritize next step
- summarize tradeoffs
