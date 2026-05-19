/**
 * BacklogFixPlan and related types
 * Represents a complete upgrade plan with all proposed operations
 */

import type { BlockedItem } from "./blocked-item.js";
import type { OperationType, RuleId } from "./types.js";

/**
 * Single operation in a fix plan
 * Represents one structural fix to be applied
 */
export interface BacklogFixOperation {
  /**
   * Unique operation ID (e.g., "OP-001", "OP-042")
   * Used to track and reference individual operations
   */
  opId: string;

  /**
   * Type of operation (add-field, normalize-date, migrate-row, etc)
   */
  type: OperationType;

  /**
   * Backlog item ID affected by this operation (e.g., "TASK-001")
   */
  itemId: string;

  /**
   * File path to be modified
   * Always within .hawp/** boundary
   */
  fileToModify: string;

  /**
   * Line range affected [startLine, endLine] (1-indexed)
   * Helps user review changes
   */
  lineRange: [number, number];

  /**
   * Human-readable description of what this operation does
   */
  description: string;

  /**
   * If this operation was blocked, the block details
   * If undefined, operation can be auto-fixed
   */
  blocked?: BlockedItem | undefined;

  /**
   * For auto-fixable operations: the actual content change needed
   * Maps from old content to new content
   */
  contentChange?: { before: string; after: string } | undefined;

  /**
   * Safety level of this operation
   */
  safety: "safe" | "needs-review" | "blocked";

  /**
   * Confidence score (0.0-1.0) in this operation
   * Higher confidence = safer to apply automatically
   */
  confidence: number;

  /**
   * Detection rule ID responsible for this operation (A1-A7, B1-B6)
   */
  ruleId?: RuleId | undefined;
}

/**
 * Complete upgrade plan for a backlog
 * Contains all detected issues and proposed fixes
 */
export interface BacklogFixPlan {
  /**
   * Unique plan identifier (e.g., "PLAN-abc123")
   * Deterministic: same backlog → same plan ID
   */
  planId: string;

  /**
   * SHA256 hash of the complete plan
   * Immutable; enables reproducibility and audit trail
   * Same backlog state → same hash
   */
  planHash: string;

  /**
   * When the scan was performed (ISO8601 format)
   */
  scannedAt: string;

  /**
   * Path to the backlog file that was scanned
   * Always ".hawp/work/BACKLOG.md" in current design
   */
  backlogPath: string;

  /**
   * Total number of plan files scanned
   */
  filesScanned: number;

  /**
   * Total number of backlog items analyzed
   */
  itemsAnalyzed: number;

  /**
   * All proposed operations (auto-fixable and blocked)
   */
  operations: BacklogFixOperation[];

  /**
   * Count of auto-fixable operations (safety = "safe")
   */
  autoFixCount: number;

  /**
   * Count of blocked operations (safety = "blocked")
   */
  blockedCount: number;

  /**
   * Estimated number of file modifications (for --dry-run reporting)
   */
  estimatedChanges: number;

  /**
   * Metadata about the scan environment
   */
  metadata: {
    version: string; // e.g., "1.0.0"
    hostname?: string;
    userId?: string;
  };
}

/**
 * Creates a BacklogFixOperation with all required fields
 */
export function createBacklogFixOperation(
  opId: string,
  type: OperationType,
  itemId: string,
  fileToModify: string,
  lineRange: [number, number],
  description: string,
  safety: "safe" | "needs-review" | "blocked",
  confidence: number,
  blocked?: BlockedItem,
  contentChange?: { before: string; after: string },
  ruleId?: RuleId,
): BacklogFixOperation {
  return {
    opId,
    type,
    itemId,
    fileToModify,
    lineRange,
    description,
    safety,
    confidence,
    blocked,
    contentChange,
    ruleId,
  };
}

/**
 * Creates a BacklogFixPlan with all required fields
 */
export function createBacklogFixPlan(
  planId: string,
  planHash: string,
  scannedAt: string,
  backlogPath: string,
  filesScanned: number,
  itemsAnalyzed: number,
  operations: BacklogFixOperation[],
  metadata: { version: string; hostname?: string; userId?: string },
): BacklogFixPlan {
  const autoFixCount = operations.filter((op) => op.safety === "safe").length;
  const blockedCount = operations.filter(
    (op) => op.safety === "blocked",
  ).length;
  const estimatedChanges = operations.filter(
    (op) => op.safety === "safe",
  ).length;

  return {
    planId,
    planHash,
    scannedAt,
    backlogPath,
    filesScanned,
    itemsAnalyzed,
    operations,
    autoFixCount,
    blockedCount,
    estimatedChanges,
    metadata,
  };
}

/**
 * Type guard to check if value is a BacklogFixPlan
 */
export function isBacklogFixPlan(value: unknown): value is BacklogFixPlan {
  if (typeof value !== "object" || value === null) {
    return false;
  }

  const plan = value as Record<string, unknown>;

  return (
    typeof plan["planId"] === "string" &&
    typeof plan["planHash"] === "string" &&
    typeof plan["scannedAt"] === "string" &&
    typeof plan["backlogPath"] === "string" &&
    typeof plan["filesScanned"] === "number" &&
    typeof plan["itemsAnalyzed"] === "number" &&
    Array.isArray(plan["operations"]) &&
    typeof plan["autoFixCount"] === "number" &&
    typeof plan["blockedCount"] === "number" &&
    typeof plan["estimatedChanges"] === "number"
  );
}
