package work

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// FixPlan is the complete detected upgrade plan.
type FixPlan struct {
	PlanID           string         `json:"planId"`
	PlanHash         string         `json:"planHash"`
	ScannedAt        string         `json:"scannedAt"`
	BacklogPath      string         `json:"backlogPath"`
	FilesScanned     int            `json:"filesScanned"`
	ItemsAnalyzed    int            `json:"itemsAnalyzed"`
	Operations       []FixOperation `json:"operations"`
	AutoFixCount     int            `json:"autoFixCount"`
	BlockedCount     int            `json:"blockedCount"`
	EstimatedChanges int            `json:"estimatedChanges"`
}

// SyncPlanStep is one concrete drift resolution for duplicate plan files.
type SyncPlanStep struct {
	ItemID         string   `json:"itemId"`
	CanonicalPlan  string   `json:"canonicalPlan"`
	DuplicatePlans []string `json:"duplicatePlans"`
	ApplySteps     []string `json:"applySteps"`
}

// ResearchItem is a verification claim that needs evidence follow-up.
type ResearchItem struct {
	ItemID            string `json:"itemId"`
	Claim             string `json:"claim"`
	FilePath          string `json:"filePath"`
	LineNumber        int    `json:"lineNumber"`
	RecommendedAction string `json:"recommendedAction"`
}

// DetectionReport is the dry-run scan result.
type DetectionReport struct {
	ReportID    string  `json:"reportId"`
	GeneratedAt string  `json:"generatedAt"`
	Plan        FixPlan `json:"plan"`
	Assessment  string  `json:"assessment"` // clean | drift-detected
	Summary     struct {
		TotalIssues   int      `json:"totalIssues"`
		AutoFixable   int      `json:"autoFixableCount"`
		Blocked       int      `json:"blockedCount"`
		FilesAffected []string `json:"filesAffected"`
		ItemsAffected []string `json:"itemsAffected"`
	} `json:"summary"`
	Recommendation struct {
		Action  string   `json:"action"`
		Reason  string   `json:"reason"`
		Details []string `json:"details,omitempty"`
	} `json:"recommendation"`
	SyncPlan      []SyncPlanStep `json:"syncPlan"`
	ResearchQueue []ResearchItem `json:"researchQueue"`
}

func stableHash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func uniqueStrings(values []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, v := range values {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

// BuildDetectionReport assembles the plan, summary, recommendation, and sync
// plan from the evaluated operations.
func BuildDetectionReport(scannedAt, backlogPath string, filesScanned, itemsAnalyzed int, operations []FixOperation) DetectionReport {
	autoFix, blocked := 0, 0
	var files, items []string
	for _, op := range operations {
		if op.Safety == "safe" {
			autoFix++
		}
		if op.Safety == "blocked" {
			blocked++
		}
		files = append(files, op.FileToModify)
		items = append(items, op.ItemID)
	}

	seed, _ := json.Marshal(struct {
		Operations []FixOperation `json:"operations"`
	}{operations})
	planHash := stableHash(string(seed))

	plan := FixPlan{
		PlanID: "PLAN-" + planHash[:8], PlanHash: planHash, ScannedAt: scannedAt,
		BacklogPath: backlogPath, FilesScanned: filesScanned, ItemsAnalyzed: itemsAnalyzed,
		Operations: operations, AutoFixCount: autoFix, BlockedCount: blocked,
		EstimatedChanges: autoFix,
	}

	report := DetectionReport{
		ReportID:    "REPORT-" + stableHash(planHash + ":" + scannedAt)[:8],
		GeneratedAt: scannedAt,
		Plan:        plan,
		SyncPlan:    buildSyncPlan(operations),
	}
	report.Summary.TotalIssues = autoFix + blocked
	report.Summary.AutoFixable = autoFix
	report.Summary.Blocked = blocked
	report.Summary.FilesAffected = uniqueStrings(files)
	report.Summary.ItemsAffected = uniqueStrings(items)

	switch {
	case autoFix+blocked == 0:
		report.Assessment = "clean"
		report.Recommendation.Action = "no-action"
		report.Recommendation.Reason = "Backlog is consistent with current templates and standards"
	case blocked > 0:
		report.Assessment = "drift-detected"
		report.Recommendation.Action = "manual-review"
		report.Recommendation.Reason = fmt.Sprintf("%d issue(s) require manual review before fixes can be applied", blocked)
		report.Recommendation.Details = []string{
			"Run: ./.hawp/bin/hawp backlog upgrade --dry-run --format json to see blocked items",
		}
		if len(report.SyncPlan) > 0 {
			report.Recommendation.Details = append(report.Recommendation.Details,
				fmt.Sprintf("Concrete sync/apply steps generated for %d duplicate-item drift case(s).", len(report.SyncPlan)))
		}
	default:
		report.Assessment = "drift-detected"
		report.Recommendation.Action = "apply-fixes"
		report.Recommendation.Reason = fmt.Sprintf("%d mechanical fix(es) ready to apply", autoFix)
		report.Recommendation.Details = []string{
			"Review changes: ./.hawp/bin/hawp backlog upgrade --dry-run --format text",
			"Apply fixes: ./.hawp/bin/hawp backlog upgrade --apply",
		}
	}
	return report
}

func buildSyncPlan(operations []FixOperation) []SyncPlanStep {
	var steps []SyncPlanStep
	for _, op := range operations {
		if op.RuleID != "B3" || op.Blocked == nil {
			continue
		}
		candidates := uniqueStrings(op.Blocked.Candidates)
		if len(candidates) < 2 {
			continue
		}
		canonical := ""
		for _, c := range candidates {
			if strings.Contains(c, "/closed/") {
				canonical = c
				break
			}
		}
		if canonical == "" {
			steps = append(steps, SyncPlanStep{
				ItemID: op.ItemID, CanonicalPlan: "manual-selection-required",
				DuplicatePlans: candidates,
				ApplySteps: []string{
					fmt.Sprintf("Review candidates for %s and choose one canonical plan file.", op.ItemID),
					"Move remaining duplicates to the closed archive or remove stale copies.",
					"Re-run: ./.hawp/bin/hawp backlog upgrade --dry-run --validate",
				},
			})
			continue
		}
		var duplicates []string
		applySteps := []string{"Keep canonical closed plan: " + canonical}
		for _, c := range candidates {
			if c == canonical {
				continue
			}
			duplicates = append(duplicates, c)
			switch {
			case strings.Contains(c, "/active/"):
				applySteps = append(applySteps, "Remove stale active copy: "+c)
			case strings.Contains(c, "/parked/"):
				applySteps = append(applySteps, "Remove stale parked copy: "+c)
			default:
				applySteps = append(applySteps, "Archive or remove duplicate copy: "+c)
			}
		}
		applySteps = append(applySteps, "Re-run: ./.hawp/bin/hawp backlog upgrade --dry-run --validate")
		steps = append(steps, SyncPlanStep{
			ItemID: op.ItemID, CanonicalPlan: canonical,
			DuplicatePlans: duplicates, ApplySteps: applySteps,
		})
	}
	return steps
}

// RenderTextReport renders the dry-run report in the TS text format.
func RenderTextReport(report DetectionReport) string {
	var b strings.Builder
	w := func(line string) { b.WriteString(line + "\n") }

	w("HAWP Backlog Upgrade Dry-Run Report")
	w("=================================")
	w("Report ID: " + report.ReportID)
	w("Generated: " + report.GeneratedAt)
	w("Assessment: " + report.Assessment)
	w("")
	w("Summary")
	w("-------")
	w(fmt.Sprintf("Total issues: %d", report.Summary.TotalIssues))
	w(fmt.Sprintf("Auto-fixable: %d", report.Summary.AutoFixable))
	w(fmt.Sprintf("Blocked: %d", report.Summary.Blocked))
	w("Files affected: " + orNone(strings.Join(report.Summary.FilesAffected, ", ")))
	w("Items affected: " + orNone(strings.Join(report.Summary.ItemsAffected, ", ")))
	w("")
	w("Recommendation")
	w("--------------")
	w(report.Recommendation.Action + ": " + report.Recommendation.Reason)
	for _, detail := range report.Recommendation.Details {
		w("- " + detail)
	}
	w("")
	w("Drift Sync/Apply Plan")
	w("---------------------")
	if len(report.SyncPlan) == 0 {
		w("No duplicate working-file drift actions required.")
	} else {
		for _, step := range report.SyncPlan {
			w(step.ItemID + ": canonical=" + step.CanonicalPlan)
			if len(step.DuplicatePlans) > 0 {
				w("  duplicates: " + strings.Join(step.DuplicatePlans, ", "))
			}
			for _, s := range step.ApplySteps {
				w("  - " + s)
			}
		}
	}
	w("")
	w("Verification Research Queue")
	w("---------------------------")
	if len(report.ResearchQueue) == 0 {
		w("No verification evidence follow-up items detected.")
	} else {
		for _, item := range report.ResearchQueue {
			w(fmt.Sprintf("%s:%d %s [source: %s]", item.ItemID, item.LineNumber, item.Claim, item.FilePath))
			w("  - " + item.RecommendedAction)
		}
	}
	w("")
	w("Operations")
	w("----------")
	if len(report.Plan.Operations) == 0 {
		b.WriteString("No modifications needed.")
		return b.String()
	}
	for _, op := range report.Plan.Operations {
		ruleLabel := ""
		if op.RuleID != "" {
			ruleLabel = "[" + op.RuleID + "] "
		}
		w(fmt.Sprintf("%s %s[%s] %s %s:%d - %s",
			op.OpID, ruleLabel, op.Safety, op.ItemID, op.FileToModify, op.LineRange[0], op.Description))
		if op.Blocked != nil {
			w(fmt.Sprintf("  blocked (%s, confidence=%g): %s", op.Blocked.Rule, op.Blocked.Confidence, op.Blocked.Reason))
			if len(op.Blocked.Candidates) > 0 {
				w("  candidates: " + strings.Join(op.Blocked.Candidates, ", "))
			}
		}
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func orNone(value string) string {
	if value == "" {
		return "(none)"
	}
	return value
}

// RenderJSONReport renders the report as indented JSON.
func RenderJSONReport(report DetectionReport) (string, error) {
	return RenderJSONValue(report)
}

// RenderJSONValue renders any value as indented JSON with a trailing newline.
func RenderJSONValue(value any) (string, error) {
	out, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out) + "\n", nil
}
