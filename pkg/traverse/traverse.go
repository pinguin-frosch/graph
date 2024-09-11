package traverse

import (
	"errors"
	"fmt"
	"graph/pkg/graph"
	"math"
)

var (
	ErrNoTraverseAlgorithm = "no traverse algorithm has been set"
)

type Traverser interface {
	getSequence(g graph.Graph, from graph.Node) (Sequence, error)
}

type TraverseManager struct {
	traverser Traverser
}

func (tm *TraverseManager) GetSequence(g graph.Graph, from graph.Node) (Sequence, error) {
	if tm.traverser == nil {
		return Sequence{}, errors.New(ErrNoTraverseAlgorithm)
	}
	s, err := tm.traverser.getSequence(g, from)
	if err != nil {
		return s, err
	}
	return s, nil
}

func (tm *TraverseManager) GetShortestSequence(g graph.Graph) (Sequence, error) {
	if tm.traverser == nil {
		return Sequence{}, errors.New(ErrNoTraverseAlgorithm)
	}
	nodes := g.GetAllNodes()
	s := NewSequence()
	s.Distance = math.MaxInt
	for _, node := range nodes {
		sequence, err := tm.traverser.getSequence(g, node)
		if err != nil {
			return s, err
		}
		if sequence.Distance < s.Distance {
			s = sequence
		}
	}
	return s, nil
}

func (tm *TraverseManager) SetTraverseAlgorithm(a Traverser) {
	tm.traverser = a
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

func (s *Sequence) Print() {
	fmt.Printf("nodes: %d\nweight: %d\nsequence: ", len(s.Sequence), s.Distance)
	for _, node := range s.Sequence {
		fmt.Printf("%s ", node.Id)
	}
	fmt.Println()
}
