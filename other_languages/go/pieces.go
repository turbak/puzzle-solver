package main

import (
	"slices"
	"strings"
)

func init() {
	registerPiece([][]string{
		{"{1}", ""},
		{"{1}", "{1}"},
		{"{1}", "{1}"},
	})

	registerPiece([][]string{
		{"{2}", "{2}", "{2}", "{2}"},
		{"{2}", "", "", ""},
	})

	registerPiece([][]string{
		{"{3}", "{3}", "{3}", ""},
		{"", "", "{3}", "{3}"},
	})

	registerPiece([][]string{
		{"{4}", "{4}", "{4}"},
		{"{4}", "{4}", "{4}"},
	})

	registerPiece([][]string{
		{"{5}", "{5}", ""},
		{"", "{5}", ""},
		{"", "{5}", "{5}"},
	})

	registerPiece([][]string{
		{"{6}", "{6}", "{6}", "{6}"},
		{"", "{6}", "", ""},
	})

	registerPiece([][]string{
		{"{7}", "{7}"},
		{"", "{7}"},
		{"{7}", "{7}"},
	})

	registerPiece([][]string{
		{"{8}", "", ""},
		{"{8}", "", ""},
		{"{8}", "{8}", "{8}"},
	})

	slices.SortFunc(allPiecesWithTranspositions, func(a, b []Piece) int {
		countA, countB := 0, 0

		for _, row := range a[0].matrix {
			for _, val := range row {
				if val != "" {
					countA += 1
				}
			}
		}

		for _, row := range b[0].matrix {
			for _, val := range row {
				if val != "" {
					countB += 1
				}
			}
		}

		return countA - countB //descending
	})
}

var allPiecesWithTranspositions [][]Piece

func registerPiece(matrix [][]string) {
	currentPiece := newPiece(matrix)
	pieces := []Piece{currentPiece}
	// number of rotations
	for range 4 {
		flipped := currentPiece.Flip()
		if !slices.ContainsFunc(pieces, func(piece Piece) bool {
			return flipped.Equal(piece)
		}) {
			pieces = append(pieces, flipped)
		}

		rotated := currentPiece.RotateClockwise()
		if slices.ContainsFunc(pieces, func(piece Piece) bool {
			return rotated.Equal(piece)
		}) {
			break
		}

		pieces = append(pieces, rotated)
		currentPiece = rotated
	}

	allPiecesWithTranspositions = append(allPiecesWithTranspositions, pieces)
}

type Piece struct {
	matrix [][]string
}

func newPiece(matrix [][]string) Piece {
	return Piece{matrix: matrix}
}

func (p *Piece) RotateClockwise() Piece {
	if len(p.matrix) == 0 {
		return Piece{}
	}

	rotated := make([][]string, len(p.matrix[0]))

	for i := range p.matrix[0] {
		rotated[i] = make([]string, len(p.matrix))
		for j := range p.matrix {
			rotated[i][len(p.matrix)-j-1] = p.matrix[j][i]
		}
	}

	return newPiece(rotated)
}

func (p *Piece) Flip() Piece {
	flipped := make([][]string, len(p.matrix))
	for i := range p.matrix {
		flipped[i] = make([]string, len(p.matrix[i]))
		for j, jFlipped := len(p.matrix[i])-1, 0; j >= 0; j, jFlipped = j-1, jFlipped+1 {
			flipped[i][jFlipped] = p.matrix[i][j]
		}
	}

	return newPiece(flipped)
}

func (p Piece) String() string {
	sb := strings.Builder{}

	for i := range p.matrix {
		for j := range p.matrix[i] {
			sb.WriteRune('[')

			if len(p.matrix[i][j]) == 0 {
				sb.WriteString("   ")
			} else {
				sb.WriteString(p.matrix[i][j])
			}

			sb.WriteRune(']')
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}

func (p Piece) Equal(other Piece) bool {
	for i := range p.matrix {
		if !slices.Equal(p.matrix[i], other.matrix[i]) {
			return false
		}
	}

	return true
}
