import { existsSync, readFileSync, statSync, readdirSync } from "node:fs";
import { join, resolve, sep } from "node:path";
import type { EvidenceCheck } from "../types";

const warn = (context: string, error: unknown): void => {
  console.error(
    `[validate] warning: ${context}: ${
      error instanceof Error ? error.message : String(error)
    }`,
  );
};

/**
 * Checks that evidence files referenced in verification sections exist
 */
export function checkEvidenceIntegrity(
  workDir: string,
  closedFiles: string[],
): EvidenceCheck {
  const result: EvidenceCheck = {
    total: 0,
    valid: 0,
    broken: [],
    status: "PASS",
  };

  for (const filePath of closedFiles) {
    try {
      const content = readFileSync(filePath, "utf-8");
      const fileName =
        filePath.split("/").pop()?.replace(".md", "") || "unknown";

      // Extract evidence links from verification section
      const links = extractEvidenceLinks(content, workDir);

      for (const link of links) {
        result.total++;

        // Check if evidence file exists
        if (existsSync(link.fullPath)) {
          result.valid++;
        } else {
          result.broken.push({
            id: fileName,
            link: link.relative,
          });
        }
      }
    } catch (error) {
      warn(`skipping unreadable closed plan ${filePath}`, error);
    }
  }

  if (result.broken.length > 0) {
    result.status = "WARN";
  }

  return result;
}

interface EvidenceLink {
  relative: string;
  fullPath: string;
}

/**
 * Extract evidence links from a plan file
 * Format: Evidence: inline or link to ../evidence/YYYY/MM/DD/<ID>-*.md
 * Links are resolved against <workDir>/evidence regardless of how deeply
 * the plan file is nested under closed/.
 */
function extractEvidenceLinks(content: string, workDir: string): EvidenceLink[] {
  const links: EvidenceLink[] = [];

  // Look for Evidence: ... links
  const lines = content.split("\n");
  for (const line of lines) {
    // Match pattern: Evidence: link to ../evidence/...
    const match = line.match(
      /Evidence:[\s]*(?:link to )?\.\.\/evidence\/([\w/.-]+\.md)/,
    );
    if (match && match[1]) {
      const relativePath = match[1];
      const evidenceRoot = resolve(workDir, "evidence");
      const fullPath = resolve(evidenceRoot, relativePath);

      // Reject links that escape the evidence folder via "..".
      if (!fullPath.startsWith(evidenceRoot + sep)) {
        console.error(
          `[validate] warning: evidence link escapes evidence folder, skipping: ${relativePath}`,
        );
        continue;
      }

      links.push({
        relative: `../evidence/${relativePath}`,
        fullPath,
      });
    }
  }

  return links;
}

/**
 * Collect all closed plan files
 */
export function collectClosedPlanFiles(closedDir: string): string[] {
  const files: string[] = [];

  try {
    const stat = statSync(closedDir);
    if (!stat.isDirectory()) return files;

    const years = readdirSync(closedDir);
    for (const year of years) {
      if (year === "README.md") continue;
      const yearPath = join(closedDir, year);
      const yearStat = statSync(yearPath);
      if (!yearStat.isDirectory()) continue;

      const months = readdirSync(yearPath);
      for (const month of months) {
        const monthPath = join(yearPath, month);
        const monthStat = statSync(monthPath);
        if (!monthStat.isDirectory()) continue;

        const days = readdirSync(monthPath);
        for (const day of days) {
          const dayPath = join(monthPath, day);
          const dayStat = statSync(dayPath);
          if (!dayStat.isDirectory()) continue;

          const planFiles = readdirSync(dayPath);
          for (const file of planFiles) {
            if (file.endsWith(".md")) {
              files.push(join(dayPath, file));
            }
          }
        }
      }
    }
  } catch (error) {
    warn(`failed to walk closed plans under ${closedDir}`, error);
  }

  return files;
}
