import { join } from "path";
import { parseArgs, getHelpText, resolveWorkDirectory } from "./cli";
import { parseBacklog, orchestrateValidation } from "./orchestrate";
import { formatReport } from "./reporter";

/**
 * Main validator entry point
 */
const main = async (): Promise<void> => {
  try {
    const args = parseArgs(process.argv.slice(2));

    if (args.help) {
      console.log(getHelpText());
      process.exit(0);
    }

    // Determine work directory
    const workDir = resolveWorkDirectory(args);
    if (!workDir) {
      console.error("Error: Could not resolve .hawp/work directory");
      process.exit(1);
    }

    console.log(`Validating: ${workDir}\n`);

    // Parse backlog
    const backlogRows = parseBacklog(join(workDir, "BACKLOG.md"));
    if (!backlogRows) {
      console.error("Error: Could not parse BACKLOG.md");
      process.exit(1);
    }

    // Run validation
    const report = await orchestrateValidation(workDir, backlogRows, {
      debugClosedTask: args.debugClosedTask === true,
    });

    // Format and print report
    const reportText = formatReport(report);
    console.log(reportText);

    // Exit with appropriate code
    process.exit(report.overallStatus === "FAIL" ? 1 : 0);
  } catch (error) {
    console.error("Validation error:", error);
    process.exit(1);
  }
};

// Run validator
main();
