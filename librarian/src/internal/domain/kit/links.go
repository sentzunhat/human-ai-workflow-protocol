package kit

import "github.com/sentzunhat/hawp/librarian/src/internal/domain/kit/markdown"

type Link = markdown.Link

func ExtractLinks(content string) []Link { return markdown.ExtractLinks(content) }
func IsLocalHref(href string) bool       { return markdown.IsLocalHref(href) }
func PathPart(href string) string        { return markdown.PathPart(href) }
