package graph

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
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
	Nodes map[string]Node                    `json:"nodes"`
	Edges map[string]map[string]map[int]Edge `json:"edges"`
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
	edges := make([]Edge, 0)
	edgesMap := g.Edges[node.Id]
	for _, edgeSubMap := range edgesMap {
		for _, edge := range edgeSubMap {
			edges = append(edges, edge)
		}
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
	// FIXME: return shortest edge and error when not found
	edges := g.Edges[from.Id][to.Id]
	for _, edge := range edges {
		return edge
	}
	return Edge{}
}

func NewGraph() Graph {
	g := Graph{}
	g.Nodes = make(map[string]Node)
	g.Edges = make(map[string]map[string]map[int]Edge)
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

	// create the map for the from node
	if _, ok := g.Edges[from.Id]; !ok {
		g.Edges[from.Id] = make(map[string]map[int]Edge)
	}
	// create the map for the to node under the from node
	if _, ok := g.Edges[from.Id][to.Id]; !ok {
		g.Edges[from.Id][to.Id] = make(map[int]Edge)
	}

	// create the map for the to node
	if _, ok := g.Edges[to.Id]; !ok {
		g.Edges[to.Id] = make(map[string]map[int]Edge)
	}
	// create the map for the from node under the to node
	if _, ok := g.Edges[to.Id][from.Id]; !ok {
		g.Edges[to.Id][from.Id] = make(map[int]Edge)
	}

	// find a valid id and add the edge
	i := 0
	for {
		if i >= math.MaxInt {
			return fmt.Errorf("couldn't generate a valid id for the edge")
		}
		if _, okFrom := g.Edges[from.Id][to.Id][i]; !okFrom {
			if _, okTo := g.Edges[to.Id][from.Id][i]; !okTo {
				edge.Id = i
				g.Edges[from.Id][to.Id][i] = edge
				g.Edges[to.Id][from.Id][i] = edge.ReversedEdge()
				break
			}
		}
		i++
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
			fmt.Printf("%s[%d](%d) ", edge.To.Id, edge.Id, edge.Weight)
		}
		fmt.Println()
	}
}
