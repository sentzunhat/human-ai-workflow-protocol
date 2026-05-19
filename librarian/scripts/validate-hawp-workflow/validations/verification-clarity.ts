import { readFileSync } from "fs";
import type { VerificationCheck } from "../types";

/**
 * Checks verification clarity: evidence markers and unproven claims
 */
export async function checkVerificationClarity(
  closedFiles: string[],
): Promise<VerificationCheck> {
  const result: VerificationCheck = {
    total: 0,
    proven: 0,
    unproven: [],
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
      for (const line of lines) {
        if (line.includes("- [x]") || line.includes("- [ ]")) {
          result.total++;

          if (line.includes("NOT YET VERIFIED")) {
            result.unproven.push({
              id: fileName,
              claim: line.replace(/^[\s-\[\]x ]+/, "").substring(0, 100),
            });
          } else if (line.includes("Evidence:") || line.includes("unproven")) {
            result.proven++;
          }
        }
      }
    } catch {
      // Ignore read errors
    }
  }

  if (result.unproven.length > 0) {
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
    if (line.trim() === "## Verification (filled at close)") {
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
