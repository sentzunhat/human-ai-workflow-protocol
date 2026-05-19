import { existsSync, readdirSync, statSync } from "fs";
import { join } from "path";
import type { BacklogCheck, BacklogRow } from "../types";
import { extractIdFromFilename } from "./id-parser";

/**
 * Checks backlog consistency: active/closed rows match actual files
 */
export async function checkBacklogConsistency(
  workDir: string,
  backlogRows: {
    active: BacklogRow[];
    closed: BacklogRow[];
    parked: BacklogRow[];
  },
): Promise<BacklogCheck> {
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
        if (id && !parkedIds.has(id)) {
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
          for (const file of files) {
            if (
              file.toLowerCase().includes(id.toLowerCase()) &&
              file.endsWith(".md")
            ) {
              return true;
            }
          }
        }
      }
    }
  } catch {
    // Ignore errors during recursive search
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

  // Date-nested layout: active/YYYY/MM/DD/<ID>.md
  const activeDir = join(workDir, "active");
  if (!existsSync(activeDir)) return false;

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
  } catch {
    // Ignore
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
      if (id && !activeIds.has(id)) {
        out.push(`${relPrefix}/${entry}`);
      }
    }
  } catch {
    // Ignore
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
