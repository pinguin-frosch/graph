package graph

import (
	"errors"
	"fmt"
	"math"
)

const (
	RepeatedVertex = "That vertex already exists!"
	RepeatedEdge   = "That edge already exists!"
)

type ResultSequence struct {
	Sequence []string
	Distance uint
}

func (r *ResultSequence) Print() {
	fmt.Printf("Sequence (%v steps and %v in weight): %v\n", len(r.Sequence), r.Distance, r.Sequence)
}

type Graph struct {
	Vertices []string `json:"vertices"`
	Edges    []*Edge  `json:"edges"`
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

func (g *Graph) GetShortestWalk() (*ResultSequence, error) {
	var shortest int = math.MaxInt
	result := ResultSequence{
		Sequence: make([]string, 0),
		Distance: 0,
	}
	for _, v := range g.Vertices {
		r, err := g.WalkFrom(v)
		if err != nil {
			return nil, err
		}
		if len(r.Sequence) < shortest {
			shortest = len(r.Sequence)
			result.Distance = r.Distance
			result.Sequence = r.Sequence
		}
	}
	return &result, nil
}

func (g *Graph) resetState() {
	for _, e := range g.Edges {
		e.visitCount = 0
	}
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
			return nil, errors.New("Received an invalid edge")
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

func (g *Graph) AddEdge(from, to string) error {
	f := g.GetVertex(from)
	t := g.GetVertex(to)
	if f == "" || t == "" {
		return errors.New(fmt.Sprintf("Cannot create edge between %v and %v!", from, to))
	}
	if g.GetEdge(from, to) != nil || g.GetEdge(to, from) != nil {
		return errors.New(RepeatedEdge)
	}
	g.Edges = append(g.Edges, &Edge{
		From: f,
		To:   t,
	})
	g.Edges = append(g.Edges, &Edge{
		From: t,
		To:   f,
	})

	return nil
}

func (g *Graph) GetEdge(from, to string) *Edge {
	f := g.GetVertex(from)
	t := g.GetVertex(to)
	if f == "" || t == "" {
		return nil
	}
	for _, e := range g.Edges {
		if e.From == from && e.To == to {
			return e
		}
	}
	return nil
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

func (g *Graph) AddVertex(key string) error {
	if g.GetVertex(key) != "" {
		return errors.New(RepeatedVertex)
	}
	g.Vertices = append(g.Vertices, key)
	return nil
}

func (g *Graph) GetVertex(key string) string {
	for _, v := range g.Vertices {
		if v == key {
			return v
		}
	}
	return ""
}

func (g *Graph) Print() {
	if len(g.Vertices) != 0 {
		for _, v := range g.Vertices {
			fmt.Printf("%v: ", v)
			edges, _ := g.GetEdges(v)
			for _, e := range edges {
				fmt.Printf("%v(%v)", e.To, e.Weight)
				if e.DeadEnd {
					fmt.Printf("*")
				}
				fmt.Printf(" ")
			}
			fmt.Println()
		}
	}
}

type Edge struct {
	From       string `json:"from"`
	To         string `json:"to"`
	DeadEnd    bool   `json:"dead_end"`
	Weight     uint   `json:"weight"`
	visitCount uint
}
