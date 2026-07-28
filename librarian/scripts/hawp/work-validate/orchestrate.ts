import { readFileSync } from "node:fs";
import { join } from "node:path";
import type {
  ValidationReport,
  BacklogRow,
  BacklogCheck,
  ClosedTaskCheck,
  DeadLinksCheck,
  EvidenceCheck,
  VerificationCheck,
} from "./types";
import { checkBacklogConsistency } from "./validations/backlog-consistency";
import { checkClosedTaskCompleteness } from "./validations/closed-task-completeness";
import { checkDeadLinks } from "./validations/dead-links";
import {
  checkEvidenceIntegrity,
  collectClosedPlanFiles,
} from "./validations/evidence-integrity";
import { checkVerificationClarity } from "./validations/verification-clarity";
import {
  extractIdFromFilename,
  extractShortUuid,
} from "./validations/id-parser";

/**
 * Options for validation orchestration
 */
export interface ValidationOptions {
  debugClosedTask?: boolean;
}

/**
 * Parse BACKLOG.md to extract rows by section
 */
export const parseBacklog = (
  backlogPath: string,
): {
  active: BacklogRow[];
  closed: BacklogRow[];
  parked: BacklogRow[];
} | null => {
  const parseTableCells = (line: string): string[] =>
    line
      .split("|")
      .slice(1, -1)
      .map((value) => value.trim());

  const normalizeHeader = (value: string): string => value.trim().toLowerCase();

  const stripCodeSpan = (value: string): string =>
    value.replace(/^`+|`+$/g, "");

  const getMappedCell = (
    cells: string[],
    headerMap: Map<string, number>,
    ...aliases: string[]
  ): string => {
    for (const alias of aliases) {
      const index = headerMap.get(alias);
      if (index !== undefined) {
        return cells[index] ?? "";
      }
    }

    return "";
  };

  try {
    const content = readFileSync(backlogPath, "utf-8");
    const lines = content.split("\n");

    const active: BacklogRow[] = [];
    const closed: BacklogRow[] = [];
    const parked: BacklogRow[] = [];

    let section: "active" | "closed" | "parked" | null = null;
    let headerMap: Map<string, number> | null = null;

    for (const line of lines) {
      if (line.includes("## Active Work")) {
        section = "active";
        headerMap = null;
        continue;
      }
      if (line.includes("## Recently Closed") || line.includes("## Done")) {
        section = "closed";
        headerMap = null;
        continue;
      }
      if (line.includes("## Blocked / Parked")) {
        section = "parked";
        headerMap = null;
        continue;
      }
      if (/^#{2,}\s/.test(line) && section !== null) {
        // Any second-level or deeper header (## or ###...) ends the current
        // section. This prevents nested `###` subsections from being parsed
        // as table rows (was a bug that caused `### Summaries` to be read as
        // active backlog rows).
        section = null;
        headerMap = null;
        continue;
      }

      if (!section || !line.startsWith("|")) {
        continue;
      }

      if (/^\|\s*-+/.test(line) || line.includes("| ---")) {
        continue;
      }

      const cells = parseTableCells(line);
      if (cells.length === 0) {
        continue;
      }

      if (headerMap === null) {
        headerMap = new Map(
          cells.map((value, index) => [normalizeHeader(value), index]),
        );
        continue;
      }

      // Rows may carry the ID in the Legacy ID cell (TASK-012), the UUID cell,
      // or both during the UUID migration. A cell holding only a placeholder
      // ("—") is skipped. Prefer the first cell whose value parses as a known
      // ID format; fall back to the first non-placeholder value.
      const currentHeaderMap = headerMap;
      const idCandidates = ["legacy id", "id", "uuid"]
        .map((alias) =>
          stripCodeSpan(getMappedCell(cells, currentHeaderMap, alias)),
        )
        .filter((value) => value !== "" && value !== "—" && value !== "-");
      const rawId =
        idCandidates.find((value) => extractIdFromFilename(value) !== null) ??
        idCandidates.find((value) => extractShortUuid(value) !== null) ??
        idCandidates[0];
      if (!rawId) {
        continue;
      }

      const normalizedId = normalizeHeader(rawId);
      if (
        normalizedId === "id" ||
        normalizedId === "legacy id" ||
        rawId.startsWith("**") ||
        /^[✓⏹️]+$/u.test(rawId)
      ) {
        continue;
      }

      const row: BacklogRow = {
        id: extractIdFromFilename(rawId) ?? extractShortUuid(rawId) ?? rawId,
        type: getMappedCell(cells, headerMap, "type"),
        title: getMappedCell(cells, headerMap, "title"),
        status: getMappedCell(cells, headerMap, "status", "reason", "closed"),
        detail: getMappedCell(cells, headerMap, "plan file", "detail"),
      };

      if (section === "active") {
        active.push(row);
      } else if (section === "closed") {
        closed.push(row);
      } else if (section === "parked") {
        parked.push(row);
      }
    }

    return { active, closed, parked };
  } catch (error) {
    console.error(
      `[validate] warning: failed to read backlog at ${backlogPath}: ${
        error instanceof Error ? error.message : String(error)
      }`,
    );
    return null;
  }
};

const countPassed = (
  backlog: BacklogCheck,
  completeness: ClosedTaskCheck,
  evidence: EvidenceCheck,
  clarity: VerificationCheck,
  deadLinks: DeadLinksCheck,
): number =>
  [backlog, completeness, evidence, clarity, deadLinks].filter(
    (check) => check.status === "PASS",
  ).length;

const countFailed = (
  backlog: BacklogCheck,
  completeness: ClosedTaskCheck,
  evidence: EvidenceCheck,
  clarity: VerificationCheck,
  deadLinks: DeadLinksCheck,
): number =>
  [backlog, completeness, evidence, clarity, deadLinks].filter(
    (check) => check.status === "FAIL",
  ).length;

const countWarnings = (
  completeness: ClosedTaskCheck,
  evidence: EvidenceCheck,
  clarity: VerificationCheck,
): number =>
  [completeness, evidence, clarity].filter((check) => check.status === "WARN")
    .length;

/**
 * Run all validation checks and build report
 */
export const orchestrateValidation = (
  workDir: string,
  backlogRows: {
    active: BacklogRow[];
    closed: BacklogRow[];
    parked: BacklogRow[];
  },
  options?: ValidationOptions,
): ValidationReport => {
  // Collect closed files
  const closedFiles = collectClosedPlanFiles(join(workDir, "closed"));

  // Run all checks
  const backlogCheck = checkBacklogConsistency(workDir, backlogRows);
  const completenessCheck = checkClosedTaskCompleteness(
    workDir,
    backlogRows.closed,
    options?.debugClosedTask ? { debug: true } : undefined,
  );
  const evidenceCheck = checkEvidenceIntegrity(workDir, closedFiles);
  const clarityCheck = checkVerificationClarity(closedFiles);
  const deadLinksCheck = checkDeadLinks(workDir);

  // Count results
  const passed = countPassed(
    backlogCheck,
    completenessCheck,
    evidenceCheck,
    clarityCheck,
    deadLinksCheck,
  );
  const failed = countFailed(
    backlogCheck,
    completenessCheck,
    evidenceCheck,
    clarityCheck,
    deadLinksCheck,
  );
  const warnings = countWarnings(
    completenessCheck,
    evidenceCheck,
    clarityCheck,
  );

  // Build and return report
  return {
    timestamp: new Date().toISOString(),
    checks: {
      backlogConsistency: backlogCheck,
      closedTaskCompleteness: completenessCheck,
      evidenceIntegrity: evidenceCheck,
      verificationClarity: clarityCheck,
      deadLinks: deadLinksCheck,
    },
    summary: {
      passed,
      failed,
      warnings,
      totalChecks: 5,
    },
    overallStatus: failed > 0 ? "FAIL" : "PASS",
  };
};
