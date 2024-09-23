package traverse

import (
	"graph/pkg/collections"
	"graph/pkg/graph"
	"slices"
)

type bfsState struct {
	nodes map[string]*nodeState
}

func newBfsState() bfsState {
	bs := bfsState{}
	bs.nodes = make(map[string]*nodeState)
	return bs
}

func Bfs(g graph.Graph, start, end graph.Node) (Sequence, error) {
	// setup initial values
	bs := newBfsState()
	nodes := g.GetAllNodes()
	for _, node := range nodes {
		if _, ok := bs.nodes[node.Id]; !ok {
			bs.nodes[node.Id] = &nodeState{}
		}
	}
	q := collections.NewQueue[graph.Node]()

	// add starting node
	q.Queue(start)

	var x graph.Node
	for !q.Empty() {
		// get first node on queue
		x = q.Dequeue()

		// found end node
		if x.Id == end.Id {
			break
		}

		// get all neighbours from x
		nodes := g.GetNodesFrom(x)
		for _, node := range nodes {
			if !bs.nodes[node.Id].visited {
				bs.nodes[node.Id].visited = true
				bs.nodes[node.Id].prev = x
				q.Queue(node)
			}
		}
	}

	// reconstruct the sequence from the queue
	s := NewSequence()
	s.Sequence = append(s.Sequence, x)
	for x.Id != start.Id {
		prev := bs.nodes[x.Id].prev
		s.Distance += g.GetEdge(x, prev).Weight
		x = prev
		s.Sequence = append(s.Sequence, x)
	}
	slices.Reverse(s.Sequence)
	return s, nil
}
