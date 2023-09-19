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
		// Find all the edges that start from the current vertex
		edges, minimum := g.GetEdges(vertex.key)

		if len(edges) == 0 {
			return nil, errors.New(fmt.Sprintf("There's nowhere to go from vertex %v\n", vertex.key))
		}

		for _, edge := range edges {
			if edge.visitCount > minimum {
				continue
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
			break
		}

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

// Returns all the edges from a certain vertex and the least any edge has been used
func (g *Graph) GetEdges(key string) ([]*Edge, uint) {
	var m uint = math.MaxUint
	edges := make([]*Edge, 0)

	// Get the minimum any edge has been used
	for _, e := range g.edges {
		if e.from.key == key {
			if e.visitCount < m {
				m = e.visitCount
			}
			edges = append(edges, e)
		}
	}

	return edges, m
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
			edges, _ := g.GetEdges(v.key)
			for _, e := range edges {
				fmt.Printf("%v ", e.to.key)
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
	visitCount uint
}
