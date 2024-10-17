package traverse

import (
	"errors"
	"fmt"
	"graph/pkg/collections"
	"graph/pkg/graph"
	"math"
	"slices"
)

var (
	ErrGraphNotEulerian = "graph is not eulerian"
)

type eulerState struct {
	visitedEdges map[string]bool
	invalidEdges map[string]bool
}

func newEulerState() eulerState {
	es := eulerState{}
	es.visitedEdges = make(map[string]bool)
	es.invalidEdges = make(map[string]bool)
	return es
}

func (es *eulerState) allEdgesVisited() bool {
	for _, v := range es.visitedEdges {
		if !v {
			return false
		}
	}
	return true
}

func isEulerianGraph(g graph.Graph) bool {
	nodes := g.GetAllNodes()
	for _, node := range nodes {
		edges := g.GetEdges(node)
		if len(edges) == 0 || len(edges)%2 == 1 {
			return false
		}
	}
	return true
}

type Pair struct {
	L graph.Node
	R graph.Node
}

func generatePairs(nodes []graph.Node) []Pair {
	pairs := make([]Pair, 0)
	for i := 0; i < len(nodes)-1; i++ {
		for j := i + 1; j < len(nodes); j++ {
			pair := Pair{L: nodes[i], R: nodes[j]}
			pairs = append(pairs, pair)
		}
	}
	return pairs
}

func (p Pair) Different(other Pair) bool {
	return p.L.Id != other.L.Id && p.L.Id != other.R.Id && p.R.Id != other.L.Id && p.R.Id != other.R.Id
}

// sort the nodes in the pair so they're consistent
func (p *Pair) Normalize() {
	if p.L.Id > p.R.Id {
		p.L, p.R = p.R, p.L
	}
}

// normalize the pairing, sort each pair and then the slice
func normalizePairing(pairing []Pair) {
	for i := range pairing {
		pairing[i].Normalize() // normalize each pair
	}
	// sort the pairing
	slices.SortFunc(pairing, func(a, b Pair) int {
		if a.Key() < b.Key() {
			return -1
		} else {
			return 1
		}
	})
}

func (p Pair) Key() string {
	return fmt.Sprintf("%s|%s", p.L.Id, p.R.Id)
}

func generatePairing(pair Pair, pairs []Pair) []Pair {
	pairing := make([]Pair, 0)

	// add pair to the pairing
	pairing = append(pairing, pair)

	// remove pairs that contains elements from the current pair
	filteredPairs := make([]Pair, 0)
	for _, p := range pairs {
		if pair.Different(p) {
			filteredPairs = append(filteredPairs, p)
		}
	}

	// repeat for the filtered pairs
	for _, p := range filteredPairs {
		subPairing := generatePairing(p, filteredPairs)
		pairing = append(pairing, subPairing...)
	}

	return pairing
}

func generateAllPairings(pairs []Pair) [][]Pair {
	// store all the pairings
	pairings := make([][]Pair, 0)

	// add pairing for each pair
	for _, pair := range pairs {
		pairing := generatePairing(pair, pairs)
		pairings = append(pairings, pairing)
	}

	return pairings
}

func removeDuplicatePairings(pairings [][]Pair) [][]Pair {
	uniquePairings := make([][]Pair, 0)
	seen := make(map[string]bool)

	for _, pairing := range pairings {
		normalizePairing(pairing)
		key := fmt.Sprint(pairing)
		if !seen[key] {
			seen[key] = true
			uniquePairings = append(uniquePairings, pairing)
		}
	}

	return uniquePairings
}

func getBestPairing(g graph.Graph, pairings [][]Pair) ([]Pair, error) {
	if len(pairings) == 0 {
		return []Pair{}, errors.New("no pairings to compare")
	}

	// keep track of the best pairing so far
	bestPairingId := 0
	bestPairingWeight := math.MaxInt

	for pairingId, pairing := range pairings {
		// calculate weight for the pairing
		pairingWeight := 0
		for _, pair := range pairing {
			s, err := Dijkstra(g, pair.L, pair.R)
			if err != nil {
				return []Pair{}, err
			}
			pairingWeight += s.Distance
		}

		// update best pairing if necessary
		if pairingWeight < bestPairingWeight {
			bestPairingWeight = pairingWeight
			bestPairingId = pairingId
		}
	}

	return pairings[bestPairingId], nil
}

func duplicateEdges(g *graph.Graph, pairing []Pair) error {
	for _, pair := range pairing {
		// get all nodes to connect the pair
		sequence, err := Dijkstra(*g, pair.L, pair.R)
		if err != nil {
			return err
		}

		// add edges to connect the pair again
		for i := 1; i < len(sequence.Sequence); i++ {
			a := sequence.Sequence[i-1]
			b := sequence.Sequence[i]
			edge, ok := g.GetShortestEdge(a, b)
			if !ok {
				return fmt.Errorf("couldn't get shortest edge between %v and %v", a.Id, b.Id)
			}
			newEdge := graph.NewEdge(a, b, edge.Weight)
			err = g.AddEdge(newEdge)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func eulerizeGraph(g *graph.Graph) error {
	// duplicate edges of deadend nodes
	deadendNodes := g.GetAllDeadendNodes()
	for _, deadendNode := range deadendNodes {
		g.AddEdge(g.GetEdges(deadendNode)[0])
	}

	// build pairs of the odd nodes
	oddNodes := g.GetAllOddNodes()
	pairs := generatePairs(oddNodes)

	// create the different ways to join the odd nodes without repetition
	allPairings := generateAllPairings(pairs)

	// remove duplicate pairings
	uniquePairings := removeDuplicatePairings(allPairings)

	// get best pairing
	bestPairing, err := getBestPairing(*g, uniquePairings)
	if err != nil {
		return err
	}

	// duplicate necessary edges
	err = duplicateEdges(g, bestPairing)
	if err != nil {
		return err
	}

	g.Print()

	return nil
}

func Euler(g graph.Graph, a graph.Node) (Sequence, error) {
	// clone to avoid modifying the original graph
	h := g.Clone()

	// check that starting node exists
	_, err := h.GetNode(a.Id)
	if err != nil {
		return Sequence{}, err
	}

	// check if graph is Eulerian
	if !isEulerianGraph(h) {
		err := eulerizeGraph(&h)
		if err != nil {
			return Sequence{}, err
		}
	}

	// setup initial values
	es := newEulerState()
	edges := h.GetAllEdges()
	for _, edge := range edges {
		es.visitedEdges[edge.Key()] = false
	}
	st := collections.NewStack[graph.Node]()

	// add starting node
	st.Push(a)

	x := a
	for !es.allEdgesVisited() {
		// get valid edges to go to the next node
		edges := h.GetEdges(x)
		validEdges := make([]graph.Edge, 0, len(edges))
		for _, edge := range edges {
			if !es.visitedEdges[edge.Key()] {
				key := fmt.Sprintf("%d|%s|%s", st.Len(), x, edge.To)
				if !es.invalidEdges[key] {
					validEdges = append(validEdges, edge)
				}
			}
		}

		// there's nowhere to go, we need to go back
		if len(validEdges) == 0 {
			// add restriction
			invalidNode, _ := st.Pop()
			previousNode, _ := st.Peek()
			key := fmt.Sprintf("%d|%s|%s", st.Len(), previousNode, invalidNode)
			es.invalidEdges[key] = true

			// mark the edge as not visited again
			edge, ok := h.GetShortestEdge(previousNode, invalidNode)
			if !ok {
				return Sequence{}, fmt.Errorf("couldn't get shortest edge between %v and %v", previousNode.Id, invalidNode.Id)
			}
			es.visitedEdges[edge.Key()] = false
			es.visitedEdges[edge.ReversedEdge().Key()] = false

			// go back
			x = previousNode
			continue
		}

		// add next node
		nextEdge := validEdges[0]
		es.visitedEdges[nextEdge.Key()] = true
		es.visitedEdges[nextEdge.ReversedEdge().Key()] = true
		nextNode := nextEdge.To
		st.Push(nextNode)
		x = nextNode
	}

	// reconstruct the sequence from the stack
	s := NewSequence()
	for !st.Empty() {
		lastNode, _ := st.Pop()
		s.Sequence = append(s.Sequence, lastNode)
	}
	for i := 0; i < len(s.Sequence)-1; i++ {
		a := s.Sequence[i]
		b := s.Sequence[i+1]
		edge, ok := h.GetShortestEdge(a, b)
		if !ok {
			return Sequence{}, fmt.Errorf("couldn't get shortest edge between %v and %v", a.Id, b.Id)
		}
		s.Distance += edge.Weight
	}

	return s, nil
}
