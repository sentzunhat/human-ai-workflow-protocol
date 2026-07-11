import { readdirSync, statSync } from "node:fs";
import { join, relative } from "node:path";

import type { KitIssue } from "../script";

const ALLOWED_UPPERCASE = new Set(["README.md"]);

const isValidName = (name: string): boolean => {
  if (ALLOWED_UPPERCASE.has(name)) return true;
  // lowercase letters, digits, hyphens, dots only
  return /^[a-z0-9][a-z0-9.-]*$/.test(name);
};

const walk = (dir: string, kitRoot: string, issues: KitIssue[]): void => {
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    const rel = relative(kitRoot, full);
    if (!isValidName(entry)) {
      issues.push({ file: rel, message: `name should be lowercase-hyphen (got "${entry}")` });
    }
    if (statSync(full).isDirectory()) {
      walk(full, kitRoot, issues);
    }
  }
};

export const checkFileNaming = (kitPath: string): KitIssue[] => {
  const issues: KitIssue[] = [];
  walk(kitPath, kitPath, issues);
  return issues;
};
