import assert from "node:assert/strict";
import { mkdtempSync, mkdirSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";

import { parseBacklog } from "../orchestrate";
import {
  extractIdFromFilename,
  extractShortUuid,
  idsMatch,
} from "../validations/id-parser";

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

test("parseBacklog falls back to the UUID cell when Legacy ID is a placeholder", () => {
  withTempDir((dir) => {
    const backlogPath = join(dir, "BACKLOG.md");
    const uuid = "361fb08e-6457-4ed5-80bd-76337b6f0e89";

    writeFileSync(
      backlogPath,
      `# Backlog

## Active Work

| UUID | Legacy ID | Type | Title | Status | Owner | Plan File | Updated |
| --- | --- | --- | --- | --- | --- | --- | --- |
| \`${uuid}\` | — | task | UUID-native item | in-progress | agent | [plan](active/${uuid}.md) | 2026-07-03 |

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
    assert.equal(parsed?.active[0]?.id, uuid);
    assert.equal(parsed?.active[0]?.type, "task");
  });
});

test("extractIdFromFilename recognizes legacy, date-prefixed, and UUID formats", () => {
  assert.equal(extractIdFromFilename("TASK-012"), "TASK-012");
  assert.equal(extractIdFromFilename("BUG-063-some-title"), "BUG-063");
  assert.equal(extractIdFromFilename("2026-04-29-BUG-001-title"), "BUG-001");
  assert.equal(
    extractIdFromFilename("361fb08e-6457-4ed5-80bd-76337b6f0e89"),
    "361fb08e-6457-4ed5-80bd-76337b6f0e89",
  );
  assert.equal(
    extractIdFromFilename("361FB08E-6457-4ED5-80BD-76337B6F0E89.md"),
    "361fb08e-6457-4ed5-80bd-76337b6f0e89",
  );
  // Bare 8-char alphanumeric (folder-per-item dir names like active/b7e2a4f9/)
  assert.equal(extractIdFromFilename("361fb08e"), "361fb08e");
  assert.equal(extractIdFromFilename("no-id-here"), null);
});

test("idsMatch and extractShortUuid handle short-UUID prefix matching", () => {
  const full = "0e1c4afa-9668-4d61-b5b6-1e27be42ca23";
  assert.equal(extractShortUuid("0e1c4afa"), "0e1c4afa");
  assert.equal(extractShortUuid("0E1C4AFA"), "0e1c4afa");
  assert.equal(extractShortUuid("0e1c4afa-9668"), null);
  assert.equal(extractShortUuid("TASK-013"), null);

  assert.ok(idsMatch("0e1c4afa", full));
  assert.ok(idsMatch(full, "0e1c4afa"));
  assert.ok(idsMatch(full, full.toUpperCase()));
  assert.ok(idsMatch("TASK-013", "TASK-013"));
  assert.ok(!idsMatch("0e1c4afa", "361fb08e-6457-4ed5-80bd-76337b6f0e89"));
  assert.ok(!idsMatch("TASK-013", full));
  // Short-vs-short must be exact, not prefix
  assert.ok(!idsMatch("0e1c4afa", "0e1c4afb"));
});

test("parseBacklog accepts a short-UUID cell as the row ID", () => {
  withTempDir((dir) => {
    const backlogPath = join(dir, "BACKLOG.md");

    writeFileSync(
      backlogPath,
      `# Backlog

## Active Work

| UUID | Legacy ID | Type | Title | Status | Owner | Plan File | Updated |
| --- | --- | --- | --- | --- | --- | --- | --- |
| \`0E1C4AFA\` | — | improvement | short display row | in-progress | agent | [plan](active/x.md) | 2026-07-03 |

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
    assert.equal(parsed?.active[0]?.id, "0e1c4afa");
  });
});
