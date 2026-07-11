import { existsSync } from "node:fs";
import { join } from "node:path";

import type { KitIssue } from "../script";

const REQUIRED = [
  "start-here.md",
  "usage/status-report.md",
  "usage/intake-workflow.md",
  "usage/init.md",
  "references/spec.md",
  "references/backlog-alignment.md",
];

export const checkRequiredFiles = (kitPath: string): KitIssue[] => {
  const issues: KitIssue[] = [];
  for (const rel of REQUIRED) {
    if (!existsSync(join(kitPath, rel))) {
      issues.push({ file: rel, message: "required kit file is missing" });
    }
  }
  return issues;
};
