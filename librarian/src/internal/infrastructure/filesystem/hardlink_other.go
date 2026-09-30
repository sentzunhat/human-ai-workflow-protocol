//go:build !aix && !darwin && !dragonfly && !freebsd && !linux && !netbsd && !openbsd && !solaris && !windows

package filesystem

import "fmt"

// RejectHardLinkedFile fails closed where the platform link count is unknown.
func RejectHardLinkedFile(path string) error {
	return fmt.Errorf("cannot verify hard-link count for %s on this platform", path)
}
