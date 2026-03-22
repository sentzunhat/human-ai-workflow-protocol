package work

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sentzunhat/hawp/librarian/src/internal/application/uuidgen"
	domainwork "github.com/sentzunhat/hawp/librarian/src/internal/domain/work"
)

// NewItemResult reports what NewItem created.
type NewItemResult struct {
	UUID         string
	Type         string
	Title        string
	PlanFilePath string // absolute path
	BacklogPath  string // absolute path
}

// validTypes mirrors the Type column values used across BACKLOG.md.
var validTypes = map[string]bool{
	"task": true, "bug": true, "improvement": true, "feature": true,
	"fix": true, "test": true, "infrastructure": true, "release": true,
	"decision": true,
}

// NewItem scaffolds a work item — the mechanical half of HAWP's intake
// step (see .hawp/kit/usage/intake-workflow.md): it generates a UUID,
// writes an investigation plan file shaped like
// .hawp/kit/templates/work-intake.md, and inserts a status="inbox" row
// into BACKLOG.md's Active Work table.
//
// It deliberately does NOT perform the investigation or write the plan
// itself — HAWP is "a shaping protocol, not a runtime" (AGENTS.md); the
// investigation and plan content require actual reasoning about the
// specific request, which is a human or AI agent's job, not a template
// fill-in. This command exists so that job doesn't also include hand-typing
// a backlog table row and a file skeleton.
func NewItem(workDir, itemType, title, inputText string) (*NewItemResult, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if itemType == "" {
		itemType = "task"
	}
	if !validTypes[itemType] {
		return nil, fmt.Errorf("type %q is not a recognized BACKLOG.md Type value (task, bug, improvement, feature, fix, test, infrastructure, release, decision)", itemType)
	}
	if inputText == "" {
		inputText = title
	}

	id, err := uuidgen.New()
	if err != nil {
		return nil, fmt.Errorf("generate uuid: %w", err)
	}

	item := domainwork.NewItemInput{
		UUID:  id,
		Type:  itemType,
		Title: title,
		Slug:  domainwork.Slugify(title),
	}

	date := time.Now().Format("2006-01-02")

	backlogPath := filepath.Join(workDir, "BACKLOG.md")
	backlogBytes, err := os.ReadFile(backlogPath)
	if err != nil {
		return nil, fmt.Errorf("read BACKLOG.md: %w", err)
	}

	updatedBacklog, err := domainwork.InsertActiveRow(string(backlogBytes), item.BacklogRow(date))
	if err != nil {
		return nil, fmt.Errorf("insert backlog row: %w", err)
	}

	planPath := filepath.Join(workDir, "active", item.PlanFileName())
	if _, err := os.Stat(planPath); err == nil {
		return nil, fmt.Errorf("plan file already exists: %s (slug collision — try a more distinct title)", planPath)
	}

	if err := os.MkdirAll(filepath.Join(workDir, "active"), 0755); err != nil {
		return nil, fmt.Errorf("create active/ directory: %w", err)
	}
	if err := os.WriteFile(planPath, []byte(item.PlanFileContent(inputText, date)), 0644); err != nil {
		return nil, fmt.Errorf("write plan file: %w", err)
	}
	// Write the plan file before the backlog row: if the process dies
	// between the two writes, an orphaned plan file with no backlog row is
	// safer to notice and clean up than a backlog row pointing at a file
	// that doesn't exist.
	if err := os.WriteFile(backlogPath, []byte(updatedBacklog), 0644); err != nil {
		return nil, fmt.Errorf("write BACKLOG.md: %w", err)
	}

	return &NewItemResult{
		UUID:         id,
		Type:         itemType,
		Title:        title,
		PlanFilePath: planPath,
		BacklogPath:  backlogPath,
	}, nil
}
