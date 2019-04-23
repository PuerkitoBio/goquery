package goquery

import "golang.org/x/net/html"

// Add adds the selector string's matching nodes to those in the current
// selection and returns a new Selection object.
// The selector string is run in the context of the document of the current
// Selection object.
// Do not assume that this method appends the elements to the existing collection
// in the order they are passed to the .Add() method.
// When all elements are members of the same document, the resulting collection
// from .Add() will be sorted in document order; that is, in order of each element's appearance in the document.
// If the collection consists of elements from different documents or ones not in any document, the sort order is undefined.
func (s *Selection) Add(selector string) *Selection {
	return s.AddNodes(findWithMatcher([]*html.Node{s.document.rootNode}, compileMatcher(selector))...)
}

// AddMatcher adds the matcher's matching nodes to those in the current
// selection and returns a new Selection object.
// The matcher is run in the context of the document of the current
// Selection object.
// Nodes ordering in the result collection follows the same rules as .Add().
func (s *Selection) AddMatcher(m Matcher) *Selection {
	return s.AddNodes(findWithMatcher([]*html.Node{s.document.rootNode}, m)...)
}

// AddSelection adds the specified Selection object's nodes to those in the
// current selection and returns a new Selection object.
// Nodes ordering in the result collection follows the same rules as .Add().
func (s *Selection) AddSelection(sel *Selection) *Selection {
	if sel == nil {
		return s.AddNodes()
	}
	return s.AddNodes(sel.Nodes...)
}

// Union is an alias for AddSelection.
func (s *Selection) Union(sel *Selection) *Selection {
	return s.AddSelection(sel)
}

// AddNodes adds the specified nodes to those in the
// current selection and returns a new Selection object.
func (s *Selection) AddNodes(nodes ...*html.Node) *Selection {
	return pushStack(s, sortWithoutDuplicates(s.Nodes, nodes, nil))
}

// AndSelf adds the previous set of elements on the stack to the current set.
// It returns a new Selection object containing the current Selection combined
// with the previous one.
// Deprecated: This function has been deprecated and is now an alias for AddBack().
func (s *Selection) AndSelf() *Selection {
	return s.AddBack()
}

// AddBack adds the previous set of elements on the stack to the current set.
// It returns a new Selection object containing the current Selection combined
// with the previous one.
func (s *Selection) AddBack() *Selection {
	return s.AddSelection(s.prevSel)
}

// AddBackFiltered reduces the previous set of elements on the stack to those that
// match the selector string, and adds them to the current set.
// It returns a new Selection object containing the current Selection combined
// with the filtered previous one
func (s *Selection) AddBackFiltered(selector string) *Selection {
	return s.AddSelection(s.prevSel.Filter(selector))
}

// AddBackMatcher reduces the previous set of elements on the stack to those that match
// the mateher, and adds them to the curernt set.
// It returns a new Selection object containing the current Selection combined
// with the filtered previous one
func (s *Selection) AddBackMatcher(m Matcher) *Selection {
	return s.AddSelection(s.prevSel.FilterMatcher(m))
}
