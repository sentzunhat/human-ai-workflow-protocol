package update

import (
	"fmt"
	"runtime"
)

// AssetName returns the release binary name for the current platform,
// matching the naming produced by `make dist` (librarian/go/Makefile).
func AssetName() (string, error) {
	os, arch := runtime.GOOS, runtime.GOARCH
	switch os + "/" + arch {
	case "darwin/arm64", "darwin/amd64",
		"linux/arm64", "linux/amd64":
		return fmt.Sprintf("hawp-%s-%s", os, arch), nil
	case "windows/amd64", "windows/arm64":
		return fmt.Sprintf("hawp-%s-%s.exe", os, arch), nil
	default:
		return "", fmt.Errorf("no release asset published for %s/%s", os, arch)
	}
}
