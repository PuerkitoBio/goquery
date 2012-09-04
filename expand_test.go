package goquery

import (
	"testing"
)

func TestAdd(t *testing.T) {
	sel := Doc().Root.Find("div.row-fluid").Add("a")
	AssertLength(t, sel.Nodes, 19)
}

func TestAddSelection(t *testing.T) {
	sel := Doc().Root.Find("div.row-fluid")
	sel2 := Doc().Root.Find("a")
	sel = sel.AddSelection(sel2)
	AssertLength(t, sel.Nodes, 19)
}

func TestAddSelectionNil(t *testing.T) {
	sel := Doc().Root.Find("div.row-fluid")
	AssertLength(t, sel.Nodes, 9)

	sel = sel.AddSelection(nil)
	AssertLength(t, sel.Nodes, 9)
}

func TestAddNodes(t *testing.T) {
	sel := Doc().Root.Find("div.pvk-gutter")
	sel2 := Doc().Root.Find(".pvk-content")
	sel = sel.AddNodes(sel2.Nodes...)
	AssertLength(t, sel.Nodes, 9)
}

func TestAddNodesNone(t *testing.T) {
	sel := Doc().Root.Find("div.pvk-gutter").AddNodes()
	AssertLength(t, sel.Nodes, 6)
}

func TestAndSelf(t *testing.T) {
	sel := Doc().Root.Find(".span12").Last().AndSelf()
	AssertLength(t, sel.Nodes, 2)
}
