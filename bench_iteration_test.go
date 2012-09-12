package goquery

import (
	"testing"
)

func BenchmarkEach(b *testing.B) {
	var tmp, n int

	b.StopTimer()
	sel := DocW().Root.Find("td")
	f := func(i int, s *Selection) {
		tmp++
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sel.Each(f)
		if n == 0 {
			n = tmp
		}
	}
	b.Logf("Each=%d", n)
}

func BenchmarkMap(b *testing.B) {
	var tmp, n int

	b.StopTimer()
	sel := DocW().Root.Find("td")
	f := func(i int, s *Selection) string {
		tmp++
		return string(tmp)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sel.Map(f)
		if n == 0 {
			n = tmp
		}
	}
	b.Logf("Map=%d", n)
}
