// Package hnhtml implements an HTML-to-Pango-markup translator that translates
// a small subset of HTML, enough to view HackerNews posts.
package hnhtml

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// ToMarkup renders the HTML source into Pango markup.
func ToMarkup(source string) string {
	n, err := html.Parse(strings.NewReader(source))
	if err != nil {
		return html.EscapeString(source)
	}

	r := markupRenderer{}
	r.buf.Grow(len(source))

	if r.traverseSiblings(n) == traverseFailed {
		return html.EscapeString(source)
	}

	return r.buf.String()
}

type traverseStatus uint8

const (
	traverseOK traverseStatus = iota
	traverseSkipChildren
	traverseFailed
)

var allowedURLSchemes = map[string]bool{
	"http":  true,
	"https": true,
}

type markupRenderer struct {
	buf   strings.Builder
	state struct {
		list  int
		quote bool
	}
}

func (r *markupRenderer) renderNode(n *html.Node) traverseStatus {
	switch n.Type {
	case html.TextNode:
		r.buf.WriteString(n.Data)
		return traverseOK

	case html.ElementNode:
		switch n.Data {
		case "h1":
			return r.wrap(n, `<span weight="bold" size="xx-large">`, "</span>")
		case "h2":
			return r.wrap(n, `<span weight="bold" size="x-large">`, "</span>")
		case "h3":
			return r.wrap(n, `<span weight="bold" size="large">`, "</span>")
		case "h4":
			return r.wrap(n, `<span weight="bold">`, "</span>")
		case "h5":
			return r.wrap(n, `<span weight="bold" size="small">`, "</span>")
		case "h6":
			return r.wrap(n, `<span weight="bold" size="x-small">`, "</span>")
		case "em", "i":
			return r.wrap(n, `<span style="italic">`, "</span>")
		case "strong", "b":
			return r.wrap(n, `<span weight="bold">`, "</span>")
		case "underline", "u":
			return r.wrap(n, `<span underline="single">`, "</span>")
		case "strike", "del":
			return r.wrap(n, `<span strikethrough="true">`, "</span>")
		case "sub":
			return r.wrap(n, `<span font_scale="0.7" rise="6pt">`, "</span>")
		case "sup":
			return r.wrap(n, `<span font_scale="0.7" rise="-2pt">`, "</span>")
		case "code":
			return r.wrap(n, `<span font_family="monospace">`, "</span>")
		case "blockquote":
			return r.wrapFunc(n,
				func() { r.state.quote = true },
				func() { r.state.quote = false },
			)
		case "p", "div":
			r.buf.WriteString("\n\n")
			if r.state.quote {
				r.buf.WriteString("> ")
			}
			return traverseOK
		case "a":
			href := nodeAttr(n, "href")
			if !allowedURLSchemes[href] {
				return traverseOK
			}
			return r.wrap(n, fmt.Sprintf(`<a href="%s">`, html.EscapeString(href)), "</a>")
		case "ol":
			return r.wrapFunc(n,
				func() { r.state.list = 1 },
				func() { r.state.list = 0 },
			)
		case "ul":
			r.state.list = 0
			return traverseOK
		case "li":
			bullet := "● "
			if r.state.list != 0 {
				bullet = strconv.Itoa(r.state.list) + ". "
				r.state.list++
			}
			r.buf.WriteString(bullet)
			return traverseOK
		case "hr":
			r.buf.WriteString("―")
			return traverseOK
		case "br":
			r.buf.WriteString("\n\n")
			return traverseOK
		case "img":
			alt := nodeAttr(n, "alt")
			src := nodeAttr(n, "src")
			if src == "" {
				return traverseOK
			}
			if alt != "" {
				fmt.Fprintf(&r.buf,
					`<a href="%s">%s</a>`,
					html.EscapeString(alt),
					html.EscapeString(src),
				)
			} else {
				fmt.Fprintf(&r.buf,
					`<a href="%s">[image]</a>`,
					html.EscapeString(src),
				)
			}
			return traverseOK
		default:
			return traverseOK
		}
	case html.ErrorNode:
		return traverseFailed
	}

	return traverseOK
}

func (r *markupRenderer) wrap(n *html.Node, head, tail string) traverseStatus {
	return r.wrapFunc(n,
		func() { r.buf.WriteString(head) },
		func() { r.buf.WriteString(tail) },
	)
}

func (r *markupRenderer) wrapFunc(n *html.Node, head, tail func()) traverseStatus {
	head()
	status := r.traverseChildren(n)
	tail()

	if status == traverseFailed {
		return traverseFailed
	}
	return traverseSkipChildren
}

func (r *markupRenderer) traverseChildren(n *html.Node) traverseStatus {
	return r.traverseSiblings(n.FirstChild)
}

func (r *markupRenderer) traverseSiblings(first *html.Node) traverseStatus {
	for n := first; n != nil; n = n.NextSibling {
		switch r.renderNode(n) {
		case traverseOK:
			// traverseChildren never returns traverseSkipChildren.
			if r.traverseChildren(n) == traverseFailed {
				return traverseFailed
			}
		case traverseSkipChildren:
			continue
		case traverseFailed:
			return traverseFailed
		}
	}

	return traverseOK
}

func nodeAttr(n *html.Node, keys ...string) string {
	for _, attr := range n.Attr {
		for _, k := range keys {
			if k == attr.Key {
				return attr.Val
			}
		}
	}
	return ""
}
