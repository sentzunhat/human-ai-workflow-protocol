package work

import (
	"regexp"
	"strings"
)

var (
	uuidPrefixRe  = regexp.MustCompile(`(?i)^([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})`)
	shortUUIDRe   = regexp.MustCompile(`(?i)^[0-9a-f]{8}$`)
	fullUUIDRe    = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	legacyIDRe    = regexp.MustCompile(`^([A-Z]+)-(\d+)`)
	numericIDRe   = regexp.MustCompile(`^\d+$`)
	datePrefixRe  = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}-([A-Za-z]+-\d+)`)
	shortIDSlugRe = regexp.MustCompile(`(?i)^([0-9a-z]{8})-[a-z]`)
	bareShortIDRe = regexp.MustCompile(`(?i)^([0-9a-z]{8})$`)
)

// ExtractShortUUID returns the lowercase short-UUID display form when the
// whole value is exactly 8 hex chars; "" otherwise.
func ExtractShortUUID(value string) string {
	if shortUUIDRe.MatchString(value) {
		return strings.ToLower(value)
	}
	return ""
}

// IDsMatch reports whether two IDs are equal case-insensitively, or one is a
// short-UUID display form prefixing the other's full UUID.
func IDsMatch(a, b string) bool {
	left, right := strings.ToLower(a), strings.ToLower(b)
	if left == right {
		return true
	}
	if shortUUIDRe.MatchString(left) && fullUUIDRe.MatchString(right) {
		return strings.HasPrefix(right, left)
	}
	if shortUUIDRe.MatchString(right) && fullUUIDRe.MatchString(left) {
		return strings.HasPrefix(left, right)
	}
	return false
}

// ExtractIDFromFilename pulls a work item ID out of a filename or row cell:
// full UUIDs (lowercased), TASK-012-style prefixes, date-prefixed
// 2026-04-29-BUG-001-title forms (uppercased), short-ID-slug (8 alphanumeric
// chars + dash + letter, e.g. b7e2a4f9-rename-...), and bare 8-char
// alphanumeric (uuid folder names like b7e2a4f9). Returns "" when unrecognized.
func ExtractIDFromFilename(filename string) string {
	if m := uuidPrefixRe.FindStringSubmatch(filename); m != nil {
		return strings.ToLower(m[1])
	}
	if m := legacyIDRe.FindString(filename); m != "" {
		return m
	}
	if numericIDRe.MatchString(filename) {
		return filename
	}
	if m := datePrefixRe.FindStringSubmatch(filename); m != nil {
		return strings.ToUpper(m[1])
	}
	if m := shortIDSlugRe.FindStringSubmatch(filename); m != nil {
		return strings.ToLower(m[1])
	}
	if m := bareShortIDRe.FindStringSubmatch(filename); m != nil {
		return strings.ToLower(m[1])
	}
	return ""
}

// matchesAnyID reports whether any known ID matches id exactly or via
// short-UUID prefix rules.
func matchesAnyID(ids map[string]struct{}, id string) bool {
	if _, ok := ids[id]; ok {
		return true
	}
	for known := range ids {
		if IDsMatch(known, id) {
			return true
		}
	}
	return false
}
