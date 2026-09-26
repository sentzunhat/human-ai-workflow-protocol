package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
)

// ExtractAll extracts every regular file and directory from a .tar.gz
// archive into destDir, preserving its internal directory structure.
// Unlike ExtractMember, this pulls out the whole tree (used for the kit +
// providers bundle, not a single named file).
func ExtractAll(archivePath, destDir string) error {
	root, err := filepath.Abs(destDir)
	if err != nil {
		return fmt.Errorf("resolve extraction root: %w", err)
	}
	parent := filepath.Dir(root)
	if err := filesystem.RejectSymlinkAncestors(parent, root); err != nil {
		return fmt.Errorf("refusing symlinked extraction root: %w", err)
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return err
	}
	if err := filesystem.RejectSymlinkAncestors(parent, root); err != nil {
		return fmt.Errorf("refusing symlinked extraction root: %w", err)
	}

	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		target, err := archiveTarget(root, header.Name)
		if err != nil {
			return err
		}
		if err := filesystem.RejectSymlinkAncestors(root, target); err != nil {
			return fmt.Errorf("refusing symlinked archive target %q: %w", header.Name, err)
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			out, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}
			out.Close()
		}
	}
}

// archiveTarget resolves a tar member below destDir without allowing an
// absolute or parent-traversal member name to escape the extraction root.
func archiveTarget(destDir, memberName string) (string, error) {
	if memberName == "" {
		return "", fmt.Errorf("archive member has an empty name")
	}
	name := filepath.FromSlash(memberName)
	if filepath.IsAbs(name) {
		return "", fmt.Errorf("archive member %q is absolute", memberName)
	}

	root, err := filepath.Abs(destDir)
	if err != nil {
		return "", fmt.Errorf("resolve extraction root: %w", err)
	}
	target := filepath.Join(root, name)
	relative, err := filepath.Rel(root, target)
	if err != nil {
		return "", fmt.Errorf("check archive member %q: %w", memberName, err)
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("archive member %q escapes extraction root", memberName)
	}
	return target, nil
}
