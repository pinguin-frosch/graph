package graph

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrInvalidNextEdge = errors.New("Received an invalid edge!")
)

func (g *Graph) GetShortestWalk() (*ResultSequence, error) {
	var shortest int = math.MaxInt
	result := ResultSequence{
		Sequence: make([]string, 0),
		Distance: 0,
	}
	for _, v := range g.Vertices {
		r, err := g.WalkFrom(v)
		if err != nil && !errors.Is(err, ErrInvalidNextEdge) {
			return nil, err
		}
		if errors.Is(err, ErrInvalidNextEdge) {
			continue
		}
		if len(r.Sequence) < shortest {
			shortest = len(r.Sequence)
			result.Distance = r.Distance
			result.Sequence = r.Sequence
		}
	}
	return &result, nil
}

func (g *Graph) WalkFrom(from string) (*ResultSequence, error) {
	g.resetState()
	result := ResultSequence{
		Distance: 0,
	}
	result.Sequence = make([]string, 0)
	totalEdges := g.GetTotalEdges()
	var usedEdges uint = 0

	vertex := g.GetVertex(from)
	if vertex == "" {
		return nil, errors.New("The starting vertex does not exist!")
	}
	result.Sequence = append(result.Sequence, vertex)

	for usedEdges < totalEdges {
		// Get the next edge to use
		edge, unique := g.GetNextEdge(vertex)
		if edge == nil {
			return nil, ErrInvalidNextEdge
		}

		// Get the other edge
		otherEdge := g.GetEdge(edge.To, edge.From)
		if otherEdge == nil {
			return nil, errors.New(fmt.Sprintf("Couldn't find the other edge, from %v to %v", edge.To, edge.From))
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

func (g *Graph) GetTotalEdges() uint {
	var total uint = 0
	matrix := make(map[string]bool)
	for _, e := range g.Edges {
		key := fmt.Sprintf("%v|%v", e.From, e.To)
		otherKey := fmt.Sprintf("%v|%v", e.To, e.From)
		if !matrix[key] && !matrix[otherKey] {
			total += 1
			matrix[key] = true
			matrix[otherKey] = true
		}
	}
	return total
}

func (g *Graph) resetState() {
	for _, e := range g.Edges {
		e.visitCount = 0
	}
}

// Returns the next edge to use in the search algorithm and wether it was the only one found
func (g *Graph) GetNextEdge(key string) (*Edge, bool) {
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

// Returns all the edges from a certain vertex, plus useful stats to later
// decide. This is not supposed to be used directly.
func (g *Graph) GetEdges(key string) ([]*Edge, EdgesResult) {
	result := EdgesResult{
		Minimum:        math.MaxUint,
		MinimumDeadEnd: math.MaxUint,
		DeadEndCount:   0,
		Weight:         0,
	}
	edges := make([]*Edge, 0)

	// Compute all the required information
	for _, e := range g.Edges {
		if e.From == key {
			switch e.DeadEnd {
			case true:
				result.DeadEndCount++
				if e.visitCount < result.MinimumDeadEnd {
					result.MinimumDeadEnd = e.visitCount
				}
			case false:
				if e.visitCount < result.Minimum {
					result.Minimum = e.visitCount
					result.Weight = e.Weight
				} else if e.visitCount == result.Minimum {
					if e.Weight < result.Weight {
						result.Weight = e.Weight
					}
				}
			}
			edges = append(edges, e)
		}
	}

	return edges, result
}

type ResultSequence struct {
	Sequence []string
	Distance uint
}

func (r *ResultSequence) Print() {
	fmt.Printf("Sequence (%v steps and %v in weight): %v\n", len(r.Sequence), r.Distance, r.Sequence)
}

type EdgesResult struct {
	// The least amount of times non dead end edges has been used
	Minimum uint
	// The least amount of times dead end edges have been used
	MinimumDeadEnd uint
	// The number of dead ends found
	DeadEndCount uint
	// The minimum weight found
	Weight uint
}
