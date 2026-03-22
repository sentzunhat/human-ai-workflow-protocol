package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// helpMentionsBaseCommand checks helpText() references the command's base
// verb (before any flags/args), so registry and help text cannot silently
// drift apart.
func helpMentionsBaseCommand(usage string) bool {
	base := strings.Fields(usage)
	if len(base) < 2 {
		return false
	}
	// usage is "hawp <name...>"; take everything up to the first flag/arg
	// token (starts with "-" or "<").
	var verb []string
	for _, tok := range base[1:] {
		if strings.HasPrefix(tok, "-") || strings.HasPrefix(tok, "<") || strings.HasPrefix(tok, "[") {
			break
		}
		verb = append(verb, tok)
	}
	return strings.Contains(helpText(), strings.Join(verb, " "))
}

func TestRegistryStaysInSyncWithHelpText(t *testing.T) {
	for _, cmd := range Registry {
		if !cmd.Available {
			continue // planned commands are listed separately in help text
		}
		if !helpMentionsBaseCommand(cmd.Usage) {
			t.Errorf("command %q (usage %q) not found in helpText() — registry and help drifted", cmd.Name, cmd.Usage)
		}
	}
}

func TestRunCommandsTextListsAvailableCommands(t *testing.T) {
	if err := Run([]string{"commands"}); err != nil {
		t.Fatalf("Run(commands) error: %v", err)
	}
	text := renderCommandsText()
	if !strings.Contains(text, "hawp uuid") || !strings.Contains(text, "hawp check") {
		t.Errorf("commands text missing expected entries:\n%s", text)
	}
	if !strings.Contains(text, ExitCodeConvention) {
		t.Error("commands text should document the exit-code convention")
	}
}

func TestRunCommandsJSONIsValidAndComplete(t *testing.T) {
	out, err := renderCommandsJSON()
	if err != nil {
		t.Fatal(err)
	}

	var parsed struct {
		ExitCodeConvention string        `json:"exitCodeConvention"`
		Commands           []CommandInfo `json:"commands"`
	}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, out)
	}
	if parsed.ExitCodeConvention == "" {
		t.Error("exitCodeConvention missing from JSON output")
	}
	if len(parsed.Commands) != len(Registry) {
		t.Errorf("JSON commands count = %d, want %d", len(parsed.Commands), len(Registry))
	}

	names := map[string]bool{}
	for _, cmd := range parsed.Commands {
		names[cmd.Name] = true
	}
	for _, want := range []string{"uuid", "check", "init", "update", "commands", "work normalize"} {
		if !names[want] {
			t.Errorf("JSON output missing command %q", want)
		}
	}
}

func TestRunCommandsRoutesJSONFlag(t *testing.T) {
	if err := Run([]string{"commands", "--json"}); err != nil {
		t.Fatalf("Run(commands --json) error: %v", err)
	}
}

func TestPlannedCommandsAreMarkedUnavailable(t *testing.T) {
	searchFound := false
	contextFound := false
	for _, cmd := range Registry {
		if cmd.Name == "search" {
			searchFound = true
			if !cmd.Available {
				t.Errorf("command search should be Available: true (Phase 4 complete)")
			}
		}
		if cmd.Name == "context" {
			contextFound = true
			if cmd.Available {
				t.Errorf("command context should be Available: false (alias for search --context)")
			}
		}
	}
	if !searchFound || !contextFound {
		t.Fatal("expected search and context entries in the registry")
	}
}
