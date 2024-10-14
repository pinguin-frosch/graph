package aograph

import (
	"fmt"
	"math"
	"slices"
	"strings"
)

type Graph struct {
	Nodes map[string]*Node
}

type Node struct {
	Id    string
	Value float64
	And   GroupNode
	Or    GroupNode
}

type GroupNode struct {
	Value float64
	Nodes map[string]*Node
}

func NewNode(id string, value float64) Node {
	node := Node{Id: id, Value: value}
	node.And.Nodes = make(map[string]*Node)
	node.Or.Nodes = make(map[string]*Node)
	return node
}

func NewGraph() Graph {
	graph := Graph{}
	graph.Nodes = make(map[string]*Node)
	return graph
}

func (g *Graph) AddNode(node Node) bool {
	if _, ok := g.Nodes[node.Id]; ok {
		return false
	}
	g.Nodes[node.Id] = &node
	return true
}

func (g *Graph) Print() {
	orderedNodeKeys := make([]string, 0)
	for _, node := range g.Nodes {
		orderedNodeKeys = append(orderedNodeKeys, node.Id)
	}
	slices.Sort(orderedNodeKeys)

	for _, key := range orderedNodeKeys {
		node := g.Nodes[key]
		fmt.Printf("Node %s, Value %f, ", node.Id, node.Value)
		fmt.Printf("Or: { Value: %f [", node.Or.Value)
		i := 0
		for _, child := range node.Or.Nodes {
			fmt.Printf("%s", child.Id)
			if i != len(node.Or.Nodes)-1 {
				fmt.Printf(" ")
			}
			i++
		}
		fmt.Printf("] }, ")
		fmt.Printf("And: { Value: %f [", node.And.Value)
		i = 0
		for _, child := range node.And.Nodes {
			fmt.Printf("%s", child.Id)
			if i != len(node.And.Nodes)-1 {
				fmt.Printf(" ")
			}
			i++
		}
		fmt.Printf("] }")
		fmt.Printf("\n")
	}
}

func (g *Graph) AddOrConnection(parentId string, childIds ...string) bool {
	for _, childId := range childIds {
		if _, ok := g.Nodes[parentId].Or.Nodes[childId]; ok {
			return false
		}
		g.Nodes[parentId].Or.Nodes[childId] = g.Nodes[childId]
	}
	return true
}

func (g *Graph) AddAndConnection(parentId string, childIds ...string) bool {
	for _, childId := range childIds {
		if _, ok := g.Nodes[parentId].And.Nodes[childId]; ok {
			return false
		}
		g.Nodes[parentId].And.Nodes[childId] = g.Nodes[childId]
	}
	return true
}

func (g *Graph) CostPerNode(node *Node) map[string]float64 {
	cost := make(map[string]float64)
	if len(node.And.Nodes) > 0 {
		childs := node.And.Nodes
		childsIds := make([]string, 0, len(childs))
		for _, child := range childs {
			childsIds = append(childsIds, child.Id)
		}
		stringKey := strings.Join(childsIds, " AND ")
		pathSum := float64(0)
		for _, child := range childs {
			pathSum += child.Value + 1 // FIXME: use a non fixed value for the edge weight
		}
		cost[stringKey] = pathSum
		node.And.Value = pathSum
	}
	if len(node.Or.Nodes) > 0 {
		childs := node.Or.Nodes
		childsIds := make([]string, 0, len(childs))
		for _, child := range childs {
			childsIds = append(childsIds, child.Id)
		}
		stringKey := strings.Join(childsIds, " OR ")
		pathSums := make([]float64, 0, len(node.Or.Nodes))
		for _, child := range childs {
			pathSums = append(pathSums, child.Value+1)
		}
		minimumCost := math.MaxFloat64
		for _, pathSum := range pathSums {
			if pathSum < minimumCost {
				minimumCost = pathSum
			}
		}
		cost[stringKey] = minimumCost
		node.Or.Value = minimumCost
	}
	return cost
}

func (g *Graph) UpdateCost() {
	// FIXME: make it so that this works using dfs
	orderedNodeKeys := make([]string, 0)
	for _, node := range g.Nodes {
		orderedNodeKeys = append(orderedNodeKeys, node.Id)
	}
	slices.Sort(orderedNodeKeys)
	slices.Reverse(orderedNodeKeys)

	for _, key := range orderedNodeKeys {
		node := g.Nodes[key]
		// Skip bottom nodes
		if len(node.And.Nodes) == 0 && len(node.Or.Nodes) == 0 {
			continue
		}

		costPerNode := g.CostPerNode(node)

		// Print results
		fmt.Printf("%s : ", node.Id)
		node.PrintConnections()
		fmt.Printf(" >>> ")
		fmt.Printf("{")
		for key, cost := range costPerNode {
			fmt.Printf("\"%s\": %f, ", key, cost)
		}
		fmt.Printf("}\n")

		// Update cost in the node
		minimunCost := math.MaxFloat64
		for _, cost := range costPerNode {
			if cost < minimunCost {
				minimunCost = cost
			}
		}
		// if minimunCost < node.Value {
		node.Value = minimunCost
		// }
	}
}

func (g *Graph) UpdateNodeAndValue(id string, value float64) {
	node := g.Nodes[id]
	node.And.Value = value
	g.Nodes[id] = node
}

func (n Node) PrintConnections() {
	orPrinted := false
	if len(n.Or.Nodes) > 0 {
		orPrinted = true
		fmt.Printf("OR: [")
		i := 0
		for _, child := range n.Or.Nodes {
			fmt.Printf("%s", child.Id)
			if i != len(n.Or.Nodes)-1 {
				fmt.Printf(" ")
			}
			i++
		}
		fmt.Printf("]")
	}
	if len(n.And.Nodes) > 0 {
		if orPrinted {
			fmt.Printf(" ")
		}
		fmt.Printf("AND: [")
		i := 0
		for _, child := range n.And.Nodes {
			fmt.Printf("%s", child.Id)
			if i != len(n.And.Nodes)-1 {
				fmt.Printf(" ")
			}
			i++
		}
		fmt.Printf("]")
	}
}
