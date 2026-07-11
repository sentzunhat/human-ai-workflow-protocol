import assert from "node:assert/strict";
import { dirname } from "node:path";
import { fileURLToPath } from "node:url";
import test from "node:test";

import { findRepoRoot } from "../../../distribution/shared/composition";
import { GENERATED_BANNER, MATERIALIZATION_TARGETS } from "../composition";
import { computeMaterializedOutputs } from "../render";

const currentDir = dirname(fileURLToPath(import.meta.url));

test("computeMaterializedOutputs renders every materialization target", () => {
  const root = findRepoRoot(currentDir);
  assert.ok(root);

  const outputs = computeMaterializedOutputs(root);
  assert.equal(outputs.length, MATERIALIZATION_TARGETS.length);

  for (const output of outputs) {
    assert.ok(
      output.content.startsWith("---\n"),
      `frontmatter expected in ${output.outputPath}`,
    );
    assert.ok(
      output.content.includes(GENERATED_BANNER.trimEnd()),
      `generated banner expected in ${output.outputPath}`,
    );
    assert.ok(output.content.endsWith("\n"));
  }
});

test("materialization targets cover generated rule providers", () => {
  const paths = MATERIALIZATION_TARGETS.map((target) => target.outputPath);
  assert.ok(paths.some((path) => path.includes(".claude/rules/")));
  assert.ok(paths.some((path) => path.includes(".cursor/rules/")));
  assert.ok(paths.some((path) => path.includes(".continue/rules/")));
  assert.ok(paths.some((path) => path.includes(".github/instructions/")));
});
