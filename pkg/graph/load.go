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
	Edges    []string `json:"edges"`
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
		if len(parts) != 2 {
			return nil, errors.New(fmt.Sprintf("Invalid edge: %v", e))
		}
		from := parts[0]
		to := parts[1]
		err := g.AddEdge(from, to)
		if err != nil {
			return nil, err
		}
	}
	return &g, nil
}
