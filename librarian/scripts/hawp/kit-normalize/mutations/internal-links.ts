import { readFileSync, readdirSync, writeFileSync } from "node:fs";
import { dirname, join, relative, resolve } from "node:path";

export interface LinkUpdate {
  file: string;
  from: string;
  to: string;
  start: number;
  end: number;
}

const LINK_RE = /\[([^\]]*)\]\(([^)]+)\)/g;
const FENCE_RE = /^```[\s\S]*?^```/gm;

const collectMarkdownFiles = (dir: string): string[] => {
  const files: string[] = [];
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    const full = join(dir, entry.name);
    if (entry.isDirectory()) {
      files.push(...collectMarkdownFiles(full));
      continue;
    }

    if (entry.isFile() && entry.name.endsWith(".md")) {
      files.push(full);
    }
  }
  return files;
};

const maskFencedCodeBlocks = (content: string): string =>
  content.replace(FENCE_RE, (match) => " ".repeat(match.length));

export const planLinkUpdates = (
  kitPath: string,
  renameMap: Map<string, string>,
): LinkUpdate[] => {
  const updates: LinkUpdate[] = [];

  for (const file of collectMarkdownFiles(kitPath)) {
    const raw = readFileSync(file, "utf-8");
    const masked = maskFencedCodeBlocks(raw);
    const fileDir = dirname(file);
    let match: RegExpExecArray | null;

    LINK_RE.lastIndex = 0;
    while ((match = LINK_RE.exec(masked)) !== null) {
      const href = match[2];
      if (href === undefined || href.startsWith("http") || href.startsWith("/") || href.startsWith("#")) {
        continue;
      }

      const anchorIndex = href.indexOf("#");
      const pathPart = anchorIndex >= 0 ? href.slice(0, anchorIndex) : href;
      const anchor = anchorIndex >= 0 ? href.slice(anchorIndex) : "";
      if (!pathPart) {
        continue;
      }

      const target = resolve(fileDir, pathPart);
      const renamedTarget = renameMap.get(target);
      if (!renamedTarget) {
        continue;
      }

      const nextHref = `${relative(fileDir, renamedTarget).replace(/\\/g, "/")}${anchor}`;
      if (nextHref !== href) {
        updates.push({
          file,
          from: href,
          to: nextHref,
          start: match.index + match[0].indexOf(href),
          end: match.index + match[0].indexOf(href) + href.length,
        });
      }
    }
  }

  return updates;
};

export const applyLinkUpdates = (
  updates: LinkUpdate[],
): number => {
  if (updates.length === 0) {
    return 0;
  }

  const perFileEdits = new Map<string, LinkUpdate[]>();
  for (const update of updates) {
    const current = perFileEdits.get(update.file) ?? [];
    current.push(update);
    perFileEdits.set(update.file, current);
  }

  let changedFiles = 0;
  for (const [file, fileUpdates] of perFileEdits.entries()) {
    const raw = readFileSync(file, "utf-8");
    let next = raw;
    for (const update of fileUpdates.sort((a, b) => b.start - a.start)) {
      next = `${next.slice(0, update.start)}${update.to}${next.slice(update.end)}`;
    }
    if (next !== raw) {
      writeFileSync(file, next, "utf-8");
      changedFiles += 1;
    }
  }

  return changedFiles;
};
