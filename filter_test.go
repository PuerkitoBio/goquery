package goquery

import (
	"testing"
)

func TestFilter(t *testing.T) {
	sel := Doc().Find(".span12").Filter(".alert")
	AssertLength(t, sel.Nodes, 1)
}

func TestFilterNone(t *testing.T) {
	sel := Doc().Find(".span12").Filter(".zzalert")
	AssertLength(t, sel.Nodes, 0)
}

func TestFilterRollback(t *testing.T) {
	sel := Doc().Find(".pvk-content")
	sel2 := sel.Filter(".alert").End()
	AssertEqual(t, sel, sel2)
}

func TestFilterFunction(t *testing.T) {
	sel := Doc().Find(".pvk-content").FilterFunction(func(i int, s *Selection) bool {
		return i > 0
	})
	AssertLength(t, sel.Nodes, 2)
}

func TestFilterFunctionRollback(t *testing.T) {
	sel := Doc().Find(".pvk-content")
	sel2 := sel.FilterFunction(func(i int, s *Selection) bool {
		return i > 0
	}).End()
	AssertEqual(t, sel, sel2)
}

func TestFilterNode(t *testing.T) {
	sel := Doc().Find(".pvk-content")
	sel2 := sel.FilterNodes(sel.Nodes[2])
	AssertLength(t, sel2.Nodes, 1)
}

func TestFilterNodeRollback(t *testing.T) {
	sel := Doc().Find(".pvk-content")
	sel2 := sel.FilterNodes(sel.Nodes[2]).End()
	AssertEqual(t, sel, sel2)
}

func TestFilterSelection(t *testing.T) {
	sel := Doc().Find(".link")
	sel2 := Doc().Find("a[ng-click]")
	sel3 := sel.FilterSelection(sel2)
	AssertLength(t, sel3.Nodes, 1)
}

func TestFilterSelectionRollback(t *testing.T) {
	sel := Doc().Find(".link")
	sel2 := Doc().Find("a[ng-click]")
	sel2 = sel.FilterSelection(sel2).End()
	AssertEqual(t, sel, sel2)
}

func TestFilterSelectionNil(t *testing.T) {
	var sel2 *Selection

	sel := Doc().Find(".link")
	sel3 := sel.FilterSelection(sel2)
	AssertLength(t, sel3.Nodes, 0)
}

func TestNot(t *testing.T) {
	sel := Doc().Find(".span12").Not(".alert")
	AssertLength(t, sel.Nodes, 1)
}

func TestNotRollback(t *testing.T) {
	sel := Doc().Find(".span12")
	sel2 := sel.Not(".alert").End()
	AssertEqual(t, sel, sel2)
}

func TestNotNone(t *testing.T) {
	sel := Doc().Find(".span12").Not(".zzalert")
	AssertLength(t, sel.Nodes, 2)
}

func TestNotFunction(t *testing.T) {
	sel := Doc().Find(".pvk-content").NotFunction(func(i int, s *Selection) bool {
		return i > 0
	})
	AssertLength(t, sel.Nodes, 1)
}

func TestNotFunctionRollback(t *testing.T) {
	sel := Doc().Find(".pvk-content")
	sel2 := sel.NotFunction(func(i int, s *Selection) bool {
		return i > 0
	}).End()
	AssertEqual(t, sel, sel2)
}

func TestNotNode(t *testing.T) {
	sel := Doc().Find(".pvk-content")
	sel2 := sel.NotNodes(sel.Nodes[2])
	AssertLength(t, sel2.Nodes, 2)
}

func TestNotNodeRollback(t *testing.T) {
	sel := Doc().Find(".pvk-content")
	sel2 := sel.NotNodes(sel.Nodes[2]).End()
	AssertEqual(t, sel, sel2)
}

func TestNotSelection(t *testing.T) {
	sel := Doc().Find(".link")
	sel2 := Doc().Find("a[ng-click]")
	sel3 := sel.NotSelection(sel2)
	AssertLength(t, sel3.Nodes, 6)
}

func TestNotSelectionRollback(t *testing.T) {
	sel := Doc().Find(".link")
	sel2 := Doc().Find("a[ng-click]")
	sel2 = sel.NotSelection(sel2).End()
	AssertEqual(t, sel, sel2)
}

func TestIntersection(t *testing.T) {
	sel := Doc().Find(".pvk-gutter")
	sel2 := Doc().Find("div").Intersection(sel)
	AssertLength(t, sel2.Nodes, 6)
}

func TestIntersectionRollback(t *testing.T) {
	sel := Doc().Find(".pvk-gutter")
	sel2 := Doc().Find("div")
	sel2 = sel.Intersection(sel2).End()
	AssertEqual(t, sel, sel2)
}

func TestHas(t *testing.T) {
	sel := Doc().Find(".container-fluid").Has(".center-content")
	AssertLength(t, sel.Nodes, 2)
	// Has() returns the high-level .container-fluid div, and the one that is the immediate parent of center-content
}

func TestHasRollback(t *testing.T) {
	sel := Doc().Find(".container-fluid")
	sel2 := sel.Has(".center-content").End()
	AssertEqual(t, sel, sel2)
}

func TestHasNodes(t *testing.T) {
	sel := Doc().Find(".container-fluid")
	sel2 := Doc().Find(".center-content")
	sel = sel.HasNodes(sel2.Nodes...)
	AssertLength(t, sel.Nodes, 2)
	// Has() returns the high-level .container-fluid div, and the one that is the immediate parent of center-content
}

func TestHasNodesRollback(t *testing.T) {
	sel := Doc().Find(".container-fluid")
	sel2 := Doc().Find(".center-content")
	sel2 = sel.HasNodes(sel2.Nodes...).End()
	AssertEqual(t, sel, sel2)
}

func TestHasSelection(t *testing.T) {
	sel := Doc().Find("p")
	sel2 := Doc().Find("small")
	sel = sel.HasSelection(sel2)
	AssertLength(t, sel.Nodes, 1)
}

func TestHasSelectionRollback(t *testing.T) {
	sel := Doc().Find("p")
	sel2 := Doc().Find("small")
	sel2 = sel.HasSelection(sel2).End()
	AssertEqual(t, sel, sel2)
}

func TestEnd(t *testing.T) {
	sel := Doc().Find("p").Has("small").End()
	AssertLength(t, sel.Nodes, 4)
}

func TestEndToTop(t *testing.T) {
	sel := Doc().Find("p").Has("small").End().End().End()
	AssertLength(t, sel.Nodes, 0)
}
