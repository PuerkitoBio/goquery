package goquery

import (
	"code.google.com/p/go.net/html"
)

// Add() adds the selector string's matching nodes to those in the current
// selection and returns a new Selection object.
// The selector string is run in the context of the document of the current
// Selection object.
func (this *Selection) Add(selector string) *Selection {
	return this.AddNodes(findWithSelector([]*html.Node{this.document.rootNode}, selector)...)
}

// AddSelection() adds the specified Selection object's nodes to those in the
// current selection and returns a new Selection object.
func (this *Selection) AddSelection(sel *Selection) *Selection {
	if sel == nil {
		return this.AddNodes()
	}
	return this.AddNodes(sel.Nodes...)
}

// Union() is an alias for AddSelection().
func (this *Selection) Union(sel *Selection) *Selection {
	return this.AddSelection(sel)
}

// AddNodes() adds the specified nodes to those in the
// current selection and returns a new Selection object.
func (this *Selection) AddNodes(nodes ...*html.Node) *Selection {
	return pushStack(this, appendWithoutDuplicates(this.Nodes, nodes))
}

// AndSelf() adds the previous set of elements on the stack to the current set.
// It returns a new Selection object containing the current Selection combined
// with the previous one.
func (this *Selection) AndSelf() *Selection {
	return this.AddSelection(this.prevSel)
}
