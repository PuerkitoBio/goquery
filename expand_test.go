package goquery

import (
	"testing"
)

func TestAdd(t *testing.T) {
	sel := Doc().Root.Find("div.row-fluid").Add("a")
	if len(sel.Nodes) != 19 {
		t.Errorf("Expected 19 nodes, found %v.", len(sel.Nodes))
	}
}

func TestAddSelection(t *testing.T) {
	sel := Doc().Root.Find("div.row-fluid")
	sel2 := Doc().Root.Find("a")
	sel = sel.AddSelection(sel2)
	if len(sel.Nodes) != 19 {
		t.Errorf("Expected 19 nodes, found %v.", len(sel.Nodes))
	}
}

func TestAddSelectionNil(t *testing.T) {
	sel := Doc().Root.Find("div.row-fluid")
	if len(sel.Nodes) != 9 {
		t.Errorf("Expected div.row-fluid to have 9 nodes, found %v.",
			len(sel.Nodes))
	}
	sel = sel.AddSelection(nil)
	if len(sel.Nodes) != 9 {
		t.Errorf("Expected add nil to keep it to 9 nodes, found %v.",
			len(sel.Nodes))
	}
}

func TestAddNodes(t *testing.T) {
	sel := Doc().Root.Find("div.pvk-gutter")
	sel2 := Doc().Root.Find(".pvk-content")
	sel = sel.AddNodes(sel2.Nodes...)
	if len(sel.Nodes) != 9 {
		t.Errorf("Expected 9 nodes, found %v.", len(sel.Nodes))
	}
}

func TestAddNodesNone(t *testing.T) {
	sel := Doc().Root.Find("div.pvk-gutter").AddNodes()
	if len(sel.Nodes) != 6 {
		t.Errorf("Expected 6 nodes, found %v.", len(sel.Nodes))
	}
}

func TestAndSelf(t *testing.T) {
	sel := Doc().Root.Find(".span12").Last().AndSelf()
	if len(sel.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, found %v.", len(sel.Nodes))
	}
}
