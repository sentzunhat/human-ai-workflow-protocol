import type { DetectionReport } from "../models";

export const renderTextReport = (report: DetectionReport): string => {
  const lines: string[] = [];

  lines.push("HAWP Backlog Upgrade Dry-Run Report");
  lines.push("=================================");
  lines.push(`Report ID: ${report.reportId}`);
  lines.push(`Generated: ${report.generatedAt}`);
  lines.push(`Assessment: ${report.assessment}`);
  lines.push("");

  lines.push("Summary");
  lines.push("-------");
  lines.push(`Total issues: ${report.summary.totalIssues}`);
  lines.push(`Auto-fixable: ${report.summary.autoFixableCount}`);
  lines.push(`Blocked: ${report.summary.blockedCount}`);
  lines.push(
    `Files affected: ${report.summary.filesAffected.join(", ") || "(none)"}`,
  );
  lines.push(
    `Items affected: ${report.summary.itemsAffected.join(", ") || "(none)"}`,
  );
  lines.push("");

  lines.push("Recommendation");
  lines.push("--------------");
  lines.push(
    `${report.recommendation.action}: ${report.recommendation.reason}`,
  );
  for (const detail of report.recommendation.details ?? []) {
    lines.push(`- ${detail}`);
  }
  lines.push("");

  lines.push("Drift Sync/Apply Plan");
  lines.push("---------------------");
  if (report.syncPlan.length === 0) {
    lines.push("No duplicate working-file drift actions required.");
  } else {
    for (const syncItem of report.syncPlan) {
      lines.push(`${syncItem.itemId}: canonical=${syncItem.canonicalPlan}`);
      if (syncItem.duplicatePlans.length > 0) {
        lines.push(`  duplicates: ${syncItem.duplicatePlans.join(", ")}`);
      }
      for (const step of syncItem.applySteps) {
        lines.push(`  - ${step}`);
      }
    }
  }
  lines.push("");

  lines.push("Verification Research Queue");
  lines.push("---------------------------");
  if (report.researchQueue.length === 0) {
    lines.push("No verification evidence follow-up items detected.");
  } else {
    for (const item of report.researchQueue) {
      lines.push(
        `${item.itemId}:${item.lineNumber} ${item.claim} [source: ${item.filePath}]`,
      );
      lines.push(`  - ${item.recommendedAction}`);
    }
  }
  lines.push("");

  lines.push("Operations");
  lines.push("----------");

  if (report.plan.operations.length === 0) {
    lines.push("No modifications needed.");
    return lines.join("\n");
  }

  for (const operation of report.plan.operations) {
    const ruleLabel = operation.ruleId ? `[${operation.ruleId}] ` : "";
    lines.push(
      `${operation.opId} ${ruleLabel}[${operation.safety}] ${operation.itemId} ${operation.fileToModify}:${operation.lineRange[0]} - ${operation.description}`,
    );

    if (operation.blocked) {
      lines.push(
        `  blocked (${operation.blocked.rule}, confidence=${operation.blocked.confidence}): ${operation.blocked.reason}`,
      );
      if (operation.blocked.candidates.length > 0) {
        lines.push(`  candidates: ${operation.blocked.candidates.join(", ")}`);
      }
    }
  }

  return lines.join("\n");
};

export const renderJsonReport = (report: DetectionReport): string =>
  `${JSON.stringify(report, null, 2)}\n`;
