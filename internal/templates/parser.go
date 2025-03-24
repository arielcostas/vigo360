package templates

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"vigo360.es/new/internal/goldmark_extensions"
)

var parser goldmark.Markdown = goldmark.New(
	goldmark.WithExtensions(extension.Footnote),
	goldmark.WithExtensions(extension.Typographer),
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithExtensions(extension.Strikethrough),
	goldmark.WithExtensions(goldmark_extensions.Extension),
)
