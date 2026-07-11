import { existsSync, readFileSync } from "node:fs";
import { join } from "node:path";
import { findUpward, normalizeForCompare } from "../../../lib";
import type {
  BuildResult,
  CompositionPlan,
  DistributionVariant,
} from "./types";

export { normalizeForCompare };

const ACTIVE_PROVIDERS = [
  "claude",
  "codex",
  "github",
  "cursor",
  "continue",
] as const;

function providerVariant(
  provider: (typeof ACTIVE_PROVIDERS)[number],
  operation: "install" | "update",
  ref: "main" | "dev",
): DistributionVariant {
  return {
    outputFile: `${provider}/${operation}/${ref}.md`,
    sectionFiles: [
      `sources/providers/${provider}/preamble-${operation}.md`,
      "sources/shared/safety.md",
      `sources/providers/${provider}/safety.md`,
      "sources/shared/repo-boundaries-kit.md",
      `sources/providers/${provider}/boundaries.md`,
      `sources/shared/${operation}.md`,
      `sources/providers/${provider}/${operation}-contract.md`,
      `sources/providers/${provider}/${operation}/${ref}.md`,
    ],
    scriptPartFiles: [
      `distribution/sources/${operation}/script-core.md`,
      `distribution/sources/providers/${provider}/script-${operation}.md`,
      `distribution/sources/${operation}/script-footer.md`,
    ],
    ref,
    provider,
    operation,
  };
}

const DISTRIBUTION_PLAN: CompositionPlan = {
  variants: ACTIVE_PROVIDERS.flatMap((provider) => [
    providerVariant(provider, "install", "main"),
    providerVariant(provider, "install", "dev"),
    providerVariant(provider, "update", "main"),
    providerVariant(provider, "update", "dev"),
  ]),
};

/** @deprecated Root-level install-main.md etc. — removed; use <provider>/install/main.md */
export const LEGACY_ROOT_GUIDES = [
  "install-main.md",
  "install-dev.md",
  "update-main.md",
  "update-dev.md",
] as const;

function extractBashBody(repoRoot: string, relativePath: string): string {
  const filePath = join(repoRoot, relativePath);
  if (!existsSync(filePath)) {
    throw new Error(`Script part not found: ${filePath}`);
  }
  const content = readFileSync(filePath, "utf-8");
  const blockRegex = /```bash\n([\s\S]*?)\n```/;
  const match = blockRegex.exec(content);
  if (!match?.[1]) {
    throw new Error(`No bash block in ${filePath}`);
  }
  return match[1].trimEnd();
}

function composeProviderScript(
  repoRoot: string,
  variant: DistributionVariant,
): string {
  const body = variant.scriptPartFiles
    .map((part) => extractBashBody(repoRoot, part))
    .join("\n\n");

  const refLineRegex = /(^|\n)REF="[^"]*"(?=\n|$)/;
  const providerLineRegex = /(^|\n)PROVIDER="[^"]*"(?=\n|$)/;

  if (!refLineRegex.test(body) || !providerLineRegex.test(body)) {
    throw new Error(
      `Composed script missing REF or PROVIDER for ${variant.outputFile}`,
    );
  }

  let substituted = body.replace(refLineRegex, `$1REF="${variant.ref}"`);
  substituted = substituted.replace(
    providerLineRegex,
    `$1PROVIDER="${variant.provider}"`,
  );

  return `\`\`\`bash\n${substituted}\n\`\`\``;
}

function generateScriptSection(
  repoRoot: string,
  variant: DistributionVariant,
): string {
  const operation =
    variant.operation === "install" ? "Install" : "Update";
  const bashBlock = composeProviderScript(repoRoot, variant);
  return (
    `## ${operation} Command (Copy/Paste)\n\n` +
    `Run this from the root of your target repository. ` +
    `No edits are required; branch and provider are already configured in the command. ` +
    `Each run fetches the latest commit from that branch.\n\n` +
    `${bashBlock}`
  );
}

function generateSourceReferenceSection(variant: DistributionVariant): string {
  const fragments = variant.sectionFiles
    .map((f) => `- \`distribution/${f}\``)
    .join("\n");
  const scriptParts = variant.scriptPartFiles
    .map((f) => `- \`${f}\``)
    .join("\n");

  return (
    `## Source Reference\n\n` +
    `This file is generated. Do not edit it directly.\n\n` +
    `- Workflow gate: pushes and pull requests on \`main\` or \`dev\` fail when generated guides drift from source.\n` +
    `- Local sync: run \`npm --prefix librarian run distribution:sync\` after editing \`distribution/sources/\` or the distribution composition code.\n\n` +
    `Generated output file:\n\n` +
    `- \`distribution/generated/${variant.outputFile}\`\n\n` +
    `Provider: \`${variant.provider}\` · Operation: \`${variant.operation}\` · Branch: \`${variant.ref}\`\n\n` +
    `Install mapping: \`core/providers/.${variant.provider}/\` → downstream paths in this guide.\n\n` +
    `This generated guide is built from:\n\n` +
    `${fragments}\n\n` +
    `Composed shell script (core + provider overlay + footer):\n\n` +
    `${scriptParts}`
  );
}

export function findRepoRoot(startDir = process.cwd()): string | null {
  return findUpward(
    startDir,
    (dir) =>
      existsSync(join(dir, "distribution", "sources", "shared")) &&
      existsSync(join(dir, "librarian", "package.json")),
  );
}

export function computeExpectedOutputs(repoRoot: string): BuildResult[] {
  const sourceRoot = join(repoRoot, "distribution");
  const outputRoot = join(repoRoot, "distribution", "generated");

  return DISTRIBUTION_PLAN.variants.map((variant) => {
    const sections = variant.sectionFiles.map((sectionFile) => {
      const sourcePath = join(sourceRoot, sectionFile);
      if (!existsSync(sourcePath)) {
        throw new Error(`Missing source fragment: ${sourcePath}`);
      }

      const content = readFileSync(sourcePath, "utf-8");
      return normalizeForJoin(content);
    });

    const scriptSection = normalizeForJoin(
      generateScriptSection(repoRoot, variant),
    );
    const refSection = normalizeForJoin(
      generateSourceReferenceSection(variant),
    );

    const joined = [...sections, scriptSection, refSection].join("\n\n");

    return {
      outputPath: join(outputRoot, variant.outputFile),
      content: ensureTrailingNewline(joined),
    };
  });
}

function normalizeForJoin(content: string): string {
  return normalizeForCompare(content).trimEnd();
}

function ensureTrailingNewline(content: string): string {
  return content.endsWith("\n") ? content : `${content}\n`;
}
