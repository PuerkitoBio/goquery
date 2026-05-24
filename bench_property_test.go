package goquery

import (
	"testing"
)

func BenchmarkAttr(b *testing.B) {
	var s string

	b.StopTimer()
	sel := DocW().Find("h1")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s, _ = sel.Attr("id")
	}
	if s != "firstHeading" {
		b.Fatalf("want firstHeading, got %q", s)
	}
}

func BenchmarkText(b *testing.B) {
	b.StopTimer()
	sel := DocW().Find("h2")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sel.Text()
	}
}

func BenchmarkLength(b *testing.B) {
	var n int

	b.StopTimer()
	sel := DocW().Find("h2")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		n = sel.Length()
	}
	if n != 14 {
		b.Fatalf("want 14, got %d", n)
	}
}

func BenchmarkHtml(b *testing.B) {
	b.StopTimer()
	sel := DocW().Find("h2")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = sel.Html()
	}
}

func BenchmarkAddClass(b *testing.B) {
	b.StopTimer()
	sel := DocW().Find("li")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sel.AddClass("foo bar baz")
	}
}

func BenchmarkRemoveClass(b *testing.B) {
	b.StopTimer()
	sel := DocW().Find("li")
	sel.AddClass("foo bar baz")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sel.RemoveClass("foo bar baz")
	}
}
