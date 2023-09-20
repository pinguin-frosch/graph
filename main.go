package main

import (
	"fmt"
	"graph/pkg/graph"
	"log"
)

func main() {
	g, err := graph.NewFromFile("graph.json")
	if err != nil {
		log.Fatal(err)
	}
	g.Print()
	sequence, distance, err := g.GetShortestWalk()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sequence (%v steps, %v m): %v\n", len(sequence), distance, sequence)
}
