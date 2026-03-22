#!/usr/bin/env node

import { runKitNormalize } from "./script";

process.exitCode = runKitNormalize(process.argv.slice(2));
