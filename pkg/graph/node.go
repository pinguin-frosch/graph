package graph

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var (
	ErrRepeatedNode   = "node is already in the graph"
	ErrNodeNotPresent = "node is not present in the graph"
)

type Node struct {
	Id string `json:"id"`
}

func (g *Graph) Degree(node Node) int {
	edges := g.GetEdges(node)
	return len(edges)
}

// returns a new node, or an error if the id is invalid
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
	slices.SortFunc(nodes, sortNodesById)
	return nodes
}

// returns all odd nodes in the graph in ascending order by id
func (g *Graph) GetAllOddNodes() []Node {
	oddNodes := make([]Node, 0)
	for _, node := range g.GetAllNodes() {
		if g.Degree(node)%2 == 1 {
			oddNodes = append(oddNodes, node)
		}
	}
	return oddNodes
}

// returns all deadend nodes in the graph in ascending order by id
func (g *Graph) GetAllDeadendNodes() []Node {
	deadendNodes := make([]Node, 0)
	for _, node := range g.GetAllNodes() {
		if g.Degree(node) == 1 {
			deadendNodes = append(deadendNodes, node)
		}
	}
	return deadendNodes
}

var sortNodesById func(a, b Node) int = func(a, b Node) int {
	if a.Id < b.Id {
		return -1
	} else {
		return 1
	}
}

// returns all nodes that are reachable from node
func (g *Graph) GetNodes(node Node) []Node {
	edges := g.GetEdges(node)
	nodes := make([]Node, 0, len(edges))
	for _, edge := range edges {
		nodes = append(nodes, edge.To)
	}
	return nodes
}

// adds a node to the graph
func (g *Graph) AddNode(node Node) error {
	if _, ok := g.Nodes[node.Id]; ok {
		return errors.New(ErrRepeatedNode)
	}
	g.Nodes[node.Id] = node
	return nil
}

// gets node with given id from the graph, or an error if node does not exist
func (g *Graph) GetNode(id string) (Node, error) {
	if node, ok := g.Nodes[id]; ok {
		return node, nil
	}
	return Node{}, errors.New(ErrNodeNotPresent)
}

// removes a node from the graph and all its edges
func (g *Graph) RemoveNode(node Node) {
	edges := g.GetEdges(node)
	for _, edge := range edges {
		g.RemoveEdges(edge.From, edge.To)
	}
	delete(g.Nodes, node.Id)
}
