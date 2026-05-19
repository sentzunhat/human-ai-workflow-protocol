/**
 * Type-level tests for backlog upgrade data models
 * Verifies that types compile correctly and compose as expected
 *
 * This file tests that:
 * - Type definitions are syntactically correct
 * - Factory functions produce correct types
 * - Type guards work correctly
 * - All objects are JSON-serializable
 * - Rules and enums are properly defined
 */

import {
  Mode,
  OutputFormat,
  ExitCode,
  createBlockedItem,
  isBlockedItem,
  createBacklogFixOperation,
  createBacklogFixPlan,
  isBacklogFixPlan,
  createDetectionReport,
  isDetectionReport,
  createFileHashRecord,
  createValidatorStateSnapshot,
  createEvidenceReport,
  assessValidatorImprovement,
  isEvidenceReport,
} from "../index.js";

// Verify enums exist and are usable
void Mode.DryRun;
void OutputFormat.Json;
void ExitCode.Success;

// Test 1: BlockedItem creation and type guard
const blocked = createBlockedItem(
  "BLOCKED-001",
  "B1",
  "TASK-001",
  0.68,
  ["task", "bug", "improvement"],
  "Ambiguous type",
  { data: "test" },
  "Manual review",
);
const isBlockedCheck: boolean = isBlockedItem(blocked);

// Test 2: JSON serialization of BlockedItem
const blockedJson = JSON.stringify(blocked);
const blockedParsed = JSON.parse(blockedJson);
const isBlockedReparsed: boolean = isBlockedItem(blockedParsed);

// Test 3: BacklogFixOperation and BacklogFixPlan
const operation = createBacklogFixOperation(
  "OP-001",
  "add-field",
  "TASK-001",
  ".hawp/work/BACKLOG.md",
  [1, 5],
  "Add field",
  "safe",
  0.95,
);

const plan = createBacklogFixPlan(
  "PLAN-abc123",
  "sha256_hash",
  new Date().toISOString(),
  ".hawp/work/BACKLOG.md",
  5,
  10,
  [operation],
  { version: "1.0.0" },
);
const isPlanCheck: boolean = isBacklogFixPlan(plan);

// Test 4: DetectionReport
const report = createDetectionReport(
  "REPORT-001",
  new Date().toISOString(),
  plan,
  [".hawp/work/BACKLOG.md"],
  ["TASK-001"],
);
const isReportCheck: boolean = isDetectionReport(report);

// Test 5: EvidenceReport with validator state
const fileRecord = createFileHashRecord(
  ".hawp/work/BACKLOG.md",
  "hash_before",
  "hash_after",
  "fix",
  new Date().toISOString(),
);

const validatorBefore = createValidatorStateSnapshot(
  "hash_before",
  5,
  new Date().toISOString(),
);

const validatorAfter = createValidatorStateSnapshot(
  "hash_after",
  3,
  new Date().toISOString(),
);

const evidence = createEvidenceReport(
  "EVIDENCE-001",
  plan.planHash,
  new Date().toISOString(),
  validatorBefore,
  validatorAfter,
  [fileRecord],
  true,
);
const isEvidenceCheck: boolean = isEvidenceReport(evidence);

// Test 6: Validator improvement assessment
const improvement = assessValidatorImprovement(validatorBefore, validatorAfter);
const improvementImproved: boolean = improvement.improved;

// Export to verify all tests used
export const typeCheckResults = {
  isBlockedCheck,
  isBlockedReparsed,
  isPlanCheck,
  isReportCheck,
  isEvidenceCheck,
  improvementImproved,
};
