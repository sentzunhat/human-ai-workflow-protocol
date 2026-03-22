package context

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	domainwork "github.com/sentzunhat/hawp/librarian/src/internal/domain/work"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/markdown"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
)

// workRoleFolders are the top-level .hawp/work/ subfolders enriched here.
// active/closed/parked resolve metadata from BACKLOG.md; the rest
// (decisions, evidence, notes, status) are archival/supporting folders
// with no backlog row — they still get a Role, just no Type/Status.
var workRoleFolders = []string{"active", "closed", "parked", "decisions", "evidence", "notes", "status"}

var closedDatePathRe = regexp.MustCompile(`/closed/(\d{4})/(\d{2})/(\d{2})/`)

func closedDateFromPath(path string) string {
	if m := closedDatePathRe.FindStringSubmatch(filepath.ToSlash(path)); m != nil {
		return m[1] + "-" + m[2] + "-" + m[3]
	}
	return ""
}

// rowByID looks up a backlog row matching id, accepting short-UUID
// prefix matches (row `0e1c4afa` <-> file `0e1c4afa-....md`).
func rowByID(rows []domainwork.BacklogRow, id string) (domainwork.BacklogRow, bool) {
	if id == "" {
		return domainwork.BacklogRow{}, false
	}
	for _, row := range rows {
		if domainwork.IDsMatch(row.ID, id) {
			return row, true
		}
	}
	return domainwork.BacklogRow{}, false
}

func resolveID(filename string) string {
	name := strings.TrimSuffix(filename, ".md")
	if id := domainwork.ExtractIDFromFilename(name); id != "" {
		return id
	}
	return domainwork.ExtractShortUUID(name)
}

// EnrichWork walks workRoot (a .hawp/work directory) and returns every
// markdown document tagged with its folder role and, for active/closed/
// parked records, metadata resolved from BACKLOG.md.
func EnrichWork(repoRoot, workRoot string) ([]Document, error) {
	backlog, err := domainwork.ParseBacklog(filepath.Join(workRoot, "BACKLOG.md"))
	if err != nil {
		return nil, err
	}
	rowsByRole := map[string][]domainwork.BacklogRow{
		"active": backlog.Active, "closed": backlog.Closed, "parked": backlog.Parked,
	}

	var documents []Document
	for _, role := range workRoleFolders {
		dir := filepath.Join(workRoot, role)
		for _, file := range markdown.CollectFiles(dir, true) {
			documents = append(documents, buildWorkDocument(repoRoot, file, role, rowsByRole[role]))
		}
	}
	return documents, nil
}

func buildWorkDocument(repoRoot, file, role string, rows []domainwork.BacklogRow) Document {
	raw, _ := os.ReadFile(file)
	id := resolveID(filepath.Base(file))
	doc := Document{
		RelPath: repo.ToRepoRelative(repoRoot, file),
		Corpus:  CorpusWork,
		Role:    role,
		ID:      id,
		Content: string(raw),
	}
	if closedDate := closedDateFromPath(file); closedDate != "" {
		doc.ClosedDate = closedDate
	}
	if row, ok := rowByID(rows, id); ok {
		doc.Type = row.Type
		if role == "closed" {
			doc.Status = "closed"
			if doc.ClosedDate == "" {
				doc.ClosedDate = strings.TrimSpace(row.Status)
			}
		} else {
			doc.Status = row.Status
		}
	}
	doc.ContextPrefix = buildWorkContextPrefix(doc)
	return doc
}

func buildWorkContextPrefix(doc Document) string {
	prefix := "[work/" + doc.Role + "]"
	if doc.ID == "" {
		return prefix
	}
	detail := doc.ID
	if doc.Type != "" {
		switch {
		case doc.Role == "closed" && doc.ClosedDate != "":
			detail += " (" + doc.Type + ", closed " + doc.ClosedDate + ")"
		case doc.Status != "":
			detail += " (" + doc.Type + ", " + doc.Status + ")"
		default:
			detail += " (" + doc.Type + ")"
		}
	}
	return prefix + " " + detail
}
