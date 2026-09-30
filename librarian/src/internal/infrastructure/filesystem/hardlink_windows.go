//go:build windows

package filesystem

import (
	"fmt"
	"os"
	"syscall"
)

// RejectHardLinkedFile refuses an existing file with multiple directory links.
func RejectHardLinkedFile(path string) error {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()

	var info syscall.ByHandleFileInformation
	if err := syscall.GetFileInformationByHandle(syscall.Handle(file.Fd()), &info); err != nil {
		return fmt.Errorf("inspect link count for %s: %w", path, err)
	}
	if info.NumberOfLinks > 1 {
		return fmt.Errorf("%s has %d hard links; refusing in-place database access", path, info.NumberOfLinks)
	}
	return nil
}
