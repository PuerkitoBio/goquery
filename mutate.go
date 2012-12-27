package goquery

import (
	"exp/html"
	esc "html"
)

// SetText() replaces the children of each selected node with the given text 
// (properly escaped of course).
// This is the same behavior as jQuery's .text() function.
func (this *Selection) SetText(s string) *Selection {
	escapedText := esc.EscapeString(s)
	for _, n := range this.Nodes {
		setNodeText(n, escapedText)
	}
	return this
}

// Replace the given node's children with the given string.
func setNodeText(node *html.Node, s string) {
	// remove all existing children
	for node.FirstChild != nil {
		node.RemoveChild(node.FirstChild)
	}
	// add the text
	node.AppendChild(&html.Node{
		Type: html.TextNode,
		Data: s,
	})
}
