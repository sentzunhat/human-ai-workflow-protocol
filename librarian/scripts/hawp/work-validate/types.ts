// Shared type definitions for HAWP workflow validation

/**
 * Status type for validation check results
 * - PASS: all requirements met, no issues
 * - FAIL: critical issue that fails validation (exits with code 1)
 * - WARN: tolerated issue that's visible but doesn't fail validation
 */
export type CheckStatus = "PASS" | "FAIL" | "WARN";

/**
 * Overall validation report status (only PASS or FAIL, no WARN)
 */
export type ReportStatus = "PASS" | "FAIL";

export interface CheckResult {
  passed: boolean;
  message: string;
  details?: string[];
  errors?: string[];
}

export interface ValidationReport {
  timestamp: string;
  checks: {
    backlogConsistency: BacklogCheck;
    closedTaskCompleteness: ClosedTaskCheck;
    evidenceIntegrity: EvidenceCheck;
    verificationClarity: VerificationCheck;
    deadLinks: DeadLinksCheck;
  };
  summary: {
    passed: number;
    failed: number;
    warnings: number;
    totalChecks: number;
  };
  overallStatus: ReportStatus;
}

export interface BacklogCheck {
  activeWork: {
    total: number;
    found: number;
    missing: string[];
  };
  recentlyClosed: {
    total: number;
    found: number;
    missing: string[];
  };
  parkedWork: {
    total: number;
    found: number;
    missing: string[];
  };
  orphanedFiles: string[];
  orphanedParked: string[];
  status: CheckStatus;
}

export interface ClosedTaskCheck {
  total: number; // plan files checked (supporting files excluded)
  skipped: number; // supporting files skipped
  withOutcome: number;
  withVerification: number;
  withCloseChecklist: number;
  failing: Array<{
    id: string;
    sections: string[];
    date: string;
    filePath: string;
  }>; // on/after cutoff → FAIL
  warnings: Array<{
    id: string;
    sections: string[];
    date: string;
    filePath: string;
  }>; // before cutoff → WARN
  supportingSkipped: Array<{
    id: string;
    date: string;
    reason: string;
    filePath: string;
  }>;
  untypedLegacy: Array<{
    id: string;
    date: string;
    reason: string;
    filePath: string;
  }>;
  untypedCurrent: Array<{
    id: string;
    date: string;
    reason: string;
    filePath: string;
  }>;
  debugFindings: Array<{
    id: string;
    date: string;
    filePath: string;
    exists: boolean;
    headingMatches: string[];
    missingSections: string[];
  }>;
  status: CheckStatus;
}

export interface EvidenceCheck {
  total: number;
  valid: number;
  broken: Array<{ id: string; link: string }>;
  status: CheckStatus;
}

export interface VerificationCheck {
  total: number;
  proven: number;
  unproven: Array<{
    id: string;
    claim: string;
    filePath: string;
    lineNumber: number;
  }>;
  ambiguous: Array<{
    id: string;
    claim: string;
    filePath: string;
    lineNumber: number;
  }>;
  status: CheckStatus;
}

export interface DeadLinksCheck {
  scanned: number;
  broken: Array<{ file: string; link: string }>;
  status: CheckStatus;
}

export interface BacklogRow {
  id: string;
  type: string;
  title: string;
  status: string;
  detail: string | undefined;
}
