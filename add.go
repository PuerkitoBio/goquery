package goquery

import (
	"exp/html"
)

// Adds matching nodes to the current selection. Returns the same Selection object.
// The new selector string is run in the context of the document of the Selection object.
func (this *Selection) Add(selector string) *Selection {
	this.Nodes = appendWithoutDuplicates(this.Nodes, findWithContext(selector, this.document.Root))
	return this
}

// Adds nodes of the specified Selection object to the current selection. Returns the same Selection object.
func (this *Selection) AddSelection(sel *Selection) *Selection {
	this.Nodes = appendWithoutDuplicates(this.Nodes, sel.Nodes)
	return this
}

// Adds nodes to the current selection. Returns the same Selection object.
// No verification about the same document origin is done.
func (this *Selection) AddNodes(nodes []*html.Node) *Selection {
	this.Nodes = appendWithoutDuplicates(this.Nodes, nodes)
	return this
}
