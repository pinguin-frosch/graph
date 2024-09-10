package graph

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
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

func NewNode(id string) Node {
	return Node{id: id}
}

// returns all nodes in the graph in ascending order by id
func (g *Graph) Nodes() []Node {
	nodes := make([]Node, len(g.nodes))
	i := 0
	for _, node := range g.nodes {
		nodes[i] = node
		i++
	}
	slices.SortFunc(nodes, func(a, b Node) int {
		if a.id < b.id {
			return -1
		} else {
			return 1
		}
	})
	return nodes
}

type Edge struct {
	from   Node
	to     Node
	weight int
}

func (e Edge) From() Node {
	return e.from
}

func (e Edge) To() Node {
	return e.to
}

// returns a new edge with the from and to fields swapped
func (e Edge) ReversedEdge() Edge {
	return Edge{e.to, e.from, e.weight}
}

func NewEdge(from, to Node, weight int) Edge {
	return Edge{from, to, weight}
}

// returns all edges that are reachable from node ordered by ascending weight
func (g *Graph) EdgesFrom(node Node) []Edge {
	edgesMap := g.edges[node.id]
	edges := make([]Edge, len(edgesMap))
	i := 0
	for _, edge := range edgesMap {
		edges[i] = edge
		i++
	}
	slices.SortFunc(edges, func(a, b Edge) int {
		if a.weight < b.weight {
			return -1
		} else if a.weight > b.weight {
			return 1
		} else {
			return 0
		}
	})
	return edges
}

func NewGraph() Graph {
	g := Graph{}
	g.nodes = make(map[string]Node)
	g.edges = make(map[string]map[string]Edge)
	return g
}

func NewGraphFromFile(filename string) (Graph, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return Graph{}, err
	}
	var g Graph
	err = json.Unmarshal(bytes, &g)
	if err != nil {
		return Graph{}, err
	}
	return g, nil
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

func (g *Graph) GetNode(id string) (Node, error) {
	if node, ok := g.nodes[id]; ok {
		return node, nil
	}
	return Node{}, errors.New(ErrNodeNotPresent)
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
		toMap[from.id] = edge.ReversedEdge()
	} else {
		g.edges[to.id] = make(map[string]Edge)
		g.edges[to.id][from.id] = edge.ReversedEdge()
	}

	return nil
}

func (g *Graph) Print() {
	nodes := g.Nodes()
	for _, node := range nodes {
		fmt.Printf("%s: ", node.id)
		edges := g.EdgesFrom(node)
		for _, edge := range edges {
			fmt.Printf("%s(%d) ", edge.to.id, edge.weight)
		}
		fmt.Println()
	}
}
