package main

import "graph/pkg/aograph"

func main() {
	g := aograph.NewGraph()

	// Add nodes
	nodes := map[string]float64{"A": -1, "B": 5, "C": 2, "D": 4, "E": 7, "F": 9, "G": 3, "H": 0, "I": 0, "J": 0}
	for id, value := range nodes {
		graphNode := aograph.NewNode(id, value)
		g.AddNode(graphNode)
	}

	// Add connections
	g.AddOrConnection("A", []string{"B"}...)
	g.AddAndConnection("A", []string{"C", "D"}...)
	g.AddOrConnection("B", []string{"E", "F"}...)
	g.AddOrConnection("C", []string{"G"}...)
	g.AddAndConnection("C", []string{"H", "I"}...)
	g.AddOrConnection("D", []string{"J"}...)

	// Update cost
	g.UpdateCost()

	g.Print()
}
