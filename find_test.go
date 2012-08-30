package goquery

import (
	"testing"
)

func TestFind(t *testing.T) {
	sel := doc.Find("div.row-fluid")
	if sel.Nodes == nil || len(sel.Nodes) != 9 {
		t.Errorf("Expected 9 matching nodes, found %v.", len(sel.Nodes))
	}
}

func TestFindInvalidSelector(t *testing.T) {
	var cnt int

	sel := doc.Find(":+ ^")
	if sel.Nodes != nil {
		t.Error("Expected a Selection object with Nodes == nil.")
	}
	sel.Each(func(i int, n *Selection) {
		cnt++
	})
	if cnt > 0 {
		t.Error("Expected Each() to not be called on empty Selection.")
	}
	sel2 := sel.Find("div")
	if sel2.Nodes != nil {
		t.Error("Expected Find() on empty Selection to return an empty Selection.")
	}
}

func TestChainedFind(t *testing.T) {
	sel := doc.Find("div.hero-unit").Find(".row-fluid")
	if sel.Nodes == nil || len(sel.Nodes) != 4 {
		t.Errorf("Expected 4 matching nodes, found %v.", len(sel.Nodes))
	}
}
