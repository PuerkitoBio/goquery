package goquery

import (
  "testing"
)

func TestChildren(t *testing.T) {
  sel := doc.Find(".pvk-content").Children()
  if len(sel.Nodes) != 13 {
    t.Errorf("Expected 13 child nodes, got %v.", len(sel.Nodes))
    for _, n := range sel.Nodes {
      t.Logf("%+v", n)
    }
  }
}

func TestChildrenFiltered(t *testing.T) {
  sel := doc.Find(".pvk-content").ChildrenFiltered(".hero-unit")
  if len(sel.Nodes) != 1 {
    t.Errorf("Expected 1 child nodes, got %v.", len(sel.Nodes))
    for _, n := range sel.Nodes {
      t.Logf("%+v", n)
    }
  }
}

func TestChildrenFilteredNone(t *testing.T) {
  sel := doc.Find(".pvk-content").ChildrenFiltered("a.btn")
  if len(sel.Nodes) != 0 {
    t.Errorf("Expected 0 child node, got %v.", len(sel.Nodes))
    for _, n := range sel.Nodes {
      t.Logf("%+v", n)
    }
  }
}
