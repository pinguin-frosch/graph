package graph

import (
	"errors"
	"fmt"
	"math"
)

type Graph struct {
	vertices     []*Vertex
	totalEdges   uint
	walkedEdges  uint
	walkSequence []string
}

func (g *Graph) WalkFrom(key string) error {
	startVertex := g.GetVertex(key)
	if startVertex == nil {
		return errors.New("The start vertex does not exist in the graph!")
	}
	g.walkSequence = append(g.walkSequence, startVertex.key)

	vertex := startVertex
	for g.totalEdges != g.walkedEdges {
		fmt.Printf("Total: %v - Current: %v\n", g.totalEdges, g.walkedEdges)

		// Encontrar la menor cantidad de veces que ha sido usada una arista
		var min uint = math.MaxUint
		for _, e := range vertex.edges {
			if e.visitCount < min {
				min = e.visitCount
			}
		}

		// Recorrer una de las menores
		for _, e := range vertex.edges {
			if e.visitCount != min {
				continue
			}
			// Ignorar si existe ya existe la arista inversa
			if g.ExistsEdge(e.to.key, vertex.key) {
			}
			if e.visitCount == 0 {
				g.walkedEdges++
			}
			e.visitCount++

			// Cambiar el vertice al siguiente
			g.walkSequence = append(g.walkSequence, e.to.key)
			vertex = e.to
			break
		}
	}

	fmt.Println(g.walkSequence)

	return nil
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

func (g *Graph) ExistsEdge(from, to string) bool {
	f := g.GetVertex(from)
	t := g.GetVertex(to)
	if f == nil || t == nil {
		return false
	}
	if !containsEdge(f.edges, to) {
		return false
	}
	return true
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
	// No aumentar la cantidad de aristas si es la misma al revés
	if !g.ExistsEdge(to, from) {
		g.totalEdges++
	}
	f.edges = append(f.edges, &Edge{to: t})
	return nil
}

func (g *Graph) Print() {
	fmt.Printf("The graph has %v vertices and %v edges.\n", len(g.vertices), g.totalEdges)
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
	visitCount uint
}
