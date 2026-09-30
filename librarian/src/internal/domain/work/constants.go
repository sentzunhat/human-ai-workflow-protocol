package work

import "github.com/sentzunhat/hawp/librarian/src/internal/domain/work/validation"

// LegacyClosedCutoff: closed files on or after this date require Outcome,
// Verification, and Close Checklist sections; earlier files are legacy.
const LegacyClosedCutoff = validation.LegacyClosedCutoff
