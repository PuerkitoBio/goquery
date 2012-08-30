package goquery

import (
	"code.google.com/p/cascadia"
	"exp/html"
)

func (this *Selection) Filter(selector string) *Selection {
	sel, e := cascadia.Compile(selector)
	if e != nil {
		// Selector doesn't compile, which means empty selection
		return &Selection{nil, this.document}
	}

	return &Selection{sel.Filter(this.Nodes), this.document}
}

func (this *Selection) FilterFunction(f func(int, *Selection) bool) *Selection {
	var matches []*html.Node

	// Check for a match for each current selection
	for i, n := range this.Nodes {
		if f(i, &Selection{[]*html.Node{n}, this.document}) {
			matches = append(matches, n)
		}
	}
	return &Selection{matches, this.document}
}

func (this *Selection) FilterNode(node *html.Node) *Selection {
	// TODO : Use Contains() on the this.Nodes array, if it contains, return node Selection, otherwise empty
	for _, n := range this.Nodes {
		if n == node {
			return &Selection{[]*html.Node{n}, this.document}
		}
	}
	return &Selection{nil, this.document}
}

func (this *Selection) FilterSelection(s *Selection) *Selection {
	// TODO : Exactly an Union() of two Selections, so maybe call it Union(), or have it as synonymous
	var matches []*html.Node

	if s == nil {
		return &Selection{nil, this.document}
	}

	// Check for a match for each current selection
	for _, n1 := range this.Nodes {
		for _, n2 := range s.Nodes {
			if n1 == n2 {
				matches = append(matches, n1)
				break
			}
		}
	}
	return &Selection{matches, this.document}
}
