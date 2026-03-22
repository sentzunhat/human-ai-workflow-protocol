#!/usr/bin/env node

import { runKitValidate } from "./script";

process.exitCode = runKitValidate(process.argv.slice(2));
