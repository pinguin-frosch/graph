package graph_test

import (
	"graph/pkg/graph"
	"testing"
)

func TestNewNode(t *testing.T) {
	invalidIds := []string{"normal1", "thisisfine()", "...777", "_._asdfñ"}
	for _, invalidId := range invalidIds {
		_, err := graph.NewNode(invalidId)
		if err == nil {
			t.Fatalf(`NewNode("%v") should return an error, %s is invalid`, invalidId, invalidId)
		}
	}
	validIds := []string{"ABC__jj", "jklIIO.aa._z", "a_b_c_d", "a.z.i.n"}
	for _, validId := range validIds {
		_, err := graph.NewNode(validId)
		if err != nil {
			t.Fatalf(`NewNode("%v") should'nt return an error, %s is valid`, validId, validId)
		}
	}
}

func TestAddNode(t *testing.T) {
	g := graph.NewGraph()
	nodeId := "a"
	a, err := graph.NewNode(nodeId)
	if err != nil {
		t.Fatalf(`NewNode("%v") failed: %v`, nodeId, err)
	}
	err = g.AddNode(a)
	if err != nil {
		t.Fatalf(`AddNode("%v") failed on empty graph: %v`, nodeId, err)
	}
	err = g.AddNode(a)
	if err == nil {
		t.Fatalf(`AddNode("%v") should fail: node with id "%v" already exists`, nodeId, nodeId)
	}
}

func TestGetNode(t *testing.T) {
	g := graph.NewGraph()
	nodeId := "a"
	_, err := g.GetNode(nodeId)
	if err == nil {
		t.Fatalf(`GetNode("%v") should fail: node with id "%v" does not exist in graph`, nodeId, nodeId)
	}
	a, _ := graph.NewNode(nodeId)
	_ = g.AddNode(a)
	_, err = g.GetNode(nodeId)
	if err != nil {
		t.Fatalf(`GetNode("%v") shouldn't fail: node with id "%v" exists in graph`, nodeId, nodeId)
	}
}

func TestAddEdge(t *testing.T) {
	g := graph.NewGraph()
	nodeAId, nodeBId := "a", "b"
	nodeA, _ := graph.NewNode(nodeAId)
	nodeB, _ := graph.NewNode(nodeBId)
	_ = g.AddNode(nodeA)
	_ = g.AddNode(nodeB)

	edgeAToA := graph.NewEdge(nodeA, nodeA, 1)
	err := g.AddEdge(edgeAToA)
	if err == nil {
		t.Fatalf(`AddEdge(edgeAToA) should fail: self edge not allowed`)
	}

	nodeCId := "c"
	nodeC, _ := graph.NewNode(nodeCId)
	edgeAToC := graph.NewEdge(nodeA, nodeC, 1)
	edgeCToA := graph.NewEdge(nodeC, nodeA, 1)
	err = g.AddEdge(edgeAToC)
	if err == nil {
		t.Fatalf(`AddEdge(edgeAToC) should fail: node with id "%v" does not exist`, nodeCId)
	}
	err = g.AddEdge(edgeCToA)
	if err == nil {
		t.Fatalf(`AddEdge(edgeAToC) should fail: node with id "%v" does not exist`, nodeCId)
	}
}

func TestRemoveEdge(t *testing.T) {
	g := graph.NewGraph()
	nodeA, _ := graph.NewNode("a")
	nodeB, _ := graph.NewNode("b")
	_ = g.AddNode(nodeA)
	_ = g.AddNode(nodeB)
	edge1 := graph.NewEdge(nodeA, nodeB, 1)
	edge5 := graph.NewEdge(nodeA, nodeB, 5)
	_ = g.AddEdge(edge1)
	_ = g.AddEdge(edge5)
	_ = g.AddEdge(edge5)

	g.RemoveEdgeWithWeight(nodeA, nodeB, 5)
	if len(g.GetEdges(nodeA)) != 2 {
		t.Fatalf("RemoveEdgeWithWeight(nodeA, nodeB, 5) removed multiple edges")
	}

	g.RemoveEdges(nodeA, nodeB)
	if len(g.GetEdges(nodeA)) != 0 {
		t.Fatalf("RemoveEdges(nodeA, nodeB) didn't remove all edges")
	}
}

func TestGetEdge(t *testing.T) {
	g := graph.NewGraph()
	nodeA, _ := graph.NewNode("a")
	nodeB, _ := graph.NewNode("b")
	nodeC, _ := graph.NewNode("c")
	_ = g.AddNode(nodeA)
	_ = g.AddNode(nodeB)
	_ = g.AddNode(nodeC)
	edge1 := graph.NewEdge(nodeA, nodeB, 1)
	edge5 := graph.NewEdge(nodeA, nodeB, 5)
	edge7 := graph.NewEdge(nodeA, nodeB, 7)
	_ = g.AddEdge(edge7)
	_ = g.AddEdge(edge5)
	_ = g.AddEdge(edge1)

	edge, ok := g.GetShortestEdge(nodeA, nodeB)
	if !ok {
		t.Fatalf(`GetEdge(nodeA, nodeB) didn't return an edge, even though there are`)
	}
	if edge.Weight != 1 {
		t.Fatalf("GetEdge(nodeA, nodeB) didn't return the shortest edge, expected weight %v, got %v", 1, edge.Weight)
	}

	edges := g.GetEdges(nodeA)
	if len(edges) != 3 {
		t.Fatalf("GetEdges(nodeA) should return %v edges, got %v", 3, len(edges))
	}

	edge, ok = g.GetShortestEdge(nodeA, nodeC)
	if ok {
		t.Fatalf("GetEdge(nodeA, nodeC) should return false, there are no edges between a and c")
	}
}

func TestRemoveNode(t *testing.T) {
	g := graph.NewGraph()
	nodeA, _ := graph.NewNode("a")
	nodeB, _ := graph.NewNode("b")
	nodeC, _ := graph.NewNode("c")
	_ = g.AddNode(nodeA)
	_ = g.AddNode(nodeB)
	_ = g.AddNode(nodeC)
	edgeAB := graph.NewEdge(nodeA, nodeB, 1)
	edgeBC := graph.NewEdge(nodeB, nodeC, 5)
	edgeAC := graph.NewEdge(nodeA, nodeC, 7)
	_ = g.AddEdge(edgeAC)
	_ = g.AddEdge(edgeBC)
	_ = g.AddEdge(edgeAB)

	g.RemoveNode(nodeA)
	bEdges := g.GetEdges(nodeB)
	if len(bEdges) != 1 {
		t.Fatalf("GetEdges(nodeB) should still have edges, from %v to %v", "b", "c")
	}

	_, err := g.GetNode(nodeA.Id)
	if err == nil {
		t.Fatalf(`GetNode(nodeA) should fail, node a does not exist anymore`)
	}

	nodes := g.GetAllNodes()
	if len(nodes) != 2 {
		t.Fatalf(`GetAllNodes() should return 2`)
	}
}
