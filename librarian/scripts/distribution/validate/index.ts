import { existsSync, readFileSync } from "fs";
import { join } from "path";
import {
  computeExpectedOutputs,
  findRepoRoot,
  normalizeForCompare,
} from "../shared/composition";

type PathLeak = {
  file: string;
  line: number;
  text: string;
};

const DOWNSTREAM_TARGET_FILES = [
  "core/.hawp/kit/instructions/da-file-tracking.md",
  "core/.hawp/kit/references/work-item-file-tracking.md",
  "core/.hawp/kit/references/install-update-safety.md",
  "core/.hawp/kit/templates/work-item-files.md",
  "core/.hawp/kit/templates/adr-template.md",
] as const;

const findDownstreamPathLeaks = (repoRoot: string): PathLeak[] => {
  const leaks: PathLeak[] = [];

  for (const relativePath of DOWNSTREAM_TARGET_FILES) {
    const absolutePath = join(repoRoot, relativePath);
    if (!existsSync(absolutePath)) {
      continue;
    }

    const lines = readFileSync(absolutePath, "utf-8").split(/\r?\n/);
    lines.forEach((line, index) => {
      if (line.includes("core/.hawp/")) {
        leaks.push({
          file: relativePath,
          line: index + 1,
          text: line.trim(),
        });
      }
    });
  }

  return leaks;
};

async function main(): Promise<void> {
  try {
    const repoRoot = findRepoRoot();
    if (!repoRoot) {
      console.error(
        "Error: Could not resolve repository root for distribution validate.",
      );
      process.exit(1);
    }

    const expectedOutputs = computeExpectedOutputs(repoRoot);
    const missing: string[] = [];
    const stale: string[] = [];
    const pathLeaks = findDownstreamPathLeaks(repoRoot);

    for (const expected of expectedOutputs) {
      if (!existsSync(expected.outputPath)) {
        missing.push(expected.outputPath);
        continue;
      }

      const current = normalizeForCompare(
        readFileSync(expected.outputPath, "utf-8"),
      );
      if (current !== expected.content) {
        stale.push(expected.outputPath);
      }
    }

    if (missing.length > 0 || stale.length > 0 || pathLeaks.length > 0) {
      console.error("distribution validation failed");

      if (missing.length > 0) {
        console.error("\nmissing generated outputs:");
        for (const file of missing) {
          console.error(`- ${file}`);
        }
      }

      if (stale.length > 0) {
        console.error("\nstale generated outputs:");
        for (const file of stale) {
          console.error(`- ${file}`);
        }
      }

      if (pathLeaks.length > 0) {
        console.error(
          "\ninvalid downstream path leaks (core/.hawp/) in source kit files:",
        );
        for (const leak of pathLeaks) {
          console.error(`- ${leak.file}:${leak.line} ${leak.text}`);
        }
      }

      console.error("\nrun npm run distribution:build");
      process.exit(1);
    }

    console.log(
      "distribution validation passed: generated outputs are current",
    );
    process.exit(0);
  } catch (error) {
    console.error("distribution validate error:", error);
    process.exit(1);
  }
}

main();
