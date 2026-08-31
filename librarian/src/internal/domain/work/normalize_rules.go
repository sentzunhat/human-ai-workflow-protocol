package work

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
)

// FixOperation is one detected fix (auto-fixable or blocked).
type FixOperation struct {
	OpID         string       `json:"opId"`
	Type         string       `json:"type"`
	ItemID       string       `json:"itemId"`
	FileToModify string       `json:"fileToModify"`
	LineRange    [2]int       `json:"lineRange"`
	Description  string       `json:"description"`
	Safety       string       `json:"safety"` // safe | blocked
	Confidence   float64      `json:"confidence"`
	RuleID       string       `json:"ruleId,omitempty"`
	Blocked      *BlockedInfo `json:"blocked,omitempty"`
}

// BlockedInfo details a blocked operation.
type BlockedInfo struct {
	ID         string   `json:"id"`
	Rule       string   `json:"rule"`
	ItemID     string   `json:"itemId"`
	Confidence float64  `json:"confidence"`
	Candidates []string `json:"candidates"`
	Reason     string   `json:"reason"`
	Recovery   string   `json:"recovery"`
}

var (
	isoDateRe        = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	canonicalIDRe    = regexp.MustCompile(`^(TASK|BUG)-\d+$`)
	numericRowIDRe   = regexp.MustCompile(`^\d+$`)
	verifySectionRe  = regexp.MustCompile(`(?ms)^##\s+Verification\b(.*?)(?:^##\s+|\z)`)
	checklistLineRe  = regexp.MustCompile(`- \[[x ]\]`)
	unprovenMarkerRe = regexp.MustCompile(`(?i)(NOT YET VERIFIED|\bunproven\b)`)
	staleTemplateRes = []*regexp.Regexp{
		regexp.MustCompile(`(?i)\]\(core/distribution/generated/`),
		regexp.MustCompile(`(?i)\]\(core/distribution/sources/`),
		regexp.MustCompile("(?i)`core/distribution/"),
	}
)

func isISODateOrEmpty(value string) bool {
	if value == "" || strings.EqualFold(value, "n/a") {
		return true
	}
	return isoDateRe.MatchString(value)
}

func isCanonicalID(value string) bool {
	return canonicalIDRe.MatchString(value) || numericRowIDRe.MatchString(value) || fullUUIDRe.MatchString(value) || ExtractShortUUID(value) != ""
}

func inferTypeFromID(id string) string {
	switch {
	case strings.HasPrefix(id, "TASK-"):
		return "task"
	case strings.HasPrefix(id, "BUG-"):
		return "bug"
	default:
		return ""
	}
}

// resolveWithinWorkRoot resolves a backlog-relative plan path, rejecting
// paths that escape the work root.
func resolveWithinWorkRoot(workRoot, relativePath string) string {
	root, err := filepath.Abs(workRoot)
	if err != nil {
		return ""
	}
	resolved := filepath.Clean(filepath.Join(root, relativePath))
	if !strings.HasPrefix(resolved, root+string(filepath.Separator)) {
		return ""
	}
	return resolved
}

func verificationSectionBody(content string) string {
	if m := verifySectionRe.FindStringSubmatch(content); m != nil {
		return m[1]
	}
	return ""
}

func isEvidenceFollowUpLine(line string) bool {
	claim := strings.TrimSpace(claimPrefixRe.ReplaceAllString(line, ""))
	return strings.HasPrefix(claim, "Research evidence for:") ||
		strings.HasPrefix(claim, "Update the original verification checklist line with Evidence:")
}

// hasUnprovenChecklistMarker reports checklist lines that carry an unproven
// marker without Evidence:. (The TS original never skipped Evidence-backed
// lines, but its section regex used `\Z` — a literal, case-insensitive "z"
// in JavaScript — which truncated the section and made the rule dead code;
// this port aligns B4 with the verification-clarity semantics instead.)
func hasUnprovenChecklistMarker(content string) bool {
	for _, line := range strings.Split(verificationSectionBody(content), "\n") {
		if !checklistLineRe.MatchString(line) || isEvidenceFollowUpLine(line) {
			continue
		}
		if strings.Contains(line, "Evidence:") {
			continue
		}
		if unprovenMarkerRe.MatchString(line) {
			return true
		}
	}
	return false
}

// AmbiguousVerificationClaims returns checklist claims in the Verification
// section with no Evidence: marker and no explicit unproven wording.
func AmbiguousVerificationClaims(content string) []string {
	var claims []string
	for _, line := range strings.Split(verificationSectionBody(content), "\n") {
		if !checklistLineRe.MatchString(line) || isEvidenceFollowUpLine(line) {
			continue
		}
		if strings.Contains(line, "Evidence:") || unprovenMarkerRe.MatchString(line) {
			continue
		}
		claim := strings.TrimSpace(claimPrefixRe.ReplaceAllString(line, ""))
		if claim != "" {
			claims = append(claims, claim)
		}
	}
	return claims
}

func hasStaleTemplateReference(content string) bool {
	for _, re := range staleTemplateRes {
		if re.MatchString(content) {
			return true
		}
	}
	return false
}

func hasHeadingNamed(content, heading string) bool {
	re := regexp.MustCompile(`(?im)^#{1,6}\s+` + regexp.QuoteMeta(heading) + `\b`)
	return re.MatchString(content)
}

func isLegacyClosedPath(path string) bool {
	m := pathDateRe.FindStringSubmatch(strings.ReplaceAll(path, "\\", "/"))
	if m == nil {
		return false
	}
	return m[1]+"-"+m[2]+"-"+m[3] < repo.LegacyClosedCutoff
}

type ruleEvaluator struct {
	repoRoot    string
	workRoot    string
	backlogPath string
	scan        *PlanScan
	planByPath  map[string]*PlanFileRecord
	opCounter   int
	ops         []FixOperation
}

func (e *ruleEvaluator) nextOpID() string {
	id := fmt.Sprintf("OP-%03d", e.opCounter)
	e.opCounter++
	return id
}

func (e *ruleEvaluator) autoFix(ruleID, opType, itemID, file string, line int, description string) {
	e.ops = append(e.ops, FixOperation{
		OpID: e.nextOpID(), Type: opType, ItemID: itemID, FileToModify: file,
		LineRange: [2]int{line, line}, Description: description,
		Safety: "safe", Confidence: 0.9, RuleID: ruleID,
	})
}

func (e *ruleEvaluator) blockedFix(ruleID, itemID, file string, line int, reason string, candidates []string, recovery string) {
	opID := e.nextOpID()
	e.ops = append(e.ops, FixOperation{
		OpID: opID, Type: "add-field", ItemID: itemID, FileToModify: file,
		LineRange: [2]int{line, line}, Description: reason,
		Safety: "blocked", Confidence: 0, RuleID: ruleID,
		Blocked: &BlockedInfo{
			ID:   "BLOCKED-" + strings.TrimPrefix(opID, "OP-"),
			Rule: ruleID, ItemID: itemID, Confidence: 0,
			Candidates: candidates, Reason: reason, Recovery: recovery,
		},
	})
}

// EvaluateRules runs the detection rule set (A1-A8, B1-B5, B7) over the
// backlog rows and plan scan, producing ordered fix operations.
func EvaluateRules(repoRoot, workRoot, backlogPath string, backlog *NormalizeBacklog, scan *PlanScan) []FixOperation {
	e := &ruleEvaluator{
		repoRoot: repoRoot, workRoot: workRoot, backlogPath: backlogPath,
		scan: scan, planByPath: map[string]*PlanFileRecord{}, opCounter: 1,
	}
	for i := range scan.Files {
		e.planByPath[scan.Files[i].Path] = &scan.Files[i]
	}

	// Structural rules.
	if !backlog.SectionPresence[SectionActive] || !backlog.SectionPresence[SectionBlocked] || !backlog.SectionPresence[SectionClosed] {
		e.autoFix("A8", "add-section-header", "BACKLOG-STRUCTURE", backlogPath, 1,
			"A8: add missing required backlog section headers")
	}
	var missingDirs []string
	for _, name := range []string{"active", "parked", "closed"} {
		if !scan.DirectoryPresence[name] {
			missingDirs = append(missingDirs, name)
		}
	}
	if len(missingDirs) > 0 {
		e.blockedFix("B5", "WORKSPACE", backlogPath, 1,
			"Non-standard work folder structure detected", missingDirs,
			"Restore missing .hawp/work directories (active, parked, closed)")
	}

	for _, row := range backlog.Rows {
		e.evaluateRow(row)
	}
	return e.ops
}

func (e *ruleEvaluator) evaluateRow(row NormalizeRow) {
	candidates := e.scan.ByID[row.ID]

	if row.Type == "" {
		if inferred := inferTypeFromID(row.ID); inferred != "" {
			e.autoFix("A1", "add-field", row.ID, e.backlogPath, row.LineNumber,
				fmt.Sprintf("A1: add missing type field inferred from ID prefix (%s)", inferred))
		} else {
			e.blockedFix("B1", row.ID, e.backlogPath, row.LineNumber,
				"Cannot infer missing type from ID",
				[]string{"task", "bug", "improvement"},
				"Set an explicit type in BACKLOG.md")
		}
	}

	if !isISODateOrEmpty(row.Updated) {
		e.autoFix("A2", "normalize-date", row.ID, e.backlogPath, row.LineNumber,
			fmt.Sprintf("A2: normalize date '%s' to YYYY-MM-DD", row.Updated))
	}

	if !isCanonicalID(row.ID) {
		e.autoFix("A3", "fix-malformed-id", row.ID, e.backlogPath, row.LineNumber,
			"A3: fix malformed backlog ID format")
	}

	var linkedPath string
	if row.PlanPath != "" {
		linkedPath = resolveWithinWorkRoot(e.workRoot, row.PlanPath)
		if linkedPath == "" {
			e.blockedFix("B2", row.ID, e.backlogPath, row.LineNumber,
				"Referenced plan path escapes the work root: "+row.PlanPath, nil,
				"Use a plan link relative to .hawp/work without '..' segments")
		} else if !repo.Exists(linkedPath) {
			e.blockedFix("B2", row.ID, e.backlogPath, row.LineNumber,
				"Referenced plan path is missing: "+row.PlanPath, nil,
				"Update BACKLOG.md link or restore the plan file")
		}
	} else if len(candidates) == 0 {
		e.blockedFix("B2", row.ID, e.backlogPath, row.LineNumber,
			"No plan link and no matching plan file were found", nil,
			"Create a plan file or add a valid plan link")
	}

	if len(candidates) > 1 {
		rel := make([]string, len(candidates))
		for i, c := range candidates {
			rel[i] = repo.ToRepoRelative(e.repoRoot, c)
		}
		e.blockedFix("B3", row.ID, e.backlogPath, row.LineNumber,
			"Multiple plan files matched the same backlog ID", rel,
			"Choose a canonical plan file and archive duplicate candidates")
	}

	var linkedPlan *PlanFileRecord
	if linkedPath != "" {
		linkedPlan = e.planByPath[linkedPath]
	}

	if row.Section == SectionClosed && linkedPlan != nil {
		content := linkedPlan.Content
		isLegacy := isLegacyClosedPath(linkedPlan.Path)
		planRel := repo.ToRepoRelative(e.repoRoot, linkedPlan.Path)

		if !isLegacy && !hasHeadingNamed(content, "Outcome") {
			e.autoFix("A4", "add-section-header", row.ID, planRel, 1,
				"A4: add missing '## Outcome' section heading")
		}
		if !isLegacy && (!hasHeadingNamed(content, "Verification") || !hasHeadingNamed(content, "Close Checklist")) {
			e.autoFix("A5", "add-scaffolding", row.ID, planRel, 1,
				"A5: add missing verification/checklist scaffolding")
		}
		if hasUnprovenChecklistMarker(content) {
			e.blockedFix("B4", row.ID, planRel, 1,
				"Closed plan still contains unproven verification markers", nil,
				"Complete verification evidence or mark unresolved risk explicitly before closing")
		}
		if claims := AmbiguousVerificationClaims(content); len(claims) > 0 {
			shown := claims
			if len(shown) > 5 {
				shown = shown[:5]
			}
			e.blockedFix("B7", row.ID, planRel, 1,
				fmt.Sprintf("Closed plan has %d ambiguous verification checklist claim(s)", len(claims)),
				shown,
				"Add `Evidence:` to each checklist claim or mark it explicitly unproven before relying on the record")
		}
		if !isLegacy && hasStaleTemplateReference(content) {
			e.autoFix("A7", "update-template-reference", row.ID, planRel, 1,
				"A7: update outdated template/path references")
		}
	}

	if row.Section == SectionActive {
		status := strings.ToLower(row.Status)
		if status == "done" || status == "wont-fix" {
			e.autoFix("A6", "migrate-row", row.ID, e.backlogPath, row.LineNumber,
				"A6: migrate completed active row to recently closed")
		}
	}
}
