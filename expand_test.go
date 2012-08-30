package goquery

import (
	"testing"
)

func TestAdd(t *testing.T) {
	var cnt int

	EnsureDocLoaded()

	sel := Doc().Find("div.row-fluid")
	cnt = len(sel.Nodes)
	sel2 := sel.Add("a")
	if sel != sel2 {
		t.Error("Expected returned Selection from Add() to be same as 'this'.")
	}
	if len(sel.Nodes) <= cnt {
		t.Error("Expected nodes to be added to Selection.")
	}
	t.Logf("%v nodes after div.row-fluid and a were added.", len(sel.Nodes))
}
