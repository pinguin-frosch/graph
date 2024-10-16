package traverse

import (
	"fmt"
	"graph/pkg/graph"
	"math"
	"slices"
)

type nodeState struct {
	visited bool
	value   int
	prev    graph.Node
}

type dijkstraState struct {
	nodes map[string]*nodeState
}

func newDijkstraState() dijkstraState {
	ds := dijkstraState{}
	ds.nodes = make(map[string]*nodeState)
	return ds
}

func (ds dijkstraState) allNodesVisited() bool {
	for _, ns := range ds.nodes {
		if !ns.visited {
			return false
		}
	}
	return true
}

func Dijkstra(g graph.Graph, a, b graph.Node) (Sequence, error) {
	// setup initial values
	ds := newDijkstraState()
	nodes := g.GetAllNodes()
	for _, node := range nodes {
		if _, ok := ds.nodes[node.Id]; !ok {
			ds.nodes[node.Id] = &nodeState{}
		}
		ds.nodes[node.Id].value = math.MaxInt
	}
	ds.nodes[a.Id].value = 0
	ds.nodes[a.Id].visited = true

	x := a
	for !ds.allNodesVisited() {
		if _, ok := ds.nodes[x.Id]; !ok {
			// find a new starting point
			maxValue := math.MaxInt
			newNode := graph.Node{}
			for _, node := range nodes {
				if ds.nodes[node.Id].visited {
					continue
				}
				if ds.nodes[node.Id].value < maxValue {
					maxValue = ds.nodes[node.Id].value
					newNode = node
				}
			}
			x = newNode
			continue
		}

		// update estimates
		nodesFromX := g.GetNodes(x)
		for _, y := range nodesFromX {
			edge, ok := g.GetShortestEdge(x, y)
			if !ok {
				return Sequence{}, fmt.Errorf("couldn't get shortest edge between %v and %v", x.Id, y.Id)
			}
			xValue := ds.nodes[x.Id].value
			yValue := ds.nodes[y.Id].value
			if (xValue + edge.Weight) < yValue {
				ds.nodes[y.Id].value = xValue + edge.Weight
				ds.nodes[y.Id].prev = x
			}
		}

		// choose the next node
		availableNodes := make([]graph.Node, 0)
		for _, y := range nodesFromX {
			if !ds.nodes[y.Id].visited {
				availableNodes = append(availableNodes, y)
			}
		}
		if len(availableNodes) == 0 {
			x = ds.nodes[x.Id].prev
			continue
		}
		slices.SortFunc(availableNodes, func(a, b graph.Node) int {
			if ds.nodes[a.Id].value < ds.nodes[b.Id].value {
				return -1
			} else if ds.nodes[a.Id].value > ds.nodes[b.Id].value {
				return 1
			} else {
				return 0
			}
		})
		y := availableNodes[0]
		ds.nodes[y.Id].visited = true
		x = y
	}

	// go back and reconstruct the sequence
	s := NewSequence()
	s.Distance = ds.nodes[b.Id].value
	s.Sequence = append(s.Sequence, b)
	for b.Id != a.Id {
		b = ds.nodes[b.Id].prev
		s.Sequence = append(s.Sequence, b)
	}
	slices.Reverse(s.Sequence)

	return s, nil
}
