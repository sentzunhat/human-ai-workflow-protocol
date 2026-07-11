/**
 * Workflow validation runner: resolves the work directory, runs all
 * validation checks, and prints the formatted report.
 * The executable boundary (index.ts) owns the process exit code.
 */

import { join } from "node:path";
import { parseArgs, getHelpText, resolveWorkDirectory } from "./cli";
import { parseBacklog, orchestrateValidation } from "./orchestrate";
import { formatReport } from "./reporter";

export const runWorkflowValidation = (argv: string[]): number => {
  try {
    const args = parseArgs(argv);

    if (args.help) {
      console.log(getHelpText());
      return 0;
    }

    const workDir = resolveWorkDirectory(args);
    if (!workDir) {
      console.error("Error: Could not resolve .hawp/work directory");
      return 1;
    }

    console.log(`Validating: ${workDir}\n`);

    const backlogRows = parseBacklog(join(workDir, "BACKLOG.md"));
    if (!backlogRows) {
      console.error("Error: Could not parse BACKLOG.md");
      return 1;
    }

    const report = orchestrateValidation(workDir, backlogRows, {
      debugClosedTask: args.debugClosedTask === true,
    });

    console.log(formatReport(report));

    return report.overallStatus === "FAIL" ? 1 : 0;
  } catch (error) {
    console.error("Validation error:", error);
    return 1;
  }
};
