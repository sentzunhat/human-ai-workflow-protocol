/**
 * Combined kit/work validation logic for `hawp backlog validate`.
 * Runs the distribution validator and the workflow validator as child
 * processes and folds their results into one PASS/FAIL outcome.
 */

import { spawnSync } from "node:child_process";
import { join, resolve } from "node:path";

import { findBacklogRepoRoot } from "../../lib";
import type { BacklogValidateArgs } from "./cli";

const runCommand = (
  command: string,
  args: string[],
  cwd: string,
): { exitCode: number; stdout: string; stderr: string } => {
  const result = spawnSync(command, args, {
    cwd,
    encoding: "utf-8",
    stdio: "pipe",
  });

  return {
    exitCode: result.status ?? 1,
    stdout: result.stdout ?? "",
    stderr: result.stderr ?? "",
  };
};

const printSection = (
  title: string,
  content: string,
  stream: "stdout" | "stderr" = "stdout",
): void => {
  const target = stream === "stderr" ? process.stderr : process.stdout;
  target.write(`\n[${title}]\n`);
  target.write(
    content.trim().length > 0 ? `${content.trim()}\n` : "(no output)\n",
  );
};

const extractWarningsCount = (workflowStdout: string): number => {
  const match = workflowStdout.match(/! Warnings:\s+(\d+)/);
  if (!match?.[1]) {
    return 0;
  }

  const parsed = Number.parseInt(match[1], 10);
  return Number.isNaN(parsed) ? 0 : parsed;
};

export const runBacklogValidateScript = (
  args: BacklogValidateArgs,
): number => {
  const repoRoot = findBacklogRepoRoot(resolve(process.cwd()));

  process.stdout.write("HAWP Backlog Validate\n");
  process.stdout.write("=====================\n");
  process.stdout.write(`repo: ${repoRoot}\n`);

  const distribution = runCommand(
    "npx",
    [
      "tsx",
      join(
        repoRoot,
        "librarian",
        "scripts",
        "librarian",
        "distribution",
        "validate",
        "index.ts",
      ),
    ],
    repoRoot,
  );

  printSection("distribution:validate", distribution.stdout, "stdout");
  if (distribution.stderr.trim().length > 0) {
    printSection("distribution:validate stderr", distribution.stderr, "stderr");
  }

  const workflow = runCommand(
    "npx",
    [
      "tsx",
      join(
        repoRoot,
        "librarian",
        "scripts",
        "hawp",
        "work-validate",
        "index.ts",
      ),
      ...args.workflowArgs,
    ],
    repoRoot,
  );

  printSection("work:validate", workflow.stdout, "stdout");
  if (workflow.stderr.trim().length > 0) {
    printSection("work:validate stderr", workflow.stderr, "stderr");
  }

  const warningsCount = extractWarningsCount(workflow.stdout);
  const failedByWarnings = args.strictWarnings && warningsCount > 0;
  const failed =
    distribution.exitCode !== 0 || workflow.exitCode !== 0 || failedByWarnings;

  if (failed) {
    process.stderr.write("\nResult: FAIL\n");
    process.stderr.write("Next:\n");
    process.stderr.write(
      "- run npm --prefix librarian run work:normalize for drift diagnosis and repair\n",
    );
    process.stderr.write(
      "- run npm --prefix librarian run distribution:build if distribution drift is reported\n",
    );
    if (failedByWarnings) {
      process.stderr.write(
        `- strict warning mode enabled and workflow reported ${warningsCount} warning(s)\n`,
      );
    }
    return 1;
  }

  process.stdout.write("\nResult: PASS\n");
  process.stdout.write(
    "Both kit/work validators completed without FAIL checks.\n",
  );
  return 0;
};
