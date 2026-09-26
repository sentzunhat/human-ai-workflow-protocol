package ingest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	appindex "github.com/sentzunhat/hawp/librarian/src/internal/application/index"
	"github.com/sentzunhat/hawp/librarian/src/internal/domain/work/identity"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
)

// buildCorpusFromRepo walks each configured path and builds an enriched corpus.
// Paths are relative to repoRoot. ".hawp/kit" and ".hawp/work" get their
// enriched walkers; all other paths are walked generically as "custom" corpus.
func buildCorpusFromRepo(repoRoot string, paths []string) (*appindex.EnrichedCorpus, error) {
	corpus := &appindex.EnrichedCorpus{}

	for _, p := range paths {
		abs, err := resolveRepoPath(repoRoot, p)
		if err != nil {
			return nil, err
		}
		switch filepath.ToSlash(p) {
		case ".hawp/kit":
			if err := walkKitFiles(abs, corpus); err != nil {
				return nil, err
			}
		case ".hawp/work":
			if err := walkWorkFiles(abs, corpus); err != nil {
				return nil, err
			}
		default:
			if err := walkCustomPath(abs, p, corpus); err != nil {
				return nil, err
			}
		}
	}

	return corpus, nil
}

func resolveRepoPath(repoRoot, configuredPath string) (string, error) {
	path := filepath.FromSlash(configuredPath)
	if path == "" || filepath.IsAbs(path) {
		return "", fmt.Errorf("configured index path %q must be a non-empty relative path", configuredPath)
	}

	abs := filepath.Join(repoRoot, path)
	rel, err := filepath.Rel(repoRoot, abs)
	if err != nil {
		return "", fmt.Errorf("resolve configured index path %q: %w", configuredPath, err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("configured index path %q escapes repository root", configuredPath)
	}
	if err := filesystem.RejectSymlinkAncestors(repoRoot, abs); err != nil {
		return "", fmt.Errorf("configured index path %q is not safely contained in repository root: %w", configuredPath, err)
	}
	return abs, nil
}

func walkKitFiles(kitPath string, corpus *appindex.EnrichedCorpus) error {
	return filepath.Walk(kitPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !info.Mode().IsRegular() || filepath.Ext(path) != ".md" {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read kit document %s: %w", path, err)
		}
		rel, err := filepath.Rel(kitPath, path)
		if err != nil {
			return err
		}

		corpus.Documents = append(corpus.Documents, appindex.EnrichedDocument{
			Path:       filepath.ToSlash(filepath.Join(".hawp/kit", rel)),
			Type:       "guide",
			Category:   "kit",
			FolderRole: "kit/" + filepath.ToSlash(filepath.Dir(rel)),
			Content:    string(content),
			Metadata:   map[string]interface{}{"file": filepath.Base(path)},
		})
		return nil
	})
}

func walkWorkFiles(workPath string, corpus *appindex.EnrichedCorpus) error {
	return filepath.Walk(workPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !info.Mode().IsRegular() || filepath.Ext(path) != ".md" {
			return err
		}

		// Skip BACKLOG for now (complex parsing)
		if filepath.Base(path) == "BACKLOG.md" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read work document %s: %w", path, err)
		}
		rel, err := filepath.Rel(workPath, path)
		if err != nil {
			return err
		}

		// Determine folder role
		parts := strings.Split(rel, string(filepath.Separator))
		folderRole := "work"
		var status *string
		var workUUID *string
		if len(parts) > 1 {
			folderRole = "work/" + parts[0]
			switch parts[0] {
			case "active", "parked":
				status = strPtr(parts[0])
				workUUID = canonicalWorkUUID(parts[1])
			case "closed":
				status = strPtr(parts[0])
				if len(parts) > 4 {
					workUUID = canonicalWorkUUID(parts[4])
				}
			}
		}

		corpus.Documents = append(corpus.Documents, appindex.EnrichedDocument{
			Path:       filepath.ToSlash(filepath.Join(".hawp/work", rel)),
			Type:       "plan",
			Category:   "work",
			FolderRole: folderRole,
			Content:    string(content),
			Status:     status,
			WorkUUID:   workUUID,
			Metadata:   map[string]interface{}{"file": filepath.Base(path)},
		})
		return nil
	})
}

// walkCustomPath walks a user-configured path (file or directory) and adds any
// .md files to the corpus under the "custom" category.
func walkCustomPath(abs, configuredPath string, corpus *appindex.EnrichedCorpus) error {
	info, err := os.Lstat(abs)
	if os.IsNotExist(err) {
		fmt.Printf("warning: configured index path not found, skipping: %s\n", configuredPath)
		return nil
	}
	if err != nil {
		return err
	}

	if !info.IsDir() {
		// Single file — index it directly if it's a .md file.
		if !info.Mode().IsRegular() || filepath.Ext(abs) != ".md" {
			return nil
		}
		content, err := os.ReadFile(abs)
		if err != nil {
			return fmt.Errorf("read custom document %s: %w", abs, err)
		}
		corpus.Documents = append(corpus.Documents, appindex.EnrichedDocument{
			Path:       filepath.ToSlash(configuredPath),
			Type:       "document",
			Category:   "custom",
			FolderRole: "custom/" + filepath.ToSlash(filepath.Dir(configuredPath)),
			Content:    string(content),
			Metadata:   map[string]interface{}{"file": filepath.Base(abs)},
		})
		return nil
	}

	return filepath.Walk(abs, func(path string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() || !fi.Mode().IsRegular() || filepath.Ext(path) != ".md" {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read custom document %s: %w", path, err)
		}
		rel, err := filepath.Rel(abs, path)
		if err != nil {
			return err
		}
		docPath := filepath.ToSlash(filepath.Join(configuredPath, rel))
		corpus.Documents = append(corpus.Documents, appindex.EnrichedDocument{
			Path:       docPath,
			Type:       "document",
			Category:   "custom",
			FolderRole: "custom/" + filepath.ToSlash(filepath.Dir(rel)),
			Content:    string(content),
			Metadata:   map[string]interface{}{"file": filepath.Base(path)},
		})
		return nil
	})
}

func canonicalWorkUUID(value string) *string {
	if short := identity.ExtractShortUUID(value); short != "" {
		return &short
	}
	if identity.IsFullUUID(value) {
		full := strings.ToLower(value)
		return &full
	}
	return nil
}

func strPtr(s string) *string {
	return &s
}
