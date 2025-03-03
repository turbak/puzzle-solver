package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"math/bits"
	"net/http"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"
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

func SolveHandler(w http.ResponseWriter, r *http.Request) {
	month := r.URL.Query().Get("month")
	day := r.URL.Query().Get("day")

	now := time.Now()
	res, err := solve(month, day)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	solutionTime := time.Since(now)

	log.Printf("Solution for %s %s took %v\n", month, day, solutionTime)

	solution := make([][]string, gridHeight)
	for i := range solution {
		solution[i] = slices.Clone(grid[i][:])
	}

	for _, piecePos := range res {
		for i := range piecePos.height {
			for j := range piecePos.width {
				if (1<<(i*gridWidth+j))&piecePos.bitmap > 0 {
					solution[piecePos.i+int(i)][piecePos.j+int(j)] = piecePos.id
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(solution)
}

type PieceAndPosition struct {
	Piece
	position
}

func solve(month string, day string) ([]PieceAndPosition, error) {
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

	resCh := make(chan []PieceAndPosition)
	wg := sync.WaitGroup{}
	wg.Add(len(allPiecesWithTranspositions[0]))

	go func() {
		wg.Wait()
		close(resCh)
	}()

	for i := range min(len(allPiecesWithTranspositions[0]), runtime.NumCPU()) {
		go func(pos int) {
			defer wg.Done()

			g := newGrid([]position{monthPos, dayPos})
			res, solved := solveHelper(g, 0, pos)
			if solved {
				resCh <- res
			}
		}(i)
	}

	res := <-resCh

	return res, nil
}

func solveHelper(g Grid, pieceIdx int, startPos int) ([]PieceAndPosition, bool) {
	if pieceIdx >= len(allPiecesWithTranspositions) {
		return make([]PieceAndPosition, 0, len(allPiecesWithTranspositions)), true
	}

	for _, pieceTransposition := range allPiecesWithTranspositions[pieceIdx][startPos:] {
		for i := uint64(0); i < gridHeight-pieceTransposition.height+1; i++ {
			for j := uint64(0); j < gridWidth-pieceTransposition.width+1; j++ {
				if !g.CanPlace(pieceTransposition, int(i), int(j)) {
					continue
				}

				gridCopy := g
				gridCopy.Place(pieceTransposition, int(i), int(j))
				if res, solved := solveHelper(gridCopy, pieceIdx+1, 0); solved {
					return append(res, PieceAndPosition{
						Piece: pieceTransposition,
						position: position{
							i: int(i),
							j: int(j),
						},
					}), true
				}
			}
		}
	}

	return nil, false
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
