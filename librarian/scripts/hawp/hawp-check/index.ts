#!/usr/bin/env node

/**
 * Executable boundary for `hawp backlog validate`.
 * Parses argv via cli.ts and delegates to script.ts; only this file sets
 * the process exit code.
 */

import { getHelpText, parseArgs } from "./cli";
import { runBacklogValidateScript } from "./script";

const args = parseArgs(process.argv.slice(2));

if (args.help) {
  process.stdout.write(getHelpText());
  process.exitCode = 0;
} else {
  process.exitCode = runBacklogValidateScript(args);
}
