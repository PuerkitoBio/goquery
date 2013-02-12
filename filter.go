package goquery

import (
	"code.google.com/p/cascadia"
	"code.google.com/p/go.net/html"
)

// Filter() reduces the set of matched elements to those that match the selector string.
// It returns a new Selection object for this subset of matching elements.
func (this *Selection) Filter(selector string) *Selection {
	return pushStack(this, winnow(this, selector, true))
}

// Not() removes elements from the Selection that match the selector string.
// It returns a new Selection object with the matching elements removed.
func (this *Selection) Not(selector string) *Selection {
	return pushStack(this, winnow(this, selector, false))
}

// FilterFunction() reduces the set of matched elements to those that pass the function's test.
// It returns a new Selection object for this subset of elements.
func (this *Selection) FilterFunction(f func(int, *Selection) bool) *Selection {
	return pushStack(this, winnowFunction(this, f, true))
}

// Not() removes elements from the Selection that pass the function's test.
// It returns a new Selection object with the matching elements removed.
func (this *Selection) NotFunction(f func(int, *Selection) bool) *Selection {
	return pushStack(this, winnowFunction(this, f, false))
}

// FilterNodes() reduces the set of matched elements to those that match the specified nodes.
// It returns a new Selection object for this subset of elements.
func (this *Selection) FilterNodes(nodes ...*html.Node) *Selection {
	return pushStack(this, winnowNodes(this, nodes, true))
}

// Not() removes elements from the Selection that match the specified nodes.
// It returns a new Selection object with the matching elements removed.
func (this *Selection) NotNodes(nodes ...*html.Node) *Selection {
	return pushStack(this, winnowNodes(this, nodes, false))
}

// FilterSelection() reduces the set of matched elements to those that match a
// node in the specified Selection object.
// It returns a new Selection object for this subset of elements.
func (this *Selection) FilterSelection(s *Selection) *Selection {
	if s == nil {
		return pushStack(this, winnowNodes(this, nil, true))
	}
	return pushStack(this, winnowNodes(this, s.Nodes, true))
}

// Not() removes elements from the Selection that match a node in the specified
// Selection object.
// It returns a new Selection object with the matching elements removed.
func (this *Selection) NotSelection(s *Selection) *Selection {
	if s == nil {
		return pushStack(this, winnowNodes(this, nil, false))
	}
	return pushStack(this, winnowNodes(this, s.Nodes, false))
}

// Intersection() is an alias for FilterSelection().
func (this *Selection) Intersection(s *Selection) *Selection {
	return this.FilterSelection(s)
}

// Has() reduces the set of matched elements to those that have a descendant
// that matches the selector.
// It returns a new Selection object with the matching elements.
func (this *Selection) Has(selector string) *Selection {
	return this.HasSelection(this.document.Find(selector))
}

// HasNodes() reduces the set of matched elements to those that have a
// descendant that matches one of the nodes.
// It returns a new Selection object with the matching elements.
func (this *Selection) HasNodes(nodes ...*html.Node) *Selection {
	return this.FilterFunction(func(_ int, s *Selection) bool {
		// Add all nodes that contain one of the specified nodes
		for _, n := range nodes {
			if s.Contains(n) {
				return true
			}
		}
		return false
	})
}

// HasSelection() reduces the set of matched elements to those that have a
// descendant that matches one of the nodes of the specified Selection object.
// It returns a new Selection object with the matching elements.
func (this *Selection) HasSelection(sel *Selection) *Selection {
	if sel == nil {
		return this.HasNodes()
	}
	return this.HasNodes(sel.Nodes...)
}

// End() ends the most recent filtering operation in the current chain and
// returns the set of matched elements to its previous state.
func (this *Selection) End() *Selection {
	if this.prevSel != nil {
		return this.prevSel
	}
	return newEmptySelection(this.document)
}

// Filter based on a selector string, and the indicator to keep (Filter) or
// to get rid of (Not) the matching elements.
func winnow(sel *Selection, selector string, keep bool) []*html.Node {
	cs := cascadia.MustCompile(selector)

	// Optimize if keep is requested
	if keep {
		return cs.Filter(sel.Nodes)
	} else {
		// Use grep
		return grep(sel, func(i int, s *Selection) bool {
			return !cs.Match(s.Get(0))
		})
	}
	return nil
}

// Filter based on an array of nodes, and the indicator to keep (Filter) or
// to get rid of (Not) the matching elements.
func winnowNodes(sel *Selection, nodes []*html.Node, keep bool) []*html.Node {
	return grep(sel, func(i int, s *Selection) bool {
		return isInSlice(nodes, s.Get(0)) == keep
	})
}

// Filter based on a function test, and the indicator to keep (Filter) or
// to get rid of (Not) the matching elements.
func winnowFunction(sel *Selection, f func(int, *Selection) bool, keep bool) []*html.Node {
	return grep(sel, func(i int, s *Selection) bool {
		return f(i, s) == keep
	})
}
