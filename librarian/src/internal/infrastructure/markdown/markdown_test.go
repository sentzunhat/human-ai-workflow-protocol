package markdown

import "testing"

func TestIsLocalHrefRejectsURISchemes(t *testing.T) {
	for _, href := range []string{
		"https://example.test",
		"http://example.test",
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
	for _, href := range []string{
		"guide.md",
		"docs/guide.md#section",
		"../README.md",
	} {
		if !IsLocalHref(href) {
			t.Errorf("IsLocalHref(%q) = false, want true", href)
		}
	}
}
