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
	sequence, err := g.WalkFrom("A")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sequence: %v\n", sequence)
}
