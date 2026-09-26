package markdown

import "testing"

func TestIsLocalHrefRejectsURISchemes(t *testing.T) {
	for _, href := range []string{
		"https://example.test",
		"mailto:hello@example.test",
		"tel:+12045550123",
		"data:text/plain,hello",
		"ftp://example.test/file",
		"file:///tmp/example.md",
	} {
		if IsLocalHref(href) {
			t.Errorf("IsLocalHref(%q) = true, want false", href)
		}
	}
	if !IsLocalHref("../usage/search.md#search-modes") {
		t.Fatal("relative kit path should remain local")
	}
}
