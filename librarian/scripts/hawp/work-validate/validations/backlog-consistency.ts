import { existsSync, readdirSync, readFileSync, statSync } from "node:fs";
import { join } from "node:path";
import type { BacklogCheck, BacklogRow } from "../types";
import { extractIdFromFilename, idsMatch } from "./id-parser";

/** True when any known row ID matches the given ID (exact or short-UUID prefix). */
const matchesAnyId = (ids: Set<string>, id: string): boolean => {
  if (ids.has(id)) return true;
  for (const known of ids) {
    if (idsMatch(known, id)) return true;
  }
  return false;
};

const warn = (context: string, error: unknown): void => {
  console.error(
    `[validate] warning: ${context}: ${
      error instanceof Error ? error.message : String(error)
    }`,
  );
};

/**
 * Checks backlog consistency: active/closed rows match actual files
 */
export function checkBacklogConsistency(
  workDir: string,
  backlogRows: {
    active: BacklogRow[];
    closed: BacklogRow[];
    parked: BacklogRow[];
  },
): BacklogCheck {
  const result: BacklogCheck = {
    activeWork: { total: 0, found: 0, missing: [] },
    recentlyClosed: { total: 0, found: 0, missing: [] },
    parkedWork: { total: 0, found: 0, missing: [] },
    orphanedFiles: [],
    orphanedParked: [],
    status: "PASS",
  };

  // Check active work
  result.activeWork.total = backlogRows.active.length;
  for (const row of backlogRows.active) {
    if (findActiveFile(workDir, row.id)) {
      result.activeWork.found++;
    } else {
      result.activeWork.missing.push(row.id);
    }
  }

  // Check recently closed
  result.recentlyClosed.total = backlogRows.closed.length;
  for (const row of backlogRows.closed) {
    // Closed files are in closed/YYYY/MM/DD/<ID>.md or closed/YYYY/MM/DD/<YYYY-MM-DD-*.md>
    const closedDir = join(workDir, "closed");
    const found = findClosedFile(closedDir, row.id);
    if (found) {
      result.recentlyClosed.found++;
    } else {
      result.recentlyClosed.missing.push(row.id);
    }
  }

  // Check parked work
  result.parkedWork.total = backlogRows.parked.length;
  for (const row of backlogRows.parked) {
    const filePath = extractLinkPath(row.detail);
    if (!filePath) {
      result.parkedWork.missing.push(row.id);
      continue;
    }
    if (existsSync(join(workDir, filePath))) {
      result.parkedWork.found++;
    } else {
      result.parkedWork.missing.push(row.id);
    }
  }

  // Check for orphaned files in active/
  const activeDir = join(workDir, "active");
  if (existsSync(activeDir)) {
    const activeIds = new Set(backlogRows.active.map((r) => r.id));
    collectOrphanedActive(activeDir, "active", activeIds, result.orphanedFiles);
  }

  // Check for orphaned files in parked/
  const parkedDir = join(workDir, "parked");
  if (existsSync(parkedDir)) {
    const files = readdirSync(parkedDir);
    const parkedIds = new Set(backlogRows.parked.map((r) => r.id));

    for (const file of files) {
      if (file.endsWith(".md")) {
        const id = extractIdFromFilename(file.replace(".md", ""));
        if (id && !matchesAnyId(parkedIds, id)) {
          result.orphanedParked.push(`parked/${file}`);
        }
      }
    }
  }

  // Set status
  if (
    result.activeWork.missing.length > 0 ||
    result.recentlyClosed.missing.length > 0 ||
    result.parkedWork.missing.length > 0 ||
    result.orphanedFiles.length > 0 ||
    result.orphanedParked.length > 0
  ) {
    result.status = "FAIL";
  }

  return result;
}

/**
 * True when a closed record's content names the given ID as a backtick'd
 * token (e.g. "`e2c4f9g5`"), the convention this repo uses in a shared
 * closed record's "Closes" list to name every backlog row it covers.
 */
function recordListsId(filePath: string, id: string): boolean {
  try {
    const content = readFileSync(filePath, "utf-8");
    const pattern = new RegExp(`\`${id.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")}\``, "i");
    return pattern.test(content);
  } catch (error) {
    warn(`error while reading closed record ${filePath}`, error);
    return false;
  }
}

/**
 * Recursively searches for a plan file in closed/ directory
 */
function findClosedFile(closedDir: string, id: string): boolean {
  if (!existsSync(closedDir)) return false;

  try {
    const years = readdirSync(closedDir);
    for (const year of years) {
      if (year === "README.md") continue;
      const yearPath = join(closedDir, year);
      const stat = statSync(yearPath);
      if (!stat.isDirectory()) continue;

      const months = readdirSync(yearPath);
      for (const month of months) {
        const monthPath = join(yearPath, month);
        const stat2 = statSync(monthPath);
        if (!stat2.isDirectory()) continue;

        const days = readdirSync(monthPath);
        for (const day of days) {
          const dayPath = join(monthPath, day);
          const stat3 = statSync(dayPath);
          if (!stat3.isDirectory()) continue;

          const files = readdirSync(dayPath);
          // Check for exact ID match (e.g., "TASK-010.md")
          if (files.includes(`${id}.md`)) {
            return true;
          }
          // Check for files containing the ID (e.g., "2026-05-01-bug-002-...")
          // or matching via short-UUID prefix (row `0e1c4afa` → `<full-uuid>.md`)
          for (const file of files) {
            if (!file.endsWith(".md")) continue;
            if (file.toLowerCase().includes(id.toLowerCase())) {
              return true;
            }
            const fileId = extractIdFromFilename(file.replace(/\.md$/i, ""));
            if (fileId && idsMatch(id, fileId)) {
              return true;
            }
            // A single closed record can document several backlog rows at
            // once (e.g. one audit session that fixed 7 findings) — a
            // filename can't literally contain every row's ID, so also
            // accept a record whose CONTENT names the ID (as a backtick'd
            // token, to avoid matching arbitrary substrings in prose).
            if (recordListsId(join(dayPath, file), id)) {
              return true;
            }
          }
        }
      }
    }
  } catch (error) {
    warn(`error while searching closed plans for ${id}`, error);
  }

  return false;
}

/**
 * Finds an active plan file supporting both flat and date-nested layouts:
 *   active/<ID>.md
 *   active/YYYY/MM/DD/<ID>.md
 */
function findActiveFile(workDir: string, id: string): boolean {
  // Flat layout (canonical)
  if (existsSync(join(workDir, "active", `${id}.md`))) return true;

  const activeDir = join(workDir, "active");
  if (!existsSync(activeDir)) return false;

  // Flat layout, short-UUID prefix match (row `0e1c4afa` → `<full-uuid>.md`)
  try {
    for (const entry of readdirSync(activeDir)) {
      if (!entry.endsWith(".md")) continue;
      const fileId = extractIdFromFilename(entry.replace(/\.md$/i, ""));
      if (fileId && idsMatch(id, fileId)) return true;
    }
  } catch (error) {
    warn(`error while scanning active plans for ${id}`, error);
  }

  // Date-nested layout: active/YYYY/MM/DD/<ID>.md

  try {
    for (const year of readdirSync(activeDir)) {
      if (!/^\d{4}$/.test(year)) continue;
      const yearPath = join(activeDir, year);
      if (!statSync(yearPath).isDirectory()) continue;

      for (const month of readdirSync(yearPath)) {
        const monthPath = join(yearPath, month);
        if (!statSync(monthPath).isDirectory()) continue;

        for (const day of readdirSync(monthPath)) {
          const dayPath = join(monthPath, day);
          if (!statSync(dayPath).isDirectory()) continue;

          const files = readdirSync(dayPath);
          if (files.includes(`${id}.md`)) return true;
          // Also match date-prefixed variants (e.g. 2026-05-06-TASK-011.md)
          if (
            files.some(
              (f) =>
                f.toLowerCase().includes(id.toLowerCase()) && f.endsWith(".md"),
            )
          )
            return true;
        }
      }
    }
  } catch (error) {
    warn(`error while searching active plans for ${id}`, error);
  }

  return false;
}

/**
 * Collects orphaned plan files from active/ covering both flat and date-nested layouts.
 * A file is orphaned if its extracted ID is not in the known active set.
 */
function collectOrphanedActive(
  dir: string,
  relPrefix: string,
  activeIds: Set<string>,
  out: string[],
): void {
  try {
    for (const entry of readdirSync(dir)) {
      const fullPath = join(dir, entry);
      if (statSync(fullPath).isDirectory()) {
        // Recurse into date-nested subdirectories only
        if (/^\d{2,4}$/.test(entry)) {
          collectOrphanedActive(
            fullPath,
            `${relPrefix}/${entry}`,
            activeIds,
            out,
          );
        }
        continue;
      }
      if (!entry.endsWith(".md")) continue;
      const id = extractIdFromFilename(entry.replace(/\.md$/i, ""));
      if (id && !matchesAnyId(activeIds, id)) {
        out.push(`${relPrefix}/${entry}`);
      }
    }
  } catch (error) {
    warn(`error while collecting orphaned plans in ${dir}`, error);
  }
}

/**
 * Extracts a file path from a Markdown link, e.g. "[plan](parked/TASK-013.md)" → "parked/TASK-013.md"
 */
function extractLinkPath(detail: string | undefined): string | null {
  if (!detail) return null;
  const match = detail.match(/\[.*?\]\((.*?)\)/);
  return match ? (match[1] ?? null) : null;
}
