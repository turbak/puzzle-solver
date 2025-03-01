package main

import (
	"testing"
	"time"
)

func BenchmarkSolve(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = solve("jan", "1")
	}
}

func TestSolve(t *testing.T) {
	start := time.Now()
	solution, err := solve("jan", "1")
	if err != nil {
		t.Fatal(err)
	}
	elapsed := time.Since(start)
	t.Logf("Solution: \n%v", solution)
	t.Logf("Elapsed: %v", elapsed)
}
