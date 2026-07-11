import { resolve } from "node:path";

import { findBacklogRepoRoot } from "../../../lib";

export interface ResolveKitPathInput {
  cwd: string;
  kitPath?: string;
}

export const resolveKitPath = (input: ResolveKitPathInput): string => {
  if (input.kitPath !== undefined) {
    return resolve(input.cwd, input.kitPath);
  }

  const repoRoot = findBacklogRepoRoot(resolve(input.cwd));
  return resolve(repoRoot, ".hawp", "kit");
};
