import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import test from "node:test";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

import { parseArgs } from "../cli";

const currentFile = fileURLToPath(import.meta.url);
const scriptPath = join(dirname(currentFile), "..", "index.ts");

test("parseArgs separates own flags from workflow passthrough args", () => {
  const args = parseArgs([
    "--strict-warnings",
    "--work-root",
    "./.hawp/work",
    "--debug-closed-task",
  ]);

  assert.equal(args.help, false);
  assert.equal(args.strictWarnings, true);
  assert.deepEqual(args.workflowArgs, [
    "--work-root",
    "./.hawp/work",
    "--debug-closed-task",
  ]);

  const helpArgs = parseArgs(["-h"]);
  assert.equal(helpArgs.help, true);
  assert.deepEqual(helpArgs.workflowArgs, []);
});

test("backlog validate help prints usage", () => {
  const result = spawnSync("npx", ["tsx", scriptPath, "--help"], {
    encoding: "utf-8",
  });

  assert.equal(result.status, 0);
  assert.match(result.stdout, /hawp backlog validate/);
  assert.match(result.stdout, /\.\/\.hawp\/bin\/hawp backlog validate/);
});
