package traverse_test

import (
	"graph/pkg/graph"
	"graph/pkg/traverse"
	"testing"
)

func TestDefaultCountTotalEdges(t *testing.T) {
	g := graph.NewGraph()
	nodes := make([]graph.Node, 0)
	ids := []string{"a", "b", "c", "d"}
	for _, id := range ids {
		node, _ := graph.NewNode(id)
		nodes = append(nodes, node)
		_ = g.AddNode(node)
	}
	expected := 0
	for i := 1; i < len(nodes); i++ {
		from, to := nodes[i-1], nodes[i]
		e := graph.NewEdge(from, to, 5)
		_ = g.AddEdge(e)
		expected++
	}
	d := traverse.NewDefault()
	total := d.CountTotalEdges(g)
	if total != expected {
		t.Errorf("CountTotalEdges(g) = %d; expected %d", total, expected)
	}
}
