import { existsSync, mkdirSync, readFileSync, writeFileSync } from "fs";
import { dirname } from "path";
import {
  computeExpectedOutputs,
  findRepoRoot,
  normalizeForCompare,
} from "../shared/composition";

async function main(): Promise<void> {
  try {
    const repoRoot = findRepoRoot();
    if (!repoRoot) {
      console.error(
        "Error: Could not resolve repository root for distribution build.",
      );
      process.exit(1);
    }

    const outputs = computeExpectedOutputs(repoRoot);
    let wroteCount = 0;

    for (const output of outputs) {
      mkdirSync(dirname(output.outputPath), { recursive: true });

      const previous = existsSync(output.outputPath)
        ? normalizeForCompare(readFileSync(output.outputPath, "utf-8"))
        : null;

      if (previous !== output.content) {
        writeFileSync(output.outputPath, output.content, "utf-8");
        wroteCount++;
        console.log(`updated ${output.outputPath}`);
      } else {
        console.log(`unchanged ${output.outputPath}`);
      }
    }

    console.log(
      `\ndistribution build complete: ${wroteCount}/${outputs.length} file(s) updated`,
    );
    process.exit(0);
  } catch (error) {
    console.error("distribution build error:", error);
    process.exit(1);
  }
}

main();
