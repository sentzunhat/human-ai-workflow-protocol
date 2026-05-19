#!/usr/bin/env node

import { spawnSync } from "node:child_process";
import { existsSync } from "node:fs";
import { dirname, join, resolve } from "node:path";

const findRepoRoot = (startDir: string): string => {
  let current = startDir;

  for (let depth = 0; depth < 12; depth += 1) {
    if (existsSync(join(current, ".hawp", "work", "BACKLOG.md"))) {
      return current;
    }

    const parent = dirname(current);
    if (parent === current) {
      break;
    }

    current = parent;
  }

  throw new Error(
    "Could not locate repo root containing .hawp/work/BACKLOG.md",
  );
};

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

const showHelp = (): void => {
  process.stdout.write(
    [
      "hawp backlog validate — validate kit/work drift in one command",
      "",
      "Usage:",
      "  ./.hawp/bin/hawp backlog validate [workflow-options]",
      "",
      "Behavior:",
      "  1) Runs distribution validation (.hawp/kit generated drift checks)",
      "  2) Runs workflow validation (.hawp/work backlog + plan checks)",
      "  3) Exits 1 if either validator fails",
      "",
      "Workflow options (passed through):",
      "  --hawp-root <path>",
      "  --work-root <path>",
      "  --debug-closed-task",
      "  --strict-warnings      Exit with code 1 if workflow warnings are present",
      "",
      "Examples:",
      "  ./.hawp/bin/hawp backlog validate",
      "  ./.hawp/bin/hawp backlog validate --work-root ./.hawp/work",
      "",
    ].join("\n"),
  );
};

const main = (): number => {
  const args = process.argv.slice(2);
  if (args.includes("--help") || args.includes("-h")) {
    showHelp();
    return 0;
  }

  const repoRoot = findRepoRoot(resolve(process.cwd()));
  const strictWarnings = args.includes("--strict-warnings");
  const workflowArgs = args.filter((arg) => arg !== "--strict-warnings");

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
        "validate-hawp-workflow",
        "index.ts",
      ),
      ...workflowArgs,
    ],
    repoRoot,
  );

  printSection("validate:workflow", workflow.stdout, "stdout");
  if (workflow.stderr.trim().length > 0) {
    printSection("validate:workflow stderr", workflow.stderr, "stderr");
  }

  const warningsCount = extractWarningsCount(workflow.stdout);
  const failedByWarnings = strictWarnings && warningsCount > 0;
  const failed =
    distribution.exitCode !== 0 || workflow.exitCode !== 0 || failedByWarnings;

  if (failed) {
    process.stderr.write("\nResult: FAIL\n");
    process.stderr.write("Next:\n");
    process.stderr.write(
      "- run ./.hawp/bin/hawp backlog upgrade --dry-run --validate for drift diagnosis\n",
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

process.exitCode = main();
