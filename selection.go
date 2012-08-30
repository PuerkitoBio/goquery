package goquery

import (
	"exp/html"
)

type Selection struct {
	Nodes    []*html.Node
	document *Document
}

func newEmptySelection(doc *Document) *Selection {
	return &Selection{nil, doc}
}

func newSingleSelection(node *html.Node, doc *Document) *Selection {
	return &Selection{[]*html.Node{node}, doc}
}
