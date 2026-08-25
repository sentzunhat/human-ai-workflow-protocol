#!/usr/bin/env node

/**
 * @deprecated Superseded by the Go CLI (`hawp work normalize`). This TypeScript
 * implementation remains for reference only and is no longer invoked by
 * the npm scripts. Use `hawp work normalize` or `npm run work:normalize` instead.
 *
 * Backlog Upgrade CLI Entry Point — delegates to TypeScript implementation.
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
