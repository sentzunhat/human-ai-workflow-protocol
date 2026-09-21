// Package validation owns reusable work-record validation policies.
package validation

import (
	"fmt"
	"github.com/sentzunhat/hawp/librarian/src/internal/domain/work/identity"
	"github.com/sentzunhat/hawp/librarian/src/internal/domain/work/model"
	"io/fs"
)

// Source supplies filesystem observations without binding validation policy
// to a concrete operating-system adapter.
type Source struct {
	Exists         func(path string) bool
	ReadDir        func(path string) ([]fs.DirEntry, error)
	ReadFile       func(path string) ([]byte, error)
	Stat           func(path string) (fs.FileInfo, error)
	ToRepoRelative func(repoRoot, absolutePath string) string
	Warn           func(format string, args ...any)
}

func warnf(source Source, format string, args ...any) {
	if source.Warn != nil {
		source.Warn(format, args...)
		return
	}
	_ = fmt.Sprintf(format, args...)
}

type CheckStatus = model.CheckStatus
type BacklogRow = model.BacklogRow
type Backlog = model.Backlog
type BacklogCheck = model.BacklogCheck
type FileFinding = model.FileFinding
type ClosedTaskCheck = model.ClosedTaskCheck
type EvidenceCheck = model.EvidenceCheck
type BrokenLink = model.BrokenLink
type Claim = model.Claim
type VerificationCheck = model.VerificationCheck
type DeadLinksCheck = model.DeadLinksCheck

const (
	StatusPass         = model.StatusPass
	StatusFail         = model.StatusFail
	StatusWarn         = model.StatusWarn
	LegacyClosedCutoff = "2026-05-10"
)

func ExtractShortUUID(value string) string                 { return identity.ExtractShortUUID(value) }
func ExtractIDFromFilename(value string) string            { return identity.ExtractIDFromFilename(value) }
func IDsMatch(a, b string) bool                            { return identity.IDsMatch(a, b) }
func matchesAnyID(ids map[string]struct{}, id string) bool { return identity.MatchesAny(ids, id) }
