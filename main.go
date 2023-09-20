package main

import (
	"fmt"
	"graph/pkg/graph"
	"log"
)

func main() {
	g, err := graph.LoadFromJson("graph.json")
	if err != nil {
		log.Fatal(err)
	}
	sequence, err := g.GetShortestWalk()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sequence (%v): %v\n", len(sequence), sequence)
}
