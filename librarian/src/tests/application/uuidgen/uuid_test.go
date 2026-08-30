package uuidgen_test

import (
	"regexp"
	"testing"

	appuuid "github.com/sentzunhat/hawp/librarian/src/internal/application/uuidgen"
)

var uuidV4Re = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func TestNewIsV4(t *testing.T) {
	seen := map[string]struct{}{}
	for i := 0; i < 50; i++ {
		uuid, err := appuuid.New()
		if err != nil {
			t.Fatal(err)
		}
		if !uuidV4Re.MatchString(uuid) {
			t.Fatalf("not a canonical lowercase UUID v4: %q", uuid)
		}
		if _, dup := seen[uuid]; dup {
			t.Fatalf("duplicate UUID generated: %q", uuid)
		}
		seen[uuid] = struct{}{}
	}
}

func TestShort(t *testing.T) {
	if got := appuuid.Short("0e1c4afa-9668-4d61-b5b6-1e27be42ca23"); got != "0e1c4afa" {
		t.Errorf("Short = %q, want 0e1c4afa", got)
	}
	if got := appuuid.Short("abc"); got != "abc" {
		t.Errorf("Short on short input = %q, want unchanged", got)
	}
}
