import { ExitCode, Mode, OutputFormat } from "./models/index.js";
import { runBacklogUpgradeScript } from "./script.js";

const supportedOutputFormats = new Set<string>(Object.values(OutputFormat));

/**
 * Parsed CLI options
 */
export interface CLIOptions {
  command: string; // "backlog"
  subcommand: string; // "upgrade"
  mode: Mode; // DryRun or Apply (defaults to DryRun)
  validate: boolean; // --validate flag
  exportPlan?: string; // --export-plan <path>
  format: OutputFormat; // text or json (defaults to text)
  output?: string; // --output <path>
  forceDirty: boolean; // --force-dirty flag
  verbose: boolean; // --verbose flag
  help: boolean; // --help flag
  version: boolean; // --version flag
}

export class CLIUsageError extends Error {
  public readonly exitCode = ExitCode.UsageError;

  public constructor(message: string) {
    super(message);
    this.name = "CLIUsageError";
  }
}

interface CLIIO {
  stdout: { write: (value: string) => void };
  stderr: { write: (value: string) => void };
}

/**
 * Parse command-line arguments into structured options
 *
 * Returns parsed options or throws CLIUsageError on invalid input
 */
export const parseArgs = async (args: string[]): Promise<CLIOptions> => {
  const opts: CLIOptions = {
    command: "",
    subcommand: "",
    mode: Mode.DryRun,
    validate: false,
    format: OutputFormat.Text,
    forceDirty: false,
    verbose: false,
    help: false,
    version: false,
  };

  // Parse arguments
  let i = 0;
  let modeSet = false;

  while (i < args.length) {
    const arg = args[i];

    if (arg === "--help" || arg === "-h") {
      opts.help = true;
      return opts;
    }

    if (arg === "--version" || arg === "-v") {
      opts.version = true;
      return opts;
    }

    if (arg === "backlog") {
      opts.command = "backlog";
      i++;
      continue;
    }

    if (arg === "upgrade") {
      opts.subcommand = "upgrade";
      i++;
      continue;
    }

    if (arg === "--dry-run") {
      if (modeSet && opts.mode !== Mode.DryRun) {
        throw new CLIUsageError(
          "Error: --dry-run and --apply are mutually exclusive. Choose one.",
        );
      }
      opts.mode = Mode.DryRun;
      modeSet = true;
      i++;
      continue;
    }

    if (arg === "--apply") {
      if (modeSet && opts.mode !== Mode.Apply) {
        throw new CLIUsageError(
          "Error: --dry-run and --apply are mutually exclusive. Choose one.",
        );
      }
      opts.mode = Mode.Apply;
      modeSet = true;
      i++;
      continue;
    }

    if (arg === "--validate") {
      opts.validate = true;
      i++;
      continue;
    }

    if (arg === "--force-dirty") {
      opts.forceDirty = true;
      i++;
      continue;
    }

    if (arg === "--verbose") {
      opts.verbose = true;
      i++;
      continue;
    }

    if (arg === "--format") {
      i++;
      const fmt = args[i];
      if (fmt && supportedOutputFormats.has(fmt)) {
        opts.format = fmt as OutputFormat;
      } else {
        throw new CLIUsageError(
          `Error: Invalid format '${fmt}'. Expected 'text' or 'json'.`,
        );
      }
      i++;
      continue;
    }

    if (arg === "--export-plan") {
      i++;
      const path = args[i];
      if (!path || path.startsWith("-")) {
        throw new CLIUsageError("Error: --export-plan requires a file path");
      }
      opts.exportPlan = path;
      i++;
      continue;
    }

    if (arg === "--output") {
      i++;
      const path = args[i];
      if (!path || path.startsWith("-")) {
        throw new CLIUsageError("Error: --output requires a file path");
      }
      opts.output = path;
      i++;
      continue;
    }

    // Unknown argument
    if (arg && arg.startsWith("-")) {
      throw new CLIUsageError(`Error: Unknown flag '${arg}'`);
    }

    i++;
  }

  // Validate required command/subcommand
  if (!opts.command || !opts.subcommand) {
    throw new CLIUsageError("Error: backlog upgrade <options> expected");
  }

  return opts;
};

/**
 * Show help message and exit
 */
export const getHelpText = (): string => `
hawp backlog upgrade — normalize and fix backlog structure

STATUS:
  dry-run detection pipeline enabled (TASK-028)
  apply mode scaffolds closed records and validates optionally

USAGE:
  hawp backlog upgrade [OPTIONS]

MODES (default: --dry-run):
  --dry-run              Select dry-run mode (default)
  --apply                Normalize closed records in place

OPTIONS:
  --validate             Run workflow validation summary after apply or dry-run
  --export-plan <path>   Export generated plan JSON to file
  --format <format>      Output format: text (default) or json
  --output <path>        Write rendered report to a file (stdout still used when omitted)
  --force-dirty          Skip git status checks in apply mode; allow dirty working tree
  --verbose              Show detailed diagnostic output

SHORTCUTS:
  --help, -h             Show this help message
  --version, -v          Show version

EXAMPLES:
  # Scan for drift without making changes (default)
  hawp backlog upgrade

  # Show results in JSON format
  hawp backlog upgrade --format json

  # Apply fixes and verify with validator
  hawp backlog upgrade --apply --validate

  # Export plan for review
  hawp backlog upgrade --export-plan upgrade-plan.json

NPM SCRIPT ALIAS (from librarian/):
  npm run workflow:normalize

SAFETY:
  - Always defaults to --dry-run
  - Dry-run does not modify project files
  - Detection output is read-only analysis of backlog/plan structure

For more information: https://github.com/beltrd/hawp
`;

/**
 * Show version and exit
 */
export const getVersionText = (): string =>
  "hawp backlog upgrade v1.1.0 (TASK-028 dry-run detection)";

export const showHelp = (): void => {
  process.stdout.write(getHelpText());
};

export const showVersion = (): void => {
  process.stdout.write(`${getVersionText()}\n`);
};

/**
 * Main CLI entry point
 *
 * Parses arguments, validates them, and returns parsed options
 */
export const runCLI = async (
  args: string[],
  io: CLIIO = { stdout: process.stdout, stderr: process.stderr },
): Promise<number> => {
  try {
    // Parse arguments
    const opts = await parseArgs(args);

    // Handle help/version first
    if (opts.help) {
      io.stdout.write(getHelpText());
      return ExitCode.Success;
    }

    if (opts.version) {
      io.stdout.write(`${getVersionText()}\n`);
      return ExitCode.Success;
    }

    // Validate command structure
    if (opts.command !== "backlog" || opts.subcommand !== "upgrade") {
      io.stderr.write("Error: invalid command structure\n");
      return ExitCode.UsageError;
    }

    if (opts.verbose) {
      io.stdout.write(
        `${JSON.stringify(
          {
            verbose: "Parsed options",
            mode: opts.mode,
            validate: opts.validate,
            format: opts.format,
            exportPlan: opts.exportPlan,
            output: opts.output,
            forceDirty: opts.forceDirty,
          },
          null,
          2,
        )}\n`,
      );
    }

    const scriptOptions = {
      mode: opts.mode,
      validate: opts.validate,
      format: opts.format,
      forceDirty: opts.forceDirty,
      verbose: opts.verbose,
      repoRoot: process.cwd(),
      ...(opts.exportPlan ? { exportPlan: opts.exportPlan } : {}),
      ...(opts.output ? { output: opts.output } : {}),
    };

    const result = await runBacklogUpgradeScript(scriptOptions);

    if (result.stdoutText.length > 0) {
      io.stdout.write(result.stdoutText);
    }

    for (const notice of result.notices) {
      io.stdout.write(`${notice}\n`);
    }

    for (const stderrLine of result.stderrLines) {
      io.stderr.write(`${stderrLine}\n`);
    }

    return result.exitCode;
  } catch (error) {
    if (error instanceof CLIUsageError) {
      io.stderr.write(`${error.message}\n`);
      return error.exitCode;
    }

    io.stderr.write(
      `CLI error: ${error instanceof Error ? error.message : String(error)}\n`,
    );
    return ExitCode.Error;
  }
};
