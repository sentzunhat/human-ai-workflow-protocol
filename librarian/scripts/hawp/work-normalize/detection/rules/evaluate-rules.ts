import { existsSync } from "node:fs";
import { resolve, sep } from "node:path";

import { LEGACY_CLOSED_CUTOFF, toRepoRelative } from "../../../../lib";
import {
  extractShortUuid,
} from "../../../work-validate/validations/id-parser";
import {
  createBacklogFixOperation,
  createBlockedItem,
  CONFIDENCE_LEVELS,
  type AutoFixRuleId,
  type BacklogFixOperation,
  type BlockedRuleId,
  type OperationType,
  type RuleId,
} from "../../models";
import type { BacklogParseResult, BacklogRow } from "../backlog-parser";
import type { PlanFileRecord, PlanScanResult } from "../plan-scanner";

interface RuleContext {
  repoRoot: string;
  workRoot: string;
  backlogPath: string;
  backlog: BacklogParseResult;
  plans: PlanScanResult;
  planByPath: Map<string, PlanFileRecord>;
}

const isIsoDate = (value: string | undefined): boolean => {
  if (!value) {
    return true;
  }
  if (value === "" || value.toLowerCase() === "n/a") {
    return true;
  }
  return /^\d{4}-\d{2}-\d{2}$/.test(value);
};

const FULL_UUID_PATTERN =
  /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

const isCanonicalId = (value: string): boolean =>
  /^(TASK|BUG)-\d+$/.test(value) ||
  FULL_UUID_PATTERN.test(value) ||
  extractShortUuid(value) !== null;

const inferTypeFromId = (id: string): "task" | "bug" | undefined => {
  if (id.startsWith("TASK-")) {
    return "task";
  }
  if (id.startsWith("BUG-")) {
    return "bug";
  }
  return undefined;
};

/**
 * Resolves a backlog-relative plan path while rejecting paths that escape
 * the work root (e.g. via "..").
 */
const resolveWithinWorkRoot = (
  workRoot: string,
  relativePath: string,
): string | undefined => {
  const root = resolve(workRoot);
  const resolved = resolve(root, relativePath);
  return resolved.startsWith(root + sep) ? resolved : undefined;
};

const createAutoFixOperation = (
  opId: string,
  ruleId: AutoFixRuleId,
  type: OperationType,
  itemId: string,
  fileToModify: string,
  lineNumber: number,
  description: string,
): BacklogFixOperation =>
  createBacklogFixOperation(
    opId,
    type,
    itemId,
    fileToModify,
    [lineNumber, lineNumber],
    description,
    "safe",
    CONFIDENCE_LEVELS.High,
    undefined,
    undefined,
    ruleId,
  );

const createBlockedOperation = (
  opId: string,
  itemId: string,
  fileToModify: string,
  lineNumber: number,
  ruleId: BlockedRuleId,
  reason: string,
  candidates: string[],
  recovery: string,
): BacklogFixOperation => {
  const blocked = createBlockedItem(
    `BLOCKED-${opId.replace("OP-", "")}`,
    ruleId,
    itemId,
    CONFIDENCE_LEVELS.None,
    candidates,
    reason,
    { fileToModify, lineNumber },
    recovery,
    fileToModify,
    lineNumber,
  );

  return createBacklogFixOperation(
    opId,
    "add-field",
    itemId,
    fileToModify,
    [lineNumber, lineNumber],
    reason,
    "blocked",
    CONFIDENCE_LEVELS.None,
    blocked,
    undefined,
    ruleId,
  );
};

const escapeRegExp = (value: string): string =>
  value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");

const hasHeading = (content: string, heading: string): boolean =>
  new RegExp(`^#{1,6}\\s+${escapeRegExp(heading)}\\b`, "im").test(content);

const isEvidenceFollowUpLine = (line: string): boolean => {
  const claim = line.replace(/^[\s-\[\]x ]+/, "").trim();
  return (
    claim.startsWith("Research evidence for:") ||
    claim.startsWith(
      "Update the original verification checklist line with Evidence:",
    )
  );
};

const hasUnprovenChecklistMarker = (content: string): boolean => {
  const sectionMatch = content.match(
    /^##\s+Verification\b([\s\S]*?)(?=^##\s+|\Z)/im,
  );
  if (!sectionMatch?.[1]) {
    return false;
  }

  return sectionMatch[1]
    .split("\n")
    .filter((line) => /- \[[x ]\]/.test(line))
    .filter((line) => !isEvidenceFollowUpLine(line))
    .some((line) => /(NOT YET VERIFIED|\bunproven\b)/i.test(line));
};

const extractAmbiguousVerificationClaims = (content: string): string[] => {
  const sectionMatch = content.match(/^##\s+Verification\b([\s\S]*?)(?=^##\s+|\Z)/im);
  if (!sectionMatch?.[1]) {
    return [];
  }

  return sectionMatch[1]
    .split("\n")
    .filter((line) => /- \[[x ]\]/.test(line))
    .filter((line) => !isEvidenceFollowUpLine(line))
    .filter(
      (line) =>
        !/Evidence:/i.test(line) &&
        !/NOT YET VERIFIED/i.test(line) &&
        !/\bunproven\b/i.test(line),
    )
    .map((line) => line.replace(/^[\s-\[\]x ]+/, "").trim())
    .filter((line) => line.length > 0);
};

const hasStaleTemplateReference = (content: string): boolean => {
  const stalePatterns = [
    /\]\(core\/distribution\/generated\//i,
    /\]\(core\/distribution\/sources\//i,
    /`core\/distribution\//i,
  ];
  return stalePatterns.some((pattern) => pattern.test(content));
};

const extractDateFromPath = (path: string): Date | undefined => {
  const dateMatch = path.match(/\/(\d{4})\/(\d{2})\/(\d{2})\//);
  if (dateMatch) {
    const [, year, month, day] = dateMatch;
    return new Date(`${year}-${month}-${day}`);
  }
  return undefined;
};

const isLegacyClosedFile = (path: string): boolean => {
  const cutoff = new Date(LEGACY_CLOSED_CUTOFF);
  const fileDate = extractDateFromPath(path);
  return fileDate ? fileDate < cutoff : false;
};

const evaluateRowRules = (
  context: RuleContext,
  row: BacklogRow,
  nextOpId: () => string,
): BacklogFixOperation[] => {
  const operations: BacklogFixOperation[] = [];
  const candidates = context.plans.byId.get(row.id) ?? [];

  if (!row.type) {
    const inferred = inferTypeFromId(row.id);
    if (inferred) {
      operations.push(
        createAutoFixOperation(
          nextOpId(),
          "A1",
          "add-field",
          row.id,
          context.backlogPath,
          row.lineNumber,
          `A1: add missing type field inferred from ID prefix (${inferred})`,
        ),
      );
    } else {
      operations.push(
        createBlockedOperation(
          nextOpId(),
          row.id,
          context.backlogPath,
          row.lineNumber,
          "B1",
          "Cannot infer missing type from ID",
          ["task", "bug", "improvement"],
          "Set an explicit type in BACKLOG.md",
        ),
      );
    }
  }

  if (!isIsoDate(row.updated)) {
    operations.push(
      createAutoFixOperation(
        nextOpId(),
        "A2",
        "normalize-date",
        row.id,
        context.backlogPath,
        row.lineNumber,
        `A2: normalize date '${row.updated ?? ""}' to YYYY-MM-DD`,
      ),
    );
  }

  if (!isCanonicalId(row.id)) {
    operations.push(
      createAutoFixOperation(
        nextOpId(),
        "A3",
        "fix-malformed-id",
        row.id,
        context.backlogPath,
        row.lineNumber,
        "A3: fix malformed backlog ID format",
      ),
    );
  }

  if (row.planPath) {
    const planPath = resolveWithinWorkRoot(context.workRoot, row.planPath);
    if (!planPath) {
      operations.push(
        createBlockedOperation(
          nextOpId(),
          row.id,
          context.backlogPath,
          row.lineNumber,
          "B2",
          `Referenced plan path escapes the work root: ${row.planPath}`,
          [],
          "Use a plan link relative to .hawp/work without '..' segments",
        ),
      );
    } else if (!existsSync(planPath)) {
      operations.push(
        createBlockedOperation(
          nextOpId(),
          row.id,
          context.backlogPath,
          row.lineNumber,
          "B2",
          `Referenced plan path is missing: ${row.planPath}`,
          [],
          "Update BACKLOG.md link or restore the plan file",
        ),
      );
    }
  } else if (candidates.length === 0) {
    operations.push(
      createBlockedOperation(
        nextOpId(),
        row.id,
        context.backlogPath,
        row.lineNumber,
        "B2",
        "No plan link and no matching plan file were found",
        [],
        "Create a plan file or add a valid plan link",
      ),
    );
  }

  if (candidates.length > 1) {
    operations.push(
      createBlockedOperation(
        nextOpId(),
        row.id,
        context.backlogPath,
        row.lineNumber,
        "B3",
        "Multiple plan files matched the same backlog ID",
        candidates.map((path) => toRepoRelative(context.repoRoot, path)),
        "Choose a canonical plan file and archive duplicate candidates",
      ),
    );
  }

  const linkedPath = row.planPath
    ? resolveWithinWorkRoot(context.workRoot, row.planPath)
    : undefined;
  const linkedPlan = linkedPath
    ? context.planByPath.get(linkedPath)
    : undefined;

  if (row.section === "recently-closed" && linkedPlan) {
    const content = linkedPlan.content;
    const isLegacy = isLegacyClosedFile(linkedPlan.path);
    const hasOutcome = hasHeading(content, "Outcome");
    const hasVerification = hasHeading(content, "Verification");
    const hasCloseChecklist = hasHeading(content, "Close Checklist");

    if (!isLegacy && !hasOutcome) {
      operations.push(
        createAutoFixOperation(
          nextOpId(),
          "A4",
          "add-section-header",
          row.id,
          toRepoRelative(context.repoRoot, linkedPlan.path),
          1,
          "A4: add missing '## Outcome' section heading",
        ),
      );
    }

    if (!isLegacy && (!hasVerification || !hasCloseChecklist)) {
      operations.push(
        createAutoFixOperation(
          nextOpId(),
          "A5",
          "add-scaffolding",
          row.id,
          toRepoRelative(context.repoRoot, linkedPlan.path),
          1,
          "A5: add missing verification/checklist scaffolding",
        ),
      );
    }

    if (hasUnprovenChecklistMarker(content)) {
      operations.push(
        createBlockedOperation(
          nextOpId(),
          row.id,
          toRepoRelative(context.repoRoot, linkedPlan.path),
          1,
          "B4",
          "Closed plan still contains unproven verification markers",
          [],
          "Complete verification evidence or mark unresolved risk explicitly before closing",
        ),
      );
    }

    const ambiguousClaims = extractAmbiguousVerificationClaims(content);
    if (ambiguousClaims.length > 0) {
      operations.push(
        createBlockedOperation(
          nextOpId(),
          row.id,
          toRepoRelative(context.repoRoot, linkedPlan.path),
          1,
          "B7",
          `Closed plan has ${ambiguousClaims.length} ambiguous verification checklist claim(s)`,
          ambiguousClaims.slice(0, 5),
          "Add `Evidence:` to each checklist claim or mark it explicitly unproven before relying on the record",
        ),
      );
    }

    if (!isLegacy && hasStaleTemplateReference(content)) {
      operations.push(
        createAutoFixOperation(
          nextOpId(),
          "A7",
          "update-template-reference",
          row.id,
          toRepoRelative(context.repoRoot, linkedPlan.path),
          1,
          "A7: update outdated template/path references",
        ),
      );
    }
  }

  if (
    row.section === "active" &&
    ["done", "wont-fix"].includes((row.status ?? "").toLowerCase())
  ) {
    operations.push(
      createAutoFixOperation(
        nextOpId(),
        "A6",
        "migrate-row",
        row.id,
        context.backlogPath,
        row.lineNumber,
        "A6: migrate completed active row to recently closed",
      ),
    );
  }

  return operations;
};

const evaluateStructuralRules = (
  context: RuleContext,
  nextOpId: () => string,
): BacklogFixOperation[] => {
  const operations: BacklogFixOperation[] = [];
  const sections = context.backlog.sectionPresence;

  if (!sections.active || !sections.blocked || !sections["recently-closed"]) {
    operations.push(
      createAutoFixOperation(
        nextOpId(),
        "A8",
        "add-section-header",
        "BACKLOG-STRUCTURE",
        context.backlogPath,
        1,
        "A8: add missing required backlog section headers",
      ),
    );
  }

  const missingDirs = Object.entries(context.plans.directoryPresence)
    .filter(([, present]) => !present)
    .map(([name]) => name);

  if (missingDirs.length > 0) {
    operations.push(
      createBlockedOperation(
        nextOpId(),
        "WORKSPACE",
        context.backlogPath,
        1,
        "B5",
        "Non-standard work folder structure detected",
        missingDirs,
        "Restore missing .hawp/work directories (active, parked, closed)",
      ),
    );
  }

  return operations;
};

export const evaluateRules = (
  repoRoot: string,
  workRoot: string,
  backlogPath: string,
  backlog: BacklogParseResult,
  plans: PlanScanResult,
): BacklogFixOperation[] => {
  const operations: BacklogFixOperation[] = [];
  let opCounter = 1;
  const nextOpId = (): string => `OP-${String(opCounter++).padStart(3, "0")}`;

  const planByPath = new Map<string, PlanFileRecord>(
    plans.files.map((file) => [file.path, file]),
  );

  const context: RuleContext = {
    repoRoot,
    workRoot,
    backlogPath,
    backlog,
    plans,
    planByPath,
  };

  operations.push(...evaluateStructuralRules(context, nextOpId));

  for (const row of backlog.rows) {
    operations.push(...evaluateRowRules(context, row, nextOpId));
  }

  return operations;
};

export type SupportedRuleId = RuleId;
