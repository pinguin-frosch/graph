package menus

import (
	"fmt"
	"graph/pkg/traverse"

	"github.com/pinguin-frosch/menu/pkg/menu"
)

var TraverseMenu *menu.Menu
var traverseManager traverse.TraverseManager

func init() {
	TraverseMenu = menu.NewMenu("traverse")
	TraverseMenu.AddOption("gs", "get sequence", func() {
		id := GraphMenu.GetString("from: ")
		node, err := Graph.GetNode(id)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		s, err := traverseManager.GetSequence(Graph, node)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		s.Print()
	})
	TraverseMenu.AddOption("gss", "get shortest sequence", func() {
		s, err := traverseManager.GetShortestSequence(Graph)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		s.Print()
	})
	TraverseMenu.AddOption("d", "dijkstra between two nodes", func() {
		fromId := GraphMenu.GetString("from: ")
		from, err := Graph.GetNode(fromId)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		toId := GraphMenu.GetString("to: ")
		to, err := Graph.GetNode(toId)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		s, err := traverse.Dijkstra(Graph, from, to)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		s.Print()
	})
	TraverseMenu.AddOption("e", "traverse graph using euler method", func() {
		fromId := GraphMenu.GetString("from: ")
		from, err := Graph.GetNode(fromId)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		s, err := traverse.Euler(Graph, from)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		s.Print()
	})
	TraverseMenu.AddOption("bfs", "traverse graph using bfs", func() {
		fromId := GraphMenu.GetString("from: ")
		from, err := Graph.GetNode(fromId)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		toId := GraphMenu.GetString("to: ")
		to, err := Graph.GetNode(toId)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		s, err := traverse.Bfs(Graph, from, to)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		s.Print()
	})
	TraverseMenu.AddOption("td", "use default traverse method", func() {
		d := traverse.NewDefault()
		traverseManager.SetTraverseAlgorithm(d)
	})
}
