package goquery

import (
	"testing"
)

func TestFirst(t *testing.T) {
	sel := doc.Find(".pvk-content").First()
	if len(sel.Nodes) != 1 {
		t.Errorf("Expected 1 node, found %v.", len(sel.Nodes))
	}
}

func TestFirstEmpty(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Error("Expected a panic, First() called on empty Selection.")
		}
	}()
	doc.Find(".pvk-zzcontentzz").First()
}

func TestLast(t *testing.T) {
	sel := doc.Find(".pvk-content").Last()
	if len(sel.Nodes) != 1 {
		t.Errorf("Expected 1 node, found %v.", len(sel.Nodes))
	}
	// Should contain Footer
	foot := doc.Find(".footer")
	if !sel.Contains(foot.Nodes[0]) {
		t.Error("Last .pvk-content should contain .footer.")
	}
}

func TestEq(t *testing.T) {
	sel := doc.Find(".pvk-content").Eq(1)
	if len(sel.Nodes) != 1 {
		t.Errorf("Expected 1 node, found %v.", len(sel.Nodes))
	}
}

func TestEqNegative(t *testing.T) {
	sel := doc.Find(".pvk-content").Eq(-1)
	if len(sel.Nodes) != 1 {
		t.Errorf("Expected 1 node, found %v.", len(sel.Nodes))
	}
	// Should contain Footer
	foot := doc.Find(".footer")
	if !sel.Contains(foot.Nodes[0]) {
		t.Error("Index -1 of .pvk-content should contain .footer.")
	}
}

func TestSlice(t *testing.T) {
	sel := doc.Find(".pvk-content").Slice(0, 2)
	if len(sel.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, found %v.", len(sel.Nodes))
	}
}

func TestSliceOutOfBounds(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Error("Expected a panic, Slice() called with out of bounds indices.")
		}
	}()
	doc.Find(".pvk-content").Slice(2, 12)
}
