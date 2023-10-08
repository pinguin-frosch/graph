package main

import (
	"graph/pkg/graph"
	"log"
)

func main() {
	g := graph.Graph{}
	err := g.ModifyInteractively()
	if err != nil {
		log.Fatal(err)
	}
}
