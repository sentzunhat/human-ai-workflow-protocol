import { existsSync, readFileSync } from "node:fs";
import { join } from "node:path";

import { walkMarkdownFiles } from "../../../lib";
import {
  extractIdFromFilename,
  extractShortUuid,
} from "../../work-validate/validations/id-parser";

export interface PlanFileRecord {
  id: string;
  path: string;
  content: string;
}

const readBacklogId = (path: string): string | undefined => {
  const content = readFileSync(path, "utf-8");
  const legacyMatch = content.match(/\*\*Backlog ID(?: \(Legacy\))?:\*\*\s*([A-Z]+-\d+)/i);
  if (legacyMatch?.[1]) {
    return legacyMatch[1].toUpperCase();
  }

  const uuidMatch = content.match(
    /\*\*UUID:\*\*\s*`?([0-9a-f]{8}(?:-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})?)`?/i,
  );
  if (uuidMatch?.[1]) {
    return extractIdFromFilename(uuidMatch[1]) ?? extractShortUuid(uuidMatch[1]) ?? undefined;
  }

  return undefined;
};

const readIdFromFilename = (path: string): string | undefined => {
  const fileName = path.split("/").pop() ?? "";
  return extractIdFromFilename(fileName) ?? extractShortUuid(fileName) ?? undefined;
};

export interface PlanScanResult {
  files: PlanFileRecord[];
  byId: Map<string, string[]>;
  directoryPresence: {
    active: boolean;
    parked: boolean;
    closed: boolean;
  };
}

export const scanPlanFiles = (workRoot: string): PlanScanResult => {
  const activeDir = join(workRoot, "active");
  const parkedDir = join(workRoot, "parked");
  const closedDir = join(workRoot, "closed");

  const active = walkMarkdownFiles(activeDir);
  const parked = walkMarkdownFiles(parkedDir);
  const closed = walkMarkdownFiles(closedDir);

  const files: PlanFileRecord[] = [];
  const byId = new Map<string, string[]>();

  for (const path of [...active, ...parked, ...closed]) {
    const content = readFileSync(path, "utf-8");
    const id = readBacklogId(path) ?? readIdFromFilename(path);
    if (!id) {
      continue;
    }

    files.push({ id, path, content });
    const existing = byId.get(id) ?? [];
    existing.push(path);
    byId.set(id, existing);
  }

  return {
    files,
    byId,
    directoryPresence: {
      active: existsSync(activeDir),
      parked: existsSync(parkedDir),
      closed: existsSync(closedDir),
    },
  };
};
