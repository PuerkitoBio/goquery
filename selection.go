package goquery

import (
	"exp/html"
)

type Selection struct {
	Nodes    []*html.Node
	document *Document
	prevSel  *Selection
}

func (this *Selection) Size() int {
	return this.Length()
}

func (this *Selection) Length() int {
	return len(this.Nodes)
}

func newEmptySelection(doc *Document) *Selection {
	return &Selection{nil, doc, nil}
}

func newSingleSelection(node *html.Node, doc *Document) *Selection {
	return &Selection{[]*html.Node{node}, doc, nil}
}
