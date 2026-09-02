# Benchmark Run: With HAWP
**Task (Input):** "review this codebase and tell me what's wrong"

**HAWP Shape:**
```ts
const shape: Shape = {
  input: "review this codebase and tell me what's wrong",
  context: "The repository is 'human-ai-workflow-protocol', a framework for managing AI workflows using HAWP. The user wants to identify potential architectural issues, maintenance risks, or inconsistencies in the current structure.",
  mission: "Perform a bounded review of the repository structure and documentation to identify high-risk areas (e.g., complexity, fragmentation, or lack of clarity) that could impede developer onboarding or maintainability.",
  constraints: "Focus only on the top-level directory structure and the .github/instructions folder. Do not analyze individual file contents unless they represent a structural pattern. Avoid speculative findings; report only what is observable from the file tree and READMEs. Use evidence-based observations.",
  output: "A structured report listing identified risks, categorized by type (e.g., Complexity, Maintenance, Onboarding), with specific file paths as evidence."
};
```

## Output
### Identified Risks

#### 1. Complexity & Cognitive Load
- **Observation**: The `distribution/sources/providers/` directory contains highly repetitive structures for `continue`, `cursor`, and `github` (e.g., `boundaries.md`, `install-contract.md`, `preamble-install.md`).
- **Evidence**: `distribution/sources/providers/continue/` vs `distribution/sources/providers/cursor/`.
- **Risk**: High risk of "configuration drift" where updates to the HAWP standard are applied to one provider but missed in others, leading to inconsistent developer experiences across different AI agents.

#### 2. Documentation Fragmentation
- **Observation**: There is a high density of instruction files in `.github/instructions/` and various `README.md` files throughout the tree.
- **Evidence**: `.github/instructions/hawp-intake.instructions.md`, `librarian/README.md`, `core/providers/shared/README.md`.
- **Risk**: Increased maintenance burden. Developers may struggle to find the "single source of truth" for workflow rules, potentially leading to violations of the HAWP protocol during manual updates.

#### 3. Maintenance Scalability
- **Observation**: The `librarian/` directory contains a complex set of scripts and tests (e.g., `backlog-upgrade`, `backlog-validate`).
- **Evidence**: `librarian/scripts/backlog-upgrade/index.ts`.
- **Risk**: As the HAWP protocol evolves, these internal tools require significant upkeep to remain compatible with changes in the `.hawp/` directory structure.

### Summary of Findings
The primary risks are **structural duplication** and **documentation fragmentation**. While the repository is highly organized, the current pattern of duplicating provider-specific logic and instructions creates a high surface area for errors during protocol updates.
