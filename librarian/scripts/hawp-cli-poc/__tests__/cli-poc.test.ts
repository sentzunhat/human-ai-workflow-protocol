import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import { resolve } from "node:path";
import test from "node:test";

const cliPath = resolve("scripts", "hawp-cli-poc", "bin", "hawp-poc.ts");

test("hawp-poc prints top-level help", () => {
  const result = spawnSync("npx", ["tsx", cliPath, "--help"], {
    cwd: resolve("."),
    encoding: "utf8",
  });

  assert.equal(result.status, 0);
  assert.match(result.stdout, /hawp-poc/);
  assert.match(result.stdout, /kit validate/);
});

test("hawp-poc routes kit validate through the command registry", () => {
  const result = spawnSync("npx", ["tsx", cliPath, "kit", "validate"], {
    cwd: resolve(".."),
    encoding: "utf8",
  });

  assert.equal(result.status, 0);
  assert.match(result.stdout, /kit:validate/);
  assert.match(result.stdout, /checks passed/);
});

test("hawp-poc returns 2 for unknown commands", () => {
  const result = spawnSync("npx", ["tsx", cliPath, "bogus"], {
    cwd: resolve("."),
    encoding: "utf8",
  });

  assert.equal(result.status, 2);
  assert.match(result.stderr, /Unknown command: bogus/);
});
