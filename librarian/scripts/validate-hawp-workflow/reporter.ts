import type { ValidationReport } from "./types";

/**
 * Format validation report to human-readable string
 */
export const formatReport = (report: ValidationReport): string => {
  const lines: string[] = [];

  lines.push(
    "======================================================================",
  );
  lines.push("HAWP Workflow Validation Report");
  lines.push(
    "======================================================================",
  );
  lines.push("");

  // Backlog Consistency
  lines.push("1. BACKLOG CONSISTENCY");
  lines.push(
    "----------------------------------------------------------------------",
  );
  const bc = report.checks.backlogConsistency;
  lines.push(`\nActive Work (${bc.activeWork.total} items):`);
  lines.push(`  Found: ${bc.activeWork.found}/${bc.activeWork.total}`);
  if (bc.activeWork.missing.length > 0) {
    lines.push(`  ✗ Missing plan files:`);
    bc.activeWork.missing.forEach((id) => {
      lines.push(`    - ${id}`);
    });
  }

  lines.push(`\nRecently Closed (${bc.recentlyClosed.total} items):`);
  lines.push(`  Found: ${bc.recentlyClosed.found}/${bc.recentlyClosed.total}`);
  if (bc.recentlyClosed.missing.length > 0) {
    lines.push(`  ✗ Missing plan files:`);
    bc.recentlyClosed.missing.forEach((id) => {
      lines.push(`    - ${id}`);
    });
  }

  if (bc.orphanedFiles.length > 0) {
    lines.push(`\nOrphaned Files (in active/ without backlog row):`);
    bc.orphanedFiles.forEach((f) => {
      lines.push(`  ✗ ${f}`);
    });
  } else {
    lines.push(`\nOrphaned Files (in active/ without backlog row):`);
    lines.push(`  (none)`);
  }

  lines.push(`\nBlocked / Parked (${bc.parkedWork.total} items):`);
  lines.push(`  Found: ${bc.parkedWork.found}/${bc.parkedWork.total}`);
  if (bc.parkedWork.missing.length > 0) {
    lines.push(`  ✗ Missing plan files:`);
    bc.parkedWork.missing.forEach((id) => {
      lines.push(`    - ${id}`);
    });
  }
  if (bc.orphanedParked.length > 0) {
    lines.push(`\nOrphaned Files (in parked/ without backlog row):`);
    bc.orphanedParked.forEach((f) => {
      lines.push(`  ✗ ${f}`);
    });
  } else {
    lines.push(`\nOrphaned Files (in parked/ without backlog row):`);
    lines.push(`  (none)`);
  }

  lines.push("");

  // Closed Task Completeness
  lines.push("2. CLOSED TASK COMPLETENESS");
  lines.push(
    "----------------------------------------------------------------------",
  );
  const ctc = report.checks.closedTaskCompleteness;
  const skippedNote =
    ctc.skipped > 0 ? `  (${ctc.skipped} supporting file(s) skipped)` : "";
  lines.push(`\nChecking ${ctc.total} plan file(s)${skippedNote}:`);
  lines.push(`  Outcome: ${ctc.withOutcome}/${ctc.total}`);
  lines.push(`  Verification: ${ctc.withVerification}/${ctc.total}`);
  lines.push(`  Close Checklist: ${ctc.withCloseChecklist}/${ctc.total}`);

  if (ctc.untypedCurrent.length > 0) {
    lines.push(
      `\n  [FAIL] Untyped closed files (2026-05-10 or later — must be tracked):`,
    );
    ctc.untypedCurrent.forEach((item) => {
      lines.push(
        `    ${item.id}: ${item.reason}  (${item.date}) [source: ${item.filePath}]`,
      );
    });
  }

  if (ctc.failing.length > 0) {
    lines.push(`\n  [FAIL] Missing sections (2026-05-10 or later — must fix):`);
    ctc.failing.forEach((item) => {
      lines.push(
        `    ${item.id}: ✗ missing ${item.sections.join(", ")}  (${item.date}) [source: ${item.filePath}]`,
      );
    });
  }

  if (ctc.untypedLegacy.length > 0) {
    lines.push(
      `\n  [WARN] Legacy untyped closed files (before 2026-05-10 — tolerated, visible):`,
    );
    ctc.untypedLegacy.forEach((item) => {
      lines.push(
        `    ${item.id}: ${item.reason}  (${item.date}) [source: ${item.filePath}]`,
      );
    });
  }

  if (ctc.warnings.length > 0) {
    lines.push(
      `\n  [WARN] Legacy files missing sections (before 2026-05-10 — tolerated):`,
    );
    ctc.warnings.forEach((item) => {
      lines.push(
        `    ${item.id}: missing ${item.sections.join(", ")}  (${item.date}) [source: ${item.filePath}]`,
      );
    });
  }

  if (ctc.debugFindings.length > 0) {
    lines.push(`\n  [DEBUG] Closed-task diagnostics for flagged files:`);
    ctc.debugFindings.forEach((item) => {
      lines.push(`    ${item.id} (${item.date})`);
      lines.push(`      - resolved path: ${item.filePath}`);
      lines.push(`      - exists: ${item.exists ? "yes" : "no"}`);
      lines.push(
        `      - first matching headings: ${item.headingMatches.length > 0 ? item.headingMatches.join(" | ") : "(none)"}`,
      );
      lines.push(
        `      - missing sections: ${item.missingSections.join(", ") || "(none)"}`,
      );
    });
  }

  if (ctc.supportingSkipped.length > 0) {
    lines.push(`\n  [INFO] Supporting files skipped by pattern:`);
    ctc.supportingSkipped.forEach((item) => {
      lines.push(`    ${item.id}: ${item.reason}  (${item.date})`);
    });
  }

  lines.push("");

  // Evidence Integrity
  lines.push("3. EVIDENCE INTEGRITY");
  lines.push(
    "----------------------------------------------------------------------",
  );
  const ev = report.checks.evidenceIntegrity;
  lines.push(`\n  Found ${ev.total} evidence links`);
  lines.push(`  ✓ ${ev.valid} valid links`);
  if (ev.broken.length > 0) {
    lines.push(`  ✗ ${ev.broken.length} broken links`);
    ev.broken.forEach((item) => {
      lines.push(`    ${item.id}: ${item.link}`);
    });
  }

  lines.push("");

  // Verification Clarity
  lines.push("4. VERIFICATION CLARITY");
  lines.push(
    "----------------------------------------------------------------------",
  );
  const vc = report.checks.verificationClarity;
  lines.push(`\n  Proven: ${vc.proven}/${vc.total}`);
  if (vc.unproven.length > 0) {
    lines.push(`  Unproven: ${vc.unproven.length}`);
    vc.unproven.forEach((item) => {
      lines.push(`    ${item.id}: ${item.claim}`);
    });
  }

  lines.push("");

  // Summary
  lines.push(
    "======================================================================",
  );
  lines.push("SUMMARY");
  lines.push(
    "======================================================================",
  );
  lines.push("");
  lines.push(`✓ Checks passed:     ${report.summary.passed}`);
  lines.push(`✗ Issues found:      ${report.summary.failed}`);
  lines.push(`! Warnings:          ${report.summary.warnings}`);
  lines.push("");
  lines.push(`Result: VALIDATION ${report.overallStatus}`);
  lines.push("");
  lines.push(
    "======================================================================",
  );

  return lines.join("\n");
};
