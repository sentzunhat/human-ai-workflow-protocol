/**
 * Backlog Upgrade Data Models — JSON-first architecture
 *
 * Exports TypeScript types and factory functions for the backlog upgrade tool.
 * Internal objects are source-of-truth; text/JSON output is rendered from them.
 */

// Shared types and enums
export {
  Mode,
  OutputFormat,
  ExitCode,
  type AutoFixRuleId,
  type BlockedRuleId,
  type RuleId,
  type OperationType,
  CONFIDENCE_THRESHOLD_FOR_AUTOFIX,
  CONFIDENCE_LEVELS,
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
  type VerificationResearchItem,
  createDetectionReport,
  isDetectionReport,
} from "./detection-report.js";
