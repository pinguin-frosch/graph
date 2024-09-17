package traverse

import (
	"fmt"
	"graph/pkg/collections"
	"graph/pkg/graph"
)

type eulerState struct {
	visitedEdges map[string]bool
	invalidEdges map[string]bool
}

func newEulerState() eulerState {
	es := eulerState{}
	es.visitedEdges = make(map[string]bool)
	es.invalidEdges = make(map[string]bool)
	return es
}

func (es *eulerState) allEdgesVisited() bool {
	for _, v := range es.visitedEdges {
		if !v {
			return false
		}
	}
	return true
}

func Euler(g graph.Graph, a graph.Node) (Sequence, error) {
	// setup initial values
	es := newEulerState()
	edges := g.GetAllEdges()
	for _, edge := range edges {
		es.visitedEdges[edge.Key()] = false
	}
	st := collections.NewStack[graph.Node]()

	// add starting node
	st.Push(a)

	x := a
	for !es.allEdgesVisited() {
		// get valid edges to go to the next node
		edges := g.GetEdgesFrom(x)
		validEdges := make([]graph.Edge, 0, len(edges))
		for _, edge := range edges {
			if !es.visitedEdges[edge.Key()] {
				key := fmt.Sprintf("%d|%s|%s", st.Len(), x, edge.To)
				if !es.invalidEdges[key] {
					validEdges = append(validEdges, edge)
				}
			}
		}

		// there's nowhere to go, we need to go back
		if len(validEdges) == 0 {
			// add restriction
			invalidNode := st.Pop()
			previousNode := st.Peek()
			key := fmt.Sprintf("%d|%s|%s", st.Len(), previousNode, invalidNode)
			es.invalidEdges[key] = true

			// mark the edge as not visited again
			es.visitedEdges[g.GetEdge(previousNode, invalidNode).Key()] = false
			es.visitedEdges[g.GetEdge(previousNode, invalidNode).ReversedEdge().Key()] = false

			// go back
			x = previousNode
			continue
		}

		// add next node
		nextEdge := validEdges[0]
		es.visitedEdges[nextEdge.Key()] = true
		es.visitedEdges[nextEdge.ReversedEdge().Key()] = true
		nextNode := nextEdge.To
		st.Push(nextNode)
		x = nextNode
	}

	// reconstruct the sequence from the stack
	s := NewSequence()
	for !st.Empty() {
		lastNode := st.Pop()
		s.Sequence = append(s.Sequence, lastNode)
	}
	for i := 0; i < len(s.Sequence)-1; i++ {
		a := s.Sequence[i]
		b := s.Sequence[i+1]
		s.Distance += g.GetEdge(a, b).Weight
	}

	return s, nil
}
