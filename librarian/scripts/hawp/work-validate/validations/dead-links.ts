import { existsSync, readFileSync, readdirSync } from "node:fs";
import { dirname, join, relative, resolve } from "node:path";

import type { DeadLinksCheck } from "../types";

const LINK_RE = /\[([^\]]*)\]\(([^)]+)\)/g;
const FENCE_RE = /^```[\s\S]*?^```/gm;

// Only scan active work — archives (closed/, evidence/, notes/, status/) may
// reference old paths and are not expected to have live links.
const ACTIVE_DIRS = ["active", "parked"];
const ACTIVE_ROOT_FILES = ["BACKLOG.md"];

const collectScanTargets = (workDir: string): string[] => {
  const files: string[] = [];

  for (const name of ACTIVE_ROOT_FILES) {
    const full = join(workDir, name);
    if (existsSync(full)) files.push(full);
  }

  for (const dir of ACTIVE_DIRS) {
    const full = join(workDir, dir);
    if (!existsSync(full)) continue;
    for (const entry of readdirSync(full, { withFileTypes: true })) {
      if (entry.isFile() && entry.name.endsWith(".md")) {
        files.push(join(full, entry.name));
      }
    }
  }

  return files;
};

export const checkDeadLinks = (workDir: string): DeadLinksCheck => {
  const files = collectScanTargets(workDir);
  const broken: Array<{ file: string; link: string }> = [];

  for (const file of files) {
    const raw = readFileSync(file, "utf-8");
    const content = raw.replace(FENCE_RE, (m) => " ".repeat(m.length));
    const rel = relative(workDir, file);

    LINK_RE.lastIndex = 0;
    let match: RegExpExecArray | null;
    while ((match = LINK_RE.exec(content)) !== null) {
      const href: string | undefined = match[2];
      if (href === undefined) continue;
      if (href.startsWith("http") || href.startsWith("/") || href.startsWith("#")) continue;
      const pathPart: string | undefined = href.split("#")[0];
      if (!pathPart) continue;
      const target = resolve(dirname(file), pathPart);
      if (!existsSync(target)) {
        broken.push({ file: rel, link: href });
      }
    }
  }

  return {
    scanned: files.length,
    broken,
    status: broken.length > 0 ? "FAIL" : "PASS",
  };
};
