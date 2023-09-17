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

	g.AddEdge("B", "A")
	g.AddEdge("B", "C")
	g.AddEdge("B", "E")

	g.AddEdge("C", "B")
	g.AddEdge("C", "D")
	g.AddEdge("C", "F")

	g.AddEdge("D", "C")

	g.AddEdge("E", "B")
	g.AddEdge("E", "F")

	g.AddEdge("F", "E")
	g.AddEdge("F", "C")

	g.Print()

	g.WalkFrom("A")
}
