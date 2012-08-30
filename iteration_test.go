package goquery

import (
	"testing"
)

func TestEach(t *testing.T) {
	var cnt int

	sel := Doc().Find(".hero-unit .row-fluid").Each(func(i int, n *Selection) {
		cnt++
		t.Logf("At index %v, node %v", i, n.Nodes[0].Data)
	}).Find("a")

	if cnt != 4 {
		t.Errorf("Expected Each() to call function 4 times, got %v times.", cnt)
	}
	if len(sel.Nodes) != 6 {
		t.Errorf("Expected 6 matching nodes, found %v.", len(sel.Nodes))
	}
}
