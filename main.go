package main

import (
	"graph/pkg/graph"
	"log"
)

func main() {
	g, err := graph.NewInteractively()
	if err != nil {
		log.Fatal(err)
	}
	g.Print()
}
