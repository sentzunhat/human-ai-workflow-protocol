// Package exitcode defines command failures that have already been rendered.
package exitcode

import "fmt"

// Error carries the process exit code for a command that already reported its
// findings to the user.
type Error struct{ Code int }

func (e Error) Error() string { return fmt.Sprintf("exit %d", e.Code) }
