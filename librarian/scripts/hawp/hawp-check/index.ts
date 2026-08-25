#!/usr/bin/env node

/**
 * @deprecated Superseded by the Go CLI (`hawp check`). This TypeScript
 * implementation remains for reference only and is no longer invoked by
 * the npm scripts. Use `hawp check` or `npm run hawp:check` instead.
 *
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
