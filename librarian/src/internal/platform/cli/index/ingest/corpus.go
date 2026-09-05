package ingest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	appindex "github.com/sentzunhat/hawp/librarian/src/internal/application/index"
)

// buildCorpusFromRepo walks each configured path and builds an enriched corpus.
// Paths are relative to repoRoot. ".hawp/kit" and ".hawp/work" get their
// enriched walkers; all other paths are walked generically as "custom" corpus.
func buildCorpusFromRepo(repoRoot string, paths []string) (*appindex.EnrichedCorpus, error) {
	corpus := &appindex.EnrichedCorpus{}

	for _, p := range paths {
		abs := filepath.Join(repoRoot, filepath.FromSlash(p))
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

func walkKitFiles(kitPath string, corpus *appindex.EnrichedCorpus) error {
	return filepath.Walk(kitPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(path) != ".md" {
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
		if err != nil || info.IsDir() || filepath.Ext(path) != ".md" {
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
		if len(parts) > 1 {
			folderRole = "work/" + parts[0]
			switch parts[0] {
			case "active", "parked", "closed":
				status = strPtr(parts[0])
			}
		}

		corpus.Documents = append(corpus.Documents, appindex.EnrichedDocument{
			Path:       filepath.ToSlash(filepath.Join(".hawp/work", rel)),
			Type:       "plan",
			Category:   "work",
			FolderRole: folderRole,
			Content:    string(content),
			Status:     status,
			Metadata:   map[string]interface{}{"file": filepath.Base(path)},
		})
		return nil
	})
}

// walkCustomPath walks a user-configured path (file or directory) and adds any
// .md files to the corpus under the "custom" category.
func walkCustomPath(abs, configuredPath string, corpus *appindex.EnrichedCorpus) error {
	info, err := os.Stat(abs)
	if os.IsNotExist(err) {
		fmt.Printf("warning: configured index path not found, skipping: %s\n", configuredPath)
		return nil
	}
	if err != nil {
		return err
	}

	if !info.IsDir() {
		// Single file — index it directly if it's a .md file.
		if filepath.Ext(abs) != ".md" {
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
		if err != nil || fi.IsDir() || filepath.Ext(path) != ".md" {
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

func strPtr(s string) *string {
	return &s
}
