/**
 * EvidenceReport — post-apply summary with immutable hashes
 * Used only after --apply mode completes
 * Enables auditing, reproducibility, rollback verification
 */

import type { HashAlgorithm } from "./types.js";

/**
 * Hash of a file before and after upgrade
 * Enables verification and rollback capability
 */
export interface FileHashRecord {
  /**
   * File path (relative to repo root)
   */
  filePath: string;

  /**
   * SHA256 hash before any modifications (empty string if file didn't exist)
   */
  hashBefore: string;

  /**
   * SHA256 hash after modifications (empty string if file was deleted)
   */
  hashAfter: string;

  /**
   * Description of what changed
   */
  operation: string;

  /**
   * When this operation was performed (ISO8601 format)
   */
  appliedAt: string;
}

/**
 * Validator state snapshot
 * Enables detection of whether upgrade improved backlog validation
 */
export interface ValidatorStateSnapshot {
  /**
   * SHA256 hash of validator state/findings
   * Deterministic: same backlog → same hash
   */
  hash: string;

  /**
   * Number of issues detected by validator
   */
  issuesFound: number;

  /**
   * When snapshot was captured (ISO8601 format)
   */
  capturedAt: string;

  /**
   * Validator version/build that created this snapshot
   */
  validatorVersion?: string | undefined;
}

/**
 * Complete evidence report after apply mode
 * Immutable artifact for audit trail and reproducibility
 */
export interface EvidenceReport {
  /**
   * Unique evidence report ID
   */
  reportId: string;

  /**
   * Hash of the plan that was executed
   * Links evidence to the exact plan that was applied
   */
  planHash: string;

  /**
   * When the plan was applied (ISO8601 format)
   */
  appliedAt: string;

  /**
   * Hash algorithm used (always "sha256")
   */
  hashAlgorithm: HashAlgorithm;

  /**
   * Validator state before apply
   * Baseline for measuring improvement
   */
  validatorStateBefore: ValidatorStateSnapshot;

  /**
   * Validator state after apply
   * Should show equal or fewer issues than before
   */
  validatorStateAfter: ValidatorStateSnapshot;

  /**
   * List of all file operations performed
   */
  fileOperations: FileHashRecord[];

  /**
   * Whether idempotency verification passed
   * Verified by running apply again and confirming no changes
   */
  idempotencyVerified: boolean;

  /**
   * Additional notes for operator/auditor
   */
  notes?: string | undefined;
}

/**
 * Validator improvement assessment
 * Quick check of whether validator findings improved
 */
export interface ValidatorImprovement {
  /**
   * Did validator issues decrease or stay the same?
   */
  improved: boolean;

  /**
   * Number of issues before
   */
  issuesBefore: number;

  /**
   * Number of issues after
   */
  issuesAfter: number;

  /**
   * Net change (negative = improvement)
   */
  netChange: number;

  /**
   * Human-readable assessment
   */
  assessment: string;
}

/**
 * Creates a FileHashRecord with all required fields
 */
export function createFileHashRecord(
  filePath: string,
  hashBefore: string,
  hashAfter: string,
  operation: string,
  appliedAt: string,
): FileHashRecord {
  return {
    filePath,
    hashBefore,
    hashAfter,
    operation,
    appliedAt,
  };
}

/**
 * Creates a ValidatorStateSnapshot with all required fields
 */
export function createValidatorStateSnapshot(
  hash: string,
  issuesFound: number,
  capturedAt: string,
  validatorVersion?: string | undefined,
): ValidatorStateSnapshot {
  return {
    hash,
    issuesFound,
    capturedAt,
    validatorVersion,
  };
}

/**
 * Creates an EvidenceReport with all required fields
 */
export function createEvidenceReport(
  reportId: string,
  planHash: string,
  appliedAt: string,
  validatorStateBefore: ValidatorStateSnapshot,
  validatorStateAfter: ValidatorStateSnapshot,
  fileOperations: FileHashRecord[],
  idempotencyVerified: boolean,
  notes?: string,
): EvidenceReport {
  return {
    reportId,
    planHash,
    appliedAt,
    hashAlgorithm: "sha256",
    validatorStateBefore,
    validatorStateAfter,
    fileOperations,
    idempotencyVerified,
    notes,
  };
}

/**
 * Assess validator improvement
 */
export function assessValidatorImprovement(
  before: ValidatorStateSnapshot,
  after: ValidatorStateSnapshot,
): ValidatorImprovement {
  const netChange = after.issuesFound - before.issuesFound;
  const improved = netChange <= 0;

  let assessment: string;
  if (netChange < 0) {
    assessment = `Validator issues reduced by ${Math.abs(netChange)}`;
  } else if (netChange === 0) {
    assessment = "Validator issues unchanged (no regression)";
  } else {
    assessment = `Validator issues increased by ${netChange} (unexpected)`;
  }

  return {
    improved,
    issuesBefore: before.issuesFound,
    issuesAfter: after.issuesFound,
    netChange,
    assessment,
  };
}

/**
 * Type guard to check if value is an EvidenceReport
 */
export function isEvidenceReport(value: unknown): value is EvidenceReport {
  if (typeof value !== "object" || value === null) {
    return false;
  }

  const report = value as Record<string, unknown>;

  return (
    typeof report["reportId"] === "string" &&
    typeof report["planHash"] === "string" &&
    typeof report["appliedAt"] === "string" &&
    (report["hashAlgorithm"] === "sha256" ||
      typeof report["hashAlgorithm"] === "string") &&
    typeof report["validatorStateBefore"] === "object" &&
    typeof report["validatorStateAfter"] === "object" &&
    Array.isArray(report["fileOperations"]) &&
    typeof report["idempotencyVerified"] === "boolean"
  );
}
