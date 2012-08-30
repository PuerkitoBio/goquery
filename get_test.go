package goquery

import (
	"testing"
)

func TestGet(t *testing.T) {
	sel := doc.Find(".pvk-content")
	node := sel.Get(1)
	if sel.Nodes[1] != node {
		t.Errorf("Expected node %v to be %v.", node, sel.Nodes[1])
	}
	node = sel.Get(-3)
	if sel.Nodes[0] != node {
		t.Errorf("Expected node %v to be %v.", node, sel.Nodes[0])
	}
	node = sel.Get(129)
	if node != nil {
		t.Error("Expected node to be nil.")
	}
}
