//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package filesystem

import (
	"fmt"
	"os"
	"syscall"
)

// RejectHardLinkedFile refuses an existing file with multiple directory links.
func RejectHardLinkedFile(path string) error {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("cannot inspect link count for %s", path)
	}
	if stat.Nlink > 1 {
		return fmt.Errorf("%s has %d hard links; refusing in-place database access", path, stat.Nlink)
	}
	return nil
}
