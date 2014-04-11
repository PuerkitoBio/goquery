package goquery

import (
	"code.google.com/p/go.net/html"
)

func getChildren(n *html.Node) (result []*html.Node) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result = append(result, c)
	}
	return result
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
	return indexInSlice(slice, node) > -1
}

// Returns the index of the target node in the slice, or -1.
func indexInSlice(slice []*html.Node, node *html.Node) int {
	if node != nil {
		for i, n := range slice {
			if n == node {
				return i
			}
		}
	}
	return -1
}

// Appends the new nodes to the target slice, making sure no duplicate is added.
// There is no check to the original state of the target slice, so it may still
// contain duplicates. The target slice is returned because append() may create
// a new underlying array.
func appendWithoutDuplicates(target []*html.Node, nodes []*html.Node) []*html.Node {
	for _, n := range nodes {
		if !isInSlice(target, n) {
			target = append(target, n)
		}
	}

	return target
}

// Loop through a selection, returning only those nodes that pass the predicate
// function.
func grep(sel *Selection, predicate func(i int, s *Selection) bool) (result []*html.Node) {
	for i, n := range sel.Nodes {
		if predicate(i, newSingleSelection(n, sel.document)) {
			result = append(result, n)
		}
	}
	return result
}

// Creates a new Selection object based on the specified nodes, and keeps the
// source Selection object on the stack (linked list).
func pushStack(fromSel *Selection, nodes []*html.Node) *Selection {
	result := &Selection{nodes, fromSel.document, fromSel}
	return result
}
