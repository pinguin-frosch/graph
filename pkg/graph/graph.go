package graph

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

var (
	ErrRepeatedNode   = "node is already in the graph"
	ErrNodeNotPresent = "node is not present in the graph"
	ErrSelfEdge       = "cannot add edge between the same node"
)

type Graph struct {
	nodes map[string]Node
	edges map[string]map[string]Edge
}

type Node struct {
	id string
}

type Edge struct {
	from   Node
	to     Node
	weight float64
}

func NewGraph() *Graph {
	g := Graph{}
	g.nodes = make(map[string]Node)
	g.edges = make(map[string]map[string]Edge)
	return &g
}

func NewGraphFromFile(filename string) (*Graph, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var g Graph
	err = json.Unmarshal(bytes, &g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (g *Graph) SaveGraphToFile() error {
	bytes, err := json.Marshal(g)
	if err != nil {
		return err
	}
	err = os.Mkdir("snapshots", 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	now := time.Now().Format("02-01-06 15:04:05")
	filename := fmt.Sprintf("snapshots/%s.json", now)
	err = os.WriteFile(filename, bytes, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("Saved as %s\n", filename)
	return nil
}

func (g *Graph) AddNode(node Node) error {
	if _, ok := g.nodes[node.id]; ok {
		return errors.New(ErrRepeatedNode)
	}
	g.nodes[node.id] = node
	return nil
}

func (g *Graph) GetNode(id string) (*Node, error) {
	if node, ok := g.nodes[id]; ok {
		return &node, nil
	}
	return nil, errors.New(ErrNodeNotPresent)
}

func (g *Graph) AddEdge(edge Edge) error {
	from := edge.from
	to := edge.to

	if from.id == to.id {
		return errors.New(ErrSelfEdge)
	}

	if _, err := g.GetNode(from.id); err != nil {
		return errors.New(ErrNodeNotPresent)
	}
	if _, err := g.GetNode(to.id); err != nil {
		return errors.New(ErrNodeNotPresent)
	}

	if fromMap, ok := g.edges[from.id]; ok {
		fromMap[to.id] = edge
	} else {
		g.edges[from.id] = make(map[string]Edge)
		g.edges[from.id][to.id] = edge
	}

	if toMap, ok := g.edges[to.id]; ok {
		toMap[from.id] = edge
	} else {
		g.edges[to.id] = make(map[string]Edge)
		g.edges[to.id][from.id] = edge
	}

	return nil
}
