package main

import (
	"graph/pkg/graph"
	"log"
)

func main() {
	g, err := graph.LoadFromJson("graph.json")
	if err != nil {
		log.Fatal(err)
	}
	g.Print()
}
