# Publication Safety Guidelines

Use this review when deciding whether HAWP content belongs in the public core.

## Keep in Public Core

- protocol definitions
- reusable templates
- fictional or clearly generic examples
- neutral prompts and review aids
- guidance that applies across many repositories

## Keep Out of Public Core

- real product audits
- internal roadmap or strategy details
- client or employer-specific workflows
- personal notes or status artifacts from private work
- secrets, credentials, logs, or local environment details

## Rewrite Rule

If an example is useful but anchored to a real project, rewrite it into a generic scenario instead of deleting it.

Examples:

- private launch audit -> generic SaaS readiness audit
- private repo migration -> documentation migration example
- personal workflow note -> project planning workflow example

## Review Questions

- Could a stranger understand this without knowing the original author or projects?
- Does any example imply a private business plan, client relationship, or internal launch?
- Does this file contain names, paths, or artifacts that should stay private?
- Is the example still practical after generic rewriting?

## Output Rule

When a file is not clearly public-safe, classify it as one of:

- keep
- rewrite
- move to private overlay
- remove
