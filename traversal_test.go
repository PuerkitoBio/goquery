package goquery

import (
	"testing"
)

func TestFind(t *testing.T) {
	sel := Doc().Find("div.row-fluid")
	if sel.Nodes == nil || len(sel.Nodes) != 9 {
		t.Errorf("Expected 9 matching nodes, found %v.", len(sel.Nodes))
	}
}

func TestFindNotSelf(t *testing.T) {
	sel := Doc().Find("h1").Find("h1")
	if len(sel.Nodes) > 0 {
		t.Errorf("Expected no node, found %v.", len(sel.Nodes))
	}
}

func TestFindInvalidSelector(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Error("Expected panic due to invalid selector.")
		}
	}()

	Doc().Find(":+ ^")
}

func TestEachEmptySelection(t *testing.T) {
	var cnt int

	sel := Doc().Find("zzzz")
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
	sel := Doc().Find("div.hero-unit").Find(".row-fluid")
	if sel.Nodes == nil || len(sel.Nodes) != 4 {
		t.Errorf("Expected 4 matching nodes, found %v.", len(sel.Nodes))
	}
}

func TestChildren(t *testing.T) {
	sel := Doc().Find(".pvk-content").Children()
	if len(sel.Nodes) != 13 {
		t.Errorf("Expected 13 child nodes, got %v.", len(sel.Nodes))
		for _, n := range sel.Nodes {
			t.Logf("%+v", n)
		}
	}
}

func TestChildrenFiltered(t *testing.T) {
	sel := Doc().Find(".pvk-content").ChildrenFiltered(".hero-unit")
	if len(sel.Nodes) != 1 {
		t.Errorf("Expected 1 child nodes, got %v.", len(sel.Nodes))
		for _, n := range sel.Nodes {
			t.Logf("%+v", n)
		}
	}
}

func TestChildrenFilteredNone(t *testing.T) {
	sel := Doc().Find(".pvk-content").ChildrenFiltered("a.btn")
	if len(sel.Nodes) != 0 {
		t.Errorf("Expected 0 child node, got %v.", len(sel.Nodes))
		for _, n := range sel.Nodes {
			t.Logf("%+v", n)
		}
	}
}
