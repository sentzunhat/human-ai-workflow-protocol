# Shared Standards Classification Rubric & Red-Flag Lexicon

**Purpose:** Classify shared standards entries as public-safe, repo-specific, or private/proprietary during TASK-058 triage.

**Date:** 2026-05-15  
**Related:** TASK-057, TASK-058

---

## Classification Levels

### Level 1: Public-Safe ✓

**Criteria:**

- Framework-agnostic or universally useful (applies across TypeScript/Node.js/database projects generally)
- No internal system references or project-specific naming
- No implementation runbooks for internal-only systems
- No security-sensitive architecture details (auth secrets, internal API key registration, protected routes internal logic)
- Portable and documented well enough for external projects to adopt
- Not explicitly marked as private or project-specific in source

**Examples:**

- Code style guidelines for TypeScript
- Testing patterns and frameworks (Vitest setup, test organization)
- Git workflow conventions (Conventional Commits)
- MongoDB schema design patterns
- Generic architecture principles (layered architecture, ports-and-adapters)
- Documentation guidelines (JSDoc, READMEs, changelog format)

**Action:** Absorb directly into `core/.hawp/kit/standards/**`

---

### Level 2: Repo-Specific Adaptation

**Criteria:**

- Useful and reusable, but tailored to this repo's specific tooling, framework, or workflow
- May reference this repo's conventions (e.g., area-based composition specific to this project's structure)
- Not proprietary, but not universally applicable without modification
- Safe to absorb with context markers or repo-specific notes

**Examples:**

- Node.js specific project structure (this repo's folder organization)
- Build pipeline specifics (esbuild, path resolution, module aliasing)
- Area-based feature composition (specific to this repo's architecture)
- Zacatl-agnostic layering patterns (but the examples may be Zacatl-specific)

**Action:** Absorb with repo-context note or adapt before absorption

---

### Level 3: Private/Proprietary ⛔

**Criteria:**

- References internal-only systems or projects by name
- Includes implementation runbooks, internal architecture playbooks, or secret-sauce recipes
- Covers internal authentication, secret management, or internal API patterns
- Tagged as project-specific, private, or org-internal
- Would expose internal infrastructure, security practices, or strategic decisions if made public
- Contains internal business logic or product-specific workflows

**Action:** DO NOT absorb. Create a split follow-up task in `.hawp/work/active/` for later evaluation or internal-only use.

---

## Red-Flag Lexicon

**Keywords and terms that indicate private/proprietary content:**

| Red Flag                                    | Category                     | Action                        |
| ------------------------------------------- | ---------------------------- | ----------------------------- |
| `Tekit`                                     | Internal project             | Mark private                  |
| `Mictlan`                                   | Internal project             | Mark private                  |
| `Zacatl`                                    | Internal project (framework) | Mark private                  |
| `mictlan/`                                  | Project path                 | Mark private                  |
| `tekit/`                                    | Project path                 | Mark private                  |
| `zacatl/`                                   | Project path                 | Mark private                  |
| `internal-runtime`                          | Internal domain              | Mark private                  |
| `product/`                                  | Product-specific             | Mark private                  |
| `providers/`                                | Internal provider registry   | Mark private                  |
| `protected-routes`                          | Auth/Security internal       | Mark private                  |
| `API key registration`                      | Security internal            | Mark private                  |
| `handler-responsibilities` (Zacatl context) | Framework-specific           | Check context; may be generic |
| `dependency-registration` (Zacatl context)  | Framework-specific           | Check context; may be generic |
| `contract-testing` (Zacatl context)         | Framework-specific           | Check context; may be generic |
| `service-boundaries` (Zacatl context)       | Framework-specific           | Check context; may be generic |
| `layered-composition` (Zacatl context)      | Framework-specific           | Check context; may be generic |
| `auth/`                                     | Authentication domain        | Mark private                  |
| `security` (path prefix)                    | Security domain              | Mark private                  |
| `private/` (folder prefix)                  | Explicit                     | Mark private                  |

---

## Classification Checklist

For each entry in shared_standards, apply this checklist:

**Quick Scan:**

- [ ] Is it in `shared_standards/private/`? → **Private** ⛔
- [ ] Is it in `shared_standards/project-specific/`? → **Repo-Specific or Private** (check context)
- [ ] Is it in `shared_standards/public/`? → **Likely Public-Safe** ✓ (verify below)

**Deep Dive (for public/ entries):**

- [ ] Does it mention Tekit, Mictlan, Zacatl, or other internal projects? → **Private** ⛔
- [ ] Does it reference `internal-runtime`, `product`, `providers`, or other internal domains? → **Private** ⛔
- [ ] Does it include implementation runbooks for internal-only systems? → **Private** ⛔
- [ ] Does it cover internal auth, secrets, or API key management? → **Private** ⛔
- [ ] Is it tagged as "project-specific" or "internal"? → **Private** ⛔
- [ ] Is it universally useful across TypeScript/Node.js/database projects? → **Public-Safe** ✓
- [ ] Could an external project adopt it as-is without modification? → **Public-Safe** ✓

**Ambiguous Cases:**

- If ambiguous, default to **Repo-Specific** and flag for human review.
- Document the ambiguity reason in the triage table.

---

## Example Classifications

### Example 1: MongoDB Schema Design Guidelines

**Source:** `shared_standards/public/standards/database/mongodb-schema-design.md`

- In `public/` folder ✓
- No internal system references ✓
- No red flags found ✓
- Universal patterns (embedding, lifecycle, cardinality) ✓

**Classification:** `public-safe` → Absorb to `core/.hawp/kit/standards/database/`

---

### Example 2: Code Style Guide

**Source:** `shared_standards/public/guidelines/code-style.md`

- In `public/` folder ✓
- TypeScript conventions only ✓
- Node 24.14.0+ requirement is repo-specific but documented ✓
- No internal references ✓

**Classification:** `public-safe` → Absorb to `core/.hawp/kit/standards/guidelines/`

---

### Example 3: Area Composition (Zacatl)

**Source:** `shared_standards/public/standards/nodejs/area-composition.md` (if references Zacatl implementation specifics)

- In `public/` folder ✓
- But title mentions "Backend Area Module Composition" with Zacatl-specific patterns
- If it includes Zacatl-only handler registration or dependency-injection specifics → **Repo-Specific** (adapt)
- If it's generic layering patterns with Zacatl as one example → **Public-Safe** (absorb as-is)

**Classification:** Depends on content; check for red flags first, then decide.

---

### Example 4: Internal Routes

**Source:** `shared_standards/private/auth/protected-routes.md`

- In `private/` folder ⛔
- Explicitly labeled private
- Title mentions "protected-routes" (auth-related) → Red flag ⛔
- Internal implementation details for route protection

**Classification:** `private` → DO NOT absorb. Create follow-up task if needed.

---

### Example 5: Zacatl Layering

**Source:** `shared_standards/public/standards/zacatl/layered-composition.md`

- In `public/` folder ✓
- But title explicitly mentions "zacatl" → Red flag (Zacatl is internal)
- If it's Zacatl-specific implementation → **Private**
- If it's generic layering with Zacatl examples → **Repo-Specific**

**Classification:** Assume `private` unless proven otherwise. Default to flag for review.

---

## Edge Cases & Resolution Guidance

### "Is this Zacatl-specific or generic?"

**Rule:** If the standard assumes Zacatl's dependency-injection container, handler lifecycle, or provider registration pattern, treat as **repo-specific or private**. If it's generic layering (application/domain/infrastructure) with optional Zacatl examples, treat as **public-safe**.

### "What if a public/ item has internal references only in one section?"

**Rule:** Extract the generic section, note the split, and create a repo-specific follow-up task for the internal parts.

### "Ambiguous folder names like `nodejs/` or `database/`?"

**Rule:** Check the file content, not the folder. Use the checklist above.

---

## Triage Table Template (for TASK-058)

Use this table format when running TASK-058:

| Source                                               | Entry                      | Red Flags           | Classification | Action | Follow-up Task | Notes                                              |
| ---------------------------------------------------- | -------------------------- | ------------------- | -------------- | ------ | -------------- | -------------------------------------------------- |
| `public/standards/database/mongodb-schema-design.md` | MongoDB embedding patterns | None                | `public-safe`  | Absorb | —              | Absorbed to `core/.hawp/kit/standards/database/`   |
| `private/auth/protected-routes.md`                   | Protected routes           | `auth/`, `private/` | `private`      | Split  | TASK-NNN       | Route protection internals; defer to separate task |
| ...                                                  | ...                        | ...                 | ...            | ...    | ...            | ...                                                |

---

## Related Documents

- Exclusion policy: `core/.hawp/kit/standards/README.md`
- Standards consolidation map: `core/.hawp/kit/standards/README.md`
- Task-058 triage execution: `.hawp/work/active/TASK-058.md`
