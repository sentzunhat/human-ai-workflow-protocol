import { existsSync } from "node:fs";
import { join, relative } from "node:path";

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

const isCanonicalId = (value: string): boolean =>
  /^(TASK|BUG)-\d+$/.test(value);

const inferTypeFromId = (id: string): "task" | "bug" | undefined => {
  if (id.startsWith("TASK-")) {
    return "task";
  }
  if (id.startsWith("BUG-")) {
    return "bug";
  }
  return undefined;
};

const toRepoRelative = (repoRoot: string, absolutePath: string): string =>
  relative(repoRoot, absolutePath).replace(/\\/g, "/");

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

const hasHeading = (content: string, heading: string): boolean =>
  content.includes(`## ${heading}`);

const hasUnprovenChecklistMarker = (content: string): boolean =>
  /- \[ \].*NOT YET VERIFIED/i.test(content);

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
  const cutoff = new Date("2026-05-10");
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
    const planPath = join(context.workRoot, row.planPath);
    if (!existsSync(planPath)) {
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
    ? join(context.workRoot, row.planPath)
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
        "A4",
        "add-section-header",
        "BACKLOG-STRUCTURE",
        context.backlogPath,
        1,
        "A4: add missing required backlog section headers",
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
