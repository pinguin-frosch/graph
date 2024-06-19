package graph

import (
	"errors"
	"fmt"
)

type DefaultTraverseAlgorithm struct{}

func (d DefaultTraverseAlgorithm) getSequence(g *Graph, from string) (*ResultSequence, error) {
	g.resetState()
	result := ResultSequence{
		Distance: 0,
	}
	result.Sequence = make([]string, 0)
	totalEdges := g.GetTotalEdges()
	var usedEdges uint = 0

	vertex := g.GetVertex(from)
	if vertex == "" {
		return nil, errors.New("starting vertex does not exist")
	}
	result.Sequence = append(result.Sequence, vertex)

	for usedEdges < totalEdges {
		// Get the next edge to use
		edge, unique := d.nextEdge(g, vertex)
		if edge == nil {
			return nil, ErrInvalidNextEdge
		}

		// Get the other edge
		otherEdge := g.GetEdge(edge.To, edge.From)
		if otherEdge == nil {
			return nil, fmt.Errorf("couldn't find the other edge, from %v to %v", edge.To, edge.From)
		}

		// Only count unused edges
		if edge.visitCount == 0 {
			usedEdges++
		}

		// Update usage, if it was the only option increase by two because it's a dead end
		if unique && edge.visitCount == 0 {
			otherEdge.visitCount += 2
			edge.visitCount += 2
		} else {
			otherEdge.visitCount += 1
			edge.visitCount += 1
		}

		result.Distance += edge.Weight
		result.Sequence = append(result.Sequence, edge.To)
		vertex = edge.To
	}

	return &result, nil
}

func (d DefaultTraverseAlgorithm) nextEdge(g *Graph, key string) (*Edge, bool) {
	edges, result := g.GetEdges(key)
	unique := len(edges) == 1

	if result.DeadEndCount == 1 {
		for _, e := range edges {
			if !e.DeadEnd {
				continue
			}
			if e.visitCount < 2 {
				return e, unique
			}
		}
	} else if result.DeadEndCount > 1 {
		for _, e := range edges {
			if !e.DeadEnd {
				continue
			}
			if e.visitCount > result.MinimumDeadEnd {
				continue
			}
			if e.visitCount < 2 {
				return e, unique
			}
		}
	}
	for _, e := range edges {
		if e.DeadEnd {
			continue
		}
		if e.visitCount > result.Minimum {
			continue
		}
		if e.Weight > result.Weight {
			continue
		}
		return e, unique
	}

	return nil, unique
}
