import { existsSync } from "node:fs";
import { join, resolve } from "node:path";
import { findUpward } from "../../lib";

/**
 * CLI argument definitions
 */
export interface CliArgs {
  help?: boolean;
  hawpRoot?: string;
  workRoot?: string;
  debugClosedTask?: boolean;
}

/**
 * Parse command-line arguments into CliArgs
 */
export const parseArgs = (argv: string[]): CliArgs => {
  const args: CliArgs = {};

  for (let i = 0; i < argv.length; i++) {
    const token = argv[i];
    if (!token) continue;

    if (token === "--help" || token === "-h") {
      args.help = true;
      continue;
    }

    if (token === "--hawp-root") {
      const value = argv[i + 1];
      if (value) args.hawpRoot = value;
      i++;
      continue;
    }

    if (token.startsWith("--hawp-root=")) {
      const value = token.split("=")[1];
      if (value) args.hawpRoot = value;
      continue;
    }

    if (token === "--work-root") {
      const value = argv[i + 1];
      if (value) args.workRoot = value;
      i++;
      continue;
    }

    if (token.startsWith("--work-root=")) {
      const value = token.split("=")[1];
      if (value) args.workRoot = value;
      continue;
    }

    if (token === "--debug-closed-task") {
      args.debugClosedTask = true;
      continue;
    }
  }

  return args;
};

/**
 * Generate help text
 */
export const getHelpText = (): string =>
  [
    "HAWP Workflow Validator",
    "",
    "Usage:",
    "  npm run workflow:validate -- [options]",
    "",
    "Default local behavior:",
    "  If no root flag is provided, the validator searches upward from the current",
    "  working directory for .hawp/work and validates that local work tree.",
    "",
    "Options:",
    "  --hawp-root <path>       Path to a .hawp directory; validates <path>/work.",
    "  --work-root <path>       Path to a .hawp/work directory.",
    "  --debug-closed-task      Include closed-task debug diagnostics for flagged files.",
    "  --help, -h               Show this help output.",
    "",
    "Exit code behavior:",
    "  0  Validation completed with no FAIL checks (PASS or WARN overall).",
    "  1  Validation failed (at least one FAIL check) or runtime/config error.",
    "",
    "Status meaning:",
    "  FAIL  Must-fix issue. Validation result is FAIL and exits with code 1.",
    "  WARN  Tolerated issue. Visible in report but does not fail validation.",
    "  INFO  Informational detail only. Not counted as WARN or FAIL.",
  ].join("\n");

/**
 * Find .hawp/work directory by searching up from current working directory
 */
const findWorkDirectory = (): string | null => {
  const root = findUpward(process.cwd(), (dir) =>
    existsSync(join(dir, ".hawp", "work")),
  );
  return root ? join(root, ".hawp", "work") : null;
};

/**
 * Resolve the work directory from CLI args or auto-discovery
 */
export const resolveWorkDirectory = (args: CliArgs): string | null => {
  if (args.workRoot) {
    const workRoot = resolve(args.workRoot);
    return existsSync(workRoot) ? workRoot : null;
  }

  if (args.hawpRoot) {
    const hawpRoot = resolve(args.hawpRoot);
    const candidate = join(hawpRoot, "work");
    return existsSync(candidate) ? candidate : null;
  }

  return findWorkDirectory();
};
