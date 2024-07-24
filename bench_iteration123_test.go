//go:build go1.23
// +build go1.23

package goquery

import "testing"

func BenchmarkEachIter123(b *testing.B) {
	var tmp, n int

	b.StopTimer()
	sel := DocW().Find("td")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for range sel.EachIter() {
			tmp++
		}
		if n == 0 {
			n = tmp
		}
	}
	if n != 59 {
		b.Fatalf("want 59, got %d", n)
	}
}

func BenchmarkEachIterWithBreak123(b *testing.B) {
	var tmp, n int

	b.StopTimer()
	sel := DocW().Find("td")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp = 0
		for range sel.EachIter() {
			tmp++
			if tmp >= 10 {
				break
			}
		}
		if n == 0 {
			n = tmp
		}
	}
	if n != 10 {
		b.Fatalf("want 10, got %d", n)
	}
}
