package query

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestTruncatePreviewPreservesUTF8(t *testing.T) {
	text := strings.Repeat("界", 60)
	got := truncatePreview(text, 150)

	if !utf8.ValidString(got) {
		t.Fatalf("truncatePreview returned invalid UTF-8: %q", got)
	}
	if got == text || !strings.HasSuffix(got, "...") {
		t.Fatalf("truncatePreview(%q) = %q, want truncated output", text, got)
	}
}
