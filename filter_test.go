package goquery

import (
	"testing"
)

func TestFilter(t *testing.T) {
	sel := Doc().Find(".span12").Filter(".alert")
	if len(sel.Nodes) != 1 {
		t.Errorf("Expected 1 node, found %v.", len(sel.Nodes))
	}
}

func TestFilterNone(t *testing.T) {
	sel := Doc().Find(".span12").Filter(".zzalert")
	if sel.Nodes != nil {
		t.Error("Expected no node (nil), found some.")
	}
}

func TestFilterFunction(t *testing.T) {
	sel := Doc().Find(".pvk-content").FilterFunction(func(i int, s *Selection) bool {
		return i > 0
	})
	if len(sel.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, found %v.", len(sel.Nodes))
	}
}

func TestFilterNode(t *testing.T) {
	sel := Doc().Find(".pvk-content")
	sel2 := sel.FilterNodes(sel.Nodes[2])
	if len(sel2.Nodes) != 1 {
		t.Errorf("Expected 1 node, found %v.", len(sel2.Nodes))
	}
}

func TestFilterSelection(t *testing.T) {
	sel := Doc().Find(".link")
	sel2 := Doc().Find("a[ng-click]")
	sel3 := sel.FilterSelection(sel2)
	if len(sel3.Nodes) != 1 {
		t.Errorf("Expected 1 node, found %v.", len(sel3.Nodes))
	}
}

func TestNot(t *testing.T) {
	sel := Doc().Find(".span12").Not(".alert")
	if len(sel.Nodes) != 1 {
		t.Errorf("Expected 1 node, found %v.", len(sel.Nodes))
	}
}

func TestNotNone(t *testing.T) {
	sel := Doc().Find(".span12").Not(".zzalert")
	if len(sel.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, found %v.", len(sel.Nodes))
	}
}

func TestNotFunction(t *testing.T) {
	sel := Doc().Find(".pvk-content").NotFunction(func(i int, s *Selection) bool {
		return i > 0
	})
	if len(sel.Nodes) != 1 {
		t.Errorf("Expected 1 nodes, found %v.", len(sel.Nodes))
	}
}

func TestNotNode(t *testing.T) {
	sel := Doc().Find(".pvk-content")
	sel2 := sel.NotNodes(sel.Nodes[2])
	if len(sel2.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, found %v.", len(sel2.Nodes))
	}
}

func TestNotSelection(t *testing.T) {
	sel := Doc().Find(".link")
	sel2 := Doc().Find("a[ng-click]")
	sel3 := sel.NotSelection(sel2)
	if len(sel3.Nodes) != 6 {
		t.Errorf("Expected 6 nodes, found %v.", len(sel3.Nodes))
	}
}

func TestUnion(t *testing.T) {
	sel := Doc().Find(".pvk-gutter")
	sel2 := Doc().Find("div").Union(sel)
	if len(sel2.Nodes) != 6 {
		t.Errorf("Expected 6 nodes, found %v.", len(sel2.Nodes))
	}
}
