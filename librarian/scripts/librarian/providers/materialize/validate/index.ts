import { existsSync, readFileSync } from "node:fs";
import { findRepoRoot } from "../../../distribution/shared/composition";
import {
  computeMaterializedOutputs,
  normalizeForCompare,
} from "../render";

function main(): void {
  try {
    const repoRoot = findRepoRoot();
    if (!repoRoot) {
      console.error(
        "Error: Could not resolve repository root for provider validate.",
      );
      process.exit(1);
    }

    const expected = computeMaterializedOutputs(repoRoot);
    const missing: string[] = [];
    const stale: string[] = [];

    for (const item of expected) {
      if (!existsSync(item.outputPath)) {
        missing.push(item.outputPath);
        continue;
      }

      const current = normalizeForCompare(
        readFileSync(item.outputPath, "utf-8"),
      );
      if (current !== normalizeForCompare(item.content)) {
        stale.push(item.outputPath);
      }
    }

    if (missing.length > 0 || stale.length > 0) {
      console.error("provider materialization validation failed");

      if (missing.length > 0) {
        console.error("\nmissing materialized outputs:");
        for (const file of missing) {
          console.error(`- ${file}`);
        }
      }

      if (stale.length > 0) {
        console.error("\nstale materialized outputs:");
        for (const file of stale) {
          console.error(`- ${file}`);
        }
      }

      console.error("\nrun npm run providers:materialize");
      process.exit(1);
    }

    console.log(
      `provider validation passed: ${expected.length} materialized file(s) are current`,
    );
    process.exit(0);
  } catch (error) {
    console.error("provider validate error:", error);
    process.exit(1);
  }
}

main();
