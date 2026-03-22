/**
 * DetectionReport — scan results from backlog analysis
 */

import type { BacklogFixPlan } from "./backlog-fix-plan.js";

export interface DriftSyncPlanStep {
  itemId: string;
  canonicalPlan: string;
  duplicatePlans: string[];
  applySteps: string[];
}

export interface VerificationResearchItem {
  itemId: string;
  claim: string;
  filePath: string;
  lineNumber: number;
  recommendedAction: string;
}

/**
 * Result of scanning and analyzing a backlog
 * Produced by the detection engine in --dry-run mode
 */
export interface DetectionReport {
  /**
   * Unique report identifier
   */
  reportId: string;

  /**
   * When the report was generated (ISO8601 format)
   */
  generatedAt: string;

  /**
   * The complete upgrade plan derived from detection
   */
  plan: BacklogFixPlan;

  /**
   * Overall assessment: is the backlog clean or does it need fixes?
   */
  assessment: "clean" | "drift-detected";

  /**
   * Summary statistics
   */
  summary: {
    totalIssues: number; // autoFixCount + blockedCount
    autoFixableCount: number;
    blockedCount: number;
    filesAffected: string[]; // unique files that would be modified
    itemsAffected: string[]; // unique item IDs (TASK-001, BUG-042, etc)
  };

  /**
   * Recommended next action for user
   */
  recommendation: {
    action: "no-action" | "apply-fixes" | "manual-review" | "export-plan";
    reason: string;
    details?: string[] | undefined;
  };

  /**
   * Concrete sync/apply plan for working-file drift (typically B3 duplicates)
   */
  syncPlan: DriftSyncPlanStep[];

  /**
   * Verification claims that need evidence research or an explicit unproven label
   */
  researchQueue: VerificationResearchItem[];
}

const buildSyncPlan = (plan: BacklogFixPlan): DriftSyncPlanStep[] => {
  const steps: DriftSyncPlanStep[] = [];

  for (const operation of plan.operations) {
    if (operation.ruleId !== "B3" || !operation.blocked) {
      continue;
    }

    const uniqueCandidates = Array.from(new Set(operation.blocked.candidates));
    if (uniqueCandidates.length < 2) {
      continue;
    }

    const closedCandidate = uniqueCandidates.find((candidate) =>
      candidate.includes("/closed/"),
    );

    if (!closedCandidate) {
      steps.push({
        itemId: operation.itemId,
        canonicalPlan: "manual-selection-required",
        duplicatePlans: uniqueCandidates,
        applySteps: [
          `Review candidates for ${operation.itemId} and choose one canonical plan file.`,
          "Move remaining duplicates to the closed archive or remove stale copies.",
          "Re-run: ./.hawp/bin/hawp backlog upgrade --dry-run --validate",
        ],
      });
      continue;
    }

    const duplicatePlans = uniqueCandidates.filter(
      (candidate) => candidate !== closedCandidate,
    );

    const applySteps = [
      `Keep canonical closed plan: ${closedCandidate}`,
      ...duplicatePlans.map((candidate) => {
        if (candidate.includes("/active/")) {
          return `Remove stale active copy: ${candidate}`;
        }
        if (candidate.includes("/parked/")) {
          return `Remove stale parked copy: ${candidate}`;
        }
        return `Archive or remove duplicate copy: ${candidate}`;
      }),
      "Re-run: ./.hawp/bin/hawp backlog upgrade --dry-run --validate",
    ];

    steps.push({
      itemId: operation.itemId,
      canonicalPlan: closedCandidate,
      duplicatePlans,
      applySteps,
    });
  }

  return steps;
};

/**
 * Creates a DetectionReport with all required fields
 */
export function createDetectionReport(
  reportId: string,
  generatedAt: string,
  plan: BacklogFixPlan,
  filesAffected: string[],
  itemsAffected: string[],
  researchQueue: VerificationResearchItem[] = [],
): DetectionReport {
  const totalIssues = plan.autoFixCount + plan.blockedCount;
  const assessment = totalIssues === 0 ? "clean" : "drift-detected";
  const syncPlan = buildSyncPlan(plan);

  // Derive recommendation from assessment
  let action: "no-action" | "apply-fixes" | "manual-review" | "export-plan";
  let reason: string;
  const details: string[] = [];

  if (assessment === "clean") {
    action = "no-action";
    reason = "Backlog is consistent with current templates and standards";
  } else if (plan.blockedCount > 0) {
    action = "manual-review";
    reason = `${plan.blockedCount} issue(s) require manual review before fixes can be applied`;
    details.push(
      `Run: ./.hawp/bin/hawp backlog upgrade --dry-run --format json to see blocked items`,
    );
    if (syncPlan.length > 0) {
      details.push(
        `Concrete sync/apply steps generated for ${syncPlan.length} duplicate-item drift case(s).`,
      );
    }
  } else {
    action = "apply-fixes";
    reason = `${plan.autoFixCount} mechanical fix(es) ready to apply`;
    details.push(
      `Review changes: ./.hawp/bin/hawp backlog upgrade --dry-run --format text`,
    );
    details.push(`Apply fixes: ./.hawp/bin/hawp backlog upgrade --apply`);
  }

  return {
    reportId,
    generatedAt,
    plan,
    assessment,
    summary: {
      totalIssues,
      autoFixableCount: plan.autoFixCount,
      blockedCount: plan.blockedCount,
      filesAffected,
      itemsAffected,
    },
    recommendation: {
      action,
      reason,
      details: details.length > 0 ? details : undefined,
    },
    syncPlan,
    researchQueue,
  };
}

/**
 * Type guard to check if value is a DetectionReport
 */
export function isDetectionReport(value: unknown): value is DetectionReport {
  if (typeof value !== "object" || value === null) {
    return false;
  }

  const report = value as Record<string, unknown>;

  return (
    typeof report["reportId"] === "string" &&
    typeof report["generatedAt"] === "string" &&
    typeof report["plan"] === "object" &&
    (report["assessment"] === "clean" ||
      report["assessment"] === "drift-detected") &&
    typeof report["summary"] === "object" &&
    typeof report["recommendation"] === "object"
  );
}
