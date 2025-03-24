// SPDX-License-Identifier: MIT
package goldmark_extensions

import (
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type render struct {
	html.Config
}

// newRenderer return new renderer for image figure.
func newRenderer(opts ...html.Option) renderer.NodeRenderer {
	var r = &render{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

var defaultRenderer = newRenderer()

// RegisterFuncs implement renderer.NodeRenderer interface.
func (r *render) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindImage, r.renderImage)
	reg.Register(ast.KindLink, r.renderLink)
}

func (r *render) renderParagraph(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	var (
		name     = "p"
		filter   = html.ParagraphAttributeFilter
		isFigure = isFigure(n)
	)
	if isFigure {
		name = "figure"
		filter = html.GlobalAttributeFilter
	}
	if entering {
		w.WriteByte('<')
		w.WriteString(name)
		if n.Attributes() != nil {
			html.RenderAttributes(w, n, filter)
		}
		w.WriteByte('>')
		if isFigure {
			w.WriteByte('\n')
		}
	} else {
		w.WriteString("</")
		w.WriteString(name)
		w.WriteString(">\n")
	}
	return ast.WalkContinue, nil
}

func (r *render) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		// special case for figure
		if node.HasChildren() && isFigure(node.Parent()) {
			w.WriteString("</figcaption>\n")
		}
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Image)
	w.WriteString("<img src=\"")
	if r.Unsafe || !html.IsDangerousURL(n.Destination) {
		w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
	}
	w.WriteString(`" alt="`)
	w.Write(n.Title)
	w.WriteString(`" loading="lazy"`)
	if n.Title != nil {
		w.WriteString(` title="`)
		r.Writer.Write(w, n.Title)
		w.WriteByte('"')
	}
	if n.Attributes() != nil {
		html.RenderAttributes(w, n, html.ImageAttributeFilter)
	}
	if r.XHTML {
		w.WriteString(" />")
	} else {
		w.WriteString(">")
	}
	// special case for figure
	if node.HasChildren() && isFigure(node.Parent()) {
		w.WriteString("\n<figcaption>")
		return ast.WalkContinue, nil
	}
	return ast.WalkSkipChildren, nil
}

// Update renderLink method to add rel="nofollow noreferrer" to external links
func (r *render) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		w.WriteString("</a>")
		return ast.WalkContinue, nil
	}
	link := node.(*ast.Link)
	w.WriteString("<a href=\"")
	w.Write(util.EscapeHTML(util.URLEscape(link.Destination, true)))
	w.WriteString("\"")
	if link.Title != nil {
		w.WriteString(" title=\"")
		w.Write(link.Title)
		w.WriteString("\"")
	}
	// Check if the link is external
	if !isRelativeURL(link.Destination) && !isInternalURL(link.Destination) {
		w.WriteString(" class=\"external\" rel=\"nofollow noreferrer\"")
	}
	w.WriteString(">")
	return ast.WalkContinue, nil
}

func isRelativeURL(url []byte) bool {
	return len(url) > 0 && url[0] == '/'
}

func isInternalURL(url []byte) bool {
	domain := os.Getenv("DOMAIN")
	return len(domain) > 0 && filepath.HasPrefix(string(url), domain)
}

func isFigure(node ast.Node) bool {
	var child = node.FirstChild()
	return node.Kind() == ast.KindParagraph &&
		child != nil &&
		child == node.LastChild() &&
		child.Kind() == ast.KindImage &&
		child.HasChildren()
}

// extension defines a goldmark.Extender for markdown image figures.
type extension struct{}

// Extend implement goldmark.Extender interface.
func (a *extension) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(defaultRenderer, 100),
		),
	)
}

// Extension is a goldmark.Extender with markdown image figures support.
var Extension goldmark.Extender = new(extension)
