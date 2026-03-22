// Package selfreplace atomically swaps the currently running executable
// for a new one, without ever leaving a partially-written or non-executable
// file at the original path.
package selfreplace

import (
	"os"
	"runtime"
)

// ErrUnsupportedOnWindows is returned by Replace on Windows, where an
// executable cannot be renamed over while it is running. Windows users
// must download and swap the binary manually (or via an external
// installer) — see `hawp update` output for the guidance shown.
var ErrUnsupportedOnWindows = &unsupportedError{}

type unsupportedError struct{}

func (*unsupportedError) Error() string {
	return "self-replace is not supported on Windows; download the new binary and replace it manually while hawp is not running"
}

// Replace atomically overwrites targetPath (typically the running
// executable, from os.Executable()) with sourcePath's content. sourcePath
// must be on the same filesystem as targetPath (same directory is
// simplest) so the final step is a plain rename, never a partial copy.
// sourcePath is consumed (renamed away) on success.
func Replace(sourcePath, targetPath string) error {
	if runtime.GOOS == "windows" {
		return ErrUnsupportedOnWindows
	}
	if err := os.Chmod(sourcePath, 0o755); err != nil {
		return err
	}
	return os.Rename(sourcePath, targetPath)
}
