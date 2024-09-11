package traverse

import (
	"fmt"
	"graph/pkg/graph"
	"slices"
)

type Default struct {
	visitedEdges map[string]int
	totalEdges   int
	usedEdges    int
}

func NewDefault() Default {
	d := Default{}
	d.visitedEdges = make(map[string]int)
	return d
}

func (d *Default) reset() {
	d.visitedEdges = make(map[string]int)
	d.usedEdges = 0
	d.totalEdges = 0
}

func (d Default) getSequence(g graph.Graph, from graph.Node) (Sequence, error) {
	d.reset()
	node := from
	s := NewSequence()
	s.Sequence = append(s.Sequence, node)
	d.totalEdges = d.CountTotalEdges(g)
	for d.usedEdges < d.totalEdges {
		edge, err := d.GetNextEdge(g, node)
		if err != nil {
			return s, err
		}
		if d.visitedEdges[edge.Key()] == 0 {
			d.usedEdges++
		}
		d.visitedEdges[edge.Key()]++
		d.visitedEdges[edge.ReversedEdge().Key()]++
		s.Distance += edge.Weight()
		s.Sequence = append(s.Sequence, edge.To())
		node = edge.To()
	}
	return s, nil
}

func (d *Default) CountTotalEdges(g graph.Graph) int {
	total := 0
	visited := make(map[string]bool)
	nodes := g.Nodes()
	for _, node := range nodes {
		edges := g.EdgesFrom(node)
		for _, edge := range edges {
			key := edge.Key()
			keyReversed := edge.ReversedEdge().Key()
			if !visited[key] && !visited[keyReversed] {
				visited[key] = true
				visited[keyReversed] = true
				total += 1
			}
		}
	}
	return total
}

func (d *Default) GetNextEdge(g graph.Graph, n graph.Node) (graph.Edge, error) {
	candidates := g.EdgesFrom(n)
	if len(candidates) < 1 {
		return graph.Edge{}, fmt.Errorf("node %s has no edges", n.Id())
	}
	slices.SortFunc(candidates, func(a, b graph.Edge) int {
		if d.visitedEdges[a.Key()] < d.visitedEdges[b.Key()] {
			return -1
		} else if d.visitedEdges[a.Key()] > d.visitedEdges[b.Key()] {
			return 1
		} else {
			if a.Weight() < b.Weight() {
				return -1
			} else if a.Weight() > b.Weight() {
				return 1
			}
		}
		return 0
	})
	return candidates[0], nil
}
