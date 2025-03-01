package main

import "fmt"

func solve(month string, day string) (Grid, error) {
	var monthPos, dayPos position

	for i, row := range grid {
		for j, val := range row {
			if val == month {
				monthPos = position{
					i: i,
					j: j,
				}
			} else if val == day {
				dayPos = position{
					i: i,
					j: j,
				}
			}
		}
	}

	g := newGrid([]position{monthPos, dayPos})
	solved := solveHelper(&g, 0)
	if !solved {
		return Grid{}, fmt.Errorf("no solution found for %v", []position{monthPos, dayPos})
	}

	return g, nil
}

func solveHelper(g *Grid, pieceIdx int) bool {
	if pieceIdx >= len(allPiecesWithTranspositions) {
		return true
	}

	for _, pieceTransposition := range allPiecesWithTranspositions[pieceIdx] {
		for i, gridRow := range g.matrix {
			for j := range gridRow {
				if !g.CanPlace(pieceTransposition, i, j) {
					continue
				}

				g.Place(pieceTransposition, i, j)
				if solved := solveHelper(g, pieceIdx+1); solved {
					return true
				}
				g.Unplace(pieceTransposition, i, j)
			}
		}
	}

	return false
}
