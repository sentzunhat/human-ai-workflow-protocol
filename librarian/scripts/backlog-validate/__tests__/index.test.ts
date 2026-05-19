import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import test from "node:test";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const currentFile = fileURLToPath(import.meta.url);
const scriptPath = join(dirname(currentFile), "..", "index.ts");

test("backlog validate help prints usage", () => {
  const result = spawnSync("npx", ["tsx", scriptPath, "--help"], {
    encoding: "utf-8",
  });

  assert.equal(result.status, 0);
  assert.match(result.stdout, /hawp backlog validate/);
  assert.match(result.stdout, /\.\/\.hawp\/bin\/hawp backlog validate/);
});
