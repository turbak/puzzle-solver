package handler

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkSolve(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = solve("jan", "1")
	}
}

func TestSolveWholeYear(t *testing.T) {
	allMonth := []string{"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"}
	allDays := make([]string, 31)
	for i := 0; i < 31; i++ {
		allDays[i] = fmt.Sprintf("%d", i+1)
	}

	avgElapsed := time.Duration(0)
	for _, month := range allMonth {
		for _, day := range allDays {
			start := time.Now()
			_, err := solve(month, day)
			if err != nil {
				t.Fatal(err)
			}
			elapsed := time.Since(start)

			avgElapsed += elapsed
		}
	}

	t.Logf("Average elapsed: %v", avgElapsed/(time.Duration(len(allMonth)*len(allDays))))
}

func TestSolve(t *testing.T) {
	start := time.Now()
	solution, err := solve("jan", "1")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Elapsed: %v", time.Since(start))
	printSolution(solution)
}
