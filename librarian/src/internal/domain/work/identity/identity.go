// Package identity contains shared value rules for HAWP work-record IDs.
package identity

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

func ExtractShortUUID(v string) string {
	if shortUUIDRe.MatchString(v) {
		return strings.ToLower(v)
	}
	return ""
}
func IDsMatch(a, b string) bool {
	l, r := strings.ToLower(a), strings.ToLower(b)
	if l == r {
		return true
	}
	if shortUUIDRe.MatchString(l) && fullUUIDRe.MatchString(r) {
		return strings.HasPrefix(r, l)
	}
	if shortUUIDRe.MatchString(r) && fullUUIDRe.MatchString(l) {
		return strings.HasPrefix(l, r)
	}
	return false
}
func ExtractIDFromFilename(v string) string {
	if m := uuidPrefixRe.FindStringSubmatch(v); m != nil {
		return strings.ToLower(m[1])
	}
	if m := legacyIDRe.FindString(v); m != "" {
		return m
	}
	if numericIDRe.MatchString(v) {
		return v
	}
	if m := datePrefixRe.FindStringSubmatch(v); m != nil {
		return strings.ToUpper(m[1])
	}
	if m := shortIDSlugRe.FindStringSubmatch(v); m != nil {
		return strings.ToLower(m[1])
	}
	if m := bareShortIDRe.FindStringSubmatch(v); m != nil {
		return strings.ToLower(m[1])
	}
	return ""
}
func IsFullUUID(v string) bool  { return fullUUIDRe.MatchString(v) }
func IsNumericID(v string) bool { return numericIDRe.MatchString(v) }
func MatchesAny(ids map[string]struct{}, id string) bool {
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
