package graph

import (
	"errors"
	"fmt"
)

const (
	RepeatedVertex = "vertex already exists"
	RepeatedEdge   = "edge already exists"
)

type Graph struct {
	Vertices []string `json:"vertices"`
	Edges    []*Edge  `json:"edges"`
}

func (g *Graph) AddEdge(from, to string) error {
	f := g.GetVertex(from)
	t := g.GetVertex(to)
	if f == "" || t == "" {
		return fmt.Errorf("cannot create edge between %v and %v", from, to)
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

func (g *Graph) AddVertex(key string) error {
	if key == "" {
		return errors.New("invalid key for vertex")
	}
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
			edges, _ := g.GetEdges(v)
			for _, e := range edges {
				fmt.Printf("%v(%v)", e.To, e.Weight)
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
