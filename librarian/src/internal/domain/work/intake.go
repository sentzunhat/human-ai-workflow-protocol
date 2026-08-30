package work

import (
	"fmt"
	"strings"
)

// NewItemInput describes a new work item to scaffold via `hawp work new`.
// This mechanically shapes the intake step's boilerplate only — per HAWP's
// "shaping protocol, not a runtime" boundary (see AGENTS.md), the actual
// investigation and plan content are left for a human or AI agent to fill
// in afterward; this just removes the chore of hand-writing the backlog row
// and plan file skeleton.
type NewItemInput struct {
	UUID  string // full UUID (caller generates, e.g. via uuidgen.New())
	Type  string // task | bug | improvement | feature | fix | test | infrastructure | release | decision
	Title string
	Slug  string // filename-safe slug derived from Title (see Slugify)
}

// shortUUID returns the 8-char display form used in backlog rows/filenames.
func (i NewItemInput) shortUUID() string {
	if len(i.UUID) > 8 {
		return i.UUID[:8]
	}
	return i.UUID
}

// PlanDirName returns the canonical active/ item directory name for this item:
// the stable 8-char UUID display form.
func (i NewItemInput) PlanDirName() string {
	return i.shortUUID()
}

// PlanFileName returns the plan file name inside the item directory.
func (i NewItemInput) PlanFileName() string {
	return "plan.md"
}

// PlanRelativePath returns the canonical repo-relative plan path for this item.
func (i NewItemInput) PlanRelativePath() string {
	return fmt.Sprintf("active/%s/%s", i.PlanDirName(), i.PlanFileName())
}

// BacklogRow renders the Active Work table row for this item, status
// "inbox" — matching the table header:
// | UUID | Legacy ID | Type | Title | Status | Owner | Plan File | Updated |
func (i NewItemInput) BacklogRow(date string) string {
	return fmt.Sprintf("| `%s` | — | %s | %s | inbox | unassigned | [plan](active/%s) | %s |",
		i.shortUUID(), i.Type, i.Title, i.PlanDirName()+"/"+i.PlanFileName(), date)
}

// PlanFileContent renders the intake investigation template shape
// (.hawp/kit/templates/work-intake.md), pre-filled with this item's known
// fields. Investigation/analysis sections are left as placeholders — filling
// them in is the next (required) step in the intake workflow, not something
// this scaffold command does on its own.
func (i NewItemInput) PlanFileContent(inputText, date string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "# %s\n\n", i.Title)
	sb.WriteString("**Backlog ID (Legacy):** — (UUID-native item)\n")
	fmt.Fprintf(&sb, "**UUID:** `%s`\n", i.UUID)
	fmt.Fprintf(&sb, "**Type:** %s\n", i.Type)
	fmt.Fprintf(&sb, "**Reported:** %s\n\n", date)
	sb.WriteString("---\n\n")
	sb.WriteString("## Input (verbatim)\n\n")
	fmt.Fprintf(&sb, "> %s\n\n", inputText)
	sb.WriteString("## Intake Summary\n\n")
	sb.WriteString("_Not yet investigated._\n\n")
	sb.WriteString("## Current Context\n\n")
	sb.WriteString("_Not yet investigated._\n\n")
	sb.WriteString("## Initial Analysis\n\n")
	sb.WriteString("**Directly verified:**\n\n- _pending_\n\n")
	sb.WriteString("**Inferred (not yet proven):**\n\n- _pending_\n\n")
	sb.WriteString("**Likely scope:**\n\n- _pending_\n\n")
	sb.WriteString("## Risk + Review Gate\n\n")
	sb.WriteString("**Risk:** _pending_ (low | medium | high)\n")
	sb.WriteString("**Gate:** _pending_ (auto-implement on low | review first on medium/high)\n\n")
	sb.WriteString("## Backlog + Plan Link\n\n")
	sb.WriteString("**Status now:** inbox\n")
	fmt.Fprintf(&sb, "**Plan file:** work/%s\n\n", i.PlanRelativePath())
	sb.WriteString("## Next Step\n\n")
	sb.WriteString("- [ ] Investigation recorded above (required before planning)\n")
	sb.WriteString("- [ ] Write or update the plan file\n")
	sb.WriteString("- [ ] Move backlog status accordingly\n")
	return sb.String()
}

// Slugify converts a title into a filename-safe slug: lowercase,
// alphanumerics kept, everything else collapsed to single hyphens, trimmed,
// capped at 60 chars. Empty input produces "item" rather than an empty slug.
func Slugify(title string) string {
	var b strings.Builder
	lastHyphen := true // treat start as "just wrote a hyphen" to avoid leading '-'
	for _, r := range strings.ToLower(title) {
		switch {
		case r >= 'a' && r <= 'z' || r >= '0' && r <= '9':
			b.WriteRune(r)
			lastHyphen = false
		default:
			if !lastHyphen {
				b.WriteRune('-')
				lastHyphen = true
			}
		}
	}
	slug := strings.TrimRight(b.String(), "-")
	if len(slug) > 60 {
		slug = strings.TrimRight(slug[:60], "-")
	}
	if slug == "" {
		slug = "item"
	}
	return slug
}

// InsertActiveRow inserts a new row into BACKLOG.md's "## Active Work"
// table, appended just before the next "## " heading (i.e. at the end of
// the Active Work section). Returns an error if that section isn't found,
// so a malformed/missing BACKLOG.md fails loudly instead of writing a row
// nobody will ever see.
func InsertActiveRow(backlogContent, row string) (string, error) {
	lines := strings.Split(backlogContent, "\n")
	activeIdx := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == "## Active Work" {
			activeIdx = i
			break
		}
	}
	if activeIdx == -1 {
		return "", fmt.Errorf(`BACKLOG.md has no "## Active Work" section`)
	}

	insertAt := len(lines)
	for i := activeIdx + 1; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "## ") {
			insertAt = i
			break
		}
	}

	// Walk back over trailing blank lines so the new row lands directly
	// after the last table row, not after a blank-line gap.
	for insertAt > activeIdx+1 && strings.TrimSpace(lines[insertAt-1]) == "" {
		insertAt--
	}

	out := make([]string, 0, len(lines)+1)
	out = append(out, lines[:insertAt]...)
	out = append(out, row)
	out = append(out, lines[insertAt:]...)
	return strings.Join(out, "\n"), nil
}
