package handler

import (
	"fmt"
	"slices"
	"strings"
)

const gridHeight = 7
const gridWidth = 7

var grid = [gridHeight][gridWidth]string{
	{"jan", "feb", "mar", "apr", "may", "jun", ""},
	{"jul", "aug", "sep", "oct", "nov", "dec", ""},
	{"1", "2", "3", "4", "5", "6", "7"},
	{"8", "9", "10", "11", "12", "13", "14"},
	{"15", "16", "17", "18", "19", "20", "21"},
	{"22", "23", "24", "25", "26", "27", "28"},
	{"29", "30", "31", "", "", "", ""},
}

type position struct {
	i, j int
}

type Grid struct {
	bitmap uint64
}

func newGrid(protectedPositions []position) Grid {
	var bitmap uint64 = 0
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == "" {
				bitmap |= 1 << (i*gridWidth + j)
			}
		}
	}

	for i := range protectedPositions {
		bitmap |= 1 << (protectedPositions[i].i*gridWidth + protectedPositions[i].j)
	}

	return Grid{
		bitmap: bitmap,
	}
}

func (g *Grid) Place(piece Piece, gridI int, gridJ int) {
	g.bitmap |= piece.bitmap << (gridI*gridWidth + gridJ)
}

func (g Grid) CanPlace(piece Piece, gridI int, gridJ int) bool {
	return (piece.bitmap<<(gridI*gridWidth+gridJ))&g.bitmap == 0
}

func (g Grid) String() string {
	sb := strings.Builder{}

	for i := range grid {
		for j := range grid[i] {
			sb.WriteRune('[')
			target := grid[i][j]
			if 1<<(i*gridWidth+j)&g.bitmap > 0 {
				target = "x"
			}
			switch len(target) {
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

	return sb.String()
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
