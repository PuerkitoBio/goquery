package goquery

import (
	"code.google.com/p/go.net/html"
)

// First() reduces the set of matched elements to the first in the set.
// It returns a new Selection object.
func (this *Selection) First() *Selection {
	return this.Eq(0)
}

// Last() reduces the set of matched elements to the last in the set.
// It returns a new Selection object.
func (this *Selection) Last() *Selection {
	return this.Eq(-1)
}

// Eq() reduces the set of matched elements to the one at the specified index.
// If a negative index is given, it counts backwards starting at the end of the
// set. It returns a new Selection object, and an empty Selection object if the
// index is invalid.
func (this *Selection) Eq(index int) *Selection {
	if index < 0 {
		index += len(this.Nodes)
	}
	return this.Slice(index, index+1)
}

// Slice() reduces the set of matched elements to a subset specified by a range
// of indices.
func (this *Selection) Slice(start int, end int) *Selection {
	if start < 0 {
		start += len(this.Nodes)
	}
	if end < 0 {
		end += len(this.Nodes)
	}
	return pushStack(this, this.Nodes[start:end])
}

// Get() retrieves the underlying node at the specified index.
// Get() without parameter is not implemented, since the node array is available
// on the Selection object.
func (this *Selection) Get(index int) *html.Node {
	if index < 0 {
		index += len(this.Nodes) // Negative index gets from the end
	}
	return this.Nodes[index]
}

// Index() returns the position of the first element within the Selection object
// relative to its sibling elements.
func (this *Selection) Index() int {
	if len(this.Nodes) > 0 {
		return newSingleSelection(this.Nodes[0], this.document).PrevAll().Length()
	}
	return -1
}

// IndexSelector() returns the position of the first element within the
// Selection object relative to the elements matched by the selector, or -1 if
// not found.
func (this *Selection) IndexSelector(selector string) int {
	if len(this.Nodes) > 0 {
		sel := this.document.Find(selector)
		return indexInSlice(sel.Nodes, this.Nodes[0])
	}
	return -1
}

// IndexOfNode() returns the position of the specified node within the Selection
// object, or -1 if not found.
func (this *Selection) IndexOfNode(node *html.Node) int {
	return indexInSlice(this.Nodes, node)
}

// IndexOfSelection() returns the position of the first node in the specified
// Selection object within this Selection object, or -1 if not found.
func (this *Selection) IndexOfSelection(s *Selection) int {
	if s != nil && len(s.Nodes) > 0 {
		return indexInSlice(this.Nodes, s.Nodes[0])
	}
	return -1
}
