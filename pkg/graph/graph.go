package graph

import (
	"errors"
	"fmt"
)

type Graph struct {
	vertices []*Vertex
}

func (g *Graph) AddVertex(key string) error {
	if containsVertex(g.vertices, key) {
		return errors.New("That vertex already exists!")
	}
	v := &Vertex{
		key: key,
	}
	g.vertices = append(g.vertices, v)
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

func (g *Graph) AddEdge(from, to string) error {
	f := g.GetVertex(from)
	t := g.GetVertex(to)
	if f == nil || t == nil {
		return errors.New(fmt.Sprintf("Invalid edge between %v and %v!", from, to))
	}
	if containsEdge(f.edges, to) {
		return errors.New("That edge already exists!")
	}
	f.edges = append(f.edges, &Edge{to: t})
	return nil
}

func (g *Graph) Print() {
	for _, v := range g.vertices {
		fmt.Printf("%v: ", v.key)
		v.Print()
		fmt.Println()
	}
}

func containsVertex(v []*Vertex, key string) bool {
	for _, v := range v {
		if v.key == key {
			return true
		}
	}
	return false
}

type Vertex struct {
	key   string
	edges []*Edge
}

func (v *Vertex) Print() {
	for _, e := range v.edges {
		fmt.Printf("%v ", e.to.key)
	}
}

func containsEdge(e []*Edge, key string) bool {
	for _, e := range e {
		if e.to.key == key {
			return true
		}
	}
	return false
}

type Edge struct {
	to         *Vertex
	visited    bool
	visitCount uint
}
