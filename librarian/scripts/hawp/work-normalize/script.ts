import { execSync } from "node:child_process";
import {
  existsSync,
  mkdirSync,
  renameSync,
  readFileSync,
  writeFileSync,
} from "node:fs";
import { basename, dirname, join, resolve } from "node:path";

import { findBacklogRepoRoot, walkMarkdownFiles } from "../../lib";
import { runDetection } from "./detection";
import { ExitCode, Mode, OutputFormat } from "./models/index.js";
import { renderJsonReport, renderTextReport } from "./output/formatters";
import {
  orchestrateValidation,
  parseBacklog as parseValidationBacklog,
} from "../work-validate/orchestrate";
import { checkVerificationClarity } from "../work-validate/validations/verification-clarity";
import type { VerificationResearchItem } from "./models";

export interface BacklogUpgradeScriptOptions {
  mode: Mode;
  validate: boolean;
  format: OutputFormat;
  output?: string;
  exportPlan?: string;
  exportResearchQueue?: string;
  forceDirty: boolean;
  verbose: boolean;
  repoRoot?: string;
}

export interface BacklogUpgradeScriptResult {
  exitCode: ExitCode;
  stdoutText: string;
  stderrLines: string[];
  notices: string[];
}

const withTrailingNewline = (value: string): string =>
  value.endsWith("\n") ? value : `${value}\n`;

const escapeRegExp = (value: string): string =>
  value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");

const hasHeading = (content: string, heading: string): boolean =>
  new RegExp(`^#{1,6}\\s+${escapeRegExp(heading)}\\b`, "im").test(content);

const inferBacklogIdFromPath = (filePath: string): string | undefined => {
  const match = filePath.match(/(TASK|BUG)-\d+/i);
  return match?.[0]?.toUpperCase();
};

const inferClosedDateFromPath = (filePath: string): string | undefined => {
  const match = filePath.match(/\/closed\/(\d{4})\/(\d{2})\/(\d{2})\//);
  if (!match) {
    return undefined;
  }

  const [, year, month, day] = match;
  return `${year}-${month}-${day}`;
};

const inferFileDatePrefix = (filePath: string): string | undefined => {
  const fileName = basename(filePath);
  const match = fileName.match(/^(\d{4})-(\d{2})-(\d{2})-/);
  if (!match) {
    return undefined;
  }

  const [, year, month, day] = match;
  return `${year}-${month}-${day}`;
};

const ensureBlankLine = (content: string): string =>
  content.endsWith("\n\n")
    ? content
    : content.endsWith("\n")
      ? `${content}\n`
      : `${content}\n\n`;

const appendSection = (
  content: string,
  heading: string,
  body: string,
): string => {
  if (hasHeading(content, heading)) {
    return content;
  }

  return `${ensureBlankLine(content)}## ${heading}\n\n${body}\n`;
};

const ensureVerificationResearchFollowUp = (
  content: string,
  filePath: string,
): { content: string; addedClaims: string[] } => {
  const clarity = checkVerificationClarity([filePath]);
  if (clarity.ambiguous.length === 0) {
    return { content, addedClaims: [] };
  }

  const verificationMatch =
    /(^##\s+Verification\b[^\n]*\n)([\s\S]*?)(?=^##\s+|$)/m.exec(content);
  if (!verificationMatch || verificationMatch.index === undefined) {
    return { content, addedClaims: [] };
  }

  const heading = verificationMatch[1] ?? "";
  const sectionBody = verificationMatch[2] ?? "";
  const fullMatch = verificationMatch[0];
  const subsectionHeading = "### Evidence Follow-Up";

  const existingClaims = new Set<string>();
  for (const match of sectionBody.matchAll(
    /^\s*-\s+\[ \]\s+Research evidence for: (.+)$/gm,
  )) {
    if (match[1]) {
      existingClaims.add(match[1].trim());
    }
  }

  const pendingClaims = clarity.ambiguous
    .map((item) => item.claim.trim())
    .filter((claim) => claim.length > 0 && !existingClaims.has(claim));

  if (pendingClaims.length === 0) {
    return { content, addedClaims: [] };
  }

  const newEntries = pendingClaims
    .map(
      (claim) =>
        `- [ ] Research evidence for: ${claim}\n- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.`,
    )
    .join("\n");

  const nextSectionBody = sectionBody.includes(subsectionHeading)
    ? `${sectionBody.replace(/\s*$/u, "")}\n${newEntries}\n`
    : `${sectionBody.replace(/\s*$/u, "")}\n\n${subsectionHeading}\n\n${newEntries}\n`;

  const replacement = `${heading}${nextSectionBody}`;
  const nextContent = `${content.slice(0, verificationMatch.index)}${replacement}${content.slice(verificationMatch.index + fullMatch.length)}`;
  return { content: nextContent, addedClaims: pendingClaims };
};

const insertBacklogId = (content: string, backlogId: string): string => {
  if (/\*\*Backlog ID:\*\*/i.test(content)) {
    return content;
  }

  const lines = content.split("\n");
  const headingIndex = lines.findIndex((line) => /^#\s+/.test(line));

  if (headingIndex >= 0) {
    lines.splice(headingIndex + 1, 0, "", `**Backlog ID:** ${backlogId}`, "");
    return lines.join("\n").replace(/\n{3,}/g, "\n\n");
  }

  return `**Backlog ID:** ${backlogId}\n\n${content}`;
};

const normalizeClosedRecord = (
  content: string,
  filePath: string,
): { content: string; addedClaims: string[] } => {
  let updated = content;
  const inferredBacklogId = inferBacklogIdFromPath(filePath);

  if (inferredBacklogId) {
    updated = insertBacklogId(updated, inferredBacklogId);
  }

  updated = appendSection(
    updated,
    "Outcome",
    "_Legacy normalization scaffold added._",
  );
  updated = appendSection(
    updated,
    "Verification",
    "_Legacy normalization scaffold added._",
  );
  updated = appendSection(
    updated,
    "Close Checklist",
    "- [ ] Legacy normalization scaffold added.",
  );

  return ensureVerificationResearchFollowUp(updated, filePath);
};

const reconcileClosedRecordPath = (
  root: string,
  absolutePath: string,
): { path: string; moved: boolean } => {
  const fileDate = inferFileDatePrefix(absolutePath);
  if (!fileDate) {
    return { path: absolutePath, moved: false };
  }

  const currentDate = inferClosedDateFromPath(absolutePath);
  if (currentDate === fileDate) {
    return { path: absolutePath, moved: false };
  }

  const fileName = basename(absolutePath);
  const [year, month, day] = fileDate.split("-") as [string, string, string];
  const targetPath = join(
    root,
    ".hawp",
    "work",
    "closed",
    year,
    month,
    day,
    fileName,
  );

  if (targetPath === absolutePath) {
    return { path: absolutePath, moved: false };
  }

  mkdirSync(dirname(targetPath), { recursive: true });
  if (!existsSync(targetPath)) {
    renameSync(absolutePath, targetPath);
    return { path: targetPath, moved: true };
  }

  return { path: absolutePath, moved: false };
};

const hasDirtyWorkingTree = (repoRoot: string): boolean => {
  try {
    const status = execSync("git status --short", {
      cwd: repoRoot,
      encoding: "utf-8",
      stdio: ["ignore", "pipe", "ignore"],
    });
    return status.trim().length > 0;
  } catch {
    // Fail closed: without git we cannot prove the tree is clean, so block
    // apply mode rather than mutate files with no safety net.
    return true;
  }
};

const applyClosedRecordNormalization = (
  root: string,
): {
  changedFiles: string[];
  skippedFiles: string[];
  researchQueue: VerificationResearchItem[];
} => {
  const workRoot = join(root, ".hawp", "work");
  const closedRoot = join(workRoot, "closed");
  const skippedFiles: string[] = [];
  const touchedFiles = new Set<string>();
  const researchQueue: VerificationResearchItem[] = [];

  for (const absolutePath of walkMarkdownFiles(closedRoot)) {
    const reconciled = reconcileClosedRecordPath(root, absolutePath);
    const currentPath = reconciled.path;

    if (reconciled.moved) {
      touchedFiles.add(currentPath);
    }

    const current = readFileSync(currentPath, "utf-8");
    const normalized = normalizeClosedRecord(current, currentPath);
    const next = normalized.content;

    if (next === current) {
      if (
        !/\*\*Backlog ID:\*\*/i.test(current) &&
        !inferBacklogIdFromPath(absolutePath)
      ) {
        skippedFiles.push(absolutePath);
      }
      continue;
    }

    writeFileSync(currentPath, next, "utf-8");
    touchedFiles.add(currentPath);

    for (const claim of normalized.addedClaims) {
      researchQueue.push({
        itemId:
          inferBacklogIdFromPath(currentPath) ??
          basename(currentPath, ".md").toUpperCase(),
        claim,
        filePath: currentPath,
        lineNumber: 0,
        recommendedAction:
          "Gather supporting proof for this verification claim, then replace the original checklist entry with an Evidence: citation or mark it explicitly unproven.",
      });
    }
  }

  return {
    changedFiles: Array.from(touchedFiles),
    skippedFiles,
    researchQueue,
  };
};

const writeOptionalFile = (outputPath: string, content: string): void => {
  const absolutePath = resolve(outputPath);
  mkdirSync(dirname(absolutePath), { recursive: true });
  writeFileSync(absolutePath, content, "utf-8");
};

const hasMixedIdColumns = (repoRoot: string): boolean => {
  const backlogPath = join(repoRoot, ".hawp", "work", "BACKLOG.md");
  const backlogContent = readFileSync(backlogPath, "utf-8");
  return /\|\s*UUID\s*\|\s*Legacy ID\s*\|/i.test(backlogContent);
};

const buildVerificationResearchQueue = (
  root: string,
): VerificationResearchItem[] => {
  const closedRoot = join(root, ".hawp", "work", "closed");
  const closedFiles = walkMarkdownFiles(closedRoot);
  const clarity = checkVerificationClarity(closedFiles);

  return clarity.ambiguous.map((item) => ({
    itemId: item.id,
    claim: item.claim,
    filePath: item.filePath,
    lineNumber: item.lineNumber,
    recommendedAction:
      "Research concrete supporting evidence for this verification claim, then update the checklist line with Evidence: ... or mark it explicitly unproven.",
  }));
};

const summarizeValidation = (
  root: string,
  workRoot: string,
  notices: string[],
): void => {
  const backlogRows = parseValidationBacklog(join(workRoot, "BACKLOG.md"));

  if (backlogRows === null) {
    notices.push(
      "Validation warning: could not parse BACKLOG.md for workflow validation.",
    );
    return;
  }

  const validationReport = orchestrateValidation(workRoot, backlogRows);
  notices.push(
    `Validation summary: VALIDATION ${validationReport.overallStatus} (${validationReport.summary.failed} issues, ${validationReport.summary.warnings} warnings).`,
  );

  if (validationReport.checks.backlogConsistency.status === "FAIL") {
    notices.push(
      "Validation warning: working files are out of sync with BACKLOG.md; reconcile active/closed plan files before kit sync.",
    );
  }

  if (hasMixedIdColumns(root)) {
    notices.push(
      "Mixed UUID/legacy backlog detected: CLI validation still resolves working files by Legacy ID. Keep migration scoped to the CLI layer until apply/sync support lands.",
    );
  }
};

export const runBacklogUpgradeScript = (
  options: BacklogUpgradeScriptOptions,
): BacklogUpgradeScriptResult => {
  const notices: string[] = [];

  if (options.verbose) {
    notices.push(
      `Script options: mode=${options.mode}, format=${options.format}, validate=${options.validate}, forceDirty=${options.forceDirty}`,
    );
  }

  try {
    const root = findBacklogRepoRoot(options.repoRoot ?? process.cwd());

    if (options.mode === Mode.Apply) {
      if (!options.forceDirty && hasDirtyWorkingTree(root)) {
        return {
          exitCode: ExitCode.UsageError,
          stdoutText: "",
          stderrLines: [
            "Error: apply mode requires a clean working tree. Re-run with --force-dirty to override.",
          ],
          notices,
        };
      }

      const { changedFiles, skippedFiles, researchQueue } =
        applyClosedRecordNormalization(root);

      notices.push(
        `Applied closed-record normalization to ${changedFiles.length} file(s).`,
      );

      if (skippedFiles.length > 0) {
        notices.push(
          `Skipped ${skippedFiles.length} ambiguous legacy file(s) without inferable Backlog ID.`,
        );
      }
      if (researchQueue.length > 0) {
        notices.push(
          `Added ${researchQueue.length} verification evidence follow-up item(s) inside Verification sections for agent-friendly research handoff.`,
        );
      }

      let stdoutText =
        changedFiles.length > 0
          ? `${changedFiles.length} closed-record file(s) normalized.`
          : "No closed-record changes were necessary.";

      if (options.validate) {
        summarizeValidation(root, join(root, ".hawp", "work"), notices);
      }

      if (options.output) {
        writeOptionalFile(options.output, stdoutText);
        notices.push(`Report written: ${options.output}`);
        stdoutText = "";
      }

      if (options.exportResearchQueue) {
        writeOptionalFile(
          options.exportResearchQueue,
          `${JSON.stringify(researchQueue, null, 2)}\n`,
        );
        notices.push(
          `Research queue exported: ${options.exportResearchQueue}`,
        );
      }

      return {
        exitCode: ExitCode.Success,
        stdoutText:
          stdoutText.length > 0 ? withTrailingNewline(stdoutText) : "",
        stderrLines: [],
        notices,
      };
    }

    if (options.forceDirty) {
      notices.push("--force-dirty has no effect in dry-run mode.");
    }

    const { report } = runDetection(root);
    const reportWithResearchQueue = {
      ...report,
      researchQueue: buildVerificationResearchQueue(root),
    };
    const rendered =
      options.format === OutputFormat.Json
        ? renderJsonReport(reportWithResearchQueue)
        : renderTextReport(reportWithResearchQueue);

    let stdoutText = withTrailingNewline(rendered);

    if (options.validate) {
      summarizeValidation(root, join(root, ".hawp", "work"), notices);
    }

    if (options.output) {
      writeOptionalFile(options.output, stdoutText);
      notices.push(`Report written: ${options.output}`);
      stdoutText = "";
    }

    if (options.exportResearchQueue) {
      writeOptionalFile(
        options.exportResearchQueue,
        `${JSON.stringify(reportWithResearchQueue.researchQueue, null, 2)}\n`,
      );
      notices.push(`Research queue exported: ${options.exportResearchQueue}`);
    }

    if (options.exportPlan) {
      writeOptionalFile(
        options.exportPlan,
        `${JSON.stringify(reportWithResearchQueue.plan, null, 2)}\n`,
      );
      notices.push(`Plan exported: ${options.exportPlan}`);
    }

    return {
      exitCode: ExitCode.Success,
      stdoutText,
      stderrLines: [],
      notices,
    };
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);
    return {
      exitCode: ExitCode.Error,
      stdoutText: "",
      stderrLines: [`Script error: ${message}`],
      notices,
    };
  }
};
