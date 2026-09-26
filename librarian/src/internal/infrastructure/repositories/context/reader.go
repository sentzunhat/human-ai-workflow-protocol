// Package contextinf provides filesystem reader adapters for
// application-level calls into domain/context functions.
package contextinf

import "os"

// ReadFile wraps os.ReadFile for dependency injection into domain
// context enrichment functions (see internal/domain/context).
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
