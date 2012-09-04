package goquery

import (
	"testing"
)

func TestFilter(t *testing.T) {
	sel := Doc().Root.Find(".span12").Filter(".alert")
	AssertLength(t, sel.Nodes, 1)
}

func TestFilterNone(t *testing.T) {
	sel := Doc().Root.Find(".span12").Filter(".zzalert")
	AssertLength(t, sel.Nodes, 0)
}

func TestFilterFunction(t *testing.T) {
	sel := Doc().Root.Find(".pvk-content").FilterFunction(func(i int, s *Selection) bool {
		return i > 0
	})
	AssertLength(t, sel.Nodes, 2)
}

func TestFilterNode(t *testing.T) {
	sel := Doc().Root.Find(".pvk-content")
	sel2 := sel.FilterNodes(sel.Nodes[2])
	AssertLength(t, sel2.Nodes, 1)
}

func TestFilterSelection(t *testing.T) {
	sel := Doc().Root.Find(".link")
	sel2 := Doc().Root.Find("a[ng-click]")
	sel3 := sel.FilterSelection(sel2)
	AssertLength(t, sel3.Nodes, 1)
}

func TestFilterSelectionNil(t *testing.T) {
	var sel2 *Selection

	sel := Doc().Root.Find(".link")
	sel3 := sel.FilterSelection(sel2)
	AssertLength(t, sel3.Nodes, 0)
}

func TestNot(t *testing.T) {
	sel := Doc().Root.Find(".span12").Not(".alert")
	AssertLength(t, sel.Nodes, 1)
}

func TestNotNone(t *testing.T) {
	sel := Doc().Root.Find(".span12").Not(".zzalert")
	AssertLength(t, sel.Nodes, 2)
}

func TestNotFunction(t *testing.T) {
	sel := Doc().Root.Find(".pvk-content").NotFunction(func(i int, s *Selection) bool {
		return i > 0
	})
	AssertLength(t, sel.Nodes, 1)
}

func TestNotNode(t *testing.T) {
	sel := Doc().Root.Find(".pvk-content")
	sel2 := sel.NotNodes(sel.Nodes[2])
	AssertLength(t, sel2.Nodes, 2)
}

func TestNotSelection(t *testing.T) {
	sel := Doc().Root.Find(".link")
	sel2 := Doc().Root.Find("a[ng-click]")
	sel3 := sel.NotSelection(sel2)
	AssertLength(t, sel3.Nodes, 6)
}

func TestIntersection(t *testing.T) {
	sel := Doc().Root.Find(".pvk-gutter")
	sel2 := Doc().Root.Find("div").Intersection(sel)
	AssertLength(t, sel2.Nodes, 6)
}

func TestHas(t *testing.T) {
	sel := Doc().Root.Find(".container-fluid").Has(".center-content")
	AssertLength(t, sel.Nodes, 2)
	// Has() returns the high-level .container-fluid div, and the one that is the immediate parent of center-content
}

func TestHasNodes(t *testing.T) {
	sel := Doc().Root.Find(".container-fluid")
	sel2 := Doc().Root.Find(".center-content")
	sel = sel.HasNodes(sel2.Nodes...)
	AssertLength(t, sel.Nodes, 2)
	// Has() returns the high-level .container-fluid div, and the one that is the immediate parent of center-content
}

func TestHasSelection(t *testing.T) {
	sel := Doc().Root.Find("p")
	sel2 := Doc().Root.Find("small")
	sel = sel.HasSelection(sel2)
	AssertLength(t, sel.Nodes, 1)
}

func TestEnd(t *testing.T) {
	sel := Doc().Root.Find("p").Has("small").End()
	AssertLength(t, sel.Nodes, 4)
}

func TestEndToTop(t *testing.T) {
	sel := Doc().Root.Find("p").Has("small").End().End().End()
	AssertLength(t, sel.Nodes, 0)
}
