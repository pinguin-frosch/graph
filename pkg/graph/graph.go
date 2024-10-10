package graph

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

var (
	ErrRepeatedNode   = "node is already in the graph"
	ErrNodeNotPresent = "node is not present in the graph"
	ErrSelfEdge       = "cannot add edge between the same node"
)

type Graph struct {
	Nodes map[string]Node            `json:"nodes"`
	Edges map[string]map[string]Edge `json:"edges"`
}

type Node struct {
	Id string `json:"id"`
}

func NewNode(id string) (Node, error) {
	copyId := id
	id = strings.Trim(id, " ")
	invalidChars := make([]string, 0)
	for _, r := range id {
		if r == '_' || r == '.' {
			continue
		}
		if ('a' > r || r > 'z') && ('A' > r || r > 'Z') {
			if !slices.Contains(invalidChars, string(r)) {
				invalidChars = append(invalidChars, string(r))
			}
		}
	}
	if len(invalidChars) != 0 {
		var zero Node
		return zero, fmt.Errorf("invalid chars in id %v: %v", copyId, invalidChars)
	}
	return Node{Id: id}, nil
}

// returns all nodes in the graph in ascending order by id
func (g *Graph) GetAllNodes() []Node {
	nodes := make([]Node, len(g.Nodes))
	i := 0
	for _, node := range g.Nodes {
		nodes[i] = node
		i++
	}
	slices.SortFunc(nodes, func(a, b Node) int {
		if a.Id < b.Id {
			return -1
		} else {
			return 1
		}
	})
	return nodes
}

func (g *Graph) GetAllEdges() []Edge {
	nodes := g.GetAllNodes()
	edges := make([]Edge, 0)
	for _, node := range nodes {
		es := g.GetEdgesFrom(node)
		for _, edge := range es {
			edges = append(edges, edge)
		}
	}
	return edges
}

type Edge struct {
	Id     int  `json:"id"`
	From   Node `json:"from"`
	To     Node `json:"to"`
	Weight int  `json:"weight"`
}

// generates a key for the edge
func (e Edge) Key() string {
	return fmt.Sprintf("[%v|%v|%v]", e.Id, e.From.Id, e.To.Id)
}

// returns a new edge with the from and to fields swapped
func (e Edge) ReversedEdge() Edge {
	return Edge{e.Id, e.To, e.From, e.Weight}
}

// retuns a new edge, the id is generated when adding it to the graph
func NewEdge(from, to Node, weight int) Edge {
	return Edge{0, from, to, weight}
}

// returns all edges that are reachable from node ordered by ascending weight
func (g *Graph) GetEdgesFrom(node Node) []Edge {
	edgesMap := g.Edges[node.Id]
	edges := make([]Edge, len(edgesMap))
	i := 0
	for _, edge := range edgesMap {
		edges[i] = edge
		i++
	}
	slices.SortFunc(edges, func(a, b Edge) int {
		if a.Weight < b.Weight {
			return -1
		} else if a.Weight > b.Weight {
			return 1
		} else {
			return 0
		}
	})
	return edges
}

// returns all nodes that are reachable from node
func (g *Graph) GetNodesFrom(node Node) []Node {
	edges := g.GetEdgesFrom(node)
	nodes := make([]Node, 0, len(edges))
	for _, edge := range edges {
		nodes = append(nodes, edge.To)
	}
	return nodes
}

func (g *Graph) GetEdge(from, to Node) Edge {
	return g.Edges[from.Id][to.Id]
}

func NewGraph() Graph {
	g := Graph{}
	g.Nodes = make(map[string]Node)
	g.Edges = make(map[string]map[string]Edge)
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
	if _, ok := g.Nodes[node.Id]; ok {
		return errors.New(ErrRepeatedNode)
	}
	g.Nodes[node.Id] = node
	return nil
}

func (g *Graph) GetNode(id string) (Node, error) {
	if node, ok := g.Nodes[id]; ok {
		return node, nil
	}
	return Node{}, errors.New(ErrNodeNotPresent)
}

func (g *Graph) AddEdge(edge Edge) error {
	from := edge.From
	to := edge.To

	if from.Id == to.Id {
		return errors.New(ErrSelfEdge)
	}

	if _, err := g.GetNode(from.Id); err != nil {
		return errors.New(ErrNodeNotPresent)
	}
	if _, err := g.GetNode(to.Id); err != nil {
		return errors.New(ErrNodeNotPresent)
	}

	if fromMap, ok := g.Edges[from.Id]; ok {
		fromMap[to.Id] = edge
	} else {
		g.Edges[from.Id] = make(map[string]Edge)
		g.Edges[from.Id][to.Id] = edge
	}

	if toMap, ok := g.Edges[to.Id]; ok {
		toMap[from.Id] = edge.ReversedEdge()
	} else {
		g.Edges[to.Id] = make(map[string]Edge)
		g.Edges[to.Id][from.Id] = edge.ReversedEdge()
	}

	return nil
}

func (g *Graph) RemoveEdge(from, to Node) {
	delete(g.Edges[from.Id], to.Id)
	delete(g.Edges[to.Id], from.Id)
}

func (g *Graph) RemoveNode(node Node) {
	edges := g.GetEdgesFrom(node)
	for _, edge := range edges {
		g.RemoveEdge(edge.From, edge.To)
	}
	delete(g.Nodes, node.Id)
}

func (g *Graph) Print() {
	nodes := g.GetAllNodes()
	for _, node := range nodes {
		fmt.Printf("%s: ", node.Id)
		edges := g.GetEdgesFrom(node)
		for _, edge := range edges {
			fmt.Printf("%s(%d) ", edge.To.Id, edge.Weight)
		}
		fmt.Println()
	}
}
