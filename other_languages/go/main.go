package main

import (
	"fmt"
	"log"
	"slices"
	"strings"
)

func main() {
	solution, err := solve("jan", "1")
	if err != nil {
		log.Fatal(err)
	}

	gridCopy := make([][]string, gridHeight)
	for i := range gridCopy {
		gridCopy[i] = slices.Clone(grid[i][:])
	}

	for _, piecePos := range solution {
		for i := range piecePos.height {
			for j := range piecePos.width {
				if (1<<(i*gridWidth+j))&piecePos.bitmap > 0 {
					gridCopy[piecePos.i+int(i)][piecePos.j+int(j)] = piecePos.id
				}
			}
		}
	}

	printSolution(solution)
}

func printSolution(solutionPos []PieceAndPosition) {
	solution := make([][]string, gridHeight)
	for i := range solution {
		solution[i] = slices.Clone(grid[i][:])
	}

	for _, piecePos := range solutionPos {
		for i := range piecePos.height {
			for j := range piecePos.width {
				if (1<<(i*gridWidth+j))&piecePos.bitmap > 0 {
					solution[piecePos.i+int(i)][piecePos.j+int(j)] = piecePos.id
				}
			}
		}
	}

	sb := strings.Builder{}

	for i := range solution {
		for j := range solution[i] {
			sb.WriteRune('[')
			target := solution[i][j]
			switch len(target) {
			case 0:
				sb.WriteString("   ")
			case 1:
				sb.WriteString("  ")
				sb.WriteString(target)
			case 2:
				sb.WriteString(" ")
				sb.WriteString(target)
			case 3:
				sb.WriteString(target)
			default:
				sb.WriteString("invalid")
			}
			sb.WriteRune(']')
		}
		sb.WriteRune('\n')
	}

	fmt.Println(sb.String())
}
