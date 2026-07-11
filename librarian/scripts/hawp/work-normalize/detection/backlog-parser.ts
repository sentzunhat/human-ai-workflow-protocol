import { readFileSync } from "node:fs";
import {
  extractIdFromFilename,
  extractShortUuid,
} from "../../work-validate/validations/id-parser";

export type BacklogSection = "active" | "blocked" | "recently-closed" | "other";

export interface BacklogRow {
  section: BacklogSection;
  lineNumber: number;
  id: string;
  type: string;
  title: string;
  status?: string;
  updated?: string;
  reason?: string;
  planPath?: string;
}

const trim = (value: string): string => value.trim();

const extractMarkdownLink = (value: string): string | undefined => {
  const match = value.match(/\(([^)]+)\)/);
  return match?.[1]?.trim();
};

const normalizeSection = (line: string): BacklogSection => {
  if (line.startsWith("## Active Work")) {
    return "active";
  }
  if (line.startsWith("## Blocked / Parked")) {
    return "blocked";
  }
  if (line.startsWith("## Recently Closed")) {
    return "recently-closed";
  }
  return "other";
};

const parseTableCells = (line: string): string[] =>
  line.split("|").slice(1, -1).map(trim);

const normalizeHeader = (value: string): string => value.trim().toLowerCase();

const stripCodeSpan = (value: string): string => value.replace(/^`+|`+$/g, "");

const createHeaderMap = (cells: string[]): Map<string, number> =>
  new Map(cells.map((value, index) => [normalizeHeader(value), index]));

const getMappedCell = (
  cells: string[],
  headerMap: Map<string, number>,
  ...aliases: string[]
): string => {
  for (const alias of aliases) {
    const index = headerMap.get(alias);
    if (index !== undefined) {
      return cells[index] ?? "";
    }
  }

  return "";
};

export interface BacklogParseResult {
  backlogPath: string;
  rows: BacklogRow[];
  sectionPresence: Record<BacklogSection, boolean>;
}

const withOptionalString = (
  key: keyof Pick<BacklogRow, "status" | "updated" | "reason" | "planPath">,
  value: string | undefined,
): Partial<BacklogRow> => {
  if (!value) {
    return {};
  }

  return { [key]: value };
};

export const parseBacklog = (backlogPath: string): BacklogParseResult => {
  const content = readFileSync(backlogPath, "utf-8");
  const lines = content.split(/\r?\n/);

  const rows: BacklogRow[] = [];
  let section: BacklogSection = "other";
  let headerMap: Map<string, number> | null = null;
  const sectionPresence: Record<BacklogSection, boolean> = {
    active: false,
    blocked: false,
    "recently-closed": false,
    other: false,
  };

  lines.forEach((line, index) => {
    if (line.startsWith("## ")) {
      section = normalizeSection(line);
      headerMap = null;
      sectionPresence[section] = true;
      return;
    }

    if (!line.startsWith("|")) {
      return;
    }

    const cells = parseTableCells(line);
    if (cells.length < 4) {
      return;
    }

    if (headerMap === null) {
      headerMap = createHeaderMap(cells);
      return;
    }

    if (/^\|\s*-+/.test(line) || line.includes("| --------")) {
      return;
    }

    const currentHeaderMap = headerMap;
    const idCandidates = ["legacy id", "id", "uuid"]
      .map((alias) =>
        stripCodeSpan(getMappedCell(cells, currentHeaderMap, alias)),
      )
      .filter((value) => value !== "" && value !== "—" && value !== "-");
    const id =
      idCandidates.find((value) => extractIdFromFilename(value) !== null) ??
      idCandidates.find((value) => extractShortUuid(value) !== null) ??
      idCandidates[0];
    if (
      !id ||
      normalizeHeader(id) === "id" ||
        normalizeHeader(id) === "legacy id"
    ) {
      return;
    }

    const normalizedId =
      extractIdFromFilename(id) ?? extractShortUuid(id) ?? id;

    if (section === "active") {
      const planPath = extractMarkdownLink(
        getMappedCell(cells, headerMap, "plan file", "detail"),
      );
      rows.push({
        section,
        lineNumber: index + 1,
        id: normalizedId,
        type: getMappedCell(cells, currentHeaderMap, "type"),
        title: getMappedCell(cells, currentHeaderMap, "title"),
        ...withOptionalString(
          "status",
          getMappedCell(cells, currentHeaderMap, "status"),
        ),
        ...withOptionalString("planPath", planPath),
        ...withOptionalString(
          "updated",
          getMappedCell(cells, currentHeaderMap, "updated"),
        ),
      });
      return;
    }

    if (section === "blocked") {
      const planPath = extractMarkdownLink(
        getMappedCell(cells, headerMap, "detail", "plan file"),
      );
      rows.push({
        section,
        lineNumber: index + 1,
        id: normalizedId,
        type: getMappedCell(cells, currentHeaderMap, "type"),
        title: getMappedCell(cells, currentHeaderMap, "title"),
        ...withOptionalString(
          "reason",
          getMappedCell(cells, currentHeaderMap, "reason", "status"),
        ),
        ...withOptionalString("planPath", planPath),
        ...withOptionalString(
          "updated",
          getMappedCell(cells, currentHeaderMap, "updated"),
        ),
      });
      return;
    }

    if (section === "recently-closed") {
      const planPath = extractMarkdownLink(
        getMappedCell(cells, headerMap, "detail", "plan file"),
      );
      rows.push({
        section,
        lineNumber: index + 1,
        id: normalizedId,
        type: getMappedCell(cells, currentHeaderMap, "type"),
        title: getMappedCell(cells, currentHeaderMap, "title"),
        ...withOptionalString(
          "updated",
          getMappedCell(cells, currentHeaderMap, "closed", "updated"),
        ),
        ...withOptionalString("planPath", planPath),
      });
      return;
    }

    const planPath = extractMarkdownLink(
      getMappedCell(cells, headerMap, "detail", "plan file"),
    );

    rows.push({
      section,
      lineNumber: index + 1,
      id: normalizedId,
      type: getMappedCell(cells, currentHeaderMap, "type"),
      title: getMappedCell(cells, currentHeaderMap, "title"),
      ...withOptionalString(
        "status",
        getMappedCell(cells, currentHeaderMap, "status", "reason", "closed"),
      ),
      ...withOptionalString("planPath", planPath),
      ...withOptionalString(
        "updated",
        getMappedCell(cells, currentHeaderMap, "updated", "closed"),
      ),
    });
  });

  return { backlogPath, rows, sectionPresence };
};
