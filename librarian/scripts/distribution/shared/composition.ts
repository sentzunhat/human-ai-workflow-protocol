import { existsSync, readFileSync } from "fs";
import { dirname, join, resolve } from "path";
import type {
  BuildResult,
  CompositionPlan,
  DistributionVariant,
} from "./types";

const DISTRIBUTION_PLAN: CompositionPlan = {
  variants: [
    {
      outputFile: "install-main.md",
      sectionFiles: [
        "sources/install/preamble.md",
        "sources/shared/safety.md",
        "sources/shared/repo-boundaries.md",
        "sources/shared/install.md",
        "sources/install/main.md",
      ],
      scriptSourceFile: "distribution/sources/install/script.md",
      ref: "main",
    },
    {
      outputFile: "install-dev.md",
      sectionFiles: [
        "sources/install/preamble.md",
        "sources/shared/safety.md",
        "sources/shared/repo-boundaries.md",
        "sources/shared/install.md",
        "sources/install/dev.md",
      ],
      scriptSourceFile: "distribution/sources/install/script.md",
      ref: "dev",
    },
    {
      outputFile: "update-main.md",
      sectionFiles: [
        "sources/update/preamble.md",
        "sources/shared/safety.md",
        "sources/shared/repo-boundaries.md",
        "sources/shared/update.md",
        "sources/update/main.md",
      ],
      scriptSourceFile: "distribution/sources/update/script.md",
      ref: "main",
    },
    {
      outputFile: "update-dev.md",
      sectionFiles: [
        "sources/update/preamble.md",
        "sources/shared/safety.md",
        "sources/shared/repo-boundaries.md",
        "sources/shared/update.md",
        "sources/update/dev.md",
      ],
      scriptSourceFile: "distribution/sources/update/script.md",
      ref: "dev",
    },
  ],
};

function extractBashBlock(
  repoRoot: string,
  scriptSourceFile: string,
  ref: string,
): string {
  const filePath = join(repoRoot, scriptSourceFile);
  if (!existsSync(filePath)) {
    throw new Error(`Script source file not found: ${filePath}`);
  }
  const content = readFileSync(filePath, "utf-8");
  const blockRegex = /```bash\n([\s\S]*?)\n```/g;
  const refLineRegex = /(^|\n)REF="[^"]*"(?=\n|$)/;
  let match: RegExpExecArray | null;
  while ((match = blockRegex.exec(content)) !== null) {
    const blockBody = match[1];
    if (
      blockBody !== undefined &&
      blockBody.includes('OWNER="') &&
      blockBody.includes('REPO="') &&
      refLineRegex.test(blockBody)
    ) {
      const body = blockBody.replace(refLineRegex, `$1REF="${ref}"`);
      return `\`\`\`bash\n${body}\n\`\`\``;
    }
  }
  throw new Error(
    `No install/update bash block with OWNER/REPO/REF found in ${filePath}`,
  );
}

function generateScriptSection(
  repoRoot: string,
  variant: DistributionVariant,
): string {
  const operation = variant.outputFile.startsWith("install")
    ? "Install"
    : "Update";
  const bashBlock = extractBashBlock(
    repoRoot,
    variant.scriptSourceFile,
    variant.ref,
  );
  return (
    `## ${operation} Command (Copy/Paste)\n\n` +
    `Run this from the root of your target repository. ` +
    `No edits are required; branch selection is already configured in the command. ` +
    `Each run fetches the latest commit from that branch.\n\n` +
    `${bashBlock}`
  );
}

function generateSourceReferenceSection(variant: DistributionVariant): string {
  return (
    `## Source Reference\n\n` +
    `This file is generated. Do not edit it directly.\n\n` +
    `- Workflow gate: pushes and pull requests on \`main\` or \`dev\` fail when generated guides drift from source.\n` +
    `- Local sync: run \`npm --prefix librarian run distribution:sync\` after editing \`distribution/sources/\` or the distribution composition code.\n\n` +
    `Generated output file:\n\n` +
    `- \`distribution/generated/${variant.outputFile}\`\n\n` +
    `This generated guide is built from:\n\n` +
    `- \`${variant.scriptSourceFile}\` — shell script (authoritative source)`
  );
}

export function findRepoRoot(startDir = process.cwd()): string | null {
  let current = resolve(startDir);

  for (let i = 0; i < 10; i++) {
    const hasDistribution = existsSync(
      join(current, "distribution", "sources", "shared"),
    );
    const hasLibrarianPackage = existsSync(
      join(current, "librarian", "package.json"),
    );
    if (hasDistribution && hasLibrarianPackage) {
      return current;
    }

    const parent = dirname(current);
    if (parent === current) {
      break;
    }
    current = parent;
  }

  return null;
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

export function normalizeForCompare(content: string): string {
  return content.replace(/\r\n/g, "\n");
}

function normalizeForJoin(content: string): string {
  return normalizeForCompare(content).trimEnd();
}

function ensureTrailingNewline(content: string): string {
  return content.endsWith("\n") ? content : `${content}\n`;
}
