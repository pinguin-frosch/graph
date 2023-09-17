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

	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.AddEdge("B", "E")
	g.AddEdge("C", "D")
	g.AddEdge("C", "F")
	g.AddEdge("E", "F")

	g.Print()
}
