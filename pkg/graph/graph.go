package graph

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type Graph struct {
	Nodes map[string]Node                    `json:"nodes"`
	Edges map[string]map[string]map[int]Edge `json:"edges"`
}

// initializes an empty graph
func NewGraph() Graph {
	g := Graph{}
	g.Nodes = make(map[string]Node)
	g.Edges = make(map[string]map[string]map[int]Edge)
	return g
}

// initializes a graph from given file
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

// saves graph state to a file
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

// prints all nodes and edges in the graph organized
func (g *Graph) Print() {
	nodes := g.GetAllNodes()
	for _, node := range nodes {
		fmt.Printf("%s: ", node.Id)
		edges := g.GetEdges(node)
		for _, edge := range edges {
			fmt.Printf("%s[%d](%d) ", edge.To.Id, edge.Id, edge.Weight)
		}
		fmt.Println()
	}
}
