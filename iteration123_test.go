//go:build go1.23
// +build go1.23

package goquery

import "testing"

func TestEachIter123(t *testing.T) {
	var cnt int

	sel := Doc().Find(".hero-unit .row-fluid")

	for i, s := range sel.EachIter() {
		cnt++
		t.Logf("At index %v, node %v", i, s.Nodes[0].Data)
	}

	sel = sel.Find("a")

	if cnt != 4 {
		t.Errorf("Expected EachIter() to call function 4 times, got %v times.", cnt)
	}
	assertLength(t, sel.Nodes, 6)
}

func TestEachIterWithBreak123(t *testing.T) {
	var cnt int

	sel := Doc().Find(".hero-unit .row-fluid")
	for i, s := range sel.EachIter() {
		cnt++
		t.Logf("At index %v, node %v", i, s.Nodes[0].Data)
		break
	}

	sel = sel.Find("a")

	if cnt != 1 {
		t.Errorf("Expected EachIter() to call function 1 time, got %v times.", cnt)
	}
	assertLength(t, sel.Nodes, 6)
}
