import assert from "node:assert/strict";
import { mkdtempSync, mkdirSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";

import { runDetection } from "../detection";

const withTempRepo = (fn: (repoRoot: string) => void): void => {
  const repoRoot = mkdtempSync(join(tmpdir(), "hawp-detect-"));
  try {
    fn(repoRoot);
  } finally {
    rmSync(repoRoot, { recursive: true, force: true });
  }
};

const ensureWorkDirs = (
  repoRoot: string,
  options: { active?: boolean; parked?: boolean; closed?: boolean } = {},
): void => {
  const { active = true, parked = true, closed = true } = options;
  const workRoot = join(repoRoot, ".hawp", "work");

  if (active) {
    mkdirSync(join(workRoot, "active"), { recursive: true });
  }
  if (parked) {
    mkdirSync(join(workRoot, "parked"), { recursive: true });
  }
  if (closed) {
    mkdirSync(join(workRoot, "closed", "2026", "05", "14"), {
      recursive: true,
    });
  }
};

const writeBacklog = (
  repoRoot: string,
  bodyRows: string,
  options: { active?: boolean; parked?: boolean; closed?: boolean } = {},
): void => {
  const workRoot = join(repoRoot, ".hawp", "work");
  ensureWorkDirs(repoRoot, options);

  const content = `# Backlog

## Active Work

| ID | Type | Title | Status | Owner | Plan File | Updated |
| --- | --- | --- | --- | --- | --- | --- |
${bodyRows}

## Blocked / Parked

| ID | Type | Title | Reason | Detail | Updated |
| --- | --- | --- | --- | --- | --- |

## Recently Closed

| ID | Type | Title | Closed | Detail |
| --- | --- | --- | --- | --- |
`;

  writeFileSync(join(workRoot, "BACKLOG.md"), content, "utf-8");
};

const writeCustomBacklog = (repoRoot: string, content: string): void => {
  const workRoot = join(repoRoot, ".hawp", "work");
  mkdirSync(workRoot, { recursive: true });
  writeFileSync(join(workRoot, "BACKLOG.md"), content, "utf-8");
};

test("runDetection emits A1 for missing type inferred from ID", () => {
  withTempRepo((repoRoot) => {
    writeBacklog(
      repoRoot,
      "| TASK-100 |  | sample title | in-progress | agent | [plan](active/TASK-100.md) | 2026-05-14 |",
    );

    writeFileSync(
      join(repoRoot, ".hawp", "work", "active", "TASK-100.md"),
      "# TASK-100\n\n**Backlog ID:** TASK-100\n",
      "utf-8",
    );

    const { report } = runDetection(repoRoot);
    const operation = report.plan.operations.find(
      (item) => item.ruleId === "A1",
    );

    assert.ok(operation, "expected A1 operation");
    assert.equal(operation?.itemId, "TASK-100");
    assert.equal(operation?.safety, "safe");
  });
});

test("runDetection emits B1 when type cannot be inferred", () => {
  withTempRepo((repoRoot) => {
    writeBacklog(
      repoRoot,
      "| WORK-200 |  | sample title | in-progress | agent | [plan](active/WORK-200.md) | 2026-05-14 |",
    );

    writeFileSync(
      join(repoRoot, ".hawp", "work", "active", "WORK-200.md"),
      "# WORK-200\n\n**Backlog ID:** WORK-200\n",
      "utf-8",
    );

    const { report } = runDetection(repoRoot);
    const operation = report.plan.operations.find(
      (item) => item.ruleId === "B1",
    );

    assert.ok(operation, "expected B1 operation");
    assert.equal(operation?.safety, "blocked");
  });
});

test("runDetection emits B3 for duplicate plan candidates", () => {
  withTempRepo((repoRoot) => {
    writeBacklog(
      repoRoot,
      "| TASK-300 | task | sample title | in-progress | agent | [plan](active/TASK-300.md) | 2026-05-14 |",
    );

    writeFileSync(
      join(repoRoot, ".hawp", "work", "active", "TASK-300.md"),
      "# active\n\n**Backlog ID:** TASK-300\n",
      "utf-8",
    );

    writeFileSync(
      join(
        repoRoot,
        ".hawp",
        "work",
        "closed",
        "2026",
        "05",
        "14",
        "TASK-300.md",
      ),
      "# closed\n\n**Backlog ID:** TASK-300\n",
      "utf-8",
    );

    const { report } = runDetection(repoRoot);
    const operation = report.plan.operations.find(
      (item) => item.ruleId === "B3",
    );

    assert.ok(operation, "expected B3 operation");
    assert.equal(operation?.safety, "blocked");
    assert.ok((operation?.blocked?.candidates.length ?? 0) >= 2);

    const syncItem = report.syncPlan.find((item) => item.itemId === "TASK-300");
    assert.ok(syncItem, "expected sync plan entry for duplicate candidates");
    assert.match(
      syncItem?.canonicalPlan ?? "",
      /closed\/2026\/05\/14\/TASK-300\.md/,
    );
    assert.equal(
      syncItem?.duplicatePlans.includes(".hawp/work/active/TASK-300.md"),
      true,
    );
  });
});

test("runDetection emits B2 when linked plan path is missing", () => {
  withTempRepo((repoRoot) => {
    writeBacklog(
      repoRoot,
      "| TASK-410 | task | missing link target | in-progress | agent | [plan](active/TASK-410.md) | 2026-05-14 |",
    );

    const { report } = runDetection(repoRoot);
    const operation = report.plan.operations.find(
      (item) => item.ruleId === "B2",
    );

    assert.ok(operation, "expected B2 operation");
    assert.equal(operation?.itemId, "TASK-410");
    assert.match(
      operation?.description ?? "",
      /Referenced plan path is missing/,
    );
  });
});

test("runDetection emits B4 for recently closed plan with unproven marker", () => {
  withTempRepo((repoRoot) => {
    writeCustomBacklog(
      repoRoot,
      `# Backlog

## Active Work

| ID | Type | Title | Status | Owner | Plan File | Updated |
| --- | --- | --- | --- | --- | --- | --- |

## Blocked / Parked

| ID | Type | Title | Reason | Detail | Updated |
| --- | --- | --- | --- | --- | --- |

## Recently Closed

| ID | Type | Title | Closed | Detail |
| --- | --- | --- | --- | --- |
| TASK-420 | task | closed item | 2026-05-14 | [plan](closed/2026/05/14/TASK-420.md) |
`,
    );
    ensureWorkDirs(repoRoot);

    writeFileSync(
      join(
        repoRoot,
        ".hawp",
        "work",
        "closed",
        "2026",
        "05",
        "14",
        "TASK-420.md",
      ),
      "# Task: sample\n\n## Outcome\n\nDone.\n\n## Verification\n\n- [ ] NOT YET VERIFIED\n\n## Close Checklist\n\n- [ ] Example\n",
      "utf-8",
    );

    const { report } = runDetection(repoRoot);
    const operation = report.plan.operations.find(
      (item) => item.ruleId === "B4",
    );

    assert.ok(operation, "expected B4 operation");
    assert.equal(operation?.itemId, "TASK-420");
    assert.match(
      operation?.description ?? "",
      /unproven verification markers/i,
    );
  });
});

test("runDetection emits B5 when required work directories are missing", () => {
  withTempRepo((repoRoot) => {
    writeBacklog(
      repoRoot,
      "| TASK-430 | task | sample title | in-progress | agent | [plan](active/TASK-430.md) | 2026-05-14 |",
      { active: true, parked: false, closed: true },
    );

    writeFileSync(
      join(repoRoot, ".hawp", "work", "active", "TASK-430.md"),
      "# TASK-430\n\n**Backlog ID:** TASK-430\n",
      "utf-8",
    );

    const { report } = runDetection(repoRoot);
    const operation = report.plan.operations.find(
      (item) => item.ruleId === "B5",
    );

    assert.ok(operation, "expected B5 operation");
    assert.equal(operation?.itemId, "WORKSPACE");
    assert.ok(operation?.blocked?.candidates.includes("parked"));
  });
});

test("runDetection emits A4 and A5 for recently closed plan missing sections", () => {
  withTempRepo((repoRoot) => {
    writeCustomBacklog(
      repoRoot,
      `# Backlog

## Active Work

| ID | Type | Title | Status | Owner | Plan File | Updated |
| --- | --- | --- | --- | --- | --- | --- |

## Blocked / Parked

| ID | Type | Title | Reason | Detail | Updated |
| --- | --- | --- | --- | --- | --- |

## Recently Closed

| ID | Type | Title | Closed | Detail |
| --- | --- | --- | --- | --- |
| TASK-435 | task | closed item | 2026-05-14 | [plan](closed/2026/05/14/TASK-435.md) |
`,
    );
    ensureWorkDirs(repoRoot);

    writeFileSync(
      join(
        repoRoot,
        ".hawp",
        "work",
        "closed",
        "2026",
        "05",
        "14",
        "TASK-435.md",
      ),
      "# Task: sample\n\n## Verification\n\nVerified.\n",
      "utf-8",
    );

    const { report } = runDetection(repoRoot);
    const hasA4 = report.plan.operations.some(
      (item) => item.ruleId === "A4" && item.itemId === "TASK-435",
    );
    const hasA5 = report.plan.operations.some(
      (item) => item.ruleId === "A5" && item.itemId === "TASK-435",
    );

    assert.equal(hasA4, true, "expected A4 for missing Outcome heading");
    assert.equal(
      hasA5,
      true,
      "expected A5 for missing Close Checklist scaffolding",
    );
  });
});

test("runDetection handles edge-case table layout under non-standard section", () => {
  withTempRepo((repoRoot) => {
    writeCustomBacklog(
      repoRoot,
      `# Backlog

## Misc Lane

| ID | Type | Title | Status | Owner | Plan File | Updated |
| --- | --- | --- | --- | --- | --- | --- |
| TASK-440 | task | custom lane row | done | agent | [plan](active/TASK-440.md) | 2026/05/14 |

## Active Work

| ID | Type | Title | Status | Owner | Plan File | Updated |
| --- | --- | --- | --- | --- | --- | --- |
`,
    );
    ensureWorkDirs(repoRoot);

    const { report } = runDetection(repoRoot);
    const hasA2 = report.plan.operations.some(
      (item) => item.itemId === "TASK-440" && item.ruleId === "A2",
    );
    const hasA3 = report.plan.operations.some(
      (item) => item.itemId === "TASK-440" && item.ruleId === "A3",
    );

    assert.equal(hasA2, true, "expected A2 on non-ISO date in custom lane");
    assert.equal(hasA3, false, "did not expect A3 for canonical ID format");
  });
});

test("runDetection supports active rows with UUID and Legacy ID columns", () => {
  withTempRepo((repoRoot) => {
    writeCustomBacklog(
      repoRoot,
      `# Backlog

## Active Work

| UUID | Legacy ID | Type | Title | Status | Owner | Plan File | Updated |
| --- | --- | --- | --- | --- | --- | --- | --- |
| \`f1c9b3a2\` | TASK-450 |  | uuid-backed active item | analyzing | agent | [plan](active/TASK-450.md) | 2026-05-14 |

## Blocked / Parked

| ID | Type | Title | Reason | Detail | Updated |
| --- | --- | --- | --- | --- | --- |

## Recently Closed

| ID | Type | Title | Closed | Detail |
| --- | --- | --- | --- | --- |
`,
    );
    ensureWorkDirs(repoRoot);

    writeFileSync(
      join(repoRoot, ".hawp", "work", "active", "TASK-450.md"),
      "# Task: sample\n\n**Backlog ID:** TASK-450\n",
      "utf-8",
    );

    const { report } = runDetection(repoRoot);
    const operation = report.plan.operations.find(
      (item) => item.ruleId === "A1",
    );

    assert.ok(operation, "expected A1 operation");
    assert.equal(operation?.itemId, "TASK-450");
  });
});

test("runDetection skips A4/A5/A7 for legacy closed files (pre-2026-05-10)", () => {
  withTempRepo((repoRoot) => {
    writeCustomBacklog(
      repoRoot,
      `# Backlog

## Active Work

| ID | Type | Title | Status | Owner | Plan File | Updated |
| --- | --- | --- | --- | --- | --- | --- |

## Blocked / Parked

| ID | Type | Title | Reason | Detail | Updated |
| --- | --- | --- | --- | --- | --- |

## Recently Closed

| ID | Type | Title | Closed | Detail |
| --- | --- | --- | --- | --- |
| TASK-001 | task | legacy closed item | 2026-04-30 | [plan](closed/2026/04/30/TASK-001.md) |
`,
    );
    ensureWorkDirs(repoRoot, { active: true, parked: true, closed: true });

    // Create legacy closed file without Outcome/Close Checklist and with stale reference
    mkdirSync(join(repoRoot, ".hawp", "work", "closed", "2026", "04", "30"), {
      recursive: true,
    });
    writeFileSync(
      join(
        repoRoot,
        ".hawp",
        "work",
        "closed",
        "2026",
        "04",
        "30",
        "TASK-001.md",
      ),
      "# Task: legacy\n\n## Verification\n\nDone.\n\nSee [reference](core/distribution/generated/file.md).\n",
      "utf-8",
    );

    const { report } = runDetection(repoRoot);
    const hasA4Legacy = report.plan.operations.some(
      (item) => item.ruleId === "A4" && item.itemId === "TASK-001",
    );
    const hasA5Legacy = report.plan.operations.some(
      (item) => item.ruleId === "A5" && item.itemId === "TASK-001",
    );
    const hasA7Legacy = report.plan.operations.some(
      (item) => item.ruleId === "A7" && item.itemId === "TASK-001",
    );

    assert.equal(hasA4Legacy, false, "did not expect A4 on legacy file");
    assert.equal(hasA5Legacy, false, "did not expect A5 on legacy file");
    assert.equal(
      hasA7Legacy,
      false,
      "did not expect A7 on legacy file (stale ref ignored)",
    );
  });
});

test("runDetection emits B7 for UUID-native closed plan with ambiguous verification claim", () => {
  withTempRepo((repoRoot) => {
    const uuid = "ddeb9eb3-6dba-4a14-8c24-46b92a225bd3";
    writeCustomBacklog(
      repoRoot,
      `# Backlog

## Active Work

| UUID | Legacy ID | Type | Title | Status | Owner | Plan File | Updated |
| --- | --- | --- | --- | --- | --- | --- | --- |

## Blocked / Parked

| ID | Type | Title | Reason | Detail | Updated |
| --- | --- | --- | --- | --- | --- |

## Recently Closed

| UUID | Legacy ID | Type | Title | Closed | Detail |
| --- | --- | --- | --- | --- | --- |
| \`ddeb9eb3\` | — | improvement | uuid closed item | 2026-07-03 | [plan](closed/2026/07/03/${uuid}.md) |
`,
    );
    ensureWorkDirs(repoRoot);
    mkdirSync(join(repoRoot, ".hawp", "work", "closed", "2026", "07", "03"), {
      recursive: true,
    });

    writeFileSync(
      join(
        repoRoot,
        ".hawp",
        "work",
        "closed",
        "2026",
        "07",
        "03",
        `${uuid}.md`,
      ),
      [
        "# UUID-native closed item",
        "",
        "**Backlog ID (Legacy):** — (UUID-native item)",
        `**UUID:** \`${uuid}\``,
        "",
        "## Outcome",
        "",
        "Done.",
        "",
        "## Verification",
        "",
        "- [x] Smoke test completed",
        "",
        "## Close Checklist",
        "",
        "- [x] Closed",
        "",
      ].join("\n"),
      "utf-8",
    );

    const { report } = runDetection(repoRoot);
    const operation = report.plan.operations.find(
      (item) => item.ruleId === "B7",
    );

    assert.ok(operation, "expected B7 operation");
    assert.equal(operation?.itemId, "ddeb9eb3");
    assert.match(operation?.description ?? "", /ambiguous verification/i);
  });
});
