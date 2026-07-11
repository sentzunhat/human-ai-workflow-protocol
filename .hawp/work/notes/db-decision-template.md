# Evidence Template: Database Decision

Use this template to record database schema decisions, compatibility choices, and migration planning notes.

## Decision Summary

**Date:** YYYY-MM-DD  
**Task ID:** TASK-XXX  
**Type:** schema design, migration, normalization, compatibility, or model equivalence  
**Scope:** what part of the data model is affected

## Context

- **Before state:** What existed before this decision.
- **Trigger:** Why the change is being made.
- **Stakeholders:** Who should review or approve the decision.

## Schema Changes

### SQL

```sql
-- Describe the SQL tables, columns, indexes, or views here.
```

### NoSQL

```text
# Describe the document shape, collection layout, or reference strategy here.
```

## Naming Rationale

- **Table or collection name:** Why this name was chosen.
- **Field or column names:** Why this naming pattern was used.
- **Embedded or referenced data:** Why the model uses that shape.

## Migration Path

- **New installs:** What gets created from scratch.
- **Existing data:** What backfill, sync, or compatibility steps are needed.
- **Deprecation:** What old names or structures remain temporarily.

## Testing Evidence

- **Rows affected:** How much data changes.
- **Validation commands:** What was run to confirm the design.
- **Observed results:** What the commands showed.

## Equivalence Notes

- **SQL equivalent:** How the model maps to relational storage.
- **NoSQL equivalent:** How the model maps to document storage.
- **Compatibility concerns:** Any differences that matter during migration.

## Reference

- **Schema standards:** [core/.hawp/kit/standards/database/README.md](../../../core/.hawp/kit/standards/database/README.md)
- **SQL standards:** [core/.hawp/kit/standards/database/sql.md](../../../core/.hawp/kit/standards/database/sql.md)
- **NoSQL standards:** [core/.hawp/kit/standards/database/nosql.md](../../../core/.hawp/kit/standards/database/nosql.md)
- **Related task:** link the relevant plan file here

## Sign-off

- [ ] Naming aligns with standards
- [ ] Compatibility path is documented
- [ ] Testing evidence is recorded
- [ ] SQL and NoSQL equivalence is noted where relevant
- [ ] Reviewer sign-off is captured
