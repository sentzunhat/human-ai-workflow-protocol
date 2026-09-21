package intake

import (
	"fmt"
	"strings"
)

// Draft is a proposed HAWP intake shape, not a work item or approval.
// Input and Context are supplied by the caller; the other fields need review.
type Draft struct {
	Input       string `json:"input"`
	Context     string `json:"context"`
	Mission     string `json:"mission"`
	Constraints string `json:"constraints"`
	Output      string `json:"output"`
	Checkpoint  string `json:"checkpoint,omitempty"`
}

// Validate checks completeness, not factual accuracy or authorization.
func (d Draft) Validate() error {
	for _, field := range []struct{ name, value string }{
		{"input", d.Input}, {"context", d.Context}, {"mission", d.Mission},
		{"constraints", d.Constraints}, {"output", d.Output},
	} {
		if strings.TrimSpace(field.value) == "" {
			return fmt.Errorf("draft %s must not be blank", field.name)
		}
	}
	return nil
}
