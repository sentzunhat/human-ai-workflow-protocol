package cli

import (
	"encoding/json"
	"fmt"
	"strings"
)

// renderCommandsText prints the registry as human-readable text, grouped
// by availability.
func renderCommandsText() string {
	var b strings.Builder
	b.WriteString("hawp commands\n=============\n")
	b.WriteString("Exit codes: " + ExitCodeConvention + "\n\n")

	for _, cmd := range Registry {
		if !cmd.Available {
			continue
		}
		fmt.Fprintf(&b, "%s\n", cmd.Usage)
		fmt.Fprintf(&b, "  %s\n", cmd.Description)
		for _, flag := range cmd.Flags {
			fmt.Fprintf(&b, "  - %s\n", flag)
		}
		fmt.Fprintf(&b, "  exit codes: %s\n\n", cmd.ExitCodes)
	}

	var planned []string
	for _, cmd := range Registry {
		if !cmd.Available {
			planned = append(planned, cmd.Usage+" — "+cmd.Description)
		}
	}
	if len(planned) > 0 {
		b.WriteString("Planned (not yet available):\n")
		for _, line := range planned {
			fmt.Fprintf(&b, "  %s\n", line)
		}
	}
	return b.String()
}

// renderCommandsJSON marshals the full registry plus the shared exit-code
// convention as a single JSON object — the canonical agent-discovery
// output for `hawp commands --json`.
func renderCommandsJSON() (string, error) {
	payload := struct {
		ExitCodeConvention string        `json:"exitCodeConvention"`
		Commands           []CommandInfo `json:"commands"`
	}{ExitCodeConvention: ExitCodeConvention, Commands: Registry}

	out, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out) + "\n", nil
}

func runCommands(args []string) error {
	if len(args) >= 1 && args[0] == "--json" {
		out, err := renderCommandsJSON()
		if err != nil {
			return err
		}
		fmt.Print(out)
		return nil
	}
	fmt.Print(renderCommandsText())
	return nil
}
