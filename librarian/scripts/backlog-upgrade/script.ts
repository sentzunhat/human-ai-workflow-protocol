import { execSync } from "node:child_process";
import {
  existsSync,
  mkdirSync,
  renameSync,
  readdirSync,
  readFileSync,
  writeFileSync,
} from "node:fs";
import { basename, dirname, join, resolve } from "node:path";

import { runDetection } from "./detection";
import { ExitCode, Mode, OutputFormat } from "./models/index.js";
import { renderJsonReport, renderTextReport } from "./output/formatters";
import {
  orchestrateValidation,
  parseBacklog as parseValidationBacklog,
} from "../validate-hawp-workflow/orchestrate";

export interface BacklogUpgradeScriptOptions {
  mode: Mode;
  validate: boolean;
  format: OutputFormat;
  output?: string;
  exportPlan?: string;
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

const normalizeClosedRecord = (content: string, filePath: string): string => {
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

  return updated;
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

const walkMarkdownFiles = (dirPath: string): string[] => {
  if (!existsSync(dirPath)) {
    return [];
  }

  const paths: string[] = [];
  for (const entry of readdirSync(dirPath, { withFileTypes: true })) {
    const absolutePath = join(dirPath, entry.name);
    if (entry.isDirectory()) {
      paths.push(...walkMarkdownFiles(absolutePath));
      continue;
    }

    if (
      !entry.isFile() ||
      !entry.name.endsWith(".md") ||
      entry.name === "README.md"
    ) {
      continue;
    }

    paths.push(absolutePath);
  }

  return paths;
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
    return false;
  }
};

const applyClosedRecordNormalization = (
  root: string,
): { changedFiles: string[]; skippedFiles: string[] } => {
  const workRoot = join(root, ".hawp", "work");
  const closedRoot = join(workRoot, "closed");
  const skippedFiles: string[] = [];
  const touchedFiles = new Set<string>();

  for (const absolutePath of walkMarkdownFiles(closedRoot)) {
    const reconciled = reconcileClosedRecordPath(root, absolutePath);
    const currentPath = reconciled.path;

    if (reconciled.moved) {
      touchedFiles.add(currentPath);
    }

    const current = readFileSync(currentPath, "utf-8");
    const next = normalizeClosedRecord(current, currentPath);

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
  }

  return { changedFiles: Array.from(touchedFiles), skippedFiles };
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

const resolveRepoRoot = (startDir: string): string => {
  let current = startDir;

  for (let depth = 0; depth < 10; depth += 1) {
    if (existsSync(join(current, ".hawp", "work", "BACKLOG.md"))) {
      return current;
    }

    const parent = dirname(current);
    if (parent === current) {
      break;
    }

    current = parent;
  }

  throw new Error(
    "Could not locate repo root containing .hawp/work/BACKLOG.md",
  );
};

export const runBacklogUpgradeScript = (
  options: BacklogUpgradeScriptOptions,
): Promise<BacklogUpgradeScriptResult> => {
  return (async () => {
    const notices: string[] = [];

    if (options.verbose) {
      notices.push(
        `Script options: mode=${options.mode}, format=${options.format}, validate=${options.validate}, forceDirty=${options.forceDirty}`,
      );
    }

    try {
      const root = resolveRepoRoot(options.repoRoot ?? process.cwd());

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

        const { changedFiles, skippedFiles } =
          applyClosedRecordNormalization(root);

        notices.push(
          `Applied closed-record normalization to ${changedFiles.length} file(s).`,
        );

        if (skippedFiles.length > 0) {
          notices.push(
            `Skipped ${skippedFiles.length} ambiguous legacy file(s) without inferable Backlog ID.`,
          );
        }

        let stdoutText =
          changedFiles.length > 0
            ? `${changedFiles.length} closed-record file(s) normalized.`
            : "No closed-record changes were necessary.";

        if (options.validate) {
          const workRoot = join(root, ".hawp", "work");
          const backlogRows = parseValidationBacklog(
            join(workRoot, "BACKLOG.md"),
          );

          if (backlogRows === null) {
            notices.push(
              "Validation warning: could not parse BACKLOG.md for workflow validation.",
            );
          } else {
            const validationReport = await orchestrateValidation(
              workRoot,
              backlogRows,
            );
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
          }
        }

        if (options.output) {
          writeOptionalFile(options.output, stdoutText);
          notices.push(`Report written: ${options.output}`);
          stdoutText = "";
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
      const rendered =
        options.format === OutputFormat.Json
          ? renderJsonReport(report)
          : renderTextReport(report);

      let stdoutText = withTrailingNewline(rendered);

      if (options.validate) {
        const workRoot = join(root, ".hawp", "work");
        const backlogRows = parseValidationBacklog(
          join(workRoot, "BACKLOG.md"),
        );

        if (backlogRows === null) {
          notices.push(
            "Validation warning: could not parse BACKLOG.md for workflow validation.",
          );
        } else {
          const validationReport = await orchestrateValidation(
            workRoot,
            backlogRows,
          );
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
        }
      }

      if (options.output) {
        writeOptionalFile(options.output, stdoutText);
        notices.push(`Report written: ${options.output}`);
        stdoutText = "";
      }

      if (options.exportPlan) {
        writeOptionalFile(
          options.exportPlan,
          `${JSON.stringify(report.plan, null, 2)}\n`,
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
  })();
};
