package graph

import (
	"errors"
	"fmt"
	"math"
)

type Graph struct {
	vertices     []*Vertex
	edges        []*Edge
	walkedEdges  uint
	walkSequence []string
}

func (g *Graph) WalkFrom(from string) error {
	vertex := g.GetVertex(from)
	if vertex == nil {
		return errors.New("The starting vertex does not exist!")
	}
	g.walkSequence = append(g.walkSequence, vertex.key)

	for g.walkedEdges != uint(len(g.edges)) {
		// Find the from that start from vertex
		from, mFrom := g.GetEdgesFrom(vertex.key)
		to, mTo := g.GetEdgesTo(vertex.key)

		if len(from) == 0 && len(to) == 0 {
			return errors.New(fmt.Sprintf("There's nowhere to go from vertex %v\n", vertex.key))
		}

		if len(from) != 0 && mFrom < mTo {
			e := g.CheckEdges(from, mFrom)
			if e == nil {
				fmt.Println(e)
				return errors.New("Received an invalid edge!")
			}
			g.walkSequence = append(g.walkSequence, e.to.key)
			vertex = e.to
			continue
		} else if len(to) != 0 && mTo <= mFrom {
			e := g.CheckEdges(to, mTo)
			if e == nil {
				fmt.Println(e)
				return errors.New("Received an invalid edge!")
			}
			g.walkSequence = append(g.walkSequence, e.from.key)
			vertex = e.from
			continue
		}

	}

	fmt.Printf("One walking sequence is %v\n", g.walkSequence)
	return nil
}

func (g *Graph) CheckEdges(edges []*Edge, minimum uint) *Edge {
	// Walk to any of the least used
	for _, e := range edges {
		if e.visitCount != minimum {
			continue
		}
		if e.visitCount == 0 {
			g.walkedEdges++
		}
		e.visitCount++
		return e
	}

	return nil
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

// Returns all the edges from a certain vertex and the minimum value found
func (g *Graph) GetEdgesFrom(from string) ([]*Edge, uint) {
	var m uint = math.MaxUint
	edges := make([]*Edge, 0)
	for _, e := range g.edges {
		if e.from.key == from {
			if e.visitCount < m {
				m = e.visitCount
			}
			edges = append(edges, e)
		}
	}
	return edges, m
}

// Returns all the edges to a certain vertex and the minimum value found
func (g *Graph) GetEdgesTo(to string) ([]*Edge, uint) {
	var m uint = math.MaxUint
	edges := make([]*Edge, 0)
	for _, e := range g.edges {
		if e.to.key == to {
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
		fmt.Println("Vertices")
		for _, v := range g.vertices {
			fmt.Printf("%v ", v.key)
		}
		fmt.Println()
	}
	if len(g.edges) != 0 {
		fmt.Println("\nEdges")
		fmt.Printf("| ")
		for _, e := range g.edges {
			fmt.Printf("%v<->%v | ", e.from.key, e.to.key)
		}
		fmt.Println()
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
