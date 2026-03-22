#!/usr/bin/env node

/**
 * Executable boundary for the HAWP workflow validator.
 * Delegates to script.ts; only this file sets the process exit code.
 */

import { runWorkflowValidation } from "./script";

process.exitCode = runWorkflowValidation(process.argv.slice(2));
