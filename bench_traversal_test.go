package goquery

import (
	"testing"
)

func BenchmarkFind(b *testing.B) {
	var n int

	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = DocB().Root.Find("dd").Length()

		} else {
			DocB().Root.Find("dd")
		}
	}
	b.Logf("Find=%d", n)
}

func BenchmarkFindWithinSelection(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find("ul")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.Find("a[class]").Length()
		} else {
			sel.Find("a[class]")
		}
	}
	b.Logf("FindWithinSelection=%d", n)
}

func BenchmarkFindSelection(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find("ul")
	sel2 := DocW().Root.Find("span")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.FindSelection(sel2).Length()
		} else {
			sel.FindSelection(sel2)
		}
	}
	b.Logf("FindSelection=%d", n)
}

func BenchmarkFindNodes(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find("ul")
	sel2 := DocW().Root.Find("span")
	nodes := sel2.Nodes
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.FindNodes(nodes...).Length()
		} else {
			sel.FindNodes(nodes...)
		}
	}
	b.Logf("FindNodes=%d", n)
}

func BenchmarkContents(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find(".toclevel-1")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.Contents().Length()
		} else {
			sel.Contents()
		}
	}
	b.Logf("Contents=%d", n)
}

func BenchmarkContentsFiltered(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find(".toclevel-1")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.ContentsFiltered("a[href=\"#Examples\"]").Length()
		} else {
			sel.ContentsFiltered("a[href=\"#Examples\"]")
		}
	}
	b.Logf("ContentsFiltered=%d", n)
}

func BenchmarkChildren(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find(".toclevel-2")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.Children().Length()
		} else {
			sel.Children()
		}
	}
	b.Logf("Children=%d", n)
}

func BenchmarkChildrenFiltered(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find("h3")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.ChildrenFiltered(".editsection").Length()
		} else {
			sel.ChildrenFiltered(".editsection")
		}
	}
	b.Logf("ChildrenFiltered=%d", n)
}

func BenchmarkParent(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find("li")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.Parent().Length()
		} else {
			sel.Parent()
		}
	}
	b.Logf("Parent=%d", n)
}

func BenchmarkParentFiltered(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find("li")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.ParentFiltered("ul[id]").Length()
		} else {
			sel.ParentFiltered("ul[id]")
		}
	}
	b.Logf("ParentFiltered=%d", n)
}
