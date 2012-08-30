package goquery

import (
	"code.google.com/p/cascadia"
	"exp/html"
)

func (this *Selection) Filter(selector string) *Selection {
	return pushStack(this, winnow(this, selector, true))
}

func (this *Selection) Not(selector string) *Selection {
	return pushStack(this, winnow(this, selector, false))
}

func (this *Selection) FilterFunction(f func(int, *Selection) bool) *Selection {
	return pushStack(this, winnowFunction(this, f, true))
}

func (this *Selection) NotFunction(f func(int, *Selection) bool) *Selection {
	return pushStack(this, winnowFunction(this, f, false))
}

func (this *Selection) FilterNodes(nodes ...*html.Node) *Selection {
	return pushStack(this, winnowNodes(this, nodes, true))
}

func (this *Selection) NotNodes(nodes ...*html.Node) *Selection {
	return pushStack(this, winnowNodes(this, nodes, false))
}

func (this *Selection) FilterSelection(s *Selection) *Selection {
	return pushStack(this, winnowNodes(this, s.Nodes, true))
}

func (this *Selection) NotSelection(s *Selection) *Selection {
	return pushStack(this, winnowNodes(this, s.Nodes, false))
}

func (this *Selection) Union(s *Selection) *Selection {
	return this.FilterSelection(s)
}

func winnow(sel *Selection, selector string, keep bool) []*html.Node {
	cs := cascadia.MustCompile(selector)

	// Optimize if keep is requested
	if keep {
		return cs.Filter(sel.Nodes)
	} else {
		// Use grep
		return grep(sel, func(i int, s *Selection) bool {
			return !cs(s.Get(0))
		})
	}
	return nil
}

func winnowNodes(sel *Selection, nodes []*html.Node, keep bool) []*html.Node {
	return grep(sel, func(i int, s *Selection) bool {
		return isInSlice(nodes, s.Get(0)) == keep
	})
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
