---
work-item: c1d2e400
type: audit
title: "Recursive audit: CLI command capabilities"
status: plan-ready
created: 2026-08-10
updated: 2026-08-10
parent: b6c4e8a2
follow-up: c1d2e401
---

# Audit: CLI Command Capabilities

## Mission

Audit nested CLI command groups for routing, transport formatting, provider
construction, and direct infrastructure access.

## Required Output

- Command-to-capability map.
- Safe mechanical split order.
- Findings that become implementation tasks in `c1d2e401`.

## Constraints

Keep command names, flags, output modes, and error behavior stable. Do not
perform broad handler rewrites during the audit.
