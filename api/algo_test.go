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

	start := time.Now()
	for _, month := range allMonth {
		for _, day := range allDays {

			_, err := solve(month, day)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
	avgElapsed := time.Since(start)

	t.Logf("Average elapsed: %v", avgElapsed/(time.Duration(len(allMonth)*len(allDays))))
	t.Logf("Total elapsed: %v", time.Since(start))
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
