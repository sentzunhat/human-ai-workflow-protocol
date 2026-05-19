import { readFileSync } from "fs";
import { join } from "path";
import type {
  ValidationReport,
  BacklogRow,
  BacklogCheck,
  ClosedTaskCheck,
  EvidenceCheck,
  VerificationCheck,
} from "./types";
import { checkBacklogConsistency } from "./validations/backlog-consistency";
import { checkClosedTaskCompleteness } from "./validations/closed-task-completeness";
import {
  checkEvidenceIntegrity,
  collectClosedPlanFiles,
} from "./validations/evidence-integrity";
import { checkVerificationClarity } from "./validations/verification-clarity";
import { extractIdFromFilename } from "./validations/id-parser";

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
      if (line.startsWith("## ") && section !== null) {
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

      const rawId = stripCodeSpan(
        getMappedCell(cells, headerMap, "legacy id", "id", "uuid"),
      );
      if (!rawId) {
        continue;
      }

      const normalizedId = normalizeHeader(rawId);
      if (normalizedId === "id" || normalizedId === "legacy id") {
        continue;
      }

      const row: BacklogRow = {
        id: extractIdFromFilename(rawId) ?? rawId,
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
  } catch {
    return null;
  }
};

/**
 * Count checks with PASS status
 */
const countPassed = (
  backlog: BacklogCheck,
  completeness: ClosedTaskCheck,
  evidence: EvidenceCheck,
  clarity: VerificationCheck,
): number =>
  [backlog, completeness, evidence, clarity].filter(
    (check) => check.status === "PASS",
  ).length;

/**
 * Count checks with FAIL status
 */
const countFailed = (
  backlog: BacklogCheck,
  completeness: ClosedTaskCheck,
  evidence: EvidenceCheck,
  clarity: VerificationCheck,
): number =>
  [backlog, completeness, evidence, clarity].filter(
    (check) => check.status === "FAIL",
  ).length;

/**
 * Count checks with WARN status
 * WARN = tolerated issues (e.g. legacy files missing sections, unverified links)
 */
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
export const orchestrateValidation = async (
  workDir: string,
  backlogRows: {
    active: BacklogRow[];
    closed: BacklogRow[];
    parked: BacklogRow[];
  },
  options?: ValidationOptions,
): Promise<ValidationReport> => {
  // Collect closed files
  const closedFiles = collectClosedPlanFiles(join(workDir, "closed"));

  // Run all checks
  const backlogCheck = await checkBacklogConsistency(workDir, backlogRows);
  const completenessCheck = await checkClosedTaskCompleteness(
    workDir,
    backlogRows.closed,
    options?.debugClosedTask ? { debug: true } : undefined,
  );
  const evidenceCheck = await checkEvidenceIntegrity(workDir, closedFiles);
  const clarityCheck = await checkVerificationClarity(closedFiles);

  // Count results
  const passed = countPassed(
    backlogCheck,
    completenessCheck,
    evidenceCheck,
    clarityCheck,
  );
  const failed = countFailed(
    backlogCheck,
    completenessCheck,
    evidenceCheck,
    clarityCheck,
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
    },
    summary: {
      passed,
      failed,
      warnings,
      totalChecks: 4,
    },
    overallStatus: failed > 0 ? "FAIL" : "PASS",
  };
};
