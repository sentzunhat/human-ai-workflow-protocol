#!/usr/bin/env node

/**
 * @deprecated Superseded by the Go CLI (`hawp work validate`). This TypeScript
 * implementation remains for reference only and is no longer invoked by
 * the npm scripts. Use `hawp work validate` or `npm run work:validate` instead.
 *
 * Executable boundary for the HAWP workflow validator.
 * Delegates to script.ts; only this file sets the process exit code.
 */

import { runWorkflowValidation } from "./script";

process.exitCode = runWorkflowValidation(process.argv.slice(2));
