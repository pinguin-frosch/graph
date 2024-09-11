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
		node := graph.NewNode(id)
		err := Graph.AddNode(node)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
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
	GraphMenu.AddOption("p", "print graph", func() {
		Graph.Print()
	})
	GraphMenu.AddOption("t", "graph traverse sub menu", func() {
		TraverseMenu.Start()
	})
}
