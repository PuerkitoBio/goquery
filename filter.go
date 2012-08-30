package goquery

import (
	"code.google.com/p/cascadia"
	"exp/html"
)

func (this *Selection) Filter(selector string) *Selection {
	sel, e := cascadia.Compile(selector)
	if e != nil {
		// Selector doesn't compile, which means empty selection
		return newEmptySelection(this.document)
	}

	return &Selection{sel.Filter(this.Nodes), this.document, nil}
}

func (this *Selection) FilterFunction(f func(int, *Selection) bool) *Selection {
	return &Selection{winnowFunction(this, f, true), this.document, nil}
}

func (this *Selection) NotFunction(f func(int, *Selection) bool) *Selection {
	return &Selection{winnowFunction(this, f, false), this.document, nil}
}

func (this *Selection) FilterNode(node *html.Node) *Selection {
	if isInSlice(this.Nodes, node) {
		return newSingleSelection(node, this.document)
	}
	return newEmptySelection(this.document)
}

func (this *Selection) Union(s *Selection) *Selection {
	return this.FilterSelection(s)
}

func (this *Selection) FilterSelection(s *Selection) *Selection {
	var matches []*html.Node

	if s == nil {
		return newEmptySelection(this.document)
	}

	// Check for a match for each current selection
	for _, n1 := range this.Nodes {
		for _, n2 := range s.Nodes {
			if n1 == n2 && !isInSlice(matches, n2) {
				matches = append(matches, n1)
				break
			}
		}
	}
	return &Selection{matches, this.document, nil}
}

func winnow(sel *Selection, selector string) []*html.Node {
	cs, e := cascadia.Compile(selector)
	if e != nil {
		// Selector doesn't compile, which means empty selection
		return nil
	}

	return cs.Filter(sel.Nodes)
}

// Identical functionality for FilterFunction() and NotFunction(), only keep changes.
func winnowFunction(sel *Selection, f func(int, *Selection) bool, keep bool) []*html.Node {
	return grep(sel, func(i int, s *Selection) bool {
		return f(i, s) == keep
	})
}

func (this *Selection) Has(selector string) *Selection {
	sel := this.document.Find(selector)

	return this.HasSelection(sel)
}

func (this *Selection) HasNode(node *html.Node) *Selection {
	return this.FilterFunction(func(_ int, s *Selection) bool {
		// Add all nodes that contain the node specified
		if s.Contains(node) {
			return true
		}
		return false
	})
}

func (this *Selection) HasSelection(sel *Selection) *Selection {
	return this.FilterFunction(func(_ int, s *Selection) bool {
		// Add all nodes that contain one of the nodes in the selection
		for _, n := range sel.Nodes {
			if s.Contains(n) {
				return true
			}
		}
		return false
	})
}
