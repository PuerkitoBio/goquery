package goquery

import (
	"testing"
)

func BenchmarkFilter(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find("li")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.Filter(".toclevel-1").Length()
		} else {
			sel.Filter(".toclevel-1")
		}
	}
	b.Logf("Filter=%d", n)
}

func BenchmarkNot(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find("li")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.Not(".toclevel-2").Length()
		} else {
			sel.Filter(".toclevel-2")
		}
	}
	b.Logf("Not=%d", n)
}

func BenchmarkFilterFunction(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find("li")
	f := func(i int, s *Selection) bool {
		return len(s.Get(0).Attr) > 0
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.FilterFunction(f).Length()
		} else {
			sel.FilterFunction(f)
		}
	}
	b.Logf("FilterFunction=%d", n)
}

func BenchmarkNotFunction(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Root.Find("li")
	f := func(i int, s *Selection) bool {
		return len(s.Get(0).Attr) > 0
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n == 0 {
			n = sel.NotFunction(f).Length()
		} else {
			sel.NotFunction(f)
		}
	}
	b.Logf("NotFunction=%d", n)
}

func BenchmarkFilterNodes(b *testing.B) {

}
