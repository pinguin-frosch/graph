package main

import (
	"graph/pkg/graph"
)

func main() {
	g := &graph.Graph{}
	g.AddVertex("A")
	g.AddVertex("B")
	g.AddVertex("C")
	g.AddVertex("D")
	g.AddVertex("E")
	g.AddVertex("F")
	g.Print()
}
