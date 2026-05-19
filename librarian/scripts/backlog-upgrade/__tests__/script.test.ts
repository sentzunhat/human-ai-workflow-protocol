import assert from "node:assert/strict";
import { execSync } from "node:child_process";
import {
  existsSync,
  mkdtempSync,
  mkdirSync,
  readFileSync,
  rmSync,
  writeFileSync,
} from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";

import { ExitCode, Mode, OutputFormat } from "../models/index.js";
import { runBacklogUpgradeScript } from "../script.js";

const withTempRepo = async (
  fn: (repoRoot: string) => void | Promise<void>,
): Promise<void> => {
  const repoRoot = mkdtempSync(join(tmpdir(), "hawp-script-"));
  try {
    await fn(repoRoot);
  } finally {
    rmSync(repoRoot, { recursive: true, force: true });
  }
};

const seedMinimalBacklog = (repoRoot: string): void => {
  const workRoot = join(repoRoot, ".hawp", "work");
  mkdirSync(join(workRoot, "active"), { recursive: true });
  mkdirSync(join(workRoot, "parked"), { recursive: true });
  mkdirSync(join(workRoot, "closed", "2026", "05", "14"), {
    recursive: true,
  });

  const backlogContent = `# Backlog

## Active Work

| ID | Type | Title | Status | Owner | Plan File | Updated |
| --- | --- | --- | --- | --- | --- | --- |
| TASK-100 | task | sample title | in-progress | agent | [plan](active/TASK-100.md) | 2026-05-14 |

## Blocked / Parked

| ID | Type | Title | Reason | Detail | Updated |
| --- | --- | --- | --- | --- | --- |

## Recently Closed

| ID | Type | Title | Closed | Detail |
| --- | --- | --- | --- | --- |
`;

  writeFileSync(join(workRoot, "BACKLOG.md"), backlogContent, "utf-8");
  writeFileSync(
    join(workRoot, "active", "TASK-100.md"),
    "# Task: sample\n\n## Verification\n\n- [ ] NOT YET VERIFIED\n",
    "utf-8",
  );
};

test("runBacklogUpgradeScript returns rendered text in dry-run mode", async () => {
  await withTempRepo(async (repoRoot) => {
    seedMinimalBacklog(repoRoot);

    const result = await runBacklogUpgradeScript({
      mode: Mode.DryRun,
      validate: false,
      format: OutputFormat.Text,
      forceDirty: false,
      verbose: false,
      repoRoot,
    });

    assert.equal(result.exitCode, ExitCode.Success);
    assert.equal(result.stderrLines.length, 0);
    assert.match(result.stdoutText, /HAWP Backlog Upgrade Dry-Run Report/);
  });
});

test("runBacklogUpgradeScript writes output and exported plan files", async () => {
  await withTempRepo(async (repoRoot) => {
    seedMinimalBacklog(repoRoot);
    const outputPath = join(repoRoot, "tmp", "report.txt");
    const exportPlanPath = join(repoRoot, "tmp", "plan.json");

    const result = await runBacklogUpgradeScript({
      mode: Mode.DryRun,
      validate: false,
      format: OutputFormat.Text,
      output: outputPath,
      exportPlan: exportPlanPath,
      forceDirty: false,
      verbose: false,
      repoRoot,
    });

    assert.equal(result.exitCode, ExitCode.Success);
    assert.equal(result.stdoutText, "");
    assert.ok(result.notices.some((line) => line.includes("Report written:")));
    assert.ok(result.notices.some((line) => line.includes("Plan exported:")));
    assert.equal(existsSync(outputPath), true);
    assert.equal(existsSync(exportPlanPath), true);

    const reportText = readFileSync(outputPath, "utf-8");
    const planText = readFileSync(exportPlanPath, "utf-8");
    assert.match(reportText, /HAWP Backlog Upgrade Dry-Run Report/);
    assert.match(planText, /"planId"/);
  });
});

test("runBacklogUpgradeScript blocks apply mode on a dirty working tree", async () => {
  const result = await runBacklogUpgradeScript({
    mode: Mode.Apply,
    validate: false,
    format: OutputFormat.Text,
    forceDirty: false,
    verbose: false,
    repoRoot: process.cwd(),
  });

  assert.equal(result.exitCode, ExitCode.UsageError);
  assert.equal(result.stdoutText, "");
  assert.equal(
    result.stderrLines[0],
    "Error: apply mode requires a clean working tree. Re-run with --force-dirty to override.",
  );
});

test("runBacklogUpgradeScript scaffolds closed records in apply mode", async () => {
  await withTempRepo(async (repoRoot) => {
    execSync("git init -q", { cwd: repoRoot });
    seedMinimalBacklog(repoRoot);

    const closedDir = join(
      repoRoot,
      ".hawp",
      "work",
      "closed",
      "2026",
      "05",
      "14",
    );
    const closedPath = join(closedDir, "TASK-700.md");
    writeFileSync(
      closedPath,
      "# Task: scaffold me\n\nLegacy content without normalized close sections.\n",
      "utf-8",
    );

    const result = await runBacklogUpgradeScript({
      mode: Mode.Apply,
      validate: false,
      format: OutputFormat.Text,
      forceDirty: true,
      verbose: false,
      repoRoot,
    });

    assert.equal(result.exitCode, ExitCode.Success);
    assert.ok(
      result.notices.some((notice) =>
        notice.includes("Applied closed-record normalization to 1 file(s)."),
      ),
    );

    const updated = readFileSync(closedPath, "utf-8");
    assert.match(updated, /\*\*Backlog ID:\*\* TASK-700/);
    assert.match(updated, /^## Outcome$/m);
    assert.match(updated, /^## Verification$/m);
    assert.match(updated, /^## Close Checklist$/m);
  });
});

test("runBacklogUpgradeScript reconciles dated closed records to matching folders", async () => {
  await withTempRepo(async (repoRoot) => {
    execSync("git init -q", { cwd: repoRoot });
    seedMinimalBacklog(repoRoot);

    const wrongFolderPath = join(
      repoRoot,
      ".hawp",
      "work",
      "closed",
      "2026",
      "05",
      "14",
      "2026-05-10-TASK-800.md",
    );
    writeFileSync(
      wrongFolderPath,
      "# Task: reconcile me\n\nLegacy content without normalized close sections.\n",
      "utf-8",
    );

    const targetPath = join(
      repoRoot,
      ".hawp",
      "work",
      "closed",
      "2026",
      "05",
      "10",
      "2026-05-10-TASK-800.md",
    );

    const result = await runBacklogUpgradeScript({
      mode: Mode.Apply,
      validate: false,
      format: OutputFormat.Text,
      forceDirty: true,
      verbose: false,
      repoRoot,
    });

    assert.equal(result.exitCode, ExitCode.Success);
    assert.equal(existsSync(wrongFolderPath), false);
    assert.equal(existsSync(targetPath), true);

    const updated = readFileSync(targetPath, "utf-8");
    assert.match(updated, /\*\*Backlog ID:\*\* TASK-800/);
    assert.match(updated, /^## Outcome$/m);
    assert.match(updated, /^## Verification$/m);
    assert.match(updated, /^## Close Checklist$/m);
  });
});

test("runBacklogUpgradeScript emits validation and mixed-id notices in dry-run", async () => {
  await withTempRepo(async (repoRoot) => {
    const workRoot = join(repoRoot, ".hawp", "work");
    mkdirSync(join(workRoot, "active"), { recursive: true });
    mkdirSync(join(workRoot, "parked"), { recursive: true });
    mkdirSync(join(workRoot, "closed", "2026", "05", "14"), {
      recursive: true,
    });

    writeFileSync(
      join(workRoot, "BACKLOG.md"),
      `# Backlog

## Active Work

| UUID | Legacy ID | Type | Title | Status | Owner | Plan File | Updated |
| --- | --- | --- | --- | --- | --- | --- | --- |
| \`f1c9b3a2\` | TASK-100 | task | sample title | in-progress | agent | [plan](active/TASK-100.md) | 2026-05-14 |

## Blocked / Parked

| ID | Type | Title | Reason | Detail | Updated |
| --- | --- | --- | --- | --- | --- |

## Recently Closed

| ID | Type | Title | Closed | Detail |
| --- | --- | --- | --- | --- |
`,
      "utf-8",
    );

    writeFileSync(
      join(workRoot, "active", "TASK-100.md"),
      "# Task: sample\n\n**Backlog ID:** TASK-100\n",
      "utf-8",
    );

    const result = await runBacklogUpgradeScript({
      mode: Mode.DryRun,
      validate: true,
      format: OutputFormat.Text,
      forceDirty: false,
      verbose: false,
      repoRoot,
    });

    assert.equal(result.exitCode, ExitCode.Success);
    assert.equal(
      result.notices.some((notice) => notice.includes("Validation summary:")),
      true,
    );
    assert.equal(
      result.notices.some((notice) =>
        notice.includes("Mixed UUID/legacy backlog detected"),
      ),
      true,
    );
  });
});

test("runBacklogUpgradeScript renders concrete drift sync/apply plan", async () => {
  await withTempRepo(async (repoRoot) => {
    const workRoot = join(repoRoot, ".hawp", "work");
    mkdirSync(join(workRoot, "active"), { recursive: true });
    mkdirSync(join(workRoot, "parked"), { recursive: true });
    mkdirSync(join(workRoot, "closed", "2026", "05", "14"), {
      recursive: true,
    });

    writeFileSync(
      join(workRoot, "BACKLOG.md"),
      `# Backlog

## Active Work

| ID | Type | Title | Status | Owner | Plan File | Updated |
| --- | --- | --- | --- | --- | --- | --- |
| TASK-300 | task | duplicate active plan | in-progress | agent | [plan](active/TASK-300.md) | 2026-05-14 |

## Blocked / Parked

| ID | Type | Title | Reason | Detail | Updated |
| --- | --- | --- | --- | --- | --- |

## Recently Closed

| ID | Type | Title | Closed | Detail |
| --- | --- | --- | --- | --- |
`,
      "utf-8",
    );

    writeFileSync(
      join(workRoot, "active", "TASK-300.md"),
      "# active\n\n**Backlog ID:** TASK-300\n",
      "utf-8",
    );
    writeFileSync(
      join(workRoot, "parked", "TASK-300.md"),
      "# parked\n\n**Backlog ID:** TASK-300\n",
      "utf-8",
    );
    writeFileSync(
      join(workRoot, "closed", "2026", "05", "14", "TASK-300.md"),
      "# closed\n\n**Backlog ID:** TASK-300\n",
      "utf-8",
    );

    const result = await runBacklogUpgradeScript({
      mode: Mode.DryRun,
      validate: false,
      format: OutputFormat.Text,
      forceDirty: false,
      verbose: false,
      repoRoot,
    });

    assert.equal(result.exitCode, ExitCode.Success);
    assert.match(result.stdoutText, /Drift Sync\/Apply Plan/);
    assert.match(
      result.stdoutText,
      /TASK-300: canonical=.*closed\/2026\/05\/14\/TASK-300\.md/,
    );
    assert.match(
      result.stdoutText,
      /Remove stale active copy: .*active\/TASK-300\.md/,
    );
    assert.match(
      result.stdoutText,
      /Remove stale parked copy: .*parked\/TASK-300\.md/,
    );
  });
});
