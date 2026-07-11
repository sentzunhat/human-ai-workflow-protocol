import { resolve } from "node:path";

import { findBacklogRepoRoot } from "../../lib";
import { parseArgs } from "./cli";
import { checkFileNaming } from "./validations/file-naming";
import { checkRequiredFiles } from "./validations/required-files";
import { checkInternalLinks } from "./validations/internal-links";

export interface KitIssue {
  file: string;
  message: string;
}

export interface KitValidateResult {
  kitPath: string;
  issues: KitIssue[];
  checks: number;
}

export const runKitValidateScript = (kitPath: string): KitValidateResult => {
  const namingIssues = checkFileNaming(kitPath);
  const requiredIssues = checkRequiredFiles(kitPath);
  const linkIssues = checkInternalLinks(kitPath);

  return {
    kitPath,
    issues: [...namingIssues, ...requiredIssues, ...linkIssues],
    checks: 3,
  };
};

export const runKitValidate = (argv: string[]): number => {
  const args = parseArgs(argv);
  const repoRoot = findBacklogRepoRoot(resolve(process.cwd()));
  const kitPath = args.kitPath ?? resolve(repoRoot, ".hawp", "kit");

  process.stdout.write("kit:validate\n");
  process.stdout.write("============\n");
  process.stdout.write(`kit: ${kitPath}\n\n`);

  const result = runKitValidateScript(kitPath);

  if (result.issues.length === 0) {
    process.stdout.write(`✓ ${result.checks} checks passed, 0 issues\n`);
    return 0;
  }

  for (const issue of result.issues) {
    process.stderr.write(`✗ ${issue.file}: ${issue.message}\n`);
  }
  process.stderr.write(`\n${result.issues.length} issue(s) found\n`);
  return 1;
};
