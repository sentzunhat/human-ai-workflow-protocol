// Package doc is the application use case for creating secondary work
// documents (status reports, evidence files, decision records, notes)
// at the canonical UUID-subfolder path:
//
//	{docType}/YYYY/MM/DD/{uuid}/{filename}.md
package doc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sentzunhat/hawp/librarian/src/internal/application/uuidgen"
)

// ValidTypes lists the accepted secondary document types.
var ValidTypes = []string{"status", "evidence", "decision", "note"}

// Result reports what CreateWorkDoc wrote.
type Result struct {
	UUID     string
	DocType  string
	Title    string
	FilePath string // absolute path to the new file
}

// CreateWorkDoc resolves the canonical date-stamped path under workDir and
// writes a minimal front-matter template.
//
// workDir is the absolute path to the .hawp/work directory. docType must be
// one of ValidTypes. workItemID ties the document to an existing work item
// (use its 8-char short ID or full UUID); pass "" to auto-generate a fresh ID.
// Using the work item's ID groups all artifacts for that item under the same
// folder name, making them discoverable by searching for the ID.
func CreateWorkDoc(docType, title, workDir, workItemID string) (*Result, error) {
	if !isValidType(docType) {
		return nil, fmt.Errorf("unknown doc type %q (want %s)", docType, strings.Join(ValidTypes, "|"))
	}
	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("title is required")
	}

	// Resolve the folder ID: use the provided work item ID (normalised to its
	// short form when a full UUID is given) or generate a fresh short UUID.
	folderID := strings.TrimSpace(workItemID)
	fullUUID := folderID
	if folderID == "" {
		id, err := uuidgen.New()
		if err != nil {
			return nil, fmt.Errorf("generate uuid: %w", err)
		}
		fullUUID = id
		folderID = uuidgen.Short(id)
	} else if len(folderID) > 8 && strings.Contains(folderID, "-") {
		// Caller passed a full UUID — use the short form for the folder to
		// stay consistent with active/ and closed/ naming.
		folderID = uuidgen.Short(folderID)
	}

	date := time.Now().Format("2006-01-02")
	dirName := folderName(docType)
	fileName := docType + ".md"

	docDir := filepath.Join(workDir, dirName, date[:4], date[5:7], date[8:10], folderID)
	if err := os.MkdirAll(docDir, 0o755); err != nil {
		return nil, fmt.Errorf("create doc directory: %w", err)
	}

	filePath := filepath.Join(docDir, fileName)
	if err := os.WriteFile(filePath, []byte(template(fullUUID, docType, title, date)), 0o644); err != nil {
		return nil, fmt.Errorf("write doc file: %w", err)
	}

	return &Result{
		UUID:     folderID,
		DocType:  docType,
		Title:    title,
		FilePath: filePath,
	}, nil
}

// folderName maps a doc type to its directory name under .hawp/work/.
func folderName(docType string) string {
	if docType == "decision" {
		return "decisions"
	}
	if docType == "note" {
		return "notes"
	}
	return docType // "status", "evidence"
}

func isValidType(t string) bool {
	for _, v := range ValidTypes {
		if t == v {
			return true
		}
	}
	return false
}

func template(id, docType, title, date string) string {
	return fmt.Sprintf(`---
uuid: %s
title: %s
type: %s
date: %s
---

`, id, title, docType, date)
}
