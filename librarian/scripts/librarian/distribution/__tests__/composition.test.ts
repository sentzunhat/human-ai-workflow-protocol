import assert from "node:assert/strict";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import test from "node:test";

import {
  computeExpectedOutputs,
  findRepoRoot,
  normalizeForCompare,
} from "../shared/composition";

const currentDir = dirname(fileURLToPath(import.meta.url));

test("findRepoRoot resolves the repository root from inside librarian", () => {
  const root = findRepoRoot(currentDir);
  assert.ok(root, "expected a repo root");
  assert.equal(findRepoRoot("/"), null);
});

test("computeExpectedOutputs produces every active provider/operation/ref variant", () => {
  const root = findRepoRoot(currentDir);
  assert.ok(root);

  const outputs = computeExpectedOutputs(root);

  const providers = ["claude", "codex", "github", "cursor", "continue"];
  const operations = ["install", "update"];
  const refs = ["main", "dev"];
  assert.equal(outputs.length, providers.length * operations.length * refs.length);

  for (const provider of providers) {
    for (const operation of operations) {
      for (const ref of refs) {
        const expectedSuffix = join(
          "distribution",
          "generated",
          provider,
          operation,
          `${ref}.md`,
        );
        const output = outputs.find((entry) =>
          entry.outputPath.endsWith(expectedSuffix),
        );
        assert.ok(output, `missing output for ${provider}/${operation}/${ref}`);
        assert.match(output.content, new RegExp(`REF="${ref}"`));
        assert.match(output.content, new RegExp(`PROVIDER="${provider}"`));
        assert.match(output.content, /set -euo pipefail/);
        assert.ok(output.content.endsWith("\n"), "trailing newline expected");
      }
    }
  }
});

test("normalizeForCompare converts CRLF to LF", () => {
  assert.equal(normalizeForCompare("a\r\nb\r\n"), "a\nb\n");
  assert.equal(normalizeForCompare("a\nb"), "a\nb");
});
