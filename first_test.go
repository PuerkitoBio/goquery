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
	sel := doc.Find(".pvk-zzcontentzz").First()
	if len(sel.Nodes) != 0 {
		t.Errorf("Expected 0 node, found %v.", len(sel.Nodes))
	}
}
