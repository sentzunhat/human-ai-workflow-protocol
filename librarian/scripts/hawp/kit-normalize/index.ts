#!/usr/bin/env node

/**
 * @deprecated Superseded by the Go CLI (`hawp kit normalize`). This TypeScript
 * implementation remains for reference only and is no longer invoked by
 * the npm scripts. Use `hawp kit normalize` or `npm run kit:normalize` instead.
 */

import { runKitNormalize } from "./script";

process.exitCode = runKitNormalize(process.argv.slice(2));
