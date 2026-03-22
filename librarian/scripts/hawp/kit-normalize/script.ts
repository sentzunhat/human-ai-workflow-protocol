import { execSync } from "node:child_process";
import { existsSync, renameSync } from "node:fs";
import { resolve } from "node:path";

import { findBacklogRepoRoot, toRepoRelative } from "../../lib";
import { getHelpText, parseArgs } from "./cli";
import {
  applyLinkUpdates,
  planLinkUpdates,
} from "./mutations/internal-links";
import { planFileRenames } from "./mutations/file-naming";

const hasDirtyWorkingTree = (repoRoot: string): boolean => {
  try {
    const status = execSync("git status --short", {
      cwd: repoRoot,
      encoding: "utf-8",
      stdio: ["ignore", "pipe", "ignore"],
    });
    return status.trim().length > 0;
  } catch {
    return true;
  }
};

const buildRenameMap = (renames: Array<{ from: string; to: string }>): Map<string, string> => {
  const map = new Map<string, string>();
  for (const rename of renames) {
    map.set(rename.from, rename.to);
  }
  return map;
};

export const runKitNormalize = (argv: string[]): number => {
  try {
    const args = parseArgs(argv);
    if (args.help) {
      process.stdout.write(`${getHelpText()}\n`);
      return 0;
    }

    const repoRoot = findBacklogRepoRoot(resolve(process.cwd()));
    const kitPath = resolve(args.kitPath ?? resolve(repoRoot, ".hawp", "kit"));

    process.stdout.write("kit:normalize\n");
    process.stdout.write("=============\n");
    process.stdout.write(`kit: ${kitPath}\n`);
    process.stdout.write(`mode: ${args.apply ? "apply" : "dry-run"}\n\n`);

    const renames = planFileRenames(kitPath);
    const renameMap = buildRenameMap(renames);
    const linkUpdates = planLinkUpdates(kitPath, renameMap);

    if (!args.apply) {
      if (renames.length === 0 && linkUpdates.length === 0) {
        process.stdout.write("No kit normalization needed.\n");
        return 0;
      }

      if (renames.length > 0) {
        process.stdout.write("Planned file renames:\n");
        for (const rename of renames) {
          process.stdout.write(
            `- ${toRepoRelative(repoRoot, rename.from)} -> ${toRepoRelative(repoRoot, rename.to)}\n`,
          );
        }
      }

      if (linkUpdates.length > 0) {
        process.stdout.write("\nPlanned link updates:\n");
        for (const update of linkUpdates) {
          process.stdout.write(
            `- ${toRepoRelative(repoRoot, update.file)}: ${update.from} -> ${update.to}\n`,
          );
        }
      }

      return 0;
    }

    if (hasDirtyWorkingTree(repoRoot)) {
      process.stderr.write(
        "Error: apply mode requires a clean working tree. Re-run from a clean tree.\n",
      );
      return 1;
    }

    for (const rename of renames.sort((left, right) => right.from.length - left.from.length)) {
      if (existsSync(rename.to)) {
        process.stderr.write(
          `Error: cannot rename ${toRepoRelative(repoRoot, rename.from)} because ${toRepoRelative(repoRoot, rename.to)} already exists.\n`,
        );
        return 1;
      }
      renameSync(rename.from, rename.to);
    }

    const appliedLinkUpdates = applyLinkUpdates(linkUpdates);

    process.stdout.write(
      `Applied ${renames.length} rename(s) and ${linkUpdates.length} link update(s) across ${appliedLinkUpdates} file(s).\n`,
    );
    return 0;
  } catch (error) {
    process.stderr.write(`kit normalize error: ${String(error)}\n`);
    return 1;
  }
};
