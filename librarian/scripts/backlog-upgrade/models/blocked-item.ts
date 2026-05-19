/**
 * BlockedItem — represents an item that cannot be automatically fixed
 * Always includes rule, confidence, candidates, and reason for blocking
 */

import type { BlockedRuleId } from "./types.js";

/**
 * Structured explanation for why an item is blocked
 * Provides clear, actionable information for human resolution
 */
export interface BlockedItem {
  /**
   * Unique identifier for this blocked item (e.g., "BLOCKED-001", "BLOCKED-042")
   */
  blockId: string;

  /**
   * The blocking rule that was triggered (B1-B6)
   * Indicates what type of issue prevents automation
   */
  rule: BlockedRuleId;

  /**
   * Confidence score (0.0-1.0) for the best candidate or fix
   * 0.0 = multiple equally likely candidates (complete ambiguity)
   * 1.0 = single clear candidate (no ambiguity)
   *
   * When confidence < CONFIDENCE_THRESHOLD_FOR_AUTOFIX, item is blocked
   */
  confidence: number;

  /**
   * List of candidate values/interpretations/options
   * For B1 (ambiguous type): ["task", "bug", "improvement"]
   * For B3 (multiple files): [path1, path2, path3]
   * Empty array means no valid candidates found
   */
  candidates: string[];

  /**
   * Human-readable explanation of why automation stopped
   * Should be concise and actionable
   *
   * Examples:
   * - "title too generic; inference confidence 68% < 90% threshold"
   * - "BUG-003.md found in 2 different date folders; cannot determine canonical"
   * - "no plan file found; backlog row references missing artifact"
   */
  reason: string;

  /**
   * Structured evidence supporting the block decision
   * Contains data used to derive the blocking decision
   * Enables debugging, future tooling, and agent understanding
   */
  evidence: Record<string, unknown>;

  /**
   * Action required to resolve this block
   * Should explain the manual step(s) needed
   *
   * Examples:
   * - "User must assign explicit type in BACKLOG.md"
   * - "User inspects evidence folder, determines canonical version, consolidates files"
   * - "User decides: delete row, create plan, or archive to decisions/legacy/"
   */
  recovery: string;

  /**
   * Backlog item ID that triggered this block (e.g., "TASK-001", "BUG-042")
   */
  itemId: string;

  /**
   * File path where the blocking issue was detected
   * Helps user locate the problem
   */
  filePath?: string | undefined;

  /**
   * Line number where issue detected (if applicable)
   */
  lineNumber?: number | undefined;
}

/**
 * Creates a BlockedItem with all required fields
 * Ensures consistency and type safety
 */
export function createBlockedItem(
  blockId: string,
  rule: BlockedRuleId,
  itemId: string,
  confidence: number,
  candidates: string[],
  reason: string,
  evidence: Record<string, unknown>,
  recovery: string,
  filePath?: string,
  lineNumber?: number,
): BlockedItem {
  return {
    blockId,
    rule,
    itemId,
    confidence,
    candidates,
    reason,
    evidence,
    recovery,
    filePath,
    lineNumber,
  };
}

/**
 * Type guard to check if value is a BlockedItem
 */
export function isBlockedItem(value: unknown): value is BlockedItem {
  if (typeof value !== "object" || value === null) {
    return false;
  }

  const item = value as Record<string, unknown>;

  return (
    typeof item["blockId"] === "string" &&
    typeof item["rule"] === "string" &&
    typeof item["itemId"] === "string" &&
    typeof item["confidence"] === "number" &&
    Array.isArray(item["candidates"]) &&
    typeof item["reason"] === "string" &&
    typeof item["recovery"] === "string" &&
    typeof item["evidence"] === "object"
  );
}
