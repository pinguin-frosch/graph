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
			fmt.Print("partial sequence: ")
			s.Print()
			return
		}
		fmt.Print("sequence: ")
		s.Print()
	})
	TraverseMenu.AddOption("td", "use default traverse method", func() {
		d := traverse.NewDefault()
		traverseManager.SetTraverseAlgorithm(d)
	})
}