package goquery

import (
	"testing"
)

func TestFind(t *testing.T) {
	sel := Doc().Root.Find("div.row-fluid")
	AssertLength(t, sel.Nodes, 9)
}

func TestFindNotSelf(t *testing.T) {
	sel := Doc().Root.Find("h1").Find("h1")
	AssertLength(t, sel.Nodes, 0)
}

func TestFindInvalidSelector(t *testing.T) {
	defer AssertPanic(t)
	Doc().Root.Find(":+ ^")
}

func TestChainedFind(t *testing.T) {
	sel := Doc().Root.Find("div.hero-unit").Find(".row-fluid")
	AssertLength(t, sel.Nodes, 4)
}

func TestChildren(t *testing.T) {
	sel := Doc().Root.Find(".pvk-content").Children()
	AssertLength(t, sel.Nodes, 5)
}

func TestContents(t *testing.T) {
	sel := Doc().Root.Find(".pvk-content").Contents()
	AssertLength(t, sel.Nodes, 13)
}

func TestChildrenFiltered(t *testing.T) {
	sel := Doc().Root.Find(".pvk-content").ChildrenFiltered(".hero-unit")
	AssertLength(t, sel.Nodes, 1)
}

func TestContentsFiltered(t *testing.T) {
	sel := Doc().Root.Find(".pvk-content").ContentsFiltered(".hero-unit")
	AssertLength(t, sel.Nodes, 1)
}

func TestChildrenFilteredNone(t *testing.T) {
	sel := Doc().Root.Find(".pvk-content").ChildrenFiltered("a.btn")
	AssertLength(t, sel.Nodes, 0)
}

func TestParent(t *testing.T) {
	sel := Doc().Root.Find(".container-fluid").Parent()
	AssertLength(t, sel.Nodes, 3)
}

func TestParentBody(t *testing.T) {
	sel := Doc().Root.Find("body").Parent()
	AssertLength(t, sel.Nodes, 1)
}

func TestParentFiltered(t *testing.T) {
	sel := Doc().Root.Find(".container-fluid").ParentFiltered(".hero-unit")
	AssertLength(t, sel.Nodes, 1)
	AssertClass(t, sel, "hero-unit")
}

func TestParents(t *testing.T) {
	sel := Doc().Root.Find(".container-fluid").Parents()
	AssertLength(t, sel.Nodes, 8)
}

func TestParentsFiltered(t *testing.T) {
	sel := Doc().Root.Find(".container-fluid").ParentsFiltered("body")
	AssertLength(t, sel.Nodes, 1)
}

func TestParentsUntil(t *testing.T) {
	sel := Doc().Root.Find(".container-fluid").ParentsUntil("body")
	AssertLength(t, sel.Nodes, 6)
}

func TestParentsUntilSelection(t *testing.T) {
	sel := Doc().Root.Find(".container-fluid")
	sel2 := Doc().Root.Find(".pvk-content")
	sel = sel.ParentsUntilSelection(sel2)
	AssertLength(t, sel.Nodes, 3)
}

func TestParentsUntilNodes(t *testing.T) {
	sel := Doc().Root.Find(".container-fluid")
	sel2 := Doc().Root.Find(".pvk-content, .hero-unit")
	sel = sel.ParentsUntilNodes(sel2.Nodes...)
	AssertLength(t, sel.Nodes, 2)
}

func TestParentsFilteredUntil(t *testing.T) {
	sel := Doc().Root.Find(".container-fluid").ParentsFilteredUntil(".pvk-content", "body")
	AssertLength(t, sel.Nodes, 2)
}

func TestParentsFilteredUntilSelection(t *testing.T) {
	sel := Doc().Root.Find(".container-fluid")
	sel2 := Doc().Root.Find(".row-fluid")
	sel = sel.ParentsFilteredUntilSelection("div", sel2)
	AssertLength(t, sel.Nodes, 3)
}

func TestParentsFilteredUntilNodes(t *testing.T) {
	sel := Doc().Root.Find(".container-fluid")
	sel2 := Doc().Root.Find(".row-fluid")
	sel = sel.ParentsFilteredUntilNodes("body", sel2.Nodes...)
	AssertLength(t, sel.Nodes, 1)
}

func TestSiblings(t *testing.T) {
	sel := Doc().Root.Find("h1").Siblings()
	AssertLength(t, sel.Nodes, 1)
}

func TestSiblings2(t *testing.T) {
	sel := Doc().Root.Find(".pvk-gutter").Siblings()
	AssertLength(t, sel.Nodes, 9)
}

func TestSiblings3(t *testing.T) {
	sel := Doc().Root.Find("body>.container-fluid").Siblings()
	AssertLength(t, sel.Nodes, 0)
}

func TestSiblingsFiltered(t *testing.T) {
	sel := Doc().Root.Find(".pvk-gutter").SiblingsFiltered(".pvk-content")
	AssertLength(t, sel.Nodes, 3)
}

func TestNext(t *testing.T) {
	sel := Doc().Root.Find("h1").Next()
	AssertLength(t, sel.Nodes, 1)
}

func TestNext2(t *testing.T) {
	sel := Doc().Root.Find(".close").Next()
	AssertLength(t, sel.Nodes, 1)
}

func TestNextNone(t *testing.T) {
	sel := Doc().Root.Find("small").Next()
	AssertLength(t, sel.Nodes, 0)
}

func TestNextFiltered(t *testing.T) {
	sel := Doc().Root.Find(".container-fluid").NextFiltered("div")
	AssertLength(t, sel.Nodes, 2)
}

func TestNextFiltered2(t *testing.T) {
	sel := Doc().Root.Find(".container-fluid").NextFiltered("[ng-view]")
	AssertLength(t, sel.Nodes, 1)
}
