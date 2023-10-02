package main

import (
	"graph/pkg/graph"
	"log"
)

func main() {
	g, err := graph.NewFromFile("graph.json")
	if err != nil {
		log.Fatal(err)
	}
	g.Print()
	result, err := g.GetShortestWalk()
	if err != nil {
		log.Fatal(err)
	}
	result.Print()
}
