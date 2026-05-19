import { existsSync, readdirSync, readFileSync } from "node:fs";
import { join } from "node:path";

export interface PlanFileRecord {
  id: string;
  path: string;
  content: string;
}

const readBacklogId = (path: string): string | undefined => {
  const content = readFileSync(path, "utf-8");
  const match = content.match(/\*\*Backlog ID:\*\*\s*([A-Z]+-\d+)/i);
  return match?.[1]?.toUpperCase();
};

const readIdFromFilename = (path: string): string | undefined => {
  const fileName = path.split("/").pop() ?? "";
  const match = fileName.match(/(TASK|BUG)-\d+/i);
  return match?.[0]?.toUpperCase();
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

    if (!entry.isFile()) {
      continue;
    }

    if (!entry.name.endsWith(".md") || entry.name === "README.md") {
      continue;
    }

    paths.push(absolutePath);
  }

  return paths;
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
