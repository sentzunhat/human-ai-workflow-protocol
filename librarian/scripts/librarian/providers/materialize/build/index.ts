import { existsSync, mkdirSync, readFileSync, writeFileSync } from "node:fs";
import { dirname } from "node:path";
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
        "Error: Could not resolve repository root for provider materialize.",
      );
      process.exit(1);
    }

    const outputs = computeMaterializedOutputs(repoRoot);
    let wroteCount = 0;

    for (const output of outputs) {
      mkdirSync(dirname(output.outputPath), { recursive: true });

      const previous = existsSync(output.outputPath)
        ? normalizeForCompare(readFileSync(output.outputPath, "utf-8"))
        : null;

      const next = normalizeForCompare(output.content);

      if (previous !== next) {
        writeFileSync(output.outputPath, output.content, "utf-8");
        wroteCount++;
        console.log(`updated ${output.outputPath}`);
      } else {
        console.log(`unchanged ${output.outputPath}`);
      }
    }

    console.log(
      `\nprovider materialize complete: ${wroteCount}/${outputs.length} file(s) updated`,
    );
    process.exit(0);
  } catch (error) {
    console.error("provider materialize error:", error);
    process.exit(1);
  }
}

main();
