#!/usr/bin/env node

/**
 * @deprecated Superseded by the Go CLI (`hawp kit validate`). This TypeScript
 * implementation remains for reference only and is no longer invoked by
 * the npm scripts. Use `hawp kit validate` or `npm run kit:validate` instead.
 */

import { runKitValidate } from "./script";

process.exitCode = runKitValidate(process.argv.slice(2));
