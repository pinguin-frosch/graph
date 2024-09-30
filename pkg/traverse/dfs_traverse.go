package traverse

import (
	"graph/pkg/collections"
	"graph/pkg/graph"
)

func DfsTraverse(g graph.Graph, start graph.Node) ([]graph.Node, error) {
	// setup initial values
	visited := make(map[string]bool)
	s := collections.NewStack[graph.Node]()

	// add starting node
	s.Push(start)

	nodes := make([]graph.Node, 0)
	var x graph.Node

	// Iterate as long as there are elements
	for !s.Empty() {

		// Get current element
		x = s.Pop()

		if !visited[x.Id] {
			// Mask as visited and add to the nodes slice
			visited[x.Id] = true
			nodes = append(nodes, x)

			// Add the valid neighbours for iteration
			neighbours := g.GetNodesFrom(x)
			for _, n := range neighbours {
				if !visited[n.Id] {
					s.Push(n)
				}
			}
		}
	}

	return nodes, nil
}
