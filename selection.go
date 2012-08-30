package goquery

import (
	"exp/html"
)

type Selection struct {
	Nodes    []*html.Node
	document *Document
}
