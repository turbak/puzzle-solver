package handler

import (
	"encoding/json"
	"net/http"
	"slices"
)

func SolveHandler(w http.ResponseWriter, r *http.Request) {
	month := r.URL.Query().Get("month")
	day := r.URL.Query().Get("day")

	res, err := solve(month, day)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	json.NewEncoder(w).Encode(solution)
}
