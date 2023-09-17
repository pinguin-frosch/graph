package graph

import (
	"errors"
	"fmt"
)

type Graph struct {
	vertices     []*Vertex
	edges        []*Edge
	walkedEdges  uint
	walkSequence []string
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
