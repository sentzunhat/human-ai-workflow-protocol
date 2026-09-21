package work

import "github.com/sentzunhat/hawp/librarian/src/internal/domain/work/identity"

// ExtractShortUUID returns the lowercase short-UUID display form when the
// whole value is exactly 8 hex chars; "" otherwise.
func ExtractShortUUID(value string) string { return identity.ExtractShortUUID(value) }

// IDsMatch reports whether two IDs are equal case-insensitively, or one is a
// short-UUID display form prefixing the other's full UUID.
func IDsMatch(a, b string) bool { return identity.IDsMatch(a, b) }

// ExtractIDFromFilename pulls a work item ID out of a filename or row cell:
// full UUIDs (lowercased), TASK-012-style prefixes, date-prefixed
// 2026-04-29-BUG-001-title forms (uppercased), short-ID-slug (8 alphanumeric
// chars + dash + letter, e.g. b7e2a4f9-rename-...), and bare 8-char
// alphanumeric (uuid folder names like b7e2a4f9). Returns "" when unrecognized.
func ExtractIDFromFilename(filename string) string { return identity.ExtractIDFromFilename(filename) }

// matchesAnyID reports whether any known ID matches id exactly or via
// short-UUID prefix rules.
func matchesAnyID(ids map[string]struct{}, id string) bool {
	return identity.MatchesAny(ids, id)
}
