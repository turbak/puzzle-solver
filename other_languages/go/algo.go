package main

import (
	"sync"
)

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

	for i := range allPiecesWithTranspositions[0] {
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
