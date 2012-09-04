package goquery

import (
	"testing"
)

func TestFirst(t *testing.T) {
	sel := Doc().Root.Find(".pvk-content").First()
	AssertLength(t, sel.Nodes, 1)
}

func TestFirstEmpty(t *testing.T) {
	defer AssertPanic(t)
	Doc().Root.Find(".pvk-zzcontentzz").First()
}

func TestLast(t *testing.T) {
	sel := Doc().Root.Find(".pvk-content").Last()
	AssertLength(t, sel.Nodes, 1)

	// Should contain Footer
	foot := Doc().Root.Find(".footer")
	if !sel.Contains(foot.Nodes[0]) {
		t.Error("Last .pvk-content should contain .footer.")
	}
}

func TestEq(t *testing.T) {
	sel := Doc().Root.Find(".pvk-content").Eq(1)
	AssertLength(t, sel.Nodes, 1)
}

func TestEqNegative(t *testing.T) {
	sel := Doc().Root.Find(".pvk-content").Eq(-1)
	AssertLength(t, sel.Nodes, 1)

	// Should contain Footer
	foot := Doc().Root.Find(".footer")
	if !sel.Contains(foot.Nodes[0]) {
		t.Error("Index -1 of .pvk-content should contain .footer.")
	}
}

func TestSlice(t *testing.T) {
	sel := Doc().Root.Find(".pvk-content").Slice(0, 2)

	AssertLength(t, sel.Nodes, 2)
}

func TestSliceOutOfBounds(t *testing.T) {
	defer AssertPanic(t)
	Doc().Root.Find(".pvk-content").Slice(2, 12)
}

func TestGet(t *testing.T) {
	sel := Doc().Root.Find(".pvk-content")
	node := sel.Get(1)
	if sel.Nodes[1] != node {
		t.Errorf("Expected node %v to be %v.", node, sel.Nodes[1])
	}
}

func TestGetNegative(t *testing.T) {
	sel := Doc().Root.Find(".pvk-content")
	node := sel.Get(-3)
	if sel.Nodes[0] != node {
		t.Errorf("Expected node %v to be %v.", node, sel.Nodes[0])
	}
}

func TestGetInvalid(t *testing.T) {
	defer AssertPanic(t)
	sel := Doc().Root.Find(".pvk-content")
	sel.Get(129)
}

func TestIndex(t *testing.T) {
	sel := Doc().Root.Find(".pvk-content")
	if i := sel.Index(); i != 1 {
		t.Errorf("Expected index of 1, got %v.", i)
	}
}

func TestIndexSelector(t *testing.T) {
	sel := Doc().Root.Find(".hero-unit")
	if i := sel.IndexSelector("div"); i != 4 {
		t.Errorf("Expected index of 4, got %v.", i)
	}
}

func TestIndexOfNode(t *testing.T) {
	sel := Doc().Root.Find("div.pvk-gutter")
	if i := sel.IndexOfNode(sel.Nodes[1]); i != 1 {
		t.Errorf("Expected index of 1, got %v.", i)
	}
}

func TestIndexOfNilNode(t *testing.T) {
	sel := Doc().Root.Find("div.pvk-gutter")
	if i := sel.IndexOfNode(nil); i != -1 {
		t.Errorf("Expected index of -1, got %v.", i)
	}
}

func TestIndexOfSelection(t *testing.T) {
	sel := Doc().Root.Find("div")
	sel2 := Doc().Root.Find(".hero-unit")
	if i := sel.IndexOfSelection(sel2); i != 4 {
		t.Errorf("Expected index of 4, got %v.", i)
	}
}
