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
	"regexp"
	"strings"
	"time"

	"github.com/sentzunhat/hawp/librarian/src/internal/application/uuidgen"
)

// validFolderID matches 8-char short UUIDs or full 36-char UUIDs.
var validFolderID = regexp.MustCompile(`^[0-9a-f]{8}(-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})?$`)

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
	} else {
		if !validFolderID.MatchString(folderID) {
			return nil, fmt.Errorf("invalid work item ID %q: must be an 8-char short UUID or full UUID", folderID)
		}
		if len(folderID) > 8 && strings.Contains(folderID, "-") {
			// Caller passed a full UUID — use the short form for the folder to
			// stay consistent with active/ and closed/ naming.
			folderID = uuidgen.Short(folderID)
		}
	}

	date := time.Now().Format("2006-01-02")
	dirName := folderName(docType)
	fileName := docType + ".md"

	docDir := filepath.Join(workDir, dirName, date[:4], date[5:7], date[8:10], folderID)
	if err := rejectSymlinkAncestors(workDir, docDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(docDir, 0o755); err != nil {
		return nil, fmt.Errorf("create doc directory: %w", err)
	}

	filePath := filepath.Join(docDir, fileName)
	if fi, err := os.Lstat(filePath); err == nil {
		// Path already exists — return it only when it is a regular file.
		// A directory or symlink here would make the return value unusable
		// or redirect writes outside the work tree.
		if !fi.Mode().IsRegular() {
			return nil, fmt.Errorf("doc path %s exists but is not a regular file", filePath)
		}
		return &Result{
			UUID:     folderID,
			DocType:  docType,
			Title:    title,
			FilePath: filePath,
		}, nil
	}
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

// rejectSymlinkAncestors walks every directory component of target that lies
// under root and returns an error if any component is a symlink. This prevents
// a symlink planted inside workDir from redirecting MkdirAll (and subsequent
// writes) outside the .hawp/work tree.
func rejectSymlinkAncestors(root, target string) error {
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return fmt.Errorf("resolve doc path: %w", err)
	}
	parts := strings.Split(rel, string(filepath.Separator))
	current := root
	for _, part := range parts {
		current = filepath.Join(current, part)
		fi, err := os.Lstat(current)
		if err != nil {
			if os.IsNotExist(err) {
				break // not yet created — MkdirAll will create it
			}
			return fmt.Errorf("stat doc path component %s: %w", current, err)
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("doc path component %s is a symlink; refusing to follow", current)
		}
	}
	return nil
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
