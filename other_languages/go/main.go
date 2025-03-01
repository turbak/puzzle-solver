package main

import (
	"fmt"
	"log"
)

func main() {
	for _, pieces := range allPiecesWithTranspositions {
		for _, piece := range pieces {
			fmt.Println(piece)
			fmt.Println()
		}
	}

	solution, err := solve("jan", "1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(solution)
}
