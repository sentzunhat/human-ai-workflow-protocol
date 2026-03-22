#!/usr/bin/env node

/**
 * Backlog Upgrade CLI Entry Point
 *
 * Main executable for the backlog upgrade command.
 * Delegates to TypeScript implementation via tsx runtime.
 */

import { runCLI } from "./cli.js";

const args = process.argv.slice(2);
void runCLI(args)
  .then((exitCode) => {
    process.exitCode = exitCode;
  })
  .catch((error: unknown) => {
    console.error("Fatal error:", error);
    process.exitCode = 1;
  });
