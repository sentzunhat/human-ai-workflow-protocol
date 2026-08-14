package context

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/context/source"
	domainwork "github.com/sentzunhat/hawp/librarian/go/internal/domain/work"
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

// EnrichWork converts acquired work files and backlog metadata into context
// documents. The source adapter owns traversal, reads, and backlog loading.
func EnrichWork(corpus source.WorkCorpus) []Document {
	if corpus.Backlog == nil {
		corpus.Backlog = &domainwork.Backlog{}
	}
	rowsByRole := map[string][]domainwork.BacklogRow{
		"active": corpus.Backlog.Active, "closed": corpus.Backlog.Closed, "parked": corpus.Backlog.Parked,
	}

	allowedRoles := map[string]bool{}
	for _, role := range workRoleFolders {
		allowedRoles[role] = true
	}

	documents := make([]Document, 0, len(corpus.Files))
	for _, file := range corpus.Files {
		role := workRole(filepath.ToSlash(file.RelPath))
		if !allowedRoles[role] {
			continue
		}
		documents = append(documents, buildWorkDocument(file, role, rowsByRole[role]))
	}
	return documents
}

func workRole(relPath string) string {
	parts := strings.Split(strings.TrimPrefix(relPath, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return ""
	}
	return parts[0]
}

func buildWorkDocument(file source.File, role string, rows []domainwork.BacklogRow) Document {
	id := resolveID(filepath.Base(file.RelPath))
	doc := Document{
		RelPath: file.RepoPath,
		Corpus:  CorpusWork,
		Role:    role,
		ID:      id,
		Content: file.Content,
	}
	if closedDate := closedDateFromPath("/" + filepath.ToSlash(file.RelPath) + "/"); closedDate != "" {
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
