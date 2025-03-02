package handler

import (
	"fmt"
	"math/bits"
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
		fieldsDiff := b[0].Area() - a[0].Area()
		if fieldsDiff != 0 {
			return fieldsDiff
		}

		areaDiff := b[0].height*b[0].width - a[0].height*a[0].width
		if areaDiff != 0 {
			return int(areaDiff)
		}

		return len(a) - len(b)
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
			return flipped.bitmap == piece.bitmap
		}) {
			pieces = append(pieces, flipped)
		}

		rotated := currentPiece.RotateClockwise()
		if slices.ContainsFunc(pieces, func(piece Piece) bool {
			return rotated.bitmap == piece.bitmap
		}) {
			break
		}

		pieces = append(pieces, rotated)
		currentPiece = rotated
	}

	allPiecesWithTranspositions = append(allPiecesWithTranspositions, pieces)
}

type Piece struct {
	id     string
	width  uint64
	height uint64
	bitmap uint64
}

func newPiece(matrix [][]string) Piece {
	var bitmap uint64
	height := len(matrix)
	if height == 0 {
		return Piece{}
	}
	width := len(matrix[0])

	idStr := ""

	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			if matrix[row][col] != "" {
				pos := row*gridWidth + col
				bitmap |= 1 << pos

				idStr = matrix[row][col]
			}
		}
	}

	return Piece{
		id:     idStr,
		width:  uint64(width),
		height: uint64(height),
		bitmap: bitmap,
	}
}

func (p *Piece) RotateClockwise() Piece {
	rotated := Piece{
		width:  p.height,
		height: p.width,
		id:     p.id,
		bitmap: 0,
	}

	for i := uint64(0); i < gridHeight; i++ {
		for j := uint64(0); j < gridWidth; j++ {
			if (1<<(j*gridWidth+i))&p.bitmap > 0 {
				rotated.bitmap |= 1 << (i*gridWidth + gridWidth - j - 1)
			}
		}
	}

	rotated.bitmap = trimTrailingZeroes(rotated.bitmap)

	return rotated
}

func (p *Piece) Flip() Piece {
	flipped := Piece{
		id:     p.id,
		width:  p.width,
		height: p.height,
		bitmap: 0,
	}

	for i := range p.height {
		var j, jFlipped uint64 = gridWidth, 0
		for ; j > 0; j, jFlipped = j-1, jFlipped+1 {
			if (1<<(i*gridWidth+j-1))&p.bitmap > 0 {
				flipped.bitmap |= 1 << (i*gridWidth + jFlipped)
			}
		}
	}

	flipped.bitmap = trimTrailingZeroes(flipped.bitmap)

	return flipped
}

func (p Piece) String() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("%b\n", p.bitmap))

	for i := range p.height {
		for j := range p.width {
			sb.WriteRune('[')
			if (1<<(i*gridWidth+j))&p.bitmap > 0 {
				sb.WriteString(p.id)
			} else {
				sb.WriteString("   ")
			}

			sb.WriteRune(']')
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}

func (p Piece) Area() int {
	return bits.OnesCount64(p.bitmap)
}

func trimTrailingZeroes(bitmap uint64) uint64 {
	rowMask := uint64(1<<gridWidth - 1)
	colMask := uint64(0)
	for i := range gridHeight {
		colMask |= 1 << (i * gridWidth)
	}

	for bitmap&rowMask == 0 {
		bitmap >>= gridWidth
	}

	for bitmap&colMask == 0 {
		bitmap >>= 1
	}

	return bitmap
}
