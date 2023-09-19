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
	g.Print()
	sequence, err := g.WalkFrom("W")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sequence: %v\n", sequence)
}
