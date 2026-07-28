// Package uuidgen generates work item UUIDs (v4) matching the behavior of
// the .hawp/bin/hawp uuid command.
package uuidgen

import (
	"crypto/rand"
	"fmt"
)

// New returns a random UUID v4 in lowercase canonical form.
func New() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // RFC 4122 variant
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

// Short returns the 8-character display form of a full UUID.
func Short(uuid string) string {
	if len(uuid) < 8 {
		return uuid
	}
	return uuid[:8]
}
