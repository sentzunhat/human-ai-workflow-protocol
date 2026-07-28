// Isolated ID parser for HAWP workflow validation
// Supports legacy TASK-NNN / BUG-NNN format and UUID-based IDs

const UUID_PATTERN =
  /^([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})/i;

const SHORT_UUID_PATTERN = /^[0-9a-f]{8}$/i;
const FULL_UUID_PATTERN =
  /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

/**
 * Extracts a short-UUID display form (exactly 8 hex chars) from a backlog cell.
 * Only whole-cell matches count — an 8-hex prefix inside a longer word is not an ID.
 */
export function extractShortUuid(value: string): string | null {
  return SHORT_UUID_PATTERN.test(value) ? value.toLowerCase() : null;
}

/**
 * Compares two work item IDs. True when they are equal (case-insensitive),
 * or when one side is a short-UUID display form (8 hex chars) and the other
 * is a full UUID starting with it.
 */
export function idsMatch(a: string, b: string): boolean {
  const left = a.toLowerCase();
  const right = b.toLowerCase();
  if (left === right) return true;
  if (SHORT_UUID_PATTERN.test(left) && FULL_UUID_PATTERN.test(right)) {
    return right.startsWith(left);
  }
  if (SHORT_UUID_PATTERN.test(right) && FULL_UUID_PATTERN.test(left)) {
    return left.startsWith(right);
  }
  return false;
}

/**
 * Extracts ID from a filename or backlog row
 * Legacy format: TASK-012, BUG-005, etc.
 * UUID format: full UUID v4, normalized to lowercase (e.g. 361fb08e-6457-4ed5-80bd-76337b6f0e89)
 *
 * @param filename - The filename or row text to extract ID from
 * @returns The extracted ID or null if not recognized
 */
export function extractIdFromFilename(filename: string): string | null {
  // UUID-based ID
  const uuidMatch = filename.match(UUID_PATTERN);
  if (uuidMatch) {
    return uuidMatch[1]!.toLowerCase();
  }

  // Standard prefix: TASK-012, BUG-005, BUG-063-some-title
  const match = filename.match(/^([A-Z]+)-(\d+)/);
  if (match) {
    return match[0]; // Return full ID like "TASK-012"
  }

  // Date-prefixed: 2026-04-29-BUG-001-title  (TYPE-NNN must follow the date directly)
  const dateMatch = filename.match(/^\d{4}-\d{2}-\d{2}-([A-Za-z]+-\d+)/);
  if (dateMatch) {
    return dateMatch[1]!.toUpperCase();
  }

  // Short-ID prefix + descriptive slug: c9a7f2e1-github-actions-pipeline,
  // h5f7c2j8-retry-v001-release. These 8-char row IDs are arbitrary
  // alphanumeric tokens, not necessarily valid hex (distinct from the full
  // UUID case, which UUID_PATTERN above already handles).
  const shortIdSlugMatch = filename.match(/^([0-9a-z]{8})-[a-z]/i);
  if (shortIdSlugMatch) {
    return shortIdSlugMatch[1]!.toLowerCase();
  }

  return null;
}
