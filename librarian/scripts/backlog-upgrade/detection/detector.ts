import { createHash } from "node:crypto";
import { existsSync } from "node:fs";
import { dirname, join, relative } from "node:path";

import {
  createBacklogFixPlan,
  createDetectionReport,
  type DetectionReport,
} from "../models";
import { parseBacklog } from "./backlog-parser";
import { evaluateRules } from "./rules/evaluate-rules";
import { scanPlanFiles } from "./plan-scanner";

const createStableHash = (value: string): string =>
  createHash("sha256").update(value).digest("hex");

const toRepoRelative = (repoRoot: string, absolutePath: string): string =>
  relative(repoRoot, absolutePath).replace(/\\/g, "/");

export interface DetectionRunResult {
  report: DetectionReport;
}

const resolveRepoRoot = (startDir: string): string => {
  let current = startDir;

  for (let depth = 0; depth < 10; depth += 1) {
    const backlogPath = join(current, ".hawp", "work", "BACKLOG.md");
    if (existsSync(backlogPath)) {
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

export const runDetection = (repoRoot: string): DetectionRunResult => {
  const resolvedRepoRoot = resolveRepoRoot(repoRoot);
  const workRoot = join(resolvedRepoRoot, ".hawp", "work");
  const backlogAbsolutePath = join(workRoot, "BACKLOG.md");
  const backlogRelativePath = toRepoRelative(
    resolvedRepoRoot,
    backlogAbsolutePath,
  );

  const backlog = parseBacklog(backlogAbsolutePath);
  const plans = scanPlanFiles(workRoot);
  const operations = evaluateRules(
    resolvedRepoRoot,
    workRoot,
    backlogRelativePath,
    backlog,
    plans,
  );

  const scannedAt = new Date().toISOString();
  const planSeed = JSON.stringify({ backlog: backlog.rows, operations });
  const planHash = createStableHash(planSeed);
  const planId = `PLAN-${planHash.slice(0, 8)}`;

  const plan = createBacklogFixPlan(
    planId,
    planHash,
    scannedAt,
    backlogRelativePath,
    plans.files.length,
    backlog.rows.length,
    operations,
    { version: "0.0.0" },
  );

  const reportHash = createStableHash(`${planHash}:${scannedAt}`);
  const report = createDetectionReport(
    `REPORT-${reportHash.slice(0, 8)}`,
    scannedAt,
    plan,
    Array.from(new Set(operations.map((operation) => operation.fileToModify))),
    Array.from(new Set(operations.map((operation) => operation.itemId))),
  );

  return { report };
};
