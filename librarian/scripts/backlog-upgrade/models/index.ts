/**
 * Backlog Upgrade Data Models — JSON-first architecture
 *
 * This module exports all TypeScript types and factory functions for the backlog upgrade tool.
 * Follows JSON-first design: internal objects are source-of-truth, text/JSON are renderers.
 *
 * All types are serializable to JSON for:
 * - CLI output (rendered to text/JSON)
 * - Storage (evidence reports)
 * - Future APIs/UIs/agents (consume same objects)
 * - Immutable audit trail (with SHA256 hashes)
 */

// Shared types and enums
export {
  Mode,
  OutputFormat,
  ExitCode,
  DetectionType,
  SafetyLevel,
  type AutoFixRuleId,
  type BlockedRuleId,
  type RuleId,
  type OperationType,
  type HashAlgorithm,
  type TimestampFormat,
  CONFIDENCE_THRESHOLD_FOR_AUTOFIX,
  CONFIDENCE_LEVELS,
  MAX_CANDIDATES_TO_SHOW,
  ALLOWED_WRITE_ROOT,
  DEFAULT_HASH_ALGORITHM,
  DEFAULT_TIMESTAMP_FORMAT,
} from "./types.js";

// Blocked item types and utilities
export {
  type BlockedItem,
  createBlockedItem,
  isBlockedItem,
} from "./blocked-item.js";

// Backlog fix plan types and utilities
export {
  type BacklogFixOperation,
  type BacklogFixPlan,
  createBacklogFixOperation,
  createBacklogFixPlan,
  isBacklogFixPlan,
} from "./backlog-fix-plan.js";

// Detection report types and utilities
export {
  type DetectionReport,
  type DriftSyncPlanStep,
  createDetectionReport,
  isDetectionReport,
} from "./detection-report.js";

// Evidence report types and utilities
export {
  type FileHashRecord,
  type ValidatorStateSnapshot,
  type EvidenceReport,
  type ValidatorImprovement,
  createFileHashRecord,
  createValidatorStateSnapshot,
  createEvidenceReport,
  assessValidatorImprovement,
  isEvidenceReport,
} from "./evidence-report.js";
