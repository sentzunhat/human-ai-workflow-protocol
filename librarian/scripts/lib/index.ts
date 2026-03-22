/**
 * Shared utilities for librarian scripts.
 *
 * Keep this module dependency-free (node builtins only) so every script area
 * can use it without coupling to another area's models.
 */

import { existsSync, readdirSync } from "node:fs";
import { dirname, join, relative, resolve } from "node:path";

/** Closed files on or after this date require Outcome, Verification, and
 *  Close Checklist sections; earlier files are tolerated as legacy. */
export const LEGACY_CLOSED_CUTOFF = "2026-05-10";

/**
 * Walks upward from startDir until the predicate matches a directory.
 * Returns the matching directory or null when nothing matches.
 */
export const findUpward = (
  startDir: string,
  predicate: (dir: string) => boolean,
  maxDepth = 12,
): string | null => {
  let current = resolve(startDir);

  for (let depth = 0; depth < maxDepth; depth += 1) {
    if (predicate(current)) {
      return current;
    }

    const parent = dirname(current);
    if (parent === current) {
      break;
    }
    current = parent;
  }

  return null;
};

/**
 * Finds the repo root containing .hawp/work/BACKLOG.md. Throws when no
 * ancestor directory qualifies.
 */
export const findBacklogRepoRoot = (startDir: string): string => {
  const root = findUpward(startDir, (dir) =>
    existsSync(join(dir, ".hawp", "work", "BACKLOG.md")),
  );

  if (!root) {
    throw new Error(
      "Could not locate repo root containing .hawp/work/BACKLOG.md",
    );
  }

  return root;
};

/** Converts an absolute path to a POSIX-style path relative to repoRoot. */
export const toRepoRelative = (
  repoRoot: string,
  absolutePath: string,
): string => relative(repoRoot, absolutePath).replace(/\\/g, "/");

/**
 * Recursively collects markdown files under dirPath, skipping README.md
 * scaffolding files. Returns absolute paths.
 */
export const walkMarkdownFiles = (dirPath: string): string[] => {
  if (!existsSync(dirPath)) {
    return [];
  }

  const paths: string[] = [];
  for (const entry of readdirSync(dirPath, { withFileTypes: true })) {
    const absolutePath = join(dirPath, entry.name);
    if (entry.isDirectory()) {
      paths.push(...walkMarkdownFiles(absolutePath));
      continue;
    }

    if (
      !entry.isFile() ||
      !entry.name.endsWith(".md") ||
      entry.name === "README.md"
    ) {
      continue;
    }

    paths.push(absolutePath);
  }

  return paths;
};

/** Normalizes line endings so generated/expected content compares cleanly. */
export const normalizeForCompare = (content: string): string =>
  content.replace(/\r\n/g, "\n");
