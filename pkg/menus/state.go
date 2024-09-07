package menus

import (
	"fmt"
	"graph/pkg/graph"

	"github.com/pinguin-frosch/menu/pkg/menu"
)

var StateMenu *menu.Menu

func init() {
	StateMenu = menu.NewMenu("state")
	StateMenu.AddOption("n", "create new graph", func() {
		Graph = graph.NewGraph()
	})
	StateMenu.AddOption("f", "new graph from file", func() {
		scanner := StateMenu.Scanner
		fmt.Print("message: ")
		scanner.Scan()
		g, err := graph.NewGraphFromFile(scanner.Text())
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		Graph = g
	})
	StateMenu.AddOption("s", "save graph state", func() {
		err := Graph.SaveGraphToFile()
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
	})
}
