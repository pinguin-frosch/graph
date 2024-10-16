package graph

import (
	"errors"
	"fmt"
	"math"
	"slices"
)

var (
	ErrSelfEdge  = "cannot add edge between the same node"
	ErrMaxIdUsed = "cannot generate a valid id for the edge, max id used"
)

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

// returns all the edges present in the graph
func (g *Graph) GetAllEdges() []Edge {
	nodes := g.GetAllNodes()
	edges := make([]Edge, 0)
	for _, node := range nodes {
		es := g.GetEdges(node)
		for _, edge := range es {
			edges = append(edges, edge)
		}
	}
	return edges
}

// returns all edges that are reachable from node ordered by ascending weight
func (g *Graph) GetEdges(node Node) []Edge {
	edges := make([]Edge, 0)
	edgesMap := g.Edges[node.Id]
	for _, edgeSubMap := range edgesMap {
		for _, edge := range edgeSubMap {
			edges = append(edges, edge)
		}
	}
	slices.SortFunc(edges, sortEdgesByWeight)
	return edges
}

var sortEdgesByWeight func(a, b Edge) int = func(a, b Edge) int {
	if a.Weight < b.Weight {
		return -1
	} else if a.Weight > b.Weight {
		return 1
	} else {
		return 0
	}
}

// returns the shortest edge between the from and to nodes, indicates if it was found
func (g *Graph) GetShortestEdge(from, to Node) (Edge, bool) {
	var zero Edge
	if len(g.GetEdges(from)) == 0 {
		return zero, false
	}
	minWeight := math.MaxInt
	edges := g.Edges[from.Id][to.Id]
	for _, edge := range edges {
		if edge.Weight < minWeight {
			minWeight = edge.Weight
		}
	}
	for _, edge := range edges {
		if edge.Weight == minWeight {
			return edge, true
		}
	}
	return zero, false
}

// adds an edge to the graph
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
			return errors.New(ErrMaxIdUsed)
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

// Removes first edge found between from and to nodes with given weight value
func (g *Graph) RemoveEdgeWithWeight(from, to Node, weight int) {
	for id, edge := range g.Edges[from.Id][to.Id] {
		if edge.Weight == weight {
			delete(g.Edges[from.Id][to.Id], id)
			delete(g.Edges[to.Id][from.Id], id)
			return
		}
	}
}

// Removes all edges between from and to nodes
func (g *Graph) RemoveEdges(from, to Node) {
	delete(g.Edges[from.Id], to.Id)
	delete(g.Edges[to.Id], from.Id)
}
