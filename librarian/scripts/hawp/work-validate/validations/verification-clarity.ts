import { readFileSync } from "node:fs";
import type { VerificationCheck } from "../types";

/**
 * Checks verification clarity: evidence markers and unproven claims
 */
export function checkVerificationClarity(
  closedFiles: string[],
): VerificationCheck {
  const result: VerificationCheck = {
    total: 0,
    proven: 0,
    unproven: [],
    ambiguous: [],
    status: "PASS",
  };

  for (const filePath of closedFiles) {
    try {
      const content = readFileSync(filePath, "utf-8");
      const fileName =
        filePath.split("/").pop()?.replace(".md", "") || "unknown";

      // Find Verification section
      const verificationSection = extractVerificationSection(content);
      if (!verificationSection) continue;

      // Count lines with evidence markers
      const lines = verificationSection.split("\n");
      for (const [index, line] of lines.entries()) {
        if (line.includes("- [x]") || line.includes("- [ ]")) {
          const claim = line.replace(/^[\s-\[\]x ]+/, "").substring(0, 100);
          if (
            claim.startsWith("Research evidence for:") ||
            claim.startsWith(
              "Update the original verification checklist line with Evidence:",
            )
          ) {
            continue;
          }

          result.total++;
          const lineNumber = index + 1;

          if (line.includes("Evidence:")) {
            result.proven++;
          } else if (
            line.includes("NOT YET VERIFIED") ||
            /\b(?:explicitly )?unproven\b/i.test(line)
          ) {
            result.unproven.push({
              id: fileName,
              claim,
              filePath,
              lineNumber,
            });
          } else {
            result.ambiguous.push({
              id: fileName,
              claim,
              filePath,
              lineNumber,
            });
          }
        }
      }
    } catch (error) {
      console.error(
        `[validate] warning: skipping unreadable closed plan ${filePath}: ${
          error instanceof Error ? error.message : String(error)
        }`,
      );
    }
  }

  if (result.unproven.length > 0 || result.ambiguous.length > 0) {
    result.status = "WARN";
  }

  return result;
}

/**
 * Extract verification section from plan file content
 */
function extractVerificationSection(content: string): string | null {
  const lines = content.split("\n");
  let inVerification = false;
  const section: string[] = [];

  for (const line of lines) {
    // Accept both "## Verification" and annotated variants such as
    // "## Verification (filled at close)".
    if (/^##\s+Verification\b/.test(line.trim())) {
      inVerification = true;
      continue;
    }

    if (inVerification) {
      if (line.startsWith("## ")) {
        // End of section
        break;
      }
      section.push(line);
    }
  }

  return section.length > 0 ? section.join("\n") : null;
}
