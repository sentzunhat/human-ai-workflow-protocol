// Isolated ID parser for HAWP workflow validation
// Supports current TASK-NNN, BUG-NNN format
// Designed for extensibility to UUID-based IDs in future

/**
 * Extracts ID from a filename or backlog row
 * Current format: TASK-012, BUG-005, etc.
 * Future format: UUID-based IDs can be added here without changing validation logic
 *
 * @param filename - The filename or row text to extract ID from
 * @returns The extracted ID or null if not recognized
 */
export function extractIdFromFilename(filename: string): string | null {
  // Standard prefix: TASK-012, BUG-005, BUG-063-some-title
  const match = filename.match(/^([A-Z]+)-(\d+)/);
  if (match) {
    return match[0]; // Return full ID like "TASK-012"
  }

  // Date-prefixed: 2026-04-29-BUG-001-title or 2026-05-01-bug-005-title
  // Date-prefixed: 2026-04-29-BUG-001-title  (TYPE-NNN must follow the date directly)
  const dateMatch = filename.match(/^\d{4}-\d{2}-\d{2}-([A-Za-z]+-\d+)/);
  if (dateMatch) {
    return dateMatch[1]!.toUpperCase();
  }

  // Future: UUID format support can be added here
  // Example: const uuidMatch = filename.match(/^([a-f0-9]{8}-[a-f0-9]{4})/);

  return null;
}

/**
 * Extracts ID from a backlog row line
 * Format: | ID | Type | Title | ...
 *
 * @param line - The backlog row line
 * @returns The extracted ID or null
 */
export function extractIdFromBacklogRow(line: string): string | null {
  const parts = line.split("|");
  if (parts.length < 2) return null;

  const idPart = parts[1]?.trim();
  if (!idPart) return null;

  // Try current format first
  const id = extractIdFromFilename(idPart);
  if (id) return id;

  // If no match, return the raw value for debugging
  return idPart;
}

/**
 * Checks if an ID is a valid sequential type-prefixed ID
 * @param id - The ID to validate
 * @returns true if valid
 */
export function isValidSequentialId(id: string): boolean {
  return /^[A-Z]+-\d+$/.test(id);
}

/**
 * Gets the type prefix from an ID
 * Example: "TASK-012" -> "TASK", "BUG-005" -> "BUG"
 *
 * @param id - The ID to extract type from
 * @returns The type prefix or null
 */
export function getTypePrefix(id: string): string | null {
  const match = id.match(/^([A-Z]+)-\d+$/);
  return match ? (match[1] ?? null) : null;
}

/**
 * Gets the numeric suffix from an ID
 * Example: "TASK-012" -> "012", "BUG-005" -> "005"
 *
 * @param id - The ID to extract number from
 * @returns The numeric suffix as string or null
 */
export function getNumericSuffix(id: string): string | null {
  const match = id.match(/^[A-Z]+-(\d+)$/);
  return match ? (match[1] ?? null) : null;
}
