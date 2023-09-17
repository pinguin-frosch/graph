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
		for _, e := range g.edges {
			fmt.Printf("%v <-> %v ", e.from.key, e.to.key)
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
