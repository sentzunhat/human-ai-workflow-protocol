/**
 * CLI adapter for `hawp backlog validate`: argument parsing and help text.
 * No execution logic lives here — see script.ts.
 */

export interface BacklogValidateArgs {
  help: boolean;
  strictWarnings: boolean;
  /** Remaining args forwarded to the workflow validator. */
  workflowArgs: string[];
}

export const parseArgs = (argv: string[]): BacklogValidateArgs => ({
  help: argv.includes("--help") || argv.includes("-h"),
  strictWarnings: argv.includes("--strict-warnings"),
  workflowArgs: argv.filter(
    (arg) => arg !== "--strict-warnings" && arg !== "--help" && arg !== "-h",
  ),
});

export const getHelpText = (): string =>
  [
    "hawp backlog validate — validate kit/work drift in one command",
    "",
    "Usage:",
    "  ./.hawp/bin/hawp backlog validate [workflow-options]",
    "",
    "Behavior:",
    "  1) Runs distribution validation (.hawp/kit generated drift checks)",
    "  2) Runs workflow validation (.hawp/work backlog + plan checks)",
    "  3) Exits 1 if either validator fails",
    "",
    "Workflow options (passed through):",
    "  --hawp-root <path>",
    "  --work-root <path>",
    "  --debug-closed-task",
    "  --strict-warnings      Exit with code 1 if workflow warnings are present",
    "",
    "Examples:",
    "  ./.hawp/bin/hawp backlog validate",
    "  ./.hawp/bin/hawp backlog validate --work-root ./.hawp/work",
    "",
  ].join("\n");
