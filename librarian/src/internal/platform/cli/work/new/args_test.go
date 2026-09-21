package newcmd

import (
	"strings"
	"testing"
)

func TestParseNewArgsValid(t *testing.T) {
	cases := []struct {
		args                         []string
		title, itemType, input, root string
	}{
		{args: []string{"Add search API"}, title: "Add search API", itemType: "task"},
		{args: []string{"Fix crash", "--type", "bug"}, title: "Fix crash", itemType: "bug"},
		{args: []string{"Ship v1", "--type", "release", "--input", "original request"}, title: "Ship v1", itemType: "release", input: "original request"},
		{args: []string{"My task", "--hawp-root", "/some/path"}, title: "My task", itemType: "task", root: "/some/path"},
	}
	for _, tc := range cases {
		t.Run(strings.Join(tc.args, " "), func(t *testing.T) {
			got, err := parseNewArgs(tc.args)
			if err != nil {
				t.Fatal(err)
			}
			if got.title != tc.title || got.itemType != tc.itemType || got.input != tc.input || got.hawpRoot != tc.root {
				t.Fatalf("got %+v", got)
			}
		})
	}
}

func TestParseNewArgsRejectsInvalidInput(t *testing.T) {
	for _, args := range [][]string{nil, {"--type", "task"}, {"My task", "--type", "chore"}, {"My task", "--unknown"}, {"My task", "extra"}, {"   "}} {
		if _, err := parseNewArgs(args); err == nil {
			t.Errorf("accepted %q", args)
		}
	}
	for _, value := range []string{"", " ", "\t\n"} {
		if _, err := parseNewArgs([]string{"Task", "--hawp-root=" + value}); err == nil {
			t.Errorf("accepted empty root %q", value)
		}
	}
}

func TestParseNewArgsRejectsControlCharactersInTitle(t *testing.T) {
	for _, title := range []string{"line\nbreak", "line\rbreak", "tab\tbreak"} {
		if _, err := parseNewArgs([]string{title}); err == nil {
			t.Errorf("accepted title with control character %q", title)
		}
	}
}
