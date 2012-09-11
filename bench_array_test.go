package goquery

import (
	"testing"
)

func BenchmarkFirst(b *testing.B) {
	b.StopTimer()
	sel := DocB().Root.Find("dd")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sel.First()
	}
}

func BenchmarkLast(b *testing.B) {
	b.StopTimer()
	sel := DocB().Root.Find("dd")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sel.Last()
	}
}

func BenchmarkEq(b *testing.B) {
	b.StopTimer()
	sel := DocB().Root.Find("dd")
	j := 0
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sel.Eq(j)
		if j++; j >= sel.Length() {
			j = 0
		}
	}
}

func BenchmarkSlice(b *testing.B) {
	b.StopTimer()
	sel := DocB().Root.Find("dd")
	j := 0
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sel.Slice(j, j+4)
		if j++; j >= (sel.Length() - 4) {
			j = 0
		}
	}
}

func BenchmarkGet(b *testing.B) {
	b.StopTimer()
	sel := DocB().Root.Find("dd")
	j := 0
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sel.Get(j)
		if j++; j >= sel.Length() {
			j = 0
		}
	}
}

func BenchmarkIndex(b *testing.B) {
	b.StopTimer()
	sel := DocB().Root.Find("#main")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sel.Index()
	}
}
