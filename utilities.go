package goquery

import (
	"exp/html"
)

// Contains() returns true if the specified Node is within,
// at any depth, one of the nodes in the Selection object.
// It is NOT inclusive, to behave like jQuery's implementation, and
// unlike Javascript's .contains(), so if the contained
// node is itself in the selection, it returns false.
func (this *Selection) Contains(n *html.Node) bool {
	return sliceContains(this.Nodes, n)
}

// Contains() returns true if the specified Node is within,
// at any depth, the root node of the Document object.
// It is NOT inclusive, to behave like jQuery's implementation, and
// unlike Javascript's .contains(), so if the contained
// node is itself in the selection, it returns false.
func (this *Document) Contains(n *html.Node) bool {
	return sliceContains([]*html.Node{this.Root}, n)
}

// Loop through all container nodes to search for the target node.
func sliceContains(container []*html.Node, contained *html.Node) bool {
	for _, n := range container {
		if nodeContains(n, contained) {
			return true
		}
	}

	return false
}

// Checks if the contained node is within the container node.
func nodeContains(container *html.Node, contained *html.Node) bool {
	// Check if the parent of the contained node is the container node, traversing
	// upward until the top is reached, or the container is found.
	for contained = contained.Parent; contained != nil; contained = contained.Parent {
		if container == contained {
			return true
		}
	}
	return false
}

// Checks if the target node is in the slice of nodes.
func isInSlice(slice []*html.Node, node *html.Node) bool {
	for _, n := range slice {
		if n == node {
			return true
		}
	}
	return false
}

// Appends the new nodes to the target slice, making sure no duplicate is added.
// There is no check to the original state of the target slice, so it may still contain
// duplicates. The target slice is returned because append() may create a new underlying array.
func appendWithoutDuplicates(target []*html.Node, nodes []*html.Node) []*html.Node {
	for _, n := range nodes {
		if !isInSlice(target, n) {
			target = append(target, n)
		}
	}

	return target
}

// Loop through a selection, returning only those nodes that pass the predicate function.
func grep(sel *Selection, predicate func(i int, s *Selection) bool) (result []*html.Node) {
	for i, n := range sel.Nodes {
		if predicate(i, newSingleSelection(n, sel.document)) {
			result = append(result, n)
		}
	}
	return
}

func pushStack(fromSel *Selection, nodes []*html.Node) (result *Selection) {
	result = &Selection{nodes, fromSel.document, fromSel}
	return
}
