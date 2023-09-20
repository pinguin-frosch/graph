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

type Graph struct {
	Vertices  []string `json:"vertices"`
	Edges     []*Edge  `json:"edges"`
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

func (g *Graph) GetShortestWalk() ([]string, error) {
	var shortest int = math.MaxInt
	var shortestSequence []string
	for _, v := range g.Vertices {
		sequence, err := g.WalkFrom(v)
		if err != nil {
			return nil, err
		}
		if len(sequence) < shortest {
			shortest = len(sequence)
			shortestSequence = sequence
		}
	}
	return shortestSequence, nil
}

func (g *Graph) resetState() {
	for _, e := range g.Edges {
		e.visitCount = 0
	}
}

func (g *Graph) WalkFrom(from string) ([]string, error) {
	totalEdges := g.GetTotalEdges()
	var usedEdges uint = 0
	g.resetState()
	sequence := make([]string, 0)
	vertex := g.GetVertex(from)
	if vertex == "" {
		return nil, errors.New("The starting vertex does not exist!")
	}
	sequence = append(sequence, vertex)

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

		sequence = append(sequence, edge.To)
		vertex = edge.To
	}

	return sequence, nil
}

func (g *Graph) AddEdge(from, to string) error {
	f := g.GetVertex(from)
	t := g.GetVertex(to)
	if f == "" || t == "" {
		return errors.New(fmt.Sprintf("Cannot create edge between %v and %v!\n", from, to))
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

// Returns all the edges from a certain vertex, the least any edge has been
// used and wether there's a dead end edge. This is not supposed to be used
// directly.
func (g *Graph) GetEdges(key string) ([]*Edge, uint, uint, int) {
	var minimum uint = math.MaxUint
	var minimumDeadEnd uint = math.MaxUint
	deadEndCount := 0
	edges := make([]*Edge, 0)

	// Compute all the required information
	for _, e := range g.Edges {
		if e.From == key {
			switch e.DeadEnd {
			case true:
				deadEndCount++
				if e.visitCount < minimumDeadEnd {
					minimumDeadEnd = e.visitCount
				}
			case false:
				if e.visitCount < minimum {
					minimum = e.visitCount
				}
			}
			edges = append(edges, e)
		}
	}

	return edges, minimum, minimumDeadEnd, deadEndCount
}

// Returns the next edge to use in the search algorithm and wether it was the only one found
func (g *Graph) GetNextEdge(key string) (*Edge, bool) {
	edges, minimum, minimumDeadEnd, deadEndCount := g.GetEdges(key)
	unique := len(edges) == 1

	if deadEndCount == 1 {
		for _, e := range edges {
			if !e.DeadEnd {
				continue
			}
			if e.visitCount < 2 {
				return e, unique
			}
		}
	} else if deadEndCount > 1 {
		for _, e := range edges {
			if !e.DeadEnd {
				continue
			}
			if e.visitCount > minimumDeadEnd {
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
		if e.visitCount > minimum {
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
			edges, _, _, _ := g.GetEdges(v)
			for _, e := range edges {
				fmt.Printf("%v", e.To)
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
