import { readdirSync, statSync } from "node:fs";
import { dirname, join } from "node:path";

export interface FileRename {
  from: string;
  to: string;
}

const ALLOWED_EXACT_NAMES = new Set(["README.md"]);

const normalizeStem = (stem: string): string =>
  stem
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/-+/g, "-")
    .replace(/^-+|-+$/g, "");

export const normalizeFileName = (fileName: string): string | null => {
  if (ALLOWED_EXACT_NAMES.has(fileName)) {
    return null;
  }

  const lastDot = fileName.lastIndexOf(".");
  const hasExtension = lastDot > 0 && lastDot < fileName.length - 1;
  const stem = hasExtension ? fileName.slice(0, lastDot) : fileName;
  const extension = hasExtension ? fileName.slice(lastDot).toLowerCase() : "";
  const normalizedStem = normalizeStem(stem);

  if (!normalizedStem) {
    return null;
  }

  const normalized = `${normalizedStem}${extension}`;
  return normalized === fileName ? null : normalized;
};

const walk = (dir: string, renames: FileRename[]): void => {
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    if (statSync(full).isDirectory()) {
      walk(full, renames);
      continue;
    }
    const normalized = normalizeFileName(entry);
    if (normalized) {
      renames.push({ from: full, to: join(dirname(full), normalized) });
    }
  }
};

export const planFileRenames = (kitPath: string): FileRename[] => {
  const renames: FileRename[] = [];
  walk(kitPath, renames);
  return renames;
};
