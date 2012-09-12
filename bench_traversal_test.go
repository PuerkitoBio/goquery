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

func BenchmarkParents(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find("th a")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.Parents().Length()
		} else {
			sel.Parents()
		}
	}
	b.Logf("Parents=%d", n)
}

func BenchmarkParentsFiltered(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find("th a")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.ParentsFiltered("tr").Length()
		} else {
			sel.ParentsFiltered("tr")
		}
	}
	b.Logf("ParentsFiltered=%d", n)
}

func BenchmarkParentsUntil(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find("th a")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.ParentsUntil("table").Length()
		} else {
			sel.ParentsUntil("table")
		}
	}
	b.Logf("ParentsUntil=%d", n)
}

func BenchmarkParentsUntilSelection(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find("th a")
	sel2 := DocW().Root.Find("#content")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.ParentsUntilSelection(sel2).Length()
		} else {
			sel.ParentsUntilSelection(sel2)
		}
	}
	b.Logf("ParentsUntilSelection=%d", n)
}

func BenchmarkParentsUntilNodes(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find("th a")
	sel2 := DocW().Root.Find("#content")
	nodes := sel2.Nodes
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.ParentsUntilNodes(nodes...).Length()
		} else {
			sel.ParentsUntilNodes(nodes...)
		}
	}
	b.Logf("ParentsUntilNodes=%d", n)
}

func BenchmarkParentsFilteredUntil(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find(".toclevel-1 a")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.ParentsFilteredUntil(":nth-child(1)", "ul").Length()
		} else {
			sel.ParentsFilteredUntil(":nth-child(1)", "ul")
		}
	}
	b.Logf("ParentsFilteredUntil=%d", n)
}

func BenchmarkParentsFilteredUntilSelection(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find(".toclevel-1 a")
	sel2 := DocW().Root.Find("ul")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.ParentsFilteredUntilSelection(":nth-child(1)", sel2).Length()
		} else {
			sel.ParentsFilteredUntilSelection(":nth-child(1)", sel2)
		}
	}
	b.Logf("ParentsFilteredUntilSelection=%d", n)
}

func BenchmarkParentsFilteredUntilNodes(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find(".toclevel-1 a")
	sel2 := DocW().Root.Find("ul")
	nodes := sel2.Nodes
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.ParentsFilteredUntilNodes(":nth-child(1)", nodes...).Length()
		} else {
			sel.ParentsFilteredUntilNodes(":nth-child(1)", nodes...)
		}
	}
	b.Logf("ParentsFilteredUntilNodes=%d", n)
}

func BenchmarkSiblings(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find("ul li:nth-child(1)")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.Siblings().Length()
		} else {
			sel.Siblings()
		}
	}
	b.Logf("Siblings=%d", n)
}

func BenchmarkSiblingsFiltered(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find("ul li:nth-child(1)")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.SiblingsFiltered("[class]").Length()
		} else {
			sel.SiblingsFiltered("[class]")
		}
	}
	b.Logf("SiblingsFiltered=%d", n)
}
