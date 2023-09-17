package main

import (
	"fmt"
	"graph/pkg/graph"
)

func main() {
	g := &graph.Graph{}
	g.AddVertex("A")
	g.AddVertex("B")

	fmt.Println(g)
}
