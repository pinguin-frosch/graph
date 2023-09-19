package graph

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type GraphFile struct {
	Vertices []string `json:"vertices"`
	// [0] = From, [1] = To, [2] = DeadEnd
	Edges []string `json:"edges"`
}

func LoadFromJson(filename string) (*Graph, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var file GraphFile
	err = json.Unmarshal(bytes, &file)
	if err != nil {
		return nil, err
	}
	g := Graph{}
	for _, v := range file.Vertices {
		err := g.AddVertex(v)
		if err != nil {
			return nil, err
		}
	}
	for _, e := range file.Edges {
		parts := strings.Split(e, "|")
		if len(parts) < 2 {
			return nil, errors.New(fmt.Sprintf("Invalid edge: %v", e))
		}
		from := parts[0]
		to := parts[1]

		err := g.AddEdge(from, to)
		if err != nil {
			return nil, err
		}

		// Set the edges as dead ends
		if len(parts) >= 3 {
			edge := g.GetEdge(from, to)
			if edge == nil {
				return nil, errors.New(fmt.Sprintf("Couldn't get edge %v<->%v\n", from, to))
			}
			edge.deadEnd = true
			otherEdge := g.GetEdge(to, from)
			if otherEdge == nil {
				return nil, errors.New(fmt.Sprintf("Couldn't get edge %v<->%v\n", to, from))
			}
			otherEdge.deadEnd = true
		}
	}
	return &g, nil
}
