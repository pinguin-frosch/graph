package graph

import (
	"errors"
	"fmt"
	"math"
)

type Graph struct {
	vertices     []*Vertex
	edges        []*Edge
	walkSequence []string
}

const DEBUG = true

func (g *Graph) WalkFrom(from string) error {
	vertex := g.GetVertex(from)
	if vertex == nil {
		return errors.New("The starting vertex does not exist!")
	}
	g.walkSequence = append(g.walkSequence, vertex.key)

	for !g.IsTraversed() {

		if DEBUG {
			fmt.Printf("Sequence: %v\n", g.walkSequence)
		}

		// Find all the edges that start from the current vertex
		edges, minimum := g.GetEdges(vertex.key)

		if len(edges) == 0 {
			return errors.New(fmt.Sprintf("There's nowhere to go from vertex %v\n", vertex.key))
		}

		for _, edge := range edges {
			if edge.visitCount > minimum {
				continue
			}

			if DEBUG {
				fmt.Printf("The selected vertex is %v\n", edge.to.key)
			}

			// Get the other edge
			otherEdge := g.GetEdge(edge.to.key, edge.from.key)
			if otherEdge == nil {
				return errors.New(fmt.Sprintf("Couldn't find the other edge, from %v to %v", edge.to.key, edge.from.key))
			}

			// Mark both edges as used once
			otherEdge.visitCount++
			edge.visitCount++

			g.walkSequence = append(g.walkSequence, edge.to.key)
			vertex = edge.to
			break
		}

	}

	fmt.Printf("One walking sequence is %v\n", g.walkSequence)
	return nil
}

func (g *Graph) IsTraversed() bool {
	ok := true
	for _, e := range g.edges {
		if DEBUG {
			fmt.Printf("%v<->%v : %v\n", e.from.key, e.to.key, e.visitCount > 0)
		}
		if e.visitCount == 0 {
			if DEBUG {
				ok = false
			} else {
				return false
			}
		}
	}
	return ok
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

	if DEBUG {
		fmt.Printf("The neighbors for %v are: ", key)
		for _, e := range edges {
			fmt.Printf("%v ", e.to.key)
		}
		fmt.Println()
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
