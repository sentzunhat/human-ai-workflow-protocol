// Package archive extracts a single named member from a .tgz or .zip
// archive — used to pull just the ONNX Runtime shared library out of its
// release tarball/zip without keeping the rest (headers, debug symbols).
package archive

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExtractMember opens archivePath (.tgz or .zip, chosen by extension) and
// writes the single member matching memberPath to destPath.
func ExtractMember(archivePath, memberPath, destPath string) error {
	switch {
	case strings.HasSuffix(archivePath, ".tgz") || strings.HasSuffix(archivePath, ".tar.gz"):
		return extractFromTarGz(archivePath, memberPath, destPath)
	case strings.HasSuffix(archivePath, ".zip"):
		return extractFromZip(archivePath, memberPath, destPath)
	default:
		return fmt.Errorf("unsupported archive type: %s", archivePath)
	}
}

func extractFromTarGz(archivePath, memberPath, destPath string) error {
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
			return fmt.Errorf("member %q not found in %s", memberPath, archivePath)
		}
		if err != nil {
			return err
		}
		if header.Typeflag != tar.TypeReg || header.Name != memberPath {
			continue
		}
		return writeMember(destPath, tr)
	}
}

func extractFromZip(archivePath, memberPath, destPath string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, file := range r.File {
		if file.Name != memberPath {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		return writeMember(destPath, rc)
	}
	return fmt.Errorf("member %q not found in %s", memberPath, archivePath)
}

func writeMember(destPath string, r io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return err
	}
	temp, err := os.CreateTemp(filepath.Dir(destPath), ".extract-*")
	if err != nil {
		return err
	}
	tempPath := temp.Name()
	defer os.Remove(tempPath)

	if _, err := io.Copy(temp, r); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tempPath, 0o755); err != nil {
		return err
	}
	return os.Rename(tempPath, destPath)
}
