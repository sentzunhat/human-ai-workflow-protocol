package cli_test

import (
	"strings"
	"testing"

	platformcli "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli"
)

func TestSearchRejectsInvalidHybridRatio(t *testing.T) {
	for _, value := range []string{"NaN", "+Inf", "-Inf", "-0.1", "1.1", "invalid"} {
		t.Run(value, func(t *testing.T) {
			err := platformcli.Run([]string{"search", "test", "--hybrid-ratio", value, "--no-update-check"})
			if err == nil || !strings.Contains(err.Error(), "hybrid-ratio") {
				t.Fatalf("expected hybrid ratio validation error, got %v", err)
			}
		})
	}
}
