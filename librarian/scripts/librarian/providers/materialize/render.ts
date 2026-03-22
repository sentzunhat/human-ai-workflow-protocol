import { existsSync, readFileSync } from "node:fs";
import { join } from "node:path";
import { normalizeForCompare } from "../../../lib";
import type { FrontmatterValue, MaterializationResult } from "./composition";
import {
  GENERATED_BANNER,
  MATERIALIZATION_TARGETS,
} from "./composition";

export { normalizeForCompare };

function serializeFrontmatter(
  frontmatter: Record<string, FrontmatterValue>,
): string {
  const lines = ["---"];

  for (const [key, value] of Object.entries(frontmatter)) {
    if (Array.isArray(value)) {
      lines.push(`${key}:`);
      for (const item of value) {
        lines.push(`  - ${JSON.stringify(item)}`);
      }
    } else if (typeof value === "boolean") {
      lines.push(`${key}: ${value}`);
    } else {
      lines.push(`${key}: ${JSON.stringify(value)}`);
    }
  }

  lines.push("---");
  return lines.join("\n");
}

function trimBody(body: string): string {
  return body.replace(/\r\n/g, "\n").trimEnd();
}

export function computeMaterializedOutputs(
  repoRoot: string,
): MaterializationResult[] {
  const sharedRoot = join(repoRoot, "core/providers/shared/behaviors");

  return MATERIALIZATION_TARGETS.map((target) => {
    const behaviorPath = join(sharedRoot, `${target.behavior}.md`);
    if (!existsSync(behaviorPath)) {
      throw new Error(`Missing shared behavior: ${behaviorPath}`);
    }

    const body = trimBody(readFileSync(behaviorPath, "utf-8"));
    const preamble = target.preamble ?? "";
    const content = [
      serializeFrontmatter(target.frontmatter),
      "",
      GENERATED_BANNER.trimEnd(),
      "",
      preamble + body,
      "",
    ].join("\n");

    return {
      outputPath: join(repoRoot, target.outputPath),
      content,
    };
  });
}
