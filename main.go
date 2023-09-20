package main

import (
	"fmt"
	"graph/pkg/graph"
	"log"
)

func main() {
	g, err := graph.NewFromFile("graph2.json")
	if err != nil {
		log.Fatal(err)
	}
	g.Print()
	sequence, err := g.GetShortestWalk()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sequence (%v): %v\n", len(sequence), sequence)
}
