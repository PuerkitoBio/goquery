package goquery

import (
	"exp/html"
	esc "html"
)

// SetText() replaces the children of each selected node with the given text 
// (properly escaped of course).
// This is the same behavior as jQuery's .text() function.
func (this *Selection) SetText(s string) *Selection {
	for _, n := range this.Nodes {
		setNodeText(n, esc.EscapeString(s))
	}
	return this
}

// Replace the given node's children with the given string.
func setNodeText(node *html.Node, s string) {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		node.RemoveChild(c)
	}
	text := &html.Node{
		Type: html.TextNode,
		Data: s,
	}
	node.AppendChild(text)
}
