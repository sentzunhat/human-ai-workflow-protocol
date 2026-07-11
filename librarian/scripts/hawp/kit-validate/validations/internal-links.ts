import { readFileSync, existsSync, statSync, readdirSync } from "node:fs";
import { join, dirname, resolve, relative } from "node:path";

import type { KitIssue } from "../script";

// Match [text](path) — relative links only
const LINK_RE = /\[([^\]]*)\]\(([^)]+)\)/g;
// Strip fenced code blocks (``` ... ```) before scanning for links
const FENCE_RE = /^```[\s\S]*?^```/gm;

const collectMarkdownFiles = (dir: string): string[] => {
  const files: string[] = [];
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    if (statSync(full).isDirectory()) {
      files.push(...collectMarkdownFiles(full));
    } else if (entry.endsWith(".md")) {
      files.push(full);
    }
  }
  return files;
};

export const checkInternalLinks = (kitPath: string): KitIssue[] => {
  const issues: KitIssue[] = [];
  for (const file of collectMarkdownFiles(kitPath)) {
    const raw = readFileSync(file, "utf-8");
    // blank out fenced code blocks so links inside them are not checked
    const content = raw.replace(FENCE_RE, (m) => m.replace(/./g, " "));
    const rel = relative(kitPath, file);
    let match: RegExpExecArray | null;
    LINK_RE.lastIndex = 0;
    while ((match = LINK_RE.exec(content)) !== null) {
      const href: string | undefined = match[2];
      if (href === undefined) continue;
      if (href.startsWith("http") || href.startsWith("/") || href.startsWith("#")) continue;
      const pathPart: string | undefined = href.split("#")[0];
      if (!pathPart) continue;
      const target = resolve(dirname(file), pathPart);
      if (!existsSync(target)) {
        issues.push({ file: rel, message: `broken link: ${href}` });
      }
    }
  }
  return issues;
};
