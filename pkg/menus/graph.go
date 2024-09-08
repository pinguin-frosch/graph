package menus

import (
	"fmt"
	"graph/pkg/graph"
	"strconv"
	"strings"

	"github.com/pinguin-frosch/menu/pkg/menu"
)

var GraphMenu *menu.Menu
var Graph *graph.Graph

func init() {
	Graph = graph.NewGraph()
	GraphMenu = menu.NewMenu("graph")
	GraphMenu.AddOption("s", "manage graph state", func() {
		StateMenu.Start()
	})
	GraphMenu.AddOption("n", "add node", func() {
		scanner := GraphMenu.Scanner
		fmt.Printf("id: ")
		scanner.Scan()
		id := strings.Trim(scanner.Text(), " ")
		node := graph.NewNode(id)
		err := Graph.AddNode(node)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
	})
	GraphMenu.AddOption("e", "add edge", func() {
		scanner := GraphMenu.Scanner
		fmt.Printf("from: ")
		scanner.Scan()
		fromId := strings.Trim(scanner.Text(), " ")
		fromNode := graph.NewNode(fromId)
		fmt.Printf("to: ")
		scanner.Scan()
		toId := strings.Trim(scanner.Text(), " ")
		toNode := graph.NewNode(toId)
		fmt.Printf("weight: ")
		scanner.Scan()
		weight, err := strconv.Atoi(strings.Trim(scanner.Text(), " "))
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
}
