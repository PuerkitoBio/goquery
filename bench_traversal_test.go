package goquery

import (
	"testing"
)

func BenchmarkFind(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DocB().Root.Find("dd")
	}
}
