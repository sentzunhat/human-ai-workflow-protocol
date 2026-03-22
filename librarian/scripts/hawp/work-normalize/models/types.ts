/**
 * Shared enums and constants for backlog upgrade tool
 * Defines mode, exit codes, rule classifications, and operation types
 */

/**
 * Execution modes for the upgrade command
 */
export enum Mode {
  DryRun = "dry-run",
  Apply = "apply",
}

/**
 * Output format options
 */
export enum OutputFormat {
  Text = "text",
  Json = "json",
}

/**
 * Exit codes for CLI operations
 */
export enum ExitCode {
  Success = 0,
  Error = 1,
  UsageError = 2,
}

/**
 * Auto-fixable rule IDs (deterministic, mechanical fixes)
 * All A-series rules are safe to apply automatically with high confidence
 */
export type AutoFixRuleId =
  | "A1" // Missing type field
  | "A2" // Normalize date format
  | "A3" // Fix malformed ID
  | "A4" // Add missing section header
  | "A5" // Add scaffolding for empty evidence sections
  | "A6" // Migrate closed work row
  | "A7" // Update outdated template references
  | "A8"; // Add missing required backlog section headers

/**
 * Blocked rule IDs (require manual resolution)
 * All B-series rules indicate ambiguity or missing information that automation cannot resolve
 */
export type BlockedRuleId =
  | "B1" // Ambiguous type inference (confidence below threshold)
  | "B2" // Orphaned records (no plan file, no evidence)
  | "B3" // Multiple plan file candidates (ambiguous consolidation)
  | "B4" // Evidence synthesis needed (missing verification content)
  | "B5" // Non-standard folder structure
  | "B6" // Evidence integrity issues (hash mismatch)
  | "B7"; // Ambiguous verification checklist wording

/**
 * All rule IDs combined
 */
export type RuleId = AutoFixRuleId | BlockedRuleId;

/**
 * Operation types for auto-fixes
 */
export type OperationType =
  | "add-field"
  | "normalize-date"
  | "fix-malformed-id"
  | "add-section-header"
  | "add-scaffolding"
  | "migrate-row"
  | "update-template-reference";

/**
 * Confidence threshold for auto-fix eligibility (0.0 - 1.0)
 * Rules with confidence >= this threshold can be auto-fixed
 * Rules below this threshold are blocked (B1 - Ambiguous type inference)
 */
export const CONFIDENCE_THRESHOLD_FOR_AUTOFIX = 0.9;

/**
 * Confidence score interpretation guidelines
 */
export const CONFIDENCE_LEVELS = {
  Certain: 1.0, // No ambiguity
  VeryHigh: 0.95, // Extremely confident
  High: 0.9, // High confidence, safe threshold
  Medium: 0.7, // Moderate confidence
  Low: 0.5, // Low confidence, risky
  VeryLow: 0.25, // Very uncertain
  None: 0.0, // Multiple equally likely candidates
} as const;
