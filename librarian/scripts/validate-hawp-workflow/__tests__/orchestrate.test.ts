import assert from "node:assert/strict";
import { mkdtempSync, mkdirSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";

import { parseBacklog } from "../orchestrate";

const withTempDir = (fn: (dir: string) => void): void => {
  const dir = mkdtempSync(join(tmpdir(), "hawp-validate-"));
  try {
    fn(dir);
  } finally {
    rmSync(dir, { recursive: true, force: true });
  }
};

test("parseBacklog prefers Legacy ID when UUID columns are present", () => {
  withTempDir((dir) => {
    mkdirSync(dir, { recursive: true });
    const backlogPath = join(dir, "BACKLOG.md");

    writeFileSync(
      backlogPath,
      `# Backlog

## Active Work

| UUID | Legacy ID | Type | Title | Status | Owner | Plan File | Updated |
| --- | --- | --- | --- | --- | --- | --- | --- |
| \`f1c9b3a2\` | TASK-013 | improvement | UUID migration | analyzing | agent | [plan](active/TASK-013.md) | 2026-05-14 |

## Blocked / Parked

| ID | Type | Title | Reason | Detail | Updated |
| --- | --- | --- | --- | --- | --- |

## Recently Closed

| ID | Type | Title | Closed | Detail |
| --- | --- | --- | --- | --- |
`,
      "utf-8",
    );

    const parsed = parseBacklog(backlogPath);

    assert.ok(parsed);
    assert.equal(parsed?.active.length, 1);
    assert.equal(parsed?.active[0]?.id, "TASK-013");
    assert.equal(parsed?.active[0]?.type, "improvement");
    assert.equal(parsed?.active[0]?.detail, "[plan](active/TASK-013.md)");
  });
});
