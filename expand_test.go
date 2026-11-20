package goquery

import (
	"fmt"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestAdd(t *testing.T) {
	sel := Doc().Find("div.row-fluid").Add("a")
	assertLength(t, sel.Nodes, 19)
}

func TestAddInvalid(t *testing.T) {
	sel1 := Doc().Find("div.row-fluid")
	sel2 := sel1.Add("")
	assertLength(t, sel1.Nodes, 9)
	assertLength(t, sel2.Nodes, 9)
	if sel1 == sel2 {
		t.Errorf("selections should not be the same")
	}
}

func TestAddRollback(t *testing.T) {
	sel := Doc().Find(".pvk-content")
	sel2 := sel.Add("a").End()
	assertEqual(t, sel, sel2)
}

func TestAddSelection(t *testing.T) {
	sel := Doc().Find("div.row-fluid")
	sel2 := Doc().Find("a")
	sel = sel.AddSelection(sel2)
	assertLength(t, sel.Nodes, 19)
}

func TestAddSelectionNil(t *testing.T) {
	sel := Doc().Find("div.row-fluid")
	assertLength(t, sel.Nodes, 9)

	sel = sel.AddSelection(nil)
	assertLength(t, sel.Nodes, 9)
}

func TestAddSelectionRollback(t *testing.T) {
	sel := Doc().Find(".pvk-content")
	sel2 := sel.Find("a")
	sel2 = sel.AddSelection(sel2).End()
	assertEqual(t, sel, sel2)
}

func TestAddNodes(t *testing.T) {
	sel := Doc().Find("div.pvk-gutter")
	sel2 := Doc().Find(".pvk-content")
	sel = sel.AddNodes(sel2.Nodes...)
	assertLength(t, sel.Nodes, 9)
}

func TestAddNodesNone(t *testing.T) {
	sel := Doc().Find("div.pvk-gutter").AddNodes()
	assertLength(t, sel.Nodes, 6)
}

func TestAddNodesRollback(t *testing.T) {
	sel := Doc().Find(".pvk-content")
	sel2 := sel.Find("a")
	sel2 = sel.AddNodes(sel2.Nodes...).End()
	assertEqual(t, sel, sel2)
}

func TestAddNodesBig(t *testing.T) {
	doc := DocW()
	sel := doc.Find("li")
	assertLength(t, sel.Nodes, 373)
	sel2 := doc.Find("xyz")
	assertLength(t, sel2.Nodes, 0)

	nodes := sel.Nodes
	sel2 = sel2.AddNodes(nodes...)
	assertLength(t, sel2.Nodes, 373)
	nodes2 := append(nodes, nodes...)
	sel2 = sel2.End().AddNodes(nodes2...)
	assertLength(t, sel2.Nodes, 373)
	nodes3 := append(nodes2, nodes...)
	sel2 = sel2.End().AddNodes(nodes3...)
	assertLength(t, sel2.Nodes, 373)
}

func TestAddNodesPreservesRelatedSelections(t *testing.T) {
	for _, size := range []int{1, minNodesForSet} {
		t.Run(fmt.Sprint(size), func(t *testing.T) {
			doc, err := NewDocumentFromReader(strings.NewReader(strings.Repeat("<i>x</i>", size+3)))
			if err != nil {
				t.Fatal(err)
			}
			all := doc.Find("i")
			original := append([]*html.Node(nil), all.Nodes...)
			prefix := all.Slice(0, size)
			left := prefix.AddNodes(original[size+1], original[0], original[size+1])
			right := prefix.AddNodes(original[size+2])

			for i, node := range original {
				if all.Nodes[i] != node {
					t.Errorf("source node %d changed", i)
				}
			}
			for name, result := range map[string]*Selection{"left": left, "right": right} {
				assertLength(t, result.Nodes, size+1)
				for i := 0; i < size; i++ {
					if result.Nodes[i] != original[i] {
						t.Errorf("%s: original node %d changed", name, i)
					}
				}
				if result.End() != prefix {
					t.Errorf("%s: End did not restore the receiver", name)
				}
			}
			if left.Nodes[size] != original[size+1] || right.Nodes[size] != original[size+2] {
				t.Error("adding to the receiver changed an earlier result")
			}
		})
	}
}

func TestAddMethodsPreserveSourceSelection(t *testing.T) {
	for _, method := range []struct {
		name  string
		index int
		add   func(*Selection, *Selection) *Selection
		want  []int
	}{
		{"Add", 0, func(s, all *Selection) *Selection { return s.Add(".c") }, []int{0, 2}},
		{"AddMatcher", 0, func(s, all *Selection) *Selection { return s.AddMatcher(compileMatcher(".c")) }, []int{0, 2}},
		{"AddSelection", 0, func(s, all *Selection) *Selection { return s.AddSelection(all.Eq(2)) }, []int{0, 2}},
		{"Union", 0, func(s, all *Selection) *Selection { return s.Union(all.Eq(2)) }, []int{0, 2}},
		{"AddBack", 1, func(s, all *Selection) *Selection { return s.AddBack() }, []int{1, 0, 2, 3}},
		{"AndSelf", 1, func(s, all *Selection) *Selection { return s.AndSelf() }, []int{1, 0, 2, 3}},
		{"AddBackFiltered", 1, func(s, all *Selection) *Selection { return s.AddBackFiltered("i") }, []int{1, 0, 2, 3}},
		{"AddBackMatcher", 1, func(s, all *Selection) *Selection { return s.AddBackMatcher(compileMatcher("i")) }, []int{1, 0, 2, 3}},
	} {
		t.Run(method.name, func(t *testing.T) {
			doc, err := NewDocumentFromReader(strings.NewReader(`<i>A</i><i>B</i><i class="c">C</i><i>D</i>`))
			if err != nil {
				t.Fatal(err)
			}
			all := doc.Find("i")
			original := append([]*html.Node(nil), all.Nodes...)
			current := all.Eq(method.index)
			result := method.add(current, all)
			if result.Length() != len(method.want) {
				t.Errorf("got %d nodes, want %d", result.Length(), len(method.want))
			} else {
				for i, index := range method.want {
					if result.Nodes[i] != original[index] {
						t.Errorf("result node %d has the wrong identity", i)
					}
				}
			}
			for i, node := range original {
				if all.Nodes[i] != node {
					t.Errorf("source node %d changed", i)
				}
			}
			if result.End() != current {
				t.Error("End did not restore the receiver")
			}
		})
	}
}

func TestAndSelf(t *testing.T) {
	sel := Doc().Find(".span12").Last().AndSelf()
	assertLength(t, sel.Nodes, 2)
}

func TestAndSelfRollback(t *testing.T) {
	sel := Doc().Find(".pvk-content")
	sel2 := sel.Find("a").AndSelf().End().End()
	assertEqual(t, sel, sel2)
}

func TestAddBack(t *testing.T) {
	sel := Doc().Find(".span12").Last().AddBack()
	assertLength(t, sel.Nodes, 2)
}

func TestAddBackRollback(t *testing.T) {
	sel := Doc().Find(".pvk-content")
	sel2 := sel.Find("a").AddBack().End().End()
	assertEqual(t, sel, sel2)
}

func TestAddBackFiltered(t *testing.T) {
	sel := Doc().Find(".span12, .footer").Find("h1").AddBackFiltered(".footer")
	assertLength(t, sel.Nodes, 2)
}

func TestAddBackFilteredRollback(t *testing.T) {
	sel := Doc().Find(".span12, .footer")
	sel2 := sel.Find("h1").AddBackFiltered(".footer").End().End()
	assertEqual(t, sel, sel2)
}
