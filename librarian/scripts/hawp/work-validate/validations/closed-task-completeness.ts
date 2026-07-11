import { readFileSync, statSync, readdirSync, existsSync } from "node:fs";
import { join } from "node:path";
import { LEGACY_CLOSED_CUTOFF } from "../../../lib";
import type { ClosedTaskCheck, BacklogRow } from "../types";
import { extractIdFromFilename } from "./id-parser";

/** Closed files on or after this date require Outcome, Verification, Close Checklist (FAIL if missing).
 *  Files before this date produce WARN only. */
const LEGACY_CUTOFF = LEGACY_CLOSED_CUTOFF;

const SUPPORTING_SUFFIXES = [
  "-summary",
  "-status",
  "-status-report",
  "-checkpoint",
  "-evidence",
];

type ClosedFileClassification =
  | { kind: "plan"; id: string }
  | { kind: "supporting"; id: string; reason: string }
  | { kind: "legacy-untyped"; id: string; reason: string }
  | { kind: "current-untyped"; id: string; reason: string };

interface ClosedFileEntry {
  filePath: string;
  date: string; // "YYYY-MM-DD" extracted from path, or "" if unknown
}

/**
 * Checks that closed plan files have required sections.
 * Supporting files (summaries, status reports, archives) are skipped.
 * Files before LEGACY_CUTOFF produce WARN; on/after produce FAIL.
 */
export function checkClosedTaskCompleteness(
  workDir: string,
  _closedRows: BacklogRow[],
  options?: { debug?: boolean },
): ClosedTaskCheck {
  const debug = options?.debug ?? false;
  const result: ClosedTaskCheck = {
    total: 0,
    skipped: 0,
    withOutcome: 0,
    withVerification: 0,
    withCloseChecklist: 0,
    failing: [],
    warnings: [],
    supportingSkipped: [],
    untypedLegacy: [],
    untypedCurrent: [],
    debugFindings: [],
    status: "PASS",
  };

  const entries = collectClosedFiles(join(workDir, "closed"));

  for (const { filePath, date } of entries) {
    const filename = filePath.split("/").pop() ?? "";

    const classification = classifyClosedFile(filename, date);

    if (classification.kind === "supporting") {
      result.skipped++;
      result.supportingSkipped.push({
        id: classification.id,
        date: date || "unknown",
        reason: classification.reason,
        filePath,
      });
      continue;
    }

    if (classification.kind === "legacy-untyped") {
      result.untypedLegacy.push({
        id: classification.id,
        date: date || "unknown",
        reason: classification.reason,
        filePath,
      });
      continue;
    }

    if (classification.kind === "current-untyped") {
      result.untypedCurrent.push({
        id: classification.id,
        date: date || "unknown",
        reason: classification.reason,
        filePath,
      });
      continue;
    }

    result.total++;
    const content = readFileSync(filePath, "utf-8");
    const headingMatches = findRequiredHeadingMatches(content);
    const missingSections: string[] = [];

    if (!headingMatches.Outcome) {
      missingSections.push("Outcome");
    } else {
      result.withOutcome++;
    }

    if (!headingMatches.Verification) {
      missingSections.push("Verification");
    } else {
      result.withVerification++;
    }

    if (!headingMatches.CloseChecklist) {
      missingSections.push("Close Checklist");
    } else {
      result.withCloseChecklist++;
    }

    if (missingSections.length > 0) {
      const id = classification.id;
      const isLegacy = !date || date < LEGACY_CUTOFF;
      if (isLegacy) {
        result.warnings.push({
          id,
          sections: missingSections,
          date: date || "unknown",
          filePath,
        });
      } else {
        result.failing.push({ id, sections: missingSections, date, filePath });
      }

      if (debug) {
        result.debugFindings.push({
          id,
          date: date || "unknown",
          filePath,
          exists: existsSync(filePath),
          headingMatches: [
            headingMatches.Outcome,
            headingMatches.Verification,
            headingMatches.CloseChecklist,
          ].filter((line): line is string => Boolean(line)),
          missingSections,
        });
      }
    }
  }

  if (result.failing.length > 0 || result.untypedCurrent.length > 0) {
    result.status = "FAIL";
  } else if (result.warnings.length > 0 || result.untypedLegacy.length > 0) {
    result.status = "WARN";
  }

  return result;
}

/**
 * Classifies a closed file into plan/supporting/untyped buckets.
 * Supporting skips are explicit-pattern only, so legacy untyped files stay visible.
 */
function classifyClosedFile(
  filename: string,
  date: string,
): ClosedFileClassification {
  const nameWithoutExt = filename.replace(/\.md$/i, "");
  const nameLower = nameWithoutExt.toLowerCase();
  const id = extractIdFromFilename(nameWithoutExt);
  const isLegacy = !date || date < LEGACY_CUTOFF;

  if (nameLower.startsWith("backlog")) {
    return {
      kind: "supporting",
      id: nameWithoutExt,
      reason: "matches BACKLOG supporting-file pattern",
    };
  }

  if (nameLower.includes("archive")) {
    return {
      kind: "supporting",
      id: nameWithoutExt,
      reason: "matches archive supporting-file pattern",
    };
  }

  if (id) {
    const idPos = nameLower.indexOf(id.toLowerCase());
    const suffix = idPos >= 0 ? nameLower.slice(idPos + id.length) : "";

    if (
      SUPPORTING_SUFFIXES.some((kw) => suffix === kw || suffix.endsWith(kw))
    ) {
      return {
        kind: "supporting",
        id: nameWithoutExt,
        reason: `matches supporting suffix pattern (${suffix || "none"})`,
      };
    }

    if (suffix.includes("-archive")) {
      return {
        kind: "supporting",
        id: nameWithoutExt,
        reason: "matches archive supporting-file pattern",
      };
    }

    return { kind: "plan", id };
  }

  if (isLegacy) {
    return {
      kind: "legacy-untyped",
      id: nameWithoutExt,
      reason: "legacy file without TASK-/BUG-style ID",
    };
  }

  return {
    kind: "current-untyped",
    id: nameWithoutExt,
    reason: "current file without TASK-/BUG-style ID",
  };
}

interface RequiredHeadingMatches {
  Outcome?: string;
  Verification?: string;
  CloseChecklist?: string;
}

function findRequiredHeadingMatches(content: string): RequiredHeadingMatches {
  const matches: RequiredHeadingMatches = {};
  const lines = content.split("\n");

  for (const line of lines) {
    const trimmed = line.trim();

    if (!matches.Outcome && /^#{1,6}\s*Outcome\b/i.test(trimmed)) {
      matches.Outcome = trimmed;
      continue;
    }

    if (!matches.Verification && /^#{1,6}\s*Verification\b/i.test(trimmed)) {
      matches.Verification = trimmed;
      continue;
    }

    if (
      !matches.CloseChecklist &&
      /^#{1,6}\s*Close Checklist\b/i.test(trimmed)
    ) {
      matches.CloseChecklist = trimmed;
    }
  }

  return matches;
}

/**
 * Recursively collect all .md files under closed/ with generic folder support.
 * Date is extracted from any /YYYY/MM/DD/ segment in path; unknown if none found.
 */
function collectClosedFiles(closedDir: string): ClosedFileEntry[] {
  const files: ClosedFileEntry[] = [];

  try {
    const stat = statSync(closedDir);
    if (!stat.isDirectory()) return files;

    collectClosedFilesRecursive(closedDir, files);
  } catch {
    // Ignore missing or unreadable directory
  }

  return files;
}

function collectClosedFilesRecursive(
  dir: string,
  out: ClosedFileEntry[],
): void {
  for (const entry of readdirSync(dir)) {
    const fullPath = join(dir, entry);
    const stat = statSync(fullPath);

    if (stat.isDirectory()) {
      collectClosedFilesRecursive(fullPath, out);
      continue;
    }

    if (!entry.endsWith(".md") || entry === "README.md") continue;

    out.push({
      filePath: fullPath,
      date: extractDateFromPath(fullPath),
    });
  }
}

function extractDateFromPath(filePath: string): string {
  const normalized = filePath.replace(/\\/g, "/");
  const match = normalized.match(/\/(\d{4})\/(\d{2})\/(\d{2})\//);
  return match ? `${match[1]}-${match[2]}-${match[3]}` : "";
}
