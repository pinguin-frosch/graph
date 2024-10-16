package menus

import (
	"fmt"
	"graph/pkg/graph"

	"github.com/pinguin-frosch/menu/pkg/menu"
)

var GraphMenu *menu.Menu
var Graph graph.Graph

func init() {
	Graph = graph.NewGraph()
	GraphMenu = menu.NewMenu("graph")
	GraphMenu.AddOption("s", "manage graph state", func() {
		StateMenu.Start()
	})
	GraphMenu.AddOption("n", "add node", func() {
		id := GraphMenu.GetString("id: ")
		node, err := graph.NewNode(id)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		err = Graph.AddNode(node)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
	})
	GraphMenu.AddOption("nr", "remove node", func() {
		id := GraphMenu.GetString("id: ")
		node, err := graph.NewNode(id)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		Graph.RemoveNode(node)
	})
	GraphMenu.AddOption("e", "add edge", func() {
		fromId := GraphMenu.GetString("from: ")
		fromNode, err := Graph.GetNode(fromId)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		toId := GraphMenu.GetString("to: ")
		toNode, err := Graph.GetNode(toId)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		weight, err := GraphMenu.GetInt("weight: ")
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		edge := graph.NewEdge(fromNode, toNode, weight)
		err = Graph.AddEdge(edge)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
	})
	GraphMenu.AddOption("err", "remove edges between two nodes", func() {
		fromId := GraphMenu.GetString("from: ")
		fromNode, err := Graph.GetNode(fromId)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		toId := GraphMenu.GetString("to: ")
		toNode, err := Graph.GetNode(toId)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		Graph.RemoveEdges(fromNode, toNode)
	})
	GraphMenu.AddOption("erw", "remove edge with weight", func() {
		fromId := GraphMenu.GetString("from: ")
		fromNode, err := Graph.GetNode(fromId)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		toId := GraphMenu.GetString("to: ")
		toNode, err := Graph.GetNode(toId)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		weight, err := GraphMenu.GetInt("weight: ")
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		Graph.RemoveEdgeWithWeight(fromNode, toNode, weight)
	})
	GraphMenu.AddOption("p", "print graph", func() {
		Graph.Print()
	})
	GraphMenu.AddOption("t", "graph traverse sub menu", func() {
		TraverseMenu.Start()
	})
}
