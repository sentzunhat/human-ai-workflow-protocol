export type FrontmatterValue = string | boolean | string[];

export interface MaterializationTarget {
  behavior: string;
  outputPath: string;
  frontmatter: Record<string, FrontmatterValue>;
  preamble?: string;
}

export interface MaterializationResult {
  outputPath: string;
  content: string;
}

export const GENERATED_BANNER =
  "<!-- Generated from core/providers/shared/behaviors — edit shared source and run npm --prefix librarian run providers:sync -->\n\n";

export const MATERIALIZATION_TARGETS: MaterializationTarget[] = [
  // GitHub — path-specific instructions (global entry stays in copilot-instructions.md)
  {
    behavior: "hawp-intake",
    outputPath: "core/providers/.github/instructions/hawp-intake.instructions.md",
    frontmatter: { applyTo: ".hawp/**,**/.hawp/**" },
    preamble:
      "Use this as a scoped, drop-in instruction when integrating HAWP into repos that already have Copilot instructions and prompts.\n\n",
  },
  {
    behavior: "hawp-backlog-alignment",
    outputPath:
      "core/providers/.github/instructions/hawp-backlog-alignment.instructions.md",
    frontmatter: { applyTo: "**" },
  },
  {
    behavior: "hawp-docs-alignment",
    outputPath:
      "core/providers/.github/instructions/hawp-docs-alignment.instructions.md",
    frontmatter: { applyTo: "docs/**,README.md,CHANGELOG.md,.hawp/kit/**,.hawp/work/**" },
  },
  // Claude Code
  {
    behavior: "hawp-core",
    outputPath: "core/providers/.claude/rules/hawp-core.md",
    frontmatter: {
      description:
        "HAWP workflow core — kit entry, status reports, backlog discipline",
    },
  },
  {
    behavior: "hawp-intake",
    outputPath: "core/providers/.claude/rules/hawp-intake.md",
    frontmatter: {
      paths: [".hawp/**"],
    },
  },
  {
    behavior: "hawp-backlog-alignment",
    outputPath: "core/providers/.claude/rules/hawp-backlog-alignment.md",
    frontmatter: {
      description: "HAWP backlog alignment and compaction guardrails",
    },
  },
  {
    behavior: "hawp-docs-alignment",
    outputPath: "core/providers/.claude/rules/hawp-docs-alignment.md",
    frontmatter: {
      paths: [".hawp/kit/**", ".hawp/work/**"],
    },
  },
  // Cursor
  {
    behavior: "hawp-core",
    outputPath: "core/providers/.cursor/rules/hawp-core.mdc",
    frontmatter: {
      description:
        "HAWP workflow core — kit entry, status reports, backlog discipline",
      alwaysApply: true,
    },
  },
  {
    behavior: "hawp-intake",
    outputPath: "core/providers/.cursor/rules/hawp-intake.mdc",
    frontmatter: {
      description: "HAWP modular intake for .hawp work folders",
      globs: "**/.hawp/**",
      alwaysApply: false,
    },
  },
  {
    behavior: "hawp-backlog-alignment",
    outputPath: "core/providers/.cursor/rules/hawp-backlog-alignment.mdc",
    frontmatter: {
      description: "HAWP backlog alignment and compaction guardrails",
      alwaysApply: true,
    },
  },
  {
    behavior: "hawp-docs-alignment",
    outputPath: "core/providers/.cursor/rules/hawp-docs-alignment.mdc",
    frontmatter: {
      description: "HAWP docs alignment when editing kit or workflow docs",
      globs: ".hawp/kit/**,.hawp/work/**",
      alwaysApply: false,
    },
  },
  // Continue — lexicographic order via filename prefix
  {
    behavior: "hawp-core",
    outputPath: "core/providers/.continue/rules/hawp-01-core.md",
    frontmatter: {
      name: "HAWP Core",
      description:
        "HAWP workflow core — kit entry, status reports, backlog discipline",
      alwaysApply: true,
    },
  },
  {
    behavior: "hawp-backlog-alignment",
    outputPath: "core/providers/.continue/rules/hawp-02-backlog-alignment.md",
    frontmatter: {
      name: "HAWP Backlog Alignment",
      description: "HAWP backlog alignment and compaction guardrails",
      alwaysApply: true,
    },
  },
  {
    behavior: "hawp-intake",
    outputPath: "core/providers/.continue/rules/hawp-03-intake.md",
    frontmatter: {
      name: "HAWP Intake",
      description: "HAWP modular intake for .hawp work folders",
      globs: ["**/.hawp/**"],
      alwaysApply: false,
    },
  },
  {
    behavior: "hawp-docs-alignment",
    outputPath: "core/providers/.continue/rules/hawp-04-docs-alignment.md",
    frontmatter: {
      name: "HAWP Docs Alignment",
      description: "HAWP docs alignment when editing kit or workflow docs",
      globs: [".hawp/kit/**", ".hawp/work/**"],
      alwaysApply: false,
    },
  },
];
