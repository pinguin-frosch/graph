package menus

import (
	"graph/pkg/graph"

	"github.com/pinguin-frosch/menu/pkg/menu"
)

var GraphMenu *menu.Menu
var Graph *graph.Graph

func init() {
	GraphMenu = menu.NewMenu("graph")
	GraphMenu.AddOption("s", "manage graph state", func() {
		StateMenu.Start()
	})
}
