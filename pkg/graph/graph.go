package graph

import (
	"errors"
	"fmt"
	"math"
)

type Graph struct {
	vertices   []*Vertex
	edges      []*Edge
	usedEdges  uint
	totalEdges uint
}

func (g *Graph) WalkFrom(from string) ([]string, error) {
	sequence := make([]string, 0)
	vertex := g.GetVertex(from)
	if vertex == nil {
		return nil, errors.New("The starting vertex does not exist!")
	}
	sequence = append(sequence, vertex.key)

	for !g.IsTraversed() {
		// Get the next edge to use
		edge := g.GetNextEdge(vertex.key)
		if edge == nil {
			return nil, errors.New("Received an invalid edge")
		}

		// Get the other edge
		otherEdge := g.GetEdge(edge.to.key, edge.from.key)
		if otherEdge == nil {
			return nil, errors.New(fmt.Sprintf("Couldn't find the other edge, from %v to %v", edge.to.key, edge.from.key))
		}

		// Only count unused edges
		if edge.visitCount == 0 {
			g.usedEdges++
		}

		// Mark both edges as used once
		otherEdge.visitCount++
		edge.visitCount++

		sequence = append(sequence, edge.to.key)
		vertex = edge.to
	}

	return sequence, nil
}

func (g *Graph) IsTraversed() bool {
	return g.usedEdges >= g.totalEdges
}

func (g *Graph) AddEdge(from, to string) error {
	f := g.GetVertex(from)
	t := g.GetVertex(to)
	if f == nil || t == nil {
		return errors.New(fmt.Sprintf("Cannot create edge between %v and %v!\n", from, to))
	}
	if g.GetEdge(from, to) != nil || g.GetEdge(to, from) != nil {
		return errors.New("That edge already exists!")
	}
	g.edges = append(g.edges, &Edge{
		from: f,
		to:   t,
	})
	g.edges = append(g.edges, &Edge{
		from: t,
		to:   f,
	})

	// Since both edges count as one we only increment by one
	g.totalEdges++

	return nil
}

func (g *Graph) GetEdge(from, to string) *Edge {
	f := g.GetVertex(from)
	t := g.GetVertex(to)
	if f == nil || t == nil {
		return nil
	}
	for _, e := range g.edges {
		if e.from.key == from && e.to.key == to {
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
	for _, e := range g.edges {
		if e.from.key == key {
			switch e.deadEnd {
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

// Returns the next edge to use in the search algorithm
func (g *Graph) GetNextEdge(key string) *Edge {
	edges, minimum, minimumDeadEnd, deadEndCount := g.GetEdges(key)

	if deadEndCount == 1 {
		for _, e := range edges {
			if !e.deadEnd {
				continue
			}
			if e.visitCount < 2 {
				return e
			}
		}
	} else if deadEndCount > 1 {
		for _, e := range edges {
			if !e.deadEnd {
				continue
			}
			if e.visitCount > minimumDeadEnd {
				continue
			}
			if e.visitCount < 2 {
				return e
			}
		}
	}
	for _, e := range edges {
		if e.deadEnd {
			continue
		}
		if e.visitCount > minimum {
			continue
		}
		return e
	}

	return nil
}

func (g *Graph) AddVertex(key string) error {
	if g.GetVertex(key) != nil {
		return errors.New("That vertex already exists!")
	}
	g.vertices = append(g.vertices, &Vertex{key: key})
	return nil
}

func (g *Graph) GetVertex(key string) *Vertex {
	for _, v := range g.vertices {
		if v.key == key {
			return v
		}
	}
	return nil
}

func (g *Graph) Print() {
	if len(g.vertices) != 0 {
		for _, v := range g.vertices {
			fmt.Printf("%v: ", v.key)
			edges, _, _, _ := g.GetEdges(v.key)
			for _, e := range edges {
				fmt.Printf("%v", e.to.key)
				if e.deadEnd {
					fmt.Printf("*")
				}
				fmt.Printf(" ")
			}
			fmt.Println()
		}
	}
}

type Vertex struct {
	key string
}

type Edge struct {
	from       *Vertex
	to         *Vertex
	deadEnd    bool
	visitCount uint
}
