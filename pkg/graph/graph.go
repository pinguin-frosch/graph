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

func (g *Graph) Print() {
	for _, v := range g.vertices {
		fmt.Println(v)
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
	key string
}
