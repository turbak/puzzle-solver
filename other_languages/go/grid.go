package main

import (
	"slices"
	"strings"
)

var grid = [][]string{
	{"jan", "feb", "mar", "apr", "may", "jun"},
	{"jul", "aug", "sep", "oct", "nov", "dec"},
	{"1", "2", "3", "4", "5", "6", "7"},
	{"8", "9", "10", "11", "12", "13", "14"},
	{"15", "16", "17", "18", "19", "20", "21"},
	{"22", "23", "24", "25", "26", "27", "28"},
	{"29", "30", "31"},
}

type position struct {
	i, j int
}

type Grid struct {
	matrix             [][]string
	protectedPositions []position
}

func newGrid(protectedPositions []position) Grid {
	gridCopy := make([][]string, len(grid))
	for i := range grid {
		gridCopy[i] = slices.Clone(grid[i])
	}

	return Grid{
		matrix:             gridCopy,
		protectedPositions: protectedPositions,
	}
}

func (g *Grid) Place(piece Piece, gridI int, gridJ int) {
	for i := range piece.matrix {
		for j := range piece.matrix[i] {
			if piece.matrix[i][j] == "" {
				continue
			}

			gridRow := i + gridI
			gridCol := j + gridJ

			g.matrix[gridRow][gridCol] = piece.matrix[i][j]
		}
	}
}

func (g *Grid) Unplace(piece Piece, gridI int, gridJ int) {
	for i := range piece.matrix {
		for j := range piece.matrix[i] {
			if piece.matrix[i][j] == "" {
				continue
			}

			gridRow := i + gridI
			gridCol := j + gridJ

			g.matrix[gridRow][gridCol] = grid[gridRow][gridCol]
		}
	}
}

func (g Grid) CanPlace(piece Piece, gridI int, gridJ int) bool {
	for i := range piece.matrix {
		for j := range piece.matrix[i] {
			if piece.matrix[i][j] == "" {
				continue
			}

			gridRow := i + gridI
			gridCol := j + gridJ

			if len(g.matrix) <= gridRow {
				return false
			}

			if len(g.matrix[gridRow]) <= gridCol {
				return false
			}

			if strings.HasPrefix(g.matrix[gridRow][gridCol], "{") {
				return false
			}

			if slices.ContainsFunc(g.protectedPositions, func(p position) bool {
				return p.i == gridRow && p.j == gridCol
			}) {
				return false
			}
		}
	}

	return true
}

func (g Grid) String() string {
	sb := strings.Builder{}

	for i := range g.matrix {
		for j := range g.matrix[i] {
			sb.WriteRune('[')
			switch len(g.matrix[i][j]) {
			case 1:
				sb.WriteString("  ")
				sb.WriteString(g.matrix[i][j])
			case 2:
				sb.WriteString(" ")
				sb.WriteString(g.matrix[i][j])
			case 3:
				sb.WriteString(g.matrix[i][j])
			default:
				sb.WriteString("invalid")
			}
			sb.WriteRune(']')
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}
