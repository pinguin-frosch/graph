package traverse

import "graph/pkg/graph"

type Traverser interface {
	getSequence(g graph.Graph, from graph.Node) (Sequence, error)
}

type Sequence struct {
	Distance int
	Sequence []graph.Node
}

func NewSequence() Sequence {
	s := Sequence{}
	s.Sequence = make([]graph.Node, 0)
	return s
}
