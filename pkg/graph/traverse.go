package graph

import (
	"errors"
	"fmt"
	"math"
)

type TraverseAlgorithm interface {
	getSequence(g *Graph, from string) (*ResultSequence, error)
}

var (
	ErrInvalidNextEdge = errors.New("received an invalid edge")
	ErrNoAlgorithmSet  = errors.New("no traverse algorithm has been set")
)

func (g *Graph) SetTraverseAlgorithm(a TraverseAlgorithm) {
	g.traverseAlgorithm = a
}

func (g *Graph) GetSequence(from string) (*ResultSequence, error) {
	if g.traverseAlgorithm == nil {
		return nil, ErrNoAlgorithmSet
	}
	r, err := g.traverseAlgorithm.getSequence(g, from)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (g *Graph) GetShortestSequence() (*ResultSequence, error) {
	var shortest = math.MaxInt
	result := ResultSequence{
		Sequence: make([]string, 0),
		Distance: 0,
	}
	for _, v := range g.Vertices {
		r, err := g.GetSequence(v)
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
